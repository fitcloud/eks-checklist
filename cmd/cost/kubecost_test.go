package cost_test

// import (
// 	"context"
// 	"testing"

// 	"eks-checklist/cmd/cost"
// 	"eks-checklist/cmd/testutils"

// 	appsv1 "k8s.io/api/apps/v1"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestGetKubecost(t *testing.T) {
// 	// testutils 패키지에서 테스트 케이스 로드
// 	testCases := testutils.LoadTestCases(t, "kubecost.yaml")

// 	for _, tc := range testCases {
// 		name := tc["name"].(string)
// 		expectFailure := tc["expect_failure"].(bool)

// 		t.Run(name, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			if !expectFailure {
// 				// kubecost deployment 생성
// 				_, err := client.AppsV1().Deployments("default").Create(context.TODO(), &appsv1.Deployment{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      "kubecost-deploy",
// 						Namespace: "default",
// 					},
// 					Spec: appsv1.DeploymentSpec{
// 						Selector: &metav1.LabelSelector{
// 							MatchLabels: map[string]string{
// 								"app": "kubecost",
// 							},
// 						},
// 						Template: corev1.PodTemplateSpec{
// 							ObjectMeta: metav1.ObjectMeta{
// 								Labels: map[string]string{
// 									"app": "kubecost",
// 								},
// 							},
// 							Spec: corev1.PodSpec{
// 								Containers: []corev1.Container{
// 									{
// 										Name:  "kubecost-container",
// 										Image: "gcr.io/kubecost1/kubecost-cost-analyzer:latest",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				}, metav1.CreateOptions{})
// 				if err != nil {
// 					t.Fatalf("failed to create kubecost deployment: %v", err)
// 				}
// 			}

// 			// 함수 실행
// 			result := cost.GetKubecost(client)

// 			// 기대값과 결과 비교
// 			if result != !expectFailure {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
// 			}
// 		})
// 	}
// }
