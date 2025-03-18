# Specify the provider and access details
provider "ksyun" {
  access_key = "your ak"
  secret_key = "your sk"
  region = "ecp instance region"
}


# update kpfs acl
resource "ksyun_kpfs_acl" "default" {
  epc_id = "c6c683f8-5bb4-4747-8516-9a61f01c4bce"
  kpfs_acl_id = "4a42284d7f354e6a9b2d7a3454e0b495"
}
