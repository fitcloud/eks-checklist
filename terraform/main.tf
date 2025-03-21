# 요구되는 테라폼 제공자 목록
terraform {
  required_version = "1.10.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.86.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.31.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.14.0"
    }
    kubectl = {
      source  = "alekc/kubectl"
      version = ">= 2.0.4"
    }
  }
}