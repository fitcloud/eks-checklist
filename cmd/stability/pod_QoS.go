package stability

import (
	"context"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func CheckQoSClass(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "애플리케이션 중요도에 따른 QoS 적용 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "Pod의 QoS 클래스는 자동으로 분석되었으며, 애플리케이션 중요도에 따라 적절한 QoS가 적용되었는지는 수동 판단해야 합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/stability/qosByPriority",
	}

	baseDir := filepath.Join(".", "output", eksCluster+"-qos-class")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	ctx := context.TODO()
	pods, err := client.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		result.FailureMsg = "Pod 목록 조회 실패: " + err.Error()
		return result
	}

	type QoSInfo struct {
		Namespace string `json:"namespace"`
		Pod       string `json:"pod"`
		QoSClass  string `json:"qosClass"`
	}

	var qosResults []QoSInfo
	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			continue
		}
		qosResults = append(qosResults, QoSInfo{
			Namespace: pod.Namespace,
			Pod:       pod.Name,
			QoSClass:  string(pod.Status.QOSClass),
		})
	}

	outputPath := filepath.Join(baseDir, "qos_class_summary.json")
	if err := common.SaveAsJSON(qosResults, outputPath); err == nil {
		result.Resources = append(result.Resources, "QoS 클래스 요약 정보: "+outputPath)
	}

	return result
}
