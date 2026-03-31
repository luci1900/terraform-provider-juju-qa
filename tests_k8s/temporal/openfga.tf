module "iam_model" {
  topic = "temporal-iam"
  source = "../../modules/model_random"
}

resource "juju_model" "openfga_model" {
  name = module.iam_model.name
}

resource "juju_application" "openfga" {
  name  = "iam"
  model_uuid = juju_model.openfga_model.uuid
  trust = true
  units = 1

  charm {
    name    = "openfga-k8s"
    channel = "latest/edge"
    base    = "ubuntu@22.04"
  }
  constraints = "arch=${var.arch} tags=${var.tags}"
}

resource "juju_offer" "openfga" {
  model_uuid       = juju_model.openfga_model.uuid
  depends_on = [juju_application.openfga]

  application_name = juju_application.openfga.name
  endpoints        = ["openfga"]
}

resource "juju_integration" "db_integration" {
  model_uuid = juju_model.openfga_model.uuid

  application {
    offer_url = juju_offer.database.url
  }

  application {
    name     = juju_application.openfga.name
    endpoint = "database"
  }
}
