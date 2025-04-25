# **Subnet ip capacity 최적화**

## Meaning
VPC(Virtual Private Cloud) 서브넷에 충분한 IP 대역이 없을 경우, 여러 가지 문제가 발생할 수 있습니다. 이러한 문제들은 VPC 내의 리소스들이 IP 주소를 할당받을 수 없게 되어, 네트워크 통신에 문제가 생기거나 새로운 리소스를 배포할 수 없게 되는 상황을 초래할 수 있습니다.

## Impact
- iP 주소 부족: 문제: VPC 내의 서브넷에 할당된 IP 주소 범위가 적다면면, 서브넷 내의 인스턴스나 기타 리소스가 IP 주소를 할당받지 못해 새로운 리소스를 생성하거나 통신이 불가능해질 수 있습니다.

## Diagnosis
EKS 클러스터에서 사용할 수 있는 Subnet과 사용가능한 IP 갯수를 확인하세요

```bash
aws eks describe-cluster --name <CLUSTER_NAME> --query "cluster.resourcesVpcConfig.subnetIds" --output text | tr '\t' '\n' | xargs aws ec2 describe-subnets --subnet-ids | jq -r '.Subnets[] | [.SubnetId, .CidrBlock, .AvailableIpAddressCount] | @tsv' | column -t -s $'\t' -N "Name, CIDR Block, Available IPs"
```

## Mitigation
사용가능한 IP 대역을 추가하세요
CIDR을 추가할 때, 이미 다른 곳에서 사용 중이지 않은 대역을 선택해 주세요 (하기 링크 참조)

### [EKS IP 최적화](https://docs.aws.amazon.com/eks/latest/best-practices/ip-opt.html)
### [Multiple CIDR ranges 사용](https://repost.aws/knowledge-center/eks-multiple-cidr-ranges)
