package security_test

// import (
// 	"context"
// 	"testing"

// 	"eks-checklist/cmd/security"
// 	"eks-checklist/cmd/testutils"

// 	corev1 "k8s.io/api/core/v1"
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestPrintAccessControl(t *testing.T) {
// 	testCases := testutils.LoadTestCases(t, "access_control.yaml")

// 	for _, tc := range testCases {
// 		name := tc["name"].(string)
// 		expectFailure := tc["expect_failure"].(bool)

// 		t.Run(name, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			if name != "AwsAuth_Missing" {
// 				configMap := &corev1.ConfigMap{
// 					ObjectMeta: v1.ObjectMeta{
// 						Name:      "aws-auth",
// 						Namespace: "kube-system",
// 					},
// 					Data: map[string]string{},
// 				}

// 				if name == "AwsAuth_With_Roles_Users_Accounts" {
// 					configMap.Data["mapRoles"] = "- rolearn: arn:aws:iam::123456789012:role/EKSAdmin\n  username: admin"
// 					configMap.Data["mapUsers"] = "- userarn: arn:aws:iam::123456789012:user/john\n  username: john"
// 					configMap.Data["mapAccounts"] = "- 123456789012"
// 				}

// 				_, err := client.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), configMap, v1.CreateOptions{})
// 				if err != nil {
// 					t.Fatalf("ConfigMap 생성 실패: %v", err)
// 				}
// 			}

// 			result := security.PrintAccessControl(client, "mock-cluster")

// 			if result != !expectFailure {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
// 			}
// 		})
// 	}
// }
