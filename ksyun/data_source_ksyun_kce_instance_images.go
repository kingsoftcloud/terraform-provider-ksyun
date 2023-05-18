/*
This data source providers a list of available instance image which support kce.

# Example Usage

```hcl

	data "ksyun_kce_instance_images" "default" {
	  output_file = "output_result"
	}

```
*/

package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKsyunKceInstanceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKceInstanceImagesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"image_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the image.",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the image.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKceInstanceImagesRead(d *schema.ResourceData, meta interface{}) error {
	kceService := KceService{meta.(*KsyunClient)}
	return kceService.ReadAndSetKceInstanceImages(d, dataSourceKsyunKceInstanceImages())
}
