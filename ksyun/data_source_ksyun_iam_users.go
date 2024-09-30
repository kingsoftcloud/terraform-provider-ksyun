/*
This data source provides a list of user resources.

# Example Usage

```hcl

		data "ksyun_iam_users" "users" {
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunIamUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunIamUsersRead,

		Schema: map[string]*schema.Schema{

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the IAM User Id.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the IAM UserId.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM User Path.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM UserName.",
						},
						"real_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM User RealName.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User CreateDate.",
						},
						"phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User Phone.",
						},
						"country_mobile_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User CountryMobileCode.",
						},
						"is_international": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN User IsInternational.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User Email.",
						},
						"phone_verified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User PhoneVerified.",
						},
						"email_verified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User EmailVerified.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User Remark.",
						},
						"krn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User Krn.",
						},
						"password_reset_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "IAN User PasswordResetRequired.",
						},
						"enable_mfa": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IAN User EnableMFA.",
						},
						"update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAN User UpdateDate.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunIamUsersRead(d *schema.ResourceData, meta interface{}) error {
	iamUserService := IamUserService{meta.(*KsyunClient)}
	return iamUserService.ReadAndSetIamUsers(d, dataSourceKsyunIamUsers())
}
