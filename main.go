package main

import (
	"os"

	wfclientset "github.com/argoproj/argo/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/k0kubun/pp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	pp.Println("##############")
	listOpts := metav1.ListOptions{}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)

	restConfig, err := clientConfig.ClientConfig()
	wfcs := wfclientset.NewForConfigOrDie(restConfig)
	wfClient := wfcs.ArgoprojV1alpha1().Workflows("argo")
	list, err := wfClient.List(listOpts)
	if err != nil {
		pp.Println(err)
	}
	pp.Println(list)
}
