package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatCreate,
		Update: resourceKsyunNatUpdate,
		Read:   resourceKsyunNatRead,
		Delete: resourceKsyunNatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nat_line_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"nat_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nat_mode": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Vpc",
					"Subnet",
				}, false),
			},

			"nat_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public",
				}, false),
			},

			"nat_ip_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"band_width": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 15000),
			},

			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "DailyPaidByTransfer",
				ValidateFunc: validation.StringInSlice([]string{
					"Monthly",
					"Peak",
					"Daily",
					"PostPaidByAdvanced95Peak",
					"DailyPaidByTransfer",
				}, false),
				DiffSuppressFunc: chargeSchemaDiffSuppressFunc,
			},

			"purchase_time": {
				Type:             schema.TypeInt,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(0, 36),
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
			},

			"nat_ip_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunNatCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNat(d, resourceKsyunNat())
	if err != nil {
		return fmt.Errorf("error on creating nat %q, %s", d.Id(), err)
	}
	return resourceKsyunNatRead(d, meta)
}

func resourceKsyunNatRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNat(d, resourceKsyunNat())
	if err != nil {
		return fmt.Errorf("error on reading nat %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNatUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyNat(d, resourceKsyunNat())
	if err != nil {
		return fmt.Errorf("error on updating nat %q, %s", d.Id(), err)
	}
	return resourceKsyunNatRead(d, meta)
}

func resourceKsyunNatDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNat(d)
	if err != nil {
		return fmt.Errorf("error on deleting nat %q, %s", d.Id(), err)
	}
	return err
}
