package stability

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
func PodReplicaSetCheck(client kubernetes.Interface) bool {
	// ReplicaSet 목록 가져오기
	replicaSets, err := client.AppsV1().ReplicaSets("").List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	// 단일 Pod(복제본 1개)인 ReplicaSet 찾기
	var failedReplicaSets []string
	for _, rs := range replicaSets.Items {
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas == 1 {
			failedReplicaSets = append(failedReplicaSets, rs.Name)
		}
	}

	// replicatset 조회가 되면 false 반환 및 목록 출력
	if len(failedReplicaSets) > 0 {
		fmt.Println("The following ReplicaSets have only 1 pod:")
		for _, name := range failedReplicaSets {
			fmt.Println("-", name)
		}
		return false
	} else {
		return true
	}
}
