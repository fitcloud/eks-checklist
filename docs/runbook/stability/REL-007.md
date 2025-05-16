# REL-007 애플리케이션에 적절한 CPU/RAM 할당

## Meaning
애플리케이션에 적절한 CPU와 RAM 리소스를 할당하여 성능과 비용 효율을 최적화합니다. 적정 자원 할당을 통해 리소스 과잉 또는 부족을 방지하고, 클러스터 안정성을 유지할 수 있습니다.

## Impact
- 리소스 부족: RAM나 CPU 부족 시 애플리케이션 성능 저하 또는 OOM(Out of Memory) 발생 가능
- 리소스 과잉: 불필요한 리소스 낭비로 비용 상승

## Diagnosis
애플리케이션의 현재 CPU와 RAM 할당이 적절한지 확인

```bash
kubectl get deployment <deployment-name> -o=jsonpath='{.spec.template.spec.containers[*].resources}'
```

## Mitigation
각 애플리케이션에 대해 **적절한 리소스 요청(Request)과 제한(Limit)**을 설정하여, 리소스 과잉/부족 문제를 해결합니다.

**example**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-app
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: example-container
          image: example-image
          resources:
            requests:
              memory: "500Mi"
              cpu: "500m"
            limits:
              memory: "1Gi"
              cpu: "1"
```
