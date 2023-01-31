provider "kubernetes" {
}

resource "kubernetes_ingress" "example_ingress" {
  metadata {
    name = "redis-ingress"
    namespace = "greenops"
  }

  spec {
    backend {
      service_name = "redisserver"
      service_port = 6379
    }

    rule {
      http {
        path {
          backend {
            service_name = "redisserver"
            service_port = 6379
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
