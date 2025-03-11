/*
Provides a Iam Policy resource.

# Example Usage

```hcl

resource "ksyun_iam_policy" "policy" {
  policy_name = "TestPolicy1"
  policy_document = "{\"Version\": \"2015-11-01\",\"Statement\": [{\"Effect\": \"Allow\",\"Action\": [\"iam:List*\"],\"Resource\": [\"*\"]}]}"
}`

```

# Import

IAM Policy can be imported using the `policy_name`, e.g.

```
$ terraform import ksyun_iam_policy.policy policy_name
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunIamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunIamPolicyCreate,
		Read:   resourceKsyunIamPolicyRead,
		Update: resourceKsyunIamPolicyUpdate,
		Delete: resourceKsyunIamPolicyDelete,
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM PolicyName.",
			},
			"policy_document": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IAM PolicyDocument.",
			},
			"policy_krn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM PolicyKrn.",
			},
		},
	}
}

func resourceKsyunIamPolicyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	iamPolicyService := IamPolicyService{meta.(*KsyunClient)}
	err = iamPolicyService.CreateIamPolicy(d, resourceKsyunIamPolicy())
	if err != nil {
		return fmt.Errorf("error on creating IAM policy %q, %s", d.Id(), err)
	}
	return
}

func resourceKsyunIamPolicyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	iamPolicyService := IamPolicyService{meta.(*KsyunClient)}
	err = iamPolicyService.UpdateIamPolicy(d, resourceKsyunIamPolicy())
	if err != nil {
		return fmt.Errorf("error on updating IAM policy %q, %s", d.Id(), err)
	}
	return
}

func resourceKsyunIamPolicyRead(d *schema.ResourceData, meta interface{}) (err error) {
	iamPolicyService := IamPolicyService{meta.(*KsyunClient)}
	err = iamPolicyService.ReadAndSetIamPolicy(d, resourceKsyunIamPolicy())
	if err != nil {
		return fmt.Errorf("error on reading IAM policy, %s", err)
	}
	return
}

func resourceKsyunIamPolicyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	iamPolicyService := IamPolicyService{meta.(*KsyunClient)}
	err = iamPolicyService.DeleteIamPolicy(d)
	if err != nil {
		return fmt.Errorf("error on deleting IAM policy %q, %s", d.Id(), err)
	}
	return
}
