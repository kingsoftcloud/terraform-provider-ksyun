/*
This data source provides a list of tag resources.

# Example Usage

```hcl

	data "ksyun_tags" "default" {
	  output_file="output_result"

	  # optional
	  # eg. key = ["tag_key1", "tag_key2", ...]
	  keys = []
	  # optional
	  # eg. value = ["tag_value1", ...]
	  values = []
	  # optional
	  # eg. resource_type = ["kec-instance", "eip", ...]
	  resource_types = []
	  # optional
	  # eg. key = ["instance_uuid", ...]
	  resource_ids = []

}
```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunTagsRead,

		Schema: map[string]*schema.Schema{

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"keys": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of tag keys.",
			},
			"values": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of tag values.",
			},
			"resource_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of resource types.",
			},
			"resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of resource ids.",
			},

			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of tag.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the tag.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag value.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunTagsRead(d *schema.ResourceData, meta interface{}) error {
	tagService := TagV1Service{meta.(*KsyunClient)}
	return tagService.ReadAndSetTags(d, dataSourceKsyunTags())
}
