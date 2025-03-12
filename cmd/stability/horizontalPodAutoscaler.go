package stability

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 클러스터에 Horizontal Pod Autoscaler가 설정정되어 있는지 확인
func CheckHpa(client kubernetes.Interface) bool {
	// 모든 Namespace Horizontal Pod Autoscaler 목록 확인인
	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(metav1.NamespaceAll).List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	return len(hpas.Items) > 0
}
