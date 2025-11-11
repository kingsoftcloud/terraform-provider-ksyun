/*
This data source provides a list of KLOG projects.

# Example Usage

```hcl

	data "ksyun_klog_projects" "default" {
		project_name="test"
		description="online"
		page="0"
		size="20"
		output_file="output_result"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKlogProjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunProjectsRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of project.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of project.",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page number start from 0.",
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page size, 1 - 500.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of project.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"projects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Project list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region of project.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of project.",
						},
						"iam_project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAMProjectName of project.",
						},
						"iam_project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The IAMProjectId of project.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the project was created.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the project was updated.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of project.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of project.",
						},
						"log_pool_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The log pool count of project.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tags of project.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The key of tag.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of tag.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunProjectsRead(d *schema.ResourceData, meta interface{}) error {
	s := KlogProjectService{meta.(*KsyunClient)}
	return s.ReadAndSetProjects(d, dataSourceKsyunKlogProjects())
}
