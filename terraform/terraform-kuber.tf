provider "kubernetes" {
}

resource "kubernetes_ingress_v1" "minio_ingress" {
  metadata {
    name = "mino-ingress"
    namespace = "argo"
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" = "/$1"
    }
  }

  spec {
    rule {
      http {
        path {
          path = "/"
          backend {
            service {
              name = "minio"
              port {
                number = 9000
              }
            }
          }
        }
      }
    }
  }
}
