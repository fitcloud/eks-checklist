# REL-008 애플리케이션 중요도에 따른 QoS 적용

## Meaning
**QoS(Quality of Service)**는 애플리케이션의 우선순위에 따라 Kubernetes 리소스 요청(Request)과 제한(Limit)을 설정하여, 리소스 경쟁 상황에서 중요한 애플리케이션이 적절한 리소스를 우선적으로 할당받도록 합니다.
Kubernetes에서는 Guaranteed, Burstable, BestEffort 3가지 QoS 클래스를 제공합니다.

## Impact
- 리소스 우선순위: 중요한 애플리케이션이 리소스 부족 상황에서 우선적으로 자원을 할당받지 않으면 성능 저하나 장애가 발생할 수 있음
- 자원 낭비: 과도한 리소스를 할당하면 비효율적인 비용 사용이 발생할 수 있음

## Diagnosis
현재 QoS 상태를 확인하여 애플리케이션 중요도에 맞는 설정이 되어 있는지 확인합니다.

```bash
kubectl get pod <pod-name> -o=jsonpath='{.status.qosClass}'
```

## Mitigation
애플리케이션의 중요도에 따라 리소스 requests와 limits를 적절히 설정하여 QoS 클래스를 조정합니다.

- Guaranteed: requests와 limits가 동일할 때, 해당 Pod는 Guaranteed 클래스를 가집니다.

- Burstable: requests가 limits보다 작을 때, Burstable 클래스를 가집니다.

- BestEffort: requests와 limits가 모두 설정되지 않으면 BestEffort로 분류됩니다.

**example**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: important-app
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: app-container
          image: app-image
          resources:
            requests:
              memory: "1Gi"
              cpu: "500m"
            limits:
              memory: "2Gi"
              cpu: "1"
```

[Kubenetes Pods Quality of Service](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/)