package security

import (
	"eks-checklist/cmd/common"
)

func CheckAccessAlarm() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "비정상 접근에 대한 알림 설정 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "EKS endpoint는 인증받은 권한만 접근해야합니다 (GuarDuty, Proemetheus + Altermanager)",
		Runbook:    "https://fitcloud.github.io/eks-checklist/security/unauthorizedAlerting",
	}

	return result
}
