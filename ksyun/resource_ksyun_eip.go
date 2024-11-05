/*
Provides an Elastic IP resource.

Example Usage

```hcl
data "ksyun_lines" "default" {
  output_file="output_result1"
  line_name="BGP"
}
resource "ksyun_eip" "default" {
  line_id ="${data.ksyun_lines.default.lines.0.line_id}"
  band_width =1
  charge_type = "PostPaidByPeak"
  purchase_time =1
  project_id=0
}
```

Import

EIP can be imported using the id, e.g.

```
$ terraform import ksyun_eip.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The id of the line.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The band width of the public address.",
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
				Description: "The charge type of the Elastic IP address.Valid Values:'PrePaidByMonth','Monthly','PostPaidByPeak','Peak','PostPaidByDay','Daily','PostPaidByTransfer','TrafficMonthly','DailyPaidByTransfer','HourlySettlement','PostPaidByHour','HourlyInstantSettlement','PostpaidByTime'. \n" +
					"**Notes:** Charge Type have a upgrade, The above-mentioned parameters, **every**, are **valid**. The changes as following:\n\n" +
					"| Previous Version | Current Version | Description | \n" +
					"| -------- | -------- | ----------- | \n" +
					"| PostPaidByPeak | Peak| Pay-as-you-go (monthly peak) | \n " +
					"| PostPaidByDay | Daily | Pay-as-you-go (daily) | \n" +
					"| PostPaidByTransfer | TrafficMonthly | Pay-as-you-go (monthly traffic) |\n " +
					"| PrePaidByMonth | Monthly | Monthly package | \n" +
					"|                | DailyPaidByTransfer | Pay-as-you-go (daily traffic) | \n" +
					"|                | HourlyInstantSettlement | Pay-as-you-go (hourly instant settlement) | \n" +
					"|                | PostPaidByHour | Pay-as-you-go (hourly billing, monthly settlement) | \n" +
					"|                | PostpaidByTime | Settlement by times |.",
			},
			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(0, 36),
				Description:      "Purchase time. If charge_type is Monthly or PrePaidByMonth, this is Required.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "The id of the project.",
				// Computed:    true,
			},
			"tags": tagsSchema(),

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the EIP.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the EIP.",
			},

			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance type to bind with the EIP.",
			},

			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Elastic IP address.",
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the EIP.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the EIP.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NetworkInterface ID.",
			},
			"internet_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "InternetGateway ID.",
			},
			"band_width_share_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the BWS which the EIP associated.",
			},
			"is_band_width_share": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "BWS EIP.",
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
