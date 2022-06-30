package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipCreate,
		Read:   resourceKsyunEipRead,
		Update: resourceKsyunEipUpdate,
		Delete: resourceKsyunEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"line_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PrePaidByMonth",
					"Monthly",
					"PostPaidByPeak",
					"Peak",
					"PostPaidByDay",
					"Daily",
					"PostPaidByTransfer",
					"TrafficMonthly",
					"DailyPaidByTransfer",
					"HourlySettlement",
					"PostPaidByHour",
					"HourlyInstantSettlement",
					"PostpaidByTime",
				}, false),
				DiffSuppressFunc: chargeSchemaDiffSuppressFunc,
			},
			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(0, 36),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  0,
			},
			"tags": tagsSchema(),

			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allocation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_band_width_share": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunEipCreate(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.CreateAddress(d, resourceKsyunEip())
	if err != nil {
		return fmt.Errorf("error on creating address %q, %s", d.Id(), err)
	}
	return resourceKsyunEipRead(d, meta)
}

func resourceKsyunEipRead(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.ReadAndSetAddress(d, resourceKsyunEip())
	if err != nil {
		return fmt.Errorf("error on reading address %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunEipUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.ModifyAddress(d, resourceKsyunEip())
	if err != nil {
		return fmt.Errorf("error on updating address %q, %s", d.Id(), err)
	}
	return resourceKsyunEipRead(d, meta)
}

func resourceKsyunEipDelete(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.RemoveAddress(d)
	if err != nil {
		return fmt.Errorf("error on deleting address %q, %s", d.Id(), err)
	}
	return err

}
