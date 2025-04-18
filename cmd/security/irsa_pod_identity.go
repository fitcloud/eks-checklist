package security

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckIRSAAndPodIdentity(clientset kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "IRSA 또는 EKS Pod Identity 기반 권한 부여",
		Manual:    false,
		Passed:    true,
		// SuccessMsg: "IRSA 또는 EKS Pod Identity 기반 권한 부여",
		FailureMsg: "일부 서비스 계정이 IRSA 또는 EKS Pod Identity를 사용하지 않고 있습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	saList, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), v1.ListOptions{
		FieldSelector: "metadata.namespace!=kube-system", // kube-system 네임스페이스 제외
	})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	// IRSA 또는 Pod Identity를 사용하지 않는 Service Account 수집
	for _, sa := range saList.Items {
		annotations := sa.Annotations

		_, hasIRSA := annotations["eks.amazonaws.com/role-arn"]
		_, hasIdentity := annotations["eks.amazonaws.com/identity"]
		_, hasAudience := annotations["eks.amazonaws.com/audience"]

		if !(hasIRSA || hasIdentity || hasAudience) {
			result.Passed = false
			result.Resources = append(result.Resources, fmt.Sprintf("ServiceAccount: %s/%s", sa.Namespace, sa.Name))
		}
	}

	return result
}
