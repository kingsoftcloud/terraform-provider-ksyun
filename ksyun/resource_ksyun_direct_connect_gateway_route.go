/*
Provides a DirectConnectGatewayRoute resource.

Example Usage

```hcl

data "ksyun_direct_connects" "test" {
  ids        = []
  name_regex = ".*test.*"
}


resource "ksyun_direct_connect_interface" "test" {
  direct_connect_id  = data.ksyun_direct_connects.test.direct_connects[0].id
  route_type         = "STATIC"
  bgp_peer           = 59019
  bgp_client_token   = "dadasd"
  reliability_method = "bfd"
  enable_ipv6        = true
  bfd_config_id      = "29e0c675-2cca-4778-b331-884fca06de17"
  vlan_id            = 111

  direct_connect_interface_name = "tf_direct_connect_test_1"
}


resource "ksyun_direct_connect_gateway" "test" {
  direct_connect_gateway_name = "tf_direct_connect_gateway_test_1"
  vpc_id                      = "a38673ae-c9b7-4f8e-b727-b6feb648805b"
}

resource "ksyun_direct_connect_gateway_route" "test" {
  direct_connect_gateway_id = ksyun_direct_connect_gateway.test.id
  destination_cidr_block    = "192.136.0.0/24"
  next_hop_type             = "Vpc"
  depends_on                = [ksyun_dc_interface_associate.test]
}
```

Import

Route can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_gateway_route.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
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
				Description: "The destination CIDR block of the route. The CIDR block must be in the format of `x.x.x.x/x`.",
			},

			"next_hop_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the next hop. Valid values: `Vpc`, `DirectConnect`, `Cen`. Default is `Vpc`. If set to `DirectConnect`, the next hop instance ID must be provided.",
			},

			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the direct connect gateway.",
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				// ForceNew:    true,
				Computed:    true,
				Description: "Priority.",
			},
			"next_hop_instance": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
				Description: "The next hop instance ID.",
			},
			"enable_ip_v6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				// ForceNew:    true,
				Description: "whether to enable IPv6. Valid values: `true`, `false`. Default is `false`.",
			},
			"bgp_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("next_hop_type").(string) != "Vpc"
				},
				Description: "BGP Status.",
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

			"route_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Route Type.",
			},
		},
	}
}

func resourceKsyunDirectConnectGatewayRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.CreateDirectConnectGatewayRoute(d, resourceKsyunDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRouteRead(d, meta)
}

func resourceKsyunDirectConnectGatewayRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ReadAndSetDirectConnectGatewayRoute(d, resourceKsyunDirectConnectGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectGatewayRouteUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.PublishDirectConnectRoute(d)
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectGatewayRouteRead(d, meta)
}

func resourceKsyunDirectConnectGatewayRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.RemoveDirectConnectGatewayRoute(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectGatewayRoute %q, %s", d.Id(), err)
	}
	return err
}
