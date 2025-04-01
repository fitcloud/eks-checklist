package security

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"
)

func CheckAccessControl(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "클러스터 접근 제어(Access entries, aws-auth 컨피그맵)",
		Manual:     true,
		Passed:     false,
		FailureMsg: "클러스터 접근 제어 설정이 되어 있으나, 적합한 설정이 되어 있는지 수동으로 확인해야 합니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	// 👉 실행 디렉토리 기준 ./result 하위 경로 생성
	baseDir := filepath.Join(".", "result", eksCluster+"-access-control")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.Passed = false
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	// ---------------------------------------
	// 1. aws-auth ConfigMap 저장
	// ---------------------------------------
	hasConfigMap := false
	configMapClient := client.CoreV1().ConfigMaps("kube-system")
	configMap, err := configMapClient.Get(context.TODO(), "aws-auth", v1.GetOptions{})
	if err == nil {
		configMapPath := filepath.Join(baseDir, "aws-auth-configmap.yaml")
		err = common.SaveK8sResourceAsYAML(configMap, configMapPath)
		if err != nil {
			result.Passed = false
			result.FailureMsg = "aws-auth ConfigMap 저장 실패: " + err.Error()
			return result
		}
		hasConfigMap = true
		result.Resources = append(result.Resources, "aws-auth ConfigMap 저장 경로: "+configMapPath)
	}

	// ---------------------------------------
	// 2. Access Entries 저장
	// ---------------------------------------
	hasAccessEntries := false
	var accessEntries []interface{}

	eksClient := eks.NewFromConfig(cfg)

	listResp, err := eksClient.ListAccessEntries(context.TODO(), &eks.ListAccessEntriesInput{
		ClusterName: &eksCluster,
	})
	if err == nil && len(listResp.AccessEntries) > 0 {
		for _, ae := range listResp.AccessEntries {
			descResp, err := eksClient.DescribeAccessEntry(context.TODO(), &eks.DescribeAccessEntryInput{
				PrincipalArn: &ae,
				ClusterName:  &eksCluster,
			})
			if err == nil {
				accessEntries = append(accessEntries, descResp.AccessEntry)
			}
		}
	}

	if len(accessEntries) > 0 {
		hasAccessEntries = true
		accessEntryPath := filepath.Join(baseDir, "access-entries.json")
		err := common.SaveAsJSON(accessEntries, accessEntryPath)
		if err != nil {
			result.Passed = false
			result.FailureMsg = "Access Entries 저장 실패: " + err.Error()
			return result
		}
		result.Resources = append(result.Resources, "Access Entries 저장 경로: "+accessEntryPath)
	}

	// ---------------------------------------
	// 3. 최종 결과 판단
	// ---------------------------------------
	if !hasConfigMap && !hasAccessEntries {
		result.Manual = false
		result.Passed = false
		result.FailureMsg = "aws-auth ConfigMap과 Access Entries 설정이 모두 존재하지 않습니다."
	}

	return result
}
