/*
Provides a lb rule resource.

# Example Usage

```hcl

	resource "ksyun_lb_rule" "default" {
		path = "/tfxun/update",
		host_header_id = "",
		backend_server_group_id=""
		listener_sync="on"
		method="RoundRobin"
		session {
			session_state = "start"
			session_persistence_period = 1000
			cookie_type = "ImplantCookie"
			cookie_name = "cookiexunqq"
		}
		health_check{
			health_check_state = "start"
			healthy_threshold = 2
			interval = 200
			timeout = 2000
			unhealthy_threshold = 2
			url_path = "/monitor"
			host_name = "www.ksyun.com"
		}
	}

```

# Import

LB Rule can be imported using the `id`, e.g.

```
$ terraform import ksyun_lb_rule.example vserver-abcdefg
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunSlbRule() *schema.Resource {
	entry := resourceKsyunHealthCheck().Schema
	for k, v := range entry {
		if k == "listener_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
			v.DiffSuppressFunc = nil
		}

		switch k {
		case "http_method":
			v.Optional = false
			v.Computed = true
			v.ValidateFunc = nil
		case "lb_type":
			delete(entry, k)
		}
	}
	return &schema.Resource{
		Create: resourceKsyunSlbRuleCreate,
		Read:   resourceKsyunSlbRuleRead,
		Update: resourceKsyunSlbRuleUpdate,
		Delete: resourceKsyunSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path of rule.",
			},
			"host_header_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of host header id.",
			},
			"backend_server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of backend server group.",
			},
			"listener_sync": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"on",
					"off",
				}, false),
				Default:     "on",
				Description: "Whether to synchronizethe the health check, the session hold and the forward algorithms of the listener.Valid Values:'on', 'off'. Default is 'on'.",
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RoundRobin",
					"LeastConnections",
				}, false),
				Default:          "RoundRobin",
				DiffSuppressFunc: lbRuleDiffSuppressFunc,
				Description:      "Forwarding mode of listener.Valid Values:'RoundRobin', 'LeastConnections'. Default is 'RoundRobin'.",
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbRuleDiffSuppressFunc,
				Description:      "health check configuration.",
			},
			"session": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Session.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
							Default:     "start",
							Description: "The state of session.Valid Values:'start', 'stop'. Default is 'start'.",
						},
						"session_persistence_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 86400),
							Default:      7200,
							Description:  "Session hold timeout.Valid Values:1-86400. Default is '7200'.",
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ImplantCookie",
								"RewriteCookie",
							}, false),
							Default:     "ImplantCookie",
							Description: "The type of the cookie.Valid Values:'ImplantCookie', 'RewriteCookie'. Default is 'ImplantCookie'.",
						},
						"cookie_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of cookie.The CookieType is valid and required when it is 'RewriteCookie'; otherwise, this value is ignored.",
						},
					},
				},
				DiffSuppressFunc: lbRuleDiffSuppressFunc,
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation time.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the rule.",
			},
		},
	}
}
func resourceKsyunSlbRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLbRule(d, resourceKsyunSlbRule())
	if err != nil {
		return fmt.Errorf("error on creating lb rule %q, %s", d.Id(), err)
	}
	return resourceKsyunSlbRuleRead(d, meta)
}

func resourceKsyunSlbRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLbRule(d, resourceKsyunSlbRule())
	if err != nil {
		return fmt.Errorf("error on reading lb rule %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSlbRuleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyLbRule(d, resourceKsyunSlbRule())
	if err != nil {
		return fmt.Errorf("error on updating lb rule %q, %s", d.Id(), err)
	}
	return resourceKsyunSlbRuleRead(d, meta)
}

func resourceKsyunSlbRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLbRule(d)
	if err != nil {
		return fmt.Errorf("error on deleting lb rule %q, %s", d.Id(), err)
	}
	return err
}
