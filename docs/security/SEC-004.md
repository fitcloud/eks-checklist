# SEC-004 데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여

## Meaning
데이터 플레인 노드(워커 노드)에 연결된 IAM 역할은 최소 권한 원칙(Least Privilege Principle)을 따라야 합니다.
기본적으로 다음 세 가지 정책만 부여되어야 하며, 그 외 정책은 보안상 위험 요소가 될 수 있습니다.

- AmazonEKSWorkerNodePolicy
- AmazonEKS_CNI_Policy
- AmazonEC2ContainerRegistryReadOnly

## Impact
- 노드 침해 시 과도한 권한 탈취 위험
- AWS 리소스 전체에 영향을 줄 수 있는 lateral movement 가능성
- 보안 및 규정 준수 기준 미달

## Diagnosis
노드 Role에 정책을 확인합니다.

```bash
aws iam list-attached-role-policies --role-name <노드 역할 이름> --query "AttachedPolicies[].PolicyName" --output text
```

## Mitigation
데이터 플레인 노드에 불필요하게 부여된 정책을 제거하세요.

불필요한 정책 제거
```bash
aws iam detach-role-policy \
  --role-name <MyNodeRole> \
  --policy-arn arn:aws:iam::aws:policy/<UnwantedPolicyName>
```
**반드시 AmazonEKSWorkerNodePolicy, AmazonEKS_CNI_Policy,AmazonEC2ContainerRegistryReadOnly 외의 정책은 제거해야 합니다.**

[Amazon EKS 데이터 플레인 권한](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/create-node-role.html)