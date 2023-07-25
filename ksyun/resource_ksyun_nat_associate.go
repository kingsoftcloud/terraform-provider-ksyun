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

    resource "ksyun_instance" "foo" {
      image_id="${data.ksyun_images.centos-7_5.images.0.image_id}"
      instance_type="N3.2B"

      subnet_id="${ksyun_subnet.foo.id}"
      instance_password="Xuan663222"
      keep_image_login=false
      charge_type="Daily"
      purchase_time=1
      security_group_id=["${ksyun_security_group.default.id}"]
      instance_name="ksyun-kec-tf-nat"
      sriov_net_support="false"
	}

	resource "ksyun_nat_associate" "foo" {
	  nat_id = "${ksyun_nat.foo.id}"
	  subnet_id = "${ksyun_subnet.test.id}"
	}
	resource "ksyun_nat_associate" "associate_ins" {
	  nat_id = "${ksyun_nat.foo.id}"
	  network_interface_id = "${ksyun_instance.foo.network_interface_id}"
	}

```

# Import

nat associate can be imported using the `id`, the id format must be `{nat_id}:{resource_id}`, resource_id range `subnet_id`, `natwork_interface_id` e.g.

## Import Subnet association
```
$ terraform import ksyun_nat_associate.example $nat_id:subnet-$subnet_id
```
## Import NetworkInterface association
```
$ terraform import ksyun_nat_associate.example $nat_id:kni-$network_interface_id
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_interface_id"},
				AtLeastOneOf:  []string{"subnet_id", "network_interface_id"},
				Description:   "The id of the Subnet. Notes: Because of there is one resource in the association, conflict with `network_interface_id`.",
			},
			"network_interface_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IsUUID,
				ConflictsWith: []string{"subnet_id"},
				AtLeastOneOf:  []string{"subnet_id", "network_interface_id"},
				Description:   "The id of network interface that belong to instance. Notes: Because of there is one resource in the association, conflict with `subnet_id`.",
			},
		},
	}
}
func resourceKsyunNatAssociationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	if _, isSubnet := d.GetOk("subnet_id"); isSubnet {
		err = vpcService.CreateNatAssociate(d, resourceKsyunNatAssociation())
		if err != nil {
			return fmt.Errorf("error on creating nat associate %q, %s", d.Id(), err)
		}
	}

	if _, isKni := d.GetOk("network_interface_id"); isKni {
		if err = vpcService.CreateNatInstanceAssociate(d, resourceKsyunNatAssociation()); err != nil {
			return err
		}
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
		return fmt.Errorf("error on deleting nat subnet associate %q, %s", d.Get("subnet_id"), err)
	}

	if err = vpcService.RemoveNatInstanceAssociate(d); err != nil {
		return fmt.Errorf("error on deleting nat network interface associated %q, %s", d.Get("network_interface_id"), err)
	}
	return err
}
