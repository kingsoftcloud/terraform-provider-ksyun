/*
Provides a Vpn Gateway Route resource under VPC resource.
**Notes:** `ksyun_vpn_gateway_route` only valid when Vpn 2.0

# Example Usage

```hcl

resource "ksyun_vpn_gateway_route" "default1" {
  vpn_gateway_id = "450a71b0-ea20-****-*****"
  next_hop_type = "vpc"
  destination_cidr_block = "10.7.255.0/30"
}

```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnGatewayRoute() *schema.Resource {
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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
				),
				ForceNew:    true,
				Description: "The destination cidr block.",
			},
			"next_hop_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vpn_tunnel",
					"vpc",
				}, false),
				Description: "The type of next hop. Valid Values: `vpn_tunnel`, `vpc`.",
			},

			"next_hop_instance_id": {
				Type: schema.TypeString,
				// Required:    true,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("next_hop_type") == "vpn_tunnel" {
						return false
					}
					return true
				},
				Description: "The instance id of next hop, which must be set when `next_hop_type` is `vpn_tunnel.",
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
	err = vpnService.CreateVpnGatewayRoute(d, resourceKsyunVpnGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on creating vpn gateway route %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnGatewayRouteRead(d, meta)
}

func resourceKsyunVpnGatewayRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpnService := NewVpnSrv(meta.(*KsyunClient))
	err = vpnService.ReadAndSetVpnGatewayRoute(d, resourceKsyunVpnGatewayRoute())
	if err != nil {
		return fmt.Errorf("error on reading vpn gateway route %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnGatewayRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := NewVpnSrv(meta.(*KsyunClient))
	err = vpcService.RemoveVpnGatewayRoute(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn gateway route %q, %s", d.Id(), err)
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
