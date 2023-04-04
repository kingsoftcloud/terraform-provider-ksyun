/*
Provides a Network ACL Associate resource.

# Example Usage

```hcl

	resource "ksyun_network_acl_associate" "test" {
	  network_acl_id = "679b6a88-67dd-4e17-a80a-985d9673050e"
	  subnet_id = "84cc79f3-dc88-4f00-a66a-c7e8d68ec615"
	}

```

# Import

Network ACL Associate can be imported using the `network_acl_id:subnet_id`, e.g.

```
$ terraform import ksyun_network_acl_associate.default $network_acl_id:$subnet_id
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunNetworkAclAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNetworkAclAssociateCreate,
		Read:   resourceKsyunNetworkAclAssociateRead,
		Delete: resourceKsyunNetworkAclAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importNetworkAclAssociate,
		},
		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the network acl.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Subnet.",
			},
		},
	}
}

func resourceKsyunNetworkAclAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNetworkAclAssociate(d, resourceKsyunNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on creating network acl associate  %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclAssociateRead(d, meta)
}

func resourceKsyunNetworkAclAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNetworkAclAssociate(d, resourceKsyunNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on reading network acl associate  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNetworkAclAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNetworkAclAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting network acl associate  %q, %s", d.Id(), err)
	}
	return err
}
