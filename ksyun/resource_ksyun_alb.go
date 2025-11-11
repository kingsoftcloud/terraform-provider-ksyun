/*
Provides a ALB resource.

# Example Usage

```hcl

	resource "ksyun_vpc" "default" {
	  vpc_name = "tf_alb_test_vpc"
	  cidr_block = "10.0.0.0/16"
	}

	resource "ksyun_alb" "default" {
	  alb_name = "tf_test_alb1"
	  alb_version = "standard"
	  alb_type = "public"
	  state = "start"
	  charge_type = "PrePaidByHourUsage"
	  vpc_id = ksyun_vpc.default.id
	  project_id = 0
	}

	data "ksyun_lines" "default" {
	  output_file="output_result1"
	  line_name="BGP"
	}

	resource "ksyun_eip" "foo" {
	  line_id =data.ksyun_lines.default.lines.0.line_id
	  band_width =1
	  charge_type = "PostPaidByDay"
	  purchase_time =1
	  project_id=0
	}

	resource "ksyun_eip_associate" "eip_bind" {
	  allocation_id = ksyun_eip.foo.id
	  instance_id   = ksyun_alb.foo.id
	  instance_type = "Slb"
	}

```

# Import

`ksyun_alb` can be imported using the id, e.g.

```
$ terraform import ksyun_alb.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunAlb() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAlbCreate,
		Read:   resourceKsyunAlbRead,
		Update: resourceKsyunAlbUpdate,
		Delete: resourceKsyunAlbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the ALB.",
			},
			"alb_version": {
				Type: schema.TypeString,
				// Optional:     true,
				// Computed:     true,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "medium", "advanced"}, false),
				Description:  "The version of the ALB. valid values:'standard', 'medium', 'advanced'.",
			},
			"alb_type": {
				Type: schema.TypeString,
				// Optional:     true,
				// Computed:     true,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "internal"}, false),
				Description:  "The type of the ALB, valid values:'public', 'internal'.",
			},

			// 2025-06-27
			"protocol_layers": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The protocol layers of the ALB, valid values: 'L4', 'L7', 'L4-L7'.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC.",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "IP version, 'ipv4' or 'ipv6'.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the project.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PrePaidByHourUsage"}, false),
				Description:  "The charge type, valid values: 'PrePaidByHourUsage'.",
			},

			"subnet_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: albInternalDiffSuppressFunc,
				Description:      "The Id of Subnet that's type is **Reserve**. It not be empty, when 'alb_type' as '**internal**'.",
			},

			"private_ip_address": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateFunc:     validation.IsIPAddress,
				DiffSuppressFunc: albInternalDiffSuppressFunc,
				Description:      "The private ip address. It not be empty, when 'alb_type' as '**internal**'.",
			},
			"enabled_quic": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable quic.",
			},

			"enable_hpa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable hpa.",
			},

			"delete_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Description:  "Whether delete protection is enabled or not. Values: `off` or `on`.",
			},

			"modification_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Description:  "Whether modification protection is enabled or not. Values: `off` or `on`.",
			},

			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The state of the ALB, Valid Values:'start', 'stop'.",
			},

			"enabled_log": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether log is enabled or not. Specific `klog_info` field when `enabled_log` is true.",
			},
			"klog_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Indicate klog info, including log-project-name and log-pool-name, that use to bind log service for this alb process.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "account id.",
						},
						"log_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "log pool name.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "log project name.",
						},
					},
				},
			},

			"tags": tagsSchema(),

			// computed values
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IP address.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The status of the ALB.",
			},
		},
	}
}

func resourceKsyunAlbCreate(d *schema.ResourceData, meta interface{}) (err error) {
	switch d.Get("alb_type") {
	case "internal":
		if !(d.HasChange("subnet_id") && d.HasChange("private_ip_address")) {
			return fmt.Errorf("subnet_id and private_ip_address must not be empty, when alb_type as internal")
		}
	}
	s := AlbService{meta.(*KsyunClient)}
	err = s.CreateAlb(d, resourceKsyunAlb())
	if err != nil {
		return
	}
	return resourceKsyunAlbRead(d, meta)
}

func resourceKsyunAlbRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbService{meta.(*KsyunClient)}
	err = s.ReadAndSetAlb(d, resourceKsyunAlb())
	if err != nil {
		return fmt.Errorf("error on reading ALB %q, %s", d.Id(), err)
	}
	return
}

func resourceKsyunAlbUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbService{meta.(*KsyunClient)}
	err = s.ModifyAlb(d, resourceKsyunAlb())
	if err != nil {
		return
	}
	return resourceKsyunAlbRead(d, meta)
}

func resourceKsyunAlbDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbService{meta.(*KsyunClient)}
	if d.Get("delete_protection") == "on" {
		return fmt.Errorf("ALB %q is protected from deletion, if you want to delete it to set `delete_protection` as off", d.Id())
	}
	err = s.RemoveAlb(d)
	if err != nil {
		return fmt.Errorf("error on deleting ALB %q, %s", d.Id(), err)
	}
	return err
}
