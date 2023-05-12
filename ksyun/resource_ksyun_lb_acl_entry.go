/*
Provides a Load Balancer acl entry resource to add content forwarding policies for Load Balancer backend resource.

# Example Usage

```hcl

	resource "ksyun_lb_acl_entry" "default" {
	  load_balancer_acl_id = "8e6d0871-da8a-481e-8bee-b3343e2a6166"
	  cidr_block = "192.168.11.2/32"
	  rule_number = 10
	  rule_action = "allow"
	  protocol = "ip"
	}

```

# Import

LB ACL entry can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_acl_entry.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunLoadBalancerAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclEntryCreate,
		Delete: resourceKsyunLoadBalancerAclEntryDelete,
		Update: resourceKsyunLoadBalancerAclEntryUpdate,
		Read:   resourceKsyunLoadBalancerAclEntryRead,
		Importer: &schema.ResourceImporter{
			State: importLoadBalancerAclEntry,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the load balancer acl.",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The information of the load balancer Acl's cidr block.",
			},
			"rule_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 32766),
				Description:  "The information of the load balancer Acl's rule priority. value range:[1-32766].",
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "allow",
				ValidateFunc: validation.StringInSlice([]string{
					"allow",
					"deny",
				}, false),
				Description: "The action of load balancer Acl rule. Valid Values:'allow', 'deny'. Default is 'allow'.",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ip",
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
				}, false),
				Description: "protocol.Valid Values:'ip'.",
			},
			"load_balancer_acl_entry_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the LB ACL entry.",
			},
		},
	}
}

func resourceKsyunLoadBalancerAclEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLoadBalancerAclEntry(d, resourceKsyunLoadBalancerAclEntry())
	if err != nil {
		return fmt.Errorf("error on reading lb acl entry %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunLoadBalancerAclEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLoadBalancerAclEntry(d, resourceKsyunLoadBalancerAclEntry())
	if err != nil {
		return fmt.Errorf("error on creating lb acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunLoadBalancerAclEntryRead(d, meta)
}

func resourceKsyunLoadBalancerAclEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyLoadBalancerAclEntry(d, resourceKsyunLoadBalancerAclEntry())
	if err != nil {
		return fmt.Errorf("error on updating lb acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunLoadBalancerAclEntryRead(d, meta)
}

func resourceKsyunLoadBalancerAclEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLoadBalancerAclEntry(d)
	if err != nil {
		return fmt.Errorf("error on deleting lb acl entry  %q, %s", d.Id(), err)
	}
	return err
}
