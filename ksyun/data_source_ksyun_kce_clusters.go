package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKceClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKceClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_manage_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"k8s_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_kmse": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
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
