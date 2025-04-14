# **ALB/NLB의 대상으로 Pod의 IP 사용**

## Meaning
AWS의 ALB(Application Load Balancer)와 NLB(Network Load Balancer)는 Kubernetes에서 외부 트래픽을 수신하고 이를 클러스터 내부의 파드로 라우팅하는 데 사용됩니다.
기본적으로 AWS LoadBalancer는 트래픽을 노드(EC2 인스턴스)에 전달하지만, EKS에서는 Pod의 IP를 직접 대상으로 사용(IP 모드)하도록 설정할 수 있습니다.
- ALB: alb.ingress.kubernetes.io/target-type = ip
- NLB: service.beta.kubernetes.io/aws-load-balancer-nlb-target-type = ip
이 설정을 통해 Kubernetes 파드가 직접 로드밸런서의 대상이 되게 설정이 가능합니다.


## Impact
기본적으로 AWS 로드 밸런서 컨트롤러는 "인스턴스" 유형을 사용하여 대상을 등록하며, 이 대상은 작업자 노드의 IP와 노드 포트가 됩니다. 이에 따른 영향은 다음과 같습니다.

- 로드 밸런서의 트래픽은 NodePort의 Worker 노드로 전달되고, 이는 iptables 규칙(노드에서 실행되는 kube-proxy에서 구성)에 의해 처리되고, ClusterIP(여전히 노드에 있음)의 서비스로 전달되고, 마지막으로 서비스는 등록된 포드를 무작위로 선택하여 트래픽을 전달합니다. 이 흐름에는 여러 홉이 포함되고, 특히 서비스가 다른 AZ에 있을 수 있는 다른 워커 노드에서 실행되는 포드를 선택하는 경우가 있으므로 추가 지연이 발생할 수 있습니다.
- 로드 밸런서는 워커 노드를 대상으로 등록하므로 대상으로 전송되는 상태 확인은 포드에서 직접 수신되지 않고 워커 노드의 NodePort에서 수신되며 상태 확인 트래픽은 위에서 설명한 것과 동일한 경로를 따릅니다.
- 모니터링과 문제 해결은 좀 더 복잡합니다. 로드 밸런서에서 전달된 트래픽이 포드로 직접 전송되지 않고, 워커 노드에서 수신된 패킷을 서비스 클러스터 IP와 포드에 신중하게 연관시켜야 패킷 경로에 대한 완벽한 엔드투엔드 가시성을 확보해 적절한 문제 해결을 수행할 수 있습니다.

### [AWS 모범사례](https://docs.aws.amazon.com/eks/latest/best-practices/load-balancing.html)

## Diagnosis
아래 명령어를 실행하면, ALB 또는 NLB가 Pod의 IP를 대상으로 하고 있지 않은 것을 확인할 수 있습니다.

```bash
# Ingress 점검 (ALB)
kubectl get ingress --all-namespaces -o json | jq -r '
  .items[]
  | {
      ns: .metadata.namespace,
      name: .metadata.name,
      targetType: (.metadata.annotations["alb.ingress.kubernetes.io/target-type"] // "")
    }
  | select(.targetType == "" or .targetType == "instance")
  | "Ingress: \(.ns)/\(.name) | target-type: \(.targetType)"
' 

# Service 점검 (NLB)
kubectl get svc --all-namespaces -o json | jq -r '
  .items[]
  | select(
      .metadata.ownerReferences != null and
      ([.metadata.ownerReferences[] | select(.kind == "Ingress")] | length > 0)
    )
  | {
      ns: .metadata.namespace,
      name: .metadata.name,
      targetType: (.metadata.annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"] // "")
    }
  | select(.targetType != "ip")
  | "Service: \(.ns)/\(.name) | target-type: \(.targetType)"
' 
```
위 명령어는 다음 조건을 만족하는 리소스를 출력합니다.
- ALB Ingress가 target-type: instance 이거나 미설정
- NLB Service가 target-type: ip가 아닌 경우

## Mitigation
Ingress 및 Service에 다음과 같은 Annotation을 명시적으로 설정해야 합니다.

### ALB (Ingress) 설정 예시
```yaml
metadata:
  annotations:
    alb.ingress.kubernetes.io/target-type: ip
```

### NLB (Service) 설정 예시
```yaml
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: ip
  ownerReferences:
  - apiVersion: networking.k8s.io/v1
      kind: Ingress
      name: nginx1-ingress
      uid: 9fd0a07d-95ff-465a-8d21-fd2a84881ceb
```

### [AWS 공식문서-EKS와 NLB 사용법](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
### [AWS Load Balancer Controller GitHub-ALB](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/#target-type)
### [AWS Load Balancer Controller GitHub-NLB](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)
