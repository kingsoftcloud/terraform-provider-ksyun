/*
Provides a DirectConnectBfdConfig resource.

Example Usage

```hcl
resource "ksyun_direct_connect_bfd_config" "test" {
  min_tx_interval   = 100
  min_rx_interval   = 200
  detect_multiplier = 3
  multi_hop         = true
}
```

Import

ksyun_direct_connect_bfd_config can be imported using the id, e.g.

```
$ terraform import ksyun_direct_connect_bfd_config.test 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunDirectConnectBfdConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDirectConnectBfdConfigCreate,
		Read:   resourceKsyunDirectConnectBfdConfigRead,
		Update: resourceKsyunDirectConnectBfdConfigUpdate,
		Delete: resourceKsyunDirectConnectBfdConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"min_tx_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The interval at which the BFD control packets are sent.",
			},
			"min_rx_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The interval at which the BFD control packets are received.",
			},
			"detect_multiplier": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Detect Multiplier.",
			},

			"multi_hop": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the multi hop.",
			},

			"bfd_config_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the BFD configuration.",
			},
		},
	}
}

func resourceKsyunDirectConnectBfdConfigCreate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.CreateDirectConnectBfdConfig(d, resourceKsyunDirectConnectBfdConfig())
	if err != nil {
		return fmt.Errorf("error on creating DirectConnectBfdConfig %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectBfdConfigRead(d, meta)
}

func resourceKsyunDirectConnectBfdConfigRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ReadAndSetDirectConnectBfdConfig(d, resourceKsyunDirectConnectBfdConfig())
	if err != nil {
		return fmt.Errorf("error on reading DirectConnectBfdConfig %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDirectConnectBfdConfigUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.ModifyDirectConnectBfdConfig(d, resourceKsyunDirectConnectBfdConfig())
	if err != nil {
		return fmt.Errorf("error on updating DirectConnectBfdConfig %q, %s", d.Id(), err)
	}
	return resourceKsyunDirectConnectBfdConfigRead(d, meta)
}

func resourceKsyunDirectConnectBfdConfigDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	err = srv.RemoveDirectConnectBfdConfig(d)
	if err != nil {
		return fmt.Errorf("error on deleting DirectConnectBfdConfig %q, %s", d.Id(), err)
	}
	return err
}
