# REL-003 동일한 역할을 하는 Pod를 다수의 노드에 분산 배포

## Meaning
Kubernetes는 스케줄링 시, 동일한 레이블을 가진 Pod가 클러스터 내 특정 노드 또는 영역에 집중되지 않도록 분산 배치하는 기능을 제공합니다. 이를 통해 노드 장애 또는 가용 영역(Zone) 장애 시에도 애플리케이션이 고가용성을 유지할 수 있도록 보장합니다.
이를 구현하려면 다음 설정 중 하나 이상을 반드시 포함해야 합니다.
- topologySpreadConstraints: 지정된 topologyKey(예: hostname, zone)에 따라 Pod가 균등하게 분산되도록 강제
- podAntiAffinity: 동일한 조건의 Pod가 동일한 노드 또는 영역에 함께 스케줄링되지 않도록 제약


## Impact
- 노드 장애 시 전체 서비스 중단 가능성: 동일한 노드에 모든 Pod가 배치되면, 노드 하나만 장애가 발생해도 서비스 전체가 중단될 수 있습니다.

- 리소스 편중과 불균형: 클러스터 리소스가 특정 노드에만 집중될 경우, 다른 노드의 자원이 비효율적으로 낭비될 수 있습니다.

- 애플리케이션 복원력 저하: 가용 영역을 고려한 분산 배치가 없을 경우, Zone 장애에 대한 내성이 떨어집니다.

## Diagnosis
다음 명령어를 통해 affinity와 topologySpreadConstraints 설정이 적절히 구성되지 않은 Pod를 탐지할 수 있습니다.

```bash
kubectl get pods --all-namespaces -o json | jq -r '
  .items[]
  | . as $pod
  | ($pod.spec.affinity != null) as $hasAffinity
  | ($pod.spec.topologySpreadConstraints // []) as $tscList
  | ($tscList | length > 0) as $tscExists
  | ($tscList
      | map(
          select(
            .maxSkew != null and
            (.maxSkew | tonumber) > 1
          )
        )
    ) as $badSkewList
  | ($tscExists and ($badSkewList | length == 0)) as $topologyValid
  | ($badSkewList | length > 0) as $hasBadSkew
  |
  if $hasBadSkew then
    $badSkewList[]
    | "Namespace: \($pod.metadata.namespace) | Pod: \($pod.metadata.name) - maxSkew 값이 \(.maxSkew) (1 초과)"
  else empty end,
  if ($hasAffinity | not) and ($topologyValid | not) then
    "Namespace: \($pod.metadata.namespace) | Pod: \($pod.metadata.name) - affinity와 유효한 topologySpreadConstraints 설정이 모두 없음"
  else empty end
'
```

## Mitigation
Pod가 클러스터에 균등하게 배치되도록 하기 위한 설정을 추가해야 합니다.
example
**TopologySpreadConstraints 사용**
```bash
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
spec:
    topologySpreadConstraints:
        - maxSkew: 1
          topologyKey: kubernetes.io/hostname
          whenUnsatisfiable: DoNotSchedule
          matchLabelKeys:
            - app
            - pod-template-hash
```
- maxSkew: 허용 가능한 분산 불균형 정도

**PodAntiAffinity 설정**
```bash
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-cache
spec:
  selector:
    matchLabels:
      app: store
  replicas: 3
  template:
    metadata:
      labels:
        app: store
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - store
            topologyKey: "kubernetes.io/hostname"
      containers:
      - name: redis-server
        image: redis:3.2-alpine
```
- topologyKey를 기반으로 동일한 app을 가진 Pod가 같은 노드에 배치되지 않도록 함

[Kubernetes - Assigning Pods to Nodes](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/)
[Kubernetes - Pod Topology Spread Constraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/)