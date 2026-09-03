---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_instance_model"
sidebar_current: "docs-ksyun-resource-instance_model"
description: |-
  Provides a KEC instance model resource.
---

# ksyun_instance_model

Provides a KEC instance model resource.

Instance Model is a launch template that allows you to quickly create instances with predefined configurations.

**Note**  Instance model cannot be modified after creation. If you need to change any parameter, you must delete and recreate the model.

## Example Usage

```hcl
# get images list
data "ksyun_images" "centos-8_0" {
  platform = "centos-8.0"
}

data "ksyun_availability_zones" "default" {
}

# vpc settings of creating instance
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_subnet" "default" {
  subnet_name       = "ksyun-subnet-tf"
  cidr_block        = "10.7.0.0/21"
  subnet_type       = "Normal"
  vpc_id            = ksyun_vpc.default.id
  availability_zone = data.ksyun_availability_zones.default.availability_zones[0].availability_zone_name
}

resource "ksyun_security_group" "default" {
  vpc_id              = ksyun_vpc.default.id
  security_group_name = "ksyun-security-group"
}

resource "ksyun_instance_model" "default" {
  model_name           = "web-server-model"
  image_id             = data.ksyun_images.centos-8_0.images[0].image_id
  instance_type        = "I2.8B"
  charge_type          = "Monthly"
  purchase_time        = 12
  security_group_ids   = [ksyun_security_group.default.id]
  subnet_id            = ksyun_subnet.default.id
  key_id               = "key-12345678"
  instance_name        = "db-server"
  instance_name_suffix = "1"
  project_id           = 1001
  allocate_address     = true
  address_bandwidth    = 5
  address_charge_type  = "Peak"

  data_disks {
    type = "SSD3.0"
    size = 500
  }

  tags {
    key   = "env"
    value = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `charge_type` - (Required, ForceNew) Charge type of the instance. Valid values: Monthly, Daily, HourlyInstantSettlement, Spot.
* `image_id` - (Required, ForceNew) The ID for the image to use for the instance.
* `instance_type` - (Required, ForceNew) The type of instance to start.
* `model_name` - (Required, ForceNew) The name of the instance model. Must be globally unique.
* `security_group_id` - (Required, ForceNew) Security Group IDs to associate with. Currently only supports 1 security group.
* `subnet_id` - (Required, ForceNew) The ID of subnet. The instance will use the subnet in the current region.
* `address_bandwidth` - (Optional, ForceNew) The bandwidth of the EIP.
* `address_charge_type` - (Optional, ForceNew) The charge type of the EIP.
* `address_project_id` - (Optional, ForceNew) The project ID of the EIP.
* `address_purchase_time` - (Optional, ForceNew) The purchase time of the EIP. Required when address_charge_type is Monthly.
* `allocate_address` - (Optional, ForceNew) Whether to allocate EIP to the instance. EIP parameters take effect when true.
* `assembled_image_data_disk_type` - (Optional, ForceNew) Data disk type for assembled image.
* `cpu` - (Optional, ForceNew) The CPU count of the instance.
* `data_disk_gb` - (Optional, ForceNew) The size of the local SSD disk. Not effective for general purpose instances. Value range: [0, 16000].
* `data_disks` - (Optional, ForceNew) The list of data disks created with instance.
* `data_guard_id` - (Optional, ForceNew) Add instance being created to a disaster tolerance group.
* `failure_auto_delete` - (Optional, ForceNew) Whether to automatically delete the instance when creation fails.
* `host_name` - (Optional, ForceNew) The hostname of the instance. OS internal computer name.
* `iam_role_name` - (Optional, ForceNew) The name of IAM role.
* `instance_name_suffix` - (Optional, ForceNew) The suffix of the instance name. Range: 0-9999, effective for batch creation.
* `instance_name` - (Optional, ForceNew) The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.
* `is_distribute_ipv6` - (Optional, ForceNew) Whether to distribute IPv6 address.
* `keep_image_login` - (Optional, ForceNew) Keep the initial settings of the custom image. Mutually exclusive with password/key_id.
* `key_id` - (Optional, ForceNew) The certificate id of the instance. Mutually exclusive with password. Not supported for other-linux.
* `line_id` - (Optional, ForceNew) The line id of the EIP.
* `local_volume_snapshot_id` - (Optional, ForceNew) Local volume snapshot ID.
* `mem` - (Optional, ForceNew) The memory size of the instance.
* `network_interface` - (Optional, ForceNew) Network interface configurations.
* `project_id` - (Optional, ForceNew) The project instance belongs to. 0 is the default project.
* `purchase_time` - (Optional, ForceNew) The duration that you will buy the resource. Required when charge_type is Monthly, value range: [1, 36].
* `sriov_net_support` - (Optional, ForceNew) Whether to support networking enhancement. Valid for I1/C1/I2(8C+) with standard images.
* `sync_tag` - (Optional, ForceNew) Whether to sync EBS tags.
* `system_disk` - (Optional, ForceNew) System disk parameters.
* `tags` - (Optional, ForceNew) Tags to bind to the instance.
* `user_data` - (Optional, ForceNew) The user data to be specified into this instance. Must be encrypted in base64 format.

The `data_disks` object supports the following:

* `delete_with_instance` - (Optional, ForceNew) Whether to delete the data disk when the instance is deleted.
* `disk_size` - (Optional, ForceNew) Data disk size. Value range: [10, 65536].
* `disk_snapshot_id` - (Optional, ForceNew) Snapshot ID for creating data disk.
* `disk_type` - (Optional, ForceNew) Data disk type.
* `snapshot_name` - (Optional, ForceNew) Snapshot name for creating data disk.

The `network_interface` object supports the following:

* `private_ip_address` - (Optional, ForceNew) Private IP address for the network interface.
* `security_group_id` - (Optional, ForceNew) Security group IDs for the network interface.
* `subnet_id` - (Optional, ForceNew) Subnet ID for the network interface.

The `system_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) System disk size. Value range: [20, 500].
* `disk_type` - (Optional, ForceNew) System disk type.

The `tags` object supports the following:

* `key` - (Required) Tag key.
* `value` - (Required) Tag value.
* `id` - (Optional, ForceNew) Tag ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The creation time of the instance model.
* `model_id` - The ID of the instance model.


## Import

Instance Model can be imported using the model id, e.g.

```
$ terraform import ksyun_instance_model.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

