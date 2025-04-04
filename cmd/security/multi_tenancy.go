package security

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"eks-checklist/cmd/common"
)

func CheckMultitenancy(client kubernetes.Interface, cfg aws.Config, eksCluster string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "멀티 테넌시 적용 유무",
		Manual:     true,
		Passed:     false,
		FailureMsg: "멀티 테넌시 격리 구성은 수동으로 점검이 필요합니다. 네임스페이스, 네트워크 정책, RBAC, 쿼터, IRSA, 우선순위 등 관련 리소스를 확인하세요.",
		Runbook:    "https://docs.aws.amazon.com/eks/latest/best-practices/tenant-isolation.html",
	}

	baseDir := filepath.Join(".", "result", eksCluster+"-multitenancy")
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		result.FailureMsg = "결과 디렉토리 생성 실패: " + err.Error()
		return result
	}

	ctx := context.TODO()

	// 1. Namespaces
	if ns, err := client.CoreV1().Namespaces().List(ctx, v1.ListOptions{}); err == nil {
		path := filepath.Join(baseDir, "namespaces.json")
		if err := common.SaveAsJSON(ns.Items, path); err == nil {
			result.Resources = append(result.Resources, "네임스페이스 목록: "+path)
		}
	}

	// 2. NetworkPolicy
	if np, err := client.NetworkingV1().NetworkPolicies("").List(ctx, v1.ListOptions{}); err == nil {
		path := filepath.Join(baseDir, "network_policies.json")
		if err := common.SaveAsJSON(np.Items, path); err == nil {
			result.Resources = append(result.Resources, "네트워크 정책 목록: "+path)
		}
	}

	// 3. RBAC (RoleBinding + ClusterRoleBinding)
	var rbacData []interface{}
	if rb, err := client.RbacV1().RoleBindings("").List(ctx, v1.ListOptions{}); err == nil {
		rbacData = append(rbacData, rb.Items)
	}
	if crb, err := client.RbacV1().ClusterRoleBindings().List(ctx, v1.ListOptions{}); err == nil {
		rbacData = append(rbacData, crb.Items)
	}
	if len(rbacData) > 0 {
		path := filepath.Join(baseDir, "rbac.json")
		if err := common.SaveAsJSON(rbacData, path); err == nil {
			result.Resources = append(result.Resources, "RBAC 설정 목록: "+path)
		}
	}

	// 4. ResourceQuota
	if rq, err := client.CoreV1().ResourceQuotas("").List(ctx, v1.ListOptions{}); err == nil {
		path := filepath.Join(baseDir, "resource_quotas.json")
		if err := common.SaveAsJSON(rq.Items, path); err == nil {
			result.Resources = append(result.Resources, "리소스 쿼터 목록: "+path)
		}
	}

	// 5. LimitRange
	if lr, err := client.CoreV1().LimitRanges("").List(ctx, v1.ListOptions{}); err == nil {
		path := filepath.Join(baseDir, "limit_ranges.json")
		if err := common.SaveAsJSON(lr.Items, path); err == nil {
			result.Resources = append(result.Resources, "LimitRange 목록: "+path)
		}
	}

	// 6. IRSA (ServiceAccounts with role annotation)
	if saList, err := client.CoreV1().ServiceAccounts("").List(ctx, v1.ListOptions{}); err == nil {
		var irsaList []interface{}
		for _, sa := range saList.Items {
			if _, ok := sa.Annotations["eks.amazonaws.com/role-arn"]; ok {
				irsaList = append(irsaList, sa)
			}
		}
		if len(irsaList) > 0 {
			path := filepath.Join(baseDir, "irsa_service_accounts.json")
			if err := common.SaveAsJSON(irsaList, path); err == nil {
				result.Resources = append(result.Resources, "IRSA 서비스 계정 목록: "+path)
			}
		}
	}

	// 7. PriorityClasses
	if pc, err := client.SchedulingV1().PriorityClasses().List(ctx, v1.ListOptions{}); err == nil {
		path := filepath.Join(baseDir, "priority_classes.json")
		if err := common.SaveAsJSON(pc.Items, path); err == nil {
			result.Resources = append(result.Resources, "PriorityClass 목록: "+path)
		}
	}

	return result
}
