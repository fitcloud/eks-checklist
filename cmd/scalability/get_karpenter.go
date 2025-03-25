package scalability

import (
	"context"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 함수 분리
func GetKarpenter(client kubernetes.Interface) bool {
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
