package stability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckHpa checks whether Deployments are using Horizontal Pod Autoscaler (HPA).
func CheckHpa(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "HPA 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 Deployment에 HPA가 적용되어 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	// 모든 Deployment 조회
	deployments, err := client.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	// 모든 HPA 조회
	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	// HPA 적용 여부를 확인
	hpaTargets := make(map[string]bool)
	for _, hpa := range hpas.Items {
		key := fmt.Sprintf("%s/%s", hpa.Namespace, hpa.Spec.ScaleTargetRef.Name)
		hpaTargets[key] = true
	}

	var withoutHPA []string

	for _, deployment := range deployments.Items {
		key := fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name)
		if !hpaTargets[key] {
			result.Passed = false
			withoutHPA = append(withoutHPA, key)
			result.Resources = append(result.Resources, fmt.Sprintf("Deployment: %s (HPA 미설정)", key))
		}
	}

	return result
}
