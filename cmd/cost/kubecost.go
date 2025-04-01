package cost

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetKubecost checks if Kubecost is deployed in the cluster.
func GetKubecost(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "EKS용 Kubecost 설치",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	deploys, err := client.AppsV1().Deployments(v1.NamespaceAll).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}
	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "kubecost") {
				result.Passed = true
				result.SuccessMsg = "Kubecost가 클러스터에 설치되어 있습니다."
				result.Resources = append(result.Resources,
					fmt.Sprintf("Namespace: %s | Deployment: %s | Image: %s",
						deploy.Namespace, deploy.Name, container.Image))
				return result
			}
		}
	}

	result.Passed = false
	result.FailureMsg = "Kubecost가 클러스터에 설치되어 있지 않습니다."
	return result
}
