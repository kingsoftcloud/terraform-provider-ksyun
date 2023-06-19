# ALB的ACL用ksyun_lb_acl定义
resource "ksyun_lb_acl" "default" {
  load_balancer_acl_name = "tf-acl"
}

resource "ksyun_lb_acl_entry" "default" {
  load_balancer_acl_id = ksyun_lb_acl.default.id
  cidr_block           = "192.168.1.1/32"
  rule_number          = 10
  rule_action          = "allow"
  protocol             = "ip"
}

resource "ksyun_lb_listener_associate_acl" "default" {
  listener_id          = ksyun_alb_listener.test.id
  load_balancer_acl_id = ksyun_lb_acl.default.id
}