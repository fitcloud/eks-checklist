package scalability

import (
	"eks-checklist/cmd/common"
)

func CheckNodeScalingPolicy() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[SCL-006] 노드 확장/축소 정책 적용",
		Manual:     true,
		Passed:     false,
		FailureMsg: "EKS Node는 AutoscaleGroup 또는 Karpenter Nodepool과 같은 동적 프로비저닝 하는 것이 좋습니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/scalability/SCL-006",
	}

	return result
}
