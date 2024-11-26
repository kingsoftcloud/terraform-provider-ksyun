---
subcategory: "KCE"
layout: "ksyun"
page_title: "ksyun: ksyun_kce_cluster_attachment"
sidebar_current: "docs-ksyun-resource-kce_cluster_attachment"
description: |-
  Provides a KCE attachment resource that attach a new instance to a cluster.
---

# ksyun_kce_cluster_attachment

Provides a KCE attachment resource that attach a new instance to a cluster.

#

## Example Usage

```hcl
data "ksyun_kce_instance_images" "test" {
}

resource "ksyun_kce_cluster_attachment" "foo" {
  cluster_id = "67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx"

  worker_config {
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.2A"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "subnet-xxxxxx"
    security_group_id = ["sg-xxxxxx"]
    charge_type       = "Daily"

  }
  advanced_setting {
    container_runtime = "containerd"
    pre_user_script   = "def"
    label {
      key   = "tf_assembly_kce"
      value = "on_configuration_files"
    }
    taints {
      key    = "key2"
      value  = "value3"
      effect = "NoSchedule"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the kce cluster.
* `worker_config` - (Required, ForceNew) The instance node configuration for attach on cluster.
* `advanced_setting` - (Optional, ForceNew) The advanced settings of the worker node.
* `instance_delete_mode` - (Optional) The instance delete mode when the instance is removed from the cluster. The value can be 'Terminate' or 'Remove'.

The `advanced_setting` object supports the following:

* `container_log_max_files` - (Optional, ForceNew) Customize the number of log files. The default value is 10.
* `container_log_max_size` - (Optional, ForceNew) Customize the maximum size of the log file. The default value is 100m.
* `container_path` - (Optional, ForceNew) The storage path of the container. The default value is /data/container. **Notes:** If this path is specified, the docker_path field will be ignored.
* `container_runtime` - (Optional, ForceNew) Container Runtime.
* `data_disk` - (Optional, ForceNew) The mount setting of data disk. **Notes:** Only impact on the first data disk.
* `docker_path` - (Optional, ForceNew) The storage path of the container. The default value is /data/docker.
* `extra_arg` - (Optional, ForceNew) The extra arguments for the kubelet. The format is key=value. For example, --kubelet-extra-args="key1=value1,key2=value2".
* `label` - (Optional, ForceNew) 
* `pre_user_script` - (Optional, ForceNew) A user script encoded in base64, which will be executed on the node **before** the Kubernetes components run. Users need to ensure the script's re-entrant and retry logic. The script and its generated logs can be found in the directory /usr/local/ksyun/kce/pre_userscript.
* `taints` - (Optional, ForceNew) Taints.
* `user_script` - (Optional, ForceNew) A user script encoded in base64, which will be executed on the node **after** the Kubernetes components run. Users need to ensure the script's re-entrant and retry logic. The script and its generated logs can be found in the directory /usr/local/ksyun/kce/pre_userscript.

The `data_disk` object supports the following:

* `auto_format_and_mount` - (Optional, ForceNew) Whether to format and mount the data disk, default value: true. If this field is filled with false, then the file_system and mount_target fields will not take effect.
* `file_system` - (Optional, ForceNew) The file system of the data disk. The default value is ext4.Valid values: ext3, ext4, xfs.
* `mount_target` - (Optional, ForceNew) The mount target of the data disk.

The `data_disks` object supports the following:

* `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on EBS disk.
* `disk_size` - (Optional, ForceNew) Data disk size. value range: [10, 16000].
* `disk_snapshot_id` - (Optional, ForceNew) When the cloud disk opens, the snapshot id is entered.
* `disk_type` - (Optional, ForceNew) Data disk type.

The `extension_network_interface` object supports the following:


The `system_disk` object supports the following:

* `disk_size` - (Optional) The size of the data disk. value range: [20, 500].
* `disk_type` - (Optional, ForceNew) System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk, `ESSD_SYSTEM_PL0`, The x7 machine type ESSD disk, `ESSD_SYSTEM_PL1`, The x7 machine type ESSD disk, `ESSD_SYSTEM_PL2`, The x7 machine type ESSD disk.

The `taints` object supports the following:

* `effect` - (Required, ForceNew) The effect of the taint. Valid values: NoSchedule, PreferNoSchedule, NoExecute.
* `key` - (Required, ForceNew) The key of the taint.
* `value` - (Required, ForceNew) The value of the taint.

The `worker_config` object supports the following:

* `charge_type` - (Required, ForceNew) charge type of the instance.
* `image_id` - (Required, ForceNew) The ID for the image to use for the instance.
* `instance_type` - (Required, ForceNew) The type of instance to start. <br> - NOTE: it's may trigger this instance to power off, if instance type will be demotion.
* `security_group_id` - (Required, ForceNew) Security Group to associate with.
* `subnet_id` - (Required, ForceNew) The ID of subnet. the instance will use the subnet in the current region.
* `auto_create_ebs` - (Optional) Whether to create EBS volumes from snapshots in the custom image, default is false.
* `data_disk_gb` - (Optional) The size of the local SSD disk.
* `data_disks` - (Optional) The list of data disks created with instance.
* `data_guard_id` - (Optional) Add instance being created to a disaster tolerance group. It will be quit the disaster tolerance group, if this field change to null.
* `dns1` - (Optional) DNS1 of the primary network interface.
* `dns2` - (Optional) DNS2 of the primary network interface.
* `force_delete` - (Optional) Indicate whether to delete instance directly or not.
* `force_reinstall_system` - (Optional) Indicate whether to reinstall system.
* `host_name` - (Optional) The hostname of the instance. only effective when image support cloud-init.
* `iam_role_name` - (Optional) name of iam role.
* `instance_name` - (Optional, ForceNew) The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.
* `instance_password` - (Optional, ForceNew) Password to an instance is a string of 8 to 32 characters.
* `instance_status` - (Optional) The state of instance.
* `keep_image_login` - (Optional) Keep the initial settings of the custom image.
* `key_id` - (Optional) The certificate id of the instance.
* `local_volume_snapshot_id` - (Optional, ForceNew) When the local data disk opens, the snapshot id is entered.
* `private_ip_address` - (Optional) Instance private IP address can be specified when you creating new instance.
* `project_id` - (Optional) The project instance belongs to.
* `purchase_time` - (Optional, ForceNew) The duration that you will buy the resource.
* `role` - (Optional) The role of instance. Valid values: Worker.
* `sriov_net_support` - (Optional, ForceNew) whether support networking enhancement.
* `sync_tag` - (Optional) Indicate whether to sync tags to instance.
* `system_disk` - (Optional) System disk parameters.
* `tags` - (Optional) the tags of the resource.
* `user_data` - (Optional) The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB. only effective when image support cloud-init.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - The ID of the kec instance. The instance will be shut down while being added to the kce cluster.


## Import

KCE worker can be imported using the id, e.g.

```
$ terraform import ksyun_kce_cluster_attachment.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

