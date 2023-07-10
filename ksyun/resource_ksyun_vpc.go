/*
Provides a VPC resource.

~> **Note**  The network segment can only be created or deleted, can not perform both of them at the same time.

Example Usage

```hcl
resource "ksyun_vpc" "example" {
  vpc_name   = "ksyun_vpc_tf"
  cidr_block = "10.1.0.2/24"
}

Import

VPC can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpc.example vpc-abc123456
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpcCreate,
		Update: resourceKsyunVpcUpdate,
		Read:   resourceKsyunVpcRead,
		Delete: resourceKsyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the vpc.",
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "The CIDR blocks of VPC.",
			},

			"provided_ipv6_cidr_block": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "whether support IPV6 CIDR blocks. <br> NOTES: providing a part of regions now.",
			},

			"is_default": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Default:     false,
				Optional:    true,
				Description: "Whether the VPC is default or not.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of creation for VPC.",
			},
			"ipv6_cidr_block_association_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the Ipv6 of this vpc bound.",
						},
					},
				},
				Description: "An Ipv6 association list of this vpc.",
			},
		},
	}
}

func resourceKsyunVpcCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpc(d, resourceKsyunVpc())
	if err != nil {
		return fmt.Errorf("error on creating vpc %q, %s", d.Id(), err)
	}
	return resourceKsyunVpcRead(d, meta)
}

func resourceKsyunVpcRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpc(d, resourceKsyunVpc())
	if err != nil {
		return fmt.Errorf("error on reading vpc %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpcUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpc(d, resourceKsyunVpc())
	if err != nil {
		return fmt.Errorf("error on updating vpc %q, %s", d.Id(), err)
	}
	return resourceKsyunVpcRead(d, meta)
}

func resourceKsyunVpcDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpc(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpc %q, %s", d.Id(), err)
	}
	return err
}
