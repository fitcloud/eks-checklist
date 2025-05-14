# NET-003 VPC CNI의 Prefix 모드 사용

## Meaning
Amazon VPC CNI는 EC2 네트워크 인터페이스(ENI)에 네트워크 Prefix를 할당하여 노드당 파드 밀도를 높이고, 사용 가능한 IP 주소 수를 증가시키는 기능을 제공합니다. 이 모드를 활성화하면 보조 IP 주소 대신 IPv4 및 IPv6 CIDR을 할당할 수 있어, 더 많은 IP 주소를 파드에 할당할 수 있습니다.

## Impact
- 파드 밀도 제한
    - Linux 노드의 보조 IP 주소 한계: VPC CNI는 각 파드에 보조 IP 주소를 할당합니다. Linux에서는 여러 보조 IP 주소를 할당할 수 있지만, /28 CIDR과 같은 Prefix를 할당하지 않으면, 특정 인스턴스에서 할당할 수 있는 IP 주소 수가 제한됩니다. 
    - Windows 노드의 IP 슬롯 부족: Windows 노드는 단일 ENI와 한정된 슬롯만 지원합니다. 보조 IP만 사용할 경우, Windows 노드에서 실행 가능한 파드 수가 제한되며, Prefix 모드를 사용하지 않으면 이 문제는 더 심각해질 수 있습니다.
- IP 주소 관리 비효율성 : 보조 IP는 각 파드에 개별적으로 할당되므로, 많은 파드를 관리하기 어렵고 비효율적입니다. Prefix 모드를 사용하면 여러 파드에 IP 범위를 효율적으로 분배할 수 있습니다.

## Diagnosis
VPC CNI의 Prefix모드가 활성화 되어 있는지 확인하세요

```bash
# aws-node DaemonSet 정보를 통해 확인합니다.
kubectl get daemonsets.apps aws-node -n kube-system -o yaml

# ENABLE_PREFIX_DELEGATION 설정 여부를 쉽게 확인합니다.
if [[ "$(kubectl get daemonset aws-node -n kube-system -o jsonpath='{.spec.template.spec.containers[?(@.name=="aws-node")].env[?(@.name=="ENABLE_PREFIX_DELEGATION")].value}')" == "true" ]]; then
  echo "✅ ENABLE_PREFIX_DELEGATION이 설정되어 있습니다."
else
  echo "❌ ENABLE_PREFIX_DELEGATION이 설정되어 있지 않습니다."
fi
```

## Mitigation
이 기능에 문제가 발생했을 때 적용할 수 있는 완화책을 설명합니다.
```bash
# 1. 워커 노드에 할당된 Prefix 확인
aws ec2 describe-instances --query 'Reservations[*].Instances[].{InstanceId: InstanceId, Prefixes: NetworkInterfaces[].Ipv4Prefixes[]}'

# 2. Prefix Delegation 활성화
kubectl set env daemonset aws-node -n kube-system ENABLE_PREFIX_DELEGATION=true

# 3. 설정 적용 후, 다시 Prefix 할당 상태 확인
aws ec2 describe-instances --query 'Reservations[*].Instances[].{InstanceId: InstanceId, Prefixes: NetworkInterfaces[].Ipv4Prefixes[]}'
```

[AWS EKS Best Practice - Amazon VPC CNI](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/vpc-cni.html)
[AWS EKS Best Practice - linux Prefix mode](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/prefix-mode-linux.html)
[AWS EKS Best Practice - window Prefix mode](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/prefix-mode-win.html)