package security

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckPVEcryption - Persistent Volume (PV)의 암호화 상태를 확인
func CheckPVEcryption(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "PV 암호화",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 PV가 암호화되지 않았거나 수동 확인이 필요합니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	pvs, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, pv := range pvs.Items {
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == "ebs.csi.aws.com" {
			encrypted, exists := pv.Spec.CSI.VolumeAttributes["encrypted"]
			if exists && encrypted == "true" {
				// 암호화됨 - PASS (추가 조치 없음)
				continue
			} else {
				// 암호화되지 않음 - FAIL
				result.Passed = false
				result.Resources = append(result.Resources,
					fmt.Sprintf("PV: %s (EBS 미암호화)", pv.Name))
			}
		} else {
			// CSI 정보가 nil인 경우엔 "N/A"로 처리
			var driver string
			if pv.Spec.CSI != nil {
				driver = pv.Spec.CSI.Driver
			} else {
				driver = "N/A"
			}
			// EBS가 아닌 PV의 경우 암호화 상태 판단 불가, 수동 확인 필요
			result.Passed = false
			result.Resources = append(result.Resources,
				fmt.Sprintf("PV: %s (암호화 여부를 수동 확인 필요, CSI Driver: %s)", pv.Name, driver))
		}
	}

	return result
}
