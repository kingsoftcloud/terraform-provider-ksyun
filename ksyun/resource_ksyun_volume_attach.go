/*
Provides an EBS attachment resource.

# Example Usage

```hcl

	resource "ksyun_volume_attach" "default" {
	  volume_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  instance_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  delete_with_instance = true
	}

```

# Import

EBS volume can be imported using the `id`, e.g.

```
$ terraform import ksyun_volume.default $volume_id:$instance_id
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunVolumeAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeAttachCreate,
		Read:   resourceKsyunVolumeAttachRead,
		Update: resourceKsyunVolumeAttachUpdate,
		Delete: resourceKsyunVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			State: importVolumeAttach,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the EBS volume.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the KEC instance to which the EBS volume is to be attached.",
			},
			"delete_with_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether to delete the EBS volume when the KEC instance to which it is attached is deleted. Default value: false.",
			},
			"volume_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the EBS volume.",
			},
			"volume_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the EBS volume.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the EBS volume was created.",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the EBS volume.",
			},
			"volume_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The category of the EBS volume.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability zone in which the EBS volume resides.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity of the EBS volume, in GB.",
			},
			"volume_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EBS volume.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the project.",
			},
		},
	}
}

func resourceKsyunVolumeAttachCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.CreateVolumeAttach(d, resourceKsyunVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on creating volume attach %q, %s", d.Id(), err)
	}
	return resourceKsyunVolumeAttachRead(d, meta)
}

func resourceKsyunVolumeAttachRead(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ReadAndSetVolumeAttach(d, resourceKsyunVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on reading volume attach %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVolumeAttachUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ModifyVolumeAttach(d, resourceKsyunVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on updating volume attach %q, %s", d.Id(), err)
	}
	return resourceKsyunVolumeAttachRead(d, meta)
}

func resourceKsyunVolumeAttachDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.RemoveVolumeAttach(d)
	if err != nil {
		return fmt.Errorf("error on deleting volume attach %q, %s", d.Id(), err)
	}
	return err
}
