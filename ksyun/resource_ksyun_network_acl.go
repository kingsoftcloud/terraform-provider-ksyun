/*
Provides a Network ACL resource under VPC resource.

# Example Usage

```hcl

	resource "ksyun_network_acl" "default" {
	  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
	  network_acl_name = "ceshi"
	}

```

# Import

Network ACL can be imported using the `id`, e.g.

```
$ terraform import ksyun_network_acl.default fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunNetworkAcl() *schema.Resource {
	entry := resourceKsyunNetworkAclEntry().Schema
	for k, v := range entry {
		if k == "network_acl_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
		}
	}
	return &schema.Resource{
		Create: resourceKsyunNetworkAclCreate,
		Read:   resourceKsyunNetworkAclRead,
		Delete: resourceKsyunNetworkAclDelete,
		Update: resourceKsyunNetworkAclUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: networkAclEntryCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the vpc.",
			},
			"network_acl_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the network ACL.",
			},
			"network_acl_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Set:         networkAclEntryHash,
				Description: "Network ACL Entries. this parameter will be deprecated, use `ksyun_network_acl_entry` instead.",
			},
			"network_acl_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the network ACL.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of network acl.",
			},
		},
	}
}

func resourceKsyunNetworkAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNetworkAcl(d, resourceKsyunNetworkAcl())
	if err != nil {
		return fmt.Errorf("error on creating network acl %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclRead(d, meta)
}

func resourceKsyunNetworkAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNetworkAcl(d, resourceKsyunNetworkAcl())
	if err != nil {
		return fmt.Errorf("error on reading network acl  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNetworkAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyNetworkAcl(d, resourceKsyunNetworkAcl())
	if err != nil {
		return fmt.Errorf("error on updating network acl %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclRead(d, meta)
}

func resourceKsyunNetworkAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNetworkAcl(d)
	if err != nil {
		return fmt.Errorf("error on deleting network acl %q, %s", d.Id(), err)
	}
	return err
}
