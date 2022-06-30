provider "ksyun" {
}

resource "ksyun_vpc" "test" {
  vpc_name   = "tf-example-vpc-01"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name      = "tf-acc-subnet1"
  	cidr_block = "10.0.5.0/24"
      subnet_type = "Normal"
      vpc_id  = "${ksyun_vpc.test.id}"
      gateway_ip = "10.0.5.1"
      dns1 = "198.18.254.41"
      dns2 = "198.18.254.40"
      availability_zone = "cn-shanghai-2a"
}
