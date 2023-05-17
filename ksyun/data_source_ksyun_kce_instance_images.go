package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKsyunKceInstanceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKceInstanceImagesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_type": {
							Type:     schema.TypeString,
							Computed: true,
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
