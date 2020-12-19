package helpers

import (
	"bytes"
	"context"
	"io"
	"os/exec"

	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient provides k8s client for kubeconfig
func GetKubeClient() (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	kubeClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClientset, nil
}

// GetKubeExtensionClient provides k8s client for CRDs
func GetKubeExtensionClient() (*apiextension.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	extClient, err := apiextension.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return extClient, nil
}

// GetKubeOperatorClient provides k8s client for operator
func GetKubeOperatorClient() (*operatorv1.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	kubeClientset, err := operatorv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClientset, nil
}

// ExecKubectl executes the given command using `kubectl`
func ExecKubectl(ctx context.Context, args ...string) ([]byte, error) {
	var stdout, stderr, combined bytes.Buffer

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	cmd.Stdout = io.MultiWriter(&stdout, &combined)
	cmd.Stderr = io.MultiWriter(&stderr, &combined)
	if err := cmd.Run(); err != nil {
		return nil, errors.Errorf("kubectl command failed (%s). output=%s", err, combined.String())
	}
	return stdout.Bytes(), nil
}

// ToYaml takes a slice of values, and returns corresponding YAML
// representation as a string slice
func ToYaml(objs []runtime.Object) ([]string, error) {
	manifests := make([]string, len(objs))
	for i, obj := range objs {
		o, err := yaml.Marshal(obj)
		if err != nil {
			return []string{}, err
		}
		manifests[i] = string(o)
	}

	return manifests, nil
}
