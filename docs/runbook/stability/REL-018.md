# REL-018 Karpenter 사용시 DaemonSet에 Priority Class 부여

## Meaning
Karpenter는 자동으로 노드를 생성하는데, DaemonSet에 PriorityClass를 설정하지 않으면 필수 컴포넌트가 새 노드에 누락될 수 있습니다.

## Impact
- DaemonSet 누락 가능성: PriorityClass가 없으면 Karpenter가 DaemonSet을 고려하지 않습니다.
- 서비스 안정성 저하: 필수 시스템 컴포넌트가 배치되지 않아 안정성에 영향을 받습니다.

## Diagnosis
다음 명령어는 Karpenter 사용 여부 검사는 포함되어 있지 않습니다. 
PriorityClass가 지정되지 않은 DaemonSet을 다음 명령어로 진단할 수 있습니다.

```bash
kubectl get daemonsets --all-namespaces -o json | jq -r '
  .items[] | select(.spec.template.spec.priorityClassName == null) |
  "Namespace: \(.metadata.namespace) | DaemonSet: \(.metadata.name)"'
```

## Mitigation
중요 DaemonSet에 priorityClassName을 추가하세요.

**example**
```yaml
spec:
  priorityClassName: system-node-critical
```

[Kubernetes 공식 문서 - DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset)
[Kubernetes 공식 문서 - pod-priority-preemption](https://kubernetes.io/ko/docs/concepts/scheduling-eviction/pod-priority-preemption)