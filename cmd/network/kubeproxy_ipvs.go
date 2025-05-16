package network

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckKubeProxyIPVSMode checks whether kube-proxy is set to use IPVS mode.
func CheckKubeProxyIPVSMode(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "[NET-008] kube-proxy에 IPVS 모드 적용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://fitcloud.github.io/eks-checklist/runbook/network/NET-008",
	}

	configMap, err := client.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "kube-proxy-config", metav1.GetOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	configContent, exists := configMap.Data["config"]
	if !exists {
		result.Passed = false
		result.FailureMsg = "kube-proxy ConfigMap에 'config' 필드가 존재하지 않습니다."
		return result
	}

	if strings.Contains(configContent, `mode: "ipvs"`) {
		result.Passed = true
		// result.SuccessMsg = "kube-proxy가 IPVS 모드로 설정되어 있습니다."
		// result.Resources = append(result.Resources, "ConfigMap: kube-system/kube-proxy-config | mode: ipvs")
		return result
	}

	// IPVS가 아닌 모드 출력
	modeValue := "알 수 없음"
	for _, line := range strings.Split(configContent, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mode:") {
			modeValue = strings.TrimSpace(strings.TrimPrefix(line, "mode:"))
			break
		}
	}

	result.Passed = false
	result.FailureMsg = fmt.Sprintf("kube-proxy가 IPVS 모드로 설정되어 있지 않습니다. 현재 모드: %s", modeValue)
	result.Resources = append(result.Resources, fmt.Sprintf("ConfigMap: kube-system/kube-proxy-config | mode: %s", modeValue))

	return result
}
