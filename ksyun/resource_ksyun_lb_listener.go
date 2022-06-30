package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunListener() *schema.Resource {
	entry := resourceKsyunHealthCheck().Schema
	for k, v := range entry {
		if k == "listener_id" {
			v.Required = false
			v.Computed = true
		} else {
			v.ForceNew = false
			v.DiffSuppressFunc = nil
		}
	}
	return &schema.Resource{
		Create: resourceKsyunListenerCreate,
		Read:   resourceKsyunListenerRead,
		Update: resourceKsyunListenerUpdate,
		Delete: resourceKsyunListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "start",
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
			},
			"listener_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TCP",
				ValidateFunc: validation.StringInSlice([]string{
					"TCP",
					"UDP",
					"HTTP",
					"HTTPS",
				}, false),
				ForceNew: true,
			},
			"certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RoundRobin",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RoundRobin",
					"LeastConnections",
					"MasterSlave",
					"QUIC_CID",
				}, false),
			},

			"enable_http2": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"tls_cipher_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TlsCipherPolicy1.0",
				ValidateFunc: validation.StringInSlice([]string{
					"TlsCipherPolicy1.0",
					"TlsCipherPolicy1.1",
					"TlsCipherPolicy1.2",
					"TlsCipherPolicy1.2-strict",
				}, false),
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},
			"http_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HTTP1.1",
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP1.0",
					"HTTP1.1",
				}, false),
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"redirect_listener_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"session": {
				Type:     schema.TypeList,
				MaxItems: 1,
				MinItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_state": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "stop",
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
						},
						"session_persistence_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(1, 86400),
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ImplantCookie",
							ValidateFunc: validation.StringInSlice([]string{
								"ImplantCookie",
								"RewriteCookie",
							}, false),
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: entry,
				},
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
				//Computed: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunListenerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateListener(d, resourceKsyunListener())
	if err != nil {
		return fmt.Errorf("error on creating listener %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerRead(d, meta)
}

func resourceKsyunListenerRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetListener(d, resourceKsyunListener())
	if err != nil {
		return fmt.Errorf("error on reading listener %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunListenerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyListener(d, resourceKsyunListener())
	if err != nil {
		return fmt.Errorf("error on updating listener %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerRead(d, meta)
}

func resourceKsyunListenerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveListener(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener %q, %s", d.Id(), err)
	}
	return err
}
