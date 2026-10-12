/*
Provides a KES Clusters data source.

# Example Usage

```hcl

	provider "ksyun" {
		region = "cn-beijing-6"
	}

	data "ksyun_kes_clusters" "default" {
		marker = "limit=10&offset=0"
	}

	output "total" {
		value = data.ksyun_kes_clusters.default.total
	}

```

# Argument Reference

* `marker` - (Optional) Pagination marker, e.g., limit=100&offset=0.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

# Attributes Reference

* `total` - Total number of KES clusters.
* `clusters` - List of KES clusters.
  - `cluster_id` - The ID of the cluster.
  - `cluster_name` - The name of the cluster.
  - `main_version` - The main version of the cluster.
  - `enable_eip` - Whether EIP is enabled.
  - `region` - The region of the cluster.
  - `availability_zone` - The availability zone of the cluster.
  - `vpc_domain_id` - The VPC domain ID.
  - `vpc_subnet_id` - The VPC subnet ID.
  - `charge_type` - The charge type.
  - `cluster_status` - The status of the cluster.
  - `create_time` - The creation time.
  - `update_time` - The update time.
  - `serving_minutes` - The serving minutes.
  - `instance_groups` - List of instance groups.
    - `id` - The ID of the instance group.
    - `instance_group_type` - The type of the instance group.
    - `resource_type` - The resource type.
    - `instance_type` - The instance type.
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKesClustersRead,
		Schema: map[string]*schema.Schema{
			"marker": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "limit=100&offset=0",
				Description: "Pagination marker, e.g., limit=100&offset=0.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of KES clusters.",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of KES clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the cluster.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cluster.",
						},
						"main_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The main version of the cluster.",
						},
						"instance_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of instance groups.",
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
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance type.",
									},
								},
							},
						},
						"enable_eip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether EIP is enabled.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region of the cluster.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone of the cluster.",
						},
						"vpc_domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC domain ID.",
						},
						"vpc_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC subnet ID.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type.",
						},
						"cluster_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the cluster.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time.",
						},
						"serving_minutes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The serving minutes.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKesClustersRead(d *schema.ResourceData, meta interface{}) (err error) {
	kesService := KesService{meta.(*KsyunClient)}
	return kesService.ReadAndSetClusters(d, dataSourceKsyunKesClusters())
}
