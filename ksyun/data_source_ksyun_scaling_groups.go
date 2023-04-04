/*
This data source provides a list of ScalingGroup resources .

# Example Usage

```hcl

	data "ksyun_scaling_groups" "default" {
	  output_file="output_result"
	  vpc_id = "246b37be-5213-49da-a971-xxxxxxxxxxxx"
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

func dataSourceKsyunScalingGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ScalingGroup IDs, all the ScalingGroup resources belong to this region will be retrieved if the ID is `\"\"`.",
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
				Description: "Total number of ScalingGroup resources that satisfy the condition.",
			},
			"scaling_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of the desired ScalingGroup.",
			},
			"scaling_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Scaling Configuration ID of the desired ScalingGroup set to.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC ID of the desired ScalingGroup set to.",
			},
			"scaling_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Security Group ID of the desired ScalingGroup set to.",
						},

						"scaling_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of the desired ScalingGroup.",
						},
						"scaling_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Scaling Configuration ID of the desired ScalingGroup set to.",
						},

						"scaling_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Scaling Configuration Name of the desired ScalingGroup set to.",
						},

						"min_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Min KEC instance size of the desired ScalingGroup set to.",
						},

						"max_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Min KEC instance size of the desired ScalingGroup set to.",
						},

						"instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The KEC instance Number of the desired ScalingGroup set to.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation of ScalingGroup, formatted in RFC3339 time string.",
						},

						"remove_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The KEC instance remove policy of the desired ScalingGroup set to.",
						},

						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPC ID of the desired ScalingGroup set to.",
						},

						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Security Group ID of the desired ScalingGroup set to.",
						},

						"security_group_id_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set:         schema.HashString,
							Description: "A list of the Security Group IDs.",
						},

						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Status of the desired ScalingGroup.",
						},

						"desired_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Desire Capacity KEC instance count of the desired ScalingGroup set to.",
						},

						"subnet_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Subnet Strategy of the desired ScalingGroup set to.",
						},

						"subnet_id_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set:         schema.HashString,
							Description: "The Subnet ID Set of the desired ScalingGroup set to.",
						},

						"slb_config_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The SLB Config Set of the desired ScalingGroup set to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"slb_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The SLB ID of the desired ScalingGroup set to.",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Listener ID of the desired ScalingGroup set to.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The weight of the desired ScalingGroup set to.",
									},
									"health_check_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The health check type of the desired ScalingGroup set to.",
									},
									"server_port_set": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Description: "The Server Port Set of the desired ScalingGroup set to.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingGroupsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunScalingGroups()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"ids":                      {mapping: "ScalingGroupId", Type: TransformWithN},
		"scaling_group_name":       {},
		"scaling_configuration_id": {},
		"vpc_id":                   {},
	}

	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroup list, %s", err)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingGroup", req)
		resp, err := client.kecconn.DescribeScalingGroup(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingGroup list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingGroupSet"].([]interface{})
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
			if r != nil && !r.MatchString(item["ScalingGroupName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, all)
	}

	err = dataSourceKsyunScalingGroupsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroup list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingGroupsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingGroups()
	targetName := "scaling_groups"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingGroupId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item["ScalingGroupId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, nil)
		},
		TargetName: targetName,
	})
	return err
}
