/*
This data source provides a list of kec local volumes in the current region.

# Example Usage

```hcl

	data "ksyun_local_volumes" "default" {
	  output_file=""
	}

```
*/

package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKsyunLocalVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLocalVolumeRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the instance which the volume belong to.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of local volumes that satisfy the condition.",
			},
			"local_volume_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of KEC local volumes. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the local volume.",
						},
						"local_volume_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the local volume.",
						},
						"local_volume_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the local volume.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the instance.",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the instance.",
						},
						"local_volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the local volume.",
						},
						"local_volume_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The category of the local volume.",
						},
						"local_volume_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The size of the local volume.",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation date.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLocalVolumeRead(d *schema.ResourceData, meta interface{}) error {
	s := LocalVolumeService{meta.(*KsyunClient)}
	return s.ReadAndSetLocalVolumes(d, dataSourceKsyunLocalVolumes())
}
