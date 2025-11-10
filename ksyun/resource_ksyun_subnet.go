/*
Provides a Subnet resource under VPC resource.

# Example Usage

```hcl

	resource "ksyun_vpc" "example" {
	  vpc_name   = "tf-example-vpc-01"
	  cidr_block = "10.0.0.0/16"
	}

	resource "ksyun_subnet" "example" {
	  subnet_name      = "tf-acc-subnet1"
	  	cidr_block = "10.0.5.0/24"
	      subnet_type = "Normal"
	      dhcp_ip_from = "10.0.5.2"
	      dhcp_ip_to = "10.0.5.253"
	      vpc_id  = "${ksyun_vpc.test.id}"
	      gateway_ip = "10.0.5.1"
	      dns1 = "198.18.254.41"
	      dns2 = "198.18.254.40"
	      availability_zone = "cn-shanghai-2a"
	}

```

# Import

Subnet can be imported using the `id`, e.g.

```
$ terraform import ksyun_subnet.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSubnetCreate,
		Update: resourceKsyunSubnetUpdate,
		Read:   resourceKsyunSubnetRead,
		Delete: resourceKsyunSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "The name of the availability zone.",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the subnet.",
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "The CIDR block assigned to the subnet.",
			},

			"subnet_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Reserve",
					"Normal",
					"Physical",
				}, false),
				Description: "The type of subnet. Valid Values:'Reserve', 'Normal', 'Physical'.",
			},

			// openapi已经不支持dhcp的参数，保留这两个值兼容老用户，实际不起作用
			"dhcp_ip_to": {
				Type:         schema.TypeString,
				Deprecated:   "This attribute is deprecated and will be removed in a future version.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIpAddress,
				Description:  "DHCP end IP.",
			},
			"dhcp_ip_from": {
				Type:         schema.TypeString,
				Deprecated:   "This attribute is deprecated and will be removed in a future version.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIpAddress,
				Description:  "DHCP start IP.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The id of the vpc.",
			},

			"gateway_ip": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "The IP of gateway.",
			},

			"dns1": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
				Computed:     true,
				Description:  "The dns of the subnet.",
			},

			"dns2": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
				Computed:     true,
				Description:  "The dns of the subnet.",
			},
			"provided_ipv6_cidr_block": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "whether support IPV6 CIDR blocks. <br> NOTES: providing a part of regions now.",
			},

			"visit_internet": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("subnet_type").(string) != "Physical" || d.Id() != ""
				},
				Description: "Whether the subnet can access the Internet. Valid, when subnet_type = Physical.",
			},

			"network_acl_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the ACL that the desired Subnet associated to.",
			},

			"nat_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the NAT that the desired Subnet associated to.",
			},

			"availability_zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the availability zone.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the subnet.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the subnet.",
			},
			"available_ip_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "number of available IPs.",
			},
			"ipv6_cidr_block_association_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the Ipv6 of this subnet bound.",
						},
					},
				},
				Description: "An Ipv6 association list of this subnet.",
			},
		},
	}
}

func resourceKsyunSubnetCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSubnet(d, resourceKsyunSubnet())
	if err != nil {
		return fmt.Errorf("error on creating vpc %q, %s", d.Id(), err)
	}
	return resourceKsyunSubnetRead(d, meta)
}

func resourceKsyunSubnetRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSubnet(d, resourceKsyunSubnet())
	if err != nil {
		return fmt.Errorf("error on reading subnet %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSubnetUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifySubnet(d, resourceKsyunSubnet())
	if err != nil {
		return fmt.Errorf("error on updating subnet %q, %s", d.Id(), err)
	}
	return resourceKsyunSubnetRead(d, meta)
}

func resourceKsyunSubnetDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveSubnet(d)
	if err != nil {
		return fmt.Errorf("error on deleting subnet %q, %s", d.Id(), err)
	}
	return err
}
