/*
Provides a Security Group Entry resource.

# Example Usage

```hcl

	resource "ksyun_security_group_entry" "default" {
	  security_group_id="7385c8ea-79f7-4e9c-b99f-517fc3726256"
	  cidr_block="10.0.0.1/32"
	  direction="in"
	  protocol="ip"
	}

```

# Import

Security Group Entry can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group_entry.example xxxxxxxx-abc123456
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunSecurityGroupEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupEntryCreate,
		Read:   resourceKsyunSecurityGroupEntryRead,
		Update: resourceKsyunSecurityGroupEntryUpdate,
		Delete: resourceKsyunSecurityGroupEntryDelete,
		Importer: &schema.ResourceImporter{
			State: importSecurityGroupEntry,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the entry.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the security group.",
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.StringIsEmpty,
					validation.IsCIDR,
				),
				Description: "The cidr block of security group rule.",
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in",
					"out",
				}, false),
				Description: "The direction of the entry, valid values:'in', 'out'.",
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"tcp",
					"udp",
					"icmp",
				}, false),
				Description: "The protocol of the entry, valid values: 'ip', 'tcp', 'udp', 'icmp'.",
			},
			"icmp_type": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
				Description:      "ICMP type.The required if protocol type is 'icmp'.",
			},
			"icmp_code": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
				Description:      "ICMP code.The required if protocol type is 'icmp'.",
			},
			"port_range_from": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
				Description:      "Port rule start port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.",
			},
			"port_range_to": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
				Description:      "Port rule end port for TCP or UDP protocol.The required if protocol type is 'tcp' or 'udp'.",
			},
			"security_group_entry_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the entry.",
			},
		},
	}
}

func resourceKsyunSecurityGroupEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSecurityGroupEntry(d, resourceKsyunSecurityGroupEntry())
	if err != nil {
		return fmt.Errorf("error on creating security group entry %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupEntryRead(d, meta)
}
func resourceKsyunSecurityGroupEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSecurityGroupEntry(d, resourceKsyunSecurityGroupEntry())
	if err != nil {
		return fmt.Errorf("error on reading security group entry %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSecurityGroupEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifySecurityGroupEntry(d, resourceKsyunSecurityGroupEntry())
	if err != nil {
		return fmt.Errorf("error on updating security group entry %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupEntryRead(d, meta)
}

func resourceKsyunSecurityGroupEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveSecurityGroupEntry(d)
	if err != nil {
		return fmt.Errorf("error on deleting security group entry %q, %s", d.Id(), err)
	}
	return err
}
