package security_test

// import (
// 	"context"
// 	"fmt"
// 	"reflect"
// 	"testing"

// 	"eks-checklist/cmd/security"
// 	"eks-checklist/cmd/testutils"

// 	"bou.ke/monkey"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/ec2"
// 	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
// 	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
// )

// func TestDataplanePrivateCheck(t *testing.T) {
// 	// YAML 파일 "dataplane_private_check.yaml"에서 테스트 케이스 로드
// 	testCases := testutils.LoadTestCases(t, "dataplane_private.yaml")
// 	for _, tc := range testCases {
// 		testName, ok := tc["name"].(string)
// 		if !ok {
// 			t.Fatalf("Test case missing 'name' field")
// 		}

// 		expectedFailureVal, ok := tc["expected_failure"]
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'expected_failure' field", testName)
// 		}
// 		expectedFailure, ok := expectedFailureVal.(bool)
// 		if !ok {
// 			t.Fatalf("Test case '%s': expected_failure is not a bool", testName)
// 		}

// 		// cluster 정보 구성
// 		clusterMap, ok := tc["cluster"].(map[string]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'cluster' field", testName)
// 		}
// 		resVpcCfg, ok := clusterMap["resourcesVpcConfig"].(map[string]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'resourcesVpcConfig' field in cluster", testName)
// 		}
// 		subnetIdsRaw, ok := resVpcCfg["subnetIds"].([]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'subnetIds' field in resourcesVpcConfig", testName)
// 		}
// 		var subnetIds []string
// 		for _, id := range subnetIdsRaw {
// 			s, ok := id.(string)
// 			if !ok {
// 				t.Fatalf("Test case '%s': subnetId is not string", testName)
// 			}
// 			subnetIds = append(subnetIds, s)
// 		}
// 		eksCluster := security.EksCluster{
// 			Cluster: &eksTypes.Cluster{
// 				ResourcesVpcConfig: &eksTypes.VpcConfigResponse{
// 					SubnetIds: subnetIds,
// 				},
// 			},
// 		}

// 		// route_tables 정보를 구성
// 		routeTablesRaw, ok := tc["route_tables"].(map[string]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'route_tables' field", testName)
// 		}
// 		// route_tables: mapping: subnet_id -> { route_table_id: string, routes: [ { DestinationCidrBlock: string, GatewayId: string } ] }
// 		type FakeRoute struct {
// 			DestinationCidrBlock string
// 			GatewayId            string
// 		}
// 		type FakeRouteTable struct {
// 			RouteTableId string
// 			Routes       []FakeRoute
// 		}
// 		fakeRouteTables := make(map[string]FakeRouteTable)
// 		for subnetID, raw := range routeTablesRaw {
// 			m, ok := raw.(map[string]interface{})
// 			if !ok {
// 				t.Fatalf("Test case '%s': route_tables entry for %s is not a map", testName, subnetID)
// 			}
// 			rtID, ok := m["route_table_id"].(string)
// 			if !ok {
// 				t.Fatalf("Test case '%s': route_table_id missing or not a string for subnet %s", testName, subnetID)
// 			}
// 			routesRaw, ok := m["routes"].([]interface{})
// 			if !ok {
// 				t.Fatalf("Test case '%s': routes missing for subnet %s", testName, subnetID)
// 			}
// 			var routes []FakeRoute
// 			for _, r := range routesRaw {
// 				rMap, ok := r.(map[string]interface{})
// 				if !ok {
// 					t.Fatalf("Test case '%s': route is not a map", testName)
// 				}
// 				dest, ok := rMap["DestinationCidrBlock"].(string)
// 				if !ok {
// 					t.Fatalf("Test case '%s': DestinationCidrBlock missing or not string", testName)
// 				}
// 				gw, ok := rMap["GatewayId"].(string)
// 				if !ok {
// 					t.Fatalf("Test case '%s': GatewayId missing or not string", testName)
// 				}
// 				routes = append(routes, FakeRoute{
// 					DestinationCidrBlock: dest,
// 					GatewayId:            gw,
// 				})
// 			}
// 			fakeRouteTables[subnetID] = FakeRouteTable{
// 				RouteTableId: rtID,
// 				Routes:       routes,
// 			}
// 		}

// 		t.Run(testName, func(t *testing.T) {
// 			// Patch DescribeRouteTables on ec2.Client to return fake route table data.
// 			patch := monkey.PatchInstanceMethod(reflect.TypeOf(new(ec2.Client)), "DescribeRouteTables",
// 				func(c *ec2.Client, ctx context.Context, input *ec2.DescribeRouteTablesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error) {
// 					// AWS SDK v2에서는 SubnetIds가 없으므로, 필터를 통해 찾아야 합니다.
// 					var subnetID string
// 					for _, filter := range input.Filters {
// 						if filter.Name != nil && *filter.Name == "association.subnet-id" {
// 							if len(filter.Values) > 0 {
// 								subnetID = filter.Values[0]
// 							}
// 							break
// 						}
// 					}
// 					if subnetID == "" {
// 						return nil, fmt.Errorf("no subnet id provided")
// 					}
// 					fakeRT, ok := fakeRouteTables[subnetID]
// 					if !ok {
// 						return nil, fmt.Errorf("no fake route table data for subnet %s", subnetID)
// 					}
// 					var routes []ec2types.Route
// 					for _, r := range fakeRT.Routes {
// 						routes = append(routes, ec2types.Route{
// 							DestinationCidrBlock: aws.String(r.DestinationCidrBlock),
// 							GatewayId:            aws.String(r.GatewayId),
// 						})
// 					}
// 					rt := ec2types.RouteTable{
// 						RouteTableId: aws.String(fakeRT.RouteTableId),
// 						Routes:       routes,
// 					}
// 					return &ec2.DescribeRouteTablesOutput{
// 						RouteTables: []ec2types.RouteTable{rt},
// 					}, nil
// 				})
// 			defer patch.Unpatch()

// 			// 더미 AWS 구성(내용은 필요 없음)
// 			cfg := aws.Config{}

// 			result := security.DataplanePrivateCheck(eksCluster, cfg)
// 			// 원래 함수는 []string(IGW에 연결된 서브넷 ID 목록)을 반환합니다.
// 			// expected_failure가 true이면 공개 서브넷이 존재해야 하므로 result가 non-empty여야 하고,
// 			// expected_failure가 false이면 result가 empty여야 합니다.
// 			if expectedFailure {
// 				if len(result) == 0 {
// 					t.Errorf("Test '%s' failed: expected failure (non-empty public subnets) but got empty result", testName)
// 				}
// 			} else {
// 				if len(result) != 0 {
// 					t.Errorf("Test '%s' failed: expected success (empty result) but got %v", testName, result)
// 				}
// 			}
// 		})
// 	}
// }
