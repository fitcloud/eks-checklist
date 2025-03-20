package network

import (
	"context"
	"eks-checklist/cmd"
	"log"
	"math"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EKS 클러스터가 사용하는 서브넷의 가용가능한 IP 주소가 10% 미만인 경우 SubnetId와 사용가능한 IP 개수를 반환하는 함수
func CheckVpcSubnetIpCapacity(eksCluster cmd.EksCluster) map[string]int {
	// EKS가 배포된 VPC의 ID 및 서브넷 가져오기
	subnetIds := eksCluster.Cluster.ResourcesVpcConfig.SubnetIds

	cfg := GetAWSConfig()
	// AWS SDK 설정
	ec2Client := ec2.NewFromConfig(cfg)

	// 서브넷 정보 가져와서 서브넷들이 사용가능한 IP 개수와 ID를 맵에 저장
	subnetIpCapacity := make(map[string]int)
	for _, subnetId := range subnetIds {
		subnet, err := ec2Client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
			SubnetIds: []string{subnetId},
		})
		if err != nil {
			log.Fatalf("unable to describe subnets, %v", err)
		}

		// 서브넷의 CIDR 블록을 파싱
		_, ipNet, err := net.ParseCIDR(*subnet.Subnets[0].CidrBlock)
		if err != nil {
			log.Fatalf("CIDR parsing failed: %v", err)
		}
		// 프리픽스 길이 계산 후 변수에 저장, AWS에서 예약된 5개 주소 제외
		// 계산식: (2^bits - ones) - 5, 예를 들면 2^(32 - 24) - 5 = 251
		ones, bits := ipNet.Mask.Size()
		totalIPs := int(math.Pow(2, float64(bits-ones))) - 5
		// 디버깅 용 총 IP 개수 출력
		log.Printf("total IPs: %d", totalIPs)

		// 서브넷의 사용가능한 IP 개수를 변수에 저장
		avaliableIp := int(*subnet.Subnets[0].AvailableIpAddressCount)
		// 디버깅 용 사용가능한 IP 개수 출력
		log.Printf("available IPs: %d", avaliableIp)

		// 사용가능한 IP 개수가 총 IP 개수의 10% 미만이면 변수에 Subnet ID와 사용가능한 IP 개수 저장
		if float64(avaliableIp) < 0.1*float64(totalIPs) {
			subnetIpCapacity[*subnet.Subnets[0].SubnetId] = avaliableIp
		}
	}
	return subnetIpCapacity
}
