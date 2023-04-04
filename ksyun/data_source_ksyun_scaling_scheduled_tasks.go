/*
This data source provides a list of ScalingScheduledTask resources in a ScalingGroup.

# Example Usage

```hcl

	data "ksyun_scaling_scheduled_tasks" "default" {
	  output_file="output_result"
	  scaling_group_id = "541241314798xxxxxx"
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunScalingScheduledTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingScheduledTasksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of resource IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ScalingScheduledTask resources that satisfy the condition.",
			},

			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A scaling group id that the desired ScalingScheduledTask belong to.",
			},

			"scaling_scheduled_task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name that the desired ScalingScheduledTask.",
			},

			"scaling_scheduled_tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ScalingGroup ID of the desired ScalingScheduledTask belong to.",
						},
						"scaling_scheduled_task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the desired ScalingScheduledTask.",
						},

						"scaling_scheduled_task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of the desired ScalingScheduledTask.",
						},

						"readjust_max_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Readjust Max Size of the desired ScalingScheduledTask.",
						},

						"readjust_min_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Readjust Min Size of the desired ScalingScheduledTask.",
						},

						"readjust_expect_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Readjust Expect Size of the desired ScalingScheduledTask.",
						},

						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Start Time of the desired ScalingScheduledTask.",
						},

						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The End Time Operator of the desired ScalingScheduledTask.",
						},

						"recurrence": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Recurrence of the desired ScalingScheduledTask.",
						},

						"repeat_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Repeat Unit of the desired ScalingScheduledTask.",
						},

						"repeat_cycle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Repeat Cycle the desired ScalingScheduledTask.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation of ScalingScheduledTask, formatted in RFC3339 time string.",
						},

						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Description of the desired ScalingScheduledTask.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingScheduledTasksRead(d *schema.ResourceData, meta interface{}) error {
	resource := dataSourceKsyunScalingScheduledTasks()
	var result []map[string]interface{}
	var all []interface{}
	var err error

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"ids":                         {mapping: "ScalingScheduledTaskId", Type: TransformWithN},
		"scaling_group_id":            {},
		"scaling_scheduled_task_name": {},
	}

	req, err = SdkRequestAutoMapping(d, resource, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingScheduledTask list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScheduledTask", req)
		resp, err := client.kecconn.DescribeScheduledTask(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingScheduledTask list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingScheduleTaskSet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range all {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["ScalingScheduledTaskName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, all)
	}

	err = dataSourceKsyunScalingScheduledTasksSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingScheduledTask list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingScheduledTasksSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingScheduledTasks()
	targetName := "scaling_scheduled_tasks"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingScheduledTaskId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, nil)
		},
		TargetName: targetName,
	})
	return err
}
