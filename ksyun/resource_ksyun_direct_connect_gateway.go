/*
Provides a DirectConnectGateWay resource.

Example Usage

```hcl
resource "ksyun_direct_connect_gateway" "test" {
  direct_connect_gateway_name = "tf_direct_connect_gateway_test_1"
  vpc_id                      = "a38673ae-c9b7-4f8e-b727-b6feb648xxxx"
}
```

Import

ksyun_direct_connect_gateway can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_gateway.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunDirectConnectGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDirectConnectGatewayCreate,
		Read:   resourceKsyunDirectConnectGatewayRead,
		Update: resourceKsyunDirectConnectGatewayUpdate,
		Delete: resourceKsyunDirectConnectGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vpc Id.",
			},
			"direct_connect_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the direct connect gateway.",
			},

			// computed fields
			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the direct connect gateway.",
			},
			// "direct_connect_interface_id": {
			// 	Type:        schema.TypeString,
			// 	Computed:    true,
			// 	Description: "The gateway attached to the direct connect interface.",
			// },
			"nat_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the NAT gateway associated with the direct connect gateway.",
			},
			"direct_connect_interface_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the direct connect interface associated with the direct connect gateway.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Band width.",
			},
			"associated_instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of associated instance.",
			},
			"cen_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the CEN account associated with the direct connect gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status.",
			},
			"cen_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Cen.",
			},
			"remote_cidr_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The set of remote cidr.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// "extra_cidr_set": {
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "ID of the project.",
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeString,
			// 	},
			// },
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version.",
			},
			"direct_connect_interface_info_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The set of direct connect associated interface info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"direct_connect_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the direct connect interface.",
						},
					},
				},
			},
		},
	}
}

func resourceKsyunDirectConnectGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.CreateDirectConnectGateway(d, resourceKsyunDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectGateway %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRead(d, meta)
}

func resourceKsyunDirectConnectGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ReadAndSetDirectConnectGateway(d, resourceKsyunDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectGateway %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ModifyDirectConnectGateway(d, resourceKsyunDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectGateway %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRead(d, meta)
}

func resourceKsyunDirectConnectGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.RemoveDirectConnectGateway(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectGateway %q, %s", d.Id(), err)
	}
	return err
}
