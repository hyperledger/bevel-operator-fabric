/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

import (
	hlfkungfusoftwareesv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FabricPeerSpecApplyConfiguration represents an declarative configuration of the FabricPeerSpec type for use
// with apply.
type FabricPeerSpecApplyConfiguration struct {
	UpdateCertificateTime    *v1.Time                                     `json:"updateCertificateTime,omitempty"`
	Affinity                 *corev1.Affinity                             `json:"affinity,omitempty"`
	ServiceMonitor           *ServiceMonitorApplyConfiguration            `json:"serviceMonitor,omitempty"`
	HostAliases              []corev1.HostAlias                           `json:"hostAliases,omitempty"`
	NodeSelector             *corev1.NodeSelector                         `json:"nodeSelector,omitempty"`
	CouchDBExporter          *FabricPeerCouchdbExporterApplyConfiguration `json:"couchDBexporter,omitempty"`
	GRPCProxy                *GRPCProxyApplyConfiguration                 `json:"grpcProxy,omitempty"`
	Replicas                 *int                                         `json:"replicas,omitempty"`
	DockerSocketPath         *string                                      `json:"dockerSocketPath,omitempty"`
	Image                    *string                                      `json:"image,omitempty"`
	ExternalBuilders         []ExternalBuilderApplyConfiguration          `json:"externalBuilders,omitempty"`
	Istio                    *FabricIstioApplyConfiguration               `json:"istio,omitempty"`
	Gossip                   *FabricPeerSpecGossipApplyConfiguration      `json:"gossip,omitempty"`
	ExternalEndpoint         *string                                      `json:"externalEndpoint,omitempty"`
	Tag                      *string                                      `json:"tag,omitempty"`
	ImagePullPolicy          *corev1.PullPolicy                           `json:"imagePullPolicy,omitempty"`
	ExternalChaincodeBuilder *bool                                        `json:"external_chaincode_builder,omitempty"`
	CouchDB                  *FabricPeerCouchDBApplyConfiguration         `json:"couchdb,omitempty"`
	FSServer                 *FabricFSServerApplyConfiguration            `json:"fsServer,omitempty"`
	ImagePullSecrets         []corev1.LocalObjectReference                `json:"imagePullSecrets,omitempty"`
	MspID                    *string                                      `json:"mspID,omitempty"`
	Secret                   *SecretApplyConfiguration                    `json:"secret,omitempty"`
	Service                  *PeerServiceApplyConfiguration               `json:"service,omitempty"`
	StateDb                  *hlfkungfusoftwareesv1alpha1.StateDB         `json:"stateDb,omitempty"`
	Storage                  *FabricPeerStorageApplyConfiguration         `json:"storage,omitempty"`
	Discovery                *FabricPeerDiscoveryApplyConfiguration       `json:"discovery,omitempty"`
	Logging                  *FabricPeerLoggingApplyConfiguration         `json:"logging,omitempty"`
	Resources                *FabricPeerResourcesApplyConfiguration       `json:"resources,omitempty"`
	Hosts                    []string                                     `json:"hosts,omitempty"`
	Tolerations              []corev1.Toleration                          `json:"tolerations,omitempty"`
	Env                      []corev1.EnvVar                              `json:"env,omitempty"`
}

// FabricPeerSpecApplyConfiguration constructs an declarative configuration of the FabricPeerSpec type for use with
// apply.
func FabricPeerSpec() *FabricPeerSpecApplyConfiguration {
	return &FabricPeerSpecApplyConfiguration{}
}

// WithUpdateCertificateTime sets the UpdateCertificateTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UpdateCertificateTime field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithUpdateCertificateTime(value v1.Time) *FabricPeerSpecApplyConfiguration {
	b.UpdateCertificateTime = &value
	return b
}

// WithAffinity sets the Affinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Affinity field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithAffinity(value corev1.Affinity) *FabricPeerSpecApplyConfiguration {
	b.Affinity = &value
	return b
}

// WithServiceMonitor sets the ServiceMonitor field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceMonitor field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithServiceMonitor(value *ServiceMonitorApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.ServiceMonitor = value
	return b
}

// WithHostAliases adds the given value to the HostAliases field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the HostAliases field.
func (b *FabricPeerSpecApplyConfiguration) WithHostAliases(values ...corev1.HostAlias) *FabricPeerSpecApplyConfiguration {
	for i := range values {
		b.HostAliases = append(b.HostAliases, values[i])
	}
	return b
}

// WithNodeSelector sets the NodeSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NodeSelector field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithNodeSelector(value corev1.NodeSelector) *FabricPeerSpecApplyConfiguration {
	b.NodeSelector = &value
	return b
}

// WithCouchDBExporter sets the CouchDBExporter field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CouchDBExporter field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithCouchDBExporter(value *FabricPeerCouchdbExporterApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.CouchDBExporter = value
	return b
}

// WithGRPCProxy sets the GRPCProxy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GRPCProxy field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithGRPCProxy(value *GRPCProxyApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.GRPCProxy = value
	return b
}

