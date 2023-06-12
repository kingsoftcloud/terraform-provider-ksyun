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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - The ID of the cluster.
* `master_config` - The configuration for the master nodes.
  * `data_disk_gb` - The size of the local SSD disk.
  * `data_disks` - The list of data disks created with instance.
    * `disk_id` - ID of the disk.
    * `disk_size` - Data disk size. value range: [10, 16000].
    * `disk_type` - Data disk type.
  * `dns1` - DNS1 of the primary network interface.
  * `dns2` - DNS2 of the primary network interface.
  * `extension_network_interface` - extension network interface information.
    * `network_interface_id` - ID of the extension network interface.
  * `has_modify_keys` - whether the certificate key has modified.
  * `has_modify_password` - whether the password has modified.
  * `has_modify_system_disk` - whether the system disk has modified.
  * `host_name` - The hostname of the instance. only effective when image support cloud-init.
  * `instance_id` - ID of the instance.
  * `instance_name` - The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.
  * `instance_password` - Password to an instance is a string of 8 to 32 characters.
  * `instance_status` - The state of instance.
  * `instance_type` - The type of instance to start.
  * `key_id` - The certificate id of the instance.
  * `local_volume_snapshot_id` - When the local data disk opens, the snapshot id is entered.
  * `network_interface_id` - ID of the network interface.
  * `private_ip_address` - Instance private IP address can be specified when you creating new instance.
  * `project_id` - The project instance belongs to.
  * `sriov_net_support` - whether support networking enhancement.
  * `system_disk` - System disk parameters.
    * `disk_size` - The size of the data disk. value range: [20, 500].
    * `disk_type` - System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk.
  * `tags` - the tags of the resource.


## Import

KCE cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kce_cluster.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

