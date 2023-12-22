/*
Provides an KcrsNamespace resource.

Example Usage

```hcl
# Create a KcrsNamespace
resource "ksyun_KcrsNamespace" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  KcrsNamespace_name = "ksc_KcrsService"
  bill_type = 1
  service_id = "KcrsNamespace_30G"
  project_id="0"
}
```

Import

KcrsNamespace can be imported using the id, e.g.

```
$ terraform import ksyun_KcrsNamespace.default KcrsService67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
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
				Description: "The id of the project.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the ID of the KcrsNamespace.",
			},

			"public": {
				Type:     schema.TypeBool,
				Required: true,

				// ForceNew:    true,
				Description: "the max ip count that can bind to the KcrsNamespace,value range: [10, 100].",
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
