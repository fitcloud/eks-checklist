<h1 align="center" style="border-bottom: none">
    <a href="https://fitcloud.github.io/eks-checklist/guidedocs/introduction/" target="_blank"><img alt="EKS-Checklist" src="/docs/images/EKS_Checklist.png" width="400"></a><br>EKS-Checklist
</h1>

<p align="center">자세한 내용은 <a href="https://fitcloud.github.io/eks-checklist/home/introduction/" target="_blank">EKS Checklist</a>을 방문하여 문서 및 가이드북을 확인하세요</p>

<div align="center"
   
[![GitHub Release][release-img]][release]
[![Downloads][download-img]][release]
[![Language][go-img]][go-link]
[![Cloud Platform][aws-img]][aws-link]
[![Best Practice][eks-best-img]][eks-best-link]
[![Best Practice][k8s-best-img]][k8s-best-link]
[![Runbook][runbook-img]][runbook-link]
[![Service][eks-img]][eks-link]

<!-- Badge Definitions -->
[release-img]: https://img.shields.io/github/v/release/fitcloud/eks-checklist?logo=github
[download-img]: https://img.shields.io/github/downloads/fitcloud/eks-checklist/total.svg
[go-img]: https://img.shields.io/badge/Language-Go-blue?logo=go
[aws-img]: https://img.shields.io/badge/Cloud_Platform-AWS-FF9900?logo=amazon-aws&logoColor=white
[eks-img]: https://img.shields.io/badge/Service-Elastic%20Kubernetes%20Service-FF9900?logo=amazon-eks&logoColor=white
[eks-best-img]: https://img.shields.io/badge/AWS_EKS-Best_Practice-FF9900?logo=amazon-eks&logoColor=white
[k8s-best-img]: https://img.shields.io/badge/Kubernetes-Best_Practice-326CE5?logo=kubernetes&logoColor=white
[runbook-img]: https://img.shields.io/badge/Runbook-Guide-2f7b6f?logo=read-the-docs&logoColor=white

<!-- Links -->
[release]: https://github.com/fitcloud/eks-checklist/releases
[go-link]: https://go.dev/
[aws-link]: https://aws.amazon.com/
[eks-link]: https://aws.amazon.com/eks/
[eks-best-link]: https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices
[k8s-best-link]: https://kubernetes.io/docs/setup/best-practices/
[runbook-link]: https://fitcloud.github.io/eks-checklist/guidedocs/introduction/

</div>

# EKS Checklist
### EKS Checklist는 Amazon EKS (Elastic Kubernetes Service) 클러스터의 정보를 확인하여 EKS 클러스터가 AWS EKS 권장 모범 사례 및 Kubernetes의 권장 모범 사례를 준수하고 있는지 점검하고, 클러스터 운영을 최적화하는 데 필요한 권장 사항을 제공합니다.

* **쉬운 셀프 점검** EKS Cluster가 AWS EKS 모법사례 및 Kubernetes 모범사례에 준수 하고 잇는지 쉽게 확인이 가능합니다.
* **보고서 제공** 사용자는 원하는 보고서 양식 html, txt, pdf 형식으로 선택하여 제공받을 수 있으며 한눈에 문제점 식별 할 수 있습니다.
* **가이드 제공** EKS Cluster의 문제점에 대한 각각의 항목에 대하여 자세한 가이드를 제공받아 문제점을 개선 할 수 있습니다.
* **다양한 OS 환경 지원** Window, Linux, MacOS 환경에 맞는 툴을 제공받아 사용 할 수 있습니다.

예) CMD Terminal(Window)

<img src="docs/images/output.png" width="750" alt="output">

## 프로젝트 목적 / 배경
Amazon EKS는 관리형 Kubernetes 서비스이지만, 클러스터를 직접 운영하다 보면 다양한 설정 실수나 비효율이 발생하기 쉽습니다. 비용,권한 설정,네트워크 구성, 확장성 부족 등의 문제는 서비스 운영에 큰 영향을 주게 됩니다.
**EKS-Checklist**는 이러한 문제를 점검하고 개선을 할 수 있도록 개발된 도구입니다. AWS 및 Kubernetes의 권장 모범 사례(Best Practices)를 기반으로 클러스터를 자동으로 분석하여 운영자가 **효율적인 클러스터 환경을 유지**할 수 있도록 개발되었습니다.

## 점검 항목
- **비용 최적화 (Cost)**: EKS 클러스터의 리소스 사용을 점검하여 불필요한 비용을 줄일 수 있는 방법을 제공합니다.
- **일반 설정 (General)**: 클러스터의 기본 설정과 환경이 적절하게 구성되었는지 확인합니다.
- **네트워크 (Network)**: VPC, 서브넷, 보안 그룹 등 네트워크 구성을 점검하여 네트워크가 올바르게 설정되었는지 확인합니다.
- **확장성 (Scalability)**: 클러스터가 필요에 따라 확장 가능하도록 설정되었는지 점검합니다.
- **보안 (Security)**: IAM 정책, 인증 및 권한 설정 등 보안 설정을 점검하여 클러스터가 안전하게 운영되고 있는지 확인합니다.
- **안정성 (Reliability)**: 클러스터의 안정성을 위한 백업, 모니터링 및 로깅 설정이 제대로 되어 있는지 점검합니다.

