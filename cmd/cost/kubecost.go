package cost

import (
	"context"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetKubecost(client kubernetes.Interface) bool {
	deploys, err := client.AppsV1().Deployments(v1.NamespaceAll).List(context.TODO(), v1.ListOptions{})
	//	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{}) 기존 코드에서 all-namespace로 변경

	if err != nil {
		panic(err.Error())
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "kubecost") {
				return true
			}
		}
	}

	return false
}
