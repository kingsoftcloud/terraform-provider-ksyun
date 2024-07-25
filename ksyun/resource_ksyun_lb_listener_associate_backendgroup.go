/*
Provides slb listener mount backend server group resource.

# Example Usage

```hcl
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc-1"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}

resource "ksyun_lb" "default" {
  vpc_id             = ksyun_vpc.test.id
  load_balancer_name = "tf-xun1"
  type               = "public"
}


data "ksyun_certificates" "default" {
}


resource "ksyun_lb_listener" "default" {
  listener_name     = "tf-xun"
  listener_port     = "8000"
  listener_protocol = "TCP"
  listener_state    = "start"
  load_balancer_id  = ksyun_lb.default.id
  method            = "LeastConnections"
  bind_type         = "BackendServerGroup"
  certificate_id    = data.ksyun_certificates.default.certificates.0.certificate_id
}

resource "ksyun_lb_backend_server_group" "default" {
  backend_server_group_name = "xuan-tf"
  vpc_id                    = ksyun_vpc.test.id
  backend_server_group_type = "Server"
  protocol                  = "TCP"
  health_check {
    host_name           = "www.ksyun.com"
    healthy_threshold   = 10
    interval            = 100
    timeout             = 300
    unhealthy_threshold = 10
  }
}

# associate backend server group with listener
resource "ksyun_lb_listener_associate_backendgroup" "mount" {
  listener_id             = ksyun_lb_listener.default.id
  backend_server_group_id = ksyun_lb_backend_server_group.default.id
}

```

# Import

resource can be imported using the id, e.g.

```
$ terraform import ksyun_lb_listener_associate_backendgroup.default $listener_id:$backend_server_group_id
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunLbListenerAssociateBackendgroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLbListenerAssociateBackendgroupCreate,
		Read:   resourceKsyunLbListenerAssociateBackendgroupRead,
		Update: nil,
		Delete: resourceKsyunLbListenerAssociateBackendgroupDelete,
		Importer: &schema.ResourceImporter{
			State: importListenerAssociateBackendgroup,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of alb backend server group.",
			},

			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of slb listener.",
			},
		},
	}
}

func resourceKsyunLbListenerAssociateBackendgroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.MountOrUnmountBackendGroup(d, resourceKsyunLbListenerAssociateBackendgroup(), true)
	if err != nil {
		return fmt.Errorf("error on mounting backend group onto listener %q, %s", d.Id(), err)
	}

	id := AssembleIds(d.Get("listener_id").(string), d.Get("backend_server_group_id").(string))
	d.SetId(id)
	return resourceKsyunLbListenerAssociateBackendgroupRead(d, meta)
}

func resourceKsyunLbListenerAssociateBackendgroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadListenerBackendGroups(d)
	if err != nil {
		return fmt.Errorf("error on reading backend group association %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunLbListenerAssociateBackendgroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.MountOrUnmountBackendGroup(d, resourceKsyunLbListenerAssociateBackendgroup(), false)
	if err != nil {
		return fmt.Errorf("error on debongding backend group from listener %q, %s", d.Id(), err)
	}
	return err
}
