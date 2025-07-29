/*
Provides a DirectConnectGateWay resource.

Example Usage

```hcl
resource "ksyun_bws" "default" {
  line_id = "5fc2595f-1bfd-481b-bf64-2d08f116d800"
  charge_type = "PostPaidByPeak"
  band_width = 12
}
```

Import

BWS can be imported using the id, e.g.

```
$ terraform import ksyun_bws.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
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
				Description: "The id of direct connect. It's meaning is the physical port.",
			},
			"direct_connect_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the direct connect interface.",
			},

			// computed fields
			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the direct connect interface.",
			},
			"nat_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of vlan in direct connect.",
			},
			"direct_connect_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "bandwidth value, value range: [1, 15000].",
			},
			"direct_connect_interface_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The charge type of the BWS. Valid values: PostPaidByPeak, PostPaidByDay, PostPaidByTransfer, DailyPaidByTransfer.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the project.",
			},
			"associated_instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"cen_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"cen_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"remote_cidr_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "ID of the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"extra_cidr_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "ID of the project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"direct_connect_interface_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "ID of the project.",
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
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.CreateDirectConnectGateway(d, resourceKsyunDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectGateway %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRead(d, meta)
}

func resourceKsyunDirectConnectGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetDirectConnectGateway(d, resourceKsyunDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectGateway %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ModifyDirectConnectGateway(d, resourceKsyunDirectConnectGateway())
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectGateway %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRead(d, meta)
}

func resourceKsyunDirectConnectGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.RemoveDirectConnectGateway(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectGateway %q, %s", d.Id(), err)
	}
	return err
}
