/*
Provides a DirectConnectTunnel resource.

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

func resourceKsyunDirectConnectTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDirectConnectTunnelCreate,
		Read:   resourceKsyunDirectConnectTunnelRead,
		Update: resourceKsyunDirectConnectTunnelUpdate,
		Delete: resourceKsyunDirectConnectTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"direct_connect_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of direct connect. It's meaning is the physical port.",
			},
			"vlan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The id of vlan in direct connect.",
			},
			"route_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "bandwidth value, value range: [1, 15000].",
			},
			"direct_connect_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The charge type of the BWS. Valid values: PostPaidByPeak, PostPaidByDay, PostPaidByTransfer, DailyPaidByTransfer.",
			},
			"direct_connect_interface_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"customer_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"local_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"ha_direct_connect_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"ha_customer_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"ha_local_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"bgp_peer": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"reliability_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"bfd_config_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"bgp_client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"customer_ipv6_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"local_ipv6_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},

			// Computed fields

			"direct_connect_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the direct connect interface.",
			},

			"ha_direct_connect_interface_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"ha_direct_connect_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"ha_vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the project.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the project.",
			},
			"customer_peer_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
			"local_peer_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the project.",
			},
		},
	}
}

func resourceKsyunDirectConnectTunnelCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcSrv := VpcService{meta.(*KsyunClient)}
	err = vpcSrv.CreateDirectConnectTunnel(d, resourceKsyunDirectConnectTunnel())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectTunnel %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectTunnelRead(d, meta)
}

func resourceKsyunDirectConnectTunnelRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetDirectConnectTunnel(d, resourceKsyunDirectConnectTunnel())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectTunnel %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectTunnelUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ModifyDirectConnectTunnel(d, resourceKsyunDirectConnectTunnel())
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectTunnel %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectTunnelRead(d, meta)
}

func resourceKsyunDirectConnectTunnelDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.RemoveDirectConnectTunnel(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectTunnel %q, %s", d.Id(), err)
	}
	return err
}
