# SCL-006 노드 확장/축소 정책 적용

## Meaning
EKS에서 Cluster Autoscaler가 자동으로 노드 확장 및 축소 정책을 설정하고, 이를 통해 클러스터 리소스를 효율적으로 관리합니다. 여기서 중요한 것은 최소/최대 노드 수가 적절히 설정되었는지 확인하는 것입니다.

## Impact
- 비효율적인 리소스 사용: 최소/최대 값이 적절하지 않으면, 리소스 부족이나 과잉 사용이 발생할 수 있음.

## Diagnosis
Cluster Autoscaler가 설정된 최소/최대 노드 수를 확인합니다

```bash
kubectl describe deployment cluster-autoscaler -n kube-system | grep "minNodes\|maxNodes"
```

## Mitigation
최소/최대 노드 수를 적절히 조정하여 리소스를 최적화합니다.

```yaml
apiVersion: autoscaling/v1
kind: AutoScalingGroup
spec:
  minSize: 3
  maxSize: 5
```