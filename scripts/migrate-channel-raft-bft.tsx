import { checkbox, confirm, input, select } from '@inquirer/prompts'
import * as k8s from '@kubernetes/client-node'
import { readFile } from 'fs/promises'
const kc = new k8s.KubeConfig()
kc.loadFromDefault()

const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

async function updateOrdererTag(ordererNames: string[], namespace: string = 'default') {
	for (const ordererName of ordererNames) {
		try {
			// Get the current FabricOrdererNode
			const res = await k8sApi.getNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes', ordererName)

			const orderer = res.body as any
			console.log(orderer)
			// Update the tag to 3.0.0-beta
			if (orderer.spec && orderer.spec.image) {
				orderer.spec.tag = '3.0.0-beta'
			} else {
				console.error(`Unable to update tag for orderer ${ordererName}: image spec not found`)
				continue
			}

			// Update the FabricOrdererNode
			await k8sApi.patchNamespacedCustomObject('hlf.kungfusoftware.es', 'v1alpha1', namespace, 'fabricorderernodes', ordererName, orderer, undefined, undefined, undefined, {
				headers: { 'Content-Type': 'application/merge-patch+json' },
			})

			console.log(`Successfully updated tag for orderer ${ordererName} to 3.0.0-beta`)
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
		return ordererList.filter((orderer: any) => orderer.spec.image.tag !== '3.0.0-beta')
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
		const pollInterval = 10000 // 10 seconds

		const startTime = Date.now()

		while (!ready && Date.now() - startTime < maxWaitTime) {
			try {
				const appsV1Api = kc.makeApiClient(k8s.AppsV1Api)
				const res = await appsV1Api.readNamespacedDeployment(orderer.name, orderer.namespace)
				const deployment = res.body

				const hasCorrectTag = deployment.spec?.template.spec?.containers.some((container) => container.image?.includes('3.0.0-beta'))
				const isReady =
					deployment.status?.conditions?.some((condition) => condition.type === 'Available' && condition.status === 'True') &&
					deployment.status?.readyReplicas === deployment.status?.replicas

				if (hasCorrectTag && isReady) {
					ready = true
					console.log(`Orderer ${orderer.name} in namespace ${orderer.namespace} is ready with tag 3.0.0-beta`)
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

async function updateChannelToBFT(channelName: string): Promise<void> {
	try {
		console.log(`Updating channel ${channelName} to use BFT consensus...`)
		const orderers = await getOrderersFromCluster('')
		// Fetch the current channel configuration
		const channel = await getChannelFromKubernetes(channelName)
		console.log(channel)
		// Update the consensus type to BFT
		if (channel.spec && channel.spec.channelConfig) {
			channel.spec.channelConfig.capabilities = ['V3_0']
			channel.spec.channelConfig.application.capabilities = ['V3_0']
			channel.spec.channelConfig.orderer.ordererType = 'BFT'
			// go through channel.spec.orderers and ask either for the orderer name or the namespace (radio, select one), or ask for the identity file path to get the certificate from
			const consenterMapping = []
			let idx = 0
			for (const orderer of channel.spec.orderers) {
				const selectedOrderer = await select({
					message: `Select the orderer ${orderer.name} (${orderer.namespace}) for the consenter ${orderer.host}:${orderer.port}`,
					choices: [...orderers.map((orderer) => ({ name: orderer.metadata.name, value: orderer.metadata.name })), { name: 'Identity file path', value: 'identity' }],
				})
				let identityCert = ''
				let mspId = ''
				if (selectedOrderer === 'identity') {
					const identity = await input({ message: 'Enter the identity file path:' })
					identityCert = (await readFile(identity)).toString('utf-8')
					// ask for the mspId
					mspId = await input({ message: 'Enter the mspId:' })
				} else {
					// get fabricorderernode and get the identity cert from `status.signCert`
					const fabricOrdererNode = await getFabricOrdererNode(selectedOrderer)
					identityCert = fabricOrdererNode.status.signCert
					mspId = fabricOrdererNode.spec.mspID
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
			}
			channel.spec.channelConfig.orderer.consenterMapping = consenterMapping
			channel.spec.channelConfig.orderer.smartBFT = {
				collectTimeout: '1s',
				decisionsPerLeader: 3,
				incomingMessageBufferSize: 200,
				leaderHeartbeatCount: 10,
				leaderHeartbeatTimeout: '1m0s',
				leaderRotation: 2,
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

async function main() {
	const channelName = await input({ message: 'Enter the channel name:' })
	const channel = await getChannelFromKubernetes(channelName)
	console.log(channel)
	// const ordererNamesInput = await input({ message: 'Enter orderer names (comma-separated):' })
	const ordererList = await getOrderersFromClusterBelow30('')
	const selectedOrderers = await checkbox({
		message: 'What orderers do you want to upgrade to 3.0.0-beta?',
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
		message: `Upgrade the following orderers to version 3.0.0-beta?\n${selectedOrderers.map((orderer) => `- ${orderer.name} (${orderer.namespace})`).join('\n')}`,
		default: true,
	})
	if (confirmed) {
		console.log('Upgrading the selected orderers...')
		await updateOrderers(selectedOrderers)
	}
	// confirm set channel to maintenance
	const stateConfirmed = await confirm({
		message: `Set channel ${channelName} to STATE_MAINTENANCE?`,
		default: true,
	})
	if (stateConfirmed) {
		await setFabricMainChannelToMaintenance(channelName)
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
// 2. Update the orderers to the version 3.0.0-beta one by one and wait for the orderers to be ready
// 3. Set channel to STATE_MAINTENANCE
// 4. Wait for the channel to be updated by checking the ${channel}-config configmap
// 5. Add consenter_mapping to the channel and update the capabilities
// 6. Set channel to STATE_NORMAL
// 7. Wait for the channel to be updated by checking the ${channel}-config configmap
// 8. Migration completed :)
