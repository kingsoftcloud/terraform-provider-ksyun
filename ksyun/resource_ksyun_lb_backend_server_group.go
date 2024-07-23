/*
Provides a lb backend server group resource.

# Example Usage

```hcl

	resource "ksyun_lb_backend_server_group" "default" {
		backend_server_group_name="xuan-tf"
		vpc_id=""
		backend_server_group_type=""
	}

```

# Import

LB backend server group can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_backend_server_group.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunBackendServerGroup() *schema.Resource {
	entry := resourceKsyunHealthCheck().Schema
	for k, v := range entry {
		if k == "listener_id" || k == "listener_protocol" || k == "is_default_host_name" || k == "host_name" {
			delete(entry, k)
		} else {
			v.ForceNew = false
			v.DiffSuppressFunc = nil
		}
		switch k {
		case "http_method":
			v.Optional = false
			v.Computed = true
			v.ValidateFunc = nil
		case "lb_type":
			delete(entry, k)
		}
	}
	entry["host_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "default",
		ValidateFunc: validation.StringIsNotEmpty,
		Description:  "hostname of the health check.",
	}

	return &schema.Resource{
		Create: resourceKsyunBackendServerGroupCreate,
		Read:   resourceKsyunBackendServerGroupRead,
		Update: resourceKsyunBackendServerGroupUpdate,
		Delete: resourceKsyunBackendServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "backend_server_group",
				Description: "The name of backend server group. Default: 'backend_server_group'.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC.",
			},
			"backend_server_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Server",
					"Mirror",
				}, false),
				Default:     "Server",
				ForceNew:    true,
				Description: "The type of backend server group. Valid values: 'Server', 'Mirror'. Default is 'Server'.",
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP",
					"UDP",
					"HTTP",
				}, false),
				Description: "The protocol of the backend server group. Valid values: 'TCP', 'UDP', 'HTTP'. Default `HTTP`.",
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbBackendServerDiffSuppressFunc,
				Description:      "Health check information, only the mirror server has this parameter.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the backend server group.",
			},
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the backend server group.",
			},
			"backend_server_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "number of backend servers.",
			},
		},
	}
}

func resourceKsyunBackendServerGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateBackendServerGroup(d, resourceKsyunBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on creating backend server group %q, %s", d.Id(), err)
	}
	return resourceKsyunBackendServerGroupRead(d, meta)
}

func resourceKsyunBackendServerGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetBackendServerGroup(d, resourceKsyunBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on reading backend server group %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBackendServerGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyBackendServerGroup(d, resourceKsyunBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on updating backend server group %q, %s", d.Id(), err)
	}
	return resourceKsyunBackendServerGroupRead(d, meta)
}

func resourceKsyunBackendServerGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveBackendServerGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting backend server group %q, %s", d.Id(), err)
	}
	return err
}
