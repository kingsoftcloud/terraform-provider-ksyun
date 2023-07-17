/*
Provides an KNAD resource.

Example Usage

```hcl
# Create an knad
resource "ksyun_knad" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  knad_name = "ksc_knad"
  bill_type = 1
  service_id = "KNAD_30G"
  project_id="0"
}
```

Import

Knad can be imported using the id, e.g.

```
$ terraform import ksyun_knad.default knad67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "The id of the project.",
			},
			"knad_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the Knad.",
			},
			"link_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the link type of the Knad. Valid Values: 'DDoS_BGP'.",
			},
			"ip_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, 100),
				Description:  "the max ip count that can bind to the Knad,value range: [10, 100].",
			},
			"band": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the band of the Knad.",
			},
			"max_band": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the max band of the Knad.",
			},
			"idc_band": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the idcband of the Knad.",
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 36),
				Description:  "Purchase time.If bill_type is 1,this is Required.",
			},
			"knad_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the name of the Knad.",
			},
			"bill_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "the bill type of the Knad. Valid Values: 1:(PrePaidByMonth),5:(DailyPaidByTransfer).",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service id of the Knad.Valid Values:'KNAD_30G','KNAD_100G','KNAD_300G','KNAD_1000G',''KNAD_2000G''.",
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
		return fmt.Errorf("error on reading knad %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKnadUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	knadService := KnadService{meta.(*KsyunClient)}
	err = knadService.ModifyKnad(d, resourceKsyunKnad())
	if err != nil {
		return fmt.Errorf("error on updating knad %q, %s", d.Id(), err)
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
