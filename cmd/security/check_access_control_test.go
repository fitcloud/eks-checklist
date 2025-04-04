package security_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	"github.com/aws/aws-sdk-go-v2/aws"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckAccessControl(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "access_control.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// aws-auth ConfigMap 생성: AwsAuth_Missing 케이스가 아닌 경우 생성
			if name != "AwsAuth_Missing" {
				configMap := &corev1.ConfigMap{
					ObjectMeta: v1.ObjectMeta{
						Name:      "aws-auth",
						Namespace: "kube-system",
					},
					Data: map[string]string{},
				}

				// 특정 케이스에 대해 추가 데이터를 설정
				if name == "AwsAuth_With_Roles_Users_Accounts" {
					configMap.Data["mapRoles"] = "- rolearn: arn:aws:iam::123456789012:role/EKSAdmin\n  username: admin"
					configMap.Data["mapUsers"] = "- userarn: arn:aws:iam::123456789012:user/john\n  username: john"
					configMap.Data["mapAccounts"] = "- 123456789012"
				}

				_, err := client.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), configMap, v1.CreateOptions{})
				if err != nil {
					t.Fatalf("ConfigMap 생성 실패: %v", err)
				}
			}

			// Dummy aws.Config 생성 (테스트에서는 별도의 설정 없이 사용)
			var cfg aws.Config

			// 실제 호출: security.PrintAccessControl이 아니라 CheckAccessControl을 사용합니다.
			result := security.CheckAccessControl(client, cfg, "mock-cluster")

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, expectPass, result.Passed)
			}
		})
	}
}
