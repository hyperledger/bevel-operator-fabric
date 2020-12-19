package peer

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
	"strings"
)

type Options struct {
	Name           string
	StorageClass   string
	Capacity       string
	NS             string
	Image          string
	Version        string
	MspID          string
	StateDB        string
	EnrollPW       string
	CAName         string
	EnrollID       string
	Hosts          []string
	BootstrapPeers []string
	Leader         bool
}

func (o Options) Validate() error {
	return nil
}

type createCmd struct {
	out      io.Writer
	errOut   io.Writer
	peerOpts Options
}

func (c *createCmd) validate() error {
	c.peerOpts.Image = helpers.DefaultPeerImage
	return c.peerOpts.Validate()
}
func (c *createCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	certAuth, err := helpers.GetCertAuthByFullName(oclient, c.peerOpts.CAName)
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	k8sIP, err := utils.GetPublicIPKubernetes(clientSet)
	if err != nil {
		return err
	}
	var bootstrapPeerURL string
	if len(c.peerOpts.BootstrapPeers) > 0 {
		var bootstrapPeerUrls []string
		for _, bp := range c.peerOpts.BootstrapPeers {
			boostrapPeer, err := helpers.GetPeerByFullName(oclient, bp)
			if err != nil {
				return err
			}
			if boostrapPeer.Status.Status != v1alpha1.RunningStatus {
				return errors.Errorf("Peer %s is not running", boostrapPeer.Name)
			}
			u, err := url.Parse(boostrapPeer.Status.URL)
			if err != nil {
				return err
			}
			chunks := strings.Split(u.Host, ":")
			ip := chunks[0]
			port := chunks[1]
			bootstrapPeerURL := fmt.Sprintf("%s:%s", ip, port)
			bootstrapPeerUrls = append(bootstrapPeerUrls, bootstrapPeerURL)
			log.Infof("Bootstrap peer url %s ip=%s port=%s", bootstrapPeerURL, ip, port)
		}
		bootstrapPeerURL = strings.Join(bootstrapPeerUrls, " ")
	} else {
		bootstrapPeerURL = ""
	}

	externalEndpoint := ""
	if len(c.peerOpts.Hosts) > 0 {
		externalEndpoint = fmt.Sprintf("%s:443", c.peerOpts.Hosts[0])
	}
	fabricPeer := &v1alpha1.FabricPeer{
		ObjectMeta: v1.ObjectMeta{
			Name:      c.peerOpts.Name,
			Namespace: c.peerOpts.NS,
		},
		Spec: v1alpha1.FabricPeerSpec{
			DockerSocketPath: "/var/run/docker.sock",
			Image:            c.peerOpts.Image,
			Istio: v1alpha1.FabricPeerIstio{
				Port: 443,
			},
			Gossip: v1alpha1.FabricPeerSpecGossip{
				ExternalEndpoint:  externalEndpoint,
				Bootstrap:         bootstrapPeerURL,
				Endpoint:          "",
				UseLeaderElection: !c.peerOpts.Leader,
				OrgLeader:         c.peerOpts.Leader,
			},
			ExternalEndpoint:         externalEndpoint,
			Tag:                      c.peerOpts.Version,
			ExternalChaincodeBuilder: true,
			CouchDB: v1alpha1.FabricPeerCouchDB{
				User:     "couchdb",
				Password: "couchdb",
			},
			MspID: c.peerOpts.MspID,
			Secret: v1alpha1.Secret{
				Enrollment: v1alpha1.Enrollment{
					Component: v1alpha1.Component{
						Cahost: certAuth.Status.Host,
						Caname: certAuth.Spec.CA.Name,
						Caport: certAuth.Status.Port,
						Catls: v1alpha1.Catls{
							Cacert: base64.StdEncoding.EncodeToString([]byte(certAuth.Status.TlsCert)),
						},
						Enrollid:     c.peerOpts.EnrollID,
						Enrollsecret: c.peerOpts.EnrollPW,
					},
					TLS: v1alpha1.TLS{
						Cahost: certAuth.Status.Host,
						Caname: certAuth.Spec.TLSCA.Name,
						Caport: certAuth.Status.Port,
						Catls: v1alpha1.Catls{
							Cacert: base64.StdEncoding.EncodeToString([]byte(certAuth.Status.TlsCert)),
						},
						Csr: v1alpha1.Csr{
							Hosts: []string{
								"127.0.0.1",
								"localhost",
								k8sIP,
							},
							CN: "",
						},
						Enrollid:     c.peerOpts.EnrollID,
						Enrollsecret: c.peerOpts.EnrollPW,
					},
				},
			},
			Service: v1alpha1.PeerService{
				Type: "NodePort",
			},
			StateDb: v1alpha1.StateDB(c.peerOpts.StateDB),
			Storage: v1alpha1.FabricPeerStorage{
				CouchDB: v1alpha1.Storage{
					Size:         "5Gi",
					StorageClass: c.peerOpts.StorageClass,
					AccessMode:   "ReadWriteOnce",
				},
				Peer: v1alpha1.Storage{
					Size:         "5Gi",
					StorageClass: c.peerOpts.StorageClass,
					AccessMode:   "ReadWriteOnce",
				},
				Chaincode: v1alpha1.Storage{
					Size:         "5Gi",
					StorageClass: c.peerOpts.StorageClass,
					AccessMode:   "ReadWriteOnce",
				},
			},
			Discovery: v1alpha1.FabricPeerDiscovery{
				Period:      "60s",
				TouchPeriod: "60s",
			},
			Logging: v1alpha1.FabricPeerLogging{
				Level:    "info",
				Peer:     "info",
				Cauthdsl: "info",
				Gossip:   "info",
				Grpc:     "info",
				Ledger:   "info",
				Msp:      "info",
				Policies: "info",
			},
			Resources: v1alpha1.FabricPeerResources{
				Peer: v1alpha1.Resources{
					Requests: v1alpha1.Requests{
						CPU:    "10m",
						Memory: "10M",
					},
					Limits: v1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4096M",
					},
				},
				CouchDB: v1alpha1.Resources{
					Requests: v1alpha1.Requests{
						CPU:    "10m",
						Memory: "10M",
					},
					Limits: v1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4096M",
					},
				},
				Chaincode: v1alpha1.Resources{
					Requests: v1alpha1.Requests{
						CPU:    "10m",
						Memory: "10M",
					},
					Limits: v1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4096M",
					},
				},
			},
			Hosts:          c.peerOpts.Hosts,
			OperationHosts: []string{},
			OperationIPs:   []string{},
		},
	}
	ctx := context.Background()
	_, err = oclient.HlfV1alpha1().FabricPeers(c.peerOpts.NS).Create(
		ctx,
		fabricPeer,
		v1.CreateOptions{},
	)
	if err != nil {
		return err
	}
	return nil
}
func newCreatePeerCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Fabric Peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.peerOpts.Name, "name", "", "name of the Fabric Peer to create")
	f.StringVar(&c.peerOpts.CAName, "ca-name", "", "name of the Fabric Peer")
	f.StringVar(&c.peerOpts.EnrollID, "enroll-id", "", "user to enroll peer certificates")
	f.StringVar(&c.peerOpts.EnrollPW, "enroll-pw", "", "password to enroll peer certificates")
	f.StringVar(&c.peerOpts.Capacity, "capacity", "5Gi", "total raw capacity of Fabric Peer in this zone, e.g. 16Ti")
	f.StringVarP(&c.peerOpts.NS, "namespace", "n", helpers.DefaultNamespace, "namespace scope for this request")
	f.StringVarP(&c.peerOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "storage class for this Fabric Peer")
	f.StringVarP(&c.peerOpts.Image, "image", "", helpers.DefaultPeerImage, "version of the Fabric Peer")
	f.StringVarP(&c.peerOpts.Version, "version", "", helpers.DefaultPeerVersion, "version of the Fabric Peer")
	f.StringVarP(&c.peerOpts.MspID, "mspid", "", "", "MSP ID of the organization")
	f.StringVarP(&c.peerOpts.StateDB, "statedb", "", "leveldb", "State database")
	f.BoolVarP(&c.peerOpts.Leader, "leader", "", false, "Force peer to be leader")
	f.StringArrayVarP(&c.peerOpts.BootstrapPeers, "bootstrap-peer", "", []string{}, "Bootstrap peers")
	f.StringArrayVarP(&c.peerOpts.Hosts, "hosts", "", []string{}, "External hosts")
	return cmd
}
