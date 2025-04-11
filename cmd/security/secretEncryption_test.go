package security_test

// import (
// 	"context"
// 	"encoding/base64"
// 	"testing"

// 	"eks-checklist/cmd/security"
// 	"eks-checklist/cmd/testutils"

// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestCheckSecretEncryption(t *testing.T) {
// 	testCases := testutils.LoadTestCases(t, "secret_encryption.yaml")

// 	for _, tc := range testCases {
// 		name := tc["name"].(string)
// 		expectPass := tc["expect_pass"].(bool)

// 		t.Run(name, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			// 시크릿 생성
// 			if secrets, ok := tc["secrets"].([]interface{}); ok {
// 				for _, s := range secrets {
// 					secretMap := s.(map[string]interface{})
// 					secretData := map[string][]byte{}

// 					if rawData, exists := secretMap["data"].(map[string]interface{}); exists {
// 						for key, value := range rawData {
// 							decoded, _ := base64.StdEncoding.DecodeString(value.(string))
// 							secretData[key] = decoded
// 						}
// 					}

// 					secret := &corev1.Secret{
// 						ObjectMeta: metav1.ObjectMeta{
// 							Name:      secretMap["name"].(string),
// 							Namespace: "default",
// 						},
// 						Data: secretData,
// 					}

// 					client.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
// 				}
// 			}

// 			result := security.CheckSecretEncryption(client)

// 			if result.Passed != expectPass {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectPass, result.Passed)
// 			}
// 		})
// 	}
// }
