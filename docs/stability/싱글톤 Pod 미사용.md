# **싱글톤 Pod 미사용 (Singleton Pod Usage Detection)**

## Meaning
싱글톤 Pod(Singleton Pod)는 클러스터 내에서 단 하나의 인스턴스만 실행되는 Pod입니다.  
이 방식은 **고가용성(HA) 및 확장성**을 저해하며, 장애 발생 시 서비스 중단을 초래할 수 있습니다.

## Impact
- 노드 장애 시 **서비스 중단 발생**
- **자동 확장(HPA) 불가능**, 트래픽 증가 대응 불가
- 롤링 업데이트 적용 불가 → 다운타임 증가
- 특정 노드에 종속될 경우 **복구 지연 발생**

## Diagnosis
싱글톤 Pod가 존재하는지 확인하세요.

```bash
kubectl get deployments,statefulsets -A | awk '$3 == 1 {print}'
kubectl get pods -A --field-selector=status.phase=Running | grep -v "ReplicaSet"
kubectl get hpa -A
kubectl describe pod <pod-name> | grep -i "nodeName"
```

## Mitigation
싱글톤 Pod를 제거하고 HA 구성을 적용하세요.

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

- HPA 적용
```bash
kubectl autoscale deployment my-app --cpu-percent=50 --min=2 --max=10
```

- 수동 실행된 Singleton Pod 제거
```bash
kubectl delete pod <pod-name> -n <namespace>
```