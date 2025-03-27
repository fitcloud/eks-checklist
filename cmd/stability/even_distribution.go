package stability

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	Red    = "\033[31m" // 빨간색
	Green  = "\033[32m" // 초록색
	Yellow = "\033[33m" // 노란색
	Reset  = "\033[0m"  // 기본 색상으로 리셋
)

// CheckPodDistributionAndAffinity 함수는 클러스터 내의 모든 Pod에 대해
// affinity 또는 topologySpreadConstraints 설정 여부를 확인하고,
// 둘 중 하나라도 적절하게 설정되어 있는지를 검사한다.
func CheckPodDistributionAndAffinity(clientset kubernetes.Interface) {
	// 모든 네임스페이스에서 Pod 목록을 가져온다.
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		fmt.Printf("FAIL : 클러스터에서 Pod 목록을 가져오는 데 실패했습니다.\n에러: %v\n", err)
		return
	}

	violations := []string{}

	for _, pod := range pods.Items {
		// affinity 설정이 존재하는지 확인
		affinityExists := pod.Spec.Affinity != nil
		topologyValid := false

		// topologySpreadConstraints 설정 여부와 maxSkew 유효성 확인
		if len(pod.Spec.TopologySpreadConstraints) > 0 {
			topologyValid = true
			for _, constraint := range pod.Spec.TopologySpreadConstraints {
				if constraint.MaxSkew > 1 {
					topologyValid = false
					violations = append(violations, fmt.Sprintf("- Pod %s (namespace: %s) - maxSkew 값이 %d (1 초과)", pod.Name, pod.Namespace, constraint.MaxSkew))
				}
			}
		}

		// affinity도 없고 유효한 topologySpreadConstraints도 없는 경우 위반사항으로 기록
		if !affinityExists && !topologyValid {
			violations = append(violations, fmt.Sprintf("- Pod %s (namespace: %s) - affinity와 유효한 topologySpreadConstraints 설정이 모두 없음", pod.Name, pod.Namespace))
		}
	}

	// 위반 사항 여부에 따라 결과 출력
	if len(violations) == 0 {
		fmt.Println(Green + "PASS: All pods are evenly distributed or have valid affinity settings." + Reset)
	} else {
		fmt.Println(Red + "FAIL: Some pods are not evenly distributed and have no valid affinity settings." + Reset)
		fmt.Println("Affected resources:")
		for _, warn := range violations {
			fmt.Println(warn)
		}
		fmt.Println("Runbook URL: https://your.runbook.url/irsa-or-pod-identity")
	}
}
