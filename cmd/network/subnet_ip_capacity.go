package network

import (
	"context"
	"fmt"
	"math"
	"net"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type EksCluster struct {
	Cluster *types.Cluster
}

// CheckVpcSubnetIpCapacity collects IP usage stats for all subnets and prints them for manual inspection.
func CheckVpcSubnetIpCapacity(eksCluster EksCluster, cfg aws.Config) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "VPC 서브넷에 충분한 IP 대역대 확보 - Manual",
		Manual:     true,
		Passed:     false,
		FailureMsg: "모든 서브넷의 IP 사용량을 출력했습니다. 사용 가능 용량이 충분한지 수동으로 확인하세요.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/network/subnetIpCapacity",
	}

	subnetIds := eksCluster.Cluster.ResourcesVpcConfig.SubnetIds
	ec2Client := ec2.NewFromConfig(cfg)

	for _, subnetId := range subnetIds {
		// 서브넷 정보 조회
		subnetOutput, err := ec2Client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
			SubnetIds: []string{subnetId},
		})
		if err != nil || len(subnetOutput.Subnets) == 0 {
			result.Resources = append(result.Resources,
				fmt.Sprintf("서브넷 %s 정보를 조회하는 데 실패했습니다.", subnetId))
			continue
		}

		subnet := subnetOutput.Subnets[0]

		// CIDR에서 전체 IP 수 계산
		_, ipNet, err := net.ParseCIDR(*subnet.CidrBlock)
		if err != nil {
			result.Resources = append(result.Resources,
				fmt.Sprintf("서브넷 %s CIDR 파싱 실패: %v", *subnet.SubnetId, err))
			continue
		}

		ones, bits := ipNet.Mask.Size()
		totalIPs := int(math.Pow(2, float64(bits-ones))) - 5 // AWS 예약 주소 제외
		availableIPs := int(*subnet.AvailableIpAddressCount)
		usageRatio := float64(availableIPs) / float64(totalIPs) * 100

		result.Resources = append(result.Resources,
			fmt.Sprintf("Subnet: %s | Total IPs: %d | Available IPs: %d | %.1f%% 사용 가능",
				*subnet.SubnetId, totalIPs, availableIPs, usageRatio))
	}

	return result
}
