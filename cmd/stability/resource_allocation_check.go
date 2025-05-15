package stability

import (
	"context"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"gopkg.in/yaml.v3"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func CheckResourceAllocation(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "애플리케이션에 적절한 CPU/RAM 할당",
		Manual:     true,
		Passed:     false,
		FailureMsg: "일부 Pod에 Request/Limit 설정이 없거나, 값의 적정성은 수동 확인이 필요합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/stability/resourceRequestsLimits",
	}

	baseDir := filepath.Join(".", "output", eksCluster+"-resource-allocation")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	yamlDir := filepath.Join(baseDir, "yamls")
	ctx := context.TODO()
	podList, err := client.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		result.FailureMsg = "Pod 목록 조회 실패: " + err.Error()
		return result
	}

	type ResourceInfo struct {
		Namespace     string `json:"namespace"`
		Pod           string `json:"pod"`
		Container     string `json:"container"`
		HasRequest    bool   `json:"hasRequest"`
		HasLimit      bool   `json:"hasLimit"`
		RequestCPU    string `json:"requestCPU,omitempty"`
		RequestMemory string `json:"requestMemory,omitempty"`
		LimitCPU      string `json:"limitCPU,omitempty"`
		LimitMemory   string `json:"limitMemory,omitempty"`
	}

	var incomplete []ResourceInfo
	ExistSetting := false

	for _, pod := range podList.Items {
		if pod.Namespace == "kube-system" {
			continue
		}

		for _, container := range pod.Spec.Containers {
			res := container.Resources
			info := ResourceInfo{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
				Container: container.Name,
			}

			if res.Requests != nil {
				info.HasRequest = true
				if cpu := res.Requests.Cpu(); cpu != nil {
					info.RequestCPU = cpu.String()
				}
				if mem := res.Requests.Memory(); mem != nil {
					info.RequestMemory = mem.String()
				}
			}
			if res.Limits != nil {
				info.HasLimit = true
				if cpu := res.Limits.Cpu(); cpu != nil {
					info.LimitCPU = cpu.String()
				}
				if mem := res.Limits.Memory(); mem != nil {
					info.LimitMemory = mem.String()
				}
			}

			if !info.HasRequest || !info.HasLimit {
				incomplete = append(incomplete, info)
			}

			// request, limits 설정이 되어 있는 Pod 일 경우, 설정값을 YAML 파일로 저장
			if res.Requests != nil && res.Limits != nil {
				ExistSetting = true

				if err := os.MkdirAll(yamlDir, os.ModePerm); err != nil {
					result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
					return result
				}

				containerSpec := map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name": container.Name,
								"resources": map[string]interface{}{
									"requests": res.Requests,
									"limits":   res.Limits,
								},
							},
						},
					},
				}

				yamlBytes, err := yaml.Marshal(containerSpec)
				if err == nil {
					yamlPath := filepath.Join(yamlDir, pod.Namespace+"-"+pod.Name+".yaml")
					os.WriteFile(yamlPath, yamlBytes, 0644)
				}
			}
		}
	}
	if ExistSetting {
		result.Resources = append(result.Resources, "Pod별 리소스 설정 YAML 디렉토리: "+yamlDir)
	} else {
		result.Manual = false
		result.FailureMsg = "모든 Pod에 Request/Limit 설정이 되어있지 않습니다."
	}

	// JSON으로 미설정 정보 저장
	jsonPath := filepath.Join(baseDir, "resource_allocation_check.json")
	if err := common.SaveAsJSON(incomplete, jsonPath); err == nil {
		result.Resources = append(result.Resources, "리소스 미설정/불완전 설정 목록: "+jsonPath)
	}

	return result
}
