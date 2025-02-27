package network

import (
	"context"
	"log"
	"math"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type EksCluster struct {
	Cluster *types.Cluster
}

// EKS 클러스터가 사용하는 서브넷의 가용가능한 IP 주소가 10% 미만인 경우 SubnetId와 사용가능한 IP 개수를 반환하는 함수
func CheckVpcSubnetIpCapacity(eksCluster EksCluster) map[string]int {
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

// 모든 namespace에서 aws-node 데몬셋을 가져와서 containers.env.key 중 ENABLE_PREFIX_DELEGATION의 값이 true 인지 false인지 확인
func CheckVpcCniPrefixMode(client kubernetes.Interface) bool {
	// instance가 nitro 기반인이 확인하는 로직이 필요하나 문제점에 봉착함
	// api 중에 노드만 가져오는게 없어서 nodegroup 단위로 검색하면 karpenter node가 누락되고
	// tag key등으로 구분해도 nodegroup이랑 karpenter랑 서로 tag가 다름
	// 그럼으로 음... 보류 질문해야할 듯

	// 모든 namespace에서 aws-node 데몬셋을 가져옴
	daemonsets, err := client.AppsV1().DaemonSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, daemonset := range daemonsets.Items {
		if daemonset.Name == "aws-node" {
			for _, container := range daemonset.Spec.Template.Spec.Containers {
				for _, env := range container.Env {
					if env.Name == "ENABLE_PREFIX_DELEGATION" {
						if env.Value == "true" {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// 모든 namespace에서 deployement를 가져와서 container.image에 aws-load-balancer-controller가 포함되어 있는지 확인
func CheckAwsLoadBalancerController(client kubernetes.Interface) bool {
	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "aws-load-balancer-controller") {
				return true
			}
		}
	}

	return false
}
