/*
This data source provides a list of namespace resources according to their instance id.

Example Usage

```hcl
data "ksyun_kcrs_namespaces" "foo" {
  output_file="kcrs_namespaces_output_result"
  instance_id = "86f14f8c-bf24-42c8-91bd-xxxxxxxx"
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKcrsNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKcrsNamespacesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,

				Description: "Kcrs Instance Id.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kcrs Instance namespace, all the Kcrs namespace belong to this instance will be retrieved if the namespaces is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of kcrs namespaces that satisfy the condition.",
			},

			"namespace_items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace.",
						},
						"public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether Public.",
						},
						"repo_count": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The count of Images in this repository.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKcrsNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	kcmService := KcrsService{meta.(*KsyunClient)}
	return kcmService.ReadAndSetKcrsNamespaces(d, dataSourceKsyunKcrsNamespaces())
}
