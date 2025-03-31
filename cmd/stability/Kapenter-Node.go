// 기존 코드

// package stability

// import (
// 	"context"
// 	"log"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic"
// )

// // CheckKarpenterNode 확인: nodeclaims.karpenter.sh 리소스를 직접 조회하여 1개 이상 있는지 확인
// // 지금 고려사항이 kapenter 버전에 따른 변경사항이 있을 수 있음
// // 예를 들면 v1beta1은 0.37.x 까지의 사항이므로 1.x.x 가는 경우에는 검토을 어떻게 할지 두번해야할지?
// func CheckKarpenterNode(client dynamic.Interface) bool {
// 	// Karpenter NodeClaim의 GVR (GroupVersionResource) 정의
// 	nodeClaimGVR := schema.GroupVersionResource{
// 		Group:    "karpenter.k8s.aws",
// 		Version:  "v1beta1",
// 		Resource: "nodeclaims",
// 	}

// 	// NodeClaims 조회
// 	nodeClaims, err := client.Resource(nodeClaimGVR).List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		log.Printf("Failed to list Karpenter NodeClaims: %v", err)
// 		return false
// 	}

// 	// NodeClaim 개수가 1개 이상이면 Karpenter 노드 존재
// 	return len(nodeClaims.Items) > 0
// }

// 변경 후 코드
package stability

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
		CheckName:  "Karpenter 기반 노드 생성",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Karpenter가 노드를 프로비저닝한 흔적(NodeClaim 리소스)이 존재하지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	if !karpenter_installed.Passed {
		result.Passed = false
		result.FailureMsg = "Karpenter가 설치되어 있지 않습니다."
		return result
	}

	// NodeClaim GVR (Karpenter v0.37.x 기준)
	nodeClaimGVR := schema.GroupVersionResource{
		Group:    "karpenter.k8s.aws",
		Version:  "v1beta1",
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
