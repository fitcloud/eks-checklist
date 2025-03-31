// 변경 전 코드

// package network

// import (
// 	"context"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // 모든 namespace에서 aws-node 데몬셋을 가져와서 containers.env.key 중 ENABLE_PREFIX_DELEGATION의 값이 true 인지 false인지 확인
// func CheckVpcCniPrefixMode(client kubernetes.Interface) bool {
// 	// instance가 nitro 기반인이 확인하는 로직이 필요하나 문제점에 봉착함
// 	// api 중에 노드만 가져오는게 없어서 nodegroup 단위로 검색하면 karpenter node가 누락되고
// 	// tag key등으로 구분해도 nodegroup이랑 karpenter랑 서로 tag가 다름
// 	// 그럼으로 음... 보류 질문해야할 듯

// 	// 모든 namespace에서 aws-node 데몬셋을 가져옴
// 	daemonsets, err := client.AppsV1().DaemonSets("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for _, daemonset := range daemonsets.Items {
// 		if daemonset.Name == "aws-node" {
// 			for _, container := range daemonset.Spec.Template.Spec.Containers {
// 				for _, env := range container.Env {
// 					if env.Name == "ENABLE_PREFIX_DELEGATION" {
// 						if env.Value == "true" {
// 							return true
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

package network

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckVpcCniPrefixMode checks if the aws-node DaemonSet has ENABLE_PREFIX_DELEGATION=true.
func CheckVpcCniPrefixMode(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "VPC CNI의 Prefix 모드 사용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	daemonsets, err := client.AppsV1().DaemonSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	found := false
	for _, ds := range daemonsets.Items {
		if ds.Name != "aws-node" {
			continue
		}
		found = true

		for _, container := range ds.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				if env.Name == "ENABLE_PREFIX_DELEGATION" {
					if env.Value == "true" {
						result.Passed = true
						// result.SuccessMsg = "aws-node에서 ENABLE_PREFIX_DELEGATION=true로 설정되어 Prefix 모드가 활성화되어 있습니다."
						// result.Resources = append(result.Resources,
						// 	fmt.Sprintf("Namespace: %s | DaemonSet: %s | Env: %s=%s", ds.Namespace, ds.Name, env.Name, env.Value))
						return result
					}
					// false로 설정된 경우
					result.Passed = false
					result.FailureMsg = "aws-node에서 ENABLE_PREFIX_DELEGATION이 false로 설정되어 Prefix 모드가 비활성화되어 있습니다."
					result.Resources = append(result.Resources,
						fmt.Sprintf("Namespace: %s | DaemonSet: %s | Env: %s=%s", ds.Namespace, ds.Name, env.Name, env.Value))
					return result
				}
			}
		}

		// ENABLE_PREFIX_DELEGATION 환경 변수가 없는 경우
		result.Passed = false
		result.FailureMsg = "aws-node DaemonSet에서 ENABLE_PREFIX_DELEGATION 환경 변수가 설정되어 있지 않습니다."
		result.Resources = append(result.Resources,
			fmt.Sprintf("Namespace: %s | DaemonSet: %s", ds.Namespace, ds.Name))
		return result
	}

	if !found {
		result.Passed = false
		result.FailureMsg = "aws-node DaemonSet을 찾을 수 없습니다."
	}

	return result
}
