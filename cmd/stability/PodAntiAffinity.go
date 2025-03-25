package stability

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 모든 Deployment의 Anti-Affinity 설정을 검사, topologyKey로 kubernetes.io/hostname이 사용되었는지 확인합니다.
func CheckDeploymentAntiAffinity(client kubernetes.Interface) bool {
	// 모든 네임스페이스의 Deployment를 조회
	deploymentList, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Deployment 리스트를 가져오는 중 오류 발생: %v\n", err)
		return false
	}

	// Deployment마다 Anti-Affinity 설정이 있는지 확인
	for _, deployment := range deploymentList.Items {
		if deployment.Spec.Template.Spec.Affinity != nil && deployment.Spec.Template.Spec.Affinity.PodAntiAffinity != nil {
			// RequiredDuringSchedulingIgnoredDuringExecution에서 topologyKey가 kubernetes.io/hostname인지 확인
			for _, rule := range deployment.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
				if rule.TopologyKey == "kubernetes.io/hostname" {
					return true
				}
			}

			// PreferredDuringSchedulingIgnoredDuringExecution에서 topologyKey가 kubernetes.io/hostname인지 확인
			for _, rule := range deployment.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
				if rule.PodAffinityTerm.TopologyKey == "kubernetes.io/hostname" {
					return true
				}
			}
		}
	}

	// Anti-Affinity가 설정되지 않았거나, topologyKey가 kubernetes.io/hostname이 아닌 경우 false 반환
	return false
}
