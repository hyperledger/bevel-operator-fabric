import { checkbox, confirm, input, select } from '@inquirer/prompts'
import * as k8s from '@kubernetes/client-node'
import { readFile } from 'fs/promises'
const kc = new k8s.KubeConfig()
kc.loadFromDefault()

const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)
const ORDERER_IMAGE_TAG = '3.0.0'
const PEER_IMAGE_TAG = '3.0.0'
async function updateOrdererTag(ordererNames: string[], namespace: string = 'default') {
	for (const ordererName of ordererNames) {
		try {
			// Get the current FabricOrdererNode
			const res = await k8sApi.getNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes', ordererName)

			const orderer = res.body as any
			// Update the tag to 3.0.0
			if (orderer.spec && orderer.spec.image) {
				orderer.spec.tag = ORDERER_IMAGE_TAG
			} else {
				console.error(`Unable to update tag for orderer ${ordererName}: image spec not found`)
				continue
			}

			// Update the FabricOrdererNode
			await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes', ordererName, orderer, undefined, undefined, undefined, {
				headers: { 'Content-Type': 'application/merge-patch+json' },
			})

			console.log(`Successfully updated tag for orderer ${ordererName} to ${ORDERER_IMAGE_TAG}`)
		} catch (err) {
			console.error(`Error updating orderer ${ordererName}:`, err)
		}
	}
}

async function getOrderersFromClusterBelow30(namespace: string): Promise<any[]> {
	const kc = new k8s.KubeConfig()
	kc.loadFromDefault()

	const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

	try {
		const res = await k8sApi.listNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes')
		const ordererList = (res.body as any).items
		return ordererList.filter((orderer: any) => orderer.spec.image.tag !== ORDERER_IMAGE_TAG)
	} catch (err) {
		console.error('Error fetching orderers from cluster:', err)
		return []
	}
}

async function getOrderersFromCluster(namespace: string): Promise<any[]> {
	const kc = new k8s.KubeConfig()
	kc.loadFromDefault()

	const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

	try {
		const res = await k8sApi.listNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes')
		const ordererList = (res.body as any).items
		return ordererList
	} catch (err) {
		console.error('Error fetching orderers from cluster:', err)
		return []
	}
}

async function updateOrderers(orderers: { name: string; namespace: string }[]) {
	for (const orderer of orderers) {
		await updateOrdererTag([orderer.name], orderer.namespace)
		console.log(`Waiting for orderer ${orderer.name} in namespace ${orderer.namespace} to be ready...`)
		// Add logic here to wait for the orderer to be ready
		// Wait for the orderer to be ready with the new tag
		let ready = false
		const maxWaitTime = 10 * 60 * 1000 // 10 minutes in milliseconds
		const pollInterval = 1000 // 1 second

		const startTime = Date.now()

		while (!ready && Date.now() - startTime < maxWaitTime) {
			try {
				const appsV1Api = kc.makeApiClient(k8s.AppsV1Api)
				const res = await appsV1Api.readNamespacedDeployment(orderer.name, orderer.namespace)
				const deployment = res.body

				const hasCorrectTag = deployment.spec?.template.spec?.containers.some((container) => container.image?.includes(ORDERER_IMAGE_TAG))
				const isReady =
					deployment.status?.conditions?.some((condition) => condition.type === 'Available' && condition.status === 'True') &&
					deployment.status?.readyReplicas === deployment.status?.replicas

				if (hasCorrectTag && isReady) {
					ready = true
					console.log(`Orderer ${orderer.name} in namespace ${orderer.namespace} is ready with tag ${ORDERER_IMAGE_TAG}`)
				} else {
					const elapsedTime = Math.floor((Date.now() - startTime) / 1000)
					console.log(`Waiting for orderer ${orderer.name} in namespace ${orderer.namespace} to be ready (${elapsedTime} seconds elapsed)...`)
					await new Promise((resolve) => setTimeout(resolve, pollInterval))
				}
			} catch (err) {
				console.error(`Error checking orderer ${orderer.name} in namespace ${orderer.namespace} status:`, err)
				await new Promise((resolve) => setTimeout(resolve, pollInterval))
			}
		}

		if (!ready) {
			console.error(`Orderer ${orderer.name} in namespace ${orderer.namespace} did not become ready within the expected time.`)
		}
	}
}

