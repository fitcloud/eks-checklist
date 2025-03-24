package stability

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckProbe - 모든 Pod을 검색하여 startProbe, livenessProbe, readinessProbe 가 모두 설정되었는지 확인
func CheckProbe(client kubernetes.Interface) bool {
	// 모든 Pod을 조회하되 kube-system 네임스페이스는 제외
	// 왜냐하면 시스템 애드온 파드들은 기본적으로 몇개씩 없음 어떻게 할지는 추후 더 고민
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{
		FieldSelector: "metadata.namespace!=kube-system", // kube-system 네임스페이스 제외
	})
	if err != nil {
		fmt.Println("Error retrieving pods:", err)
		return false
	}

	// 모든 pod가 프로브를 가있는지 확인하는 변수
	allProbesSet := true

	// 각 Pod의 컨테이너에서 Probe 설정을 확인
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			missingProbes := []string{}

			// 각 프로브가 설정되어 있지 않으면 해당 프로브를 missingProbes 배열에 추가
			if container.StartupProbe == nil {
				missingProbes = append(missingProbes, "startupProbe")
			}
			if container.LivenessProbe == nil {
				missingProbes = append(missingProbes, "livenessProbe")
			}
			if container.ReadinessProbe == nil {
				missingProbes = append(missingProbes, "readinessProbe")
			}

			// 빠진 프로브가 하나라도 있으면 해당 Pod 이름과 빠진 프로브들을 출력
			if len(missingProbes) > 0 {
				fmt.Printf("Pod %s is missing the following probes: %v\n", pod.Name, missingProbes)
				// 프로브가 빠져 있으면 allProbesSet을 false로 설정
				allProbesSet = false
			}
		}
	}

	// 모든 Pod이 프로브를 설정했으면 true 반환, 아니면 false 반환
	return allProbesSet
}
