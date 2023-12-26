/*
Provides an access token resource under kcrs repository instance.

Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest"
	instance_type = "basic"
}

# Create a KcrsToken
resource "ksyun_kcrs_token" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	token_type = "Day"
	token_time = 10
	desc = "test"
	enable = true
}
```

Import

KcrsToken can be imported using `instance_id:token_id`, e.g.

```
$ terraform import ksyun_kcrs_token.foo ${instance_id}:${token_id}
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
				Description: "Instance id of repository.",
			},
			"token_type": {
				Type:     schema.TypeString,
				Required: true,

				ValidateFunc: validation.StringInSlice([]string{"Hour", "Day", "NeverExpire"},
					false),
				Description: "Token type.",
			},

			"token_time": {
				Type:     schema.TypeInt,
				Required: true,

				ValidateFunc: validation.IntBetween(1, 9999),
				// ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("token_type") == "NeverExpire" {
						return true
					}
					return false
				},

				Description: "The validation time of token. If the `token_type` is 'NeverExpire', this field is invalid.",
			},
			"desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for this token.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable this token.",
			},

			// computed
			"expire_time": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The expired time for this token.",
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

	if d.HasChanges("token_type", "token_time", "desc") {
		req["TokenType"] = d.Get("token_type")
		req["TokenTime"] = d.Get("token_time")
		req["Desc"] = d.Get("desc")
	}

	if len(req) > 0 {
		req["InstanceId"] = d.Get("instance_id")
		req["TokenId"] = d.Id()
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
	req["TokenId"] = d.Id()

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
