/*
Provides a resource to create a PrivateDns zone_vpc_attachment

# Example Usage

```hcl
provider "ksyun" {
  region = "cn-guangzhou-1"
}

resource "ksyun_private_dns_zone" "foo" {
	zone_name = "tf-pdns-binding.com"
	zone_ttl = 360
	charge_type = "TrafficMonthly"
}

resource "ksyun_vpc" "foo" {
  vpc_name      = "tf-pdns-binding-vpc"
  cidr_block    = "10.7.0.0/21"
}

resource "ksyun_private_dns_zone_vpc_attachment" "example" {
	zone_id = ksyun_private_dns_zone.foo.id
	vpc_set {
		region_name = "cn-guangzhou-1"
		vpc_id = ksyun_vpc.foo.id
    }
}
```

# Import

Private Dns zone_vpc_attachment can be imported using the id, e.g.

```
terraform import ksyun_private_dns_zone_vpc_attachment.example ${zone_id}:${vpc_id}
```
*/

package ksyun

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunPrivateDnsZoneVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunPrivateDnsZoneVpcAttachmentCreate,
		Read:   resourceKsyunPrivateDnsZoneVpcAttachmentRead,
		Delete: resourceKsyunPrivateDnsZoneVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importPrivateDnsZoneVpcAttachment,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Private Dns Zone ID.",
			},

			"vpc_set": {
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"vpc_set"},
				Type:         schema.TypeList,
				Description:  "New add vpc info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vpc Id.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vpc region.",
						},
					},
				},
			},
		},
	}
}
func resourceKsyunPrivateDnsZoneVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneVpcAttachmentService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneVpcAttachmentService.BindZoneVpc(d, resourceKsyunPrivateDnsZoneVpcAttachment())
	if err != nil {
		var mErr error
		if strings.Contains(err.Error(), "bind vpc status error") {
			unBindErr := sPrivateDnsZoneVpcAttachmentService.UnbindZoneVpc(d, resourceKsyunPrivateDnsZoneVpcAttachment())
			if unBindErr != nil {
				mErr = multierror.Append(fmt.Errorf("an error caused when cleaning the error attachment %q, %s", d.Id(), unBindErr))
			}
		}

		mErr = multierror.Append(fmt.Errorf("error on creating private_dns_zone_vpc_attachment %q, %s", d.Id(), err))
		return mErr
	}
	return resourceKsyunPrivateDnsZoneVpcAttachmentRead(d, meta)
}

func resourceKsyunPrivateDnsZoneVpcAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneVpcAttachmentService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneVpcAttachmentService.ReadAndSetZoneVpcAttachment(d, resourceKsyunPrivateDnsZoneVpcAttachment())
	if err != nil {
		return fmt.Errorf("error on reading private_dns_zone_vpc_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunPrivateDnsZoneVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	sPrivateDnsZoneVpcAttachmentService := DnsService{meta.(*KsyunClient)}
	err = sPrivateDnsZoneVpcAttachmentService.UnbindZoneVpc(d, resourceKsyunPrivateDnsZoneVpcAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting private_dns_zone_vpc_attachment %q, %s", d.Id(), err)
	}
	return err
}
