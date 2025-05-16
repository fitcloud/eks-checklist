# SEC-012 데이터 플레인 사설망

## Meaning
Data plane를 Private 서브넷에 배치하면 인터넷에 대한 노출을 최소화하여 공격 위험을 줄일 수 있습니다. 만약 worker node를 공용 서브넷에 배치하기로 결정했다면, 제한적인 AWS 보안 그룹 규칙을 설정하여 노출을 최소화하는 것이 중요합니다 (비권장)

## Impact
- 보안 위험 증가: Public 서브넷에 배치된 worker node는 외부에서 직접 접근할 수 있게 되어, 악의적인 공격자에게 노출될 가능성이 높아집니다. 예를 들어, DDoS(분산 서비스 거부 공격)나 기타 인터넷 기반 공격을 받을 수 있습니다.
- 관리 복잡성 증가: public 서브넷에서 직접 접근할 수 있기 때문에, 이를 보호하기 위한 보안 그룹 및 네트워크 ACL 설정이 더 복잡해집니다. 이런 설정을 잘못 구성하면 보안 위험이 커질 수 있습니다.
- 공용 IP 사용에 따른 비용: public 서브넷에 배치된 노드는 공용 IP 주소를 사용해야 할 수 있으며, 이는 비용적으로 불리할 수 있습니다. 공용 IP를 계속 사용하는 데는 추가적인 비용이 발생할 수 있습니다.
- 침해 사고 발생 시 영향 범위 확대: worker node가 public 서브넷에 있을 경우, 침해 사고가 발생하면 해당 node뿐만 아니라 다른 리소스에도 영향을 미칠 수 있습니다. 내부 네트워크와의 경계를 명확히 하기가 어려워, 피해가 확산될 수 있습니다.

## Diagnosis
EKS 클러스터가 사용하는 Subnet이 Public or Private인지 확인하세요

``` bash
aws ec2 describe-route-tables --filters "Name=association.subnet-id,Values=$(aws eks describe-cluster --name eks-checklist --query 'cluster.resourcesVpcConfig.subnetIds' --output text | tr '\t' ',')" | jq -r '.RouteTables[] | {RouteTableId, Routes: .Routes[] | select(.DestinationCidrBlock == "0.0.0.0/0") } | {RouteTableId, DestinationCidrBlock: .Routes.DestinationCidrBlock, GatewayId: (if .Routes.GatewayId then .Routes.GatewayId else .Routes.NatGatewayId end), PublicStatus: (if .Routes.GatewayId and (.Routes.GatewayId | test("^igw-")) then "public" elif .Routes.NatGatewayId then "private" else "private" end)}'
```

**출력 example**
```bash
Internet Gateway 존재
  "DestinationCidrBlock": "0.0.0.0/0",
  "GatewayId": "igw-0123456789",
  "PublicStatus": "public"
 ------------------------------------------ 
Natgateway 존재
  "DestinationCidrBlock": "0.0.0.0/0",
  "GatewayId": "nat-0123456789",
  "PublicStatus": "private"
```

## Mitigation
Private Subnet을 사용하세요
Internet Gateway 존재 시 작업 진행

**example**
```bash
aws ec2 create-route --route-table-id <ROUTE-TABLE-ID> --destination-cidr-block 0.0.0.0/0 --nat-gateway-id <NAT-GATEWAY-ID>
```

[EKS 인프라 보안](https://docs.aws.amazon.com/eks/latest/best-practices/protecting-the-infrastructure.html)