# module "eks_blueprints_addons" {
#   source  = "aws-ia/eks-blueprints-addons/aws"
#   version = "1.20.0"

#   cluster_name      = module.eks.cluster_name
#   cluster_endpoint  = module.eks.cluster_endpoint
#   cluster_version   = module.eks.cluster_version
#   oidc_provider_arn = module.eks.oidc_provider_arn

#   enable_aws_load_balancer_controller = false
#   enable_cluster_autoscaler           = true
#   cluster_autoscaler = {
#     name          = "cluster-autoscaler"
#     chart_version = "9.29.0"
#     repository    = "https://kubernetes.github.io/autoscaler"
#     namespace     = "kube-system"
#     # values        = [templatefile("${path.module}/values.yaml", {})]
#   }
#   enable_karpenter                      = false
#   enable_kube_prometheus_stack          = false
#   enable_metrics_server                 = false
#   enable_external_dns                   = false
#   enable_cert_manager                   = false
#   cert_manager_route53_hosted_zone_arns = ["arn:aws:route53:::hostedzone/Z066090219MWO7J1N8E5U"]

#   tags = {
#     Environment = "eks-checklist"
#   }
# }