/*
This data source providers a list of kce cluster resources according to their instance ID.

# Example Usage

```hcl

	data "ksyun_kce_clusters" "default" {
	  output_file = "output_result"
	}

```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKceClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKceClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the cluster.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"cluster_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of kce clusters.",
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
						"cluster_manage_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The management mode of the master node.",
						},
						"k8s_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kubernetes version.",
						},
						"cluster_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the cluster.",
						},
						"pod_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The pod CIDR block.",
						},
						"service_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service CIDR block.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the VPC.",
						},
						"vpc_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block of the VPC.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the cluster.",
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of nodes.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The updated time.",
						},
						"enable_kmse": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to support KMSE.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network type of the cluster.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKceClustersRead(d *schema.ResourceData, meta interface{}) error {
	kceService := KceService{meta.(*KsyunClient)}
	return kceService.ReadAndSetKceClusters(d, dataSourceKsyunKceClusters())
}
