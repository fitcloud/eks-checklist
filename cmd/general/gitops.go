package general

import (
	"eks-checklist/cmd/common"
)

func CheckGitOps() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "GitOps 적용",
		Manual:     true,
		Passed:     false,
		FailureMsg: "Git 기반의 배포 도구를 사용하여 워크로드의 일관성을 유지하는것이 좋습니다 (ArgoCD, FluxCD,l Jenkins X …etc)",
		Runbook:    "https://fitcloud.github.io/eks-checklist/general/gitOps",
	}

	return result
}
