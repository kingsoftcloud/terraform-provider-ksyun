---
subcategory: "EIP"
layout: "ksyun"
page_title: "ksyun: ksyun_eips"
sidebar_current: "docs-ksyun-datasource-eips"
description: |-
  This data source provides a list of EIP resources (Elastic IP address) according to their EIP ID.
---

# ksyun_eips

This data source provides a list of EIP resources (Elastic IP address) according to their EIP ID.

#

## Example Usage

```hcl
data "ksyun_eips" "default" {
  output_file = "output_result"

  ids                  = []
  project_id           = []
  instance_type        = []
  network_interface_id = []
  internet_gateway_id  = []
  band_width_share_id  = []
  line_id              = []
  public_ip            = []
}
```

## Argument Reference

The following arguments are supported:

* `band_width_share_id` - (Optional) A list of BandWidthShare IDs.
* `ids` - (Optional) A list of Elastic IP IDs, all the EIPs belong to this region will be retrieved if the ID is `""`.
* `instance_type` - (Optional) A list of Instance Type.
* `internet_gateway_id` - (Optional) A list of InternetGateway IDs.
* `ip_version` - (Optional) IP Version.
* `line_id` - (Optional) A list of Line IDs.
* `network_interface_id` - (Optional) A list of NetworkInterface IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `project_id` - (Optional) One or more project IDs.
* `public_ip` - (Optional) A list of EIP address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eips` - An information list of EIP. Each element contains the following attributes:
  * `allocation_id` - ID of the EIP.
  * `band_width_share_id` - the ID of the BWS which the EIP associated.
  * `band_width` - bandwidth of the EIP.
  * `create_time` - creation time of the EIP.
  * `id` - ID of the EIP.
  * `instance_id` - The instance id to bind with the EIP.
  * `instance_type` - The instance type to bind with the EIP.
  * `internet_gateway_id` - InternetGateway ID.
  * `ip_version` - IP Version.
  * `line_id` - Line ID.
  * `network_interface_id` - NetworkInterface ID.
  * `project_id` - project ID.
  * `public_ip` - EIP address.
  * `state` - state of the EIP.
* `total_count` - Total number of Elastic IPs that satisfy the condition.


