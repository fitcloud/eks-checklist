# SCL-004 중요 Pod에 노드 삭제 방지용 Label 부여

## Meaning
중요한 워크로드가 올라간 노드가 Karpenter나 Cluster Autoscaler에 의해 자동으로 삭제되지 않도록, 해당 노드에 삭제 방지용 label을 부여합니다.

## Impact
- 중요 서비스가 의도치 않게 중단될 수 있음
- 복구 지연 및 운영 리스크 증가

## Diagnosis
중요 Pod가 실행 중인 노드에 삭제 방지용 label이 설정되어 있는지 확인

```bash
kubectl get nodes --show-labels | grep do-not-disrupt
```

## Mitigation
중요 노드에 karpenter.sh/do-not-disrupt=true label 추가

nodepool 또는 Cluster Autoscaler 설정에서 해당 노드는 스케일 인 대상 제외 처리

```bash
kubectl label node <NODE_NAME> karpenter.sh/do-not-disrupt: "true"
```
