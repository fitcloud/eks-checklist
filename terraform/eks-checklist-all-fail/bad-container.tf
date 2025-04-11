# EBS CSI 드라이버를 사용하는 스토리지 클래스
resource "kubernetes_storage_class" "ebs_sc" {
  # EBS CSI 드라이버가 EKS Addon을 통해서 생성될 경우
  count = lookup(module.eks.cluster_addons, "aws-ebs-csi-driver", null) != null ? 1 : 0

  metadata {
    name = "ebs-sc"
    annotations = {
      "storageclass.kubernetes.io/is-default-class" : "true"
    }
  }
  storage_provisioner = "ebs.csi.aws.com"
  volume_binding_mode = "WaitForFirstConsumer"
  parameters = {
    type      = "gp3"
    encrypted = "true"
  }
}

# 기본값으로 생성된 스토리지 클래스 해제
resource "kubernetes_annotations" "default_storageclass" {
  count = lookup(module.eks.cluster_addons, "aws-ebs-csi-driver", null) != null ? 1 : 0

  api_version = "storage.k8s.io/v1"
  kind        = "StorageClass"
  force       = "true"

  metadata {
    name = "gp2"
  }
  annotations = {
    "storageclass.kubernetes.io/is-default-class" = "false"
  }

  depends_on = [
    kubernetes_storage_class.ebs_sc
  ]
}

resource "kubernetes_persistent_volume_claim" "nginx_pvc" {
  metadata {
    name = "nginx-pvc"
  }

  spec {
    access_modes = ["ReadWriteOnce"]

    resources {
      requests = {
        storage = "5Gi"
      }
    }

    storage_class_name = "ebs-sc"
  }
}

resource "kubernetes_pod" "nginx" {
  metadata {
    name = "nginx"
  }

  spec {
    container {
      name  = "nginx"
      image = "nginx:latest"

      volume_mount {
        name       = "nginx-storage"
        mount_path = "/usr/share/nginx/html"
      }
    }

    volume {
      name = "nginx-storage"

      persistent_volume_claim {
        claim_name = "nginx-pvc"
      }
    }
  }
}
