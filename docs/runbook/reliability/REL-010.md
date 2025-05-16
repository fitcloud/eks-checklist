# REL-010 반영구 저장소에 애플리케이션 로그 저장

## Meaning
애플리케이션의 로그를 반영구 저장소(예: Openserach, CloudWatch Logs 등)에 저장하여 장기 보존 및 분석을 가능하게 합니다. 이를 통해 로그를 안전하게 보관하고, 필요시 추후 분석할 수 있습니다.

## Impact
- 로그 유실: 로컬 스토리지나 일시적인 저장소에 로그를 저장하면, 시스템 장애 시 로그가 유실될 위험이 있음
- 분석 불가능: 장기 보존된 로그를 활용한 문제 분석이나 리포팅이 어려워짐

## Diagnosis
현재 로그가 반영구 저장소에 저장되고 있는지 확인

```bash
kubectl logs <pod-name> --follow
```
로그가 제대로 반영구 저장소에 전송되지 않는다면, 설정을 점검해야 합니다.

## Mitigation
애플리케이션 로그를 CloudWatch Logs나 Opensearch와 같은 반영구 저장소에 전송하도록 설정합니다.

**Example (Fluent Bit 사용 시)**

```yaml
[OUTPUT]
    Name            es
    Match           *
    Region          your-region
    Total_File_Size 50M
```