async function waitForChannelConsensusTypeBFT(channelName: string) {
	console.log(`Waiting for ${channelName} ConsensusType to be BFT...`)
	const kc = new k8s.KubeConfig()
	kc.loadFromDefault()
	const k8sApi = kc.makeApiClient(k8s.CoreV1Api)

	const maxWaitTime = 10 * 60 * 1000 // 10 minutes in milliseconds
	const pollInterval = 1000 // 1 seconds
	const startTime = Date.now()

	while (Date.now() - startTime < maxWaitTime) {
		try {
			const res = await k8sApi.readNamespacedConfigMap(`${channelName}-config`, 'default')
			const configMap = res.body
			const channelJson = JSON.parse(configMap.data!['channel.json'])
			const consensusType = channelJson.channel_group.groups.Orderer.values.ConsensusType.value.type

			if (consensusType === 'BFT') {
				console.log(`Channel ${channelName} ConsensusType is now BFT`)
				return
			}

			console.log(`Waiting for ${channelName} ConsensusType to be BFT. Current type: ${consensusType}`)
			await new Promise((resolve) => setTimeout(resolve, pollInterval))
		} catch (err) {
			console.error(`Error checking ${channelName}-config configmap:`, err)
			await new Promise((resolve) => setTimeout(resolve, pollInterval))
		}
	}

	console.error(`Timeout: ${channelName} ConsensusType did not change to BFT within 10 minutes`)
	throw new Error(`Timeout waiting for ${channelName} ConsensusType to be BFT`)
}

async function waitForChannelStateUpdate(channelName: string, expectedState: string) {
	console.log(`Waiting for ${channelName} to be updated...`)
	const kc = new k8s.KubeConfig()
	kc.loadFromDefault()
	const k8sApi = kc.makeApiClient(k8s.CoreV1Api)

	const maxWaitTime = 5 * 60 * 1000 // 5 minutes in milliseconds
	const pollInterval = 1000 // 1 second
	const startTime = Date.now()

	while (Date.now() - startTime < maxWaitTime) {
		try {
			const res = await k8sApi.readNamespacedConfigMap(`${channelName}-config`, 'default')
			const configMap = res.body
			const channelJson = JSON.parse(configMap.data!['channel.json'])
			const state = channelJson.channel_group.groups.Orderer.values.ConsensusType.value.state

			if (state === expectedState) {
				console.log(`Channel ${channelName} is now in ${expectedState}`)
				return
			}

			console.log(`Waiting for ${channelName} to be in ${expectedState}...`)
			await new Promise((resolve) => setTimeout(resolve, pollInterval))
		} catch (err) {
			console.error(`Error checking ${channelName}-config configmap:`, err)
			await new Promise((resolve) => setTimeout(resolve, pollInterval))
		}
	}

	console.error(`Timeout: ${channelName} did not enter STATE_MAINTENANCE within 5 minutes`)
	// Add logic here to check the ${channel}-config configmap
}

async function setFabricMainChannelToNormal(channelName: string, namespace: string = '') {
	try {
		console.log(`Setting ${channelName} to STATE_NORMAL...`)

		const kc = new k8s.KubeConfig()
		kc.loadFromDefault()

		const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

		// Fetch the current FabricMainChannel object
		const res = await k8sApi.getNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricmainchannels', channelName)

		const channel = res.body as any

		// Update the orderer state to STATE_NORMAL
		if (channel.spec && channel.spec.channelConfig && channel.spec.channelConfig.orderer) {
			channel.spec.channelConfig.orderer.state = 'STATE_NORMAL'
		} else {
			console.error(`Unable to update state for channel ${channelName}: channelConfig.orderer not found`)
			return
		}

		// Update the FabricMainChannel
		await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricmainchannels', channelName, channel, undefined, undefined, undefined, {
			headers: { 'Content-Type': 'application/merge-patch+json' },
		})

		// Wait for the channel to be updated
		await waitForChannelStateUpdate(channelName, 'STATE_NORMAL')
		console.log(`Successfully set ${channelName} to STATE_NORMAL`)
	} catch (err) {
		console.error(`Error setting ${channelName} to STATE_NORMAL:`, err)
		throw err
	}
}

