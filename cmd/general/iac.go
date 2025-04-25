package general

import (
	"eks-checklist/cmd/common"
)

func CheckIAC() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "코드형 인프라 (EKS 클러스터, 애플리케이션 배포)",
		Manual:     true,
		Passed:     false,
		FailureMsg: "클러스터 및 Application은 IaC 방식으로 관리하는 것이 좋습니다 (Terraform ,CDK, CloudFormation)",
		Runbook:    "https://fitcloud.github.io/eks-checklist/general/iac",
	}

	return result
}
