package stability

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckClusterAutoscalerEnabled checks whether the Cluster Autoscaler is deployed in the cluster.
func CheckClusterAutoscalerEnabled(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Cluster Autoscaler 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Cluster Autoscaler가 설치되어 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	deployments, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, deploy := range deployments.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "cluster-autoscaler") {
				result.Passed = true
				// result.SuccessMsg = fmt.Sprintf("Deployment '%s/%s'에 Cluster Autoscaler가 설치되어 있습니다.", deploy.Namespace, deploy.Name)
				result.Resources = append(result.Resources,
					fmt.Sprintf("Namespace: %s | Deployment: %s | Image: %s", deploy.Namespace, deploy.Name, container.Image))
				return result
			}
		}
	}

	result.Passed = false

	return result
}
