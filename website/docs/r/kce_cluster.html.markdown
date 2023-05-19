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

#

## Example Usage

```hcl
data "ksyun_kce_instance_images" "test" {
  output_file = "output_result"
}

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf_test_cluster"
  cluster_desc        = "description..."
  cluster_manage_mode = "DedicatedCluster"
  vpc_id              = ksyun_vpc.tf_test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.254.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.19.3"
  reserve_subnet_id   = ksyun_subnet.tf_test_reserve_subnet.id

  master_config {
    role          = "Master_Etcd"
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.tf_test_subnet.id
    security_group_id = [ksyun_security_group.default.id]
    charge_type       = "Daily"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required) The name of the cluster.
* `k8s_version` - (Required, ForceNew) kubernetes version, valid values:"v1.17.6", "v1.19.3", "v1.21.3".
* `master_config` - (Required) The configuration for the master nodes.
* `network_type` - (Required, ForceNew) The network type of the cluster. valid values: 'Flannel', 'Canal'.
* `pod_cidr` - (Required, ForceNew) The pod CIDR block.
* `reserve_subnet_id` - (Required, ForceNew) The ID of the reserve subnet.
* `service_cidr` - (Required, ForceNew) The service CIDR block.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `cluster_desc` - (Optional) The description of the cluster.
* `cluster_manage_mode` - (Optional, ForceNew) The management mode of the master node.
* `master_etcd_separate` - (Optional, ForceNew) The deployment method for the Master and Etcd components of the cluster. if set to True, Deploy the Master and Etcd components on dedicated nodes. if set to false, Deploy the Master and Etcd components on shared nodes.
* `max_pod_per_node` - (Optional, ForceNew) The maximum number of pods that can be run on each node. valid values: 16, 32, 64, 128, 256.
* `public_api_server` - (Optional, ForceNew) Whether to expose the apiserver to the public network. If not needed, do not fill in this option. If selected, a public SLB and EIP will be created to enable public access to the cluster's API server. Users need to pass the Elastic IP creation pass-through parameter, which should be a JSON-formatted string.

The `data_disks` object supports the following:

* `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on EBS disk.
* `disk_size` - (Optional, ForceNew) Data disk size. value range: [10, 16000].
* `disk_snapshot_id` - (Optional, ForceNew) When the cloud disk opens, the snapshot id is entered.
* `disk_type` - (Optional, ForceNew) Data disk type.

The `extension_network_interface` object supports the following:


The `master_config` object supports the following:

* `charge_type` - (Required, ForceNew) charge type of the instance.
* `count` - (Required) 
* `image_id` - (Required) The ID for the image to use for the instance.
* `role` - (Required) Node role. when the MasterEtcdSeparate field is set to false, both the Worker and Master_Etcd roles need to be specified.when the MasterEtcdSeparate field is set to true, the Master, Etcd, and Worker roles need to be specified simultaneously.
* `security_group_id` - (Required) Security Group to associate with.
* `subnet_id` - (Required) The ID of subnet. the instance will use the subnet in the current region.
* `auto_create_ebs` - (Optional) Whether to create EBS volumes from snapshots in the custom image, default is false.
* `data_disk_gb` - (Optional) The size of the local SSD disk.
* `data_disks` - (Optional) The list of data disks created with instance.
* `data_guard_id` - (Optional, ForceNew) Add instance being created to a disaster tolerance group.
* `dns1` - (Optional) DNS1 of the primary network interface.
* `dns2` - (Optional) DNS2 of the primary network interface.
* `force_delete` - (Optional) Indicate whether to delete instance directly or not.
* `force_reinstall_system` - (Optional) Indicate whether to reinstall system.
* `host_name` - (Optional) The hostname of the instance. only effective when image support cloud-init.
* `iam_role_name` - (Optional) name of iam role.
* `instance_name` - (Optional) The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.
* `instance_password` - (Optional) Password to an instance is a string of 8 to 32 characters.
* `instance_status` - (Optional) The state of instance.
* `instance_type` - (Optional) The type of instance to start.
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

The `system_disk` object supports the following:

* `disk_size` - (Optional) The size of the data disk. value range: [20, 500].
* `disk_type` - (Optional, ForceNew) System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - The ID of the cluster.


## Import

KCE cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kce_cluster.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

