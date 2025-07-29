/*
Provides a DirectConnectGatewayRoute resource.

Example Usage

```hcl
resource "ksyun_bws" "default" {
  line_id = "5fc2595f-1bfd-481b-bf64-2d08f116d800"
  charge_type = "PostPaidByPeak"
  band_width = 12
}
```

Import

BwS can be imported using the id, e.g.

```
$ terraform import ksyun_bws.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunDirectConnectGatewayRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDirectConnectGatewayRouteCreate,
		Read:   resourceKsyunDirectConnectGatewayRouteRead,
		Update: resourceKsyunDirectConnectGatewayRouteUpdate,
		Delete: resourceKsyunDirectConnectGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"destination_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of direct connect. It's meaning is the physical port.",
			},

			"next_hop_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the direct connect interface.",
			},

			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the direct connect interface.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The id of vlan in direct connect.",
			},
			"next_hop_instance": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "bandwidth value, value range: [1, 15000].",
			},
			"enable_ip_v6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "whether to enable IPv6. Valid values: `true`, `false`. Default is `false`.",
			},

			// common fields
			"direct_connect_gateway_route_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the direct connect interface.",
			},
			"next_hop_instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the next hop instance.",
			},
			"as_path": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "AS Path of the route.",
			},
			"direct_connect_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Direct Connect ID.",
			},
			"bgp_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BGP Status",
			},
			"route_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Route Type",
			},
		},
	}
}

func resourceKsyunDirectConnectGatewayRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.CreateDirectConnectGatewayRoute(d, resourceKsyunDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRouteRead(d, meta)
}

func resourceKsyunDirectConnectGatewayRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetDirectConnectGatewayRoute(d, resourceKsyunDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectGatewayRouteUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ModifyDirectConnectGatewayRoute(d, resourceKsyunDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRouteRead(d, meta)
}

func resourceKsyunDirectConnectGatewayRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.RemoveDirectConnectGatewayRoute(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return err
}
