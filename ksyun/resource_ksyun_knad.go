package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKnad() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKnadCreate,
		Read:   resourceKsyunKnadRead,
		Update: resourceKsyunKnadUpdate,
		Delete: resourceKsyunKnadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"knad_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"band": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_band": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"idc_band": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"knad_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bill_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
func resourceKsyunKnadCreate(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.CreateKnad(d, resourceKsyunKnad())
	if err != nil {
		return fmt.Errorf("error on creating knad %q, %s", d.Id(), err)
	}
	return resourceKsyunKnadRead(d, meta)
}

func resourceKsyunKnadRead(d *schema.ResourceData, meta interface{}) (err error) {

	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.ReadAndSetKnad(d, resourceKsyunKnad())
	if err != nil {
		return fmt.Errorf("error on reading address %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKnadUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.ModifyKnad(d, resourceKsyunKnad())
	if err != nil {
		return fmt.Errorf("error on updating address %q, %s", d.Id(), err)
	}
	return resourceKsyunKnadRead(d, meta)
}

func resourceKsyunKnadDelete(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.RemoveKnad(d)
	if err != nil {
		return fmt.Errorf("error on deleting knad %q, %s", d.Id(), err)
	}
	return err

}
