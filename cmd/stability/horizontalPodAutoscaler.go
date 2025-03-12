package stability

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckHpa(client kubernetes.Interface) bool {
	// 모든 Namespace Horizontal Pod Autoscaler 목록 확인인
	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(metav1.NamespaceAll).List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	return len(hpas.Items) > 0
}
