/*
Query available storage cluster information by UID.

# Example Usage

```hcl

	data "ksyun_kpfs_clusters" "default" {
		output_file="output_result"
		region         = "cn-qingyangtest-1"
		s_roce_cluster = "QYYC01-Sroce-Cluster-01"
		store_class  = "KPFS-P-S01"
		store_pool_type = "KPFS-P1"
		avail_zone     = "cn-qingyangtest-1a"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKpfsClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKpfsClustersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region of the KPFS cluster.",
			},
			"avail_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The availability zone of the KPFS cluster.",
			},
			"store_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The storage classes supported by the KPFS cluster.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"store_pool_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The storage pool type of the KPFS cluster.",
			},
			"s_roce_cluster": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SRoCE cluster name of the KPFS cluster.",
			},
			"cluster_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique code of the KPFS cluster.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of KPFS clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region of the KPFS cluster.",
						},
						"avail_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone of the KPFS cluster.",
						},
						"store_classes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The storage classes supported by the KPFS cluster.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"store_pool_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The storage pool type of the KPFS cluster.",
						},
						"s_roce_cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The SRoCE cluster name of the KPFS cluster.",
						},
						"cluster_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique code of the KPFS cluster.",
						},
					},
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
		},
	}
}

func dataSourceKsyunKpfsClustersRead(d *schema.ResourceData, meta interface{}) error {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	return kpfsService.ReadKpfsClusterList(d, dataSourceKsyunKpfsClusters())
}
