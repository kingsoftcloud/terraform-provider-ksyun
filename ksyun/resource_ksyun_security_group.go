/*
Provides a Security Group resource.

# Example Usage

```hcl

	resource "ksyun_security_group" "default" {
	  vpc_id = "26231a41-4c6b-4a10-94ed-27088d5679df"
	  security_group_name="xuan-tf--s"
	}

```

# Import

Security Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_security_group.example xxxxxxxx-abc123456
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunSecurityGroup() *schema.Resource {
	entry := resourceKsyunSecurityGroupEntry().Schema
	for k, v := range entry {
		if k == "security_group_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
		}
	}
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupCreate,
		Update: resourceKsyunSecurityGroupUpdate,
		Read:   resourceKsyunSecurityGroupRead,
		Delete: resourceKsyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The Id of the vpc.",
			},

			"security_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the security group.",
			},

			"security_group_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      securityGroupEntryHash,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Description: "Network security group Entries. this parameter will be deprecated, use `ksyun_security_group_entry` instead.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the security group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of creation of security group.",
			},
		},
	}
}

func resourceKsyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSecurityGroup(d, resourceKsyunSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on creating security group  %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSecurityGroup(d, resourceKsyunSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on reading security group  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifySecurityGroup(d, resourceKsyunSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on updating security group  %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveSecurityGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting security group  %q, %s", d.Id(), err)
	}
	return err
}
