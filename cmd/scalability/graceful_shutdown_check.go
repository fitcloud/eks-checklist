package scalability

import (
	"context"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func CheckGracefulShutdown(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Application에 Graceful shutdown 적용 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "Graceful shutdown 처리는 컨테이너 종료 이벤트 처리 여부를 코드 및 설정에서 수동 점검해야 합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/scalability/gracefulShutdown",
	}

	baseDir := filepath.Join(".", "output", eksCluster+"-graceful-shutdown")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	ctx := context.TODO()
	podList, err := client.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		result.FailureMsg = "Pod 목록 조회 실패: " + err.Error()
		return result
	}

	type ShutdownInfo struct {
		Namespace                  string `json:"namespace"`
		Pod                        string `json:"pod"`
		Container                  string `json:"container"`
		TerminationGracePeriodSecs int64  `json:"terminationGracePeriodSeconds,omitempty"`
		HasPreStopHook             bool   `json:"hasPreStopHook"`
	}

	var shutdownData []ShutdownInfo
	for _, pod := range podList.Items {
		grace := int64(30)
		if pod.Spec.TerminationGracePeriodSeconds != nil {
			grace = *pod.Spec.TerminationGracePeriodSeconds
		}

		for _, container := range pod.Spec.Containers {
			hasPreStop := container.Lifecycle != nil && container.Lifecycle.PreStop != nil
			shutdownData = append(shutdownData, ShutdownInfo{
				Namespace:                  pod.Namespace,
				Pod:                        pod.Name,
				Container:                  container.Name,
				TerminationGracePeriodSecs: grace,
				HasPreStopHook:             hasPreStop,
			})
		}
	}

	path := filepath.Join(baseDir, "graceful_shutdown_settings.json")
	if err := common.SaveAsJSON(shutdownData, path); err == nil {
		result.Resources = append(result.Resources, "Graceful Shutdown 관련 설정 정보: "+path)
	}

	return result
}
