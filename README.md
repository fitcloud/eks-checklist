# EKS Checklist

## 개요
EKS Checklist는 Amazon Elastic Kubernetes Service(EKS) 클러스터의 상태를 점검하고, 필수적인 검증을 수행하는 도구입니다.  
이 프로젝트는 Go 언어로 개발되었으며, 샘플 EKS 클러스터는 Terraform을 이용해 생성됩니다.

## 개발 환경 구성
### 1. 환경 변수 설정 (Windows)
Windows 환경에서 `Powershell(Admin)`을 실행하고 다음 명령어를 입력하여 환경 변수를 설정합니다.
```powershell
[System.Environment]::SetEnvironmentVariable('HOME', $env:USERPROFILE,[System.EnvironmentVariableTarget]::Machine)
```

### 2. 개발 도구
- IDE : VSCode를 사용하며, Dev Container Extension을 활용하여 개발 환경을 자동 구성합니다.

- VSCode를 실행하면 루트 디렉토리의 .devcontainer/devcontainer.json 파일에 정의된 환경이 자동으로 세팅됩니다.

### 3. 실행 방법
```sh
go run main.go
```

## 테스트 환경 구성
EKS 클러스터를 이용한 테스트가 필요하며, Terraform을 사용해 환경을 구성합니다.
### 1. Terraform을 이용한 EKS 클러스터 생성
Terraform을 이용해 테스트 환경을 구성하려면 terraform 디렉터리로 이동 후 다음 명령어를 실행합니다.

#### (1).  테라폼 초기화 (필수 모듈 다운로드)
```sh
terraform init
```

#### (2).  인프라 생성 (EKS 클러스터 구축)
```sh
terraform apply
```

### 2. Kubeconfig 설정
EKS 클러스터가 생성되면 kubectl이 클러스터에 접근할 수 있도록 kubeconfig를 설정합니다.
```sh
aws eks update-kubeconfig --name eks-checklist
```

## Git Flow
본 프로젝트는 Git Flow 전략을 따르며, dev 브랜치를 기준으로 기능별 feature 브랜치를 생성하여 작업합니다.

### 1. 브랜치 네이밍 규칙
```
feature/<대분류>-<기능>
```

예시:
```
feature/network-targetip
```

### 2. 개발 프로세스
#### (1). ```dev``` 브랜치에서 ```feature``` 브랜치를 생성하여 기능 개발을 진행합니다.

#### (2). 기능 구현 후 ```dev``` 브랜치로 Pull Request(PR)를 생성합니다.

#### (3). 코드 리뷰를 거친 후 dev 브랜치에 머지(Merge)합니다.

#### (4). 다른 기능 브랜치는 최신 ```dev``` 브랜치를 반영하여 지속적으로 동기화합니다.
