package scalability_test

// import (
// 	"context"
// 	"testing"

// 	"eks-checklist/cmd/scalability"
// 	"eks-checklist/cmd/testutils"

// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestCheckSpotNodeTerminationHandler(t *testing.T) {
// 	testCases := testutils.LoadTestCases(t, "spot_handler.yaml")

// 	for _, tc := range testCases {
// 		name := tc["name"].(string)
// 		expectFailure := tc["expect_failure"].(bool)

// 		t.Run(name, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			if !expectFailure {
// 				// termination-handler 파드 생성
// 				_, err := client.CoreV1().Pods("default").Create(context.TODO(), &corev1.Pod{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      "spot-termination-handler-abc",
// 						Namespace: "default",
// 					},
// 					Spec: corev1.PodSpec{
// 						Containers: []corev1.Container{
// 							{
// 								Name:  "termination-handler",
// 								Image: "amazon/spot-termination-handler:latest",
// 							},
// 						},
// 					},
// 				}, metav1.CreateOptions{})
// 				if err != nil {
// 					t.Fatalf("failed to create termination-handler pod: %v", err)
// 				}
// 			}

// 			// 함수 실행
// 			result := scalability.CheckSpotNodeTerminationHandler(client)

// 			if result != !expectFailure {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
// 			}
// 		})
// 	}
// }