async function setFabricMainChannelToMaintenance(channelName: string, namespace: string = '') {
	try {
		console.log(`Setting ${channelName} to STATE_MAINTENANCE...`)

		const kc = new k8s.KubeConfig()
		kc.loadFromDefault()

		const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

		// Fetch the current FabricMainChannel object
		const res = await k8sApi.getNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricmainchannels', channelName)

		const channel = res.body as any

		// Update the orderer state to STATE_MAINTENANCE
		if (channel.spec && channel.spec.channelConfig && channel.spec.channelConfig.orderer) {
			channel.spec.channelConfig.orderer.state = 'STATE_MAINTENANCE'
		} else {
			console.error(`Unable to update state for channel ${channelName}: channelConfig.orderer not found`)
			return
		}

		// Update the FabricMainChannel
		await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricmainchannels', channelName, channel, undefined, undefined, undefined, {
			headers: { 'Content-Type': 'application/merge-patch+json' },
		})
		// wait for the channel to be updated
		await waitForChannelStateUpdate(channelName, 'STATE_MAINTENANCE')
		console.log(`Successfully set ${channelName} to STATE_MAINTENANCE`)
	} catch (err) {
		console.error(`Error setting ${channelName} to STATE_MAINTENANCE:`, err)
		throw err
	}
}

async function getFabricOrdererNode(ordererName: string, namespace: string = 'default'): Promise<any> {
	try {
		console.log(`Fetching FabricOrdererNode ${ordererName} from namespace ${namespace}...`)

		const kc = new k8s.KubeConfig()
		kc.loadFromDefault()

		const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

		const res = await k8sApi.getNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes', ordererName)

		const fabricOrdererNode = res.body as any

		if (!fabricOrdererNode.status || !fabricOrdererNode.status.signCert) {
			throw new Error(`FabricOrdererNode ${ordererName} does not have a signCert in its status`)
		}

		console.log(`Successfully fetched FabricOrdererNode ${ordererName}`)
		return fabricOrdererNode
	} catch (err) {
		console.error(`Error fetching FabricOrdererNode ${ordererName}:`, err)
		throw err
	}
}

async function getChannelFromKubernetes(channelName: string): Promise<any> {
	try {
		console.log(`Fetching channel ${channelName} from Kubernetes...`)

		const kc = new k8s.KubeConfig()
		kc.loadFromDefault()

		const customApi = kc.makeApiClient(k8s.CustomObjectsApi)

		const group = 'hlf.kungfusoftware.es'
		const version = 'v1alpha1'
		const plural = 'fabricmainchannels'

		const response = await customApi.getNamespacedCustomObject(group, version, '', plural, channelName)

		console.log(`Successfully fetched channel ${channelName}`)
		return response.body
	} catch (error) {
		console.error(`Error fetching channel ${channelName}:`, error)
		throw error
	}
}

