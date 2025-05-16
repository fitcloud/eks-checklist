package reliability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckProbe - 모든 Pod을 검색하여 startupProbe, livenessProbe, readinessProbe 가 모두 설정되었는지 확인
func CheckProbe(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[REL-005] Probe(Startup, Readiness, Liveness) 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 컨테이너에 startup/liveness/readiness probe가 누락되어 있습니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/reliability/REL-005",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			continue // kube-system은 검사 제외
		}

		for _, container := range pod.Spec.Containers {
			var missing []string
			if container.StartupProbe == nil {
				missing = append(missing, "startupProbe")
			}
			if container.LivenessProbe == nil {
				missing = append(missing, "livenessProbe")
			}
			if container.ReadinessProbe == nil {
				missing = append(missing, "readinessProbe")
			}

			if len(missing) > 0 {
				result.Passed = false
				result.Resources = append(result.Resources,
					fmt.Sprintf("Namespace: %s | Pod: %s | Container: %s (미설정: %v)", pod.Namespace, pod.Name, container.Name, missing))
			}
		}
	}

	return result
}
