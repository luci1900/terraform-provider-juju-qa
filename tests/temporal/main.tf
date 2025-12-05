terraform {
  required_providers {
    juju = {
      source = "registry.terraform.io/juju/juju"
      version = "1.0.0"
    }
  }
}

provider "juju" {}

module "main_model" {
  topic = "temporal"
  source = "../../modules/model_random"
}

output "model_name" {
  value = module.main_model.name
}

resource "juju_model" "this" {
  name = module.main_model.name
}

# application and integration top-level
resource "juju_application" "temporal_k8s" {
  name = "temporal"
  model_uuid = resource.juju_model.this.uuid

  charm {
    name = "temporal-k8s"
  }

  config = {
    num-history-shards = 2
  }
}


resource "juju_integration" "temporal_db" {
  model_uuid = resource.juju_model.this.uuid
 
  application {
    offer_url = juju_offer.database.url
  }

  application {
    name     = juju_application.temporal_k8s.name
    endpoint = "db"
  }
}

resource "juju_integration" "temporal_visibility_db" {
    model_uuid = resource.juju_model.this.uuid

  application {
    offer_url = juju_offer.database.url
  }

  application {
    name     = juju_application.temporal_k8s.name
    endpoint = "visibility"
  }
}

resource "juju_application" "temporal_admin_k8s" {
  name = "temporal-admin"
  model_uuid = juju_model.this.uuid

  charm {
    name = "temporal-admin-k8s"
  }
}

resource "juju_integration" "temporal_admin" {

  application {
    name     = juju_application.temporal_k8s.name
    endpoint = "admin"
  }

  application {
    name     = juju_application.temporal_admin_k8s.name
    endpoint = "admin"
  }
  model_uuid = juju_model.this.uuid
}


resource "juju_application" "temporal_k8s_ui" {
  name = "temporalui"
  model_uuid = resource.juju_model.this.uuid

  charm {
    name = "temporal-ui-k8s"
  }
}

resource "juju_integration" "temporal_ui" {
  model_uuid = juju_model.this.uuid

  application {
    name     = juju_application.temporal_k8s.name
    endpoint = "ui"
  }

  application {
    name     = juju_application.temporal_k8s_ui.name
    endpoint = "ui"
  }
}
