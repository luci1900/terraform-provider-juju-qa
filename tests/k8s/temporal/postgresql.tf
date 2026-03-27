module "db_model" {
  topic = "temporal-db"
  source = "../../modules/model_random"
}

resource "juju_model" "db" {
  name = module.db_model.name
}

resource "juju_application" "postgresql" {
  name  = "db"
  model_uuid = juju_model.db.uuid
  trust = true
  units = 1

  charm {
    name    = "postgresql-k8s"
    channel = "14/stable"
    base    = "ubuntu@22.04"
  }

  config = {
    plugin_pg_trgm_enable   = true
    plugin_uuid_ossp_enable = true
    plugin_btree_gin_enable = true
  }
  constraints = "arch=${var.arch} tags=${var.tags}"
}

resource "juju_offer" "database" {
  model_uuid       = juju_model.db.uuid
  depends_on = [juju_application.postgresql]

  application_name = juju_application.postgresql.name
  endpoints        = ["database"]
}
