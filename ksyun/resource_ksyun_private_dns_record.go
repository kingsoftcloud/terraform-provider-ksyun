/*
Provides a Private Dns Record resource under PDNS Zone resource.

# Example Usage

```hcl
resource "ksyun_private_dns_zone" "foo" {
	zone_name = "tf-pdns-zone-pdns.com"
	zone_ttl = 360
	charge_type = "TrafficMonthly"
}

resource "ksyun_private_dns_record" "foo" {
	record_name = "tf-pdns-record"
	record_ttl = 360
	zone_id = ksyun_private_dns_zone.foo.id
	type = "CNAME"
	record_value = "tf-record.com"
}

```
*/

package ksyun

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunPrivateDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunPrivateDnsRecordCreate,
		Read:   resourceKsyunPrivateDnsRecordRead,
		Update: resourceKsyunPrivateDnsRecordUpdate,
		Delete: resourceKsyunPrivateDnsRecordDelete,

		Schema: map[string]*schema.Schema{
			"record_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The record name of private dns.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Private Dns Zone ID.",
			},

			"record_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(60, 86400),
				Description: "Record cache time. The smaller the value, the faster the record will take effect." +
					" Value range: 60~86400s.",
			},

			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: pdnsZoneRecordDiffSuppressFunc,
				Description:      "Record type. Valid values: \"A\", \"AAAA\", \"CNAME\", \"MX\", \"TXT\", \"SRV\".",
			},

			"record_value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Record value, such as IP: 192.168.10.2," +
					" CNAME: cname.ksyun.com, and MX: mail.ksyun.com..",
			},

			"priority": {
				Type:             schema.TypeInt,
				Computed:         true,
				Optional:         true,
				ForceNew:         true,
				Description:      "Record priority. Value range: [SRV|0~65535], [MX|1~100]. Required, when type is `SRV` or `MX`.",
				DiffSuppressFunc: pdnsZoneRecordDiffSuppressFunc,
				ValidateFunc:     validation.IntBetween(0, 65535),
			},
			"weight": {
				Type:             schema.TypeInt,
				Computed:         true,
				ForceNew:         true,
				Description:      "Record weight. Value range: 0~65535. Required, when type is `SRV`.",
				Optional:         true,
				DiffSuppressFunc: pdnsZoneRecordDiffSuppressFunc,
				ValidateFunc:     validation.IntBetween(0, 65535),
			},
			"port": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				Description:      "The port of record in which is associated with domain or ip. Required, when type is `SRV`.",
				DiffSuppressFunc: pdnsZoneRecordDiffSuppressFunc,
				ValidateFunc:     validation.IntBetween(0, 65535),
			},
		},
	}
}
func resourceKsyunPrivateDnsRecordCreate(d *schema.ResourceData, meta interface{}) (err error) {
	if err := requiredCheckWithType(d); err != nil {
		return err
	}

	sPrivateDnsRecordService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsRecordService.CreatePrivateDnsRecord(d, resourceKsyunPrivateDnsRecord())
	if err != nil {
		return fmt.Errorf("error on creating PrivateDnsRecord %q, %s", d.Id(), err)
	}
	return resourceKsyunPrivateDnsRecordRead(d, meta)
}

func resourceKsyunPrivateDnsRecordRead(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsRecordService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsRecordService.ReadAndSetPrivateDnsRecord(d, resourceKsyunPrivateDnsRecord())
	if err != nil {
		return fmt.Errorf("error on reading PrivateDnsRecord %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunPrivateDnsRecordUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsRecordService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsRecordService.ModifyPrivateDnsRecord(d, resourceKsyunPrivateDnsRecord())
	if err != nil {
		return fmt.Errorf("error on updating PrivateDnsRecord %q, %s", d.Id(), err)
	}
	return resourceKsyunPrivateDnsRecordRead(d, meta)
}

func resourceKsyunPrivateDnsRecordDelete(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsRecordService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsRecordService.DeletePrivateDnsRecord(d, resourceKsyunPrivateDnsRecord())
	if err != nil {
		return fmt.Errorf("error on deleting PrivateDnsRecord %q, %s", d.Id(), err)
	}
	return err
}

func requiredCheckWithType(d *schema.ResourceData) error {
	t := d.Get("type")
	if !(t == "SRV" || t == "MX") {
		return nil
	}

	if t == "SRV" {
		if !d.HasChanges("weight", "port", "priority") {
			return errors.New("'weight', 'port', 'priority' the three fields need to specify, when type is 'SRV'")
		}
	}
	if t == "MX" {
		if !d.HasChanges("priority") {
			return errors.New("'priority' the field needs to specify, when type is 'MX'")
		}
		priority := d.Get("priority").(int)
		if priority < 1 || priority > 100 {
			return errors.New("'priority' valid value between 1 and 100, when type is 'MX'")
		}
	}
	return nil
}
