# GEN-002 GitOps 적용

## Meaning
GitOps는 인프라 및 애플리케이션 상태를 Git 저장소에 선언적으로 정의하고, 이를 기반으로 클러스터의 상태를 지속적으로 동기화하는 운영 방식입니다. Git을 단일 진실 소스로 사용하여, 변경 이력 추적, 검토 프로세스 통합(PR 기반 배포), 자동화된 롤백 등을 가능하게 합니다. ArgoCD, FluxCD 등이 대표적인 GitOps 도구입니다.

EKS 환경에서는 GitOps를 통해 클러스터에 배포되는 애플리케이션 구성 및 설정을 안정적으로 관리할 수 있습니다.

## Impact
- 수동 배포 오류: 수동 배포나 스크립트 실행 오류로 인한 환경 불일치 발생 가능
- 추적 불가능한 변경: 배포 이력 미기록 및 협업 구조 부재
- 롤백 어려움: 문제 발생 시 이전 상태로 쉽게 복구 어려움
- 보안/컴플라이언스 미준수: 배포 프로세스에 승인 및 감사 기능이 누락될 수 있음

## Diagnosis
GitOps 적용 여부를 아래 항목으로 확인합니다

ArgoCD, FluxCD 등 GitOps 도구 설치 및 사용 여부 (argocd, flux CLI 등)

Git 리포지토리에 kustomization.yaml, helm-release.yaml 등 선언적 리소스 정의 존재 여부

클러스터 상태와 Git 리포 상태 비교 시 불일치 여부 (argocd app diff 등)

**example**

```bash
# ArgoCD 설치 확인
kubectl get pods -n argocd

# 등록된 GitOps 앱 확인
argocd app list
```

## Mitigation

애플리케이션 구성을 Git 저장소에 선언적으로 정의합니다.

**구성 요소 example**

Git 저장소: 선언적 리소스 정의 (Deployment, Service, HelmRelease 등)

GitOps 도구: Git 상태와 클러스터 상태를 지속적으로 비교 및 동기화

PR 기반 배포: Git 변경 승인 → 자동 동기화

**Before**

수동 kubectl apply 또는 CI 스크립트로 배포

**After**

Git PR 머지 → ArgoCD 또는 FluxCD가 자동 배포 및 Sync

참고 링크
[ArgoCD 공식 문서](https://argo-cd.readthedocs.io/en/stable/)
[FluxCD 공식 문서](https://fluxcd.io/)