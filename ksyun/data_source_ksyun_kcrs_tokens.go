/*
This data source provides a list of token resources according to their instance id.

Example Usage

```hcl
data "ksyun_kcrs_tokens" "foo" {
  output_file="kcrs_tokens_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-xxxxxxx"
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKcrsTokens() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKcrsTokensRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,

				Description: "Kcrs Instance Id.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of kcrs tokens that satisfy the condition.",
			},

			"tokens": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the token.",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Enable.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description for this token.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Time.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expired Time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKcrsTokensRead(d *schema.ResourceData, meta interface{}) error {
	kcmService := KcrsService{meta.(*KsyunClient)}
	return kcmService.ReadAndSetKcrsTokens(d, dataSourceKsyunKcrsTokens())
}
