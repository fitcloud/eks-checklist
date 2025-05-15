package security

import (
	"context"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func CheckImageStaticAnalysis(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "컨테이너 이미지 정적 분석 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "컨테이너 이미지의 보안 취약점 여부는 수동으로 정적 분석 도구를 사용해 확인해야 합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/security/imageScanning",
	}

	// 결과 디렉토리 생성
	baseDir := filepath.Join(".", "output", eksCluster+"-image-analysis")
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

	type ImageInfo struct {
		Namespace string `json:"namespace"`
		Pod       string `json:"pod"`
		Container string `json:"container"`
		Image     string `json:"image"`
	}

	var images []ImageInfo
	for _, pod := range podList.Items {
		for _, container := range pod.Spec.Containers {
			images = append(images, ImageInfo{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
				Container: container.Name,
				Image:     container.Image,
			})
		}
		for _, initContainer := range pod.Spec.InitContainers {
			images = append(images, ImageInfo{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
				Container: initContainer.Name + " (init)",
				Image:     initContainer.Image,
			})
		}
	}

	// 이미지 목록 저장
	outputPath := filepath.Join(baseDir, "container_images.json")
	if err := common.SaveAsJSON(images, outputPath); err == nil {
		result.Resources = append(result.Resources, "컨테이너 이미지 목록: "+outputPath)
	}

	return result
}
