/*
Provides a ALB rule group resource.

# Example Usage

```hcl

```

# Import

BWS can be imported using the id, e.g.

```
$ terraform import ksyun_alb_rule_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunAlbRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAlbRuleGroupCreate,
		Read:   resourceKsyunAlbRuleGroupRead,
		Update: resourceKsyunAlbRuleGroupUpdate,
		Delete: resourceKsyunAlbRuleGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alb_listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the ALB listener.",
			},
			"alb_rule_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the rule group.",
			},
			"alb_rule_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the ALB rule group.",
			},
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the backend server group.",
			},
			"alb_rule_set": {
				Type:        schema.TypeList,
				MinItems:    1,
				Required:    true,
				Description: "Rule set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alb_rule_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"domain", "url"}, false),
							Description:  "Rule type. valid values: 'domain', 'url'.",
						},
						"alb_rule_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule value.",
						},
					},
				},
			},
			"listener_sync": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether to synchronize the health check, session persistence, and load balancing algorithm of the listener. valid values: 'on', 'off'.",
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RoundRobin",
				ValidateFunc: validation.StringInSlice([]string{
					"RoundRobin",
					"LeastConnections",
				}, false),
				Description: "Forwarding mode of listener. Valid Values:'RoundRobin', 'LeastConnections'.",
			},
			"session_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Description: "The state of session. Valid Values:'start', 'stop'.",
			},
			"session_persistence_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 86400),
				Description:  "Session hold timeout. Valid Values:1-86400.",
			},
			"cookie_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ImplantCookie",
					"RewriteCookie",
				}, false),
				Description: "The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.",
			},
			"cookie_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of cookie.",
			},
			"health_check_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Description: "Status maintained by health examination.Valid Values:'start', 'stop'.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
				Description:  "Interval of health examination.Valid Values:1-3600. Default is 5.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
				Description:  "Health check timeout.Valid Values:1-3600. Default is 4.",
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "Health threshold.Valid Values:1-10. Default is 5.",
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "Unhealthy threshold.Valid Values:1-10. Default is 4.",
			},
			"url_path": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Link to HTTP type listener health check.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The service host name of the health check, which is available only for the HTTP or HTTPS health check.",
			},
		},
	}
}

func resourceKsyunAlbRuleGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbRuleGroup{meta.(*KsyunClient)}
	err = s.CreateAlbRuleGroup(d, resourceKsyunAlbRuleGroup())
	if err != nil {
		return fmt.Errorf("error on creating ALB rule group %q, %s", d.Id(), err)
	}
	return resourceKsyunAlbRuleGroupRead(d, meta)
}
func resourceKsyunAlbRuleGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbRuleGroup{meta.(*KsyunClient)}
	err = s.ReadAndSetRuleGroup(d, resourceKsyunAlbRuleGroup())
	if err != nil {
		return fmt.Errorf("error on reading ALB rule group %q, %s", d.Id(), err)
	}
	return
}
func resourceKsyunAlbRuleGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbRuleGroup{meta.(*KsyunClient)}
	err = s.ModifyRuleGroup(d, resourceKsyunAlbRuleGroup())
	if err != nil {
		return fmt.Errorf("error on updating rule group %q, %s", d.Id(), err)
	}
	err = resourceKsyunAlbRuleGroupRead(d, meta)
	return
}
func resourceKsyunAlbRuleGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbRuleGroup{meta.(*KsyunClient)}
	err = s.RemoveRuleGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting rule group %q, %s", d.Id(), err)
	}
	return
}
