/*
Provides a route resource under VPC resource.

# Example Usage

```hcl

	resource "ksyun_vpc" "example" {
	  vpc_name   = "tf-example-vpc-01"
	  cidr_block = "10.0.0.0/16"
	}

	resource "ksyun_route" "example" {
	  destination_cidr_block = "10.0.0.0/16"
	  route_type = "InternetGateway"
	  vpc_id = "${ksyun_vpc.example.id}"
	}

```

# Import

route can be imported using the `id`, e.g.

```
$ terraform import ksyun_route.example xxxx-xxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRouteCreate,
		Read:   resourceKsyunRouteRead,
		Delete: resourceKsyunRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The id of the vpc.",
			},

			"destination_cidr_block": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "The CIDR block assigned to the route.",
			},

			"route_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"InternetGateway",
					"Tunnel",
					"Host",
					"Peering",
					"DirectConnect",
					"Vpn",
				}, false),
				Description: "The type of route.Valid Values:'InternetGateway', 'Tunnel', 'Host', 'Peering', 'DirectConnect', 'Vpn'.",
			},
			"tunnel_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The id of the tunnel If route_type is Tunnel, This Field is Required.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The id of the VM, If route_type is Host, This Field is Required.",
			},
			"vpc_peering_connection_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The id of the Peering, If route_type is Peering, This Field is Required.",
			},
			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The id of the DirectConnectGateway, If route_type is DirectConnect, This Field is Required.",
			},
			"vpn_tunnel_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The id of the Vpn, If route_type is Vpn, This Field is Required.",
			},
			"next_hop_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of next hop.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the gateway.",
						},

						"gateway_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the gateway.",
						},
					},
				},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of creation of the route.",
			},
		},
	}
}

func resourceKsyunRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateRoute(d, resourceKsyunRoute())
	if err != nil {
		return fmt.Errorf("error on creating route %q, %s", d.Id(), err)
	}
	return resourceKsyunRouteRead(d, meta)
}

func resourceKsyunRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetRoute(d, resourceKsyunRoute())
	if err != nil {
		return fmt.Errorf("error on reading route %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveRoute(d)
	if err != nil {
		return fmt.Errorf("error on deleting route %q, %s", d.Id(), err)
	}
	return err
}
