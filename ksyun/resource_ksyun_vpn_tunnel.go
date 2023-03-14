/*
Provides a Vpn Tunnel resource.

# Example Usage

```hcl

	resource "ksyun_vpn_tunnel" "default" {
	  vpn_tunnel_name   = "ksyun_vpn_tunnel_tf_1"
	  type = "Ipsec"
	  vpn_gateway_id = "9b3d361e-f65b-464b-947a-fafb5cfb10d2"
	  customer_gateway_id = "7f5a5c91-4814-41bf-b9d6-d9d811f4df0f"
	  ike_dh_group = 2
	  pre_shared_key = "123456789abcd"
	}

```

# Import

Vpn Tunnel can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_tunnel.default $id
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpnTunnelCreate,
		Update: resourceKsyunVpnTunnelUpdate,
		Read:   resourceKsyunVpnTunnelRead,
		Delete: resourceKsyunVpnTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_tunnel_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the vpn tunnel.",
			},

			"type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"GreOverIpsec",
					"Ipsec",
				}, false),
				Required:    true,
				ForceNew:    true,
				Description: "The bandWidth of the vpn tunnel.Valid Values:'GreOverIpsec','Ipsec'.",
			},

			"vpn_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The vpn_gre_ip of the vpn tunnel.If type is GreOverIpsec, Required.",
			},

			"ha_vpn_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The ha_vpn_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required.",
			},

			"customer_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The customer_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required.",
			},

			"ha_customer_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The ha_customer_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required.",
			},

			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vpn_gateway_id of the vpn tunnel.",
			},

			"customer_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The customer_gateway_id of the vpn tunnel.",
			},

			"pre_shared_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The pre_shared_key of the vpn tunnel.",
			},

			"ike_authen_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"md5",
					"sha",
				}, false),
				Computed:    true,
				Description: "The ike_authen_algorithm of the vpn tunnel.Valid Values:'md5','sha'.",
			},

			"ike_dh_group": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					1,
					2,
					5,
				}),
				Computed:    true,
				Description: "The ike_dh_group of the vpn tunnel.Valid Values:1,2,5.",
			},

			"ike_encry_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"3des",
					"aes",
					"des",
				}, false),
				Computed:    true,
				Description: "The ike_encry_algorithm of the vpn tunnel.Valid Values:'3des','aes','des'.",
			},

			"ipsec_encry_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"esp-3des",
					"esp-aes",
					"esp-des",
					"esp-null",
					"esp-seal",
				}, false),
				Computed:    true,
				Description: "The ipsec_encry_algorithm of the vpn tunnel.Valid Values:'esp-3des','esp-aes','esp-des','esp-null','esp-seal'.",
			},

			"ipsec_authen_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"esp-md5-hmac",
					"esp-sha-hmac",
				}, false),
				Computed:    true,
				Description: "The ipsec_authen_algorithm of the vpn tunnel.Valid Values:'esp-md5-hmac','esp-sha-hmac'.",
			},

			"ipsec_lifetime_traffic": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(2560, 4608000),
				Computed:     true,
				Description:  "The ipsec_lifetime_traffic of the vpn tunnel.",
			},

			"ipsec_lifetime_second": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(120, 2592000),
				Computed:     true,
				Description:  "The ipsec_lifetime_second of the vpn tunnel.",
			},
		},
	}
}

func resourceKsyunVpnTunnelCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpnTunnel(d, resourceKsyunVpnTunnel())
	if err != nil {
		return fmt.Errorf("error on creating vpn tunnel  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnTunnelRead(d, meta)
}

func resourceKsyunVpnTunnelRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpnTunnel(d, resourceKsyunVpnTunnel())
	if err != nil {
		return fmt.Errorf("error on reading vpn tunnel  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnTunnelUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpnTunnel(d, resourceKsyunVpnTunnel())
	if err != nil {
		return fmt.Errorf("error on updating vpn tunnel  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnTunnelRead(d, meta)
}

func resourceKsyunVpnTunnelDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpnTunnel(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn tunnel  %q, %s", d.Id(), err)
	}
	return err
}
