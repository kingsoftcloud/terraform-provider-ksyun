/*
Provides slb listener mount backend server group resource.

~> **NOTE:** This resource is **deprecated**. Use `backend_server_group_mounted` of `ksyun_lb_listener` instead. See [ksyun_lb_listener](https://registry.terraform.io/providers/kingsoftcloud/ksyun/latest/docs/resources/lb_listener) for more details.

# Example Usage

```hcl

```

# Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_lb_listener_associate_backendgroup.default $listener_id:$backend_server_group_id
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunLbListenerAssociateBackendgroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLbListenerAssociateBackendgroupCreate,
		Read:   resourceKsyunLbListenerAssociateBackendgroupRead,
		Update: nil,
		Delete: resourceKsyunLbListenerAssociateBackendgroupDelete,
		Importer: &schema.ResourceImporter{
			State: importListenerAssociateBackendgroup,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of alb backend server group.",
			},

			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of slb listener.",
			},
		},
	}
}

func resourceKsyunLbListenerAssociateBackendgroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.MountOrUnmountBackendGroup(d, resourceKsyunLbListenerAssociateBackendgroup(), true)
	if err != nil {
		return fmt.Errorf("error on mounting backend group onto listener %q, %s", d.Id(), err)
	}

	id := AssembleIds(d.Get("listener_id").(string), d.Get("backend_server_group_id").(string))
	d.SetId(id)
	return resourceKsyunLbListenerAssociateBackendgroupRead(d, meta)
}

func resourceKsyunLbListenerAssociateBackendgroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadListenerBackendGroups(d)
	if err != nil {
		return fmt.Errorf("error on reading backend group association %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunLbListenerAssociateBackendgroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.MountOrUnmountBackendGroup(d, resourceKsyunLbListenerAssociateBackendgroup(), false)
	if err != nil {
		return fmt.Errorf("error on debongding backend group from listener %q, %s", d.Id(), err)
	}
	return err
}
