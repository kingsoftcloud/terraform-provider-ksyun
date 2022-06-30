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
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
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
			},

			// openapi已经不支持dhcp的参数，保留这两个值兼容老用户，实际不起作用
			"dhcp_ip_to": {
				Type:         schema.TypeString,
				Deprecated:   "This attribute is deprecated and will be removed in a future version.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIpAddress,
			},
			"dhcp_ip_from": {
				Type:         schema.TypeString,
				Deprecated:   "This attribute is deprecated and will be removed in a future version.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIpAddress,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"gateway_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			"dns1": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
				Computed:     true,
			},

			"dns2": {
				Type:         schema.TypeString,
				ForceNew:     false,
				Optional:     true,
				ValidateFunc: validateIpAddress,
				Computed:     true,
			},
			"network_acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"availability_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_ip_number": {
				Type:     schema.TypeString,
				Computed: true,
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
