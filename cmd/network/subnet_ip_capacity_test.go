package network_test

import (
	"context"
	"fmt"
	"math"
	"net"
	"reflect"
	"testing"

	"eks-checklist/cmd/network"
	"eks-checklist/cmd/testutils"

	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// FakeSubnet holds fake data for a subnet.
type FakeSubnet struct {
	CidrBlock               string
	AvailableIpAddressCount int64
	SubnetId                string
}

func TestCheckVpcSubnetIpCapacity_YAML(t *testing.T) {
	// YAML 파일 "subnet_ip_capacity.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "subnet_ip_capacity.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		// expected_result는 YAML에서 mapping으로 주어짐 (예: { "subnet-1": 10 })
		expectedResultRaw := tc["expected_result"].(map[string]interface{})
		expectedResult := make(map[string]int)
		for k, v := range expectedResultRaw {
			switch val := v.(type) {
			case float64:
				expectedResult[k] = int(val)
			case int:
				expectedResult[k] = val
			default:
				t.Fatalf("expected_result value is of unexpected type: %T", val)
			}
		}

		// cluster 정보 생성
		clusterMap := tc["cluster"].(map[string]interface{})
		resourcesVpcConfig := clusterMap["resourcesVpcConfig"].(map[string]interface{})
		subnetIdsRaw := resourcesVpcConfig["subnetIds"].([]interface{})
		var subnetIds []string
		for _, id := range subnetIdsRaw {
			subnetIds = append(subnetIds, id.(string))
		}
		eksCluster := network.EksCluster{
			Cluster: &eksTypes.Cluster{
				ResourcesVpcConfig: &eksTypes.VpcConfigResponse{
					SubnetIds: subnetIds,
				},
			},
		}

		// 테스트 케이스에서 fake 서브넷 데이터를 생성 (mapping: subnet_id → FakeSubnet)
		fakeSubnets := make(map[string]FakeSubnet)
		subnetsRaw := tc["subnets"].([]interface{})
		for _, s := range subnetsRaw {
			sMap := s.(map[string]interface{})
			subnetId := sMap["subnet_id"].(string)
			cidrBlock := sMap["cidr_block"].(string)

			var availableIp int64
			switch v := sMap["available_ip"].(type) {
			case float64:
				availableIp = int64(v)
			case int:
				availableIp = int64(v)
			default:
				t.Fatalf("available_ip is of unexpected type: %T", v)
			}

			fakeSubnets[subnetId] = FakeSubnet{
				CidrBlock:               cidrBlock,
				AvailableIpAddressCount: availableIp,
				SubnetId:                subnetId,
			}
		}

		t.Run(testName, func(t *testing.T) {
			// Patch ec2.NewFromConfig를 사용해 EC2 클라이언트를 생성한 후 DescribeSubnets 메서드를 monkey patch 합니다.
			patch := monkey.PatchInstanceMethod(reflect.TypeOf(new(ec2.Client)), "DescribeSubnets",
				func(c *ec2.Client, ctx context.Context, input *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
					if len(input.SubnetIds) == 0 {
						return nil, fmt.Errorf("no subnet id provided")
					}
					id := input.SubnetIds[0]
					fakeSub, ok := fakeSubnets[id]
					if !ok {
						return nil, fmt.Errorf("subnet not found: %s", id)
					}
					// CIDR 블록을 파싱하여 총 IP 개수를 계산하는 로직과 동일하게 동작하도록 함
					_, ipNet, err := net.ParseCIDR(fakeSub.CidrBlock)
					if err != nil {
						return nil, fmt.Errorf("CIDR parsing failed for %s: %v", fakeSub.CidrBlock, err)
					}
					ones, bits := ipNet.Mask.Size()
					totalIPs := int(math.Pow(2, float64(bits-ones))) - 5
					_ = totalIPs
					// 여기서는 fakeSub.AvailableIpAddressCount를 사용
					return &ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
							{
								CidrBlock:               aws.String(fakeSub.CidrBlock),
								AvailableIpAddressCount: aws.Int32(int32(fakeSub.AvailableIpAddressCount)),
								SubnetId:                aws.String(fakeSub.SubnetId),
							},
						},
					}, nil
				})
			defer patch.Unpatch()

			// 더미 AWS 구성(내용은 필요 없음)
			cfg := aws.Config{}

			result := network.CheckVpcSubnetIpCapacity(eksCluster, cfg)
			if !reflect.DeepEqual(result, expectedResult) {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectedResult, result)
			}
		})
	}
}
