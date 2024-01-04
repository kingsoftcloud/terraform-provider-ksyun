terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
      version = "1.14.0"
    }
  }
}


provider "ksyun" {
  region = "cn-beijing-6"
}

variable "suffix" {
  default = "tfkcrs"
}