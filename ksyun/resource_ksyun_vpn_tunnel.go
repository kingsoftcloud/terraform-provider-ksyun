/*
Provides a Vpn Tunnel resource.

# Example Usage

```hcl

# create Vpn Tunnel with Vpn 1.0
resource "ksyun_vpn_tunnel" "tunnel-vpn1" {
  vpn_tunnel_name   = "tf_vpn_tunnel_vpn1"
  type = "Ipsec"
  vpn_gateway_id = "9b3d361e-f65b-464b-947a-fafb5cfb10d2"
  customer_gateway_id = "7f5a5c91-4814-41bf-b9d6-d9d811f4df0f"
  ike_dh_group = 2
  pre_shared_key = "123456789abcd"
}

# create Vpn Tunnel with Vpn 2.0
resource "ksyun_vpn_tunnel" "tunnel-vpn2" {
  vpn_gateway_version = "2.0" # choose vpn gateway version
  vpn_tunnel_name   = "tf_vpn_tunnel_vpn2"
  type = "Ipsec"
  ike_version = "v1"
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
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnTunnel() (r *schema.Resource) {
	r = &schema.Resource{
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
					"RouteIpsec",
				}, false),
				Required:    true,
				ForceNew:    true,
				Description: "The bandWidth of the vpn tunnel. Valid Values: VPN-v1: 'GreOverIpsec' or 'Ipsec'; VPN-v2: `RouteIpsec` or `Ipsec`.",
			},

			"ha_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "active_active",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active_active",
					"active_standby",
				}, false),
				Description: "The high-availability mode of vpn tunnel. Valid values: `active_active` valid only when type as `Ipsec`; `active_active` and `active_standby` valid only when type as `RouteIpsec`.",
			},
			"open_health_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The switch of vpn tunnel health check. **Notes: that's valid only when vpn-v2.0 and tunnel type is `RouteIpsec`**.",
			},

			"local_peer_ip": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
				),
				Description: "The local IP in Kingsoft Cloud with CIDR indicated.",
			},

			"customer_peer_ip": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
				),
				Description: "The IP of customer with CIDR indicated.",
			},

			"vpn_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The vpn_gre_ip of the vpn tunnel. If type is GreOverIpsec and Vpn-Gateway-Version is 1.0, Required.",
			},

			"ha_vpn_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The ha_vpn_gre_ip of the vpn tunnel.If type is GreOverIpsec,Required and Vpn-Gateway-Version is 1.0, Required.",
			},

			"customer_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The customer_gre_ip of the vpn tunnel.If type is GreOverIpsec and Vpn-Gateway-Version is 1.0, Required.",
			},

			"ha_customer_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
				Description: "The ha_customer_gre_ip of the vpn tunnel.If type is GreOverIpsec and Vpn-Gateway-Version is 1.0, Required.",
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
			"ike_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1",
					"v2",
				}, false),
				Description: "the version of Ike.",
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
			"vpn_gateway_version": {
				Type: schema.TypeString,
				// Computed:    true,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"1.0", "2.0"},
					false,
				),
				Default:     "1.0",
				Description: "The version of vpn gateway. The version must be identical with `vpn_gate_way_version` of `ksyun_vpn_gateway`.",
			},
			// computed parameters
			"vpn_m_tunnel_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vpn first tunnel created time.",
			},
			"vpn_s_tunnel_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vpn second tunnel created time.",
			},
			"vpn_m_tunnel_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vpn first tunnel state.",
			},
			"vpn_s_tunnel_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vpn second tunnel state.",
			},
			"vpn_tunnel_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vpn tunnel created time.",
			},
		},
	}

	for k, s := range r.Schema {
		if isV2Attr, isV1Attr := stringSliceContains(vpnV2Attribute, k), stringSliceContains(vpnV1Attribute, k); isV2Attr || isV1Attr {
			s.DiffSuppressFunc = vpnV2ParamsDiffSuppressFunc

			if isV1Attr {
				s.Description += " Notes: it's valid when vpn gateway version is 1.0."
			}
			if isV2Attr {
				s.Description += " Notes: it's valid when vpn gateway version is 2.0."
			}
		}
	}
	return r
}

func resourceKsyunVpnTunnelCreate(d *schema.ResourceData, meta interface{}) (err error) {
	if err := checkVpnType(d); err != nil {
		return err
	}

	vpcService := VpcService{meta.(*KsyunClient)}

	vpnGateway, readGWErr := vpcService.ReadVpnGateway(d, d.Get("vpn_gateway_id").(string))
	if readGWErr != nil {
		return readGWErr
	}
	if v, ok := vpnGateway["VpnGatewayVersion"]; !ok {
		return fmt.Errorf("an error caused by checking vpn version of gateway and tunnel")
	} else {
		if v != d.Get("vpn_gateway_version") {
			return fmt.Errorf("vpn_gateway_version is not identical with `ksyun_vpn_gateway`, should keep same with ksyun_vpn_gateway.vpn_gateway_version")
		}
	}

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
	if err := checkVpnType(d); err != nil {
		return err
	}
	// check cannot be changed field
	// if err := bannedPartialParamsChanges(d); err != nil {
	// 	return err
	// }

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

// checkVpnVersionParams check parameters when vpn version diff
func checkVpnType(d *schema.ResourceData) error {
	var (
		isV2               = d.Get("vpn_gateway_version") == "2.0"
		vpnType            = d.Get("type")
		haMode             = d.Get("ha_mode")
		_, hcoExist        = d.GetOk("open_health_check")
		_, ikeVersionExist = d.GetOk("ike_version")

		isAllowHealthCheckOpen = false
		errs                   []error
	)

	if isV2 {
		if !ikeVersionExist {
			errs = append(errs, fmt.Errorf("ike_version cannot be blank, when vpn_gateway_version is 2.0. Should be set it value"))
		}

		switch vpnType {
		case "GreOverIpsec":
			errs = append(errs, fmt.Errorf("type GreOverIpsec and Ipsec is valid with vpn1.0, RouteIpsec and ipsec is valid with vpn2.0"))
		case "RouteIpsec":
			// allow open_health_check field valid
			isAllowHealthCheckOpen = true

			_, localExist := d.GetOk("local_peer_ip")
			_, customerExist := d.GetOk("customer_peer_ip")
			if !localExist || !customerExist {
				errs = append(errs, fmt.Errorf("customer_peer_ip and local_peer_ip cannot be blank, when vpn_gateway_version is 2.0 and vpn type is RouteIpsec"))
			}
		case "Ipsec":
			if haMode == "active_standby" {
				errs = append(errs, fmt.Errorf("value of ha_mode filed only as active_active, when vpn_gateway_version is 2.0 and vpn type is Ipsec"))
			}
		}

	} else {
		if vpnType == "RouteIpsec" {
			errs = append(errs, fmt.Errorf("type RouteIpsec and Ipsec is valid with vpn2.0, GreOverIpsec and Ipsec is valid with vpn1.0"))
		}

		if haMode == "active_standby" {
			errs = append(errs, fmt.Errorf("value of ha_mode filed only as active_active, when vpn_gateway_version is 2.0 and vpn type is Ipsec"))
		}

	}

	if !isAllowHealthCheckOpen && hcoExist {
		errs = append(errs, fmt.Errorf("open_health_check is valid, when vpn_gateway_version is 2.0 and vpn type is RouteIpsec"))
	}

	if errs != nil && len(errs) > 0 {
		return multierror.Append(nil, errs...)
	}

	return nil
}

func bannedPartialParamsChanges(d *schema.ResourceData) error {
	bannedList := []string{"customer_peer_ip", "local_peer_ip"}
	if d.HasChanges(bannedList...) {
		return fmt.Errorf("%s cannot be changed. if you need change it, you should create manually it as a new resource", strings.Join(bannedList, ", "))
	}
	return nil
}
