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

# 암호화 안됬으며 읽기 쓰기 둘다 하는 EBS 볼륨을 사용하는 PVC
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

## 볼륨 사용하는 레플리카 1인 nginx
resource "kubernetes_deployment" "nginx" {
  metadata {
    name = "nginx"
    labels = {
      app = "nginx"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "nginx"
      }
    }

    template {
      metadata {
        labels = {
          app = "nginx"
        }
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
  }
}

## 싱글톤 파드 아파치
resource "kubernetes_pod" "apache_singleton" {
  metadata {
    name = "apache-singleton"
    labels = {
      app = "apache"
    }
  }

  spec {
    container {
      name  = "apache"
      image = "httpd:2.4" # Apache Docker 이미지

      port {
        container_port = 80
      }
    }
  }
}

# ## endpoint slice을 사용하지 않는 서비스
# resource "kubernetes_service" "nginx" {
#   metadata {
#     name = "nginx"
#   }

#   spec {
#     selector = {
#       app = kubernetes_pod.nginx.metadata[0].labels["app"]
#     }
#     port {
#       port     = 80
#       protocol = "TCP"
#     }
#     type = "NodePort"

#     # endpoint slice을 사용하지 않는 서비스
#     publish_not_ready_addresses = true
#   }
# }
