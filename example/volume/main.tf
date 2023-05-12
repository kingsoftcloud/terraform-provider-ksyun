# Specify the provider and access details
provider "ksyun" {
  access_key = "ak"
  secret_key = "sk"
  region     = "cn-shanghai-3"
}

resource "ksyun_volume" "default" {
  volume_name       = "test"
  volume_type       = "SSD3.0"
  size              = 15
  charge_type       = "Daily"
  availability_zone = "cn-shanghai-3a"
  volume_desc       = "test"

  ## 传入快照ID，用快照创建EBS盘
  ## 注意：如果使用的整机镜像创建主机，API会自动根据镜像中包含的快照创建数据盘，不需在tf配置中定义数据盘
  # snapshot_id = "snapshot_id"
}

