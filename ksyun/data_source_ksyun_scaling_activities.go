/*
This data source provides a list of ScalingActivity resources in a ScalingGroup.

# Example Usage

```hcl

	data "ksyun_scaling_activities" "default" {
	  output_file="output_result"
	  scaling_group_id = "541241314798505984"
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunScalingActivities() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingActivitiesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ScalingActivity resources that satisfy the condition.",
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A ScalingGroup ID that the desired ScalingActivity belong to.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Start Time that the desired ScalingActivity set to.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The End Time that the desired ScalingActivity set to.",
			},

			"scaling_activities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the desired ScalingActivity.",
						},
						"cause": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cause of the desired ScalingActivity.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the desired ScalingActivity.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start time the desired ScalingActivity.",
						},
						"scaling_activity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the desired ScalingActivity.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end time of the desired ScalingActivity.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The type the desired ScalingActivity.",
						},
						"error_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The error code of the desired ScalingActivity.",
						},
						"success_instance_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set:         schema.HashString,
							Description: "The success KEC Instance ID List of the desired ScalingActivity.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingActivitiesRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunScalingActivities()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"scaling_group_id": {},
		"end_time":         {},
		"start_time":       {},
	}
	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingActivity list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingActivity", req)
		resp, err := client.kecconn.DescribeScalingActivity(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingActivity list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingActivitySet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, all)

	err = dataSourceKsyunScalingActivitiesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingActivity list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingActivitiesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingActivities()
	targetName := "scaling_activities"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingActivityId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, scalingActivitySpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}

func scalingActivitySpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["Desciption"] = SdkResponseMapping{Field: "description"}
	specialMapping["SuccInsList"] = SdkResponseMapping{Field: "success_instance_list"}
	return specialMapping
}
