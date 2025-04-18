package stability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckDaemonSetPriorityClass checks if all DaemonSets have a PriorityClass assigned.
func CheckDaemonSetPriorityClass(karpenter_installed common.CheckResult, client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Karpenter 사용시 DaemonSet에 Priority Class 부여",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 DaemonSet에 PriorityClass가 설정되어 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	if !karpenter_installed.Passed {
		result.Passed = false
		result.FailureMsg = "Karpenter가 설치되어 있지 않습니다."
		return result
	}

	daemonSets, err := client.AppsV1().DaemonSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	hasMissing := false
	for _, ds := range daemonSets.Items {
		if ds.Spec.Template.Spec.PriorityClassName == "" {
			hasMissing = true
			result.Resources = append(result.Resources,
				fmt.Sprintf("Namespace: %s | DaemonSet: %s (PriorityClass 미설정)", ds.Namespace, ds.Name))
		}
	}

	if hasMissing {
		result.Passed = false
		result.FailureMsg = "일부 DaemonSet에 PriorityClass가 설정되어 있지 않습니다."
	} else {
		result.Passed = true
		// result.SuccessMsg = "모든 DaemonSet에 PriorityClass가 설정되어 있습니다."
	}

	return result
}
