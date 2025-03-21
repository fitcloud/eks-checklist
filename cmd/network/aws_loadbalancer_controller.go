package network

import (
	"context"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 모든 namespace에서 deployement를 가져와서 container.image에 aws-load-balancer-controller가 포함되어 있는지 확인
func CheckAwsLoadBalancerController(client kubernetes.Interface) bool {
	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "aws-load-balancer-controller") {
				return true
			}
		}
	}

	return false
}
