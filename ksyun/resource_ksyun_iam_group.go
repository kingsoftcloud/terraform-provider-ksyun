/*
Provides a Iam Group resource.

# Example Usage

```hcl

resource "ksyun_iam_group" "group" {
  group_name = "GroupNameTest"
  description = "desc"
}

```

# Import

IAM Group can be imported using the `group_name`, e.g.

```
$ terraform import ksyun_iam_group.group group_name
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunIamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunIamGroupCreate,
		Read:   resourceKsyunIamGroupRead,
		Delete: resourceKsyunIamGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IAM GroupName.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IAM Group Description.",
			},
		},
	}
}

func resourceKsyunIamGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	iamGroupService := IamGroupService{meta.(*KsyunClient)}
	err = iamGroupService.CreateIamGroup(d, resourceKsyunIamGroup())
	if err != nil {
		return fmt.Errorf("error on creating IAM group %q, %s", d.Id(), err)
	}
	return resourceKsyunIamGroupRead(d, meta)
}

func resourceKsyunIamGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunIamGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	iamGroupService := IamGroupService{meta.(*KsyunClient)}
	err = iamGroupService.ReadAndSetIamGroup(d, resourceKsyunIamGroup())
	if err != nil {
		return fmt.Errorf("error on reading IAM group, %s", err)
	}
	return
}

func resourceKsyunIamGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	iamGroupService := IamGroupService{meta.(*KsyunClient)}
	err = iamGroupService.DeleteIamGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting IAM group %q, %s", d.Id(), err)
	}
	return
}
