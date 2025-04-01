package security

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Secret 객체의 암호화 여부를 확인하는 함수
func CheckSecretEncryption(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Secret 객체 암호화",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 Secret 객체가 암호화되지 않은 채로 저장되어 있습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	secrets, err := client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	// Secret 객체가 없으면 PASS 처리 (암호화 우려 없음)
	if len(secrets.Items) == 0 {
		result.Passed = true
		return result
	}

	// 암호화 여부 판단
	for _, secret := range secrets.Items {
		for key, value := range secret.Data {
			if len(value) > 0 {
				// base64 인코딩된 데이터가 존재하는 경우 암호화 미적용으로 판단
				result.Passed = false
				resourceInfo := fmt.Sprintf("Namespace: %s | Secret: %s | Key: %s (base64 데이터 발견)", secret.Namespace, secret.Name, key)
				result.Resources = append(result.Resources, resourceInfo)
			}
		}
	}

	return result
}
