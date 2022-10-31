<p>Packages:</p>
<ul>
<li>
<a href="#hlf.kungfusoftware.es%2fv1alpha1">hlf.kungfusoftware.es/v1alpha1</a>
</li>
</ul>
<h2 id="hlf.kungfusoftware.es/v1alpha1">hlf.kungfusoftware.es/v1alpha1</h2>
Resource Types:
<ul><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCA">FabricCA</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricChaincode">FabricChaincode</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricExplorer">FabricExplorer</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannel">FabricFollowerChannel</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannel">FabricMainChannel</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfig">FabricNetworkConfig</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsole">FabricOperationsConsole</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPI">FabricOperatorAPI</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUI">FabricOperatorUI</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNode">FabricOrdererNode</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingService">FabricOrderingService</a>
</li><li>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeer">FabricPeer</a>
</li></ul>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCA">FabricCA
</h3>
<p>
<p>FabricCA is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricCA</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">
FabricCASpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#nodeselector-v1-core">
Kubernetes core/v1.NodeSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>serviceMonitor</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">
ServiceMonitor
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>db</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCADatabase">
FabricCADatabase
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
<p>Hosts for the Fabric CA</p>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpecService">
FabricCASpecService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>version</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>debug</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>clrSizeLimit</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>rootCA</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCATLSConf">
FabricCATLSConf
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ca</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">
FabricCAItemConf
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tlsCA</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">
FabricCAItemConf
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cors</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Cors">
Cors
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>metrics</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAMetrics">
FabricCAMetrics
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricChaincode">FabricChaincode
</h3>
<p>
<p>FabricChaincode is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricChaincode</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricChaincodeSpec">
FabricChaincodeSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>packageId</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>credentials</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.TLS">
TLS
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricExplorer">FabricExplorer
</h3>
<p>
<p>FabricExplorer is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricExplorer</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricExplorerSpec">
FabricExplorerSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannel">FabricFollowerChannel
</h3>
<p>
<p>FabricFollowerChannel is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricFollowerChannel</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">
FabricFollowerChannelSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
<p>Name of the channel</p>
</td>
</tr>
<tr>
<td>
<code>mspId</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization to join the channel</p>
</td>
</tr>
<tr>
<td>
<code>orderers</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelOrderer">
[]FabricFollowerChannelOrderer
</a>
</em>
</td>
<td>
<p>Orderers to fetch the configuration block from</p>
</td>
</tr>
<tr>
<td>
<code>peersToJoin</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelPeer">
[]FabricFollowerChannelPeer
</a>
</em>
</td>
<td>
<p>Peers to join the channel</p>
</td>
</tr>
<tr>
<td>
<code>externalPeersToJoin</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelExternalPeer">
[]FabricFollowerChannelExternalPeer
</a>
</em>
</td>
<td>
<p>Peers to join the channel</p>
</td>
</tr>
<tr>
<td>
<code>anchorPeers</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelAnchorPeer">
[]FabricFollowerChannelAnchorPeer
</a>
</em>
</td>
<td>
<p>Anchor peers defined for the current organization</p>
</td>
</tr>
<tr>
<td>
<code>hlfIdentity</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.HLFIdentity">
HLFIdentity
</a>
</em>
</td>
<td>
<p>Identity to use to interact with the peers and the orderers</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannel">FabricMainChannel
</h3>
<p>
<p>FabricMainChannel is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricMainChannel</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">
FabricMainChannelSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
<p>Name of the channel</p>
</td>
</tr>
<tr>
<td>
<code>identities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelIdentity">
map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelIdentity
</a>
</em>
</td>
<td>
<p>HLF Identities to be used to create and manage the channel</p>
</td>
</tr>
<tr>
<td>
<code>adminPeerOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAdminPeerOrganizationSpec">
[]FabricMainChannelAdminPeerOrganizationSpec
</a>
</em>
</td>
<td>
<p>Organizations that manage the <code>application</code> configuration of the channel</p>
</td>
</tr>
<tr>
<td>
<code>peerOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPeerOrganization">
[]FabricMainChannelPeerOrganization
</a>
</em>
</td>
<td>
<p>Peer organizations that are external to the Kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>externalPeerOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalPeerOrganization">
[]FabricMainChannelExternalPeerOrganization
</a>
</em>
</td>
<td>
<p>External peer organizations that are inside the kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>channelConfig</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConfig">
FabricMainChannelConfig
</a>
</em>
</td>
<td>
<p>Configuration about the channel</p>
</td>
</tr>
<tr>
<td>
<code>adminOrdererOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec">
[]FabricMainChannelAdminOrdererOrganizationSpec
</a>
</em>
</td>
<td>
<p>Organizations that manage the <code>orderer</code> configuration of the channel</p>
</td>
</tr>
<tr>
<td>
<code>ordererOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererOrganization">
[]FabricMainChannelOrdererOrganization
</a>
</em>
</td>
<td>
<p>External orderer organizations that are inside the kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>externalOrdererOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalOrdererOrganization">
[]FabricMainChannelExternalOrdererOrganization
</a>
</em>
</td>
<td>
<p>Orderer organizations that are external to the Kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>orderers</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConsenter">
[]FabricMainChannelConsenter
</a>
</em>
</td>
<td>
<p>Consenters are the orderer nodes that are part of the channel consensus</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfig">FabricNetworkConfig
</h3>
<p>
<p>FabricNetworkConfig is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricNetworkConfig</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfigSpec">
FabricNetworkConfigSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>organization</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>internal</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>organizations</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secretName</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsole">FabricOperationsConsole
</h3>
<p>
<p>FabricOperationsConsole is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricOperationsConsole</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleSpec">
FabricOperationsConsoleSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>auth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleAuth">
FabricOperationsConsoleAuth
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>couchDB</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleCouchDB">
FabricOperationsConsoleCouchDB
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>config</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>ingress</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">
Ingress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hostUrl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPI">FabricOperatorAPI
</h3>
<p>
<p>FabricOperatorAPI is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricOperatorAPI</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPISpec">
FabricOperatorAPISpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ingress</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">
Ingress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>auth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIAuth">
FabricOperatorAPIAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hlfConfig</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIHLFConfig">
FabricOperatorAPIHLFConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorUI">FabricOperatorUI
</h3>
<p>
<p>FabricOperatorUI is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricOperatorUI</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUISpec">
FabricOperatorUISpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>logoUrl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>auth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUIAuth">
FabricOperatorUIAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ingress</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">
Ingress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>apiUrl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOrdererNode">FabricOrdererNode
</h3>
<p>
<p>FabricOrdererNode is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricOrdererNode</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">
FabricOrdererNodeSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>grpcProxy</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.GRPCProxy">
GRPCProxy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>updateCertificateTime</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>serviceMonitor</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">
ServiceMonitor
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hostAliases</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#hostalias-v1-core">
[]Kubernetes core/v1.HostAlias
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#nodeselector-v1-core">
Kubernetes core/v1.NodeSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>genesis</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>bootstrapMethod</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.BootstrapMethod">
BootstrapMethod
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>channelParticipationEnabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNodeService">
OrdererNodeService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secret</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Secret">
Secret
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>adminIstio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOrderingService">FabricOrderingService
</h3>
<p>
<p>FabricOrderingService is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricOrderingService</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">
FabricOrderingServiceSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>enrollment</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererEnrollment">
OrdererEnrollment
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodes</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNode">
[]OrdererNode
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererService">
OrdererService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>systemChannel</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererSystemChannel">
OrdererSystemChannel
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeer">FabricPeer
</h3>
<p>
<p>FabricPeer is the Schema for the hlfs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code>
string</td>
<td>
<code>
hlf.kungfusoftware.es/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code>
string
</td>
<td><code>FabricPeer</code></td>
</tr>
<tr>
<td>
<code>metadata</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">
FabricPeerSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>updateCertificateTime</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>serviceMonitor</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">
ServiceMonitor
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hostAliases</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#hostalias-v1-core">
[]Kubernetes core/v1.HostAlias
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#nodeselector-v1-core">
Kubernetes core/v1.NodeSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>couchDBexporter</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchdbExporter">
FabricPeerCouchdbExporter
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>grpcProxy</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.GRPCProxy">
GRPCProxy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>dockerSocketPath</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>externalBuilders</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ExternalBuilder">
[]ExternalBuilder
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>gossip</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpecGossip">
FabricPeerSpecGossip
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>externalEndpoint</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>external_chaincode_builder</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>couchdb</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchDB">
FabricPeerCouchDB
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>fsServer</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFSServer">
FabricFSServer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secret</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Secret">
Secret
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.PeerService">
PeerService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>stateDb</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.StateDB">
StateDB
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerStorage">
FabricPeerStorage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>discovery</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerDiscovery">
FabricPeerDiscovery
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>logging</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerLogging">
FabricPeerLogging
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerResources">
FabricPeerResources
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ApplicationCapabilities">ApplicationCapabilities
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ChannelConfig">ChannelConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>V2_0</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.BootstrapMethod">BootstrapMethod
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>)
</p>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.CA">CA
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cert</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>user</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>password</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.CARef">CARef
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>caName</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caNamespace</code>
<em>
string
</em>
</td>
<td>
<p>FabricCA Namespace of the organization</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Catls">Catls
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Component">Component</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.TLS">TLS</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cacert</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ChannelCapabilities">ChannelCapabilities
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ChannelConfig">ChannelConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>V2_0</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ChannelConfig">ChannelConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererSystemChannel">OrdererSystemChannel</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>batchTimeout</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxMessageCount</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>absoluteMaxBytes</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>preferredMaxBytes</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ordererCapabilities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererCapabilities">
OrdererCapabilities
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>applicationCapabilities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ApplicationCapabilities">
ApplicationCapabilities
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>channelCapabilities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ChannelCapabilities">
ChannelCapabilities
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>snapshotIntervalSize</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tickInterval</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>electionTick</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>heartbeatTick</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxInflightBlocks</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Component">Component
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Enrollment">Enrollment</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererEnrollment">OrdererEnrollment</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cahost</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caname</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caport</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>catls</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Catls">
Catls
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>enrollid</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>enrollsecret</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Condition">Condition
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ConditionType">
ConditionType
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>reason</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ConditionReason">
ConditionReason
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>lastTransitionTime</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ConditionReason">ConditionReason
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Condition">Condition</a>)
</p>
<p>
<p>ConditionReason is intended to be a one-word, CamelCase representation of
the category of cause of the current status. It is intended to be used in
concise output, such as one-line kubectl get output, and in summarizing
occurrences of causes.</p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ConditionType">ConditionType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Condition">Condition</a>)
</p>
<p>
<p>ConditionType is the type of the condition and is typically a CamelCased
word or short phrase.</p>
<p>Condition types should indicate state in the &ldquo;abnormal-true&rdquo; polarity. For
example, if the condition indicates when a policy is invalid, the &ldquo;is valid&rdquo;
case is probably the norm, so the condition should be called &ldquo;Invalid&rdquo;.</p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Cors">Cors
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>origins</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Csr">Csr
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNodeEnrollmentTLS">OrdererNodeEnrollmentTLS</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.TLS">TLS</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>cn</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.DeploymentStatus">DeploymentStatus
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAStatus">FabricCAStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricChaincodeStatus">FabricChaincodeStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricExplorerStatus">FabricExplorerStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelStatus">FabricFollowerChannelStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelStatus">FabricMainChannelStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfigStatus">FabricNetworkConfigStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleStatus">FabricOperationsConsoleStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIStatus">FabricOperatorAPIStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUIStatus">FabricOperatorUIStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeStatus">FabricOrdererNodeStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceStatus">FabricOrderingServiceStatus</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerStatus">FabricPeerStatus</a>)
</p>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Enrollment">Enrollment
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Secret">Secret</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>component</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Component">
Component
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tls</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.TLS">
TLS
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ExternalBuilder">ExternalBuilder
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>path</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>propagateEnvironment</code>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAAffiliation">FabricCAAffiliation
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>departments</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCABCCSP">FabricCABCCSP
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>sw</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCABCCSPSW">
FabricCABCCSPSW
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCABCCSPSW">FabricCABCCSPSW
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCABCCSP">FabricCABCCSP</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>hash</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>security</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACFG">FabricCACFG
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>identities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACFGIdentities">
FabricCACFGIdentities
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>affiliations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACFGAffilitions">
FabricCACFGAffilitions
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACFGAffilitions">FabricCACFGAffilitions
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACFG">FabricCACFG</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>allowRemove</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACFGIdentities">FabricCACFGIdentities
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACFG">FabricCACFG</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>allowRemove</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACRL">FabricCACRL
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>expiry</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACSR">FabricCACSR
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cn</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>names</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCANames">
[]FabricCANames
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ca</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACSRCA">
FabricCACSRCA
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACSRCA">FabricCACSRCA
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACSR">FabricCACSR</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>expiry</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pathLength</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAClientAuth">FabricCAClientAuth
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricTLSCACrypto">FabricTLSCACrypto</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
string
</em>
</td>
<td>
<p>NoClientCert, RequestClientCert, RequireAnyClientCert, VerifyClientCertIfGiven and RequireAndVerifyClientCert.</p>
</td>
</tr>
<tr>
<td>
<code>cert_file</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCACrypto">FabricCACrypto
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>key</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cert</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>chain</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCADatabase">FabricCADatabase
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>datasource</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIdentity">FabricCAIdentity
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCARegistry">FabricCARegistry</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pass</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>type</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>affiliation</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>attrs</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIdentityAttrs">
FabricCAIdentityAttrs
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIdentityAttrs">FabricCAIdentityAttrs
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIdentity">FabricCAIdentity</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>hf.Registrar.Roles</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hf.Registrar.DelegateRoles</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hf.Registrar.Attributes</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hf.Revoker</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hf.IntermediateCA</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hf.GenCRL</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hf.AffiliationMgr</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediate">FabricCAIntermediate
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>parentServer</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateParentServer">
FabricCAIntermediateParentServer
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateEnrollment">FabricCAIntermediateEnrollment
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>hosts</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>profile</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>label</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateParentServer">FabricCAIntermediateParentServer
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediate">FabricCAIntermediate</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>url</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caName</code>
<em>
string
</em>
</td>
<td>
<p>FabricCA Name of the organization</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateTLS">FabricCAIntermediateTLS
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>certFiles</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>client</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateTLSClient">
FabricCAIntermediateTLSClient
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateTLSClient">FabricCAIntermediateTLSClient
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediateTLS">FabricCAIntermediateTLS</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>certFile</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>keyFile</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cfg</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACFG">
FabricCACFG
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>subject</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASubject">
FabricCASubject
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>csr</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACSR">
FabricCACSR
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>signing</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigning">
FabricCASigning
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>crl</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACRL">
FabricCACRL
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>registry</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCARegistry">
FabricCARegistry
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>intermediate</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIntermediate">
FabricCAIntermediate
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>bccsp</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCABCCSP">
FabricCABCCSP
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>affiliations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAAffiliation">
[]FabricCAAffiliation
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>ca</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACrypto">
FabricCACrypto
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tlsCa</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricTLSCACrypto">
FabricTLSCACrypto
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAMetrics">FabricCAMetrics
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provider</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>statsd</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAMetricsStatsd">
FabricCAMetricsStatsd
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAMetricsStatsd">FabricCAMetricsStatsd
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAMetrics">FabricCAMetrics</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>network</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>address</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>writeInterval</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>prefix</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCANames">FabricCANames
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCACSR">FabricCACSR</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>C</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ST</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>O</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>L</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>OU</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCARegistry">FabricCARegistry
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>max_enrollments</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>identities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAIdentity">
[]FabricCAIdentity
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASigning">FabricCASigning
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningDefault">
FabricCASigningDefault
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>profiles</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningProfiles">
FabricCASigningProfiles
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASigningDefault">FabricCASigningDefault
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigning">FabricCASigning</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>expiry</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>usage</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASigningProfiles">FabricCASigningProfiles
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigning">FabricCASigning</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ca</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningSignProfile">
FabricCASigningSignProfile
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tls</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningTLSProfile">
FabricCASigningTLSProfile
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASigningSignProfile">FabricCASigningSignProfile
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningProfiles">FabricCASigningProfiles</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>usage</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>expiry</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caconstraint</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningSignProfileConstraint">
FabricCASigningSignProfileConstraint
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASigningSignProfileConstraint">FabricCASigningSignProfileConstraint
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningSignProfile">FabricCASigningSignProfile</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>isCA</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxPathLen</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASigningTLSProfile">FabricCASigningTLSProfile
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASigningProfiles">FabricCASigningProfiles</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>usage</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>expiry</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCA">FabricCA</a>)
</p>
<p>
<p>FabricCASpec defines the desired state of FabricCA</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#nodeselector-v1-core">
Kubernetes core/v1.NodeSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>serviceMonitor</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">
ServiceMonitor
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>db</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCADatabase">
FabricCADatabase
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
<p>Hosts for the Fabric CA</p>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpecService">
FabricCASpecService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>version</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>debug</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>clrSizeLimit</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>rootCA</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCATLSConf">
FabricCATLSConf
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ca</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">
FabricCAItemConf
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tlsCA</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">
FabricCAItemConf
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cors</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Cors">
Cors
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>metrics</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAMetrics">
FabricCAMetrics
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASpecService">FabricCASpecService
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#servicetype-v1-core">
Kubernetes core/v1.ServiceType
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCAStatus">FabricCAStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCA">FabricCA</a>)
</p>
<p>
<p>FabricCAStatus defines the observed state of FabricCA</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodePort</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tls_cert</code>
<em>
string
</em>
</td>
<td>
<p>TLS Certificate to connect to the FabricCA</p>
</td>
</tr>
<tr>
<td>
<code>ca_cert</code>
<em>
string
</em>
</td>
<td>
<p>Root certificate for Sign certificates generated by FabricCA</p>
</td>
</tr>
<tr>
<td>
<code>tlsca_cert</code>
<em>
string
</em>
</td>
<td>
<p>Root certificate for TLS certificates generated by FabricCA</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCASubject">FabricCASubject
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCATLSConf">FabricCATLSConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cn</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>C</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ST</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>O</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>L</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>OU</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricCATLSConf">FabricCATLSConf
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>subject</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASubject">
FabricCASubject
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricChaincodeSpec">FabricChaincodeSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricChaincode">FabricChaincode</a>)
</p>
<p>
<p>FabricChaincodeSpec defines the desired state of FabricChaincode</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>packageId</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>credentials</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.TLS">
TLS
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricChaincodeStatus">FabricChaincodeStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricChaincode">FabricChaincode</a>)
</p>
<p>
<p>FabricChaincodeStatus defines the observed state of FabricChaincode</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricExplorerSpec">FabricExplorerSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricExplorer">FabricExplorer</a>)
</p>
<p>
<p>FabricExplorerSpec defines the desired state of FabricExplorer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricExplorerStatus">FabricExplorerStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricExplorer">FabricExplorer</a>)
</p>
<p>
<p>FabricExplorerStatus defines the observed state of FabricExplorer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFSServer">FabricFSServer
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelAnchorPeer">FabricFollowerChannelAnchorPeer
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">FabricFollowerChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
<p>Host of the anchor peer</p>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<p>Port of the anchor peer</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelExternalPeer">FabricFollowerChannelExternalPeer
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">FabricFollowerChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>url</code>
<em>
string
</em>
</td>
<td>
<p>FabricPeer URL of the peer</p>
</td>
</tr>
<tr>
<td>
<code>tlsCACert</code>
<em>
string
</em>
</td>
<td>
<p>FabricPeer TLS CA certificate of the peer</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelOrderer">FabricFollowerChannelOrderer
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">FabricFollowerChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>url</code>
<em>
string
</em>
</td>
<td>
<p>URL of the orderer, e.g.: &ldquo;grpcs://xxxxx:443&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>certificate</code>
<em>
string
</em>
</td>
<td>
<p>TLS Certificate of the orderer node</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelPeer">FabricFollowerChannelPeer
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">FabricFollowerChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
<p>FabricPeer Name of the peer inside the kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code>
<em>
string
</em>
</td>
<td>
<p>FabricPeer Namespace of the peer inside the kubernetes cluster</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">FabricFollowerChannelSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannel">FabricFollowerChannel</a>)
</p>
<p>
<p>FabricFollowerChannelSpec defines the desired state of FabricFollowerChannel</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
<p>Name of the channel</p>
</td>
</tr>
<tr>
<td>
<code>mspId</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization to join the channel</p>
</td>
</tr>
<tr>
<td>
<code>orderers</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelOrderer">
[]FabricFollowerChannelOrderer
</a>
</em>
</td>
<td>
<p>Orderers to fetch the configuration block from</p>
</td>
</tr>
<tr>
<td>
<code>peersToJoin</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelPeer">
[]FabricFollowerChannelPeer
</a>
</em>
</td>
<td>
<p>Peers to join the channel</p>
</td>
</tr>
<tr>
<td>
<code>externalPeersToJoin</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelExternalPeer">
[]FabricFollowerChannelExternalPeer
</a>
</em>
</td>
<td>
<p>Peers to join the channel</p>
</td>
</tr>
<tr>
<td>
<code>anchorPeers</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelAnchorPeer">
[]FabricFollowerChannelAnchorPeer
</a>
</em>
</td>
<td>
<p>Anchor peers defined for the current organization</p>
</td>
</tr>
<tr>
<td>
<code>hlfIdentity</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.HLFIdentity">
HLFIdentity
</a>
</em>
</td>
<td>
<p>Identity to use to interact with the peers and the orderers</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelStatus">FabricFollowerChannelStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannel">FabricFollowerChannel</a>)
</p>
<p>
<p>FabricFollowerChannelStatus defines the observed state of FabricFollowerChannel</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricIstio">FabricIstio
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPISpec">FabricOperatorAPISpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.GRPCProxy">GRPCProxy</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>ingressGateway</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec">FabricMainChannelAdminOrdererOrganizationSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAdminPeerOrganizationSpec">FabricMainChannelAdminPeerOrganizationSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAnchorPeer">FabricMainChannelAnchorPeer
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
<p>Host of the peer</p>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<p>Port of the peer</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelApplicationConfig">FabricMainChannelApplicationConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConfig">FabricMainChannelConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>capabilities</code>
<em>
[]string
</em>
</td>
<td>
<p>Capabilities of the application channel configuration</p>
</td>
</tr>
<tr>
<td>
<code>policies</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig">
map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Policies of the application channel configuration</p>
</td>
</tr>
<tr>
<td>
<code>acls</code>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ACLs of the application channel configuration</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConfig">FabricMainChannelConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>application</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelApplicationConfig">
FabricMainChannelApplicationConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Application configuration of the channel</p>
</td>
</tr>
<tr>
<td>
<code>orderer</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererConfig">
FabricMainChannelOrdererConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Orderer configuration of the channel</p>
</td>
</tr>
<tr>
<td>
<code>capabilities</code>
<em>
[]string
</em>
</td>
<td>
<p>Capabilities for the channel</p>
</td>
</tr>
<tr>
<td>
<code>policies</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig">
map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Policies for the channel</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConsensusState">FabricMainChannelConsensusState
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererConfig">FabricMainChannelOrdererConfig</a>)
</p>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConsenter">FabricMainChannelConsenter
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
<p>Orderer host of the consenter</p>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<p>Orderer port of the consenter</p>
</td>
</tr>
<tr>
<td>
<code>tlsCert</code>
<em>
string
</em>
</td>
<td>
<p>TLS Certificate of the orderer node</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelEtcdRaft">FabricMainChannelEtcdRaft
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererConfig">FabricMainChannelOrdererConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>options</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelEtcdRaftOptions">
FabricMainChannelEtcdRaftOptions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelEtcdRaftOptions">FabricMainChannelEtcdRaftOptions
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelEtcdRaft">FabricMainChannelEtcdRaft</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tickInterval</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>electionTick</code>
<em>
uint32
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>heartbeatTick</code>
<em>
uint32
</em>
</td>
<td>
<p>HeartbeatTick is the number of ticks that must pass between heartbeats</p>
</td>
</tr>
<tr>
<td>
<code>maxInflightBlocks</code>
<em>
uint32
</em>
</td>
<td>
<p>MaxInflightBlocks is the maximum number of in-flight blocks that may be sent to followers at any given time.</p>
</td>
</tr>
<tr>
<td>
<code>snapshotIntervalSize</code>
<em>
uint32
</em>
</td>
<td>
<p>Maximum size of each raft snapshot file.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalOrdererNode">FabricMainChannelExternalOrdererNode
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererOrganization">FabricMainChannelOrdererOrganization</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
<p>Admin host of the orderer node</p>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<p>Admin port of the orderer node</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalOrdererOrganization">FabricMainChannelExternalOrdererOrganization
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization</p>
</td>
</tr>
<tr>
<td>
<code>tlsRootCert</code>
<em>
string
</em>
</td>
<td>
<p>TLS Root certificate authority of the orderer organization</p>
</td>
</tr>
<tr>
<td>
<code>signRootCert</code>
<em>
string
</em>
</td>
<td>
<p>Root certificate authority for signing</p>
</td>
</tr>
<tr>
<td>
<code>ordererEndpoints</code>
<em>
[]string
</em>
</td>
<td>
<p>Orderer endpoints for the organization in the channel configuration</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalPeerOrganization">FabricMainChannelExternalPeerOrganization
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization</p>
</td>
</tr>
<tr>
<td>
<code>tlsRootCert</code>
<em>
string
</em>
</td>
<td>
<p>TLS Root certificate authority of the orderer organization</p>
</td>
</tr>
<tr>
<td>
<code>signRootCert</code>
<em>
string
</em>
</td>
<td>
<p>Root certificate authority for signing</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelIdentity">FabricMainChannelIdentity
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secretNamespace</code>
<em>
string
</em>
</td>
<td>
<p>Secret namespace</p>
</td>
</tr>
<tr>
<td>
<code>secretName</code>
<em>
string
</em>
</td>
<td>
<p>Secret name</p>
</td>
</tr>
<tr>
<td>
<code>secretKey</code>
<em>
string
</em>
</td>
<td>
<p>Key inside the secret that holds the private key and certificate to interact with the network</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererBatchSize">FabricMainChannelOrdererBatchSize
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererConfig">FabricMainChannelOrdererConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxMessageCount</code>
<em>
int
</em>
</td>
<td>
<p>The number of transactions that can fit in a block.</p>
</td>
</tr>
<tr>
<td>
<code>absoluteMaxBytes</code>
<em>
int
</em>
</td>
<td>
<p>The absolute maximum size of a block, including all metadata.</p>
</td>
</tr>
<tr>
<td>
<code>preferredMaxBytes</code>
<em>
int
</em>
</td>
<td>
<p>The preferred maximum size of a block, including all metadata.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererConfig">FabricMainChannelOrdererConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConfig">FabricMainChannelConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ordererType</code>
<em>
string
</em>
</td>
<td>
<p>OrdererType of the consensus, default &ldquo;etcdraft&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>capabilities</code>
<em>
[]string
</em>
</td>
<td>
<p>Capabilities of the channel</p>
</td>
</tr>
<tr>
<td>
<code>policies</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig">
map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Policies of the orderer section of the channel</p>
</td>
</tr>
<tr>
<td>
<code>batchTimeout</code>
<em>
string
</em>
</td>
<td>
<p>Interval of the ordering service to create a block and send to the peers</p>
</td>
</tr>
<tr>
<td>
<code>batchSize</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererBatchSize">
FabricMainChannelOrdererBatchSize
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>state</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConsensusState">
FabricMainChannelConsensusState
</a>
</em>
</td>
<td>
<p>State about the channel, can only be <code>STATE_NORMAL</code> or <code>STATE_MAINTENANCE</code>.</p>
</td>
</tr>
<tr>
<td>
<code>etcdRaft</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelEtcdRaft">
FabricMainChannelEtcdRaft
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererNode">FabricMainChannelOrdererNode
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererOrganization">FabricMainChannelOrdererOrganization</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
<p>Name of the orderer node</p>
</td>
</tr>
<tr>
<td>
<code>namespace</code>
<em>
string
</em>
</td>
<td>
<p>Kubernetes namespace of the orderer node</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererOrganization">FabricMainChannelOrdererOrganization
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization</p>
</td>
</tr>
<tr>
<td>
<code>caName</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FabricCA Name of the organization</p>
</td>
</tr>
<tr>
<td>
<code>caNamespace</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FabricCA Namespace of the organization</p>
</td>
</tr>
<tr>
<td>
<code>tlsCACert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>TLS Root certificate authority of the orderer organization</p>
</td>
</tr>
<tr>
<td>
<code>signCACert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Root certificate authority for signing</p>
</td>
</tr>
<tr>
<td>
<code>ordererEndpoints</code>
<em>
[]string
</em>
</td>
<td>
<p>Orderer endpoints for the organization in the channel configuration</p>
</td>
</tr>
<tr>
<td>
<code>orderersToJoin</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererNode">
[]FabricMainChannelOrdererNode
</a>
</em>
</td>
<td>
<p>Orderer nodes within the kubernetes cluster to be added to the channel</p>
</td>
</tr>
<tr>
<td>
<code>externalOrderersToJoin</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalOrdererNode">
[]FabricMainChannelExternalOrdererNode
</a>
</em>
</td>
<td>
<p>External orderers to be added to the channel</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPeerOrganization">FabricMainChannelPeerOrganization
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
<p>MSP ID of the organization</p>
</td>
</tr>
<tr>
<td>
<code>caName</code>
<em>
string
</em>
</td>
<td>
<p>FabricCA Name of the organization</p>
</td>
</tr>
<tr>
<td>
<code>caNamespace</code>
<em>
string
</em>
</td>
<td>
<p>FabricCA Namespace of the organization</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPoliciesConfig">FabricMainChannelPoliciesConfig
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
string
</em>
</td>
<td>
<p>Type of policy, can only be <code>ImplicitMeta</code> or <code>Signature</code>.</p>
</td>
</tr>
<tr>
<td>
<code>rule</code>
<em>
string
</em>
</td>
<td>
<p>Rule of policy</p>
</td>
</tr>
<tr>
<td>
<code>modPolicy</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelSpec">FabricMainChannelSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannel">FabricMainChannel</a>)
</p>
<p>
<p>FabricMainChannelSpec defines the desired state of FabricMainChannel</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
<p>Name of the channel</p>
</td>
</tr>
<tr>
<td>
<code>identities</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelIdentity">
map[string]github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1.FabricMainChannelIdentity
</a>
</em>
</td>
<td>
<p>HLF Identities to be used to create and manage the channel</p>
</td>
</tr>
<tr>
<td>
<code>adminPeerOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAdminPeerOrganizationSpec">
[]FabricMainChannelAdminPeerOrganizationSpec
</a>
</em>
</td>
<td>
<p>Organizations that manage the <code>application</code> configuration of the channel</p>
</td>
</tr>
<tr>
<td>
<code>peerOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelPeerOrganization">
[]FabricMainChannelPeerOrganization
</a>
</em>
</td>
<td>
<p>Peer organizations that are external to the Kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>externalPeerOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalPeerOrganization">
[]FabricMainChannelExternalPeerOrganization
</a>
</em>
</td>
<td>
<p>External peer organizations that are inside the kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>channelConfig</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConfig">
FabricMainChannelConfig
</a>
</em>
</td>
<td>
<p>Configuration about the channel</p>
</td>
</tr>
<tr>
<td>
<code>adminOrdererOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec">
[]FabricMainChannelAdminOrdererOrganizationSpec
</a>
</em>
</td>
<td>
<p>Organizations that manage the <code>orderer</code> configuration of the channel</p>
</td>
</tr>
<tr>
<td>
<code>ordererOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelOrdererOrganization">
[]FabricMainChannelOrdererOrganization
</a>
</em>
</td>
<td>
<p>External orderer organizations that are inside the kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>externalOrdererOrganizations</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelExternalOrdererOrganization">
[]FabricMainChannelExternalOrdererOrganization
</a>
</em>
</td>
<td>
<p>Orderer organizations that are external to the Kubernetes cluster</p>
</td>
</tr>
<tr>
<td>
<code>orderers</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannelConsenter">
[]FabricMainChannelConsenter
</a>
</em>
</td>
<td>
<p>Consenters are the orderer nodes that are part of the channel consensus</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricMainChannelStatus">FabricMainChannelStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricMainChannel">FabricMainChannel</a>)
</p>
<p>
<p>FabricMainChannelStatus defines the observed state of FabricMainChannel</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfigSpec">FabricNetworkConfigSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfig">FabricNetworkConfig</a>)
</p>
<p>
<p>FabricNetworkConfigSpec defines the desired state of FabricNetworkConfig</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>organization</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>internal</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>organizations</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secretName</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfigStatus">FabricNetworkConfigStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricNetworkConfig">FabricNetworkConfig</a>)
</p>
<p>
<p>FabricNetworkConfigStatus defines the observed state of FabricNetworkConfig</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleAuth">FabricOperationsConsoleAuth
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleSpec">FabricOperationsConsoleSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>scheme</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>username</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>password</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleCouchDB">FabricOperationsConsoleCouchDB
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleSpec">FabricOperationsConsoleSpec</a>)
</p>
<p>
<p>FabricOperationsConsoleSpec defines the desired state of FabricOperationsConsole</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>username</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>password</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleSpec">FabricOperationsConsoleSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsole">FabricOperationsConsole</a>)
</p>
<p>
<p>FabricOperationsConsoleSpec defines the desired state of FabricOperationsConsole</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>auth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleAuth">
FabricOperationsConsoleAuth
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>couchDB</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleCouchDB">
FabricOperationsConsoleCouchDB
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>config</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>ingress</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">
Ingress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hostUrl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleStatus">FabricOperationsConsoleStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsole">FabricOperationsConsole</a>)
</p>
<p>
<p>FabricOperationsConsoleStatus defines the observed state of FabricOperationsConsole</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIAuth">FabricOperatorAPIAuth
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPISpec">FabricOperatorAPISpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>oidcJWKS</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>oidcIssuer</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIHLFConfig">FabricOperatorAPIHLFConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPISpec">FabricOperatorAPISpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>user</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>networkConfig</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPINetworkConfig">
FabricOperatorAPINetworkConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPINetworkConfig">FabricOperatorAPINetworkConfig
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIHLFConfig">FabricOperatorAPIHLFConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secretName</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>key</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPISpec">FabricOperatorAPISpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPI">FabricOperatorAPI</a>)
</p>
<p>
<p>FabricOperatorAPISpec defines the desired state of FabricOperatorAPI</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ingress</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">
Ingress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>auth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIAuth">
FabricOperatorAPIAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hlfConfig</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIHLFConfig">
FabricOperatorAPIHLFConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPIStatus">FabricOperatorAPIStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPI">FabricOperatorAPI</a>)
</p>
<p>
<p>FabricOperatorAPIStatus defines the observed state of FabricOperatorAPI</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorUIAuth">FabricOperatorUIAuth
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUISpec">FabricOperatorUISpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>oidcAuthority</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>oidcClientId</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>oidcScope</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorUISpec">FabricOperatorUISpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUI">FabricOperatorUI</a>)
</p>
<p>
<p>FabricOperatorUISpec defines the desired state of FabricOperatorUI</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>logoUrl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>auth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUIAuth">
FabricOperatorUIAuth
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ingress</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">
Ingress
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>apiUrl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOperatorUIStatus">FabricOperatorUIStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUI">FabricOperatorUI</a>)
</p>
<p>
<p>FabricOperatorUIStatus defines the observed state of FabricOperatorUI</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNode">FabricOrdererNode</a>)
</p>
<p>
<p>FabricOrdererNodeSpec defines the desired state of FabricOrdererNode</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>grpcProxy</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.GRPCProxy">
GRPCProxy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>updateCertificateTime</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>serviceMonitor</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">
ServiceMonitor
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hostAliases</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#hostalias-v1-core">
[]Kubernetes core/v1.HostAlias
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#nodeselector-v1-core">
Kubernetes core/v1.NodeSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>genesis</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>bootstrapMethod</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.BootstrapMethod">
BootstrapMethod
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>channelParticipationEnabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNodeService">
OrdererNodeService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secret</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Secret">
Secret
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>adminIstio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeStatus">FabricOrdererNodeStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNode">FabricOrdererNode</a>)
</p>
<p>
<p>FabricOrdererNodeStatus defines the observed state of FabricOrdererNode</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>lastCertificateUpdate</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>signCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tlsCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>signCaCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tlsCaCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tlsAdminCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>operationsPort</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>adminPort</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">FabricOrderingServiceSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingService">FabricOrderingService</a>)
</p>
<p>
<p>FabricOrderingServiceSpec defines the desired state of FabricOrderingService</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>enrollment</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererEnrollment">
OrdererEnrollment
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodes</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNode">
[]OrdererNode
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererService">
OrdererService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>systemChannel</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererSystemChannel">
OrdererSystemChannel
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceStatus">FabricOrderingServiceStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingService">FabricOrderingService</a>)
</p>
<p>
<p>FabricOrderingServiceStatus defines the observed state of FabricOrderingService</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchDB">FabricPeerCouchDB
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>user</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>password</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>externalCouchDB</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerExternalCouchDB">
FabricPeerExternalCouchDB
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchdbExporter">FabricPeerCouchdbExporter
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerDiscovery">FabricPeerDiscovery
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>period</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>touchPeriod</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerExternalCouchDB">FabricPeerExternalCouchDB
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchDB">FabricPeerCouchDB</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerLogging">FabricPeerLogging
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>level</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>peer</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cauthdsl</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>gossip</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>grpc</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ledger</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>msp</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>policies</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerResources">FabricPeerResources
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>peer</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>couchdb</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>chaincode</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>couchdbExporter</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>proxy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeer">FabricPeer</a>)
</p>
<p>
<p>FabricPeerSpec defines the desired state of FabricPeer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>updateCertificateTime</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>affinity</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>serviceMonitor</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">
ServiceMonitor
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>hostAliases</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#hostalias-v1-core">
[]Kubernetes core/v1.HostAlias
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#nodeselector-v1-core">
Kubernetes core/v1.NodeSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>couchDBexporter</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchdbExporter">
FabricPeerCouchdbExporter
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>grpcProxy</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.GRPCProxy">
GRPCProxy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>replicas</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>dockerSocketPath</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>externalBuilders</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ExternalBuilder">
[]ExternalBuilder
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>gossip</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpecGossip">
FabricPeerSpecGossip
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>externalEndpoint</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>external_chaincode_builder</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>couchdb</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerCouchDB">
FabricPeerCouchDB
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>fsServer</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFSServer">
FabricFSServer
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>mspID</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secret</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Secret">
Secret
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>service</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.PeerService">
PeerService
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>stateDb</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.StateDB">
StateDB
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storage</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerStorage">
FabricPeerStorage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>discovery</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerDiscovery">
FabricPeerDiscovery
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>logging</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerLogging">
FabricPeerLogging
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerResources">
FabricPeerResources
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
[]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tolerations</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>env</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerSpecGossip">FabricPeerSpecGossip
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>externalEndpoint</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>bootstrap</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>endpoint</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>useLeaderElection</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>orgLeader</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerStatus">FabricPeerStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeer">FabricPeer</a>)
</p>
<p>
<p>FabricPeerStatus defines the observed state of FabricPeer</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code>
<em>
github.com/kfsoftware/hlf-operator/pkg/status.Conditions
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>lastCertificateUpdate</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>signCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tlsCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>tlsCaCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>signCaCert</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricPeerStorage">FabricPeerStorage
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>couchdb</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>peer</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>chaincode</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Storage">
Storage
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.FabricTLSCACrypto">FabricTLSCACrypto
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAItemConf">FabricCAItemConf</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>key</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cert</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>clientAuth</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCAClientAuth">
FabricCAClientAuth
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.GRPCProxy">GRPCProxy
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>image</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tag</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>istio</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricIstio">
FabricIstio
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.HLFIdentity">HLFIdentity
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricFollowerChannelSpec">FabricFollowerChannelSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secretName</code>
<em>
string
</em>
</td>
<td>
<p>Secret name</p>
</td>
</tr>
<tr>
<td>
<code>secretNamespace</code>
<em>
string
</em>
</td>
<td>
<p>Secret namespace</p>
</td>
</tr>
<tr>
<td>
<code>secretKey</code>
<em>
string
</em>
</td>
<td>
<p>Key inside the secret that holds the private key and certificate to interact with the network</p>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Ingress">Ingress
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleSpec">FabricOperationsConsoleSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorAPISpec">FabricOperatorAPISpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperatorUISpec">FabricOperatorUISpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>className</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>annotations</code>
<em>
map[string]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tls</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#ingresstls-v1beta1-networking">
[]Kubernetes networking/v1beta1.IngressTLS
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>hosts</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.IngressHost">
[]IngressHost
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.IngressHost">IngressHost
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Ingress">Ingress</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>paths</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.IngressPath">
[]IngressPath
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.IngressPath">IngressPath
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.IngressHost">IngressHost</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>path</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>pathType</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.MetricsProvider">MetricsProvider
(<code>string</code> alias)</h3>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererCapabilities">OrdererCapabilities
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ChannelConfig">ChannelConfig</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>V2_0</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererEnrollment">OrdererEnrollment
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">FabricOrderingServiceSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>component</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Component">
Component
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tls</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.TLS">
TLS
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererNode">OrdererNode
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">FabricOrderingServiceSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>host</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>port</code>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>enrollment</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNodeEnrollment">
OrdererNodeEnrollment
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererNodeEnrollment">OrdererNodeEnrollment
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNode">OrdererNode</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tls</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNodeEnrollmentTLS">
OrdererNodeEnrollmentTLS
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererNodeEnrollmentTLS">OrdererNodeEnrollmentTLS
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererNodeEnrollment">OrdererNodeEnrollment</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>csr</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Csr">
Csr
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererNodeService">OrdererNodeService
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#servicetype-v1-core">
Kubernetes core/v1.ServiceType
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodePortOperations</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodePortRequest</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererService">OrdererService
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">FabricOrderingServiceSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceType">
ServiceType
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrdererSystemChannel">OrdererSystemChannel
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">FabricOrderingServiceSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>config</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ChannelConfig">
ChannelConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.OrgCertsRef">OrgCertsRef
</h3>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.PeerService">PeerService
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#servicetype-v1-core">
Kubernetes core/v1.ServiceType
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Secret">Secret
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enrollment</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Enrollment">
Enrollment
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Service">Service
</h3>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.ServiceType">
ServiceType
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ServiceMonitor">ServiceMonitor
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>labels</code>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>sampleLimit</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>interval</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>scrapeTimeout</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.ServiceType">ServiceType
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererService">OrdererService</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.Service">Service</a>)
</p>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.StateDB">StateDB
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerSpec">FabricPeerSpec</a>)
</p>
<p>
</p>
<h3 id="hlf.kungfusoftware.es/v1alpha1.Storage">Storage
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricCASpec">FabricCASpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOperationsConsoleCouchDB">FabricOperationsConsoleCouchDB</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrdererNodeSpec">FabricOrdererNodeSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricOrderingServiceSpec">FabricOrderingServiceSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricPeerStorage">FabricPeerStorage</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>size</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>storageClass</code>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>accessMode</code>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#persistentvolumeaccessmode-v1-core">
Kubernetes core/v1.PersistentVolumeAccessMode
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="hlf.kungfusoftware.es/v1alpha1.TLS">TLS
</h3>
<p>
(<em>Appears on:</em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Enrollment">Enrollment</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.FabricChaincodeSpec">FabricChaincodeSpec</a>, 
<a href="#hlf.kungfusoftware.es/v1alpha1.OrdererEnrollment">OrdererEnrollment</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cahost</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caname</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>caport</code>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>catls</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Catls">
Catls
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>csr</code>
<em>
<a href="#hlf.kungfusoftware.es/v1alpha1.Csr">
Csr
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>enrollid</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>enrollsecret</code>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>464adb2</code>.
</em></p>
