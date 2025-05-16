# REL-015 PV 사용시 volume affinity 위반 사항 체크

## Meaning
Volume Affinity는 Pod가 특정 **Persistent Volume (PV)**과 연결될 수 있도록 제한하는 설정입니다. 이 설정을 통해, 특정 노드에서만 PV가 사용되도록 강제할 수 있으며, 이를 통해 성능 최적화 및 리소스 분리를 할 수 있습니다. Volume Affinity 위반은 PV가 잘못된 노드에서 사용되거나, 요구하는 규칙을 충족하지 않는 경우 발생할 수 있습니다.

## Impact
- 성능 저하: 잘못된 노드에서 PV가 사용되면 성능 이슈나 장애가 발생할 수 있음
- 리소스 충돌: 여러 Pod가 동일한 리소스를 사용하게 되어 충돌이 발생할 수 있음

## Diagnosis
Volume Affinity 규칙이 제대로 적용되고 있는지 확인합니다. PV에 **nodeAffinity**가 설정되어 있는지 확인합니다.

```bash
kubectl get pv <pv-name> -o=jsonpath='{.spec.nodeAffinity}'
```

## Mitigation
Volume Affinity가 위반된 경우, **nodeAffinity**를 올바르게 설정하여 특정 노드에만 PV가 바인딩되도록 합니다.

**example**

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: example-pv
spec:
  capacity:
    storage: 10Gi
  volumeMode: Filesystem
  persistentVolumeReclaimPolicy: Retain
  accessModes:
    - ReadWriteOnce
  storageClassName: standard
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - node-1
                - node-2
```

