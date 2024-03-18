/*
Provides an Health Check resource.

# Example Usage

```hcl

	resource "ksyun_healthcheck" "default" {
	  listener_id = "537e2e7b-0007-4a75-9749-882167dbc93d"
	  health_check_state = "stop"
	  healthy_threshold = 2
	  interval = 20
	  timeout = 200
	  unhealthy_threshold = 2
	  url_path = "/monitor"
	  is_default_host_name = true
	  host_name = "www.ksyun.com"
	}

```

# Import

HealthCheck can be imported using the id, e.g.

```
$ terraform import ksyun_healthcheck.default ${lb_type}:67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunHealthCheckCreate,
		Read:   resourceKsyunHealthCheckRead,
		Update: resourceKsyunHealthCheckUpdate,
		Delete: resourceKsyunHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			State: importHealthcheck,
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the listener.",
			},
			"health_check_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Default:     "start",
				Description: "Status maintained by health examination.Valid Values:'start', 'stop'.",
			},

			"health_check_connect_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of connecting for health check.",
			},

			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Default:      5,
				Description:  "Health threshold.Valid Values:1-10. Default is 5.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
				Default:      5,
				Description:  "Interval of health examination.Valid Values:1-3600. Default is 5.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
				Default:      4,
				Description:  "Health check timeout.Valid Values:1-3600. Default is 4.",
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Default:      4,
				Description:  "Unhealthy threshold.Valid Values:1-10. Default is 4.",
			},
			"url_path": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: lbHealthCheckDiffSuppressFunc,
				ValidateFunc:     validation.StringIsNotEmpty,
				Description:      "Link to HTTP type listener health check.",
			},
			"is_default_host_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the host name is default or not.",
			},
			"host_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbHealthCheckDiffSuppressFunc,
				Description:      "The service host name of the health check, which is available only for the HTTP or HTTPS health check.",
			},

			"http_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The http requests' method. Valid Value: GET|HEAD.",
				ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD"},
					false),
				Computed: true,
			},

			"lb_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Slb",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"Alb", "Slb"},
					false),
				Description: "The type of listener. Valid Value: `Alb` and `Slb`. Default: `Slb`.",
			},

			"health_check_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the health check.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "protocol of the listener.",
			},
		},
	}
}
func resourceKsyunHealthCheckCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateHealthCheck(d, resourceKsyunHealthCheck())
	if err != nil {
		return fmt.Errorf("error on creating health check %q, %s", d.Id(), err)
	}
	return resourceKsyunHealthCheckRead(d, meta)
}

func resourceKsyunHealthCheckRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetHealCheck(d, resourceKsyunHealthCheck())
	if err != nil {
		return fmt.Errorf("error on reading health check %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunHealthCheckUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyHealthCheck(d, resourceKsyunHealthCheck())
	if err != nil {
		return fmt.Errorf("error on updating health check %q, %s", d.Id(), err)
	}
	return resourceKsyunHealthCheckRead(d, meta)
}

func resourceKsyunHealthCheckDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveHealthCheck(d)
	if err != nil {
		return fmt.Errorf("error on deleting health check %q, %s", d.Id(), err)
	}
	return err
}
