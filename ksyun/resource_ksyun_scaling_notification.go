/*
Provides a ScalingNotification resource.

# Example Usage

```hcl

	resource "ksyun_scaling_notification" "foo" {
	  scaling_group_id = "541241314798505984"
	  scaling_notification_types = ["1","3"]
	}

```

# Import

```
$ terraform import ksyun_scaling_notification.example scaling-notification-abc123456
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"strings"
	"time"
)

func resourceKsyunScalingNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingNotificationCreate,
		Read:   resourceKsyunScalingNotificationRead,
		Delete: resourceKsyunScalingNotificationDelete,
		Update: resourceKsyunScalingNotificationUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				var err error
				items := strings.Split(d.Id(), ":")
				if len(items) != 2 {
					return nil, fmt.Errorf("id must split with %s and size %v", ":", 2)
				}
				err = d.Set("scaling_group_id", items[0])
				if err != nil {
					return nil, err
				}
				err = d.Set("scaling_notification_id", items[1])
				if err != nil {
					return nil, err
				}
				d.SetId(items[1] + ":" + items[0])
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{

			"scaling_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The ScalingGroup ID of the desired ScalingNotification belong to.",
			},

			"scaling_notification_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The List Types of the desired ScalingNotification.Valid Value '1', '2', '3', '4', '5', '6'.",
			},

			"scaling_notification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the notification.",
			},
		},
	}
}

func resourceKsyunScalingNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingNotification()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, resourceKsyunScalingNotificationExtra())
	if err != nil {
		return fmt.Errorf("error on creating ScalingNotification, %s", err)
	}
	//query first
	resp, err = conn.DescribeScalingNotification(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingNotification %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingNotificationSet"].([]interface{})
		if ok && len(items) > 0 {
			d.SetId((items[0]).(map[string]interface{})["ScalingNotificationId"].(string) + ":" + req["ScalingGroupId"].(string))
			//process update
			return resourceKsyunScalingNotificationUpdate(d, meta)
		}

	}

	action := "CreateScalingNotification"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateScalingNotification(&req)
	if err != nil {
		return fmt.Errorf("error on creating ScalingNotification, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ScalingNotificationId"].(string) + ":" + req["ScalingGroupId"].(string))
	}
	return resourceKsyunScalingNotificationRead(d, meta)
}

func resourceKsyunScalingNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingNotification()

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, resourceKsyunScalingNotificationExtra())
	if err != nil {
		return fmt.Errorf("error on modifying ScalingNotification, %s", err)
	}
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingNotificationId"] = strings.Split(d.Id(), ":")[0]
	action := "ModifyScalingNotification"
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.ModifyScalingNotification(&req)
	if err != nil {
		return fmt.Errorf("error on modifying ScalingNotification, %s", err)
	}
	return resourceKsyunScalingNotificationRead(d, meta)
}

func resourceKsyunScalingNotificationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingNotificationId.1"] = strings.Split(d.Id(), ":")[0]
	action := "DescribeScalingNotification"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeScalingNotification(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingNotification %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingNotificationSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunScalingNotification(), items[0], nil)
	}
	return nil
}

func resourceKsyunScalingNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingNotificationId"] = strings.Split(d.Id(), ":")[0]
	action := "DeleteScalingNotification"
	otherErrorRetry := 10

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.ModifyScalingNotification(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return OtherErrorProcess(&otherErrorRetry, fmt.Errorf("error on  deleting ScalingNotification %q, %s", d.Id(), err1))
		}
	})

}

func resourceKsyunScalingNotificationExtra() map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = make(map[string]SdkRequestMapping)
	extra["scaling_notification_types"] = SdkRequestMapping{
		Field: "NotificationType.",
		FieldReqFunc: func(item interface{}, s string, source string, m *map[string]interface{}) error {
			if x, ok := item.(*schema.Set); ok {
				for i, value := range (*x).List() {
					if d, ok := value.(string); ok {
						(*m)[s+strconv.Itoa(i+1)] = d
					}
				}
			}
			return nil
		},
	}
	return extra
}
