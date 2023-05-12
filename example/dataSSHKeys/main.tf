# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  keys
data "ksyun_ssh_keys" "default" {
#  根据key的名称查找
  key_name="ydx_zh_pub"

#  查找多个名称
#  key_names = ["ydx_zh_pub", "test"]

# 根据key的id（指纹）查询
#  ids = ["b433f65e-29e7-4b77-b40d-xxxxxxxxxxxxx", "7f5e9a85-9f89-4988-9792-xxxxxxxxxxxxx"]

  output_file="output_result"
}

# 返回值示例

#[
#  {
#    "create_time": "2022-08-01 10:36:52",
#    "key_id": "b433f65e-29e7-4b77-b40d-142xxxxxxxxx",
#    "key_name": "test",
#    "public_key": "ssh-rsa xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
#  }
# ]

