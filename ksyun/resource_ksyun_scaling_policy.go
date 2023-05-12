/*
Provides a ScalingPolicy resource.

# Example Usage

```hcl

	resource "ksyun_scaling_policy" "foo" {
	  scaling_group_id = "541241314798505984"
	  threshold = 20
	}

```

# Import

```
$ terraform import ksyun_scaling_policy.example scaling-policy-abc123456
```
*/
package ksyun

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

var ksyunScalingPolicyMetricDimensionName = []string{
	"cpu_usage",
	"mem_usage",
	"net_outtraffic",
	"net_intraffic",
	"listener_outtraffic",
	"listener_intraffic",
}

var ksyunScalingPolicyMetricFunction = []string{
	"avg",
	"min",
	"max",
}

var ksyunScalingPolicyMetricComparisonOperator = []string{
	"Greater",
	"EqualOrGreater",
	"Less",
	"EqualOrLess",
	"Equal",
	"NotEqual",
}

var ksyunScalingPolicyMetricAdjustmentType = []string{
	"TotalCapacity",
	"QuantityChangeInCapacity",
	"PercentChangeInCapacity",
}

func resourceKsyunScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingPolicyCreate,
		Read:   resourceKsyunScalingPolicyRead,
		Delete: resourceKsyunScalingPolicyDelete,
		Update: resourceKsyunScalingPolicyUpdate,
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
				err = d.Set("scaling_policy_id", items[1])
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
			"scaling_policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tf-scaling-policy",
				Description: "The Name of the desired ScalingPolicy.",
			},

			"dimension_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "cpu_usage",
				ValidateFunc: validation.StringInSlice(ksyunScalingPolicyMetricDimensionName, false),
				Description:  "The Dimension Name of the desired ScalingPolicy.Valid Value 'cpu_usage', 'mem_usage', 'net_outtraffic', 'net_intraffic', 'listener_outtraffic', 'listener_intraffic'.",
			},

			"comparison_operator": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Greater",
				ValidateFunc: validation.StringInSlice(ksyunScalingPolicyMetricComparisonOperator, false),
				Description:  "The Comparison Operator of the desired ScalingPolicy.Valid Value 'Greater', 'EqualOrGreater', 'Less', 'EqualOrLess', 'Equal', 'NotEqual'.",
			},

			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     50,
				Description: "The Threshold of the desired ScalingPolicy.",
			},

			"repeat_times": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 10),
				Description:  "The Repeat Times of the desired ScalingPolicy.Valid Value 1-10.",
			},

			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(60, 999999),
				Description:  "The Period of the desired ScalingPolicy.Min is 60.",
			},

			"function": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "avg",
				ValidateFunc: validation.StringInSlice(ksyunScalingPolicyMetricFunction, false),
				Description:  "The Function Model of the desired ScalingPolicy.Valid Value 'avg', 'min', 'max'.",
			},

			"adjustment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "QuantityChangeInCapacity",
				ValidateFunc: validation.StringInSlice(ksyunScalingPolicyMetricAdjustmentType, false),
				Description:  "The Adjustment Type of the desired ScalingPolicy.Valid Value 'TotalCapacity', 'QuantityChangeInCapacity', 'PercentChangeInCapacity'.",
			},

			"adjustment_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The Adjustment Value of the desired ScalingPolicy.Valid Value -100 ~ 100.",
			},

			"cool_down": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntAtLeast(60),
				Description:  "The Cool Down of the desired ScalingPolicy.Min is 60.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},

			"scaling_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the scaling policy.",
			},
		},
	}
}

func resourceKsyunScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingPolicy()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, resourceKsyunScalingPolicyExtra(d))
	if err != nil {
		return fmt.Errorf("error on creating ScalingPolicy, %s", err)
	}
	//zero process
	if _, ok := req["AdjustmentValue"]; !ok {
		req["AdjustmentValue"] = 0
	}

	action := "CreateScalingPolicy"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateScalingPolicy(&req)
	if err != nil {
		return fmt.Errorf("error on creating ScalingPolicy, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ReturnSet"].(map[string]interface{})["ScalingPolicyId"].(string) + ":" + req["ScalingGroupId"].(string))
	}
	return resourceKsyunScalingPolicyRead(d, meta)
}

func resourceKsyunScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingPolicy()

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, resourceKsyunScalingPolicyExtra(d))
	if err != nil {
		return fmt.Errorf("error on modifying ScalingPolicy, %s", err)
	}
	if len(req) > 0 {
		req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
		req["ScalingPolicyId"] = strings.Split(d.Id(), ":")[0]
		action := "ModifyScalingPolicy"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifyScalingPolicy(&req)
		if err != nil {
			return fmt.Errorf("error on modifying ScalingPolicy, %s", err)
		}
	}
	return resourceKsyunScalingPolicyRead(d, meta)
}

func resourceKsyunScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingPolicyId.1"] = strings.Split(d.Id(), ":")[0]
	action := "DescribeScalingPolicy"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeScalingPolicy(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingPolicy %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingPolicySet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunScalingPolicy(), items[0], nil)
		SdkResponseAutoResourceData(d, resourceKsyunScalingPolicy(), items[0].(map[string]interface{})["Metric"], nil)
	}
	return nil
}

func resourceKsyunScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingPolicyId"] = strings.Split(d.Id(), ":")[0]
	action := "DeleteScalingPolicy"
	otherErrorRetry := 10

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.DeleteScalingPolicy(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return OtherErrorProcess(&otherErrorRetry, fmt.Errorf("error on  deleting ScalingPolicy %q, %s", d.Id(), err1))
		}
	})

}

func resourceKsyunScalingPolicyExtra(d *schema.ResourceData) map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = make(map[string]SdkRequestMapping)
	fieldReqFunc := func(item interface{}, s string, source string, m *map[string]interface{}) error {
		if _, ok := (*m)[s]; !ok {
			jsonMap := make(map[string]interface{})
			jsonMap["comparisonOperator"] = d.Get("comparison_operator")
			jsonMap["dimensionName"] = d.Get("dimension_name")
			jsonMap["threshold"] = d.Get("threshold")
			jsonMap["repeatTimes"] = d.Get("repeat_times")
			jsonMap["function"] = d.Get("function")
			jsonMap["period"] = d.Get("period")
			str, err := json.Marshal(jsonMap)
			if err != nil {
				return err
			}
			(*m)[s] = string(str)
		}
		return nil
	}
	extra["dimension_name"] = SdkRequestMapping{
		Field:        "Metric",
		FieldReqFunc: fieldReqFunc,
	}
	extra["comparison_operator"] = SdkRequestMapping{
		Field:        "Metric",
		FieldReqFunc: fieldReqFunc,
	}
	extra["threshold"] = SdkRequestMapping{
		Field:        "Metric",
		FieldReqFunc: fieldReqFunc,
	}
	extra["repeat_times"] = SdkRequestMapping{
		Field:        "Metric",
		FieldReqFunc: fieldReqFunc,
	}
	extra["function"] = SdkRequestMapping{
		Field:        "Metric",
		FieldReqFunc: fieldReqFunc,
	}
	extra["period"] = SdkRequestMapping{
		Field:        "Metric",
		FieldReqFunc: fieldReqFunc,
	}
	return extra
}
