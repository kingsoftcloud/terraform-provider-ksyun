package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunBandWidthShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBandWidthShareCreate,
		Read:   resourceKsyunBandWidthShareRead,
		Update: resourceKsyunBandWidthShareUpdate,
		Delete: resourceKsyunBandWidthShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"line_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"band_width_share_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 15000),
			},
			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PostPaidByPeak",
					"PostPaidByDay",
					"PostPaidByTransfer",
				}, false),
				DiffSuppressFunc: chargeSchemaDiffSuppressFunc,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceKsyunBandWidthShareCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.CreateBandWidthShare(d, resourceKsyunBandWidthShare())
	if err != nil {
		return fmt.Errorf("error on creating bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunBandWidthShareRead(d, meta)
}

func resourceKsyunBandWidthShareRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetBandWidthShare(d, resourceKsyunBandWidthShare())
	if err != nil {
		return fmt.Errorf("error on reading bandWidthShare %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBandWidthShareUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ModifyBandWidthShare(d, resourceKsyunBandWidthShare())
	if err != nil {
		return fmt.Errorf("error on updating bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunBandWidthShareRead(d, meta)
}

func resourceKsyunBandWidthShareDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.RemoveBandWidthShare(d)
	if err != nil {
		return fmt.Errorf("error on deleting bandWidthShare %q, %s", d.Id(), err)
	}
	return err

}
