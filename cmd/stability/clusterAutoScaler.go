package stability

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
func CheckClusterAutoscalerEnabled(client kubernetes.Interface) bool {
	// kube-system 네임스페이스의 모든 Deployment 목록 가져오기
	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	// 모든 Deployment를 순회하면서 컨테이너 이미지가 "cluster-autoscaler"를 포함하는지 확인
	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "cluster-autoscaler") {
				return true // Cluster Autoscaler가 실행 중
			}
		}
	}

	return false // Cluster Autoscaler가 없음
}

func CheckHpa(client kubernetes.Interface) bool {
	// 모든 Namespace Horizontal Pod Autoscaler 목록 확인인
	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(metav1.NamespaceAll).List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	return len(hpas.Items) > 0
}
