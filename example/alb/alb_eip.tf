data "ksyun_lines" "default" {
}

resource "ksyun_eip" "test" {
  line_id       = data.ksyun_lines.default.lines.0.line_id
  band_width    = 1
  charge_type   = "PostPaidByDay"
  purchase_time = 1
  project_id    = 0
}
resource "ksyun_eip_associate" "eip_bind" {
  allocation_id = ksyun_eip.test.id
  instance_id   = ksyun_alb.test.id
  instance_type = "Slb"
}