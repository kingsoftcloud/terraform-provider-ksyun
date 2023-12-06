/*
Provides a Vpn Gateway resource under VPC resource.

# Example Usage

```hcl

# create vpn gateway with vpn 1.0 version
resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name   = "ksyun_vpn_gw_tf1"
  band_width = 10
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  charge_type = "Daily"
  # vpn_gateway_version = "1.0"
}

# create vpn gateway with vpn 2.0 version
resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name   = "ksyun_vpn_gw_tf1"
  band_width = 10
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  charge_type = "Daily"
  vpn_gateway_version = "2.0"
}
```

# Import

Vpn Gateway can be imported using the `id`, e.g.

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

func resourceKsyunVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpnGatewayCreate,
		Update: resourceKsyunVpnGatewayUpdate,
		Read:   resourceKsyunVpnGatewayRead,
		Delete: resourceKsyunVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the vpn gateway.",
			},

			"band_width": {
				Type: schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{
					5,
					10,
					20,
					50,
					100,
					200,
				}),
				Required:    true,
				Description: "The bandWidth of the vpn gateway.Valid Values:5,10,20,50,100,200.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the vpc.",
			},
			"vpn_gateway_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"1.0",
					"2.0",
				}, false),
				Default:     "1.0",
				Description: "the version of vpn gateway. Default `1.0`.",
			},

			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Monthly",
					"Daily",
				}, false),
				Description: "The charge type of the vpn gateway.Valid Values:'Monthly','Daily'.",
			},

			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(0, 36),
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
				Description:      "The purchase time of the vpn gateway.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project id  of the vpn gateway.Default is 0.",
			},

			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the vpn gateway.",
			},
		},
	}
}

func resourceKsyunVpnGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpnGateway(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on creating vpn gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnGatewayRead(d, meta)
}

func resourceKsyunVpnGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpnGateway(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on reading vpn gateway  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpnGateway(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on updating vpn gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnGatewayRead(d, meta)
}

func resourceKsyunVpnGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpnGateway(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn gateway  %q, %s", d.Id(), err)
	}
	return err
}
