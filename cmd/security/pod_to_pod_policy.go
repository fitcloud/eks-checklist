package security

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckPodToPodNetworkPolicy checks whether NetworkPolicies exist for pod-to-pod communication.
func CheckPodToPodNetworkPolicy(client kubernetes.Interface, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[SEC-009] Pod-to-Pod 접근 제어",
		Manual:     true,
		Passed:     true,
		FailureMsg: "Pod 간 접근 제어를 위한 NetworkPolicy가 설정되어 있지만 정책이 적합하게 설정되어 있는지 수동으로 확인해야합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/security/SEC-009",
	}

	// 1. NetworkPolicy 목록 조회
	npList, err := client.NetworkingV1().NetworkPolicies("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	if len(npList.Items) == 0 {
		result.Passed = false
		result.FailureMsg = "Pod 간 접근 제어를 위한 NetworkPolicy가 존재하지 않습니다."
		return result
	}

	// 2. 결과 저장 디렉토리 생성
	baseDir := filepath.Join(".", "output", eksCluster+"-pod-network-policy")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.Passed = false
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	// 3. 각 NetworkPolicy를 YAML로 저장
	for _, np := range npList.Items {
		filename := fmt.Sprintf("%s-%s.yaml", np.Namespace, np.Name)
		filePath := filepath.Join(baseDir, filename)

		err := common.SaveK8sResourceAsYAML(&np, filePath)
		if err != nil {
			result.Resources = append(result.Resources, fmt.Sprintf("저장 실패: %s/%s (%v)", np.Namespace, np.Name, err))
		} else {
			result.Resources = append(result.Resources, fmt.Sprintf("저장됨: %s", filePath))
		}
	}

	return result
}
