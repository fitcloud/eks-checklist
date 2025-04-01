package scalability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckNodeGroupUsage는 Karpenter 전용 노드 그룹 또는 Fargate 사용 여부를 검사합니다.
func CheckNodeGroupUsage(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Karpenter 전용 노드 그룹 혹은 Fargate 사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Karpenter 전용 노드 그룹 또는 Fargate가 사용되고 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	var isKarpenterNode bool
	var isFargateNode bool

	for _, node := range nodes.Items {
		// Karpenter 라벨 존재 시
		if _, found := node.Labels["karpenter.sh/provisioner-name"]; found {
			isKarpenterNode = true
			result.Resources = append(result.Resources, fmt.Sprintf("Karpenter Node: %s", node.Name))
		}

		// Fargate 라벨 존재 시
		if profile, found := node.Labels["eks.amazonaws.com/fargate-profile"]; found && profile != "" {
			isFargateNode = true
			result.Resources = append(result.Resources, fmt.Sprintf("Fargate Node: %s (profile: %s)", node.Name, profile))
		}
	}

	if isKarpenterNode || isFargateNode {
		result.Passed = true
	} else {
		result.Passed = false
	}

	return result
}
