
# ALB的健康检查使用ksyun_healthcheck定义
resource "ksyun_healthcheck" "default" {
  listener_id = ksyun_alb_listener.test.id
  health_check_state = "start"
  healthy_threshold = 2
  interval = 20
  timeout = 200
  unhealthy_threshold = 2
  url_path = "/monitor"
  is_default_host_name = false
  host_name = "www.ksyun.com"
}