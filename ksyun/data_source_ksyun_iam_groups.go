/*
This data source provides a list of group resources.

# Example Usage

```hcl

		data "ksyun_iam_groups" "groups" {
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunIamGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunIamGroupsRead,

		Schema: map[string]*schema.Schema{

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IAM GroupId.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Group Path.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM GroupName.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Group Description.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN Group CreateDate.",
						},
						"krn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN Group Krn.",
						},
						"user_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN Group UserCount.",
						},
						"policy_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN Group PolicyCount.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunIamGroupsRead(d *schema.ResourceData, meta interface{}) error {
	iamGroupService := IamGroupService{meta.(*KsyunClient)}
	return iamGroupService.ReadAndSetIamGroups(d, dataSourceKsyunIamGroups())
}
