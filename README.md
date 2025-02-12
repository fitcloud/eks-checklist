# EKS Checklist

## 개발 환경 구성
1. VSCode를 IDE로 사용한다고 가정하여 Dev Container에 애플리케이션 개발/빌드/테스트에 필요한 모든 도구를 구성하여 사용
2. VSCode를 실행하면 루트 경로의 `.devcontainer/devcontainer.json`에 명시된 구성대로 Dev Container가 자동으로 생성됨

## 테스트 환경 구성
1. 애플리케이션 기능 테스트를 위해서는 EKS 클러스터가 필요하고 EKS 클러스터는 Terraform을 통해서 구성 가능
2. `terraform` 경로에 있는 TF 파일들로 테스트 환경 구성
    - 테라폼 코드 실행에 필요한 제공자 및 모듈 다운로드
        ```
        terraform init
        ```
    - 테라롬 코드로 인프라 구성
        ```
        terraform apply
        ```
3. EKS 클러스터 생성이 완료되면 아래의 명령어를 실행해서 kubeconfig 파일 생성
    ```
    aws eks update-kubeconfig --name eks-checklist
    ```