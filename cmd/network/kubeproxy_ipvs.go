// 변경 전 코드

// package network

// import (
// 	"context"
// 	"log"
// 	"strings"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // CheckKubeProxyIPVSMode - kube-proxy의 ConfigMap에서 IPVS 모드 적용 여부와 현재 모드 출력
// func CheckKubeProxyIPVSMode(client kubernetes.Interface) bool {
// 	// kube-system 네임스페이스에서 kube-proxy의 ConfigMap 조회
// 	configMap, err := client.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "kube-proxy-config", metav1.GetOptions{})
// 	if err != nil {
// 		log.Println("Error fetching kube-proxy config map:", err)
// 		return false
// 	}

// 	// configMap에서 mode 값을 가져오기
// 	mode, exists := configMap.Data["config"]
// 	if !exists {
// 		log.Println("No config found in the kube-proxy config map.")
// 		return false
// 	}

// 	// config에서 mode 값이 "ipvs"로 설정되어 있는지 확인
// 	if strings.Contains(mode, "mode: \"ipvs\"") {
// 		log.Println("kube-proxy mode is set to IPVS.")
// 		return true
// 	} else {
// 		// IPVS가 아니면 현재 설정된 mode를 추출하여 출력
// 		lines := strings.Split(mode, "\n")
// 		for _, line := range lines {
// 			if strings.HasPrefix(line, "mode:") {
// 				// "mode: ipvs" 처럼 나올 때 값을 추출
// 				parts := strings.Fields(line)
// 				if len(parts) > 1 {
// 					log.Printf("kube-proxy mode is not set to IPVS. Current mode is: %s\n", parts[1])
// 					return false
// 				}
// 			}
// 		}
// 	}

// 	// "mode"가 없으면 "Unknown" 출력
// 	log.Println("kube-proxy mode is not set to IPVS. Current mode is: Unknown")
// 	return false
// }

// 변경 후 코드
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
		CheckName: "kube-proxy에 IPVS 모드 적용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
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
