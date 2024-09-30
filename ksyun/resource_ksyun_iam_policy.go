/*
Provides a Iam Policy resource.

# Example Usage

```hcl

resource "ksyun_iam_policy" "policy" {
  policy_name = "TestPolicy3"
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
				Description: "IAM PolicyName.",
			},
			"policy_document": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IAM PolicyDocument.",
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
	return
}

func resourceKsyunIamPolicyRead(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunIamPolicyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return
}
