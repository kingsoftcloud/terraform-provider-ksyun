/*
Provides a Iam Project resource.

# Example Usage

```hcl

resource "ksyun_iam_project" "project" {
  project_name = "ProjectNameTest"
  project_desc = "ProjectDesc"
}

```

# Import

IAM Project can be imported using the `project_name`, e.g.

```
$ terraform import ksyun_iam_project.project project_name
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunIamProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunIamProjectCreate,
		Read:   resourceKsyunIamProjectRead,
		Update: resourceKsyunIamProjectUpdate,
		Delete: resourceKsyunIamProjectDelete,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IAM ProjectName.",
			},
			"project_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IAM ProjectDesc.",
			},
		},
	}
}

func resourceKsyunIamProjectCreate(d *schema.ResourceData, meta interface{}) (err error) {
	iamProjectService := IamProjectService{meta.(*KsyunClient)}
	err = iamProjectService.CreateIamProject(d, resourceKsyunIamProject())
	if err != nil {
		return fmt.Errorf("error on creating IAM project %q, %s", d.Id(), err)
	}
	return resourceKsyunIamProjectRead(d, meta)
}

func resourceKsyunIamProjectUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunIamProjectRead(d *schema.ResourceData, meta interface{}) (err error) {
	iamProjectService := IamProjectService{meta.(*KsyunClient)}
	err = iamProjectService.ReadAndSetIamProject(d, resourceKsyunIamProject())
	if err != nil {
		return fmt.Errorf("error on reading IAM project, %s", err)
	}
	return
}

func resourceKsyunIamProjectDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return
}
