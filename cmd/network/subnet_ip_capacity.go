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

// CheckVpcSubnetIpCapacity checks whether each subnet has at least 10% of its IP address capacity available.
func CheckVpcSubnetIpCapacity(eksCluster EksCluster, cfg aws.Config) common.CheckResult {
	result := common.CheckResult{
		CheckName: "VPC 서브넷에 충분한 IP 대역대 확보",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	subnetIds := eksCluster.Cluster.ResourcesVpcConfig.SubnetIds
	ec2Client := ec2.NewFromConfig(cfg)

	hasLowCapacity := false

	for _, subnetId := range subnetIds {
		// 서브넷 정보 조회
		subnetOutput, err := ec2Client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
			SubnetIds: []string{subnetId},
		})
		if err != nil || len(subnetOutput.Subnets) == 0 {
			result.Passed = false
			result.FailureMsg = fmt.Sprintf("서브넷 %s 정보를 조회하는 데 실패했습니다.", subnetId)
			return result
		}

		subnet := subnetOutput.Subnets[0]

		// CIDR에서 전체 IP 수 계산
		_, ipNet, err := net.ParseCIDR(*subnet.CidrBlock)
		if err != nil {
			result.Passed = false
			result.FailureMsg = fmt.Sprintf("서브넷 %s CIDR 파싱 실패: %v", *subnet.SubnetId, err)
			return result
		}

		ones, bits := ipNet.Mask.Size()
		totalIPs := int(math.Pow(2, float64(bits-ones))) - 5 // AWS 예약 5개 제외
		availableIPs := int(*subnet.AvailableIpAddressCount)

		usageRatio := float64(availableIPs) / float64(totalIPs)

		// 10% 미만이면 FAIL로 표시
		if usageRatio < 0.1 {
			hasLowCapacity = true

			// 결과 리소스로 기록
			result.Resources = append(result.Resources,
				fmt.Sprintf("Subnet: %s | Total IPs: %d | Available IPs: %d | %.1f%% 사용 가능",
					*subnet.SubnetId, totalIPs, availableIPs, usageRatio*100))
		}
	}

	if hasLowCapacity {
		result.Passed = false
		result.FailureMsg = "일부 VPC 서브넷의 사용 가능한 IP 주소가 전체의 10% 미만입니다."
	} else {
		result.Passed = true
	}

	return result
}
