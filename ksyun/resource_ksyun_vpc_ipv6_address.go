/*
Provides a Ipv6PublicIpAddress resource.

~> **Note**  The network segment can only be created or deleted, can not perform both of them at the same time.

Example Usage

```hcl
resource "ksyun_vpc_ipv6_address" "example" {
  network_interface_id   = "54098712-6a7b-4561-a84a-d36e9bb0b206"
  ipv6_public_ip_address = "2401:1d40:100:100:25:a:0:4e8"
  band_width=1
  charge_type=Daily
}

Import

Ipv6PublicIpAddress can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpc_ipv6_address.example 12395712-6a7b-4561-a84a-d36e5fb0b206
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunVpcIpv6Address() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpcIpv6AddressCreate,
		Update: resourceKsyunVpcIpv6AddressUpdate,
		Read:   resourceKsyunVpcIpv6AddressRead,
		Delete: resourceKsyunVpcIpv6AddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"network_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network interface id.",
			},
			"ipv6_public_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ipv6 public ip address.",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The band width of the Ipv6PublicIp.",
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

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of creation for Ipv6PublicIp.",
			},
		},
	}
}

func resourceKsyunVpcIpv6AddressCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpcIpv6Address(d, resourceKsyunVpcIpv6Address())
	if err != nil {
		return fmt.Errorf("error on creating vpcIpv6Address %q, %s", d.Id(), err)
	}
	return resourceKsyunVpcIpv6AddressRead(d, meta)
}

func resourceKsyunVpcIpv6AddressRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpcIpv6Address(d, resourceKsyunVpcIpv6Address())
	if err != nil {
		return fmt.Errorf("error on reading vpcIpv6Address %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpcIpv6AddressUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpcIpv6Address(d, resourceKsyunVpcIpv6Address())
	if err != nil {
		return fmt.Errorf("error on updating vpcIpv6Address %q, %s", d.Id(), err)
	}
	return resourceKsyunVpcIpv6AddressRead(d, meta)
}

func resourceKsyunVpcIpv6AddressDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpcIpv6Address(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpcIpv6Address %q, %s", d.Id(), err)
	}
	return err
}
