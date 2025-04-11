# EKS 클러스터
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.33.1"

  cluster_name    = local.project
  cluster_version = var.eks_cluster_version

  # EKS 클러스터 API 엔드포인트 접근 제어
  cluster_endpoint_public_access = var.eks_cluster_endpoint_public_access

  # 클러스터를 생성한 IAM 객체에서 쿠버네티스 어드민 권한 할당
  enable_cluster_creator_admin_permissions = true

  # 보안 그룹을 생성할 VPC
  vpc_id = module.vpc.vpc_id

  # 컨트롤 플레인으로 연결된 ENI를 생성할 서브넷
  control_plane_subnet_ids = module.vpc.private_subnets

  # 로깅 비활성화
  cluster_enabled_log_types   = []
  create_cloudwatch_log_group = false

  # Secret 암호화 비활성화
  cluster_encryption_config = {}

  # 노드그룹을 사용할 경우에만 보안그룹 생성
  create_node_security_group = false

  # EKS Auto Mode 미사용
  enable_auto_mode_custom_tags = false

  # EKS Addons 추가
  cluster_addons = {
    coredns                = {}
    eks-pod-identity-agent = {}
    kube-proxy             = {}
    vpc-cni = {
      // before_compute = true
      configuration_values = jsonencode({
        env = {
          #  Prefix 모드 사용 
          ENABLE_PREFIX_DELEGATION = "true"
        }
      })
    }
  }

  # EKS 노드 그룹 t3.medium 인스턴스 타입 단 1개 생성
  eks_managed_node_groups = {
    nodegroup-1 = {
      instance_types   = ["t3.medium"]
      desired_capacity = 1
      min_size         = 1
      max_size         = 3
      volume_size      = 20
      subnet_ids       = module.vpc.private_subnets

      tags = {
        "Name" = "${local.project}-nodegroup-1"
      }
    }
  }
}