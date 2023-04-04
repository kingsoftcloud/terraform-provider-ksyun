/*
This data source provides a list of ScalingPolicy resources in a ScalingGroup.

# Example Usage

```hcl

	data "ksyun_scaling_policies" "default" {
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

func dataSourceKsyunScalingPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of policy IDs.",
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
				Description: "Total number of ScalingPolicy resources that satisfy the condition.",
			},

			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A scaling group id that the desired ScalingPolicy belong to.",
			},

			"scaling_policies_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name that the desired ScalingPolicy.",
			},

			"scaling_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ScalingGroup ID of the desired ScalingPolicy belong to.",
						},
						"scaling_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the desired ScalingPolicy.",
						},

						"scaling_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of the desired ScalingPolicy.",
						},

						"adjustment_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Adjustment Type of the desired ScalingPolicy.",
						},

						"adjustment_value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Adjustment Value of the desired ScalingPolicy.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation of ScalingPolicy, formatted in RFC3339 time string.",
						},

						"cool_down": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Cool Down of the desired ScalingPolicy.",
						},

						"dimension_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Dimension Name of the desired ScalingPolicy.",
						},

						"comparison_operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Comparison Operator of the desired ScalingPolicy.",
						},

						"threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Threshold of the desired ScalingPolicy.",
						},

						"repeat_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Repeat Times of the desired ScalingPolicy.",
						},

						"period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Period of the desired ScalingPolicy.",
						},

						"function": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Function Model of the desired ScalingPolicy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	resource := dataSourceKsyunScalingPolicies()
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
		"ids":                   {mapping: "ScalingPolicyId", Type: TransformWithN},
		"scaling_group_id":      {},
		"scaling_policies_name": {},
	}

	req, err = SdkRequestAutoMapping(d, resource, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingPolicy list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingPolicy", req)
		resp, err := client.kecconn.DescribeScalingPolicy(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingPolicy list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingPolicySet"].([]interface{})
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
			if r != nil && !r.MatchString(item["ScalingPolicyName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, all)
	}

	err = dataSourceKsyunScalingPoliciesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingPolicy list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingPoliciesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingPolicies()
	targetName := "scaling_policies"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingPolicyId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			var compute map[string]interface{}
			if item["Metric"] != nil {
				compute, _ = SdkMapMapping(item["Metric"].(map[string]interface{}), SdkSliceData{
					SliceMappingFunc: func(m map[string]interface{}) map[string]interface{} {
						return SdkResponseAutoMapping(resource, targetName, m, nil, nil)
					},
				})
			}
			return SdkResponseAutoMapping(resource, targetName, item, compute, nil)
		},
		TargetName: targetName,
	})
	return err
}
