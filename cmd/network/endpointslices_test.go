package network_test

import (
	"context"
	"strconv"
	"testing"

	"eks-checklist/cmd/network"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEndpointSlicesCheck_YAML(t *testing.T) {
	// 생성한 YAML 파일 "endpointslices.yaml"로부터 테스트 케이스를 로드합니다.
	testCases := testutils.LoadTestCases(t, "endpointslices.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)
		endpoints := tc["endpoints"].([]interface{})
		endpointSlices := tc["endpointSlices"].([]interface{})

		t.Run(testName, func(t *testing.T) {
			// Fake Kubernetes 클라이언트 생성
			client := fake.NewSimpleClientset()

			// YAML에 정의된 Endpoints 객체 생성
			for _, ep := range endpoints {
				epDef := ep.(map[string]interface{})
				ns := epDef["namespace"].(string)
				name := epDef["name"].(string)
				epObj := &corev1.Endpoints{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns,
					},
				}
				_, err := client.CoreV1().Endpoints(ns).Create(context.TODO(), epObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Failed to create Endpoints object: %v", err)
				}
			}

			// YAML에 정의된 EndpointSlices 객체 생성
			for idx, es := range endpointSlices {
				esDef := es.(map[string]interface{})
				ns := esDef["namespace"].(string)
				serviceName := esDef["service_name"].(string)
				esObj := &discoveryv1.EndpointSlice{
					ObjectMeta: metav1.ObjectMeta{
						// strconv.Itoa(idx)를 사용하여 int를 문자열로 변환
						Name:      "slice-" + serviceName + "-" + strconv.Itoa(idx),
						Namespace: ns,
						Labels: map[string]string{
							"kubernetes.io/service-name": serviceName,
						},
					},
				}
				_, err := client.DiscoveryV1().EndpointSlices(ns).Create(context.TODO(), esObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Failed to create EndpointSlice object: %v", err)
				}
			}

			// EndpointSlicesCheck 함수 실행 후 반환값 검증
			result := network.EndpointSlicesCheck(client)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}
