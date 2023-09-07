# ALB
resource "ksyun_alb" "test" {
  alb_name    = "tf-alb-test"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.test.id
  project_id  = 0
  enabled_log = false
#  klog_info {
#    log_pool_name = "log_pool_name"
#    project_name = "project_name"
#  }
}

data "ksyun_certificates" "default" {
}

# ALB监听器
resource "ksyun_alb_listener" "test" {
  alb_id             = ksyun_alb.test.id
  alb_listener_name  = "alb-test-listener"
  protocol           = "HTTPS"
  port               = 8099
  alb_listener_state = "start"
    certificate_id     = data.ksyun_certificates.default.certificates.0.certificate_id
  http_protocol      = "HTTP1.1"
  session {
    cookie_type                = "RewriteCookie"
    cookie_name                = "KLBRSID1"
    session_state              = "start"
    session_persistence_period = 3100
  }
}

# 监听器的转发策略
resource "ksyun_alb_rule_group" "default" {
  alb_listener_id         = ksyun_alb_listener.test.id
  alb_rule_group_name     = "tf_alb_rule_group"
  backend_server_group_id = ksyun_lb_backend_server_group.test.id
  alb_rule_set {
    alb_rule_type  = "domain"
    alb_rule_value = "www.ksyun.com"
  }
  alb_rule_set {
    alb_rule_type  = "url"
    alb_rule_value = "/test/path"
  }
  listener_sync = "on"
}


data "ksyun_certificates" "listener_cert" {
  name_regex = "test"
}

resource "ksyun_alb_listener_cert_group" "default" {
  alb_listener_id = ksyun_alb_listener.test.id
  certificate {
    certificate_id = data.ksyun_certificates.listener_cert.certificates.0.certificate_id
  }
}