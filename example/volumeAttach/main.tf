# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region     = "cn-shanghai-3"
}

resource "ksyun_volume_attach" "default" {
  volume_id   = "9778a85d-ea38-4521-a0d1-987538cbdc40"
  instance_id = "37ed557e-c092-418d-97bf-9642315de2f1"
}