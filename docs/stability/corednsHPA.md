# REL-016 CoreDNS에 HPA 적용

## Meaning
CoreDNS는 클러스터 내 서비스 디스커버리(DNS) 및 이름 해석 서비스를 제공하는 핵심 구성요소입니다.
HPA(Horizontal Pod Autoscaler)를 적용하면 CoreDNS의 리소스 사용량(예: CPU, 메모리)에 따라 동적으로 Pod 수를 조정하여 성능 및 가용성을 향상시킬 수 있습니다.

## Impact
- 트래픽 증가 대응 불가: 클러스터 트래픽이 증가할 경우 DNS 요청 처리에 지연 발생.
- 서비스 가용성 저하: Pod 수가 고정되어 부하 증가 시 CoreDNS가 다운되면 클러스터 내 모든 서비스가 영향을 받을 수 있습니다
- 서비스 복구 지연: CoreDNS가 복구되기 전까지 클러스터의 정상 작동이 불가능

## Diagnosis
CoreDNS HPA 적용 상태를 확인합니다
```bash
kubectl get deployments -n kube-system --no-headers | awk '$1=="coredns"{print $1}' | while read name; do kubectl get hpa -n kube-system --no-headers | grep $name >/dev/null && echo "$name: HPA is configured" || echo "$name: HPA not configured"; done
```

## Mitigation
HPA 설정
CoreDNS에 HPA를 설정하여 트래픽 증가에 따라 Pod 수를 동적으로 확장합니다

**example**
```bash
kubectl autoscale deployment coredns -n kube-system --cpu-percent=50 --min=<min-number> --max=<max-number>
kubectl describe hpa coredns -n kube-system
```

[HPA 구성](https://kubernetes.io/ko/docs/tasks/run-application/horizontal-pod-autoscale/)
