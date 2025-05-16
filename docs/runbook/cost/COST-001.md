# COST-001 kubecost 설치

## Meaning
Kubecost는 클러스터의 비용 최적화 및 리소스 사용량 분석을 도와주는 도구입니다. 설치가 되어 있지 않으면, 리소스 효율성을 확보하거나 비용을 투명하게 관리하는 것이 어려울 수 있습니다.

## Impact
- 비용 추적 부족 : 쿠버네티스 상의 비용 데이터에 대한 실시간 가시성을 제공하지 못해 예산 초과 가능성이 높아질 수 있습니다.
- 리소스 최적화 실패 : 모니터링을 하지 않으면 워크로드별로 리소스를 최적화하지 못해, CPU와 메모리가 과소 또는 과잉 프로비저닝될 위험이 있습니다.
- 미래 예측 및 계획의 어려움 : EKS 내부 자체 비용 예측 및 리포팅 기능이 없어, 사용량 증가나 워크로드 확장에 대한 사전 계획이 어렵습니다.
- 비용 측정 구조의 어려움 : pod, node, namespace 및 각 기타 항목별 비용 측정을 하기 어렵습니다.

## Diagnosis
Kubecost가 설치되어 있는지 확인하세요

```bash
kubectl get deployments -n kubecost | grep "kubecost" >/dev/null && echo "Kubecost is installed" || echo "Kubecost is not installed"
```

## Mitigation
Kubecost 설치
Helm을 사용하여 Kubecost 설치 및 확인 진행

```bash
helm repo add kubecost https://kubecost.github.io/cost-analyzer/
helm repo update
helm install kubecost kubecost/cost-analyzer --namespace kubecost --create-namespace
kubectl get deployments -n kubecost | grep "kubecost"
```
[Kubecost 설치](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/cost-monitoring-kubecost.html)
