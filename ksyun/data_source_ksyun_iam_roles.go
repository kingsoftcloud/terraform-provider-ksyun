/*
This data source provides a list of role resources.

# Example Usage

```hcl

	data "ksyun_iam_roles" "roles" {
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunIamRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunIamRolesRead,

		Schema: map[string]*schema.Schema{

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IAM RoleId.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM RoleName.",
						},
						"krn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Role Krn.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM Role Description.",
						},
						"trust_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN Role TrustType.",
						},
						"trust_accounts": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN Role TrustAccounts.",
						},
						"trust_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN Role TrustProvider.",
						},
						"service_role_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN Role ServiceRoleType.",
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
func dataSourceKsyunIamRolesRead(d *schema.ResourceData, meta interface{}) error {
	iamRoleService := IamRoleService{meta.(*KsyunClient)}
	return iamRoleService.ReadAndSetIamRoles(d, dataSourceKsyunIamRoles())
}
