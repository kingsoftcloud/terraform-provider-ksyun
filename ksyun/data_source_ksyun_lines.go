/*
This data source provides a list of line resources supported.

# Example Usage

```hcl

	data "ksyun_lines" "default" {
	  output_file="output_result"
	  line_name="BGP"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunLines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLinesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of lines, all the lines belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"line_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the line.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of lines that satisfy the condition.",
			},

			"lines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All the lines according the argument.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the line.",
						},
						"line_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the line.",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the line.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLinesRead(d *schema.ResourceData, meta interface{}) error {
	eipService := EipService{meta.(*KsyunClient)}
	return eipService.ReadAndSetLines(d, dataSourceKsyunLines())
}
