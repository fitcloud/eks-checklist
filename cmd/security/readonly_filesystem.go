package security

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ReadnonlyFilesystemCheck checks whether containers use readOnlyRootFilesystem.
func ReadnonlyFilesystemCheck(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[SEC-014] 읽기 전용 파일시스템 사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 컨테이너가 readOnlyRootFilesystem=true 설정을 사용하지 않고 있습니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/security/SEC-014",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	nodeOSCache := make(map[string]string)

	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			continue // kube-system 네임스페이스는 검사 제외
		}

		nodeName := pod.Spec.NodeName
		nodeOS, ok := nodeOSCache[nodeName]
		if !ok {
			node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, v1.GetOptions{})
			if err != nil {
				nodeOS = "unknown"
			} else if osLabel, exists := node.Labels["kubernetes.io/os"]; exists {
				nodeOS = osLabel
			} else {
				nodeOS = "unknown"
			}
			nodeOSCache[nodeName] = nodeOS
		}

		for _, container := range pod.Spec.Containers {
			if nodeOS == "windows" {
				continue
			}

			sc := container.SecurityContext
			if sc == nil || sc.ReadOnlyRootFilesystem == nil || !*sc.ReadOnlyRootFilesystem {
				result.Passed = false
				result.Resources = append(result.Resources,
					fmt.Sprintf("Namespace: %s | Pod: %s | Container: %s", pod.Namespace, pod.Name, container.Name))
			}
		}
	}

	return result
}