## 필수 조건
EKS-Checklist를 사용하기 전에 다음 항목들이 설치되고 올바르게 설정되어 있어야 합니다.
1. **AWS CLI**: AWS CLI가 설치되어 있고 적절한 권한으로 설정되어 있어야 합니다. 설치 및 설정 방법은 [여기](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)를 참고하세요.
2. **AWS EKS 클러스터 접근 권한**: EKS-Checklist는 AWS EKS 클러스터에 접근할 수 있는 권한을 필요로 합니다.

## 설치 방법
EKS-Checklist는 GitHub의 릴리즈 페이지에서 원하는 버전을 다운로드하여 사용할 수 있습니다.
### GitHub Releases에서 다운로드
1. [GitHub Releases](https://github.com/fitcloud/eks-checklist/releases) 페이지로 이동합니다.
2. 운영 체제에 맞는 바이너리를 다운로드합니다.

## 사용법
### 실행 방법
```bash
{바이너리 파일명} {--flags}
예: eks-checklist --help
```
- `--context` : 사용할 kubeconfig 컨텍스트 이름
- `--kubeconfig` : 사용할 kubeconfig 파일 경로 (기본값: `"C:\\Users\\사용자이름\\.kube\\config"`)
- `--filter` : 출력 결과 필터링 옵션 (`all`, `pass`, `fail`, `manual`)
- `--output` : 출력 형식 지정 (`text`, `html`) — 기본값: `text`
- `--profile` : 사용할 AWS CLI 프로파일 이름
- `--sort` : 결과를 상태별(PASS / FAIL / MANUAL)로 정렬
- `-h`, `--help` : 도움말 출력
### macOS (Darwin)
1. [Releases 페이지](https://github.com/fitcloud/eks-checklist/releases)에서 macOS용 바이너리를 다운로드
   예: `eks-checklist-darwin-amd64`
2. 실행 권한 부여 및 설치:
```bash
chmod +x eks-checklist-darwin-amd64
sudo mv eks-checklist-darwin-amd64 /usr/local/bin/eks-checklist
```
3. 실행 예시:
```bash
eks-checklist
```
### Linux
1. [Releases 페이지](https://github.com/fitcloud/eks-checklist/releases)에서 Linux용 바이너리 다운로드
   예: `eks-checklist-linux-amd64`
2. 실행 권한 부여 및 설치:
```bash
chmod +x eks-checklist-linux-amd64
sudo mv eks-checklist-linux-amd64 /usr/local/bin/eks-checklist
```
3. 실행 예시:
```bash
eks-checklist
```
### Windows
1. [Releases 페이지](https://github.com/fitcloud/eks-checklist/releases)에서 `.exe` 파일 다운로드
   예: `eks-checklist-windows-amd64.exe`
2. 적절한 폴더에 저장 (예: `C:\\Program Files\\EKS-Checklist\\`)
3. 명령 프롬프트 또는 PowerShell에서 실행:
```bash
eks-checklist.exe
```
### Docker
**다음 세 디렉토리를 마운트하여 사용합니다**:
- `~/.kube`  : kubeconfig file 디렉터리  
- `~/.aws`   : AWS credentials 디렉터리  
- `./output` : 옵션에 따른 레포트 결과물이 저장되는 디텍토리  

1. EKS-Checklist Public Image로 docker 실행:
```bash
docker run -v ~/.kube:/root/.kube -v ~/.aws:/root/.aws -v ./output:/output public.ecr.aws/x5b3c7k0/eks-checklist:latest
```
실행 시 결과물이 바로 나타나며 종료됩니다 --detach 옵션 사용 시 docker logs 명령어를 사용해 결과물을 확인 할 수 있습니다.

### EKS Pod
**조건**: EC2, VPC, EKS 등 AWS 리소스를 조회하기 때문에 IAM Role 권한이 필요합니다.
여기서는 **eksctl을 활용해 IRSA(ServiceAccount)로 Pod에 IAM Role을 주입**하여 사용합니다.

**예) eksctl**: 기본 namespace에 "eks-checklist-sa" 이름의 ServiceAccount를 만들고 두 개의 권한을 부여합니다
 - AmazonEKSWorkerNodePolicy: EKS Cluster 조회 권한
 - AmazonEC2ReadOnlyAccess: EC2, VPC, Subnet 등 조회 권한

1. IAM ServiceAccount 생성 (IRSA): 
```bash
eksctl create iamserviceaccount \
   --name eks-checklist-sa \
   --namespace default \
   --cluster <CLUSTER_NAME> \
   --attach-policy-arn arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy \
   --attach-policy-arn arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess \
   --approve \
   --override-existing-serviceaccounts
```
2. HTML 보고서 추출 Job 배포:
```bash
kubectl apply -f https://raw.githubusercontent.com/fitcloud/eks-checklist/refs/heads/main/manifest/output-html-job.yaml
```
3. 보고서 가져오기:
```bash
kubectl cp <eks-checklist-job-POD-NAME>:/output/<REPORT_NAME> ./eks-checklist-report.html
```
4. 정리 (리소스 삭제):
```bash
eksctl delete iamserviceaccount --cluster <CLUSTER_NAME> --name eks-checklist-sa
kubectl delete -f output-html-job.yaml
```
