/*
This data source provides a list of Bare Metal Raid Attributes resources according to their Bare Metal Raid Attribute ID.

# Example Usage

```hcl
# Get  bare metal_raid_attributes

	data "ksyun_bare_metal_raid_attributes" "default" {
	  output_file="output_result"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunBareMetalRaidAttributes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBareMetalRaidAttributesRead,
		Schema: map[string]*schema.Schema{
			"host_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Bare Metal Raid Attribute Host Types.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by name of Bare Metal Raid template.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Bare Metal Raid Attributes that satisfy the condition.",
			},
			"raid_attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Total number of Bare Metal Raid Attributes that satisfy the condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Bare Metal Raid template.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for Bare Metal Raid template.",
						},
						"raid_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the raid.",
						},
						"host_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "host type of the Bare Metal.",
						},
						"disk_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of disks that used raid template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the disk.",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the disk.",
									},
									"raid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "raid level.",
									},
									"disk_attribute": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "attribute of the disk.",
									},
									"disk_count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "count of disks.",
									},
									"disk_space": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "space of the data disk.",
									},
									"space": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "available Space.",
									},
									"system_disk_space": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "space of the system disk.",
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

func dataSourceKsyunBareMetalRaidAttributesRead(d *schema.ResourceData, meta interface{}) error {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	return bareMetalService.ReadAndSetRaidAttributes(d, dataSourceKsyunBareMetalRaidAttributes())
}
