/*
Provides an EIP Association resource for associating Elastic IP to UHost Instance, Load Balancer, etc.

Example Usage

```hcl
resource "ksyun_eip_associate" "slb" {
  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
  instance_type="Slb"
  instance_id="7fae85e4-ab1a-415c-aef9-03a402c79d97"
}
resource "ksyun_eip_associate" "server" {
  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
  instance_type="Ipfwd"
  instance_id="566567677-6766-4743-afb7-7c7081214092"
  network_interface_id="87945980-59659-04548-759045803"
}
```

Import

EIP Association can be imported using the id, e.g.

```
$ terraform import ksyun_eip_associate.default ${allocation_id}:${instance_id}:${network_interface_id}
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipAssociationCreate,
		Read:   resourceKsyunEipAssociationRead,
		Delete: resourceKsyunEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: importAddressAssociate,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of EIP.",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Ipfwd",
					"Slb",
				}, false),
				Description: "The type of the instance.Valid Values:'Ipfwd', 'Slb'.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the instance.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The id of the network interface.",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP version of the EIP.",
			},
			"internet_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "InternetGateway ID.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the project.",
			},
			"line_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the line.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The band width of the public address.",
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the EIP.",
			},

			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Elastic IP address.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the EIP.",
			},
			"band_width_share_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the BWS which the EIP associated.",
			},
			"is_band_width_share": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BWS EIP.",
			},
		},
	}
}
func resourceKsyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.CreateAddressAssociate(d, resourceKsyunEipAssociation())
	if err != nil {
		return fmt.Errorf("error on creating address association %q, %s", d.Id(), err)
	}
	return resourceKsyunEipAssociationRead(d, meta)
}

func resourceKsyunEipAssociationRead(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.ReadAndSetAddressAssociate(d, resourceKsyunEipAssociation())
	if err != nil {
		return fmt.Errorf("error on reading address association %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.RemoveAddressAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting address association %q, %s", d.Id(), err)
	}
	return err
}
