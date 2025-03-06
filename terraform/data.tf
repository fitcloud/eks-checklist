# EKS 클러스터 인증 토큰
data "aws_eks_cluster_auth" "this" {
  name = module.eks.cluster_name
}

# AWS 지역 정보 불러오기
data "aws_region" "current" {}

# 현재 설정된 AWS 리전에 있는 가용영역 정보 불러오기
data "aws_availability_zones" "azs" {}

# 현재 Terraform을 실행하는 IAM 객체
data "aws_caller_identity" "current" {}