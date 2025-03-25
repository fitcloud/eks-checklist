package security

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// 각 서브넷이 IGW에 연결되어 있다면 반환하는 함수
func DataplanePrivateCheck(eksCluster EksCluster, cfg aws.Config) []string {
	subnetIds := eksCluster.Cluster.ResourcesVpcConfig.SubnetIds
	ec2Client := ec2.NewFromConfig(cfg)

	var publicSubnets []string

	for _, subnetId := range subnetIds {
		// 라우트 테이블 조회
		rtOut, err := ec2Client.DescribeRouteTables(context.TODO(), &ec2.DescribeRouteTablesInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("association.subnet-id"),
					Values: []string{subnetId},
				},
			},
		})
		if err != nil {
			log.Printf("Failed to describe route table for subnet %s: %v", subnetId, err)
			continue
		}

		// IGW 연결 여부 확인
		for _, rt := range rtOut.RouteTables {
			for _, route := range rt.Routes {
				if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == "0.0.0.0/0" {
					if route.GatewayId != nil && strings.HasPrefix(*route.GatewayId, "igw-") {
						publicSubnets = append(publicSubnets, subnetId)
						log.Printf("Subnet %s is connected to IGW via route table %s", subnetId, *rt.RouteTableId)
						break
					}
				}
			}
		}
	}

	return publicSubnets
}
