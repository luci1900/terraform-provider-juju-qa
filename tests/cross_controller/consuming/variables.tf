variable "shared_tags" {
    type        = string
    description = "Shared tags for all app constraints"
    default    = "juju"
}

variable "physical_tags" {
    type        = string
    description = "Tags for app constraints intended for physical machines"
    default     = ""
}
