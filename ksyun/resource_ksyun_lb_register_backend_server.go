/*
Provides a lb register backend server resource.

# Example Usage

```hcl

	resource "ksyun_lb_register_backend_server" "default" {
		backend_server_group_id="xxxx"
		backend_server_ip="192.168.5.xxx"
		backend_server_port="8081"
		weight=10
	}

```

# Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_lb_register_backend_server.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunRegisterBackendServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRegisterBackendServerCreate,
		Read:   resourceKsyunRegisterBackendServerRead,
		Update: resourceKsyunRegisterBackendServerUpdate,
		Delete: resourceKsyunRegisterBackendServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of backend server group.",
			},
			"backend_server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The IP of backend server.",
			},
			"backend_server_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The port of backend server.Valid Values:1-65535.",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
				Default:      1,
				Description:  "The weight of backend service.Valid Values:0-255.",
			},
			"register_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration ID of binding server group.",
			},
			"real_server_ip": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IP of real server.",
			},
			"real_server_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port of real server.Valid Values:1-65535.",
			},
			"real_server_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of real server.Valid Values:'Host'.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of instance.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of network interface.",
			},
			"real_server_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of real server.Values:'healthy','unhealthy'.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the backend service was created.",
			},
		},
	}
}
func resourceKsyunRegisterBackendServerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateBackendServer(d, resourceKsyunRegisterBackendServer())
	if err != nil {
		return fmt.Errorf("error on creating backend server %q, %s", d.Id(), err)
	}
	return resourceKsyunRegisterBackendServerRead(d, meta)
}

func resourceKsyunRegisterBackendServerRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetBackendServer(d, resourceKsyunRegisterBackendServer())
	if err != nil {
		return fmt.Errorf("error on reading backend server %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunRegisterBackendServerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyBackendServer(d, resourceKsyunRegisterBackendServer())
	if err != nil {
		return fmt.Errorf("error on updating backend server %q, %s", d.Id(), err)
	}
	return resourceKsyunRegisterBackendServerRead(d, meta)
}

func resourceKsyunRegisterBackendServerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveBackendServer(d)
	if err != nil {
		return fmt.Errorf("error on deleting backend server %q, %s", d.Id(), err)
	}
	return err
}
