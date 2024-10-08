/*
Provides a Iam User resource.

# Example Usage

```hcl

resource "ksyun_iam_user" "user" {
  user_name = "username01"
  real_name = "realname01"
  phone = "13800000000"
  email = "test@ksyun.com"
  remark = "remark"
  password = "password"
  password_reset_required = 0
  open_login_protection = 1
  open_security_protection = 1
  view_all_project = 0
}

```

# Import

IAM User can be imported using the `user_name`, e.g.

```
$ terraform import ksyun_iam_user.user user_name
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunIamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunIamUserCreate,
		Read:   resourceKsyunIamUserRead,
		Delete: resourceKsyunIamUserDelete,
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM UserName.",
			},
			"real_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM RealName.",
			},
			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM Phone.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM Email.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM Remark.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM Password.",
			},
			"password_reset_required": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Does IAM user login reset password.",
			},
			"open_login_protection": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Does IAM user enable login protection.",
			},
			"open_security_protection": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Does IAM user enable operation protection.",
			},
			"view_all_project": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Can IAM users view all projects.",
			},
		},
	}
}

func resourceKsyunIamUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	iamUserService := IamUserService{meta.(*KsyunClient)}
	err = iamUserService.CreateIamUser(d, resourceKsyunIamUser())
	if err != nil {
		return fmt.Errorf("error on creating IAM user %q, %s", d.Id(), err)
	}
	return resourceKsyunIamUserRead(d, meta)
}

func resourceKsyunIamUserUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunIamUserRead(d *schema.ResourceData, meta interface{}) (err error) {
	iamUserService := IamUserService{meta.(*KsyunClient)}
	err = iamUserService.ReadAndSetIamUser(d, resourceKsyunIamUser())
	if err != nil {
		return fmt.Errorf("error on reading IAM user, %s", err)
	}
	return
}

func resourceKsyunIamUserDelete(d *schema.ResourceData, meta interface{}) (err error) {
	iamUserService := IamUserService{meta.(*KsyunClient)}
	err = iamUserService.DeleteIamUser(d)
	if err != nil {
		return fmt.Errorf("error on deleting IAM user %q, %s", d.Id(), err)
	}
	return
}
