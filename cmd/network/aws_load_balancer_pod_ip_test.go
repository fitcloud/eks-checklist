package network_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/network"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckAwsLoadBalancerPodIp_YAML(t *testing.T) {
	// 생성한 YAML 파일 "aws_load_balancer_pod_ip.yaml"에서 테스트 케이스를 로드합니다.
	testCases := testutils.LoadTestCases(t, "aws_load_balancer_pod_ip.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)
		ingressesRaw := tc["ingresses"].([]interface{})
		servicesRaw := tc["services"].([]interface{})

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 Ingress 객체 생성
			for _, ig := range ingressesRaw {
				igDef := ig.(map[string]interface{})
				ns := igDef["namespace"].(string)
				name := igDef["name"].(string)
				annotations := make(map[string]string)
				if rawAnn, ok := igDef["annotations"].(map[string]interface{}); ok {
					for key, value := range rawAnn {
						annotations[key] = value.(string)
					}
				}

				ingObj := &netv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:        name,
						Namespace:   ns,
						Annotations: annotations,
					},
				}
				_, err := client.NetworkingV1().Ingresses(ns).Create(context.TODO(), ingObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Failed to create Ingress %s/%s: %v", ns, name, err)
				}
			}

			// YAML에 정의된 Service 객체 생성
			for _, svcRaw := range servicesRaw {
				svcDef := svcRaw.(map[string]interface{})
				ns := svcDef["namespace"].(string)
				name := svcDef["name"].(string)
				annotations := make(map[string]string)
				if rawAnn, ok := svcDef["annotations"].(map[string]interface{}); ok {
					for key, value := range rawAnn {
						annotations[key] = value.(string)
					}
				}
				var ownerRefs []metav1.OwnerReference
				if rawOwnerRefs, ok := svcDef["ownerReferences"].([]interface{}); ok {
					for _, or := range rawOwnerRefs {
						orMap := or.(map[string]interface{})
						ownerRefs = append(ownerRefs, metav1.OwnerReference{
							Kind: orMap["kind"].(string),
							Name: orMap["name"].(string),
							// UID는 테스트에서는 필요 없으므로 빈 문자열로 처리
							UID: "",
						})
					}
				}
				svcObj := &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:            name,
						Namespace:       ns,
						Annotations:     annotations,
						OwnerReferences: ownerRefs,
					},
				}
				_, err := client.CoreV1().Services(ns).Create(context.TODO(), svcObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Failed to create Service %s/%s: %v", ns, name, err)
				}
			}

			// CheckAwsLoadBalancerPodIp 함수 실행 및 반환값 검증
			result := network.CheckAwsLoadBalancerPodIp(client)
			if result != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result)
			}
		})
	}
}
