# SCL-005 Application에 Graceful shutdown 적용

## Meaning
애플리케이션이 종료될 때 정상적으로 연결 종료, 상태 저장, 리소스 정리 등을 수행하도록 Graceful Shutdown 처리를 구현해야 합니다. Kubernetes에서는 SIGTERM 시그널 수신 후 terminationGracePeriodSeconds 동안 정리 시간이 주어집니다.

## Impact
- 처리 중 요청 유실
- DB 트랜잭션 오류
- 상태 불일치 등 장애 발생 가능

## Diagnosis
Deployment 또는 Pod 설정에서 terminationGracePeriodSeconds가 명시되어 있는지 확인

```bash
kubectl get deploy -o jsonpath='{.items[*].spec.template.spec.terminationGracePeriodSeconds}'
```

또한, 컨테이너 애플리케이션이 SIGTERM을 수신했을 때 shutdown 로직을 실행하는지 확인 필요 (코드 레벨).

## Mitigation
앱에서 SIGTERM 수신 시 정리 작업 수행 (process.on('SIGTERM'), trap, etc.)

Kubernetes에 terminationGracePeriodSeconds 설정 추가 (예: 30초)

```yaml
spec:
  terminationGracePeriodSeconds: 30
```

[Kubenetes Pod Lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)