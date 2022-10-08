

provider "kubernetes" {
  config_path = "~/.kube/config"
}
resource "kubernetes_namespace" "tracker" {
  metadata {
    name = "tracker"
  }
}
resource "kubernetes_deployment" "tracker" {
  metadata {
    name      = "tracker"
    namespace = kubernetes_namespace.tracker.metadata.0.name
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "tracker"
      }
    }
    template {
      metadata {
        labels = {
          app = "tracker"
        }
      }
      spec {
        container {
          //INSERT YOUR PROJECT ID HERE <-----------------
          image = "gcr.io/first-vigil-296706/tracker"
          name  = "tracker"
          port {
            container_port = 80
          }
        }
      }
    }
  }
}
resource "kubernetes_service" "tracker" {
  metadata {
    name      = "tracker"
    namespace = kubernetes_namespace.tracker.metadata.0.name
  }
  spec {
    selector = {
      app = kubernetes_deployment.tracker.spec.0.template.0.metadata.0.labels.app
    }
    type = "NodePort"
    port {
      node_port   = 30201
      port        = 80
      target_port = 80
    }
  }
}


