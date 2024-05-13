/*
Provides a Security Group Entry resource that can manage a list of diverse cidr_block.

# Example Usage

```hcl

	resource "ksyun_security_group_entry_lite" "default" {
	  security_group_id="7385c8ea-xxxx-xxxx-xxxx-517fc3726256"
	  cidr_block=["10.0.0.1/32", "10.111.222.1/32"]
	  direction="in"
	  protocol="ip"
	}

```

# Import

-> **NOTE:** This resource cannot be imported. if you need import security group entry, you are supposed to use `ksyun_security_group_entry`

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
)

func resourceKsyunSecurityGroupEntryLite() *schema.Resource {
	entry := resourceKsyunSecurityGroupEntry().Schema
	for k := range entry {
		if k == "security_group_entry_id" {
			delete(entry, k)
		}
	}
	entry["cidr_block"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validation.Any(
				validation.StringIsEmpty,
				validation.IsCIDR,
			),
		},
		DiffSuppressFunc: securityGroupEntryLiteDiffSuppress,

		Description: "The cidr block list of security group rule.",
	}
	entry["security_group_entry_id_list"] = &schema.Schema{
		Type: schema.TypeString,
		// Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "The security group entry id of this lite managed.",
	}

	return &schema.Resource{
		Create: resourceKsyunSecurityGroupEntryLiteCreate,
		Read:   resourceKsyunSecurityGroupEntryLiteRead,
		Update: resourceKsyunSecurityGroupEntryLiteUpdate,
		Delete: resourceKsyunSecurityGroupEntryLiteDelete,
		Importer: &schema.ResourceImporter{
			State: denyImport,
		},
		Schema: entry,
	}
}

func resourceKsyunSecurityGroupEntryLiteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	cidrBlocks, _ := helper.GetSchemaListWithString(d, "cidr_block")
	if IsDuplicationInSlice(cidrBlocks) {
		return fmt.Errorf("cidr block list contains duplication, please check the duplication")
	}

	entryLiteId := buildSecurityGroupEntryLiteId(d)
	d.SetId(entryLiteId)

	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSecurityGroupEntryLite(d, resourceKsyunSecurityGroupEntryLite())
	if err != nil {
		return fmt.Errorf("error on creating security group entry lite %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupEntryLiteRead(d, meta)
}

func resourceKsyunSecurityGroupEntryLiteRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSecurityGroupEntryLite(d, resourceKsyunSecurityGroupEntryLite())
	if err != nil {
		return fmt.Errorf("error on reading security group entry lite %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSecurityGroupEntryLiteUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}

	if d.HasChange("description") {
		err = vpcService.ModifySecurityGroupEntryLite(d, resourceKsyunSecurityGroupEntryLite())
		if err != nil {
			return fmt.Errorf("error on updating security group entry lite %q, %s", d.Id(), err)
		}
	}

	return resourceKsyunSecurityGroupEntryLiteRead(d, meta)
}

func resourceKsyunSecurityGroupEntryLiteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveSecurityGroupEntryLite(d)
	if err != nil {
		return fmt.Errorf("error on deleting security group entry lite %q, %s", d.Id(), err)
	}
	return err
}

func buildSecurityGroupEntryLiteId(d *schema.ResourceData) string {
	var varSuffix string
	// groupId-direction-protocol-port_range_from(icmp_type)-port_range_to(icmp_code)
	protocol := d.Get("protocol").(string)
	switch protocol {
	case "icmp":
		varSuffix = fmt.Sprintf("%d-%d", d.Get("icmp_type").(int), d.Get("icmp_code").(int))
	default:
		varSuffix = fmt.Sprintf("%d-%d", d.Get("port_range_from").(int), d.Get("port_range_to").(int))
	}
	return fmt.Sprintf("%s_%s-%s-%s", d.Get("security_group_id").(string), d.Get("direction").(string), protocol, varSuffix)
}
