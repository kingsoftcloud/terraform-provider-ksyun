variable "volume_size" {
  default = 20
}
resource "ksyun_instance" "default" {
  image_id="IMG-cf8fe1f2-68f2-4483-ae33-ba463c203278"
  instance_type="N3.1A"
//  local_volume_snapshot_id = "1111111111"
  key_id=["0dd40a91-5712-4b8a-a4d7-f9d9412c9d2b","9c45b560-e51d-4aee-9e99-0e292476692d"]
  system_disk{
    disk_type="SSD3.0"
    disk_size= 40
  }
  data_disks {
    disk_type = "SSD3.0"
    disk_size = var.volume_size
    delete_with_instance = true
  }
//  data_disks {
//    disk_type = "SSD3.0"
//    disk_size = 200
//  }
//  data_disk_gb=0
  #only support part type
  subnet_id="05c45fcf-405e-441a-8842-a26809595a3e"
//  instance_password="Aa123456"
  keep_image_login=false
  charge_type="Daily"
//  purchase_time=1
  security_group_id=["35ac2642-1958-4ed7-b02c-dc86f27bc9d9","7e2f45b5-e79d-4612-a7fc-fe74a50b639a"]
//  private_ip_address=""
  instance_name="xym-new-1"
//  sriov_net_support="false"
  project_id=0
//  data_guard_id=""
//  force_delete =true
//  user_data=""
//  iam_role_name = "KsyunKECImageImportDefaultRole"
  force_reinstall_system = false
  instance_status = "active"
//  dns1 = "198.18.254.30"
//  dns2 = "198.18.254.31"
  tags = {
    "xym-test" ="sdsds",
  }
}
