package reliability

import (
	"eks-checklist/cmd/common"
)

func CheckApplicationLogs() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[REL-010] 반영구 저장소에 애플리케이션 로그 저장",
		Manual:     true,
		Passed:     false,
		FailureMsg: "application의 로그는 Opensearch, Cloudwatch Logs 등 영구 저장소에 수집하는 것이 좋습니다",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/reliability/REL-010",
	}

	return result
}
