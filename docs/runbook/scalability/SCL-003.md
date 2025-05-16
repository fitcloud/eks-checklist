# SCL-003 Spot 노드 사용시 Spot 중지 핸들러 적용

## Meaning
Spot 인스턴스는 AWS에서 저비용으로 제공되지만, 언제든지 회수(interruption)될 수 있어 사전에 종료 알림을 감지하고 애플리케이션을 안전하게 종료·이동시키는 핸들러가 필요합니다.
이 점검은 클러스터 내에 “termination-handler”라는 이름이 포함된 Pod(Spot Termination Handler)가 배포되어 있는지를 확인합니다.

## Impact
- 예기치 않은 중단: Spot 인스턴스 종료 시 알림을 받지 못하면 Pod가 강제 종료되어, 서비스 장애나 데이터 손실이 발생할 수 있습니다.
- Graceful Shutdown 실패: 종료 전 처리가 이루어지지 않아 요청이 유실되고, 애플리케이션 무결성이 훼손될 수 있습니다.
- 운영 복잡성 증가: 수동으로 Spot 종료를 감지·대응해야 하므로 운영 부담이 커지고, 자동 확장/복구 프로세스가 비효율적으로 작동합니다.

## Diagnosis
클러스터 전체 네임스페이스에서 “termination-handler” 문자열을 포함한 Pod가 있는지 확인

```bash
kubectl get pods -A --no-headers | awk '{print $1, $2}' | grep termination-handler
```

## Mitigation
Spot Termination Handler를 설치하여 Spot 인스턴스 종료 알림을 자동으로 처리하도록 설정하세요.
```bash
kubectl apply -f https://github.com/aws/aws-node-termination-handler/releases/latest/download/aws-node-termination-handler.yaml
```

[AWS Node Termination handler github](https://github.com/aws/aws-node-termination-handler)