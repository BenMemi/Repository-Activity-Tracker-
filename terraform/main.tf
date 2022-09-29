

provider "google" {
credentials = file("/Users/memi/Repository-Activity-Tracker-/first-vigil-296706-f83f2a2f26d0.json") //swap to your service Key!!!

  project = "first-vigil-296706"
  region  = "us-central1"
  zone    = "us-central1-c"
}

resource "google_compute_instance" "vm_instance" {
  name                      = "balancer-github-tracker"
  machine_type              = "e2-micro"
  allow_stopping_for_update = false
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1804-lts"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral public IP
    }
  }
  metadata_startup_script = "${file("start_script.sh")}"
}



