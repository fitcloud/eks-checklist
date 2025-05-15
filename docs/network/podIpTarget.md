# NET-006 ALB/NLB의 대상으로 Pod의 IP 사용

## Meaning
ALB/NLB가 트래픽을 EC2 인스턴스(Node) 대신 Pod IP로 직접 전달하도록 설정하면 불필요한 경유를 줄여 성능을 향상시킬 수 있습니다.

## Impact
- 트래픽이 NodePort를 거치면 Hop이 증가
- 상태 체크가 노드 단위로 동작하여 정확성 저하
- 모니터링 및 디버깅이 복잡해질 수 있음

## Diagnosis
존재하는 ALB 중 target-type이 ip가 아닌 리소스를 출력하세요

```bash
kubectl get ingress -A -o json | jq '.items[] | select(.metadata.annotations["alb.ingress.kubernetes.io/target-type"] != "ip")'
```

존재하는 NLB 중 target-type이 ip가 아닌 리소스를 출력하세요

```bash
kubectl get svc -A -o json | jq '.items[] | select(.metadata.annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"] != "ip")'
```


## Mitigation
Ingress 및 Service에 다음과 같은 Annotation을 명시적으로 설정해야 합니다.

**ALB (Ingress) 설정 예시**
```yaml
metadata:
  annotations:
    alb.ingress.kubernetes.io/target-type: ip
```

**NLB (Service) 설정 예시**
```yaml
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: ip
```
[AWS 모범사례](https://docs.aws.amazon.com/eks/latest/best-practices/load-balancing.html)
[AWS 공식문서-EKS와 NLB 사용법](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
[AWS Load Balancer Controller GitHub-ALB](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/#target-type)
[AWS Load Balancer Controller GitHub-NLB](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)
