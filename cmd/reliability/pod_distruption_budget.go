package reliability

import (
	"eks-checklist/cmd/common"
)

func CheckPDB() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[REL-006] 중요 워크로드에 대한 PDB(Pod Distruption Budget) 적용",
		Manual:     true,
		Passed:     false,
		FailureMsg: "중요 워크로드 application은 PDB 설정을 통해 가용성을 지키는 것이 좋습니다",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/reliability/REL-006",
	}

	return result
}
