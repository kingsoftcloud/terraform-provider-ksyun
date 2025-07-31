---
subcategory: "KCE"
layout: "ksyun"
page_title: "ksyun: ksyun_kce_cluster"
sidebar_current: "docs-ksyun-resource-kce_cluster"
description: |-
  Provides a KCE cluster resource.
---

# ksyun_kce_cluster

Provides a KCE cluster resource.

~> **NOTE:** We recommend that uses the `ksyun_kce_cluster` resource to create a cluster with few `worker_config`.
If you want to manage more worker instances in this cluster, to use the `ksyun_kce_cluster_attach_existence` or `ksyun_kce_cluster_attachment` resource to attach the worker instances to the cluster. The reason is that the `worker_config` is unchangeable and may cause the cluster to be re-created because it is marked *ForceNew*.

#

## Example Usage

## basic dependency resources

```hcl
data "ksyun_kce_instance_images" "test" {
  output_file = "output_result"
}

data "ksyun_kce_instance_images" "test" {
}

variable "az" {
  default = "cn-beijing-6e"
}

variable "suffix" {
  default = "-kce-complete"
}
```

## create a ManagementCluster

```hcl
resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf-modification${var.suffix}"
  cluster_desc        = "description...modification"
  cluster_manage_mode = "ManagedCluster"
  vpc_id              = ksyun_vpc.test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.252.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.23.17"
  reserve_subnet_id   = ksyun_subnet.reserve.id

  managed_cluster_multi_master {
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = ksyun_security_group.test.id
  }

  worker_config {
    count         = 2
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_worker"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }
}
```

## create a DedicatedCluster

