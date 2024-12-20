/*
Provides a ALB rule group resource.

# Example Usage

```hcl
# network and security group configuration
resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "test" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}

resource "ksyun_security_group" "test" {
  vpc_id              = ksyun_vpc.test.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "test" {
  security_group_id = ksyun_security_group.test.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}

# --------------------------------------------------------
# alb-rule-group relational configuration

# ksyun alb configuration
resource "ksyun_alb" "test" {
  alb_name    = "tf-alb-test"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.test.id
  project_id  = 0
  enabled_log = false
  ip_version  = "ipv4"
}

# ksyun alb listener configuration
resource "ksyun_alb_listener" "test" {
  alb_id             = ksyun_alb.test.id
  alb_listener_name  = "alb-test-listener"
  protocol           = "HTTP"
  port               = 8088
  alb_listener_state = "start"
  http_protocol      = "HTTP1.1"
  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }
}

# --------------------------------------------
# backend server group and kec instance configuration
# backend server group configuration
resource "ksyun_lb_backend_server_group" "test" {
  backend_server_group_name = "tf_bsg"
  vpc_id                    = ksyun_vpc.test.id
  backend_server_group_type = "Server"
}
resource "ksyun_lb_register_backend_server" "default" {
  backend_server_group_id=ksyun_lb_backend_server_group.test.id
  backend_server_ip= ksyun_instance.test.0.private_ip_address
  backend_server_port=8090
  weight=10
}
resource "ksyun_lb_register_backend_server" "default2" {
  backend_server_group_id=ksyun_lb_backend_server_group.test.id
  backend_server_ip= ksyun_instance.test.1.private_ip_address
  backend_server_port=8090
  weight=10
}

# kec instance creating
data "ksyun_images" "default" {
  output_file  = "output_result"
  name_regex   = "centos-7.0"
  is_public    = true
  image_source = "system"
}

data "ksyun_ssh_keys" "test" {

}

resource "ksyun_instance" "test" {
  count             = 2
  security_group_id = [
    ksyun_security_group.test.id
  ]
  subnet_id = ksyun_subnet.test.id
  key_id    = [data.ksyun_ssh_keys.test.keys.0.key_id]

  instance_type = "S6.1A"
  charge_type   = "Daily"
  instance_name = "tf-alb-test-vm"
  project_id    = 0

  image_id = data.ksyun_images.default.images.0.image_id
}


# ksyun_alb_rule_group configuration
resource "ksyun_alb_rule_group" "default" {
  alb_listener_id         = ksyun_alb_listener.test.id
  alb_rule_group_name     = "tf_alb_rule_group"
  backend_server_group_id = ksyun_lb_backend_server_group.test.id
  alb_rule_set {
    alb_rule_type  = "url"
    alb_rule_value = "/test/path"
  }
  alb_rule_set {
    alb_rule_type  = "domain"
    alb_rule_value = "www.ksyun.com"
  }
  listener_sync = "on"
}
```

# Import

`ksyun_alb_rule_group` can be imported using the id, e.g.

```
$ terraform import ksyun_alb_rule_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var hqcValueSchema = map[string]*schema.Schema{
	"key": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The key of querying.",
	},
	"value": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "The value of querying.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
}

var rewriteConfigSchema = map[string]*schema.Schema{
	"http_host": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The host of the rewrite.",
	},
	"url": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The url of the rewrite.",
	},
	"query_string": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The query string of the rewrite.",
	},
}

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

			"alb_rule_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the ALB rule group.",
			},
			"backend_server_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: albRuleGroupTypeDiffSuppressFunc,

				ConflictsWith: []string{"redirect_alb_listener_id", "fixed_response_config"},
				AtLeastOneOf:  []string{"backend_server_group_id", "redirect_alb_listener_id", "fixed_response_config"},
				Description:   "The ID of the backend server group. Conflict with 'backend_server_group_id' and 'fixed_response_config'.",
			},
			"alb_rule_set": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule set, define strategies for being load-balance of backend server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alb_rule_type": {
							Type:     schema.TypeString,
							Required: true,
							// ValidateFunc: validation.StringInSlice([]string{"domain", "url"}, false),
							Description: "Rule type. valid values: (domain|url|method|sourceIp|header|query|cookie).",
						},
						"alb_rule_value": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: albRuleGroupDiffSuppressFunc,
							Description:      "Rule value. It works, when `alb_rule_type` is domain or url.",
						},
						"method_value": {
							Type:             schema.TypeList,
							Optional:         true,
							Description:      "The method value of the rule. It works, when `alb_rule_type` is method.",
							DiffSuppressFunc: albRuleGroupDiffSuppressFunc,
							Elem:             &schema.Schema{Type: schema.TypeString},
						},
						"source_ip_value": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: albRuleGroupDiffSuppressFunc,
							Elem:             &schema.Schema{Type: schema.TypeString},
							Description:      "The source ip value of the rule. It works, when `alb_rule_type` is sourceIp.",
						},
						"header_value": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: albRuleGroupDiffSuppressFunc,
							Description:      "The header value of the rule. It works, when `alb_rule_type` is header.",
							Elem:             &schema.Resource{Schema: hqcValueSchema},
						},
						"cookie_value": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: albRuleGroupDiffSuppressFunc,
							Description:      "The cookie value of the rule. It works, when `alb_rule_type` is cookie.",
							Elem:             &schema.Resource{Schema: hqcValueSchema},
						},
						"query_value": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: albRuleGroupDiffSuppressFunc,
							Description:      "The query value of the rule. It works, when `alb_rule_type` is query.",
							Elem:             &schema.Resource{Schema: hqcValueSchema},
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
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "The state of session. Valid Values:'start', 'stop'. Should set it value, when `listener_sync` is off.",
			},
			"session_persistence_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 86400),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Session hold timeout. Valid Values:1-86400. Should set it value, when `listener_sync` is off.",
			},
			"cookie_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ImplantCookie",
					"RewriteCookie",
				}, false),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.",
			},
			"cookie_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The name of cookie. Should set it value, when `listener_sync` is off and `cookie_type` is `RewriteCookie`.",
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
			},
			"health_check_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Status maintained by health examination.Valid Values:'start', 'stop'. Should set it value, when `listener_sync` is off.",
			},

			"http_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The http requests' method. Valid Value: GET|HEAD. It works, when `health_protocol` is HTTP.",
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD"},
					false),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
			},

			"health_port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "The port of connecting for health check. It works, when `listener_sync` is off.",
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
			},

			"health_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol of connecting for health check. It works, when `listener_sync` is off.",
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "TCP"},
					false),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
			},

			"interval": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 1000),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Interval of health examination.Valid Values:1-3600. Should set it value, when `listener_sync` is off.",
			},
			"timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 3600),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Health check timeout.Valid Values:1-3600. Should set it value, when `listener_sync` is off.",
			},
			"healthy_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 10),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Health threshold.Valid Values:1-10. Should set it value, when `listener_sync` is off.",
			},
			"unhealthy_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.IntBetween(1, 10),
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Unhealthy threshold.Valid Values:1-10. Should set it value, when `listener_sync` is off.",
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:         true,
				Default:          "/",
				ValidateFunc:     validation.StringIsNotEmpty,
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "Link to HTTP type listener health check. Should set it value, when `listener_sync` is off.",
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				// Computed:         true,
				Default:          "",
				DiffSuppressFunc: AlbRuleGroupSyncOffDiffSuppressFunc,
				Description:      "The service host name of the health check, which is available only for the HTTP or HTTPS health check. Should set it value, when `listener_sync` is off.",
			},

			"redirect_alb_listener_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: albRuleGroupTypeDiffSuppressFunc,

				ConflictsWith: []string{"backend_server_group_id", "fixed_response_config"},
				AtLeastOneOf:  []string{"backend_server_group_id", "redirect_alb_listener_id", "fixed_response_config"},
				Description:   "The id of redirect alb listener. Conflict with 'backend_server_group_id' and 'fixed_response_config'.",
			},
			"redirect_http_code": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "301",
				DiffSuppressFunc: albRuleGroupTypeDiffSuppressFunc,
				Description:      "The http code of redirecting. Valid Values: 301|302|307.",
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				// Default:  "ForwardGroup",
				Description: "The type of rule group, Valid Values: ForwardGroup|Redirect|FixedResponse. Default: ForwardGroup. \n" +
					"**Notes**: The type is supposed to be of consistency with backend instance. `ForwardGroup -> backend_server_group_id`," +
					" `Redirect -> redirect_alb_listener_id`, `FixedResponse -> fixed_response_config`.",
			},

			"fixed_response_config": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: albRuleGroupTypeDiffSuppressFunc,
				ConflictsWith:    []string{"backend_server_group_id", "redirect_alb_listener_id"},
				AtLeastOneOf:     []string{"backend_server_group_id", "redirect_alb_listener_id", "fixed_response_config"},
				Description:      "The config of fixed response. Conflict with 'backend_server_group_id' and 'fixed_response_config'.",
				Elem:             fixedResponseConfigResourceElem(),
			},

			"alb_rule_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the rule group.",
			},
			"rewrite_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: rewriteConfigSchema,
				},
				DiffSuppressFunc: albRuleGroupTypeDiffSuppressFunc,
				Description:      "The config of rewrite.",
			},
		},
	}
}

func resourceKsyunAlbRuleGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	if indispensableErr := checkIndispensableParams(d); indispensableErr != nil {
		return indispensableErr
	}

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
	// check indispensable parameters such as session_state, health_check_state
	if indispensableErr := checkIndispensableParams(d); indispensableErr != nil {
		return indispensableErr
	}

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

func checkIndispensableParams(d *schema.ResourceData) error {

	errFormat := "`%s` cannot be blank, when `listener_sync` is off and `%s` is start. Should be set it value"
	if d.Get("listener_sync").(string) == "off" {
		if state, ok := d.GetOk("session_state"); !ok {
			return fmt.Errorf("`session_state` cannot be blank, when `listener_sync` is off. Should be set it value")
		} else if state == "start" {
			sessionIndispensables := albRuleGroupSessionNecessary
			var errAssembly []string
			for _, sI := range sessionIndispensables {
				if _, ok := d.GetOk(sI); !ok {
					errAssembly = append(errAssembly, sI)
				}
			}
			if errAssembly != nil && len(errAssembly) > 0 {
				return fmt.Errorf(errFormat, strings.Join(errAssembly, ", "), "session_state")
			}

			if d.Get("cookie_type").(string) == "RewriteCookie" {
				if _, ok := d.GetOk("cookie_name"); !ok {
					return fmt.Errorf("`cookie_name` cannot be blank, when `listener_sync` is off and `cookie_type` is RewriteCookie. Should be set it value")
				}
			}
		}
		if state, ok := d.GetOk("health_check_state"); !ok {
			return fmt.Errorf("`health_check_state` cannot be blank, when `listener_sync` is off. Should be set it value")
		} else if state == "start" {
			healthIndispensables := albRuleGroupHealthNecessary
			var errAssembly []string
			for _, hI := range healthIndispensables {
				if stringSliceContains([]string{"url_path", "host_name"}, hI) {
					continue
				}
				if _, ok := d.GetOk(hI); !ok {
					errAssembly = append(errAssembly, hI)
				}
			}
			if errAssembly != nil && len(errAssembly) > 0 {
				return fmt.Errorf(errFormat, strings.Join(errAssembly, ", "), "health_check_state")
			}
		}
	}
	return nil
}
