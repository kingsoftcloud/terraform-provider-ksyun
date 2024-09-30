/*
This data source provides a list of project resources.

# Example Usage

```hcl

	data "ksyun_iam_projects" "projects" {
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunIamProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunIamProjectsRead,

		Schema: map[string]*schema.Schema{

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"projects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the IAM ProjectId.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Project AccountId.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Project ProjectName.",
						},
						"project_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Project ProjectDesc.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN Project Status.",
						},
						"krn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN Project Krn.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN Role CreateDate.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunIamProjectsRead(d *schema.ResourceData, meta interface{}) error {
	iamProjectService := IamProjectService{meta.(*KsyunClient)}
	return iamProjectService.ReadAndSetIamProjects(d, dataSourceKsyunIamProjects())
}
