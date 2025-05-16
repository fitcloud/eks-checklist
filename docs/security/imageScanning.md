# SEC-013 컨테이너 이미지 정적 분석

## Meaning
컨테이너 이미지에는 종종 수많은 패키지와 라이브러리가 포함되며, 이 중 일부는 **보안 취약점(CVE)**를 내포할 수 있습니다.
**정적 분석(Static Analysis)**은 이미지를 실행하기 전, 이미지 내부의 패키지 및 구성 파일을 분석하여 알려진 보안 취약점을 식별하는 프로세스입니다.

이 과정을 통해 다음을 방지할 수 있습니다

1. 이미지 빌드 시점부터의 취약점 포함

2. 운영 환경으로의 보안 위협 전파

3. 파이프라인의 보안 공백
## Impact
- 취약한 이미지 배포: CVE 포함 이미지가 실제 환경에 배포될 수 있음
- 공격 벡터 증가: 미패치된 취약점으로 인해 공격에 노출될 가능성
- 감사/보안 테스트 실패: 기업 보안 정책 및 외부 감사 기준 미충족

## Diagnosis
컨테이너 이미지에 대해 정적 분석이 수행되고 있는지 확인합니다

1. CI/CD 파이프라인에서 이미지 스캐너가 포함되어 있는지 확인합니다
2. GitHub Actions, GitLab CI, ArgoCD 등에서 정적 분석 단계 확인
3. 취약점 레벨 기준(Low/Medium/High/Critical)에 따른 차단 정책 존재 여부

대표적인 도구
- Trivy (Aqua Security)
- Amazon ECR 이미지 스캐닝

**example Trivy CLI로 이미지 스캔**
```bash
trivy image nginx:1.25.2
```

result
```text
nginx:1.25.2 (debian 11)

Total: 12 (UNKNOWN: 0, LOW: 4, MEDIUM: 5, HIGH: 2, CRITICAL: 1)
```

**Amazon ECR 사용 시**

```bash
aws ecr describe-image-scan-findings \
  --repository-name my-repo \
  --image-id imageTag=latest
```
스캔 결과는 High, Critical 여부에 따라 알림 또는 배포 차단 설정 가능

## Mitigation
정적 분석 도구 도입 및 자동화

1. 스캔 결과 기반 정책 수립
2. 정기적 리포트 생성 (예: 주간 기준)
3. 이미지 최적화 및 최소화(Alpine, Distroless 이미지 활용) - 사용하지 않는 패키지 제거로 공격 표면 축소

**Before**

아무런 검사 없이 이미지 배포
운영 중 취약점 발견 후 긴급 대응

**After**

빌드 시 자동 취약점 검사
배포 전 차단 정책으로 보안 선제 대응

[Container 이미지 Scan](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/windows-images.html)