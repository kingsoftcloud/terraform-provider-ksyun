/*
This data source provides a list of Load Balancer Rule resources according to their Load Balancer Rule ID.

# Example Usage

```hcl

	data "ksyun_lb_acls" "default" {
	  output_file="output_result"
	  ids=[]

}
```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunSlbAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSlbAclsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of LB Rule IDs, all the LB Rules belong to the Load Balancer listener will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of LB Rules that satisfy the condition.",
			},
			"lb_acls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for LB ACL.",
						},
						"load_balancer_acl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the LB ACL.",
						},
						"load_balancer_acl_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the LB ACL.",
						},
						"load_balancer_acl_entry_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of ACL entries.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_acl_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the ACL.",
									},
									"load_balancer_acl_entry_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the ACL entry.",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The information of Acl's cidr block.",
									},
									"rule_number": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "rule priority.",
									},
									"rule_action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rule action, allow or deny.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "rul protocol.",
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
func dataSourceKsyunSlbAclsRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetLoadBalancerAcls(d, dataSourceKsyunSlbAcls())
}
