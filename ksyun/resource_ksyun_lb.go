/*
Provides a Load Balancer resource.

# Example Usage

```hcl

	resource "ksyun_lb" "default" {
	  vpc_id = "74d0a45b-472d-49fc-84ad-221e21ee23aa"
	  load_balancer_name = "tf-xun1"
	  type = "public"
	}

```

# Import

LB can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb.example fdeba8ca-8aa6-4cd0-8ffa-52ca9e9fef42
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunLb() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLbCreate,
		Read:   resourceKsyunLbRead,
		Update: resourceKsyunLbUpdate,
		Delete: resourceKsyunLbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The ID of the VPC linked to the Load Balancers.",
			},
			"load_balancer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the load balancer.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public",
					"internal",
				}, false),
				ForceNew:    true,
				Description: "The type of load balancer.Valid Values:'public', 'internal'.",
			},
			"subnet_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: loadBalancerDiffSuppressFunc,
				Description:      "The id of the subnet.only Internal type is Required.",
			},
			"private_ip_address": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: loadBalancerDiffSuppressFunc,
				Description:      "The internal Load Balancers can set an private ip address in Reserve Subnet.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "ID of the project.",
			},

			"load_balancer_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "start",
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Description: "The Load Balancers state.Valid Values:'start', 'stop'.",
			},

			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
				ValidateFunc: validation.StringInSlice([]string{
					"all",
					"ipv4",
					"ipv6",
				}, false),
				ForceNew:    true,
				Description: "IP version, valid values: 'all', 'ipv4', 'ipv6'.",
			},

			"tags": tagsSchema(),

			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of Public IP. It is `\"\"` if `internal` is `true`.",
			},

			"load_balancer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the LB.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "associate or disassociate.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of creation for load balancer.",
			},
			"is_waf": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether it is a waf LB or not.",
			},
			"access_logs_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Default is `false`, Setting the value to `true` to enable the service.",
			},
			"access_logs_s3_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Bucket for storing access logs.",
			},
		},
	}
}
func resourceKsyunLbCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLoadBalancer(d, resourceKsyunLb())
	if err != nil {
		return fmt.Errorf("error on creating lb %q, %s", d.Id(), err)
	}
	return resourceKsyunLbRead(d, meta)
}

func resourceKsyunLbRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLoadBalancer(d, resourceKsyunLb())
	if err != nil {
		return fmt.Errorf("error on reading lb %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunLbUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyLoadBalancer(d, resourceKsyunLb())
	if err != nil {
		return fmt.Errorf("error on updating lb %q, %s", d.Id(), err)
	}
	return resourceKsyunLbRead(d, meta)
}

func resourceKsyunLbDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLoadBalancer(d)
	if err != nil {
		return fmt.Errorf("error on deleting lb %q, %s", d.Id(), err)
	}
	return err
}
