package argo

import (
	"os"

	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/locona/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	"k8s.io/client-go/tools/clientcmd"
)

type argo struct {
	client v1alpha1.WorkflowInterface
}

func New(ns string) (*argo, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	wfcs := wfclientset.NewForConfigOrDie(restConfig)
	wfClient := wfcs.ArgoprojV1alpha1().Workflows(ns)
	return &argo{client: wfClient}, nil
}
