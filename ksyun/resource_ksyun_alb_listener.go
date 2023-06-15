/*
Provides a ALB Listener resource.

# Example Usage

```hcl

	resource "ksyun_alb_listener" "default" {
	}

```

# Import

ALB Listener can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener.example vserver-abcdefg
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunAlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAlbListenerCreate,
		Read:   resourceKsyunAlbListenerRead,
		Update: resourceKsyunAlbListenerUpdate,
		Delete: resourceKsyunAlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the ALB.",
			},
			"alb_listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the listener.",
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
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
				}, false),
				Description: "The protocol of listener. Valid Values: 'HTTP', 'HTTPS'.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The protocol port of listener.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of certificate.",
			},
			"tls_cipher_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TlsCipherPolicy1.0",
					"TlsCipherPolicy1.1",
					"TlsCipherPolicy1.2",
					"TlsCipherPolicy1.2-strict",
					"TlsCipherPolicy1.2-moststrict",
				}, false),
				Description: "TLS cipher policy, valid values:'TlsCipherPolicy1.0','TlsCipherPolicy1.1','TlsCipherPolicy1.2','TlsCipherPolicy1.2-strict','TlsCipherPolicy1.2-moststrict'.",
			},
			"alb_listener_state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"start", "stop"}, false),
				Description:  "The state of listener.Valid Values:'start', 'stop'.",
			},

			"redirect_alb_listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the redirect ALB listener.",
			},

			"session": {
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "session.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_state": {
							Type:     schema.TypeString,
							Optional: true,
							//Default:  "stop",
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
							Description: "The state of session. Valid Values:'start', 'stop'.",
						},
						"session_persistence_period": {
							Type:     schema.TypeInt,
							Optional: true,
							//Default:      3600,
							ValidateFunc: validation.IntBetween(1, 86400),
							Description:  "Session hold timeout. Valid Values:1-86400.",
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							//Default:  "ImplantCookie",
							ValidateFunc: validation.StringInSlice([]string{
								"ImplantCookie",
								"RewriteCookie",
							}, false),
							Description: "The type of cookie, valid values: 'ImplantCookie','RewriteCookie'.",
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							//Computed:    true,
							Description: "The name of cookie.",
						},
					},
				},
				DiffSuppressFunc: AlbListenerDiffSuppressFunc,
			},

			"enable_http2": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "whether enable to HTTP2.",
			},
			"http_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP1.0",
					"HTTP1.1",
				}, false),
				Description: "Backend Protocol, valid values:'HTTP1.0','HTTP1.1'.",
			},

			"alb_listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of listener.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},

			//HealthCheckId

			//"alb_listener_acl_id"
		},
	}
}

func resourceKsyunAlbListenerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerService{meta.(*KsyunClient)}
	err = s.CreateListener(d, resourceKsyunAlbListener())
	if err != nil {
		return fmt.Errorf("error on creating ALB listener %q, %s", d.Id(), err)
	}
	return resourceKsyunAlbListenerRead(d, meta)
}
func resourceKsyunAlbListenerRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerService{meta.(*KsyunClient)}
	err = s.ReadAndSetListener(d, resourceKsyunAlbListener())
	if err != nil {
		return fmt.Errorf("error on reading ALB listener %q, %s", d.Id(), err)
	}
	return
}
func resourceKsyunAlbListenerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerService{meta.(*KsyunClient)}
	err = s.ModifyListener(d, resourceKsyunAlbListener())
	if err != nil {
		return fmt.Errorf("error on updating listener %q, %s", d.Id(), err)
	}
	err = resourceKsyunAlbListenerRead(d, meta)
	return
}
func resourceKsyunAlbListenerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerService{meta.(*KsyunClient)}
	err = s.RemoveListener(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener %q, %s", d.Id(), err)
	}
	return
}
