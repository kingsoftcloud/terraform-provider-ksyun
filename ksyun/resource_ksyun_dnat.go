/*
Provides a Dnat resource under VPC resource.

# Example Usage

```hcl

data "ksyun_images" "centos-7_5" {
  platform= "centos-7.5"
}
data "ksyun_availability_zones" "default" {
}

resource "ksyun_security_group" "default" {
  vpc_id = "${ksyun_vpc.foo.id}"
  security_group_name="ksyun-security-group-nat"
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
  project_id=100012
}

resource "ksyun_nat" "foo" {
  nat_name = "ksyun-nat-tf"
  nat_mode = "Subnet"
  nat_type = "public"
  band_width = 1
  charge_type = "DailyPaidByTransfer"
  vpc_id = "${ksyun_vpc.foo.id}"
  nat_ip_number = 3
}
resource "ksyun_vpc" "foo" {
	vpc_name        = "tf-vpc-nat"
	cidr_block = "10.0.5.0/24"
}

resource "ksyun_subnet" "foo" {
  subnet_name      = "tf-acc-nat-subnet1"
  cidr_block = "10.0.5.0/24"
  subnet_type = "Normal"
  vpc_id  = "${ksyun_vpc.foo.id}"
  gateway_ip = "10.0.5.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${data.ksyun_availability_zones.default.availability_zones.0.availability_zone_name}"
}

resource "ksyun_dnat" "foo" {
	nat_id = ksyun_nat.foo.id
	dnat_name = "tf_dnat"
	ip_protocol = "Any"
	nat_ip =ksyun_nat.foo.nat_ip_set[0].nat_ip
	private_ip_address = ksyun_instance.foo.private_ip_address
	description = "test"
	private_port = "Any"
	public_port = "Any"
}
```

# Import

dnat rules can be imported using the `id`, e.g.

```
$ terraform import ksyun_dnat.foo $dnat_id
```
*/

package ksyun

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunDnat() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDnatCreate,
		Read:   resourceKsyunDnatRead,
		Delete: resourceKsyunDnatDelete,
		Update: resourceKsyunDnatUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the Nat.",
			},

			"dnat_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the name of dnat rule.",
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Any", "TCP", "UDP"}, false),
				Description:  "the protocol of dnat port, Valid Options: `Any`, `TCP` and `UDP`. <br> Notes: `public_port` and `private_port` must be set as `Any`, when `ip_protocol` is `Any`. Instead, you should set ports.",
			},
			"nat_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "the nat ip of nat.",
			},
			"public_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Any",
				ValidateFunc: validatePort,
				Description:  "the public port of the internet.",
			},
			"private_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "the private ip of instance in the identical vpc.",
			},
			"private_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Any",
				ValidateFunc: validatePort,
				Description:  "the private port that be accessed in vpc.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the description of this dnat rule.",
			},

			"dnat_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the Subnet. Notes: Because of there is one resource in the association, conflict with `network_interface_id`.",
			},
		},
	}
}
func resourceKsyunDnatCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{client: meta.(*KsyunClient)}
	ipProtocol := d.Get("ip_protocol")
	if !reflect.DeepEqual(ipProtocol, AnyPortType) {
		pubPort, pubPortExist := d.GetOk("public_port")
		infraPort, infraPortExist := d.GetOk("private_port")
		if !pubPortExist || !infraPortExist {
			return fmt.Errorf("`public_port` and `private_port` must be set, when `ip_protocol` is not `Any`")
		}
		if reflect.DeepEqual(pubPort, AnyPortType) || reflect.DeepEqual(infraPort, AnyPortType) {
			return fmt.Errorf("the value of `public_port` and `private_port` must be set, when `ip_protocol` is not `Any`. e.g. hcl: public_port=6379")
		}
	}
	if err = vpcService.CreateDnat(d, resourceKsyunDnat()); err != nil {
		return fmt.Errorf("error on creating dnat %s", err)
	}
	return resourceKsyunDnatRead(d, meta)
}
func resourceKsyunDnatUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{client: meta.(*KsyunClient)}
	if err = vpcService.ModifyDnat(d, resourceKsyunDnat()); err != nil {
		return fmt.Errorf("error on updating dnat %s", err)
	}
	return resourceKsyunDnatRead(d, meta)
}

func resourceKsyunDnatRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetDnats(d, resourceKsyunDnat())
	if err != nil {
		return fmt.Errorf("error on reading  dnat %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunDnatDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveDnat(d, resourceKsyunDnat())
	if err != nil {
		return fmt.Errorf("error on deleting nat subnet associate %q, %s", d.Get("subnet_id"), err)
	}

	return err
}
