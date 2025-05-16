# SEC-008 비정상 접근에 대한 알림 설정

## Meaning
운영 중인 클러스터 및 AWS 계정에 대한 **비정상 접근 시도(Unauthorized Access, Credential Compromise 등)**를 감지하고, 이를 실시간으로 알림받는 체계를 갖추는 것은 보안 운영의 핵심입니다. 특히 EKS 환경에서는 API 서버 접근, IAM 오용, 악성 IP 접속 등 다양한 위협이 존재할 수 있습니다.

AWS에서는 CloudTrail, CloudWatch, Security Hub, SNS 등을 연계하여, 비정상 이벤트 탐지와 자동 알림이 가능합니다.

## Impact
- 사전 탐지 실패: 비정상 행위를 사전에 탐지하지 못하면 침해 사고로 이어질 수 있음
- 알림 부재: 이상 징후를 포착했더라도 알림 시스템이 없으면 대응 불가
- 감사 추적 불가: 장기적으로 어떤 위험이 있었는지 로그 기반 분석 불가능

## Diagnosis

다음 항목을 기반으로 비정상 접근 탐지 및 알림 설정 여부를 확인합니다:
- CloudTrail: 모든 이벤트 기록 여부, S3 저장 여부
- CloudWatch 알람 구성 여부: 특정 API 호출, 실패 로그 등에 대한 알람 구성
- SNS 또는 ChatOps 연계 여부(Slack, 이메일 등)

명령어 또는 콘솔 진단 예시:

```bash
# CloudTrail 로그 조회
aws cloudtrail lookup-events --lookup-attributes AttributeKey=EventName,AttributeValue=ConsoleLogin

# CloudWatch 알람 조회
aws cloudwatch describe-alarms

# SNS 구독 확인
aws sns list-subscriptions
```

## Mitigation
아래와 같은 방식으로 비정상 접근 탐지 및 알림 체계를 구성합니다.

CloudTrail 설정
모든 API 호출을 기록하고, S3 + CloudWatch Logs에 연동하여 지속적인 분석이 가능하도록 구성

CloudWatch 알람 설정
로그인 실패, EC2 권한 변경, EKS 관련 설정 변경 감지

**Before**

이상행위 탐지 시 수동 확인 필요

**After**

CloudTrail + CloudWatch Logs + SNS로 실시간 탐지 및 알림 구성