// WithReplicas sets the Replicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Replicas field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithReplicas(value int) *FabricPeerSpecApplyConfiguration {
	b.Replicas = &value
	return b
}

// WithDockerSocketPath sets the DockerSocketPath field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DockerSocketPath field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithDockerSocketPath(value string) *FabricPeerSpecApplyConfiguration {
	b.DockerSocketPath = &value
	return b
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithImage(value string) *FabricPeerSpecApplyConfiguration {
	b.Image = &value
	return b
}

// WithExternalBuilders adds the given value to the ExternalBuilders field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ExternalBuilders field.
func (b *FabricPeerSpecApplyConfiguration) WithExternalBuilders(values ...*ExternalBuilderApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithExternalBuilders")
		}
		b.ExternalBuilders = append(b.ExternalBuilders, *values[i])
	}
	return b
}

// WithIstio sets the Istio field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Istio field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithIstio(value *FabricIstioApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Istio = value
	return b
}

// WithGossip sets the Gossip field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Gossip field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithGossip(value *FabricPeerSpecGossipApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Gossip = value
	return b
}

// WithExternalEndpoint sets the ExternalEndpoint field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExternalEndpoint field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithExternalEndpoint(value string) *FabricPeerSpecApplyConfiguration {
	b.ExternalEndpoint = &value
	return b
}

// WithTag sets the Tag field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Tag field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithTag(value string) *FabricPeerSpecApplyConfiguration {
	b.Tag = &value
	return b
}

// WithImagePullPolicy sets the ImagePullPolicy field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ImagePullPolicy field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithImagePullPolicy(value corev1.PullPolicy) *FabricPeerSpecApplyConfiguration {
	b.ImagePullPolicy = &value
	return b
}

// WithExternalChaincodeBuilder sets the ExternalChaincodeBuilder field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExternalChaincodeBuilder field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithExternalChaincodeBuilder(value bool) *FabricPeerSpecApplyConfiguration {
	b.ExternalChaincodeBuilder = &value
	return b
}

// WithCouchDB sets the CouchDB field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CouchDB field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithCouchDB(value *FabricPeerCouchDBApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.CouchDB = value
	return b
}

// WithFSServer sets the FSServer field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FSServer field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithFSServer(value *FabricFSServerApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.FSServer = value
	return b
}

// WithImagePullSecrets adds the given value to the ImagePullSecrets field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ImagePullSecrets field.
func (b *FabricPeerSpecApplyConfiguration) WithImagePullSecrets(values ...corev1.LocalObjectReference) *FabricPeerSpecApplyConfiguration {
	for i := range values {
		b.ImagePullSecrets = append(b.ImagePullSecrets, values[i])
	}
	return b
}

// WithMspID sets the MspID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MspID field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithMspID(value string) *FabricPeerSpecApplyConfiguration {
	b.MspID = &value
	return b
}

// WithSecret sets the Secret field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Secret field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithSecret(value *SecretApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Secret = value
	return b
}

// WithService sets the Service field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Service field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithService(value *PeerServiceApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Service = value
	return b
}

// WithStateDb sets the StateDb field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the StateDb field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithStateDb(value hlfkungfusoftwareesv1alpha1.StateDB) *FabricPeerSpecApplyConfiguration {
	b.StateDb = &value
	return b
}

// WithStorage sets the Storage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Storage field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithStorage(value *FabricPeerStorageApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Storage = value
	return b
}

// WithDiscovery sets the Discovery field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Discovery field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithDiscovery(value *FabricPeerDiscoveryApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Discovery = value
	return b
}

// WithLogging sets the Logging field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Logging field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithLogging(value *FabricPeerLoggingApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Logging = value
	return b
}

// WithResources sets the Resources field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Resources field is set to the value of the last call.
func (b *FabricPeerSpecApplyConfiguration) WithResources(value *FabricPeerResourcesApplyConfiguration) *FabricPeerSpecApplyConfiguration {
	b.Resources = value
	return b
}

// WithHosts adds the given value to the Hosts field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Hosts field.
func (b *FabricPeerSpecApplyConfiguration) WithHosts(values ...string) *FabricPeerSpecApplyConfiguration {
	for i := range values {
		b.Hosts = append(b.Hosts, values[i])
	}
	return b
}

// WithTolerations adds the given value to the Tolerations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tolerations field.
func (b *FabricPeerSpecApplyConfiguration) WithTolerations(values ...corev1.Toleration) *FabricPeerSpecApplyConfiguration {
	for i := range values {
		b.Tolerations = append(b.Tolerations, values[i])
	}
	return b
}

// WithEnv adds the given value to the Env field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Env field.
func (b *FabricPeerSpecApplyConfiguration) WithEnv(values ...corev1.EnvVar) *FabricPeerSpecApplyConfiguration {
	for i := range values {
		b.Env = append(b.Env, values[i])
	}
	return b
}