async function updateChannelCapabilities(channelName: string): Promise<void> {
	try {
		console.log(`Updating channel ${channelName} capabilities to V3_0...`)
		const channel = await getChannelFromKubernetes(channelName)

		if (channel.spec && channel.spec.channelConfig) {
			channel.spec.channelConfig.capabilities = ['V3_0']
		} else {
			console.error(`Channel ${channelName} configuration is not in the expected format.`)
			throw new Error('Invalid channel configuration')
		}

		const kc = new k8s.KubeConfig()
		kc.loadFromDefault()
		const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

		await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', '', 'fabricmainchannels', channelName, channel, undefined, undefined, undefined, {
			headers: { 'Content-Type': 'application/merge-patch+json' },
		})
		await waitForChannelCapabilitiesUpdate(channelName, ['V3_0'])
		console.log(`Successfully updated channel ${channelName} capabilities to V3_0`)
	} catch (error) {
		console.error(`Error updating channel ${channelName} capabilities:`, error)
		throw error
	}
}
async function waitForChannelCapabilitiesUpdate(channelName: string, expectedCapabilities: string[]): Promise<void> {
	console.log(`Waiting for channel ${channelName} capabilities to update...`)
	const kc = new k8s.KubeConfig()
	kc.loadFromDefault()
	const k8sApi = kc.makeApiClient(k8s.CoreV1Api)

	const maxWaitTime = 5 * 60 * 1000 // 5 minutes in milliseconds
	const pollInterval = 1000 // 1 second
	const startTime = Date.now()

	while (Date.now() - startTime < maxWaitTime) {
		try {
			const res = await k8sApi.readNamespacedConfigMap(`${channelName}-config`, 'default')
			const configMap = res.body
			const channelJson = JSON.parse(configMap.data!['channel.json'])
			const currentCapabilities = Object.keys(channelJson.channel_group.values.Capabilities.value.capabilities || {})

			if (arraysEqual(currentCapabilities, expectedCapabilities)) {
				console.log(`Channel ${channelName} capabilities have been updated successfully.`)
				return
			}

			console.log(`Waiting for ${channelName} capabilities to update. Current capabilities: ${currentCapabilities}`)
			await new Promise((resolve) => setTimeout(resolve, pollInterval))
		} catch (err) {
			console.error(`Error checking ${channelName}-config configmap:`, err)
			await new Promise((resolve) => setTimeout(resolve, pollInterval))
		}
	}

	console.error(`Timeout: ${channelName} capabilities did not update within 5 minutes`)
	throw new Error(`Timeout waiting for ${channelName} capabilities to update`)
}

function arraysEqual(arr1: string[], arr2: string[]): boolean {
	if (arr1.length !== arr2.length) return false
	return arr1.every((value, index) => value === arr2[index])
}

async function updateChannelToBFT(channelName: string): Promise<void> {
	try {
		console.log(`Updating channel ${channelName} to use BFT consensus...`)
		const orderers = await getOrderersFromCluster('')
		// Fetch the current channel configuration
		const channel = await getChannelFromKubernetes(channelName)
		console.log(channel)
		// Update the consensus type to BFT
		if (channel.spec && channel.spec.channelConfig) {
			channel.spec.channelConfig.orderer.ordererType = 'BFT'
			// go through channel.spec.orderers and ask either for the orderer name or the namespace (radio, select one), or ask for the identity file path to get the certificate from
			const consenterMapping = []
			let idx = 1
			const selectedOrderers = new Set()
			for (const orderer of channel.spec.orderers as {
				host: string
				port: number
				tlsCert: string
			}[]) {
				const availableOrderers = orderers.filter((o) => !selectedOrderers.has(o.metadata.name))
				const choices = [
					...availableOrderers.map((o) => ({
						name: `${o.metadata.name} (${o.metadata.namespace})`,
						value: `${o.metadata.name}.${o.metadata.namespace}`,
					})),
					{ name: 'Identity file path', value: 'identity' },
				]
				const selectedOrderer = await select({
					message: `Select the orderer ${orderer.host} for the consenter ${orderer.host}:${orderer.port}`,
					choices: choices,
				})
				console.log('selectedOrderer', selectedOrderer)
				let identityCert = ''
				let mspId = ''
				if (selectedOrderer === 'identity') {
					const identity = await input({ message: 'Enter the identity file path:' })
					identityCert = (await readFile(identity)).toString('utf-8')
					mspId = await input({ message: 'Enter the mspId:' })
				} else {
					const [name, namespace] = selectedOrderer.split('.')
					const fabricOrdererNode = await getFabricOrdererNode(name, namespace)
					identityCert = fabricOrdererNode.status.signCert
					mspId = fabricOrdererNode.spec.mspID
					selectedOrderers.add(selectedOrderer)
				}

				if (!identityCert) {
					throw new Error(`Identity cert not found for orderer ${selectedOrderer}`)
				}
				if (!mspId) {
					throw new Error(`MspId not found for orderer ${selectedOrderer}`)
				}

				consenterMapping.push({
					client_tls_cert: orderer.tlsCert,
					host: orderer.host,
					id: idx,
					identity: identityCert,
					msp_id: mspId,
					port: orderer.port,
					server_tls_cert: orderer.tlsCert,
				})
				idx++
			}
			channel.spec.channelConfig.orderer.consenterMapping = consenterMapping
			channel.spec.channelConfig.orderer.smartBFT = {
				collectTimeout: '1s',
				decisionsPerLeader: 3,
				incomingMessageBufferSize: 200,
				leaderHeartbeatCount: 10,
				leaderHeartbeatTimeout: '1m0s',
				leaderRotation: 0,
				requestAutoRemoveTimeout: '3m',
				requestBatchMaxBytes: 10485760,
				requestBatchMaxCount: 100,
				requestBatchMaxInterval: '50ms',
				requestComplainTimeout: '20s',
				requestForwardTimeout: '2s',
				requestMaxBytes: 10485760,
				requestPoolSize: 100000,
				speedUpViewChange: false,
				syncOnStart: true,
				viewChangeResendInterval: '5s',
				viewChangeTimeout: '20s',
			}
		} else {
			console.error(`Channel ${channelName} configuration is not in the expected format.`)
			throw new Error('Invalid channel configuration')
		}

		// Update the FabricMainChannel resource
		const kc = new k8s.KubeConfig()
		kc.loadFromDefault()
		const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

		await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', '', 'fabricmainchannels', channelName, channel, undefined, undefined, undefined, {
			headers: { 'Content-Type': 'application/merge-patch+json' },
		})
		await waitForChannelConsensusTypeBFT(channelName)
		console.log(`Successfully updated channel ${channelName} to use BFT consensus`)
	} catch (error) {
		console.error(`Error updating channel ${channelName} to BFT:`, error)
		throw error
	}
}

