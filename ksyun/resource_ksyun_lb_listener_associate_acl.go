/*
Associate a Load Balancer Listener resource with acl.

# Example Usage

```hcl

	resource "ksyun_lb_listener_associate_acl" "default" {
	  listener_id = "b330eae5-11a3-4e9e-bf7d-xxxxxxxxxxxx"
	  load_balancer_acl_id = "7e94fa82-05c7-496c-ae5e-xxxxxxxxxxxx"
	}

```

LB Listener assocaite acl resource can be imported using the `listener_id`+`load_balancer_acl_id`, e.g.

```
$ terraform import ksyun_lb_listener_associate_acl.default ${listener_id}:${load_balancer_acl_id}
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunListenerAssociateAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerAssociateAclCreate,
		Read:   resourceKsyunListenerAssociateAclRead,
		Delete: resourceKsyunListenerAssociateAclDelete,
		Importer: &schema.ResourceImporter{
			State: importLoadBalancerAclAssociate,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the listener.",
			},
			"load_balancer_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the load balancer acl.",
			},
			"lb_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of listener. Valid Value: `Alb` and `Slb`. Default: `Slb`.",
			},
		},
	}
}
func resourceKsyunListenerAssociateAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLoadBalancerAclAssociate(d, resourceKsyunListenerAssociateAcl())
	if err != nil {
		return fmt.Errorf("error on creating listener acl associate %q, %s", d.Id(), err)
	}
	_ = d.Set("lb_type", "Slb")
	return resourceKsyunListenerAssociateAclRead(d, meta)
}

func resourceKsyunListenerAssociateAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLoadBalancerAclAssociate(d, resourceKsyunListenerAssociateAcl())
	if err != nil {
		return fmt.Errorf("error on reading  listener acl associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunListenerAssociateAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLoadBalancerAclAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener acl associate %q, %s", d.Id(), err)
	}
	return err
}
