/*
Provides a Security Group Entry resource.

# Example Usage

```hcl

	resource "ksyun_security_group_entry_set" "default" {
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
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunSecurityGroupEntrySet() *schema.Resource {
	entry := resourceKsyunSecurityGroupEntry().Schema
	for k, v := range entry {
		if k == "security_group_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
		}
	}
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupEntrySetCreate,
		Read:   resourceKsyunSecurityGroupEntrySetRead,
		Update: resourceKsyunSecurityGroupEntrySetUpdate,
		Delete: resourceKsyunSecurityGroupEntrySetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the security group.",
			},
			"security_group_entries": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      securityGroupEntryHash,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Description: "Network security group Entries.",
			},
		},
	}
}

/*
###

this development of resource is suspended, because of the remote have a default entry that's direction is out, so that
provider cannot manage that default entry when creating.

###

*/

func resourceKsyunSecurityGroupEntrySetCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSecurityGroupEntrySet(d, resourceKsyunSecurityGroupEntrySet())
	if err != nil {
		return fmt.Errorf("error on creating security group set %q, %s", d.Id(), err)
	}
	d.SetId(d.Get("security_group_id").(string))
	return resourceKsyunSecurityGroupEntrySetRead(d, meta)
}

func resourceKsyunSecurityGroupEntrySetRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSecurityGroup(d, resourceKsyunSecurityGroupEntrySet())
	if err != nil {
		return fmt.Errorf("error on reading security group set %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSecurityGroupEntrySetUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifySecurityGroupSet(d, resourceKsyunSecurityGroupEntrySet())
	if err != nil {
		return fmt.Errorf("error on updating security group set %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupEntrySetDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}

	apiProcess := NewApiProcess(context.Background(), d, vpcService.client, true)

	gId := d.Get("security_group_id").(string)
	for _, entries := range d.Get("security_group_entries").(*schema.Set).List() {
		entry := entries.(map[string]interface{})
		entryId := entry["security_group_entry_id"].(string)
		call, err := vpcService.RemoveSecurityGroupEntryCommonCall(gId, entryId)
		if err != nil {
			return fmt.Errorf("error on deleting security group entry set %q, %s", entryId, err)
		}
		apiProcess.PutCalls(call)
	}
	err = apiProcess.Run()
	if err != nil {
		return fmt.Errorf("error on deleting security group set %q, %s", d.Id(), err)
	}
	return err
}
