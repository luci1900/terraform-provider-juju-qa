terraform {
  required_providers {
    juju = {
      source = "registry.terraform.io/juju/juju"
      version = "1.1.0-rc1"
    }
  }
}

provider "juju" {
}

resource "juju_model" "this" {
  name = "tfqa-cross-controller-offering"
}

resource "juju_application" "source" {
  model_uuid = juju_model.this.uuid
  name = "dummy-source"

  charm {
    name    = "juju-qa-dummy-source"
  }

  config = {
    token = "abc"
  }

  constraints = "arch=${var.arch} tags=${var.tags}"
}

resource "juju_offer" "source" {
  model_uuid       = juju_model.this.uuid
  depends_on = [juju_application.source]

  application_name = juju_application.source.name
  endpoints       = ["sink"]
}
