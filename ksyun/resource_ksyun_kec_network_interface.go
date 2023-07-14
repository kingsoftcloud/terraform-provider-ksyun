/*
Provides a KEC network interface resource.

# Example Usage

```hcl

	resource "ksyun_kec_network_interface" "default" {
	 subnet_id = "81530211-2785-47a8-b2a0-ae13120fa97d"
	 security_group_ids = ["7e2f45b5-e79d-4612-a7fc-fe74a50b639a","35ac2642-1958-4ed7-b02c-dc86f27bc9d9"]
	 network_interface_name = "Ksc_NetworkInterface"
	}

```

# Import

Instance can be imported using the id, e.g.

```
$ terraform import ksyun_kec_network_interface.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunKecNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceKecNetworkInterfaceCreate,
		Update: resourceKecNetworkInterfaceUpdate,
		Read:   resourceKecNetworkInterfaceRead,
		Delete: resourceKecNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: kecNetworkInterfaceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"network_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the network interface.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the subnet which the network interface belongs to.",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Private IP.",
			},
			"security_group_ids": {
				Type:        schema.TypeSet,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Set:         schema.HashString,
				Description: "A list of security group IDs.",
			},

			"secondary_private_ips": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"secondary_private_ip_address_count"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Secondary Private IP.",
						},
					},
				},
				Description: "Assign secondary private ips to the network interface. <br> Notes: `secondary_private_ips` conflict with `secondary_private_ip_address_count`.",
			},

			"secondary_private_ip_address_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"secondary_private_ips"},
				ValidateFunc:  validation.IntBetween(1, 99),
				Description:   "The count of secondary private id address automatically assigned. <br> Notes:  `secondary_private_ip_address_count` conflict with `secondary_private_ips`.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance id to bind with the network interface.",
			},
		},
	}
}
func resourceKecNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.createNetworkInterface(d, resourceKsyunKecNetworkInterface())
	if err != nil {
		return fmt.Errorf("error on creating network interface %q, %s", d.Id(), err)
	}
	return resourceKecNetworkInterfaceRead(d, meta)
}

func resourceKecNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.modifyNetworkInterface(d, resourceKsyunKecNetworkInterface())
	if err != nil {
		return fmt.Errorf("error on updating network interface %q, %s", d.Id(), err)
	}
	return resourceKecNetworkInterfaceRead(d, meta)
}

func resourceKecNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.readAndSetNetworkInterface(d, resourceKsyunKecNetworkInterface())
	if err != nil {
		return fmt.Errorf("error on reading network interface %q, %s", d.Id(), err)
	}
	return err
}

func resourceKecNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.RemoveNetworkInterface(d)
}
