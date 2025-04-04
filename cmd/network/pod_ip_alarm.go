package network

import (
	"eks-checklist/cmd/common"
)

func CheckPodIPAlarm() common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Pod에 부여할 IP 부족시 알림 설정 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "EKS Node가 배포되는 서브넷의 할당가능한 IP 개수가 부족하면 알람을 받도록 설정하세요 (CloudWatch Alarm, Prometheus …etc)",
		Runbook:    "https://fitcloud.github.io/eks-checklist/컨테이너_이미지_태그에_latest_미사용/",
	}

	return result
}
