/*
Provides a namespace resource under kcrs repository instance.

Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest"
	instance_type = "basic"
}


# Create a namespace under the repository instance
resource "ksyun_kcrs_namespace" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	namespace = "tftest"
	public = false
}
```

Import

KcrsNamespace can be imported using `instance_id:namespace_name`, e.g.

```
$ terraform import ksyun_kcrs_namespace.foo ${instance_id}:${namespace_name}
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunKcrsNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKcrsNamespaceCreate,
		Read:   resourceKsyunKcrsNamespaceRead,
		Update: resourceKsyunKcrsNamespaceUpdate,
		Delete: resourceKsyunKcrsNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: importKcrsNamespace,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "Instance id of repository.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of namespace.",
			},

			"public": {
				Type:     schema.TypeBool,
				Required: true,

				// ForceNew:    true,
				Description: "Whether to be public this namespace.",
			},
		},
	}
}
func resourceKsyunKcrsNamespaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsNamespaceService := KcrsService{meta.(*KsyunClient)}
	err = kcrsNamespaceService.CreateKcrsNamespace(d, resourceKsyunKcrsNamespace())
	if err != nil {
		return fmt.Errorf("error on creating kcrs namespace %q, %s", d.Id(), err)
	}
	return resourceKsyunKcrsNamespaceRead(d, meta)
}

func resourceKsyunKcrsNamespaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsNamespaceService := KcrsService{meta.(*KsyunClient)}
	err = kcrsNamespaceService.ReadAndSetKcrsNamespace(d, resourceKsyunKcrsInstance())
	if err != nil {
		return fmt.Errorf("error on reading kcrs Namespace %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKcrsNamespaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*KsyunClient)
	conn := client.kcrsconn
	if d.HasChange("public") {
		req := make(map[string]interface{}, 3)
		req["InstanceId"] = d.Get("instance_id")
		req["Namespace"] = d.Get("namespace")
		req["Public"] = helper.StringBoolean(d.Get("public").(bool))

		action := "DescribeInstance"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifyNamespaceType(&req)
		if err != nil {
			return err
		}
	}

	return resourceKsyunKcrsNamespaceRead(d, meta)
}

func resourceKsyunKcrsNamespaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsNamespaceService := KcrsService{meta.(*KsyunClient)}
	err = kcrsNamespaceService.RemoveKcrsNamespace(d)
	if err != nil {
		return fmt.Errorf("error on deleting kcrs Namespace %q, %s", d.Id(), err)
	}
	return err

}
