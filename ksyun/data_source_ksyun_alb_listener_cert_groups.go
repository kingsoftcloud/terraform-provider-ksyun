/*
This data source provides a list of ALB listener cert group resources according to their ID.

# Example Usage

```hcl

	data "ksyun_alb_listener_cert_groups" "default" {
		output_file="output_result"
		ids=[]
		alb_listener_id=[]
	}

```
*/
package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKsyunAlbListenerCertGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAlbListenerCertGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ALB Listener cert group IDs, all the ALB Listener cert groups belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"alb_listener_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more ALB Listener IDs.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ALB listener cert groups that satisfy the condition.",
			},
			"listener_cert_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ALB Listener cert groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the ALB Listener cert group.",
						},
						"alb_listener_cert_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ALB listener cert group.",
						},
						"alb_listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the ALB listener.",
						},
						"alb_listener_cert_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information list of ALB Listener certs. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The creation time.",
									},
									"certificate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the certificate.",
									},
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the certificate.",
									},
									"cert_authority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "certificate authority.",
									},
									"common_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The common name on the certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expire time of the certificate.",
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

func dataSourceKsyunAlbListenerCertGroupsRead(d *schema.ResourceData, meta interface{}) error {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	return s.ReadAndSetCertGroups(d, dataSourceKsyunAlbListenerCertGroups())
}
