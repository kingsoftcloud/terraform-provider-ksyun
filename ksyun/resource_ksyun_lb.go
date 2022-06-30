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
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ValidateFunc: validation.StringInSlice([]string{
					"public",
					"internal",
				}, false),
				ForceNew: true,
			},
			"subnet_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: loadBalancerDiffSuppressFunc,
			},
			"private_ip_address": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: loadBalancerDiffSuppressFunc,
			},

			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"load_balancer_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "start",
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
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
				ForceNew: true,
			},

			"tags": tagsSchema(),

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"load_balancer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_waf": {
				Type:     schema.TypeBool,
				Computed: true,
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
