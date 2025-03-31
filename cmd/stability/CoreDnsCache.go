// 변경 전 코드

// package stability

// import (
// 	"context"
// 	"log"
// 	"strings"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // CheckCoreDNSCache: CoreDNS의 캐시 적용 여부 검사
// func CheckCoreDNSCache(client kubernetes.Interface) bool {
// 	// kube-system 네임스페이스의 coredns ConfigMap 조회
// 	configMap, err := client.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "coredns", v1.GetOptions{})
// 	if err != nil {
// 		log.Printf("Failed to get CoreDNS ConfigMap: %v\n", err)
// 		return false
// 	}

// 	// Corefile 설정 가져오기
// 	corefile, ok := configMap.Data["Corefile"]
// 	if !ok {
// 		log.Println("CoreDNS Corefile not found in ConfigMap")
// 		return false
// 	}

// 	// Corefile에서 "cache" 플러그인 존재 여부 확인
// 	// 단순 있는걸로 판단할건지?? 고민 되긴함
// 	if strings.Contains(corefile, "cache") {
// 		return true
// 	}

// 	return false
// }

// 변경 후 코드
package stability

import (
	"context"
	"strings"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckCoreDNSCache checks whether the "cache" plugin is enabled in the CoreDNS Corefile.
func CheckCoreDNSCache(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "DNS 캐시 적용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	configMap, err := client.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "coredns", v1.GetOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	corefile, ok := configMap.Data["Corefile"]
	if !ok {
		result.Passed = false
		result.FailureMsg = "CoreDNS ConfigMap에 Corefile 항목이 존재하지 않습니다."
		return result
	}

	if strings.Contains(corefile, "cache") {
		result.Passed = true
		// result.SuccessMsg = "CoreDNS Corefile에 'cache' 플러그인이 설정되어 있습니다."
		// result.Resources = append(result.Resources, "ConfigMap: kube-system/coredns (cache plugin detected)")
	} else {
		result.Passed = false
		result.FailureMsg = "CoreDNS Corefile에 'cache' 플러그인이 설정되어 있지 않습니다."
		result.Resources = append(result.Resources, "ConfigMap: kube-system/coredns (cache plugin not found)")
	}

	return result
}
