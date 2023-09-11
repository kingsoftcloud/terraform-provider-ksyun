/*
Provides a Vpn Gateway Route resource under VPC resource.

# Example Usage

```hcl

```

# Import

Vpn Gateway Route can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_gateway.default $id
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnGatewayRouteRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpnGatewayRouteCreate,
		Read:   resourceKsyunVpnGatewayRouteRead,
		Delete: resourceKsyunVpnGatewayRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the vpn gateway.",
			},

			"destination_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the vpc.",
			},
			"next_hop_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vpn_tunnel",
					"vpc",
				}, false),
				Description: "the type of next hop.",
			},

			"next_hop_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The charge type of the vpn gateway.Valid Values:'Monthly','Daily'.",
			},
		},
	}
}

func resourceKsyunVpnGatewayRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	if err := checkNestNecessaryParams(d); err != nil {
		return err
	}

	vpcService := VpcService{meta.(*KsyunClient)}

	vpnGateway, readGWErr := vpcService.ReadVpnGateway(d, d.Get("vpn_gateway_id").(string))
	if readGWErr != nil {
		return readGWErr
	}
	if v, ok := vpnGateway["VpnGatewayVersion"]; !ok {
		return fmt.Errorf("an error caused by checking vpn version of gateway")
	} else {
		if v != "2.0" {
			return fmt.Errorf("the use of ksyun_vpn_gateway_route is only supported on vpn2.0")
		}
	}

	vpnService := NewVpnSrv(meta.(*KsyunClient))
	err = vpnService.CreateVpnGatewayRoute(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on creating vpn gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnGatewayRouteRead(d, meta)
}

func resourceKsyunVpnGatewayRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpnService := NewVpnSrv(meta.(*KsyunClient))
	err = vpnService.ReadAndSetVpnGatewayRoute(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on reading vpn gateway  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnGatewayRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := NewVpnSrv(meta.(*KsyunClient))
	err = vpcService.RemoveVpnGatewayRoute(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn gateway  %q, %s", d.Id(), err)
	}
	return err
}

func checkNestNecessaryParams(d *schema.ResourceData) error {
	switch d.Get("next_hop_type") {
	case "vpn_tunnel":
		if _, ok := d.GetOk("next_hop_instance_id"); !ok {
			return fmt.Errorf("next_hop_instance_id cannot be blank, when next_hop_type is vpn_tunnel")
		}
	}
	return nil
}
