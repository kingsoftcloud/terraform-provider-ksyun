/*
This data source provides a list of Bare Metal Image resources according to their Bare Metal Image ID.

# Example Usage

```hcl
# Get  bare metal_images

	data "ksyun_bare_metal_images" "default" {
	  output_file="output_result"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunBareMetalImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBareMetalImagesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Bare Metal Images IDs, all the Bare Metal Images belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"image_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Bare Metal Images Types.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by name of Bare Metal Image.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Bare Metal Images that satisfy the condition.",
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for Bera Metal Image.",
						},
						"os_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS type of the Image.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS name of the Image.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of the Image.",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the Image.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunBareMetalImagesRead(d *schema.ResourceData, meta interface{}) error {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	return bareMetalService.ReadAndSetImages(d, dataSourceKsyunBareMetalImages())
}
