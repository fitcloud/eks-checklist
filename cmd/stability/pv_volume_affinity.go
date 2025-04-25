package stability

import (
	"context"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func CheckVolumeAffinity(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "PV 사용시 volume affinity 위반 사항 체크 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "PV와 관련된 nodeAffinity 조건을 자동 수집하였으며, Pod 스케줄링 위치와의 일치 여부는 수동으로 점검해야 합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/stability/volumeAffinityCheck",
	}

	ctx := context.TODO()

	// 1. PVC 목록
	pvcList, err := client.CoreV1().PersistentVolumeClaims("").List(ctx, v1.ListOptions{})
	if err != nil {
		result.FailureMsg = "PVC 목록 조회 실패: " + err.Error()
		return result
	}

	if pvcList.Items == nil {
		result.Passed = true
		result.Manual = false
		return result
	}

	baseDir := filepath.Join(".", "result", eksCluster+"-volume-affinity")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	// 2. PV 목록 (Map 형태로 구성)
	pvMap := map[string]corev1.PersistentVolume{}
	if pvList, err := client.CoreV1().PersistentVolumes().List(ctx, v1.ListOptions{}); err == nil {
		for _, pv := range pvList.Items {
			pvMap[pv.Name] = pv
		}
	}

	// 3. Pod 목록
	podList, _ := client.CoreV1().Pods("").List(ctx, v1.ListOptions{})

	type VolumeAffinityInfo struct {
		Namespace    string                     `json:"namespace"`
		PVCName      string                     `json:"pvcName"`
		PVName       string                     `json:"pvName"`
		NodeAffinity *corev1.VolumeNodeAffinity `json:"nodeAffinity,omitempty"`
		PodName      string                     `json:"podName,omitempty"`
		PodNode      string                     `json:"podNode,omitempty"`
	}

	var results []VolumeAffinityInfo

	for _, pvc := range pvcList.Items {
		if pvc.Status.Phase != corev1.ClaimBound || pvc.Spec.VolumeName == "" {
			continue
		}
		pv, ok := pvMap[pvc.Spec.VolumeName]
		if !ok {
			continue
		}

		// Pod 연결 찾기 (volume에서 pvc 사용 중인 pod 추적)
		var podName, podNode string
		for _, pod := range podList.Items {
			if pod.Namespace != pvc.Namespace {
				continue
			}
			for _, vol := range pod.Spec.Volumes {
				if vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == pvc.Name {
					podName = pod.Name
					podNode = pod.Spec.NodeName
					break
				}
			}
		}

		results = append(results, VolumeAffinityInfo{
			Namespace:    pvc.Namespace,
			PVCName:      pvc.Name,
			PVName:       pv.Name,
			NodeAffinity: pv.Spec.NodeAffinity,
			PodName:      podName,
			PodNode:      podNode,
		})
	}

	outputPath := filepath.Join(baseDir, "pv_affinity_violations.json")
	if err := common.SaveAsJSON(results, outputPath); err == nil {
		result.Resources = append(result.Resources, "PV affinity 정보: "+outputPath)
	}

	return result
}
