# SCL-007 다양한 인스턴스 타입 사용

## Meaning
현재 사용중인 Node의 Instace type을 확인하여 단일 Instacne type 사용 시 다양한 Instance type 사용 변경으로 애플리케이션의 요구 사항에 맞게 최적화하고, 클러스터의 리소스를 효율적으로 활용 합니다.

## Impact
- 성능 최적화 부족: 각 애플리케이션은 고유한 성능 요구 사항을 가지기 때문에, 인스턴스 타입을 하나로 고정하면 애플리케이션의 요구 사항을 충족시키기 어려워 성능 최적화에 제한이 생깁니다.
- 비용 비효율성: 고성능 인스턴스 타입을 사용하게 되면, 고성능이 필요하지 않은 애플리케이션도 동일한 인스턴스를 사용하게 되어 불필요한 자원 소비가 발생하고, 이로 인해 비용이 증가하게 됩니다.

## Diagnosis
현재 사용중인 Instance Type을 확인하세요

```bash
kubectl get nodes -o custom-columns="NODE_NAME:.metadata.name,INSTANCE_TYPE:.metadata.labels.beta\.kubernetes\.io/instance-type"
```

**출력 example**
```bash
NODE_NAME                                         INSTANCE_TYPE
ip-10-20-12-101.ap-northeast-2.compute.internal   t3.medium
ip-10-20-13-65.ap-northeast-2.compute.internal    t3.medium
```

## Mitigation
다양한 Instance Type을 사용하세요

[최적의 Amazon EC2 Node Instance 유형 선택](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/choosing-instance-type.html)

[EKS Fargate](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/fargate.html)

[AWS 관리형 노드 그룹](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/create-managed-node-group.html)

[Karpenter Nodepool](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/karpenter.html)