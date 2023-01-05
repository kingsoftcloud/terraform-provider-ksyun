/*
This data source provides a list of Certificate resources (KCM) according to their ID.

Example Usage

```hcl
data "ksyun_certificates" "default" {
  output_file="output_result"
  ids = ["c7b2ba05-9302-4933-8588-a66f920ff57d"]
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunCertificatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Certificate IDs, all the Certificates belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by certificate name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of certificates that satisfy the condition.",
			},

			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the certificate.",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the certificate.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	kcmService := KcmService{meta.(*KsyunClient)}
	return kcmService.ReadAndSetCertificates(d, dataSourceKsyunCertificates())
}
