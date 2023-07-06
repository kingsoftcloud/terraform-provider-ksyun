---
subcategory: "Instance(KEC)"
layout: "ksyun"
page_title: "ksyun: ksyun_instance"
sidebar_current: "docs-ksyun-resource-instance"
description: |-
  Provides a KEC instance resource.
---

# ksyun_instance

Provides a KEC instance resource.

**Note**  At present, 'Monthly' instance cannot be deleted and must wait it to be outdated and released automatically.

## Example Usage

```hcl
resource "ksyun_instance" "default" {
  image_id      = "${data.ksyun_images.centos-7_5.images.0.image_id}"
  instance_type = "S4.1A"
  system_disk {
    disk_type = "SSD3.0"
    disk_size = 30
  }
  data_disk_gb         = 0
  subnet_id            = "${ksyun_subnet.default.id}"
  instance_password    = "Xuan663222"
  keep_image_login     = false
  charge_type          = "Daily"
  purchase_time        = 1
  security_group_id    = ["${ksyun_security_group.default.id}"]
  private_ip_address   = ""
  instance_name        = "xuan-tf-combine"
  instance_name_suffix = ""
  sriov_net_support    = false
  project_id           = 0
  data_guard_id        = ""
  key_id               = ["${ksyun_ssh_key.default.id}"]
  force_delete         = true
}
```

## Argument Reference

The following arguments are supported:

* `charge_type` - (Required, ForceNew) charge type of the instance.
* `image_id` - (Required) The ID for the image to use for the instance.
* `security_group_id` - (Required) Security Group to associate with.
* `subnet_id` - (Required) The ID of subnet. the instance will use the subnet in the current region.
* `auto_create_ebs` - (Optional) Whether to create EBS volumes from snapshots in the custom image, default is false.
* `data_disk_gb` - (Optional) The size of the local SSD disk.
* `data_disks` - (Optional) The list of data disks created with instance.
* `data_guard_id` - (Optional) Add instance being created to a disaster tolerance group. It will be quit the disaster tolerance group, if this field change to null.
* `dns1` - (Optional) DNS1 of the primary network interface.
* `dns2` - (Optional) DNS2 of the primary network interface.
* `force_delete` - (Optional, **Deprecated**) this field is Deprecated and no effect for change Indicate whether to delete instance directly or not.
* `force_reinstall_system` - (Optional) Indicate whether to reinstall system.
* `host_name` - (Optional) The hostname of the instance. only effective when image support cloud-init.
* `iam_role_name` - (Optional) name of iam role.
* `instance_name` - (Optional) The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.
* `instance_password` - (Optional) Password to an instance is a string of 8 to 32 characters.
* `instance_status` - (Optional) The state of instance.
* `instance_type` - (Optional) The type of instance to start. <br> - NOTE: it's may trigger this instance to power off, if instance type will be demotion.
* `keep_image_login` - (Optional) Keep the initial settings of the custom image.
* `key_id` - (Optional) The certificate id of the instance.
* `local_volume_snapshot_id` - (Optional, ForceNew) When the local data disk opens, the snapshot id is entered.
* `private_ip_address` - (Optional) Instance private IP address can be specified when you creating new instance.
* `project_id` - (Optional) The project instance belongs to.
* `purchase_time` - (Optional, ForceNew) The duration that you will buy the resource.
* `sriov_net_support` - (Optional, ForceNew) whether support networking enhancement.
* `system_disk` - (Optional) System disk parameters.
* `tags` - (Optional) the tags of the resource.
* `user_data` - (Optional, ForceNew) The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB. only effective when image support cloud-init.

The `data_disks` object supports the following:

* `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on EBS disk.
* `disk_size` - (Optional, ForceNew) Data disk size. value range: [10, 16000].
* `disk_snapshot_id` - (Optional, ForceNew) When the cloud disk opens, the snapshot id is entered.
* `disk_type` - (Optional, ForceNew) Data disk type.

The `system_disk` object supports the following:

* `disk_size` - (Optional) The size of the data disk. value range: [20, 500].
* `disk_type` - (Optional, ForceNew) System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `extension_network_interface` - extension network interface information.
  * `network_interface_id` - ID of the extension network interface.
* `has_modify_keys` - whether the certificate key has modified.
* `has_modify_password` - whether the password has modified.
* `has_modify_system_disk` - whether the system disk has modified.
* `instance_id` - ID of the instance.
* `network_interface_id` - ID of the network interface.


## Import

Instance can be imported using the id, e.g.

```
$ terraform import ksyun_instance.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

