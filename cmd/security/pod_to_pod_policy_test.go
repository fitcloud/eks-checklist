package security_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckPodToPodNetworkPolicy(t *testing.T) {
	// YAML 파일 "pod_to_pod_policy.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "pod_to_pod_policy.yaml")
	// 예시로 eksCluster 이름을 "test-cluster"로 사용합니다.
	eksClusterName := "test-cluster"
	t.Logf("Loaded %d test cases", len(testCases))

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		npListRaw, ok := tc["networkpolicies"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'networkpolicies' field", testName)
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 각 NetworkPolicy 객체 생성
			for _, npRaw := range npListRaw {
				npMap, ok := npRaw.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': networkpolicy is not a map", testName)
				}
				ns, ok := npMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': networkpolicy missing 'namespace'", testName)
				}
				name, ok := npMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': networkpolicy missing 'name'", testName)
				}

				// podSelector (matchLabels만 처리)
				selector := metav1.LabelSelector{}
				if psRaw, exists := npMap["podSelector"]; exists {
					if psMap, ok := psRaw.(map[string]interface{}); ok {
						if mlRaw, exists := psMap["matchLabels"]; exists {
							if ml, ok := mlRaw.(map[string]interface{}); ok {
								selector.MatchLabels = make(map[string]string)
								for key, val := range ml {
									selector.MatchLabels[key] = val.(string)
								}
							}
						}
					}
				}

				// Ingress 규칙 (간단하게 from 항목만 처리)
				var ingressRules []networkingv1.NetworkPolicyIngressRule
				if ingRaw, exists := npMap["ingress"]; exists {
					if ingList, ok := ingRaw.([]interface{}); ok {
						for _, ingItem := range ingList {
							ingMap, ok := ingItem.(map[string]interface{})
							if !ok {
								continue
							}
							var fromList []networkingv1.NetworkPolicyPeer
							if fromRaw, exists := ingMap["from"]; exists {
								if fromArr, ok := fromRaw.([]interface{}); ok {
									for _, f := range fromArr {
										fMap, ok := f.(map[string]interface{})
										if !ok {
											continue
										}
										peer := networkingv1.NetworkPolicyPeer{}
										if psRaw, exists := fMap["podSelector"]; exists {
											if psMap, ok := psRaw.(map[string]interface{}); ok {
												peer.PodSelector = &metav1.LabelSelector{
													MatchLabels: make(map[string]string),
												}
												if mlRaw, exists := psMap["matchLabels"]; exists {
													if ml, ok := mlRaw.(map[string]interface{}); ok {
														for key, val := range ml {
															peer.PodSelector.MatchLabels[key] = val.(string)
														}
													}
												}
											}
										}
										if nsRaw, exists := fMap["namespaceSelector"]; exists {
											if nsMap, ok := nsRaw.(map[string]interface{}); ok {
												peer.NamespaceSelector = &metav1.LabelSelector{
													MatchLabels: make(map[string]string),
												}
												if mlRaw, exists := nsMap["matchLabels"]; exists {
													if ml, ok := mlRaw.(map[string]interface{}); ok {
														for key, val := range ml {
															peer.NamespaceSelector.MatchLabels[key] = val.(string)
														}
													}
												}
											}
										}
										fromList = append(fromList, peer)
									}
								}
							}
							ingressRules = append(ingressRules, networkingv1.NetworkPolicyIngressRule{
								From: fromList,
							})
						}
					}
				}

				// Egress 규칙은 테스트 단순화를 위해 빈 배열 사용
				npObj := &networkingv1.NetworkPolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns,
					},
					Spec: networkingv1.NetworkPolicySpec{
						PodSelector: selector,
						Ingress:     ingressRules,
						Egress:      []networkingv1.NetworkPolicyEgressRule{},
					},
				}

				_, err := client.NetworkingV1().NetworkPolicies(ns).Create(context.TODO(), npObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create NetworkPolicy %s/%s: %v", testName, ns, name, err)
				}
			}

			// CheckPodToPodNetworkPolicy 함수 실행 (eksCluster 이름을 인자로 전달)
			result := security.CheckPodToPodNetworkPolicy(client, eksClusterName)

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}
