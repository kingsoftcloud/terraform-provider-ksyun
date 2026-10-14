/*
Provides a KES vector instance groups (KES8) data source.

# Example Usage

```hcl

	provider "ksyun" {
		region = "cn-beijing-6"
	}

	data "ksyun_kes_instance_groups_vector" "default" {
		cluster_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}

	output "first_instance_group_id" {
		value = data.ksyun_kes_instance_groups_vector.default.instance_groups.0.id
	}

```

# Argument Reference

* `cluster_id` - (Required) The ID of the KES cluster whose instance groups are queried.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

# Attributes Reference

* `instance_groups` - List of instance groups of the cluster.
  - `id` - The ID of the instance group.
  - `instance_group_type` - The type of the instance group.
  - `resource_type` - The resource type.
  - `instance_type_code` - The instance type code.
  - `instance_count` - The instance count.
  - `volume_type` - The volume type.
  - `volume_size` - The volume size.
  - `volume_count` - The volume count.
  - `vpc_id` - The VPC ID.
  - `vpc_subnet_id` - The VPC subnet ID.
  - `avalability_zone` - The availability zone.
  - `multi_instance_count` - The multi instance count.
  - `instances` - List of instances.
    - `id` - The ID of the instance.
    - `instance_group_id` - The instance group ID of the instance.
    - `instance_id` - The instance ID.
    - `instance_name` - The instance name.
    - `internal_ip` - The internal IP.
    - `volumes` - List of volume identifiers.
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKesInstanceGroupsVector() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKesInstanceGroupsVectorRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the KES cluster whose instance groups are queried.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"instance_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of instance groups of the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the instance group.",
						},
						"instance_group_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the instance group.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"instance_type_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance type code.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instance count.",
						},
						"volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The volume type.",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The volume size.",
						},
						"volume_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The volume count.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID.",
						},
						"vpc_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC subnet ID.",
						},
						"avalability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone.",
						},
						"multi_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The multi instance count.",
						},
						"instances": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of instances.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the instance.",
									},
									"instance_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance group ID of the instance.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance name.",
									},
									"internal_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The internal IP.",
									},
									"volumes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "List of volume identifiers.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKesInstanceGroupsVectorRead(d *schema.ResourceData, meta interface{}) (err error) {
	kesService := KesService{meta.(*KsyunClient)}
	return kesService.ReadAndSetInstanceGroups(d, dataSourceKsyunKesInstanceGroupsVector())
}
