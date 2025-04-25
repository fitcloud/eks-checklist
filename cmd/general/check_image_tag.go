package general

import (
	"context"
	"eks-checklist/cmd/common"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckImageTag(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "컨테이너 이미지 태그에 latest 미사용",
		Manual:    false,
		Passed:    true,
		// SuccessMsg: "모든 컨테이너 이미지는 latest 태그를 사용 중이지 않습니다.",
		FailureMsg: "일부 컨테이터 이미지가 latest 태그를 사용 중입니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/general/noLatestTag",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if strings.Contains(container.Image, "latest") {
				result.Passed = false
				result.Resources = append(result.Resources, "Namespace: "+pod.Namespace+" | Pod: "+pod.Name+" | Container: "+container.Name+" | Image: "+container.Image)
			}
		}
	}

	return result
}
