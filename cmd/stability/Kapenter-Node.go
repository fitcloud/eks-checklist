package stability

import (
	"context"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// CheckKarpenterNode 확인: nodeclaims.karpenter.sh 리소스를 직접 조회하여 1개 이상 있는지 확인
// 다이나믹 클라이언트 어떻게 처리할지 고민좀 해봐야할 듯 k8s로 안되서 일단 패스
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
