/*
Provides a Cloud Firewall Service Group resource.

# Example Usage

```hcl
resource "ksyun_kfw_service_group" "default" {
  cfw_instance_id   = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  service_group_name = "test-service-group"
  service_infos      = ["TCP:1-65535/1-65535", "UDP:22/33"]
  description        = "test service group"
}
```

# Import

Cloud Firewall Service Group can be imported using the `service_group_id`, e.g.

```
$ terraform import ksyun_kfw_service_group.default service_group_id
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKfwServiceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKfwServiceGroupCreate,
		Read:   resourceKsyunKfwServiceGroupRead,
		Update: resourceKsyunKfwServiceGroupUpdate,
		Delete: resourceKsyunKfwServiceGroupDelete,
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
			"service_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service Group ID.",
			},
			"service_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service group name.",
			},
			"service_infos": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Service information. Format: protocol:src_port_min-src_port_max/dest_port_min-dest_port_max. Example: TCP:1-100/80-80, UDP:22/33, ICMP.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
			"citation_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of references.",
			},
		},
	}
}

func resourceKsyunKfwServiceGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.createKfwServiceGroup(d, resourceKsyunKfwServiceGroup())
}

func resourceKsyunKfwServiceGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.readAndSetKfwServiceGroup(d, resourceKsyunKfwServiceGroup(), false)
}

func resourceKsyunKfwServiceGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.modifyKfwServiceGroup(d, resourceKsyunKfwServiceGroup())
}

func resourceKsyunKfwServiceGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.removeKfwServiceGroup(d, resourceKsyunKfwServiceGroup())
}
