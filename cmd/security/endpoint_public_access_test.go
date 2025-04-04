package security_test

import (
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestCheckEndpointPublicAccess(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "endpoint_public_access.yaml")

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)
		endpointPublicAccess := tc["endpoint_public_access"].(bool)

		t.Run(testName, func(t *testing.T) {
			eksCluster := security.EksCluster{
				Cluster: &types.Cluster{
					Name: strPtr(testName),
					ResourcesVpcConfig: &types.VpcConfigResponse{
						EndpointPublicAccess: endpointPublicAccess,
					},
				},
			}

			result := security.CheckEndpointPublicAccess(eksCluster)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}
