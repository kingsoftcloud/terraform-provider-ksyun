/*
Provides alb register alb backend server group resource.

# Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_alb_backend_server_group" "foo" {
  name="tf-alb-bsg"
  vpc_id=ksyun_vpc.test.id
  upstream_keepalive="adaptation"
  backend_server_type="Host"
}

resource "ksyun_alb_register_backend_server" "foo" {
  backend_server_group_id=ksyun_alb_backend_server_group.foo.id
  backend_server_ip=ksyun_instance.test.private_ip_address
  port = 8080
  weight=40
}

```

# Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_alb_register_backend_server.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunRegisterAlbBackendServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRegisterAlbBackendServerCreate,
		Read:   resourceKsyunRegisterAlbBackendServerRead,
		Update: resourceKsyunRegisterAlbBackendServerUpdate,
		Delete: resourceKsyunRegisterAlbBackendServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of alb backend server group.",
			},
			"direct_connect_gateway_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_interface_id"},
				Description:   "The ID of direct connect gateway.",
			},
			"network_interface_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"direct_connect_gateway_id"},
				Description:   "The ID of network interface.",
			},

			"backend_server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The IP of alb backend server.",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The port of alb backend server. Valid Values:1-65535.",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
				Default:      20,
				Description:  "The weight of backend service. Valid Values:0-255.",
			},
			"backend_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration ID of binding server group.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of instance.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the backend service was created.",
			},
		},
	}
}

func resourceKsyunRegisterAlbBackendServerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.CreateAlbBackendServer(d, resourceKsyunRegisterAlbBackendServer())
	if err != nil {
		return fmt.Errorf("error on creating alb backend server %q, %s", d.Id(), err)
	}
	return resourceKsyunRegisterAlbBackendServerRead(d, meta)
}

func resourceKsyunRegisterAlbBackendServerRead(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.ReadAndSetAlbBackendServer(d, resourceKsyunRegisterAlbBackendServer())
	if err != nil {
		return fmt.Errorf("error on reading alb backend server %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunRegisterAlbBackendServerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.ModifyAlbBackendServer(d, resourceKsyunRegisterAlbBackendServer())
	if err != nil {
		return fmt.Errorf("error on updating alb backend server %q, %s", d.Id(), err)
	}
	return resourceKsyunRegisterAlbBackendServerRead(d, meta)
}

func resourceKsyunRegisterAlbBackendServerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.RemoveAlbBackendServer(d)
	if err != nil {
		return fmt.Errorf("error on deleting alb backend server %q, %s", d.Id(), err)
	}
	return err
}
