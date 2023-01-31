provider "kubernetes" {
}

resource "kubernetes_ingress" "example_ingress" {
  metadata {
    name = "redis-ingress"
    namespace = "argo"
  }

  spec {
    backend {
      service_name = "minio"
      service_port = 9000
    }

    rule {
      http {
        path {
          backend {
            service_name = "minio"
            service_port = 9000
          }

          path = "/"
        }
      }
    }

    tls {
      secret_name = "ingress-tls"
    }
  }
}