async function getPeersFromClusterBelow30(namespace: string): Promise<any[]> {
	const kc = new k8s.KubeConfig()
	kc.loadFromDefault()

	const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

	try {
		const res = await k8sApi.listNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricpeers')
		const peerList = (res.body as any).items
		return peerList.filter((peer: any) => peer.spec.image.tag !== PEER_IMAGE_TAG)
	} catch (err) {
		console.error('Error fetching peers from cluster:', err)
		return []
	}
}

async function updatePeerTag(peerNames: string[], namespace: string = 'default') {
	for (const peerName of peerNames) {
		try {
			const res = await k8sApi.getNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricpeers', peerName)

			const peer = res.body as any
			if (peer.spec && peer.spec.image) {
				peer.spec.tag = PEER_IMAGE_TAG
			} else {
				console.error(`Unable to update tag for peer ${peerName}: image spec not found`)
				continue
			}

			await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricpeers', peerName, peer, undefined, undefined, undefined, {
				headers: { 'Content-Type': 'application/merge-patch+json' },
			})

			console.log(`Successfully updated tag for peer ${peerName} to ${PEER_IMAGE_TAG}`)
		} catch (err) {
			console.error(`Error updating peer ${peerName}:`, err)
		}
	}
}

async function updatePeers(peers: { name: string; namespace: string }[]) {
	for (const peer of peers) {
		await updatePeerTag([peer.name], peer.namespace)
		console.log(`Waiting for peer ${peer.name} in namespace ${peer.namespace} to be ready...`)

		let ready = false
		const maxWaitTime = 10 * 60 * 1000 // 10 minutes in milliseconds
		const pollInterval = 1000 // 1 second

		const startTime = Date.now()

		while (!ready && Date.now() - startTime < maxWaitTime) {
			try {
				const appsV1Api = kc.makeApiClient(k8s.AppsV1Api)
				const res = await appsV1Api.readNamespacedDeployment(peer.name, peer.namespace)
				const deployment = res.body

				const hasCorrectTag = deployment.spec?.template.spec?.containers.some((container) => container.image?.includes(PEER_IMAGE_TAG))
				const isReady =
					deployment.status?.conditions?.some((condition) => condition.type === 'Available' && condition.status === 'True') &&
					deployment.status?.readyReplicas === deployment.status?.replicas

				if (hasCorrectTag && isReady) {
					ready = true
					console.log(`Peer ${peer.name} in namespace ${peer.namespace} is ready with tag ${PEER_IMAGE_TAG}`)
				} else {
					const elapsedTime = Math.floor((Date.now() - startTime) / 1000)
					console.log(`Waiting for peer ${peer.name} in namespace ${peer.namespace} to be ready (${elapsedTime} seconds elapsed)...`)
					await new Promise((resolve) => setTimeout(resolve, pollInterval))
				}
			} catch (err) {
				console.error(`Error checking peer ${peer.name} in namespace ${peer.namespace} status:`, err)
				await new Promise((resolve) => setTimeout(resolve, pollInterval))
			}
		}

		if (!ready) {
			console.error(`Peer ${peer.name} in namespace ${peer.namespace} did not become ready within the expected time.`)
		}
	}
}

