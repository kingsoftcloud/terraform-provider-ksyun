/*
Provides a Nat resource under VPC resource.

# Example Usage

```hcl

	resource "ksyun_vpc" "test" {
	  vpc_name = "ksyun-vpc-tf"
	  cidr_block = "10.0.0.0/16"
	}

	resource "ksyun_nat" "foo" {
	  nat_name = "ksyun-nat-tf"
	  nat_mode = "Vpc"
	  nat_type = "public"
	  band_width = 1
	  charge_type = "DailyPaidByTransfer"
	  vpc_id = "${ksyun_vpc.test.id}"
	}

```

# Import

nat can be imported using the `id`, e.g.

```
$ terraform import ksyun_nat.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/
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
				Computed: true,
				// ForceNew:    true,
				Description: "ID of the line.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the VPC.",
			},
			"nat_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the NAT.",
			},

			"nat_mode": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Vpc",
					"Subnet",
				}, false),
				Description: "Mode of the NAT, valid values: 'Vpc', 'Subnet'.",
			},

			"nat_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public",
				}, false),
				Description: "Type of the NAT, valid values: 'public'.",
			},

			"nat_ip_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 20),
				Description:  "The Counts of Nat Ip, value range:[1, 20], Default is 1.",
			},

			"band_width": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 15000),
				// Default:      1,
				Description: "The BandWidth of Nat Ip, value range:[1, 15000], Default is 1.",
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
				Description:      "charge type, valid values: 'Monthly', 'Peak', 'Daily', 'PostPaidByAdvanced95Peak', 'DailyPaidByTransfer'. Default is DailyPaidByTransfer.",
			},

			"purchase_time": {
				Type:             schema.TypeInt,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(0, 36),
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
				Description:      "The PurchaseTime of the Nat, value range [1, 36]. If charge_type is Monthly this Field is Required.",
			},

			"tags": tagsSchema(),

			"nat_ip_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The nat ip list of the desired Nat.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT IP address.",
						},
						"nat_ip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the NAT IP.",
						},
					},
				},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of creation of Nat.",
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
