# REL-018 Karpenter 사용시 DaemonSet에 Priority Class 부여

## Meaning
Karpenter는 필요에 따라 노드를 자동으로 생성하거나 제거해주는 AWS 기반의 노드 프로비저너입니다.
이때 새로 생성된 노드에 어떤 Pod을 우선적으로 스케줄링할지 판단할 기준이 필요한데,
DaemonSet에 PriorityClass가 명시되어 있지 않으면 Karpenter가 해당 워크로드를 고려하지 못할 수 있습니다.
따라서 DaemonSet에 PriorityClass를 부여하면, 노드 스케줄링 및 자원 부족 상황에서 DaemonSet이 더 안정적으로 배치될 수 있도록 보장할 수 있습니다.
(system-node-critical과 system-cluster-critical은 쿠버네티스에서 기본적으로 제공하는 PriorityClass입니다.)

## Impact
- 노드 생성 직후 DaemonSet 누락 가능성: PriorityClass가 없으면, Karpenter는 해당 DaemonSet의 스케줄링을 고려하지 않고 노드를 생성하여, 필수 데몬이 새 노드에 누락될 수 있습니다.

- 서비스 안정성 저하: CNI, 로그 수집기, 보안 에이전트 등 필수 시스템 컴포넌트가 노드에 배치되지 않으면, 전체 클러스터 안정성 및 관측 가능성이 저하됩니다.

- 스케일링 불안정: 노드의 리소스 요청이 과소 평가되어, 과소 프로비저닝이 발생할 수 있습니다.

## Diagnosis
다음 명령어는 Karpenter 사용 여부 검사는 포함되어 있지 않습니다. 
PriorityClass가 지정되지 않은 DaemonSet을 다음 명령어로 진단할 수 있습니다.

```bash
kubectl get daemonsets --all-namespaces -o json | jq -r '
  .items[] | select(.spec.template.spec.priorityClassName == null) |
  "Namespace: \(.metadata.namespace) | DaemonSet: \(.metadata.name) (PriorityClass 미설정)"
'
```
**출력 예시**
- Namespace: kube-system | DaemonSet: log-agent (PriorityClass 미설정)
PriorityClass가 설정되지 않은 것은 다음과 같이 출력되게 됩니다. 
PriorityClass가 설정되어 있다면 출력되지 않을 것 입니다.

## Mitigation
DaemonSet에 priorityClassName을 명시합니다.
example
```yaml
spec:
  priorityClassName: high-priority
```

사용자 정의 PriorityClass가 필요한 경우 다음과 같이 정의할 수 있습니다.
- value : 값이 클수록, 우선순위가 높습니다.
- globalDefault : PriorityClass가 설정되어 있지 않은 리소스에 기본 적용합니다.
- preemptionPolicy : Never로 설정된 경우, 해당 프라이어리티클래스의 파드는 비-선점입니다.

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: high-priority
value: 1000000
globalDefault: false
description: "우선순위가 높은 DaemonSet용 PriorityClass"
```
[Kubernetes 공식 문서 - DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset)
[Kubernetes 공식 문서 - pod-priority-preemption](https://kubernetes.io/ko/docs/concepts/scheduling-eviction/pod-priority-preemption)