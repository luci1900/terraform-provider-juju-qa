module "ingress_model" {
  topic = "temporal-ingress"
  source = "../../modules/model_random"
}

resource "juju_model" "traefik" {
  name = module.ingress_model.name
}

resource "juju_application" "traefik" {
  name  = "ingress"
  model_uuid = juju_model.traefik.uuid
  trust = true
  units = 1

  charm {
    name    =  "traefik-k8s"
    channel = "latest/edge"
    base    = "ubuntu@20.04"
  }

  config = {
    external_hostname = ""
    routing_mode      = "path"
  }
  constraints = "arch=${var.arch} tags=${var.tags}"
}

resource "juju_offer" "ingress" {
  model_uuid       = juju_model.traefik.uuid
  depends_on = [juju_application.traefik]

  application_name = juju_application.traefik.name
  name             = "${juju_application.traefik.name}-ingress"
  endpoints        = ["ingress"]
}
