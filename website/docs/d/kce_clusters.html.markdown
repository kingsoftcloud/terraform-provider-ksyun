---
subcategory: "KCE"
layout: "ksyun"
page_title: "ksyun: ksyun_kce_clusters"
sidebar_current: "docs-ksyun-datasource-kce_clusters"
description: |-
  This data source providers a list of kce cluster resources according to their instance ID.
---

# ksyun_kce_clusters

This data source providers a list of kce cluster resources according to their instance ID.

#

## Example Usage

```hcl
data "ksyun_kce_clusters" "default" {
  output_file = "output_result"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) The id of the cluster.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_set` - a list of kce clusters.
  * `cluster_desc` - The description of the cluster.
  * `cluster_id` - The ID of the cluster.
  * `cluster_manage_mode` - The management mode of the master node.
  * `cluster_name` - The name of the cluster.
  * `cluster_type` - The type of the cluster.
  * `create_time` - The creation time.
  * `enable_kmse` - Whether to support KMSE.
  * `k8s_version` - Kubernetes version.
  * `network_type` - The network type of the cluster.
  * `node_num` - The number of nodes.
  * `pod_cidr` - The pod CIDR block.
  * `service_cidr` - The service CIDR block.
  * `status` - The status of the cluster.
  * `update_time` - The updated time.
  * `vpc_cidr` - The CIDR block of the VPC.
  * `vpc_id` - The ID of the VPC.


