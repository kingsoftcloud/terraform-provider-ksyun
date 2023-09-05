data "ksyun_images" "default" {
  output_file  = "output_result"
  name_regex   = "centos-7.0"
  is_public    = true
  image_source = "system"
}

data "ksyun_ssh_keys" "test" {

}

resource "ksyun_instance" "test" {
  count             = 2
  security_group_id = [
    ksyun_security_group.test.id
  ]
  subnet_id = ksyun_subnet.test.id
  key_id    = [data.ksyun_ssh_keys.test.keys.0.key_id]

  instance_type = "N3.1A"
  charge_type   = "Daily"
  instance_name = "tf-alb-test-vm"
  project_id    = 0

  image_id = data.ksyun_images.default.images.0.image_id
}
