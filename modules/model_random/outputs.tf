output "name" {
  value = "${var.ns}-${var.topic}-${random_string.run.result}"
}
