/*
Provides a Network ACL Entry resource under Network ACL resource.

# Example Usage

```hcl

	resource "ksyun_network_acl_entry" "test" {
	  description = "测试1"
	  cidr_block = "10.0.16.0/24"
	  rule_number = 16
	  direction = "in"
	  rule_action = "deny"
	  protocol = "ip"
	  network_acl_id = "679b6a88-67dd-4e17-a80a-985d9673050e"
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunNetworkAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNetworkAclEntryCreate,
		Read:   resourceKsyunNetworkAclEntryRead,
		Delete: resourceKsyunNetworkAclEntryDelete,
		Update: resourceKsyunNetworkAclEntryUpdate,
		Importer: &schema.ResourceImporter{
			State: importNetworkAclEntry,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the network acl entry.",
			},
			"network_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the network acl.",
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.StringIsEmpty,
					validation.IsCIDR,
				),
				Description: "The cidr_block of the network acl entry.",
			},
			"rule_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 32766),
				Description:  "The rule_number of the network acl entry. value range:[1,32766].",
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in",
					"out",
				}, false),
				ForceNew:    true,
				Description: "The direction of the network acl entry. Valid Values: 'in','out'.",
			},
			"rule_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow",
					"deny",
				}, false),
				ForceNew:    true,
				Description: "The rule_action of the network acl entry.Valid Values: 'allow','deny'.",
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"tcp",
					"udp",
					"icmp",
				}, false),
				ForceNew:    true,
				Description: "The protocol of the network acl entry.Valid Values: 'ip','icmp','tcp','udp'.",
			},
			"icmp_type": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
				Description:      "The icmp_type of the network acl entry.If protocol is icmp, Required.",
			},
			"icmp_code": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
				Description:      "The icmp_code of the network acl entry.If protocol is icmp, Required.",
			},
			"port_range_from": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
				Description:      "The port_range_from of the network acl entry.If protocol is tcp or udp,Required.",
			},
			"port_range_to": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
				Description:      "The port_range_to of the network acl entry.If protocol is tcp or udp,Required.",
			},
			"network_acl_entry_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the network acl entry.",
			},
		},
	}
}

func resourceKsyunNetworkAclEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNetworkAclEntry(d, resourceKsyunNetworkAclEntry())
	if err != nil {
		return fmt.Errorf("error on creating network acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclEntryRead(d, meta)
}

func resourceKsyunNetworkAclEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNetworkAclEntry(d, resourceKsyunNetworkAclEntry())
	if err != nil {
		return fmt.Errorf("error on reading network acl entry  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNetworkAclEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyNetworkAclEntry(d, resourceKsyunNetworkAclEntry())
	if err != nil {
		return fmt.Errorf("error on updating network acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclEntryRead(d, meta)
}

func resourceKsyunNetworkAclEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNetworkAclEntry(d)
	if err != nil {
		return fmt.Errorf("error on deleting network acl entry %q, %s", d.Id(), err)
	}
	return err
}
