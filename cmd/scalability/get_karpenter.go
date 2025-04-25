package scalability

import (
	"context"
	"strings"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetKarpenter checks if the Karpenter deployment is installed in the cluster.
func GetKarpenter(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Karpenter 사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Karpenter 관련 이미지가 포함된 Deployment를 찾을 수 없습니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/scalability/karpenterEnabled",
	}

	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "karpenter") {
				result.Passed = true
				return result
			}
		}
	}

	result.Passed = false
	return result
}
