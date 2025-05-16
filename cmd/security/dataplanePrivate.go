package security

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

// DataplanePrivateCheck checks whether all subnets used by the EKS data plane are private.
func DataplanePrivateCheck(eksCluster EksCluster, cfg aws.Config) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[SEC-012] 데이터 플레인 사설망",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 서브넷이 IGW(인터넷 게이트웨이)와 연결되어 있어 퍼블릭 상태입니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/security/SEC-012",
	}

	eksClient := eks.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)

	// 노드 그룹 목록 조회
	nodeGroupsOutput, err := eksClient.ListNodegroups(context.TODO(), &eks.ListNodegroupsInput{
		ClusterName: aws.String(*eksCluster.Cluster.Name),
	})
	if err != nil {
		result.Passed = false
		result.FailureMsg = fmt.Sprintf("노드 그룹 목록 조회 실패: %v", err)
		return result
	}

	subnetIDSet := make(map[string]struct{})
	var vpcID string

	// 각 노드 그룹의 서브넷 ID 수집
	for _, nodeGroupName := range nodeGroupsOutput.Nodegroups {
		nodeGroupOutput, err := eksClient.DescribeNodegroup(context.TODO(), &eks.DescribeNodegroupInput{
			ClusterName:   aws.String(*eksCluster.Cluster.Name),
			NodegroupName: aws.String(nodeGroupName),
		})
		if err != nil {
			result.Passed = false
			result.FailureMsg = fmt.Sprintf("노드 그룹 '%s' 상세 정보 조회 실패: %v", nodeGroupName, err)
			return result
		}

		for _, subnetID := range nodeGroupOutput.Nodegroup.Subnets {
			subnetIDSet[subnetID] = struct{}{}
		}
	}

	// 서브넷 ID 목록 생성
	var subnetIDs []string
	for id := range subnetIDSet {
		subnetIDs = append(subnetIDs, id)
	}

	// VPC ID가 설정되지 않은 경우 클러스터의 VPC ID 사용
	if vpcID == "" {
		vpcID = *eksCluster.Cluster.ResourcesVpcConfig.VpcId
	}

	// VPC의 모든 라우트 테이블 조회
	rtOut, err := ec2Client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	})
	if err != nil {
		result.Passed = false
		result.FailureMsg = fmt.Sprintf("라우트 테이블 조회 실패: %v", err)
		return result
	}

	// 서브넷 ID -> 연결된 라우트 테이블 매핑
	subnetToRouteTable := make(map[string]ec2types.RouteTable)

	for _, rt := range rtOut.RouteTables {
		for _, assoc := range rt.Associations {
			if assoc.SubnetId != nil {
				subnetToRouteTable[*assoc.SubnetId] = rt
			}
		}
	}

	var publicSubnets []string

	for _, subnetID := range subnetIDs {
		rt, exists := subnetToRouteTable[subnetID]

		// 서브넷에 직접 연결된 라우트 테이블이 없으면, 메인 라우트 테이블 사용
		if !exists {
			for _, rtCandidate := range rtOut.RouteTables {
				for _, assoc := range rtCandidate.Associations {
					if assoc.Main != nil && *assoc.Main {
						rt = rtCandidate
						exists = true
						break
					}
				}
				if exists {
					break
				}
			}
		}

		// IGW로 가는 0.0.0.0/0 경로가 있는지 확인
		for _, route := range rt.Routes {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == "0.0.0.0/0" {
				if route.GatewayId != nil {
					if strings.HasPrefix(*route.GatewayId, "igw-") {
						publicSubnets = append(publicSubnets, subnetID)
						break
					}
				} else {
					fmt.Printf("서브넷 %s의 0.0.0.0/0 경로에 GatewayId가 없습니다.\n", subnetID)
				}
			}
		}
	}

	if len(publicSubnets) > 0 {
		result.Passed = false
		for _, subnet := range publicSubnets {
			result.Resources = append(result.Resources, fmt.Sprintf("Public Subnet: %s", subnet))
		}
	}

	return result
}
