# REL-002-2개 이상의 Pod 복제본 사용

## Meaning
ReplicaSet은 Kubernetes에서 Pod의 복제본을 관리하는 리소스입니다. ReplicaSet은 지정된 수의 Pod 복제본이 항상 실행되도록 보장합니다.

## Impact
- 고가용성: 복제본이 2개 이상인 경우, 하나의 Pod가 장애가 발생해도 다른 Pod가 트래픽을 처리할 수 있습니다. 이를 통해 애플리케이션의 가용성이 향상됩니다.
- 내결함성: ReplicaSet은 **내결함성(fault tolerance)**을 보장하는 중요한 역할을 합니다. Pod가 장애를 겪을 경우, 클러스터 내 다른 Pod가 이를 대신 처리하게 되어, 장애 발생 시에도 애플리케이션이 정상적으로 작동할 수 있도록 합니다
- 확장성: Pod 중 하나가 실패하거나 삭제되면, ReplicaSet은 자동으로 새로운 Pod를 생성하여 지정된 수의 복제본을 유지합니다. 이를 통해 자동 복구가 이루어집니다. 이 기능은 관리자가 수동으로 Pod를 복구할 필요 없이 시스템을 안정적으로 유지할 수 있게 해줍니다.
- 업데이트 관리: ReplicaSet은 애플리케이션의 새로운 버전을 배포할 때 Rolling Update 전략을 사용하여 기존 Pod를 점진적으로 교체합니다. 이를 통해 서비스 중단 없이 무중단 배포가 가능합니다

## Diagnosis
2개 이상 pod 복제본을 사용하는지 확인하세요

```bash
#전체 replicaset이 적용된 deployments 출력
kubectl get deployments --all-namespaces -o custom-columns="NAMESPACE:.metadata.namespace, DEPLOYMENT:.metadata.name, REPLICAS:.spec.replicas"
#1개 이하의 replicaset이 적용된 deployments 출력
kubectl get deployments --all-namespaces -o json | jq -r '.items[] | select(.spec.replicas <= 1) | "\(.metadata.namespace) | \(.metadata.name) | \(.spec.replicas)"'
```

## Mitigation
단일 Pod를 제거하고 2개 이상의 replicaset 구성을 적용하세요.

- Deployment로 변환 후 replicas: 2 이상 설정
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: app-container
          image: my-app-image
```
[Kubernetes Replicaset](https://kubernetes.io/ko/docs/concepts/workloads/controllers/replicaset/)