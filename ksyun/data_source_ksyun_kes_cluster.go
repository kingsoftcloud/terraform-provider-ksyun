/*
Provides a KES Cluster data source.

# Example Usage

```hcl

	provider "ksyun" {
		region = "cn-beijing-6"
	}

	data "ksyun_kes_cluster" "default" {
		cluster_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}

	output "cluster_name" {
		value = data.ksyun_kes_cluster.default.clusters.0.cluster_name
	}

```

# Argument Reference

* `cluster_id` - (Required) The ID of the KES cluster to query.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

# Attributes Reference

* `clusters` - A list containing a single KES cluster detail.
  - `cluster_id` - The ID of the cluster.
  - `cluster_name` - The name of the cluster.
  - `cluster_type` - The type of the cluster.
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
  - `proxy_port` - The proxy port of the cluster.
  - `instance_groups` - List of instance groups.
    - `id` - The ID of the instance group.
    - `instance_group_type` - The type of the instance group.
    - `resource_type` - The resource type.
    - `instance_type` - The instance type.
    - `volume_size` - The volume size.
    - `volume_type` - The volume type.
  - `tags` - List of tags.
    - `tag_key` - The key of the tag.
    - `tag_value` - The value of the tag.
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKesCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKesClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the KES cluster to query.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list containing a single KES cluster detail.",
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
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the cluster.",
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
									"volume_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The volume size.",
									},
									"volume_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The volume type.",
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
						"proxy_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The proxy port of the cluster.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key of the tag.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the tag.",
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

func dataSourceKsyunKesClusterRead(d *schema.ResourceData, meta interface{}) (err error) {
	kesService := KesService{meta.(*KsyunClient)}
	return kesService.ReadAndSetCluster(d, dataSourceKsyunKesCluster())
}
