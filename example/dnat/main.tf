provider "ksyun" {
}


variable "suffix" {
  default = "-dnat-example"
}

data "ksyun_images" "centos-7_5" {
  platform= "centos-7.5"
}


resource "ksyun_vpc" "foo" {
  vpc_name        = "tf-${var.suffix}-vpc"
  cidr_block = "10.7.0.0/21"
}

data "ksyun_availability_zones" "default" {
}

resource "ksyun_subnet" "foo" {
  subnet_name      = "tf-${var.suffix}-subnet"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Normal"
  vpc_id  = ksyun_vpc.foo.id
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = data.ksyun_availability_zones.default.availability_zones[0].availability_zone_name
}



resource "ksyun_security_group" "default" {
  vpc_id = ksyun_vpc.foo.id
  security_group_name="tf-${var.suffix}-security-group"
}


resource "ksyun_nat" "foo" {
  nat_name = "tf-${var.suffix}-nat"
  nat_mode = "Subnet"
  nat_type = "public"
  band_width = 200
  charge_type = "DailyPaidByTransfer"
  vpc_id = ksyun_vpc.foo.id
  nat_ip_number= 3
}


resource "ksyun_instance" "foo" {
  image_id=data.ksyun_images.centos-7_5.images[0].image_id
  instance_type="N3.2B"

  subnet_id=ksyun_subnet.foo.id
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=[ksyun_security_group.default.id]
  instance_name="tf-${var.suffix}-kec"
  sriov_net_support="false"
}



resource "ksyun_nat_instance_bandwidth_limit" "foo" {
  nat_id = ksyun_nat.foo.id
  network_interface_id = ksyun_instance.foo.network_interface_id
  bandwidth_limit = 18
}

resource "ksyun_dnat" "foo" {
  nat_id = ksyun_nat.foo.id
  dnat_name =  "tf-${var.suffix}-dnat-acceptance"
  ip_protocol = "Any"
  nat_ip = ksyun_nat.foo.nat_ip_set[0].nat_ip
  private_ip_address = ksyun_instance.foo.private_ip_address
  description = "tf-${var.suffix}-dnat-1"
}

