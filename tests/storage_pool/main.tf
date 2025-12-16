terraform {
  required_providers {
    juju = {
      source = "registry.terraform.io/juju/juju"
      version = "1.0.0"
    }
  }
}

provider "juju" {
}

module "model" {
  topic = "storage-pool"
  source = "../../modules/model_random"
}

output "model_name" {
  value = module.model.name
}

resource "juju_model" "this" {
  name = module.model.name
}

resource "juju_storage_pool" "this" {
  name             = "pool"
  model_uuid       = juju_model.this.uuid
  storage_provider = "kubernetes"
}

resource "juju_application" "this" {
  model_uuid = juju_model.this.uuid
  name       = "db"
  charm {
    name    = "postgresql-k8s"
    channel = "14/stable"
    base    = "ubuntu@22.04"
  }

  storage_directives = {
    "pgdata" = "1M,${juju_storage_pool.this.name}"
  }
  constraints = "arch=${var.arch} tags=${var.tags}"
}
