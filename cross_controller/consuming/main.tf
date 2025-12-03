terraform {
  required_providers {
    juju = {
      source = "registry.terraform.io/juju/juju"
      version = "1.1.0-rc1"
    }
  }
}

variable "offering_controller_name" {
  type = string
}

variable "offering_controller_addresses" {
  type = string
}

variable "offering_controller_username" {
  type = string
}

variable "offering_controller_password" {
  type = string
} 

variable "offering_controller_ca_cert" {
  type = string 
}

provider "juju" {
  offering_controllers = {
      "${var.offering_controller_name}" = {
      controller_addresses = var.offering_controller_addresses
      username             = var.offering_controller_username
      password             = var.offering_controller_password
      ca_certificate       = var.offering_controller_ca_cert
    }
  }
}

resource "juju_model" "this" {
  name = "tfqa-cross-controller"
}

resource "juju_application" "sink" {
  name = "dummy-sink"
  model_uuid = juju_model.this.uuid

  charm {
    name    = "juju-qa-dummy-sink"
  }

}

resource "juju_integration" "sink-source" {
  model_uuid = juju_model.this.uuid
  depends_on = [juju_application.sink]

  application {
    offering_controller = "${var.offering_controller_name}"
    offer_url = "admin/tfqa-cross-controller-offering.dummy-source"
    endpoint = "sink"
  }

  application {
    name     = juju_application.sink.name
    endpoint = "source"
  }
}
