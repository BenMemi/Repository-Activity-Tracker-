

provider "google" {
  credentials = file("/Users/memi/Repository-Activity-Tracker-/first-vigil-296706-f83f2a2f26d0.json")

  project = "first-vigil-296706"
  region  = "us-central1"
  zone    = "us-central1-c"
}

resource "google_compute_instance" "vm_instance" {
  name                      = "balancer_github_tracker"
  machine_type              = "e2-micro"
  allow_stopping_for_update = false
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    network = google_compute_network.vpc_tr_network.self_link
    access_config {
    }
  }

  metadata_startup_script = "${file("install_docker.sh")}"
}

resource "google_compute_network" "vpc_tr_network" {
  name                    = "github-tracker-network"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "ssh-rule" {
  name    = "ssh-github-tracker"
  network = google_compute_network.vpc_tr_network.self_link
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  source_ranges = ["0.0.0.0/0"]

}

