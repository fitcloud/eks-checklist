package security

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckPodToPodNetworkPolicy는 네트워크 정책이 존재하는지 확인하고, selector와 ingress/egress의 정보를 JSON 형식으로 출력합니다.
func CheckPodToPodNetworkPolicy(client kubernetes.Interface) bool {
	// 모든 네임스페이스의 NetworkPolicy를 조회
	npList, err := client.NetworkingV1().NetworkPolicies("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("네트워크 정책을 가져오는 중 오류 발생: %v\n", err)
		return false
	}

	// NetworkPolicy가 하나라도 존재하면 접근 제어가 설정되었다고 판단
	if len(npList.Items) > 0 {
		fmt.Println("발견된 Network Policies:")
		for _, np := range npList.Items {
			// NetworkPolicy의 podSelector 출력
			npOutput := map[string]interface{}{
				"NetworkPolicyName": np.Name,
				"PodSelector": map[string]interface{}{
					"MatchLabels":      np.Spec.PodSelector.MatchLabels,
					"MatchExpressions": np.Spec.PodSelector.MatchExpressions,
				},
				"Ingress": []interface{}{},
				"Egress":  []interface{}{},
			}

			// Ingress 규칙에 출력
			for _, ingress := range np.Spec.Ingress {
				if len(ingress.From) > 0 {
					ingressFrom := []map[string]interface{}{}
					for _, from := range ingress.From {
						ingressFrom = append(ingressFrom, map[string]interface{}{
							"PodSelector":       from.PodSelector,
							"NamespaceSelector": from.NamespaceSelector,
						})
					}
					npOutput["Ingress"] = ingressFrom
				}
			}

			// Egress 규칙에 출력
			for _, egress := range np.Spec.Egress {
				if len(egress.To) > 0 {
					egressTo := []map[string]interface{}{}
					for _, to := range egress.To {
						egressTo = append(egressTo, map[string]interface{}{
							"PodSelector":       to.PodSelector,
							"NamespaceSelector": to.NamespaceSelector,
						})
					}
					npOutput["Egress"] = egressTo
				}
			}

			// JSON 형식으로 출력
			npJSON, err := json.MarshalIndent(npOutput, "", "  ")
			if err != nil {
				fmt.Printf("JSON 출력 오류: %v\n", err)
			} else {
				fmt.Println(string(npJSON))
			}
		}
	}
	// NetworkPolicy가 하나라도 존재하면 접근 제어가 설정되었다고 판단
	return len(npList.Items) > 0
}
