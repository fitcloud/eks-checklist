# **싱글톤 Pod 미사용(Singleton Pod Usage Detection)**

## Meaning
쿠버네티스 클러스터에서는 **고가용성(High Availability, HA)**과 **확장성(Scalability)**을 보장하기 위해 **싱글톤(Singleton) Pod 패턴을 지양해야 합니다**.  
싱글톤 Pod란 **클러스터 내에서 단 하나의 인스턴스만 실행되는 Pod**를 의미하며, 이는 장애 발생 시 서비스 중단을 초래할 가능성이 큽니다.

✅ **올바른 설계 원칙**  
- 클러스터 내에서 동일한 Pod가 여러 개 배포될 수 있도록 구성해야 합니다.  
- 단일 Pod 장애 시 즉시 복구되도록 **ReplicaSet, Deployment, HPA(Auto Scaling)** 등을 활용해야 합니다.  
- 특정 기능이 단일 Pod에 의존하지 않도록 **StatefulSet 또는 Multi-Instance 설계**를 고려해야 합니다.  

❌ **잘못된 패턴 예시 (Singleton Pod)**
- 특정 서비스가 단 하나의 Pod에서만 실행됨
- 수동으로 `kubectl run <pod>`으로 실행된 단일 Pod
- `replicas: 1`로 설정된 Deployment 또는 StatefulSet
- Singleton Pod 장애 시 자동 복구되지 않음

---

## Impact
EKS 클러스터에서 싱글톤 Pod를 사용하면 **다음과 같은 문제**를 유발할 수 있습니다.

🔴 **1. 서비스 가용성 저하**
- 싱글톤 Pod가 실행 중인 노드가 장애를 일으키면 서비스가 즉시 중단됩니다.
- EKS의 **자동 복구 기능이 적용되지 않아** 지속적인 운영이 어려워집니다.

🔴 **2. 확장성 문제**
- 하나의 Pod만 존재하기 때문에 트래픽 증가 시 **자동 확장이 불가능**합니다.
- Horizontal Pod Autoscaler(HPA)를 적용할 수 없어, 성능 문제가 발생할 가능성이 높습니다.

🔴 **3. 장애 복구 지연**
- 싱글톤 Pod가 실행 중인 노드가 장애를 일으키면, **다른 노드로 재배치하는 과정에서 서비스 중단**이 발생할 수 있습니다.
- `kubectl delete pod <pod>`를 수동으로 실행해야 하는 비효율적인 운영 방식이 될 수 있습니다.

🔴 **4. 롤링 업데이트 불가능**
- Deployment에서 여러 개의 Pod를 운영하면 **무중단 롤링 업데이트(Rolling Update)**가 가능하지만,  
  싱글톤 Pod는 **업데이트 시 반드시 다운타임이 발생**합니다.

---

## Diagnosis
EKS 클러스터에서 **싱글톤 Pod가 사용되고 있는지 점검하는 방법**은 다음과 같습니다.

1️⃣ **싱글톤 Pod 감지 (replicas가 1인 Deployment / StatefulSet 찾기)**
```bash
kubectl get deployments,statefulsets -A | awk '$3 == 1 {print}'
```
- replicas 값이 1로 설정된 Deployment 또는 StatefulSet을 찾습니다.

2️⃣ **수동으로 실행된 싱글톤 Pod 찾기**
```bash
kubectl get pods -A --field-selector=status.phase=Running | grep -v "ReplicaSet"
```
- ReplicaSet이 없는 Pod를 찾습니다.
- `kubectl run` 또는 `kubectl create pod` 명령어로 직접 생성된 Pod는 재시작 정책이 적용되지 않을 가능성이 높습니다.

3️⃣ **Auto Scaling이 적용되지 않은 Deployment 찾기**
```bash
kubectl get hpa -A
```
- Horizontal Pod Autoscaler(HPA)가 적용되지 않은 Deployment는 확장성 문제가 발생할 수 있습니다.

4️⃣ **싱글톤 Pod의 노드 종속성 확인**
```bash
kubectl describe pod <pod-name> | grep -i "nodeName"
```
- 특정 노드에 고정적으로 스케줄링된 싱글톤 Pod가 존재하는지 확인합니다.
- 노드 장애 시 재배포가 불가능한 경우 즉시 수정해야 합니다.

---

## Mitigation
EKS에서 싱글톤 Pod를 제거하고, 고가용성(HA) 및 확장성을 고려한 아키텍처로 변경하는 방법은 다음과 같습니다.

✅ 1. **Deployment**로 변환하여 `replicas: 2` 이상 설정 싱글톤 Pod가 있다면 **Deployment로 변환**하고 **replicas를 2개 이상으로 설정**하세요.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 3  # 최소 2개 이상으로 설정
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

✅ 2. **HPA(Horizontal Pod Autoscaler) 적용** 트래픽 증가 시 자동으로 확장되도록 설정하세요.
```bash
kubectl autoscale deployment my-app --cpu-percent=50 --min=2 --max=10
```

✅ 3. **PodDisruptionBudget(PDB) 설정** 싱글톤 Pod는 중단 시 장애가 발생할 가능성이 크므로, PodDisruptionBudget(PDB)를 활용하여 **최소한 하나 이상의 Pod가 항상 실행되도록** 보장하세요.
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: my-app-pdb
spec:
  minAvailable: 1  # 항상 최소 1개 이상 유지
  selector:
    matchLabels:
      app: my-app
```

✅ 4. **Pod Anti-Affinity 설정** Pod가 특정 노드에 종속되지 않도록 Anti-Affinity 정책을 설정하세요.
```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchLabels:
            app: my-app
        topologyKey: "kubernetes.io/hostname"

```

✅ 5. **ReplicaSet이 없는 단일 Pod 제거** 잘못 생성된 싱글톤 Pod를 제거하려면 다음 명령어를 실행하세요.
```bash
kubectl delete pod <pod-name> -n <namespace>
```

---

## Conclusion
✅ **EKS에서 싱글톤 Pod는 사용하는 것은 권장드리지 않습니다.**

✅ **고가용성(HA) 및 확장성을 위해 Deployment와 HPA를 활용하세요.**

✅ **replicas: 2 이상을 설정하여 장애 시에도 서비스가 지속되도록 구성하세요.**

✅ **수동으로 실행된 Singleton Pod를 감지하고, 자동화된 방식으로 운영하세요.**