async function main() {
	const channelName = await input({ message: 'Enter the channel name:' })
	const channel = await getChannelFromKubernetes(channelName)
	// const ordererNamesInput = await input({ message: 'Enter orderer names (comma-separated):' })
	const ordererList = await getOrderersFromClusterBelow30('')
	const selectedOrderers = await checkbox({
		message: `What orderers do you want to upgrade to ${ORDERER_IMAGE_TAG}?`,
		choices: ordererList.map((orderer: any) => ({
			name: orderer.metadata.name,
			value: {
				name: orderer.metadata.name,
				namespace: orderer.metadata.namespace,
			},
			checked: true,
		})),
	})
	// console.log('selectedOrderers', selectedOrderers)
	// ask for confirmation on to upgrade the selected orderers
	const confirmed = await confirm({
		message: `Upgrade the following orderers to version ${ORDERER_IMAGE_TAG}?\n${selectedOrderers.map((orderer) => `- ${orderer.name} (${orderer.namespace})`).join('\n')}`,
		default: true,
	})
	if (confirmed) {
		console.log('Upgrading the selected orderers...')
		await updateOrderers(selectedOrderers)
	}

	// Add peer upgrade step
	const peerList = await getPeersFromClusterBelow30('')
	const selectedPeers = await checkbox({
		message: `What peers do you want to upgrade to ${PEER_IMAGE_TAG}?`,
		choices: peerList.map((peer: any) => ({
			name: peer.metadata.name,
			value: {
				name: peer.metadata.name,
				namespace: peer.metadata.namespace,
			},
			checked: true,
		})),
	})

	const peerConfirmed = await confirm({
		message: `Upgrade the following peers to version ${PEER_IMAGE_TAG}?\n${selectedPeers.map((peer) => `- ${peer.name} (${peer.namespace})`).join('\n')}`,
		default: true,
	})

	if (peerConfirmed) {
		console.log('Upgrading the selected peers...')
		await updatePeers(selectedPeers)
	}
	// confirm set channel to maintenance
	const stateConfirmed = await confirm({
		message: `Set channel ${channelName} to STATE_MAINTENANCE?`,
		default: true,
	})
	if (stateConfirmed) {
		await setFabricMainChannelToMaintenance(channelName)
	}

	const capabilitiesConfirmed = await confirm({
		message: `Update channel ${channelName} capabilities to V3_0?`,
		default: true,
	})
	if (capabilitiesConfirmed) {
		await updateChannelCapabilities(channelName)
	}

	const bftConfirmed = await confirm({
		message: `Update channel ${channelName} to use BFT consensus?`,
		default: true,
	})
	if (bftConfirmed) {
		await updateChannelToBFT(channelName)
	}

	const stateNormalConfirmed = await confirm({
		message: `Set channel ${channelName} to STATE_NORMAL?`,
		default: true,
	})
	if (stateNormalConfirmed) {
		await setFabricMainChannelToNormal(channelName)
	}
}

main().catch(console.error)

// 1. Ask for backup of the orderers
// 2. Update the orderers to the version 3.0.0 one by one and wait for the orderers to be ready
// 3. Set channel to STATE_MAINTENANCE
// 4. Wait for the channel to be updated by checking the ${channel}-config configmap
// 5. Add consenter_mapping to the channel and update the capabilities
// 6. Set channel to STATE_NORMAL
// 7. Wait for the channel to be updated by checking the ${channel}-config configmap
// 8. Migration completed :)
