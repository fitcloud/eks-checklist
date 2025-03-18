package cmd

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getKubeconfig(kubeconfigPath string, awsProfile string) rest.Config {
	var kubeconfig *string

	if home := homedir.HomeDir(); kubeconfigPath == "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", kubeconfigPath, "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// AWS_PROFILE 설정
	if awsProfile != "" {
		os.Setenv("AWS_PROFILE", awsProfile)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return *config
}

func getEksClusterName(kubeconfig rest.Config) string {
	clusterNameIdx := slices.Index(kubeconfig.ExecProvider.Args, "--cluster-name") + 1

	return kubeconfig.ExecProvider.Args[clusterNameIdx]
}

func createK8sClient(kubeconfig rest.Config) kubernetes.Interface {
	client, err := kubernetes.NewForConfig(&kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return client
}

func getKarpenter(client kubernetes.Interface) bool {
	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "karpenter") {
				return true
			}
		}
	}

	return false
}
