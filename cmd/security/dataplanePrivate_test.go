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
	"github.com/aws/aws-sdk-go-v2/service/eks"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestDataplanePrivateCheck(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "dataplane_private.yaml")

	type FakeRoute struct {
		DestinationCidrBlock string
		GatewayId            string
	}
	type FakeRouteTable struct {
		RouteTableId string
		Routes       []FakeRoute
	}

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		// ── 1) YAML에서 cluster.name, vpcId, subnetIds 읽어오기
		clusterMap := tc["cluster"].(map[string]interface{})
		clusterName := clusterMap["name"].(string)
		resVpc := clusterMap["resourcesVpcConfig"].(map[string]interface{})
		vpcId := resVpc["vpcId"].(string)
		rawSubnets := resVpc["subnetIds"].([]interface{})
		var subnetIds []string
		for _, s := range rawSubnets {
			subnetIds = append(subnetIds, s.(string))
		}

		// ── 2) EksCluster 객체 생성 (Name, VpcId, SubnetIds 모두 세팅)
		eksCluster := security.EksCluster{
			Cluster: &eksTypes.Cluster{
				Name: aws.String(clusterName),
				ResourcesVpcConfig: &eksTypes.VpcConfigResponse{
					VpcId:     aws.String(vpcId),
					SubnetIds: subnetIds,
				},
			},
		}

		// ── 3) route_tables 페이크 데이터 준비
		routeTablesRaw := tc["route_tables"].(map[string]interface{})
		fakeRouteTables := make(map[string]FakeRouteTable)
		for sid, raw := range routeTablesRaw {
			mm := raw.(map[string]interface{})
			rtID := mm["route_table_id"].(string)
			routesRaw := mm["routes"].([]interface{})
			var fr FakeRouteTable
			fr.RouteTableId = rtID
			for _, rr := range routesRaw {
				rmap := rr.(map[string]interface{})
				fr.Routes = append(fr.Routes, FakeRoute{
					DestinationCidrBlock: rmap["DestinationCidrBlock"].(string),
					GatewayId:            rmap["GatewayId"].(string),
				})
			}
			fakeRouteTables[sid] = fr
		}

		t.Run(testName, func(t *testing.T) {
			// ── 4) EKS 클라이언트 NewFromConfig 패치
			var eksClient *eks.Client
			p1 := monkey.Patch(eks.NewFromConfig,
				func(cfg aws.Config, opts ...func(*eks.Options)) *eks.Client {
					eksClient = &eks.Client{}
					return eksClient
				},
			)
			defer p1.Unpatch()

			// ── 5) ListNodegroups 패치
			p2 := monkey.PatchInstanceMethod(
				reflect.TypeOf(eksClient), "ListNodegroups",
				func(_ *eks.Client, _ context.Context, _ *eks.ListNodegroupsInput, _ ...func(*eks.Options)) (*eks.ListNodegroupsOutput, error) {
					return &eks.ListNodegroupsOutput{Nodegroups: []string{clusterName}}, nil
				},
			)
			defer p2.Unpatch()

			// ── 6) DescribeNodegroup 패치 (포인터 타입 주의)
			p3 := monkey.PatchInstanceMethod(
				reflect.TypeOf(eksClient), "DescribeNodegroup",
				func(_ *eks.Client, _ context.Context, _ *eks.DescribeNodegroupInput, _ ...func(*eks.Options)) (*eks.DescribeNodegroupOutput, error) {
					return &eks.DescribeNodegroupOutput{
						Nodegroup: &eksTypes.Nodegroup{Subnets: subnetIds},
					}, nil
				},
			)
			defer p3.Unpatch()

			// ── 7) EC2 클라이언트 NewFromConfig 패치
			var ec2Client *ec2.Client
			p4 := monkey.Patch(ec2.NewFromConfig,
				func(cfg aws.Config, opts ...func(*ec2.Options)) *ec2.Client {
					ec2Client = &ec2.Client{}
					return ec2Client
				},
			)
			defer p4.Unpatch()

			// ── 8) DescribeRouteTables 패치
			p5 := monkey.PatchInstanceMethod(
				reflect.TypeOf(ec2Client), "DescribeRouteTables",
				func(_ *ec2.Client, _ context.Context, in *ec2.DescribeRouteTablesInput, _ ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error) {
					var tables []ec2types.RouteTable
					for sid, frt := range fakeRouteTables {
						assoc := ec2types.RouteTableAssociation{SubnetId: aws.String(sid)}
						var routes []ec2types.Route
						for _, r := range frt.Routes {
							routes = append(routes, ec2types.Route{
								DestinationCidrBlock: aws.String(r.DestinationCidrBlock),
								GatewayId:            aws.String(r.GatewayId),
							})
						}
						tables = append(tables, ec2types.RouteTable{
							RouteTableId: aws.String(frt.RouteTableId),
							Associations: []ec2types.RouteTableAssociation{assoc},
							Routes:       routes,
						})
					}
					return &ec2.DescribeRouteTablesOutput{RouteTables: tables}, nil
				},
			)
			defer p5.Unpatch()

			// ── 9) 실제 함수 호출 및 검증
			got := security.DataplanePrivateCheck(eksCluster, aws.Config{})
			if got.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, got.Passed)
			}
		})
	}
}
