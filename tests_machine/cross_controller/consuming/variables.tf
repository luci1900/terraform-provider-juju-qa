variable "tags" {
    type        = string
    description = "Tags for all app constraints, including for physical machines"
    default    = ""
}

variable "arch" {
    type        = string
    description = "CPU architecture for app constraints"
    default     = "arm64"
}
