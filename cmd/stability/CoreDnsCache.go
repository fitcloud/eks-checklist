package stability

import (
	"context"
	"log"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckCoreDNSCache: CoreDNS의 캐시 적용 여부 검사
func CheckCoreDNSCache(client kubernetes.Interface) bool {
	// kube-system 네임스페이스의 coredns ConfigMap 조회
	configMap, err := client.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "coredns", v1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get CoreDNS ConfigMap: %v\n", err)
		return false
	}

	// Corefile 설정 가져오기
	corefile, ok := configMap.Data["Corefile"]
	if !ok {
		log.Println("CoreDNS Corefile not found in ConfigMap")
		return false
	}

	// Corefile에서 "cache" 플러그인 존재 여부 확인
	// 단순 있는걸로 판단할건지?? 고민 되긴함
	if strings.Contains(corefile, "cache") {
		return true
	}

	return false
}
