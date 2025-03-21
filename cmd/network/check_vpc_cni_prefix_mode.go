package network

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 모든 namespace에서 aws-node 데몬셋을 가져와서 containers.env.key 중 ENABLE_PREFIX_DELEGATION의 값이 true 인지 false인지 확인
func CheckVpcCniPrefixMode(client kubernetes.Interface) bool {
	// instance가 nitro 기반인이 확인하는 로직이 필요하나 문제점에 봉착함
	// api 중에 노드만 가져오는게 없어서 nodegroup 단위로 검색하면 karpenter node가 누락되고
	// tag key등으로 구분해도 nodegroup이랑 karpenter랑 서로 tag가 다름
	// 그럼으로 음... 보류 질문해야할 듯

	// 모든 namespace에서 aws-node 데몬셋을 가져옴
	daemonsets, err := client.AppsV1().DaemonSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, daemonset := range daemonsets.Items {
		if daemonset.Name == "aws-node" {
			for _, container := range daemonset.Spec.Template.Spec.Containers {
				for _, env := range container.Env {
					if env.Name == "ENABLE_PREFIX_DELEGATION" {
						if env.Value == "true" {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
