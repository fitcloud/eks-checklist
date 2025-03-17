package security

import (
	"context"
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PrintAccessControl는 aws-auth ConfigMap을 확인하고 EKS 클러스터 정보를 출력합니다.
func PrintAccessControl(client kubernetes.Interface, eksCluster string) bool {
	// 'aws-auth' ConfigMap을 'kube-system' 네임스페이스에서 가져옵니다.
	configMapClient := client.CoreV1().ConfigMaps("kube-system")
	configMap, err := configMapClient.Get(context.TODO(), "aws-auth", v1.GetOptions{})
	if err != nil {
		log.Printf("aws-auth ConfigMap을 가져오는 중 오류 발생: %v", err)
		return false
	}

	// aws-auth ConfigMap의 내용을 출력합니다.
	fmt.Println("aws-auth ConfigMap:")
	for key, value := range configMap.Data {
		fmt.Printf("%s: %s\n", key, value)
	}

	// mapRoles 항목이 있는 경우 출력
	if roles, exists := configMap.Data["mapRoles"]; exists {
		fmt.Println("\nmapRoles:")
		// 각 역할을 줄바꿈으로 분리하여 출력
		for _, role := range strings.Split(roles, "\n") {
			if role != "" {
				fmt.Printf("- %s\n", role)
			}
		}
	} else {
		fmt.Println("\nmapRoles 항목이 없습니다.")
	}

	// mapUsers 항목이 있는 경우 출력
	if users, exists := configMap.Data["mapUsers"]; exists {
		fmt.Println("\nmapUsers:")
		// 각 사용자를 줄바꿈으로 분리하여 출력
		for _, user := range strings.Split(users, "\n") {
			if user != "" {
				fmt.Printf("- %s\n", user)
			}
		}
	} else {
		fmt.Println("\nmapUsers 항목이 없습니다.")
	}

	return true
}
