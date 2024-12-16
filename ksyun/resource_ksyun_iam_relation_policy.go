/*
Provides a Iam Policy resource.

# Example Usage

```hcl
resource "ksyun_iam_relation_policy" "user" {
  name = "iam_user_name"
  policy_name = "IAMReadOnlyAccess"
  relation_type = 1
  policy_type = "system"
}`

resource "ksyun_iam_relation_policy" "user" {
  name = "iam_role_name"
  policy_name = "IAMReadOnlyAccess"
  relation_type = 2
  policy_type = "system"
}`

```

# Import

IAM Policy can be imported using the `policy_name`, e.g.

```
$ terraform import ksyun_iam_relation_policy.user
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunIamRelationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunIamRelationPolicyCreate,
		Read:   resourceKsyunIamRelationPolicyRead,
		Delete: resourceKsyunIamRelationPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM UserName or RoleName according to relation type.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM PolicyName.",
			},
			"relation_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "relation type 1 is the user,relation type 2 is the role.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "policy type system is the system policy,policy type custom is the custom policy.",
			},
		},
	}
}

func resourceKsyunIamRelationPolicyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	iamRelationPolicyService := IamRelationPolicyService{meta.(*KsyunClient)}
	err = iamRelationPolicyService.CreateIamRelationPolicy(d, resourceKsyunIamRelationPolicy())
	if err != nil {
		return fmt.Errorf("error on creating IAM reliaton policy %q, %s", d.Id(), err)
	}
	return
}

func resourceKsyunIamRelationPolicyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunIamRelationPolicyRead(d *schema.ResourceData, meta interface{}) (err error) {
	iamRelationPolicyService := IamRelationPolicyService{meta.(*KsyunClient)}
	err = iamRelationPolicyService.ReadAndSetIamRelationPolicy(d, resourceKsyunIamRelationPolicy())
	if err != nil {
		return fmt.Errorf("error on reading IAM reliaton policy, %s", err)
	}
	return
}

func resourceKsyunIamRelationPolicyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	iamRelationPolicyService := IamRelationPolicyService{meta.(*KsyunClient)}
	err = iamRelationPolicyService.DeleteIamRelationPolicy(d)
	if err != nil {
		return fmt.Errorf("error on deleting IAM reliaton policy %q, %s", d.Id(), err)
	}
	return
}
