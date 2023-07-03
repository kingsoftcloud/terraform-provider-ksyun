/*
This data source provides a list of EBS volumes.

# Example Usage

```hcl

	data "ksyun_volumes" "default" {
	  output_file="output_result"
	  ids=[]
	  volume_category=""
	  volume_status=""
	  volume_type=""
	  availability_zone=""
	}

```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunVolumesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of EBS volumes that satisfy the condition.",
			},
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				MaxItems:    100,
				Description: "A list of EBS IDs, all the EBS resources belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"volume_category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The category to which the EBS volume belongs.",
			},
			"volume_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the EBS volume.",
			},
			"volume_create_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The time when the EBS volume was created.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the EBS volume.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The availability zone in which the EBS volume resides.",
			},
			"volumes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of EBS volumes. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the EBS volume.",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the EBS volume.",
						},
						"volume_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the EBS volume.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity of the EBS volume.",
						},
						"volume_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the EBS volume.",
						},
						"volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the EBS volume.",
						},
						"volume_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The category of the EBS volume.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the KEC instance to which the EBS volume is to be attached.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the EBS volume was created.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone in which the EBS volume resides.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the project.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunVolumesRead(d *schema.ResourceData, meta interface{}) error {
	ebsService := EbsService{meta.(*KsyunClient)}
	return ebsService.ReadAndSetVolumes(d, dataSourceKsyunVolumes())
}
