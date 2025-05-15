# GEN-001 코드형 인프라

## Meaning
EKS 클러스터 및 관련 인프라 자원(VPC, IAM, 노드 그룹 등)을 수동으로 AWS Console에서 구성하는 것은 반복성과 가시성이 부족합니다. 코드형 인프라(IaC, Infrastructure as Code)를 사용하면 선언적이고 버전 관리 가능한 방식으로 클러스터와 자원을 정의할 수 있어 팀 협업 및 운영 안정성 측면에서 매우 유리합니다.

Terraform, AWS CDK 등의 도구는 이러한 IaC 구현에 널리 사용됩니다

## Impact
- 재현 불가 문제: 수동 설정 시 동일한 인프라 환경을 재현하기 어려움
- 협업 및 변경 이력 부족: 변경 사항이 기록되지 않아 추적 및 협업 어려움
- 자동화 파이프라인 부재: CI/CD 또는 GitOps 파이프라인에 통합 불가

## Diagnosis
EKS 클러스터 및 리소스가 IaC 도구를 통해 배포되었는지 다음 항목을 확인합니다

Git 리포지토리에 Terraform / CDK 등의 IaC 코드 존재 여부

AWS 리소스 태그에 IaC 관련 태그(managed-by: terraform 등) 여부

**example**

```bash
# Git 레포에서 Terraform 코드 유무 확인
find . -name "*.tf"
```

## Mitigation
클러스터 및 관련 리소스의 수명 주기를 IaC 도구를 통해 관리하도록 변경합니다.

**example**

EKS 클러스터: Terraform aws_eks_cluster 리소스로 정의

**Before**

수동으로 AWS 콘솔에서 클러스터 생성 및 리소스 추가

**After**

Terraform 또는 CDK를 통해 클러스터 구성 및 배포

Git 리포지토리에 모든 인프라 코드 관리

[Terraform AWS EKS](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_cluster)
[AWS CDK for EKS](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_eks-readme.html)