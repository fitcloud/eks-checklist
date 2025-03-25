package security

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Secret을 확인 후, Secret에 대한 데이터가 그냥 바로 가져와지면 Base64인거임
// 원래는 EncryptionConfiguration 이거를 검사해야할 것 같은데 아직 보류
func CheckSecretEncryption(client kubernetes.Interface) bool {
	secrets, err := client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Error fetching secrets: %v", err)
		return false
	}

	// 시크릿이 없는 경우
	if len(secrets.Items) == 0 {
		log.Println("secrets not found")
		return true
	}

	// 각 Secret의 데이터가 base64로 되어 있는지 확인
	for _, secret := range secrets.Items {
		for key, value := range secret.Data {
			// base64로 인코딩된 값이 존재하면 암호화가 적용되지 않았다고 판단
			if len(value) > 0 {
				log.Printf("Secret %s contains base64 data for key %s\n", secret.Name, key)
				return false
			}
		}
	}

	return true
}
