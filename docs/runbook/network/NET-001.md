# NET-001 VPC 서브넷에 충분한 IP 대역대 확보

## Meaning
서브넷에 사용 가능한 IP가 부족하면, Pod이나 노드가 생성되지 않는 문제가 발생할 수 있습니다.

## Impact
- Pod Pending 상태 지속
- 노드 오토스케일 실패
- 전체 서비스 확장 불가

## Diagnosis
EKS 클러스터에서 사용할 수 있는 Subnet과 사용가능한 IP 갯수를 확인하세요

```bash
aws ec2 describe-subnets --subnet-ids $(aws eks describe-cluster --name <CLUSTER_NAME> --query "cluster.resourcesVpcConfig.subnetIds" --output text) --query 'Subnets[*].{ID:SubnetId, CIDR:CidrBlock, AvailableIPs:AvailableIpAddressCount}' --output table
```

## Mitigation
IP 부족 시, 서브넷을 추가하거나 VPC에 새로운 CIDR 블록을 추가하세요.

[EKS IP 최적화](https://docs.aws.amazon.com/eks/latest/best-practices/ip-opt.html)
[AWS VPC](https://docs.aws.amazon.com/ko_kr/vpc/latest/userguide/vpc-ip-addressing.html)
[Multiple CIDR ranges 사용](https://repost.aws/knowledge-center/eks-multiple-cidr-ranges)
