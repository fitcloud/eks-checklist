package reliability

import (
	"eks-checklist/cmd/common"
)

func CheckNodeScalingPolicy() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[REL-009] 인프라 및 애플리케이션 모니터링 스택 적용",
		Manual:     true,
		Passed:     false,
		FailureMsg: "EKS와 워크로드 전체를 확인할 수 있는 모니터링 스택이 있는 것이 좋습니다 ( kube-prometheus-stack, cloudwatch ..etc)",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/reliability/REL-009",
	}

	return result
}
