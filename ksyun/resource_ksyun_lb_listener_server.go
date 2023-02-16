/*
Provides a Load Balancer Listener server resource.

# Example Usage

```hcl

	resource "ksyun_lb_listener_server" "default" {
	  listener_id = "3a520244-ddc1-41c8-9d2b-xxxxxxxxxxxx"
	  real_server_ip = "10.0.77.20"
	  real_server_port = 8000
	  real_server_type = "host"
	  instance_id = "3a520244-ddc1-41c8-9d2b-xxxxxxxxxxxx"
	  weight = 10
	}

```

# Import

LB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_listener.example vserver-abcdefg
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunInstancesWithListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstancesWithListenerCreate,
		Read:   resourceKsyunInstancesWithListenerRead,
		Update: resourceKsyunInstancesWithListenerUpdate,
		Delete: resourceKsyunInstancesWithListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The id of the listener.",
			},
			"real_server_ip": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The IP of real server.",
			},
			"real_server_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The port of real server.Valid Values:1-65535.",
			},
			"real_server_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"host",
					"DirectConnectGateway",
					"VpnTunnel",
				}, false),
				Default:     "host",
				Description: "The type of real server.Valid Values:'host', 'DirectConnectGateway', 'VpnTunnel'.",
			},
			"instance_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				DiffSuppressFunc: lbRealServerDiffSuppressFunc,
				Description:      "The ID of instance.",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
				Default:      1,
				Description:  "The weight of backend service.Valid Values:1-255.",
			},
			"master_slave_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Master",
					"Slave",
				}, false),
				Default:          "Master",
				DiffSuppressFunc: lbRealServerDiffSuppressFunc,
				Description:      "whether real server is master of salve. when listener method is MasterSlave, this field is supported.",
			},
			"real_server_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the real server.",
			},

			"register_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration ID of real server.",
			},

			"listener_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Forwarding mode of listener.",
			},
		},
	}
}
func resourceKsyunInstancesWithListenerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateRealServer(d, resourceKsyunInstancesWithListener())
	if err != nil {
		return fmt.Errorf("error on creating real server %q, %s", d.Id(), err)
	}
	return resourceKsyunInstancesWithListenerRead(d, meta)
}

func resourceKsyunInstancesWithListenerRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetRealServer(d, resourceKsyunInstancesWithListener())
	if err != nil {
		return fmt.Errorf("error on reading real server %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunInstancesWithListenerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyRealServer(d, resourceKsyunInstancesWithListener())
	if err != nil {
		return fmt.Errorf("error on updating real server %q, %s", d.Id(), err)
	}
	return resourceKsyunInstancesWithListenerRead(d, meta)
}

func resourceKsyunInstancesWithListenerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveRealServer(d)
	if err != nil {
		return fmt.Errorf("error on deleting real server %q, %s", d.Id(), err)
	}
	return err
}
