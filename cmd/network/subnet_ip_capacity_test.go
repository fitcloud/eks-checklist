package network_test

import (
	"eks-checklist/cmd"
	"eks-checklist/cmd/network" // network 패키지 임포트
	"eks-checklist/cmd/testutils"

	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEc2Client은 ec2.Client를 모킹하기 위한 구조체입니다.
type MockEc2Client struct {
	mock.Mock
}

func (m *MockEc2Client) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ec2.DescribeSubnetsOutput), args.Error(1)
}

func TestCheckVpcSubnetIpCapacity(t *testing.T) {
	// LoadTestCases에서 반환되는 값은 []map[string]interface{}
	testCases := testutils.LoadTestCases(t, "subnet_ip_capacity_test.yaml")

	// MockEc2Client 설정
	mockEc2Client := new(MockEc2Client)

	// 반환된 testCases를 순차적으로 처리
	for _, tc := range testCases {
		t.Run(tc["name"].(string), func(t *testing.T) {
			// 서브넷 데이터를 가져와서 가짜 클러스터를 만든다.
			subnetData := tc["subnet_data"].(map[string]interface{})
			subnetID := subnetData["subnet_id"].(string)
			cidrBlock := subnetData["cidr_block"].(string)
			availableIpCount := int(subnetData["available_ip_count"].(int64))

			// Mocking DescribeSubnets 호출
			mockEc2Client.On("DescribeSubnets", &ec2.DescribeSubnetsInput{
				SubnetIds: []string{subnetID},
			}).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []types.Subnet{
					{
						SubnetId:                aws.String(subnetID),
						CidrBlock:               aws.String(cidrBlock),
						AvailableIpAddressCount: aws.Int32(int32(availableIpCount)),
					},
				},
			}, nil)

			// 가짜 EKS 클러스터 데이터 설정
			eksCluster := cmd.EksCluster{
				Cluster: *types.EksCluster{
					ResourcesVpcConfig: &ec2.VpcConfigResponse{
						SubnetIds: []string{subnetID},
					},
				},
			}

			// CheckVpcSubnetIpCapacity 호출
			subnetIpCapacity := network.CheckVpcSubnetIpCapacity(eksCluster)

			// 테스트 조건에 맞게 결과 검증
			lowIPs := float64(availableIpCount) < 0.1*float64(251) // 251은 /24 서브넷의 IP 개수에서 5를 뺀 값
			if tc["expect_failure"].(bool) {
				assert.True(t, lowIPs, "Expected failure but condition did not match")
			} else {
				assert.False(t, lowIPs, "Expected success but condition did not match")
			}
		})
	}
}
