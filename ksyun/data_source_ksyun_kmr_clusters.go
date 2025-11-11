/*
Provides a KMR Clusters data source.

# Example Usage

```hcl

	data "ksyun_kmr_clusters" "default" {
	  marker = "limit=10&offset=0"
	}

	output "total" {
	  value = data.ksyun_kmr_clusters.default.total
	}

```

# Argument Reference

* `marker` - (Optional) Pagination marker, e.g., limit=100&offset=0.
* `output_file` - (Optional) File name where to save data source results.

# Attributes Reference

* `total` - Total number of KMR clusters.
* `clusters` - List of KMR clusters.
  - `cluster_id` - The ID of the cluster.
  - `cluster_name` - The name of the cluster.
  - `main_version` - The main version of the cluster.
  - `enable_eip` - Whether EIP is enabled.
  - `region` - The region of the cluster.
  - `vpc_domain_id` - The VPC domain ID.
  - `charge_type` - The charge type.
  - `cluster_status` - The status of the cluster.
  - `create_time` - The creation time.
  - `update_time` - The update time.
  - `serving_minutes` - The serving minutes.
  - `instance_groups` - List of instance groups.
  - `id` - The ID of the instance group.
  - `instance_group_type` - The type of the instance group.
  - `instance_type` - The instance type.
  - `resource_type` - The resource type.
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKmrClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKmrClustersRead,
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
				Description: "Total number of KMR clusters.",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of KMR clusters.",
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
						"vpc_domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC domain ID.",
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

func dataSourceKsyunKmrClustersRead(d *schema.ResourceData, meta interface{}) error {
	kmrService := KmrService{meta.(*KsyunClient)}
	return kmrService.ReadAndSetClusters(d, dataSourceKsyunKmrClusters())
}
