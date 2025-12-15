/*
Provides a Cloud Firewall Instance resource.

# Example Usage

```hcl
resource "ksyun_kfw_instance" "default" {
  instance_name  = "test-kfw-instance"
  instance_type  = "Advanced"
  bandwidth      = 50
  total_eip_num   = 50
  charge_type     = "Monthly"
  project_id      = "0"
  purchase_time   = 1
}
```

# Import

Cloud Firewall Instance can be imported using the `cfw_instance_id`, e.g.

```
$ terraform import ksyun_kfw_instance.default cfw_instance_id
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunCfwInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunCfwInstanceCreate,
		Read:   resourceKsyunCfwInstanceRead,
		Update: resourceKsyunCfwInstanceUpdate,
		Delete: resourceKsyunCfwInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cfw_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Cloud Firewall Instance.",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
				Description:  "The name of the Cloud Firewall Instance. Length 0-64 characters, supports Chinese, English, numbers.",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Advanced",
					"Enterprise",
				}, false),
				Description: "Instance type. Valid values: Advanced, Enterprise.",
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, 5000),
				Description:  "Bandwidth (10-5000M). Must be a multiple of 5M. Advanced: minimum 10M, Enterprise: minimum 50M.",
			},
			"total_eip_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 500),
				Description:  "Total number of protected IPs. Range: 1-500. Advanced: minimum 20, Enterprise: minimum 50.",
			},
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Monthly",
					"Daily",
				}, false),
				Description: "Billing type. Valid values: Monthly (prepaid), Daily (pay-as-you-go, trial).",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID. Length 0-36 characters, supports letters, numbers, hyphens(-).",
			},
			"purchase_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 36),
				Description:  "Purchase duration. Required when charge_type is Monthly, range: 1-36 months. Required when charge_type is Daily and ProductWhat is 2, range: 1-14 days.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status (1-creating, 2-running, 3-modifying, 4-stopped, 5-abnormal, 6-unsubscribing).",
			},
			"used_eip_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of protected IPs in use.",
			},
			"total_acl_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ACL rules that can be added.",
			},
			"ips_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IPS status (0-stopped, 1-enabled).",
			},
			"av_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "AV status (0-stopped, 1-enabled).",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
		},
	}
}

func resourceKsyunCfwInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kfwService := KfwService{meta.(*KsyunClient)}
	err = kfwService.createKfwInstance(d, resourceKsyunCfwInstance())
	if err != nil {
		return fmt.Errorf("error on creating KFW Instance: %s", err)
	}
	return resourceKsyunCfwInstanceRead(d, meta)
}

func resourceKsyunCfwInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	kfwService := KfwService{meta.(*KsyunClient)}
	err = kfwService.readAndSetKfwInstance(d, resourceKsyunCfwInstance(), false)
	if err != nil {
		return fmt.Errorf("error on reading KFW Instance: %s", err)
	}
	return err
}

func resourceKsyunCfwInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kfwService := KfwService{meta.(*KsyunClient)}
	err = kfwService.modifyKfwInstance(d, resourceKsyunCfwInstance())
	if err != nil {
		return fmt.Errorf("error on updating KFW Instance: %s", err)
	}
	return resourceKsyunCfwInstanceRead(d, meta)
}

func resourceKsyunCfwInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kfwService := KfwService{meta.(*KsyunClient)}
	err = kfwService.removeKfwInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting KFW Instance: %s", err)
	}
	return nil
}
