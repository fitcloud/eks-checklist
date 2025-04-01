package stability

import (
	"context"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckCoreDNSHpa checks if CoreDNS has an HPA set in the kube-system namespace.
func CheckCoreDNSHpa(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "CoreDNS에 HPA 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "CoreDNS에 Horizontal Pod Autoscaler(HPA)가 설정되어 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers("kube-system").List(context.TODO(), v1.ListOptions{
		FieldSelector: "metadata.name=coredns",
	})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	if len(hpas.Items) > 0 {
		result.Passed = true
		// result.SuccessMsg = "CoreDNS에 Horizontal Pod Autoscaler가 설정되어 있습니다."
		// result.Resources = append(result.Resources,
		// 	fmt.Sprintf("Namespace: %s | HPA: %s", hpa.Namespace, hpa.Name))
	} else {
		result.Passed = false
		result.FailureMsg = "CoreDNS에 Horizontal Pod Autoscaler(HPA)가 설정되어 있지 않습니다."
	}

	return result
}
