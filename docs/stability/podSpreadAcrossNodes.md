# REL-003 동일한 역할을 하는 Pod를 다수의 노드에 분산 배포

## Meaning
Pod가 특정 노드/Zone에 몰리지 않도록 topologySpreadConstraints 또는 podAntiAffinity 설정을 통해 분산 배포해야 합니다. 이는 노드/Zone 장애 발생 시 서비스의 고가용성을 확보하는 데 필수입니다.

## Impact
- 하나의 노드 장애로 전체 서비스 중단 위험
- 리소스 불균형 및 낭비
- Zone 장애에 대한 복원력 부족

## Diagnosis
아래 명령어로 affinity 및 spread 설정이 없는 Pod를 탐지할 수 있습니다.

```bash
kubectl get pods -A -o json | jq -r '
.items[] | select(
  (.spec.affinity | not) and 
  ((.spec.topologySpreadConstraints // []) | length == 0)
) | "Namespace: \(.metadata.namespace) | Pod: \(.metadata.name) - 분산 설정 없음"'
```

## Mitigation
다음 중 하나 이상 설정을 추가하세요.

**TopologySpreadConstraints 사용**
```yaml
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: kubernetes.io/hostname
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        app: your-app
```
maxSkew: 허용 가능한 분산 불균형 정도

**PodAntiAffinity 설정**
```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchLabels:
            app: your-app
        topologyKey: kubernetes.io/hostname
```
topologyKey를 기반으로 동일한 app을 가진 Pod가 같은 노드에 배치되지 않도록 함

[Kubernetes - Assigning Pods to Nodes](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/)
[Kubernetes - Pod Topology Spread Constraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/)