# REL-006 중요 워크로드에 대한 PDB(Pod Distruption Budget) 적용

## Meaning
**Pod Disruption Budget (PDB)**는 주요 서비스의 가용성을 보장하기 위해 Pod이 예기치 않게 종료되지 않도록 설정하는 정책입니다. 이를 통해 롤링 업데이트, 노드 종료 등에서 안정적인 서비스 운영을 유지할 수 있습니다

## Impact
- 중요 서비스 중단: PDB가 설정되지 않으면, 중요 워크로드가 자동으로 종료되어 서비스 중단이나 장애가 발생할 수 있음

- 불필요한 다운타임: 예기치 않은 Pod 종료로 인한 다운타임 발생 가능성

## Diagnosis
**PDB가 설정된 워크로드를 확인**

```bash
kubectl get poddisruptionbudgets -A
```

PDB가 없거나 적절히 설정되지 않은 경우, 중요 워크로드에 PDB를 설정해야 합니다.

## Mitigation
중요 워크로드에 대해 PDB를 설정하여, 최소한의 Pod 가용성을 유지하도록 합니다.

**example**

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: critical-app-pdb
  namespace: default
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: critical-app
```

[Kubenetes Pods distruption](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)