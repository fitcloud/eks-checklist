# REL-004 HPA(Horizontal Pod Autoscaler) 적용

## Meaning
HPA(Horizontal Pod Autoscaler)는 Kubernetes Deployment의 리소스 사용량(CPU, 메모리 등)을 기반으로 Pod 수를 자동으로 조정하는 기능입니다. Deployment에 HPA를 적용하면 트래픽 변동에 유연하게 대응할 수 있으며, 안정성과 리소스 효율성을 동시에 확보할 수 있습니다.

## Impact
- 트래픽 증가 대응 불가 - 트래픽이 급증하면 고정된 Pod 수로 인해 애플리케이션 성능이 저하되거나 요청이 실패할 수 있습니다.
- 리소스 낭비 - 트래픽 감소 시에도 Pod 수가 고정되어 불필요한 리소스를 소비, 비용이 증가합니다.
- 서비스 가용성 저하 - 과도한 부하로 인해 Pod가 다운되면 서비스 중단 가능성이 높아지고 복구 작업이 지연될 수 있습니다.
- 운영 관리 어려움 - 관리자가 트래픽 상황을 수동으로 모니터링하고 조정해야 하므로 운영 복잡성이 증가합니다.

## Diagnosis

모든 Deployment 및 HPA 확인

```bash
kubectl get deployments -A --no-headers | awk '{print $1, $2}' | while read namespace name; do kubectl get hpa -n $namespace --no-headers | grep $name >/dev/null || echo "$namespace/$name: HPA not configured"; done
```

## Mitigation

**example**

HPA 적용 - Deployment에 HPA를 설정하여 autosacle를 설정합니다
```bash
kubectl autoscale deployment <deployment-name> --cpu-percent=50 --min=<min-number> --max=<max-number>
```

[HPA 구성](https://kubernetes.io/ko/docs/tasks/run-application/horizontal-pod-autoscale/)