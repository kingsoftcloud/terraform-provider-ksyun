/*
Provides an alb backend server group resource.

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

```

# Import

ALB backend server group can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_backend_server_group.default fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunAlbBackendServerGroup() *schema.Resource {
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
		}
	}
	entry["host_name"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		// Default:      "default",
		ValidateFunc: validation.StringIsNotWhiteSpace,
		Description:  "hostname of the health check.",
	}

	return &schema.Resource{
		Create: resourceKsyunAlbBackendServerGroupCreate,
		Read:   resourceKsyunAlbBackendServerGroupRead,
		Update: resourceKsyunAlbBackendServerGroupUpdate,
		Delete: resourceKsyunAlbBackendServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ksc_bsg",
				Description: "The name of alb backend server group. Default: 'ksc_bsg'.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the VPC.",
			},
			"upstream_keepalive": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"adaptation",
					"keepalive",
					"shortconnection",
				}, false),
				Description: "The upstream keepalive type. Valid Value: `adaptation`, `keepalive`, `shortconnection`.",
			},
			"backend_server_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Host",
					"DirectConnect",
				}, false),
				Default:     "Host",
				ForceNew:    true,
				Description: "The type of backend server. Valid values: 'Host', 'DirectConnect'. Default is 'Host'.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The protocol of backend server. Valid values: 'HTTP', 'gRPC'. Default is 'HTTP'.",
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Deprecated:       "Alb does not support health checks at this time. If you need a health check configuration on this server group, you are supposed to use 'ksyun_alb_rule_group'",
				Removed:          "This parameter is removed in the latest version. Please use 'ksyun_alb_rule_group' to configure health check.",
				DiffSuppressFunc: lbBackendServerDiffSuppressFunc,
				Description:      "Health check information.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the alb backend server group.",
			},
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the alb backend server group.",
			},
			"backend_server_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "number of backend servers.",
			},
		},
	}
}
func resourceKsyunAlbBackendServerGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.CreateAlbBackendServerGroup(d, resourceKsyunAlbBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on creating alb backend server group %q, %s", d.Id(), err)
	}
	return resourceKsyunAlbBackendServerGroupRead(d, meta)
}

func resourceKsyunAlbBackendServerGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.ReadAndSetAlbBackendServerGroup(d, resourceKsyunAlbBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on reading alb backend server group %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunAlbBackendServerGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.ModifyAlbBackendServerGroup(d, resourceKsyunAlbBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on updating alb backend server group %q, %s", d.Id(), err)
	}
	return resourceKsyunAlbBackendServerGroupRead(d, meta)
}

func resourceKsyunAlbBackendServerGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	albService := AlbService{meta.(*KsyunClient)}
	err = albService.RemoveAlbBackendServerGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting alb backend server group %q, %s", d.Id(), err)
	}
	return err
}
