package stability

import (
	"context"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// CheckKarpenterNode 확인: nodeclaims.karpenter.sh 리소스를 직접 조회하여 1개 이상 있는지 확인
// 지금 고려사항이 kapenter 버전에 따른 변경사항이 있을 수 있음
// 예를 들면 v1beta1은 0.37.x 까지의 사항이므로 1.x.x 가는 경우에는 검토을 어떻게 할지 두번해야할지?
func CheckKarpenterNode(client dynamic.Interface) bool {
	// Karpenter NodeClaim의 GVR (GroupVersionResource) 정의
	nodeClaimGVR := schema.GroupVersionResource{
		Group:    "karpenter.k8s.aws",
		Version:  "v1beta1",
		Resource: "nodeclaims",
	}

	// NodeClaims 조회
	nodeClaims, err := client.Resource(nodeClaimGVR).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list Karpenter NodeClaims: %v", err)
		return false
	}

	// NodeClaim 개수가 1개 이상이면 Karpenter 노드 존재
	return len(nodeClaims.Items) > 0
}
