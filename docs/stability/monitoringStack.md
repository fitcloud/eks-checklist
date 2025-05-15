# REL-009 인프라 및 애플리케이션 모니터링 스택 적용

## Meaning
kube-prometheus-stack는 Prometheus와 Grafana를 포함한 모니터링 솔루션으로, 인프라 및 애플리케이션 모니터링을 실시간으로 설정하고 관리할 수 있습니다.

## Impact
- 모니터링 부재: 시스템의 성능 문제나 장애를 빠르게 파악할 수 없으며, 다운타임이 증가할 수 있음

## Diagnosis
kube-prometheus-stack가 설치되어 있는지 확인

```bash
kubectl get pods -n monitoring -l release=kube-prometheus-stack
```

## Mitigation
kube-prometheus-stack가 설치되지 않았다면, Helm을 통해 설치합니다.

```bash
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack
```

[AWS EKS Prometheus를 사용한 클러스터 모니터링](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/prometheus.html)