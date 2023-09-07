# 定义ALB后端服务器组


resource "ksyun_lb_backend_server_group" "test" {
  backend_server_group_name = "tf_bsg"
  vpc_id                    = ksyun_vpc.test.id
  backend_server_group_type = "Server"
}

locals {
  private_ip_address = ksyun_instance.test.*.private_ip_address
}
resource "ksyun_lb_register_backend_server" "default" {
  backend_server_group_id=ksyun_lb_backend_server_group.test.id
  backend_server_ip= ksyun_instance.test.0.private_ip_address
  backend_server_port=8090
  weight=10
}
resource "ksyun_lb_register_backend_server" "default2" {
  backend_server_group_id=ksyun_lb_backend_server_group.test.id
  backend_server_ip= ksyun_instance.test.1.private_ip_address
  backend_server_port=8090
  weight=10
}