
resource "ksyun_vpc" "foo" {
  vpc_name      = "tf-${var.suffix}-vpc"
  cidr_block    = "10.7.0.0/21"
}

data "ksyun_availability_zones" "default" {
}

resource "ksyun_subnet" "foo" {
  subnet_name = "tf-${var.suffix}-subnet"
  cidr_block  = "10.7.0.0/21"
  subnet_type = "Reserve"
  vpc_id      = ksyun_vpc.foo.id
  gateway_ip  = "10.7.0.1"
  dns1        = "198.18.254.41"
  dns2        = "198.18.254.40"
  availability_zone = data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name
}


resource "ksyun_kcrs_vpc_attachment" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	vpc_id = ksyun_vpc.foo.id
	reserve_subnet_id = ksyun_subnet.foo.id
	enable_vpc_domain_dns = true
}

# data "ksyun_certificates" "foo" {
#  ids=[]
#  certificate_source = "Submission"
#  output_file = "certificates"
# }

# Create a CustomizedDomain under the repository instance
# resource "ksyun_kcrs_customized_domain" "foo" {
# 	instance_id = ksyun_kcrs_instance.foo.id
# 	realm_name = "tftest.com"
# 	cert_id = data.ksyun_certificates.foo.certificates[0].certificate_id
# }