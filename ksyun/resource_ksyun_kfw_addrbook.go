/*
Provides a Cloud Firewall Address Book resource.

# Example Usage

```hcl
resource "ksyun_kfw_addrbook" "default" {
  cfw_instance_id = "cd8e449c-8ac2-4ce2-a5fa-1448b1309a84"
  addrbook_name    = "test-addrbook"
  ip_version       = "IPv4"
  ip_address       = ["10.1.1.11", "10.2.2.21"]
  description      = "test address book"
}
```

# Import

Cloud Firewall Address Book can be imported using the `addrbook_id`, e.g.

```
$ terraform import ksyun_kfw_addrbook.default addrbook_id
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunKfwAddrbook() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKfwAddrbookCreate,
		Read:   resourceKsyunKfwAddrbookRead,
		Update: resourceKsyunKfwAddrbookUpdate,
		Delete: resourceKsyunKfwAddrbookDelete,
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
			"addrbook_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address Book ID.",
			},
			"addrbook_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address book name.",
			},
			"ip_version": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IPv4",
					"IPv6",
				}, false),
				Description: "IP version. Valid values: IPv4, IPv6.",
			},
			"ip_address": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IP addresses.",
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
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
		},
	}
}

func resourceKsyunKfwAddrbookCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.createKfwAddrbook(d, resourceKsyunKfwAddrbook())
}

func resourceKsyunKfwAddrbookRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.readAndSetKfwAddrbook(d, resourceKsyunKfwAddrbook(), false)
}

func resourceKsyunKfwAddrbookUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.modifyKfwAddrbook(d, resourceKsyunKfwAddrbook())
}

func resourceKsyunKfwAddrbookDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := KfwService{client: meta.(*KsyunClient)}
	return s.removeKfwAddrbook(d, resourceKsyunKfwAddrbook())
}
