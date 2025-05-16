# REL-001 싱글톤 Pod 미사용

## Meaning
EKS에서 고가용성(High Availability)을 구현하는 것은 여러 인스턴스를 실행하여 서비스가 중단되지 않도록 하는 중요한 전략입니다. 고가용성을 구현하려면 StatefulSet, Deployment, 또는 ReplicaSet을 사용하여 **복제본(replica)**을 여러 개 설정하는 방법이 있습니다. 이 방식은 애플리케이션의 가용성을 보장하고, 하나의 Pod가 실패했을 때 다른 Pod가 그 역할을 대신할 수 있게 해줍니다.

## Impact
- 단일 실패 지점: 싱글톤 Pod는 클러스터 내에서 단 하나만 실행되므로, 만약 이 Pod가 장애를 겪거나 다운되면 애플리케이션이 완전히 중단될 수 있습니다.
- 스케일링의 제한: 싱글톤 Pod는 기본적으로 단 하나만 실행되기 때문에, 트래픽이 급격히 증가하는 상황에서 자동으로 스케일링하여 대응하기 어렵습니다.
- 업그레이드 및 배포 어려움: 업그레이드 및 배포에서 어려움이 있을 수 있습니다. 예를 들어, 싱글톤 Pod의 새로운 버전으로의 업데이트가 필요할 때, 다운타임을 최소화하기 위한 배포 전략을 수립하는 데 어려움이 따를 수 있습니다

## Diagnosis
싱글톤 Pod가 존재하는지 확인하세요.

```bash
kubectl get pods --all-namespaces -o json | jq -r '.items[] | select(.metadata.ownerReferences | length == 0) | "namespace: \(.metadata.namespace) | podName: \(.metadata.name)"'
```

## Mitigation
싱글톤 Pod를 제거하고 구성을 적용하세요

[EKS Singleton Pod 실행 방지](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/application.html#_recommendations)

[Kubernetes 워크로드 리소스](https://kubernetes.io/ko/docs/concepts/workloads/controllers/)