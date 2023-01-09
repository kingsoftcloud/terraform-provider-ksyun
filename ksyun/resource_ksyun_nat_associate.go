/*
Provides a Nat Associate resource under VPC resource.

# Example Usage

```hcl

	resource "ksyun_vpc" "test" {
	  vpc_name = "ksyun-vpc-tf"
	  cidr_block = "10.0.0.0/16"
	}

	resource "ksyun_nat" "foo" {
	  nat_name = "ksyun-nat-tf"
	  nat_mode = "Subnet"
	  nat_type = "public"
	  band_width = 1
	  charge_type = "DailyPaidByTransfer"
	  vpc_id = "${ksyun_vpc.test.id}"
	}

	resource "ksyun_subnet" "test" {
	  subnet_name      = "tf-acc-subnet1"
	  cidr_block = "10.0.5.0/24"
	  subnet_type = "Normal"
	  dhcp_ip_from = "10.0.5.2"
	  dhcp_ip_to = "10.0.5.253"
	  vpc_id  = "${ksyun_vpc.test.id}"
	  gateway_ip = "10.0.5.1"
	  dns1 = "198.18.254.41"
	  dns2 = "198.18.254.40"
	  availability_zone = "cn-beijing-6a"
	}

	resource "ksyun_nat_associate" "foo" {
	  nat_id = "${ksyun_nat.foo.id}"
	  subnet_id = "${ksyun_subnet.test.id}"
	}

```

# Import

nat associate can be imported using the `id`, e.g.

```
$ terraform import ksyun_nat_associate.example $nat_id:$subnet_id
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunNatAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatAssociationCreate,
		Read:   resourceKsyunNatAssociationRead,
		Delete: resourceKsyunNatAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: importNatAssociate,
		},

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Nat.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Subnet.",
			},
		},
	}
}
func resourceKsyunNatAssociationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNatAssociate(d, resourceKsyunNatAssociation())
	if err != nil {
		return fmt.Errorf("error on creating nat associate %q, %s", d.Id(), err)
	}
	return resourceKsyunNatAssociationRead(d, meta)
}

func resourceKsyunNatAssociationRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNatAssociate(d, resourceKsyunNatAssociation())
	if err != nil {
		return fmt.Errorf("error on reading nat associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNatAssociationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNatAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting nat associate %q, %s", d.Id(), err)
	}
	return err
}
