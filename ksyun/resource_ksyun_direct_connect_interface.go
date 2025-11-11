/*
Provides a DirectConnectInterface resource.

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
```

Import

DCInterface can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_interface.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunDirectConnectInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDirectConnectInterfaceCreate,
		Read:   resourceKsyunDirectConnectInterfaceRead,
		Update: resourceKsyunDirectConnectInterfaceUpdate,
		Delete: resourceKsyunDirectConnectInterfaceDelete,
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
				ForceNew:    true,
				Description: "Route Type. Valid values: `BGP`, `STATIC`. Default is `BGP`. If set to `STATIC`, the customer must provide the BGP peer IP address and local peer IP address.",
			},
			"direct_connect_interface_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the direct connect interface. It is used to identify the direct connect interface.",
			},
			"direct_connect_interface_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The account ID of the direct connect interface. It is used to create a direct connect interface in another account.",
			},
			"customer_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Customer peer IP address. It is used to establish a BGP session with the customer.",
			},
			"local_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Local peer IP address. It is used to establish a BGP session with the customer.",
			},
			"ha_direct_connect_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Ha direct connect ID. It is used to create a high availability direct connect interface.",
			},
			"ha_customer_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Ha customer peer IP address.",
			},
			"ha_local_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Ha customer peer IP address.",
			},
			"bgp_peer": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The BGP peer IP address. It is used to establish a BGP session with the customer.",
			},
			"reliability_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Reliability method. Valid values: `bfd`, `nqa`. Default is `nqa`. If set to `BFD`, BFD configuration must be provided.",
			},
			"bfd_config_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the BFD configuration.",
			},
			"bgp_client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Bgp client token is used to ensure the idempotency of the request. It can be any string, but it must be unique for each request.",
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed:    true,
				Default:     false,
				Description: "Enable IPv6. Valid values: `true`, `false`. Default is `false`.",
			},
			"customer_ipv6_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Customer IPv6 peer IP address.",
			},
			"local_ipv6_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Local IPv6 peer IP address.",
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
				Description: "Ha direct connect interface name.",
			},
			"ha_direct_connect_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ha direct connect interface ID.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account ID of the direct connect interface.",
			},
			"ha_vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The id of vlan in direct connect.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Priority of the direct connect interface.",
			},
			"customer_peer_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer peer IPv6 address.",
			},
			"local_peer_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Local peer IPv6 address.",
			},
		},
	}
}

func resourceKsyunDirectConnectInterfaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcSrv := VpcService{meta.(*KsyunClient)}
	err = vpcSrv.CreateDirectConnectInterface(d, resourceKsyunDirectConnectInterface())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectInterface %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectInterfaceRead(d, meta)
}

func resourceKsyunDirectConnectInterfaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ReadAndSetDirectConnectInterface(d, resourceKsyunDirectConnectInterface())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectInterface %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectInterfaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ModifyDirectConnectInterface(d, resourceKsyunDirectConnectInterface())
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectInterface %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectInterfaceRead(d, meta)
}

func resourceKsyunDirectConnectInterfaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.RemoveDirectConnectInterface(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectInterface %q, %s", d.Id(), err)
	}
	return err
}
