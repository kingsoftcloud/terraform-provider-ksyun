/*
Provides a KEC network interface attachment resource

# Example Usage

```hcl

	resource "ksyun_kec_network_interface_attachment" "default" {
	 network_interface_id = "ebd74f60-04f1-4b67-91e0-xxxxxxxxxxxx"
	 instance_id = "110d1ce0-113e-4019-8b39-xxxxxxxxxxxx"
	}

```

# Import

KEC network interface attachment can be imported using the id, e.g.

```
$ terraform import ksyun_kec_network_interface_attachment.default ${network_interface_id}:${instance_id}
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKecNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKecNetworkInterfaceAttachmentCreate,
		Read:   resourceKsyunKecNetworkInterfaceAttachmentRead,
		Delete: resourceKsyunKecNetworkInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importKecNetworkInterfaceAttachment,
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the network interface.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the instance.",
			},
			"network_interface_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the network interface.",
			},
		},
	}
}

func resourceKsyunKecNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.readAndSetNetworkInterfaceAttachment(d, resourceKsyunKecNetworkInterfaceAttachment())
	if err != nil {
		return fmt.Errorf("error on reading network interface attachement %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKecNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.createNetworkInterfaceAttachment(d, resourceKsyunKecNetworkInterfaceAttachment())
	if err != nil {
		return fmt.Errorf("error on creating network interface attachement %q, %s", d.Id(), err)
	}
	return resourceKsyunKecNetworkInterfaceAttachmentRead(d, meta)
}

func resourceKsyunKecNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.modifyNetworkInterfaceAttachment(d, resourceKsyunKecNetworkInterfaceAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting network interface attachement %q, %s", d.Id(), err)
	}
	return err
}