```hcl
resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf-modification${var.suffix}"
  cluster_desc        = "description...modification"
  cluster_manage_mode = "DedicateCluster"
  vpc_id              = ksyun_vpc.test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.252.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.23.17"
  reserve_subnet_id   = ksyun_subnet.reserve.id

  managed_cluster_multi_master {
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = ksyun_security_group.test.id
  }

  master_config {
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_master"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }

  worker_config {
    count         = 2
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_worker"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }

  ## config a different machine type
  worker_config {
    count         = 2
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4C"
    instance_name = "tf_kce_worker"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) The name of the cluster.
* `k8s_version` - (Required, ForceNew) The latest three kubernetes version. Current valid values:"v1.25.7", "v1.23.17", "v1.21.3". **Notes:** The version is updated in real time with the K8s official. Therefore, you can view the maintaining strategies in [Kingsoft Cloud K8s Version Strategies](https://docs.ksyun.com/documents/43229?type=3) and get the latest versions.
* `network_type` - (Required, ForceNew) The network type of the cluster. valid values: 'Flannel', 'Canal', 'Calico'.
* `pod_cidr` - (Required, ForceNew) The pod CIDR block.
* `reserve_subnet_id` - (Required, ForceNew) The ID of the reserve subnet.
* `service_cidr` - (Required, ForceNew) The service CIDR block.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `cluster_desc` - (Optional) The description of the cluster.
* `cluster_manage_mode` - (Optional, ForceNew) The management mode of the master node.
* `component` - (Optional) The component of the cluster.
* `expose_public_api_server` - (Optional) Whether to expose the apiserver to the public network, when cluster type is ManagedCluster, it works.
* `managed_cluster_multi_master` - (Optional) The configuration for the managed cluster multi master. If the cluster_manage_mode is ManagedCluster, this field is **required**.
* `master_config` - (Optional, ForceNew) The configuration for the master nodes. If the cluster_manage_mode is DedicatedCluster, this field is **required**. **Notes:** work_config block is identified by the **instance_type, subnet_id, security_group_id, role, image_id**. If the unique identification is the same, the instance config block is conflict, and then **cause an error**.If the unique identification is changed, that leads to the cluster **re-creation**.
* `master_etcd_separate` - (Optional, ForceNew) The deployment method for the Master and Etcd components of the cluster. if set to True, Deploy the Master and Etcd components on dedicated nodes. if set to false, Deploy the Master and Etcd components on shared nodes.
* `max_pod_per_node` - (Optional, ForceNew) The maximum number of pods that can be run on each node. valid values: 16, 32, 64, 128, 256.
* `public_api_server` - (Optional, ForceNew) Whether to expose the apiserver to the public network. If not needed, do not fill in this option. If selected, a public SLB and EIP will be created to enable public access to the cluster's API server. Users need to pass the Elastic IP creation pass-through parameter, which should be a JSON-formatted string.
* `worker_config` - (Optional, ForceNew) The configuration for the worker nodes. If the cluster_manage_mode is ManagedCluster, this field is **required**. **Notes:** work_config block is identified by the **instance_type, subnet_id, security_group_id, role, image_id**. If the unique identification is the same, the instance config block is conflict, and then **cause an error**.If the unique identification is changed, that leads to the cluster **re-creation**.

The `advanced_setting` object supports the following:

* `container_log_max_files` - (Optional, ForceNew) Customize the number of log files. The default value is 10.
* `container_log_max_size` - (Optional, ForceNew) Customize the maximum size of the log file. The default value is 100m.
* `container_path` - (Optional, ForceNew) The storage path of the container. The default value is /data/container. **Notes:** If this path is specified, the docker_path field will be ignored.
* `container_runtime` - (Optional, ForceNew) Container Runtime. Valid Values: `docker`, `containerd`.
* `data_disk` - (Optional, ForceNew) The mount setting of data disk. **Notes:** Only impact on the first data disk.
* `docker_path` - (Optional, ForceNew) The storage path of the container. The default value is /data/docker.
* `extra_arg` - (Optional, ForceNew) The extra arguments for the kubelet. The format is key=value. For example, --kubelet-extra-args="key1=value1,key2=value2".
* `label` - (Optional, ForceNew) 
* `pre_user_script` - (Optional, ForceNew) A user script encoded in base64, which will be executed on the node **before** the Kubernetes components run. Users need to ensure the script's re-entrant and retry logic. The script and its generated logs can be found in the directory /usr/local/ksyun/kce/pre_userscript.
* `taints` - (Optional, ForceNew) Taints.
* `user_script` - (Optional, ForceNew) A user script encoded in base64, which will be executed on the node **after** the Kubernetes components run. Users need to ensure the script's re-entrant and retry logic. The script and its generated logs can be found in the directory /usr/local/ksyun/kce/pre_userscript.

The `component` object supports the following:

* `name` - (Required) The name of the component.
* `namespace` - (Required) The namespace of the component.
* `release_name` - (Required) The name of the release.
* `config` - (Optional) The config of the component.
* `resources` - (Optional) The resources of the component.
* `version` - (Optional) The version of the component.

The `config` object supports the following:

* `config_string` - (Optional) The config string of the component.
* `virtual_kubelet` - (Optional) The config of the virtual kubelet.

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


The `kcilet_heartbeat` object supports the following:

* `failuret_threshold` - (Optional) The failure threshold of the virtual kubelet.
* `period` - (Optional) The period of the virtual kubelet.

The `limits` object supports the following:

* `cpu` - (Optional) The cpu of the component.
* `memory` - (Optional) The memory of the component.

The `managed_cluster_multi_master` object supports the following:

* `security_group_id` - (Required) The ID of the security group for the managed cluster masters.
* `subnet_id` - (Required) The ID of the subnet for the managed cluster masters.

The `master_config` object supports the following:

* `charge_type` - (Required, ForceNew) charge type of the instance.
* `count` - (Required, ForceNew) The number of master nodes. The count of master nodes must be 3 or 5.
* `image_id` - (Required, ForceNew) The ID for the image to use for the instance.
* `instance_type` - (Required, ForceNew) The type of instance to start. <br> - NOTE: it's may trigger this instance to power off, if instance type will be demotion.
* `security_group_id` - (Required, ForceNew) Security Group to associate with.
* `subnet_id` - (Required, ForceNew) The ID of subnet. the instance will use the subnet in the current region.
* `advanced_setting` - (Optional, ForceNew) Advanced settings.
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
* `instance_password` - (Optional) Password to an instance is a string of 8 to 32 characters.
* `instance_status` - (Optional) The state of instance.
* `keep_image_login` - (Optional) Keep the initial settings of the custom image.
* `key_id` - (Optional) The certificate id of the instance.
* `local_volume_snapshot_id` - (Optional, ForceNew) When the local data disk opens, the snapshot id is entered.
* `private_ip_address` - (Optional) Instance private IP address can be specified when you creating new instance.
* `project_id` - (Optional) The project instance belongs to.
* `purchase_time` - (Optional, ForceNew) The duration that you will buy the resource.
* `role` - (Optional) 
* `sriov_net_support` - (Optional, ForceNew) whether support networking enhancement.
* `sync_tag` - (Optional) Indicate whether to sync tags to instance.
* `system_disk` - (Optional) System disk parameters.
* `tags` - (Optional) the tags of the resource.
* `user_data` - (Optional) The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB. only effective when image support cloud-init.

The `openapi` object supports the following:

* `access_key` - (Optional) The access key of the virtual kubelet.
* `aksk_config_map` - (Optional) The aksk config map of the virtual kubelet.
* `region` - (Optional) The region of the virtual kubelet.
* `secret_key` - (Optional) The secret key of the virtual kubelet.

The `requests` object supports the following:

* `cpu` - (Optional) The cpu of the component.
* `memory` - (Optional) The memory of the component.

The `resources` object supports the following:

* `limits` - (Optional) The limits of the component.
* `requests` - (Optional) The requests of the component.

The `server` object supports the following:

* `listen_port` - (Optional) The port of the virtual kubelet server.
* `server_cert_file` - (Optional) The server certificate file of the virtual kubelet.
* `server_key_file` - (Optional) The server key file of the virtual kubelet.

The `system_disk` object supports the following:

* `disk_size` - (Optional) The size of the data disk. value range: [20, 500].
* `disk_type` - (Optional, ForceNew) System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk, `ESSD_SYSTEM_PL0`, The x7 machine type ESSD disk, `ESSD_SYSTEM_PL1`, The x7 machine type ESSD disk, `ESSD_SYSTEM_PL2`, The x7 machine type ESSD disk.

The `taints` object supports the following:

* `effect` - (Required) The effect of the taint. Valid values: NoSchedule, PreferNoSchedule, NoExecute.
* `key` - (Required) The key of the taint.
* `value` - (Required) The value of the taint.

The `taints` object supports the following:

* `effect` - (Required, ForceNew) The effect of the taint. Valid values: NoSchedule, PreferNoSchedule, NoExecute.
* `key` - (Required, ForceNew) The key of the taint.
* `value` - (Required, ForceNew) The value of the taint.

The `virtual_kubelet` object supports the following:

* `advanced_setting` - (Optional) The advanced settings of the virtual kubelet.
* `allow_privileged` - (Optional) Whether to allow privileged containers.
* `anonymous_auth` - (Optional) Whether to enable anonymous authentication.
* `batch_create_enable` - (Optional) Whether to enable batch creation of pods in the virtual kubelet.
* `capacity` - (Optional) The capacity of the virtual kubelet.
* `cluster_dns` - (Optional) The DNS of the virtual kubelet.
* `cluster_domain` - (Optional) The domain of the virtual kubelet.
* `cluster_id` - (Optional) The ID of the cluster.
* `custom_labels` - (Optional) The custom labels of the virtual kubelet.
* `disable_taint` - (Optional) Whether to disable taint.
* `enable_node_lease` - (Optional) Whether to enable node lease.
* `image_cache_enable` - (Optional) Whether to enable image cache in the virtual kubelet.
* `instance_settings` - (Optional) The instance settings of the virtual kubelet.
* `kci_pod_deletion_cost` - (Optional) The pod deletion cost of the virtual kubelet.
* `kcilet_heartbeat` - (Optional) The heartbeat of the virtual kubelet.
* `kcilet_kubeconfig_path` - (Optional) The kubeconfig path of the virtual kubelet.
* `kubeconfig` - (Optional) The config string of the virtual kubelet.
* `log_level` - (Optional) The log level of the virtual kubelet.
* `metrics_server` - (Optional) The metrics server of the virtual kubelet.
* `namespace` - (Optional) The namespace of the virtual kubelet.
* `nodename` - (Optional) The name of the virtual kubelet node.
* `openapi` - (Optional) The openapi of the virtual kubelet.
* `pod_sync_bucket_size` - (Optional) The bucket size of pod synchronization in the virtual kubelet.
* `pod_sync_rate` - (Optional) The rate of pod synchronization in the virtual kubelet.
* `pod_sync_workers` - (Optional) The number of workers for pod synchronization in the virtual kubelet.
* `server` - (Optional) The server of the virtual kubelet.
* `startup_timeout` - (Optional) The startup timeout of the virtual kubelet.
* `summary_sync_interval` - (Optional) The interval of summary synchronization in the virtual kubelet.
* `taints` - (Optional) Taints.
* `version` - (Optional) The version of the virtual kubelet.

The `worker_config` object supports the following:

* `charge_type` - (Required, ForceNew) charge type of the instance.
* `count` - (Required, ForceNew) The number of worker nodes.
* `image_id` - (Required, ForceNew) The ID for the image to use for the instance.
* `instance_type` - (Required, ForceNew) The type of instance to start. <br> - NOTE: it's may trigger this instance to power off, if instance type will be demotion.
* `security_group_id` - (Required, ForceNew) Security Group to associate with.
* `subnet_id` - (Required, ForceNew) The ID of subnet. the instance will use the subnet in the current region.
* `advanced_setting` - (Optional, ForceNew) Advanced settings.
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
* `instance_password` - (Optional) Password to an instance is a string of 8 to 32 characters.
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
* `cluster_id` - The ID of the cluster.
* `kube_config_intranet` - The configuration for the private kubernetes cluster.
* `kube_config` - The configuration for the kubernetes cluster.
* `master_id_list` - The ID list of the master nodes.
* `worker_id_list` - The ID list of the worker nodes.


## Import

KCE cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kce_cluster.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

