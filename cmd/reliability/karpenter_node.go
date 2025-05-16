package reliability

import (
	"context"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// CheckKarpenterNode checks whether there are any Karpenter NodeClaims provisioned in the cluster.
func CheckKarpenterNode(karpenter_installed common.CheckResult, client dynamic.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[REL-013] Karpenter 기반 노드 생성",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/reliability/REL-013",
	}

	if !karpenter_installed.Passed {
		result.Passed = false
		result.FailureMsg = "Karpenter가 설치되어 있지 않습니다."
		return result
	}

	// NodeClaim GVR (Karpenter v0.37.x 기준)
	nodeClaimGVR := schema.GroupVersionResource{
		Group:    "karpenter.k8s.aws",
		Version:  "v1",
		Resource: "nodeclaims",
	}

	nodeClaims, err := client.Resource(nodeClaimGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	if len(nodeClaims.Items) == 0 {
		result.Passed = false
		result.FailureMsg = "Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다."
		return result
	}

	result.Passed = true

	return result
}
