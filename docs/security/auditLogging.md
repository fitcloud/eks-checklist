# SEC-007 Audit 로그 활성화

## Meaning
Audit Log는 Amazon EKS 클러스터에서 발생하는 모든 API 요청과 이벤트를 기록하는 기능입니다. 
Audit Log를 비활성화하면 클러스터의 보안 및 운영 상태를 추적하기 어려워질 수 있습니다.

## Impact
- 보안 취약점 증가: API 요청 기록이 없으므로 악의적인 활동을 탐지하기 어려움
- 문제 해결 지연: 이벤트 기록이 없어 장애 원인을 분석하기 어려움.
- 운영 투명성 부족: 클러스터 활동에 대한 명확한 기록이 없어 관리 효율성 저하

## Diagnosis
Audit Log가 비활성화되어 있는지 확인하려면 아래 명령어를 사용하거나 제공된 코드를 실행하세요

```bash
aws eks describe-cluster --name <클러스터 이름> --query "cluster.logging.clusterLogging[?types[?contains(@, 'audit')]].{AuditLog:enabled}" --output table
```

## Mitigation
EKS Cluster Audit log 활성화

```bash
aws eks update-cluster-config --region <region> --name <클러스터 이름> --logging '{"clusterLogging":[{"types":["audit"],"enabled":true}]}'
```

[EKS 감사 및 로깅](hhttps://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/auditing-and-logging.html)
[EKS Auditlog 설정 - AWS CLI 부분 참조](https://docs.aws.amazon.com/eks/latest/userguide/control-plane-logs.html#enabling-control-plane-log-export)
[Kubernetes Audit](https://kubernetes.io/docs/tasks/debug/debug-cluster/audit/)