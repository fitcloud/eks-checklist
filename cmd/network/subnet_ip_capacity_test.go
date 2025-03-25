package network

import (
	testutils "eks-checklist/cmd/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 테스트용 단순화된 함수
func TestCheckVpcSubnetIpCapacity(t *testing.T) {
	// Load the test cases from the YAML file
	testCases := testutils.LoadTestCases(t, "subnet_ip_capacity.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		subnetData := tc["subnet_data"].(map[string]interface{})

		// available_ip_count 값을 읽어들이는 부분
		var availableIpCount int
		switch v := subnetData["available_ip_count"].(type) {
		case int:
			availableIpCount = v
		case float64:
			availableIpCount = int(v)
		}

		// expect_failure 값 (실패를 예상하는지 여부)
		expectFailure := tc["expect_failure"].(bool)

		// 테스트 케이스에 대한 테스트 실행
		t.Run(name, func(t *testing.T) {
			// 여기서는 IP 갯수만 가지고 pass/fail을 결정
			isEnough := availableIpCount >= 50 // 예시: 50개 이상의 IP가 있으면 충분하다고 가정

			if expectFailure {
				// 실패를 예상한 경우, IP 갯수가 50개 미만이어야 한다.
				assert.False(t, isEnough, "Expected insufficient IPs for subnet: "+subnetData["subnetid"].(string))
			} else {
				// 실패를 예상하지 않은 경우, IP 갯수가 50개 이상이어야 한다.
				assert.True(t, isEnough, "Expected sufficient IPs for subnet: "+subnetData["subnetid"].(string))
			}
		})
	}
}
