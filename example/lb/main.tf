# Specify the provider and access details
provider "ksyun" {
}
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-example-vpc-02"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-acc-subnet1"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Reserve"
  availability_zone = "cn-beijing-6a"
  vpc_id            = "${ksyun_vpc.test.id}"
}

# Create Load Balancer
resource "ksyun_lb" "default" {
  vpc_id              = "${ksyun_vpc.test.id}"
  load_balancer_name  = "tf-xun1"
  type                = "internal"
  subnet_id           = "${ksyun_subnet.test.id}"
  load_balancer_state = "start"
  private_ip_address  = "10.0.1.2"


  ## 设置LB日志（日志功能文档说明：https://docs.ksyun.com/documents/36890）
  ## 只有公网LB支持日志功能，私网LB开启此功能会报错
  ## 部分不支持此功能的region也会报错

  ## 日志开关，true为开启日志功能，默认为false
  access_logs_enabled = false

  ## 存放日志的ks3 bucket地址
  access_logs_s3_bucket = "xxx.ks3-cn-beijing.ksyuncs.com"
}
