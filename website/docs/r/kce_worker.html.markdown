---
subcategory: "KCE"
layout: "ksyun"
page_title: "ksyun: ksyun_kce_worker"
sidebar_current: "docs-ksyun-resource-kce_worker"
description: |-
  Provides a KCE worker resource.
---

# ksyun_kce_worker

Provides a KCE worker resource.

#

## Example Usage

```hcl
resource "ksyun_kce_worker" "default" {
  cluster_id        = ksyun_kce_cluster.test_cluster.id
  instance_id       = ksyun_instance.traffic_analysis.0.id
  image_id          = data.ksyun_kce_instance_images.test.image_set.0.image_id
  instance_password = "1235Test$"
  data_disk {
    auto_format_and_mount = true
    file_system           = "ext4"
    mount_target          = "/data"
  }
  container_runtime = "docker"
  docker_path       = "/data/docker_new"
  user_script       = "abc"
  pre_user_script   = "def"
  schedulable       = false

  container_log_max_size  = 200
  container_log_max_files = 20
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the kce cluster.
* `image_id` - (Required, ForceNew) The ID of the image which support KCE.
* `instance_id` - (Required, ForceNew) The ID of the kec instance. The instance will be shut down while being added to the kce cluster.
* `container_log_max_files` - (Optional, ForceNew) Specify custom data to configure a node, namely, specify the script to run after you deploy the node. You must ensure the reentrancy and retry logic of the script. You can view the script and its log files in the /usr/local/ksyun/kce/userscript directory of the node.
* `container_log_max_size` - (Optional, ForceNew) The maximum size of a container log file. When the size of a container log file reaches this limit, a new container log file is generated for data writing. The default value is 100 MB.
* `container_path` - (Optional, ForceNew) The storage path of the container. If not specified, the default is /data/container. Note: when this parameter is passed, the DockerPath parameter is invalid.
* `container_runtime` - (Optional, ForceNew) container runtime instruction.
* `data_disk` - (Optional, ForceNew) Data Disk config.
* `docker_path` - (Optional, ForceNew) The storage path of the container. If not specified, the default is /data/docker.
* `extra_arg` - (Optional, ForceNew) Custom parameters for k8s components on the node.
* `instance_password` - (Optional, ForceNew) The password of the instance.
* `pre_user_script` - (Optional, ForceNew) The user script in base64 encoding. This script will be executed on the node before the k8s component runs. Users need to ensure the re-entry and retry logic of the script. The script and the generated log file can be found in the /usr/local/ksyun/kce/pre_userscript directory.
* `schedulable` - (Optional) Whether the node can be normally scheduled after being added to the cluster. The default is true.
* `user_script` - (Optional, ForceNew) The user script in base64 encoding. This script will be executed on the node after the k8s component runs. Users need to ensure the re-entry and retry logic of the script. The script and the generated log file can be found in the /usr/local/ksyun/kce/userscript directory.

The `data_disk` object supports the following:

* `auto_format_and_mount` - (Optional, ForceNew) Whether to format and mount the data disk, with the default value of true. If set to false, the FileSystem and MountTarget fields will not take effect.
* `file_system` - (Optional, ForceNew) The file system of the data disk, with optional values of ext3, ext4, and xfs. The default value is ext4. If the disk already has a file system, no processing will be performed. If there is no file system, it will be formatted according to the user's definition, only taking effect on the first disk.
* `mount_target` - (Optional, ForceNew) The mounting point of the data disk, which will be mounted and only take effect on the first disk.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

KCE worker can be imported using the id, e.g.

```
$ terraform import ksyun_kce_worker.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

