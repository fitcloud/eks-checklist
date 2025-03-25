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

	// Access Entries 출력
	fmt.Println("\nAccess Entries:")

	// mapRoles 항목 출력
	if roles, exists := configMap.Data["mapRoles"]; exists {
		fmt.Println("\n- IAM Roles:")
		for _, role := range strings.Split(roles, "\n") {
			if role != "" {
				fmt.Printf("  - %s\n", role)
			}
		}
	} else {
		fmt.Println("\n- IAM Roles: 없음")
	}

	// mapUsers 항목 출력
	if users, exists := configMap.Data["mapUsers"]; exists {
		fmt.Println("\n- IAM Users:")
		for _, user := range strings.Split(users, "\n") {
			if user != "" {
				fmt.Printf("  - %s\n", user)
			}
		}
	} else {
		fmt.Println("\n- IAM Users: 없음")
	}

	// mapAccounts (AWS 계정 기반 액세스) 항목 출력
	if accounts, exists := configMap.Data["mapAccounts"]; exists {
		fmt.Println("\n- AWS Accounts:")
		for _, account := range strings.Split(accounts, "\n") {
			if account != "" {
				fmt.Printf("  - %s\n", account)
			}
		}
	} else {
		fmt.Println("\n- AWS Accounts: 없음")
	}

	return true
}
