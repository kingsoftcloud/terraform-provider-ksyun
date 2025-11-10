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
				Description: "Pagination marker, e.g., limit=100&offset=0",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of KMR clusters",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of KMR clusters",
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
						"main_version": {
							Type:     schema.TypeString,
							Computed: true,
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
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_status": {
							Type:     schema.TypeString,
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
						"serving_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
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
