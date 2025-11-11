/*
Provides a Monitor Alarm Policy resource.

# Example Usage

```hcl

	resource "ksyun_monitor_alarm_policy" "foo" {
	  policy_name        = "tf-test-policy"
	  product_type       = 0
	  policy_type        = 0
	  resource_bind_type = 3

	  trigger_rules {
	    compare       = ">"
		effect_bt     = "00:00"
		effect_et     = "23:59"
	    interval      = 5
		item_key      = "cpu.utilizition.total"
	    item_name     = "CPU利用率"
		max_count     = 3
	    method        = "avg"
	    period        = "5m"
	    points        = 2
	    trigger_value = "90"
	    units         = "%"
	  }

	  instance_ids = []

	  user_notice {
	    contact_way  = 2
	    contact_flag = 2
	    contact_id   = 4423
	  }

	  url_notice = ["https://example.com/webhook"]
	}

```

# Import

Monitor Alarm Policy can be imported using the `policy_id`, e.g.

```
$ terraform import ksyun_monitor_alarm_policy.foo policy_id
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunMonitorAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunMonitorAlarmPolicyCreate,
		Read:   resourceKsyunMonitorAlarmPolicyRead,
		Delete: resourceKsyunMonitorAlarmPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the alarm policy.",
			},
			"product_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Cloud service category, see [Cloud Service Category](https://docs.ksyun.com/documents/42040).",
			},
			"policy_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Policy type, 0: Normal policy, 1: Default policy.",
			},
			"resource_bind_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Policy associated resource type: 1: All, 2: Project group, 3: Selected instances.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Project group ID, required when ResourceBindType=2.",
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Instance ID list, required when ResourceBindType=3.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"trigger_rules": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Trigger alarm rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Statistical period, unit: minutes; e.g., 1 minute, 5 minutes, 1 hour correspond to 1m, 5m, 60m respectively.",
						},
						"method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Statistical method, only supports avg, max, min, sum.",
						},
						"compare": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Comparison method, supports >, <, =.",
						},
						"trigger_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger alarm threshold.",
						},
						"item_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Monitoring item name, e.g., CPU utilization.",
						},
						"item_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Monitoring item, e.g., cpu.utilizition.total.",
						},
						"units": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Monitoring item unit, e.g., %.",
						},
						"points": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Consecutive periods.",
						},
						"effect_bt": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Policy effective start time, e.g., 12:00.",
						},
						"effect_et": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Policy effective end time, e.g., 12:00.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm interval, unit: minutes, e.g., 5, 10, 30.",
						},
						"max_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of alarm notifications, value range: 1~5.",
						},
					},
				},
			},
			"user_notice": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Alarm receiving methods.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact_way": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm receiving method, 1: Send email, 2: Send SMS, 3: Send email and SMS.",
						},
						"contact_flag": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm contact type, 1: Contact group, 2: Contact person.",
						},
						"contact_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Contact ID or contact group ID.",
						},
					},
				},
			},
			"url_notice": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    5,
				Description: "Alarm callback Webhook addresses, up to 5 can be added.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the alarm policy.",
			},
		},
	}
}

func resourceKsyunMonitorAlarmPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	monitorService := MonitorService{meta.(*KsyunClient)}
	err := monitorService.CreateAlarmPolicy(d, resourceKsyunMonitorAlarmPolicy())
	if err != nil {
		return fmt.Errorf("error on creating Monitor Alarm Policy %q, %s", d.Id(), err)
	}
	return resourceKsyunMonitorAlarmPolicyRead(d, meta)
}

func resourceKsyunMonitorAlarmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	monitorService := MonitorService{meta.(*KsyunClient)}
	err := monitorService.ReadAndSetAlarmPolicy(d, resourceKsyunMonitorAlarmPolicy())
	if err != nil {
		return fmt.Errorf("error on reading Monitor Alarm Policy, %s", err)
	}
	return nil
}

func resourceKsyunMonitorAlarmPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	monitorService := MonitorService{meta.(*KsyunClient)}
	err := monitorService.DeleteAlarmPolicy(d)
	if err != nil {
		return fmt.Errorf("error on deleting Monitor Alarm Policy %q, %s", d.Id(), err)
	}
	return nil
}
