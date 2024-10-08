/*
Provides a Iam Role resource.

# Example Usage

```hcl

resource "ksyun_iam_role" "role" {
  role_name = "role_name_test"
  trust_accounts = "2000096256"
  description = "desc"
}

```

# Import

IAM Role can be imported using the `role_name`, e.g.

```
$ terraform import ksyun_iam_role.role role_name
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunIamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunIamRoleCreate,
		Read:   resourceKsyunIamRoleRead,
		Delete: resourceKsyunIamRoleDelete,
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM RoleName.",
			},
			"trust_accounts": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM TrustAccounts.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM Description.",
			},
		},
	}
}

func resourceKsyunIamRoleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	iamRoleService := IamRoleService{meta.(*KsyunClient)}
	err = iamRoleService.CreateIamRole(d, resourceKsyunIamRole())
	if err != nil {
		return fmt.Errorf("error on creating IAM role %q, %s", d.Id(), err)
	}
	return resourceKsyunIamRoleRead(d, meta)
}

func resourceKsyunIamRoleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunIamRoleRead(d *schema.ResourceData, meta interface{}) (err error) {
	iamRoleService := IamRoleService{meta.(*KsyunClient)}
	err = iamRoleService.ReadAndSetIamRole(d, resourceKsyunIamRole())
	if err != nil {
		return fmt.Errorf("error on reading IAM role, %s", err)
	}
	return
}

func resourceKsyunIamRoleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	iamRoleService := IamRoleService{meta.(*KsyunClient)}
	err = iamRoleService.DeleteIamRole(d)
	if err != nil {
		return fmt.Errorf("error on deleting IAM role %q, %s", d.Id(), err)
	}
	return
}
