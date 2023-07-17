resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}

resource "ksyun_security_group" "test" {
  vpc_id              = ksyun_vpc.test.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "test" {
  security_group_id = ksyun_security_group.test.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}