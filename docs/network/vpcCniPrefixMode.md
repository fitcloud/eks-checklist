# NET-003 VPC CNI의 Prefix 모드 사용

## Meaning
Prefix 모드를 활성화하면 ENI에 IPv4 프리픽스 블록을 할당하여, 노드당 더 많은 Pod에 IP를 부여할 수 있습니다.

## Impact
- Prefix 모드를 사용하지 않으면 노드당 파드 수 제한
- Windows 노드의 경우 특히 심각한 IP 부족 현상 발생 가능.

## Diagnosis
VPC CNI의 Prefix모드가 활성화 되어 있는지 확인하세요

```bash
kubectl get ds aws-node -n kube-system -o jsonpath='{.spec.template.spec.containers[?(@.name=="aws-node")].env[?(@.name=="ENABLE_PREFIX_DELEGATION")].value}'
```
출력 값이 "true"이면 Prefix 모드가 활성화된 것입니다.

## Mitigation
관리형 애드온 및 자체관리형 애드온을 설치하세요
```bash
kubectl set env ds aws-node -n kube-system ENABLE_PREFIX_DELEGATION=true
```
관리형 애드온 및 자체관리형 애드온인 경우 설치 시 'ENABLE_PREFIX_DELEGATION=true' 인수를 넘깁니다

[AWS EKS Best Practice - Amazon VPC CNI](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/vpc-cni.html)
[AWS EKS Best Practice - linux Prefix mode](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/prefix-mode-linux.html)
[AWS EKS Best Practice - window Prefix mode](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/prefix-mode-win.html)