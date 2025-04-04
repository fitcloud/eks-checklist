package security_test

import (
	"context"
	"reflect"
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestDataplanePrivateCheck(t *testing.T) {
	// YAML 파일 "dataplane_private.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "dataplane_private.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectedPass := tc["expect_pass"].(bool)

		clusterMap := tc["cluster"].(map[string]interface{})
		resVpcCfg := clusterMap["resourcesVpcConfig"].(map[string]interface{})
		subnetIdsRaw := resVpcCfg["subnetIds"].([]interface{})
		var subnetIds []string
		for _, id := range subnetIdsRaw {
			subnetIds = append(subnetIds, id.(string))
		}
		eksCluster := security.EksCluster{
			Cluster: &eksTypes.Cluster{
				ResourcesVpcConfig: &eksTypes.VpcConfigResponse{
					SubnetIds: subnetIds,
				},
			},
		}

		routeTablesRaw := tc["route_tables"].(map[string]interface{})
		// route_tables: mapping: subnet_id -> { route_table_id: string, routes: [ { DestinationCidrBlock: string, GatewayId: string } ] }
		type FakeRoute struct {
			DestinationCidrBlock string
			GatewayId            string
		}
		type FakeRouteTable struct {
			RouteTableId string
			Routes       []FakeRoute
		}
		fakeRouteTables := make(map[string]FakeRouteTable)
		for subnetID, raw := range routeTablesRaw {
			m := raw.(map[string]interface{})
			rtID := m["route_table_id"].(string)
			routesRaw := m["routes"].([]interface{})
			var routes []FakeRoute
			for _, r := range routesRaw {
				rMap := r.(map[string]interface{})
				routes = append(routes, FakeRoute{
					DestinationCidrBlock: rMap["DestinationCidrBlock"].(string),
					GatewayId:            rMap["GatewayId"].(string),
				})
			}
			fakeRouteTables[subnetID] = FakeRouteTable{
				RouteTableId: rtID,
				Routes:       routes,
			}
		}

		t.Run(testName, func(t *testing.T) {
			patch := monkey.PatchInstanceMethod(reflect.TypeOf(new(ec2.Client)), "DescribeRouteTables",
				func(c *ec2.Client, ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error) {
					var subnetID string
					for _, filter := range input.Filters {
						if filter.Name != nil && *filter.Name == "association.subnet-id" {
							if len(filter.Values) > 0 {
								subnetID = filter.Values[0]
							}
							break
						}
					}
					fakeRT := fakeRouteTables[subnetID]
					var routes []ec2types.Route
					for _, r := range fakeRT.Routes {
						routes = append(routes, ec2types.Route{
							DestinationCidrBlock: aws.String(r.DestinationCidrBlock),
							GatewayId:            aws.String(r.GatewayId),
						})
					}
					rt := ec2types.RouteTable{
						RouteTableId: aws.String(fakeRT.RouteTableId),
						Routes:       routes,
					}
					return &ec2.DescribeRouteTablesOutput{
						RouteTables: []ec2types.RouteTable{rt},
					}, nil
				})
			defer patch.Unpatch()

			cfg := aws.Config{}
			result := security.DataplanePrivateCheck(eksCluster, cfg)
			if result.Passed != expectedPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectedPass, result.Passed)
			}
		})
	}
}
