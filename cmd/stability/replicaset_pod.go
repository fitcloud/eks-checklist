// 변경 전 코드

// package stability

// import (
// 	"context"
// 	"fmt"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
// func PodReplicaSetCheck(client kubernetes.Interface) bool {
// 	// ReplicaSet 목록 가져오기
// 	replicaSets, err := client.AppsV1().ReplicaSets("").List(context.TODO(), v1.ListOptions{})

// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// 단일 Pod(복제본 1개)인 ReplicaSet 찾기
// 	var failedReplicaSets []string
// 	for _, rs := range replicaSets.Items {
// 		if rs.Spec.Replicas != nil && *rs.Spec.Replicas == 1 {
// 			failedReplicaSets = append(failedReplicaSets, rs.Name)
// 		}
// 	}

// 	// replicatset 조회가 되면 false 반환 및 목록 출력
// 	if len(failedReplicaSets) > 0 {
// 		fmt.Println("The following ReplicaSets have only 1 pod:")
// 		for _, name := range failedReplicaSets {
// 			fmt.Println("-", name)
// 		}
// 		return false
// 	} else {
// 		return true
// 	}
// }

package stability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodReplicaSetCheck checks that ReplicaSets are configured with more than 1 pod (replica).
func PodReplicaSetCheck(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "2개 이상의 Pod 복제본 사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 ReplicaSet이 복제본을 1개만 사용하고 있습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	replicaSets, err := client.AppsV1().ReplicaSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, rs := range replicaSets.Items {
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas == 1 {
			result.Passed = false
			resource := fmt.Sprintf("Namespace: %s | ReplicaSet: %s (Replicas: 1)", rs.Namespace, rs.Name)
			result.Resources = append(result.Resources, resource)
		}
	}

	return result
}
