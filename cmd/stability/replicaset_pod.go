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
		if rs.Namespace == "kube-system" {
			continue // kube-system 네임스페이스는 검사 제외
		}
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas == 1 {
			result.Passed = false
			resource := fmt.Sprintf("Namespace: %s | ReplicaSet: %s (Replicas: 1)", rs.Namespace, rs.Name)
			result.Resources = append(result.Resources, resource)
		}
	}

	return result
}
