package security

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// DataplanePrivateCheck checks whether all subnets used by the EKS data plane are private.
func DataplanePrivateCheck(eksCluster EksCluster, cfg aws.Config) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "데이터 플레인 사설망",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 서브넷이 IGW(인터넷 게이트웨이)와 연결되어 있어 퍼블릭 상태입니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	subnetIds := eksCluster.Cluster.ResourcesVpcConfig.SubnetIds
	ec2Client := ec2.NewFromConfig(cfg)

	var publicSubnets []string

	for _, subnetId := range subnetIds {
		rtOut, err := ec2Client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("association.subnet-id"),
					Values: []string{subnetId},
				},
			},
		})
		if err != nil {
			result.Passed = false
			result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
			return result
		}

		for _, rt := range rtOut.RouteTables {
			for _, route := range rt.Routes {
				if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == "0.0.0.0/0" {
					if route.GatewayId != nil && strings.HasPrefix(*route.GatewayId, "igw-") {
						publicSubnets = append(publicSubnets, subnetId)
						break
					}
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
