/*
This data source provides a list of instance resources according to their id.

Example Usage

```hcl
data "ksyun_kcrs_instances" "foo" {
  output_file="kcrs_instance_output_result"
  ids = []
}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunKcrsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKcrsInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Kcrs Instance IDs, all the Kcrs Instances belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more project IDs. If its value is none, returns instance all of project.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by instance name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of kcrs instances that satisfy the condition.",
			},

			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of the instance.",
						},
						"instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status.",
						},
						"internal_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal endpoint dns.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created Time.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge Type.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project Id.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKcrsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	kcmService := KcrsService{meta.(*KsyunClient)}
	return kcmService.ReadAndSetKcrsInstances(d, dataSourceKsyunKcrsInstances())
}
