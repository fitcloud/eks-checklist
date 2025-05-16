# REL-002 2개 이상의 Pod 복제본 사용

## Meaning
ReplicaSet은 항상 지정된 수의 Pod가 실행되도록 보장하며, 최소 2개 이상을 구성해야 고가용성 및 자동 복구가 가능합니다.

## Impact
- 단일 Pod 장애 시 서비스 전체 중단
- 무중단 배포 불가
- 장애 복구 및 확장성 저하

## Diagnosis
다음 명령어로 복제본이 1개 이하인 Deployment를 확인하세요

```bash
kubectl get deployments -A -o json | jq -r '.items[] | select(.spec.replicas <= 1) | "\(.metadata.namespace) | \(.metadata.name) | \(.spec.replicas)"'
```

## Mitigation
Deployment의 .spec.replicas 값을 2 이상으로 수정하세요

```yaml
spec:
  replicas: 2
```
[Kubernetes Replicaset](https://kubernetes.io/ko/docs/concepts/workloads/controllers/replicaset/)