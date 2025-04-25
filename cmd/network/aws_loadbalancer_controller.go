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
		Runbook:   "https://fitcloud.github.io/eks-checklist/network/albController",
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
