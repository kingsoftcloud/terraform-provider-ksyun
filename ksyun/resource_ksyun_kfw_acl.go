/*
Provides a Cloud Firewall ACL Rule resource.

# Example Usage

```hcl
resource "ksyun_kfw_acl" "default" {
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  acl_name        = "test-acl-rule"
  direction       = "in"
  src_type        = "ip"
  src_ips         = ["10.0.0.11", "10.0.0.21"]
  dest_type       = "ip"
  dest_ips        = ["10.0.0.31"]
  service_type    = "service"
  service_infos   = ["TCP:1-65535/1-65535"]
  app_type        = "any"
  policy          = "accept"
  status          = "start"
  priority_position = "after1"
  description     = "test acl rule"
}
```

# Import

Cloud Firewall ACL Rule can be imported using the `acl_id`, e.g.

```
$ terraform import ksyun_cfw_acl.default acl_id
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceKsyunKfwAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKfwAclCreate,
		Read:   resourceKsyunKfwAclRead,
		Update: resourceKsyunKfwAclUpdate,
		Delete: resourceKsyunKfwAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cfw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud Firewall Instance ID.",
			},
			"acl_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ACL Rule ID.",
			},
			"acl_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ACL rule name.",
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in",
					"out",
				}, false),
				Description: "Direction. Valid values: in (inbound), out (outbound).",
			},
			"src_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"addrbook",
					"zone",
					"any",
				}, false),
				Description: "Source address type. Valid values: ip, addrbook, zone, any.",
			},
			"src_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Source IP addresses.",
			},
			"src_addrbooks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Source address book IDs.",
			},
			"src_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Description: "Source zones (geographic regions).",
			},
			"dest_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"addrbook",
					"any",
				}, false),
				Description: "Destination address type. Valid values: ip, addrbook, any.",
			},
			"dest_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Destination IP addresses.",
			},
			"dest_addrbooks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Destination address book IDs.",
			},
			"dest_host": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Destination domain names.",
			},
			"dest_hostbook": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Destination host book IDs.",
			},
			"service_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"service",
					"servicegroup",
					"any",
				}, false),
				Description: "Service type. Valid values: service, servicegroup, any.",
			},
			"service_infos": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/80-80, UDP:22/33, ICMP.",
			},
			"service_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Service group IDs.",
			},
			"app_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"app",
					"any",
				}, false),
				Description: "Application type. Valid values: app, any.",
			},
			"app_value": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Application values.",
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"accept",
					"deny",
				}, false),
				Description: "Action. Valid values: accept, deny.",
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Description: "Status. Valid values: start, stop.",
			},
			"priority_position": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Priority position. Format: after+priority or before+priority. Example: after+1, before+1.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Priority value.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
			"hit_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Hit count.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
		},
	}
}

func resourceKsyunKfwAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*KsyunClient)
	service := KfwService{client: client}
	return service.createKfwAcl(d, resourceKsyunKfwAcl())
}

func resourceKsyunKfwAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*KsyunClient)
	service := KfwService{client: client}
	return service.readAndSetKfwAcl(d, resourceKsyunKfwAcl(), false)
}

func resourceKsyunKfwAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*KsyunClient)
	service := KfwService{client: client}
	return service.modifyKfwAcl(d, resourceKsyunKfwAcl())
}

func resourceKsyunKfwAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*KsyunClient)
	service := KfwService{client: client}
	return service.removeKfwAcl(d, resourceKsyunKfwAcl())
}

func timestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
