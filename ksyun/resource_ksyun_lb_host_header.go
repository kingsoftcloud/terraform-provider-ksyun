/*
Provides a lb host header resource.

# Example Usage

```hcl

	resource "ksyun_lb_host_header" "default" {
		listener_id = "xxxx"
		host_header = "tf-xuan"
		certificate_id = ""
	}

```

EIP can be imported using the id, e.g.

```
terraform import ksyun_lb_host_header.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunListenerHostHeader() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerHostHeaderCreate,
		Read:   resourceKsyunListenerHostHeaderRead,
		Update: resourceKsyunListenerHostHeaderUpdate,
		Delete: resourceKsyunListenerHostHeaderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the listener.",
			},
			"host_header": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The host header.",
			},
			"certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: hostHeaderDiffSuppressFunc,
				Description:      "The ID of the certificate, HTTPS type listener creates this parameter which is not default.",
			},
			"host_header_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host header id.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the host header was created.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol of the listener.",
			},
		},
	}
}
func resourceKsyunListenerHostHeaderCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateHostHeader(d, resourceKsyunListenerHostHeader())
	if err != nil {
		return fmt.Errorf("error on creating host header %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerHostHeaderRead(d, meta)
}

func resourceKsyunListenerHostHeaderRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetHostHeader(d, resourceKsyunListenerHostHeader())
	if err != nil {
		return fmt.Errorf("error on reading host header %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunListenerHostHeaderUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyHostHeader(d, resourceKsyunListenerHostHeader())
	if err != nil {
		return fmt.Errorf("error on updating host header %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerHostHeaderRead(d, meta)
}

func resourceKsyunListenerHostHeaderDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveHostHeader(d)
	if err != nil {
		return fmt.Errorf("error on deleting host header %q, %s", d.Id(), err)
	}
	return err
}
