# EKS-Checklist

**EKS-Checklist**는 Amazon EKS (Elastic Kubernetes Service) 클러스터의 설정과 상태를 자동으로 점검하여, 운영자가 클러스터를 **최적화**, **보안 강화**, **비용 절감**할 수 있도록 지원하는 도구입니다.

> 이 도구는 Go 언어로 작성되었으며, AWS SDK for Go, Kubernetes Go Client, 그리고 CLI 명령어 프레임워크인 Cobra를 활용하여 제작되었습니다.

---

## ✅ 프로젝트 목적

Amazon EKS는 관리형 Kubernetes 환경을 제공하지만, 실제 운영에서는 다음과 같은 문제가 자주 발생합니다:

- 리소스 과다 사용으로 인한 비용 증가
- 불완전한 IAM 설정으로 인한 보안 위험
- 잘못된 네트워크 구성으로 인한 서비스 장애
- 오토스케일링 누락 등으로 인한 확장성 부족

**EKS-Checklist**는 이와 같은 문제를 사전에 식별하고 개선할 수 있도록 설계된 도구입니다. AWS 및 Kubernetes의 **베스트 프랙티스(Best Practices)**를 기반으로 클러스터 상태를 분석하고, 다음과 같은 항목에 대한 진단을 제공합니다:

---

## 🔍 점검 항목

| 카테고리        | 설명 |
|----------------|------|
| **비용 최적화 (Cost)**     | 과도한 리소스 사용, 미사용 노드, 고비용 인스턴스 등을 탐지하고 절감 방안을 제시합니다. |
| **일반 설정 (General)**   | 클러스터 버전, 태그 구성, 메타데이터 등 기본 구성이 적절한지 확인합니다. |
| **네트워크 (Network)**    | VPC, Subnet, 보안 그룹, ENI, IP 할당 등을 분석합니다. |
| **확장성 (Scalability)**  | HPA, Cluster Autoscaler, 노드그룹 구성 등 확장성 관련 설정을 점검합니다. |
| **보안 (Security)**       | IAM 정책, 인증 구성, API 서버 접근 제어 등을 점검합니다. |
| **안정성 (Stability)**    | 로그, 모니터링, 백업 설정 등을 분석하여 클러스터 안정성 수준을 진단합니다. |

---

## 📋 요구 사항 (Prerequisites)

도구를 사용하기 위해 다음 환경이 준비되어 있어야 합니다:

1. **AWS CLI**
   - 설치: [공식 문서](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) 참고
   - 인증: `aws configure` 명령어로 설정 (Access Key, Secret, Region 등)

2. **kubectl**
   - 클러스터와 연결된 `kubeconfig`가 설정되어 있어야 함
   - [kubectl 설치 가이드](https://kubernetes.io/docs/tasks/tools/)

3. **EKS 클러스터 접근 권한**
   - IAM Role 또는 User가 EKS 클러스터 및 리소스에 접근 가능한 권한이 있어야 합니다.

---

## 📦 설치 방법

### 방법 1: GitHub Releases에서 바이너리 다운로드

1. GitHub의 [Releases 페이지](https://github.com/fitcloud/eks-checklist/releases)로 이동합니다.
2. 운영 체제에 맞는 바이너리 파일을 다운로드합니다:
   - macOS: `eks-checklist-darwin-amd64`
   - Linux: `eks-checklist-linux-amd64`
   - Windows: `eks-checklist-windows-amd64.exe`

---

## 💻 플랫폼별 설치 예시

### Linux

```bash
wget https://github.com/fitcloud/eks-checklist/releases/download/{version}/eks-checklist-linux-amd64
chmod +x eks-checklist-linux-amd64
sudo mv eks-checklist-linux-amd64 /usr/local/bin/eks-checklist
eks-checklist --profile my-aws-profile
```
## MacOS

```bash
curl -LO https://github.com/fitcloud/eks-checklist/releases/download/{version}/eks-checklist-darwin-amd64
chmod +x eks-checklist-darwin-amd64
sudo mv eks-checklist-darwin-amd64 /usr/local/bin/eks-checklist
eks-checklist --profile my-aws-profile
```
## Window

1. .exe 파일을 다운로드하여 예: C:\Program Files\EKS-Checklist\에 저장합니다.
2. 명령 프롬프트 또는 PowerShell에서 다음과 같이 실행합니다:

```bash
cd "C:\Program Files\EKS-Checklist\"
eks-checklist-windows-amd64.exe --profile my-aws-profile
```

## 🚀 사용 방법

### 기본 사용 예시
```bash
eks-checklist --context my-cluster --profile dev --output text --out all
```
### 주요 옵션 설명

| 옵션                | 설명 |
|---------------------|------|
| `--context`         | 사용할 kubeconfig context 이름 |
| `--kubeconfig`      | kubeconfig 파일 경로 (기본: 사용자 홈 디렉토리 경로) |
| `--profile`         | 사용할 AWS CLI 프로파일 이름 |
| `--output`          | 출력 형식 지정 (`text`, `html`) |
| `--out`             | 결과 필터링 옵션 (`all`, `pass`, `fail`, `manual`) |
| `--sort`            | 결과를 상태별 정렬 (`pass`, `fail`, `manual`) |
| `--help` 또는 `-h` | 도움말 출력 |

## 출력 예시
도구 실행 결과는 다음과 같은 방식으로 정리됩니다:
<img src="docs/images/output.png" width="750" alt="output">
