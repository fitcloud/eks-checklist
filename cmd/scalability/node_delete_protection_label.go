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

func CheckImportantPodProtection(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "중요 Pod에 노드 삭제 방지용 Label 부여 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "중요한 Pod가 실행 중인 노드에 삭제 방지용 라벨이 설정되었는지 수동으로 확인해야 합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/scalability/podEvictionProtection",
	}

	baseDir := filepath.Join(".", "output", eksCluster+"-important-pod-protection")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	ctx := context.TODO()

	// Pod 목록 조회
	podList, err := client.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		result.FailureMsg = "Pod 목록 조회 실패: " + err.Error()
		return result
	}

	// Node 정보 조회
	nodeMap := map[string]map[string]string{}
	nodeList, err := client.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err == nil {
		for _, node := range nodeList.Items {
			nodeMap[node.Name] = node.Labels
		}
	}

	type PodNodeProtection struct {
		Namespace  string            `json:"namespace"`
		Pod        string            `json:"pod"`
		Node       string            `json:"node"`
		NodeLabels map[string]string `json:"nodeLabels"`
	}

	var protectionList []PodNodeProtection
	for _, pod := range podList.Items {
		// 중요 네임스페이스만 필터링 (또는 모든 Pod 수집해도 됨)
		if pod.Status.Phase == "Running" && pod.Spec.NodeName != "" {
			protectionList = append(protectionList, PodNodeProtection{
				Namespace:  pod.Namespace,
				Pod:        pod.Name,
				Node:       pod.Spec.NodeName,
				NodeLabels: nodeMap[pod.Spec.NodeName],
			})
		}
	}

	outputPath := filepath.Join(baseDir, "pod_node_labels.json")
	if err := common.SaveAsJSON(protectionList, outputPath); err == nil {
		result.Resources = append(result.Resources, "Pod 실행 노드의 Label 정보: "+outputPath)
	}

	return result
}
