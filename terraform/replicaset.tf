resource "kubernetes_manifest" "example" {
  manifest = {
    apiVersion = "apps/v1"
    kind       = "ReplicaSet"
    metadata = {
      name      = "single-replica"
      namespace = "default"
      labels = {
        app = "test-app"
      }
    }
    spec = {
      replicas = 2  # ✅ FAIL 조건: 1개의 Pod만 생성
      selector = {
        matchLabels = {
          app = "test-app"
        }
      }
      template = {
        metadata = {
          labels = {
            app = "test-app"
          }
        }
        spec = {
          containers = [
            {
              name  = "nginx"
              image = "nginx:latest"
              ports = [
                {
                  containerPort = 80
                }
              ]
            }
          ]
        }
      }
    }
  }
}