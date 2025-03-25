package stability

import (
	"context"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DaemonSet의 PriorityClass가 존재하는지 확인
func CheckDaemonSetPriorityClass(client kubernetes.Interface) bool {
	// 모든 namespace에 DaemonSet을 조회
	daemonset, err := client.AppsV1().DaemonSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// DaemonSet의 PriorityClass가 존재하는지 확인
	for _, ds := range daemonset.Items {
		if ds.Spec.Template.Spec.PriorityClassName == "" {
			// PriorityClass가 존재하지 않는 deamonset을 출력
			log.Printf("DaemonSet %s does not have a PriorityClass\n", ds.Name)
			return false
		}
	}
	return true
}
