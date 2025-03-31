// 변경 전 코드

// package network

// import (
// 	"context"
// 	"strings"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // 모든 namespace에서 deployement를 가져와서 container.image에 aws-load-balancer-controller가 포함되어 있는지 확인
// func CheckAwsLoadBalancerController(client kubernetes.Interface) bool {
// 	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for _, deploy := range deploys.Items {
// 		for _, container := range deploy.Spec.Template.Spec.Containers {
// 			if strings.Contains(container.Image, "aws-load-balancer-controller") {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }

// 변경 후 코드
package network

import (
	"context"
	"strings"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckAwsLoadBalancerController checks if AWS Load Balancer Controller is installed via Deployment.
func CheckAwsLoadBalancerController(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "AWS Load Balancer Controller 사용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "aws-load-balancer-controller") {
				result.Passed = true
				// result.SuccessMsg = "AWS Load Balancer Controller가 설치되어 있습니다."
				// result.Resources = append(result.Resources,
				// 	fmt.Sprintf("Namespace: %s | Deployment: %s | Image: %s",
				// 		deploy.Namespace, deploy.Name, container.Image))
				return result
			}
		}
	}

	result.Passed = false
	result.FailureMsg = "AWS Load Balancer Controller가 설치되어 있지 않습니다."
	return result
}
