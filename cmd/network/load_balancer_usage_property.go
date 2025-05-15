package network

import (
	"context"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func CheckLoadBalancerUsage(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "사용 사례에 맞는 로드밸런서 사용(ALB or NLB) - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "모든 Ingress 리소스를 수집하였습니다. 각 서비스에 적합한 ALB 또는 NLB 사용 여부는 수동으로 점검해야 합니다.",
		Runbook:    "https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/",
	}

	baseDir := filepath.Join(".", "output", eksCluster+"-loadbalancer-usage")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	ctx := context.TODO()
	ingList, err := client.NetworkingV1().Ingresses("").List(ctx, v1.ListOptions{})
	if err != nil {
		result.FailureMsg = "Ingress 목록 조회 실패: " + err.Error()
		return result
	}

	type IngressInfo struct {
		Namespace        string            `json:"namespace"`
		Name             string            `json:"name"`
		IngressClassName *string           `json:"ingressClassName,omitempty"`
		Annotations      map[string]string `json:"annotations,omitempty"`
	}

	var ingresses []IngressInfo
	for _, ing := range ingList.Items {
		ingresses = append(ingresses, IngressInfo{
			Namespace:        ing.Namespace,
			Name:             ing.Name,
			IngressClassName: ing.Spec.IngressClassName,
			Annotations:      ing.Annotations,
		})
	}

	outputPath := filepath.Join(baseDir, "ingresses.json")
	if err := common.SaveAsJSON(ingresses, outputPath); err == nil {
		result.Resources = append(result.Resources, "Ingress 목록: "+outputPath)
	}

	return result
}
