# NET-007 Pod Readiness Gate 적용

## Meaning
PodReadinessGate는 Kubernetes에서 Pod의 준비 상태를 세밀하게 관리할 수 있는 기능입니다. 이 기능을 사용하면 Pod의 준비 상태가 특정 조건을 만족할 때까지 Pod를 Ready 상태로 간주하지 않도록 설정할 수 있습니다.이 기능을 서비스가 실제로 트래픽을 받을 준비가 되었을 때만 요청을 처리하도록 보장할 수 있습니다

## Impact
- 서비스 장애: Pod가 실제로 준비되지 않았는데도 Ready 상태로 표시되어 트래픽을 받을 수도 있습니다.준비되지 않은 상태에서 요청을 받아버리면 서비스 장애가 발생 할 수 있습니다.
- 배포 문제: PodReadinessGate를 사용하지 않으면, 준비되지 않은 Pod가 배포 과정에서 정상적으로 트래픽을 처리한다고 잘못 판단되어, 배포가 완료되었다고 표시될 수 있습니다

## Diagnosis
PodReadinessGate가 적용되지 않은 Pod가 있는지 확인하세요
```bash
kubectl get namespaces --show-labels | awk 'NR==1 {print "NAME", "READINESSGATE"; 
next} {if ($4 ~ /elbv2.k8s.aws\/pod-readiness-gate-inject=enabled/) print $1, "Enabled"; else print $1, "Disabled"}'
```

## Mitigation
PodReadinessGate를 적용하세요

**example**
```bash
$ kubectl create namespace readiness
namespace/readiness created

$ kubectl label namespace readiness elbv2.k8s.aws/pod-readiness-gate-inject=enabled
namespace/readiness labeled

$ kubectl describe namespace readiness
Name:         readiness
Labels:       elbv2.k8s.aws/pod-readiness-gate-inject=enabled
Annotations:  <none>
Status:       Active

```
[AWS-load-balancer-controller PodReadinessgate 문서](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.1/deploy/pod_readiness_gate/)

