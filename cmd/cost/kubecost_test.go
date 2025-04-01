package cost_test

import (
	"context"
	"strings"
	"testing"

	"eks-checklist/cmd/cost"
	"eks-checklist/cmd/testutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetKubecost(t *testing.T) {
	// YAML 파일 "kubecost.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "kubecost.yaml")

	for _, tc := range testCases {

		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		// YAML 파일에서 추가로 node_name, node_labels 읽기
		nodeName, _ := tc["node_name"].(string)
		rawLabels, _ := tc["node_labels"].([]interface{})
		labels := make(map[string]string)
		for _, l := range rawLabels {
			// 예: "app: kubecost" 형식의 문자열을 ":"로 분리
			parts := strings.SplitN(l.(string), ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				labels[key] = value
			}
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// expect_pass 값에 따라 올바른 이미지 할당:
			// expect_pass가 true이면 Kubecost 이미지가 있어야 함.
			// false이면 Kubecost 이미지가 없는 다른 이미지를 사용.
			var containerImage string
			if expectPass {
				containerImage = "gcr.io/kubecost1/kubecost-cost-analyzer:latest"
			} else {
				containerImage = "nginx:1.21"
			}

			// Deployment 이름은 node_name 값에 따라 설정 (없으면 기본 이름 사용)
			deployName := "kubecost-deploy"
			if nodeName != "" {
				deployName = nodeName + "-deploy"
			}

			// YAML 파일에서 읽은 라벨을 Pod의 라벨로 사용
			_, err := client.AppsV1().Deployments("default").Create(context.TODO(), &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deployName,
					Namespace: "default",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: labels,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "main-container",
									Image: containerImage,
								},
							},
						},
					},
				},
			}, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("failed to create deployment: %v", err)
			}

			// 함수 실행
			result := cost.GetKubecost(client)

			// 결과 검증: result.Passed가 expect_pass와 동일해야 함.
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}
