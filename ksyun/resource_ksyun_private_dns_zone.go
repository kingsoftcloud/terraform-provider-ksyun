/*
Provides a Private Dns Zone resource.

# Example Usage

```hcl
resource "ksyun_private_dns_zone" "foo" {
	zone_name = "tf-pdns-zone-pdns.com"
	zone_ttl = 360
	charge_type = "TrafficMonthly"
}

```

# Import

Private Dns Record can be imported using the `id`, e.g.

```
$ terraform import ksyun_private_dns_zone.foo fdeba8ca-8aa6-4cd0-8ffa-xxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunPrivateDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunPrivateDnsZoneCreate,
		Read:   resourceKsyunPrivateDnsZoneRead,
		Update: resourceKsyunPrivateDnsZoneUpdate,
		Delete: resourceKsyunPrivateDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The zone name of private dns.",
			},
			"zone_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(60, 86400),
				Description: "The zone cache time. The smaller the value, the faster the record will take effect." +
					" Value range: 60~86400s.",
			},

			"project_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},

			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TrafficMonthly",
				}, false),
				DiffSuppressFunc: chargeSchemaDiffSuppressFunc,
				Description:      "The charge type of the Private Dns Zone. Values: `TrafficMonthly`.",
			},

			"bind_vpc_set": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "This zone have bound VPC set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of binding VPC.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC name.",
						},
					},
				},
			},
		},
	}
}
func resourceKsyunPrivateDnsZoneCreate(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneService.CreatePrivateDnsZone(d, resourceKsyunPrivateDnsZone())
	if err != nil {
		return fmt.Errorf("error on creating PrivateDnsZone %q, %s", d.Id(), err)
	}
	return resourceKsyunPrivateDnsZoneRead(d, meta)
}

func resourceKsyunPrivateDnsZoneRead(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneService.ReadAndSetPrivateDnsZone(d, resourceKsyunPrivateDnsZone())
	if err != nil {
		return fmt.Errorf("error on reading PrivateDnsZone %q, %s", d.Id(), err)
	}

	return err
}

func resourceKsyunPrivateDnsZoneUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneService.ModifyPrivateDnsZone(d, resourceKsyunPrivateDnsZone())
	if err != nil {
		return fmt.Errorf("error on updating PrivateDnsZone %q, %s", d.Id(), err)
	}

	// for waiting backend cache consistence
	time.Sleep(3 * time.Second)
	return resourceKsyunPrivateDnsZoneRead(d, meta)
}

func resourceKsyunPrivateDnsZoneDelete(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneService.DeletePrivateDnsZone(d)
	if err != nil {
		return fmt.Errorf("error on deleting PrivateDnsZone %q, %s", d.Id(), err)
	}
	return err
}
