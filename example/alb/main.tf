terraform {
  required_providers {
    ksyun = {
      source = "kingsoftcloud/ksyun"
    }
  }
}

variable region {
  type    = string
  default = "cn-beijing-6"
}
variable az {
  type    = string
  default = "cn-beijing-6a"
}
provider "ksyun" {
  region = var.region
}




