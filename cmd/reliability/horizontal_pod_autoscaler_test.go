package reliability_test

import (
	"context"
	"eks-checklist/cmd/reliability"
	"eks-checklist/cmd/testutils"
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckHpa(t *testing.T) {
	// YAML 파일 "hpa_check.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "horizontal_pod_autoscaler.yaml")

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		// deployments 필드 파싱
		deploymentsRaw, ok := tc["deployments"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'deployments' field", testName)
		}
		// hpas 필드 파싱
		hpasRaw, ok := tc["hpas"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'hpas' field", testName)
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 Deployment 객체 생성
			for _, d := range deploymentsRaw {
				dMap, ok := d.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': deployment is not a map", testName)
				}
				ns, ok := dMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': deployment missing 'namespace'", testName)
				}
				depName, ok := dMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': deployment missing 'name'", testName)
				}
				dep := &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      depName,
						Namespace: ns,
					},
				}
				_, err := client.AppsV1().Deployments(ns).Create(context.TODO(), dep, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create deployment %s/%s: %v", testName, ns, depName, err)
				}
			}

			// YAML에 정의된 HPA 객체 생성
			for _, h := range hpasRaw {
				hMap, ok := h.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': hpa is not a map", testName)
				}
				ns, ok := hMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': hpa missing 'namespace'", testName)
				}
				scaleTarget, ok := hMap["scale_target"].(string)
				if !ok {
					t.Fatalf("Test case '%s': hpa missing 'scale_target'", testName)
				}
				// HPA 객체 생성 시 Deployment를 ScaleTarget으로 지정
				hpa := &autoscalingv1.HorizontalPodAutoscaler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("hpa-for-%s", scaleTarget),
						Namespace: ns,
					},
					Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
						ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
							Kind:       "Deployment",
							Name:       scaleTarget,
							APIVersion: "apps/v1",
						},
					},
				}
				_, err := client.AutoscalingV1().HorizontalPodAutoscalers(ns).Create(context.TODO(), hpa, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create hpa in %s for target %s: %v", testName, ns, scaleTarget, err)
				}
			}

			// 함수 실행 및 결과 비교
			result := reliability.CheckHpa(client)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v\nFailureMsg: %s\nResources: %v", testName, expectPass, result.Passed, result.FailureMsg, result.Resources)
			} else {
				t.Logf("Test '%s' passed", testName)
			}
		})
	}
}
