# REL-014 다수의 가용 영역에 데이터 플레인 노드 배포

## Meaning
Amazon EKS 클러스터에서 데이터 플레인(Data Plane) 노드를 여러 가용 영역(Availability Zone)에 분산 배포하면, 클러스터의 고가용성과 장애 복구 능력이 향상됩니다. 단일 AZ에 문제가 발생해도, 다른 AZ에 배포된 노드가 서비스를 지속할 수 있어 서비스 중단을 최소화할 수 있습니다.

## Impact
- AZ 장애 시 서비스 중단 위험: 모든 노드가 하나의 가용 영역에 있을 경우, 해당 AZ 장애 시 전체 클러스터 운영에 문제가 생길 수 있습니다.
- 자동 복구 제한: EKS의 자동 복구 기능은 여러 AZ에 걸쳐 노드를 구성했을 때 더 효과적으로 동작합니다.

## Diagnosis
다수의 가용 영역에 데이터 플레인 노드가 배포되어 있는지 확인하세요

```bash
kubectl describe node | grep "topology.kubernetes.io/zone"
```

## Mitigation
다수의 가용 영역에 데이터 플레인 노드를 배포하세요

**example**

nodegroup 사용 시
```bash
eksctl create nodegroup \
  --cluster <Cluster_Name> \
  --name <Nodegroup_Name> \
  --node-type <Instance_Type> \
  --nodes 3 \
  --zones ap-northeast-2a,ap-northeast-2b,ap-northeast-2c
```

Karpenter 사용 시
```yaml
apiVersion: karpenter.sh/v1beta1
kind: NodePool
metadata:
  name: Nodepool
spec:
  template:
    spec:
      requirements:
        - key: topology.kubernetes.io/zone
          operator: In
          values: ["ap-northeast-2a", "ap-northeast-2b", "ap-northeast-2c"]
```

[EKS 모범사례 Multiple AZ](https://docs.aws.amazon.com/eks/latest/best-practices/data-plane.html)
[Karpenter 공식 문서](https://karpenter.sh/docs/concepts/nodepools/)