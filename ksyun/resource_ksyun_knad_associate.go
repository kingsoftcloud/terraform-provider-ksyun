package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKnadAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKnadAssociateCreate,
		Read:   resourceKsyunKnadAssociateRead,
		Update: resourceKsyunKnadAssociateUpdate,
		Delete: resourceKsyunKnadAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importKnadAssociate,
		},

		Schema: map[string]*schema.Schema{
			"knad_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeSet,
				Optional: true,
				//ForceNew: true,
				//Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceKsyunKnadAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.AssociateKnad(d, resourceKsyunKnadAssociate())
	if err != nil {
		return fmt.Errorf("error on associate knad %q, %s", d.Id(), err)
	}
	return resourceKsyunKnadAssociateRead(d, meta)
}

func resourceKsyunKnadAssociateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.AssociateKnad(d, resourceKsyunKnadAssociate())
	if err != nil {
		return fmt.Errorf("error on associate knad %q, %s", d.Id(), err)
	}
	return resourceKsyunKnadAssociateRead(d, meta)
}

func resourceKsyunKnadAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.ReadAndSetAssociateKnad(d, resourceKsyunKnadAssociate())
	if err != nil {
		return fmt.Errorf("error on reading knad associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKnadAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.DisassociateKnad(d)
	if err != nil {
		return fmt.Errorf("error on disAssociate knad %q, %s", d.Id(), err)
	}
	return err
}
