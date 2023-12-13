/*
Provides an KcrsToken resource.

Example Usage

```hcl
# Create a KcrsToken
resource "ksyun_KcrsToken" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  KcrsToken_name = "ksc_KcrsService"
  bill_type = 1
  service_id = "KcrsToken_30G"
  project_id="0"
}
```

Import

KcrsToken can be imported using the id, e.g.

```
$ terraform import ksyun_KcrsToken.default KcrsService67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunKcrsToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKcrsTokenCreate,
		Read:   resourceKsyunKcrsTokenRead,
		Update: resourceKsyunKcrsTokenUpdate,
		Delete: resourceKsyunKcrsTokenDelete,
		Importer: &schema.ResourceImporter{
			State: importKcrsToken,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the project.",
			},
			"token_type": {
				Type:     schema.TypeString,
				Required: true,

				ValidateFunc: validation.StringInSlice([]string{"Hour", "Day", "NeverExpire"},
					false),
				Description: "the ID of the KcrsToken.",
			},
			// "charge_type": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	Default:  "HourlyInstantSettlement",
			// 	ValidateFunc: validation.StringInSlice([]string{"HourlyInstantSettlement"},
			// 		false),
			// 	Description: "the link type of the KcrsToken. Valid Values: 'DDoS_BGP'.",
			// },
			"token_time": {
				Type:     schema.TypeBool,
				Required: true,

				// ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("token_type") == "NeverExpire" {
						return true
					}
					return false
				},

				Description: "the max ip count that can bind to the KcrsToken,value range: [10, 100].",
			},
			"desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// computed
			"expire_time": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
func resourceKsyunKcrsTokenCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsTokenService := KcrsService{meta.(*KsyunClient)}
	err = kcrsTokenService.CreateKcrsToken(d, resourceKsyunKcrsToken())
	if err != nil {
		return fmt.Errorf("error on creating kcrs token %q, %s", d.Id(), err)
	}
	return resourceKsyunKcrsTokenRead(d, meta)
}

func resourceKsyunKcrsTokenRead(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsTokenService := KcrsService{meta.(*KsyunClient)}
	err = kcrsTokenService.ReadAndSetKcrsInstanceToken(d, resourceKsyunKcrsInstance())
	if err != nil {
		return fmt.Errorf("error on reading kcrs token %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKcrsTokenUpdate(d *schema.ResourceData, meta interface{}) (err error) {

	var (
		client  = meta.(*KsyunClient)
		kcrsSrv = KcrsService{client: client}
		conn    = client.kcrsconn
		req     = make(map[string]interface{})
	)

	if d.HasChanges("token_type", "token_time") {
		req["TokenType"] = d.Get("token_type")
		req["TokenTime"] = d.Get("token_time")
	}

	if d.HasChange("desc") {
		req["Desc"] = d.Get("desc")
	}
	if len(req) > 0 {
		req["InstanceId"] = d.Get("instance_id")
		_, actionErr := conn.ModifyInstanceTokenInformation(&req)
		if actionErr != nil {
			err = multierror.Append(err, fmt.Errorf("error on updating kcrs token information %q, %s", d.Id(), actionErr))
		}
	}
	if d.HasChange("enable") {
		if actionErr := kcrsSrv.modifyInstanceTokenStatus(d, resourceKsyunKcrsToken()); actionErr != nil {
			err = multierror.Append(err, fmt.Errorf("error on updating kcrs token status %q, %s", d.Id(), actionErr))
		}
	}

	if err != nil {
		return err
	}

	return resourceKsyunKcrsTokenRead(d, meta)
}

func resourceKsyunKcrsTokenDelete(d *schema.ResourceData, meta interface{}) (err error) {

	var (
		req              = make(map[string]interface{})
		kcrsTokenService = KcrsService{meta.(*KsyunClient)}
	)
	instanceId := d.Get("instance_id").(string)
	req["InstanceId"] = instanceId
	req["TokenId"] = d.Get("TokenId")

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		conn := kcrsTokenService.client.kcrsconn
		_, err := conn.DeleteInstanceToken(&req)
		if err != nil {
			if _, readErr := kcrsTokenService.ReadKcrsInstanceToken(d, instanceId); err != nil && notFoundError(readErr) {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error on deleting kcrs Token %q, %s", d.Id(), err)
	}
	return err

}
