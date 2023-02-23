/*
This data source provides a list of ScalingConfiguration resources.

# Example Usage

```hcl

	data "ksyun_scaling_configurations" "default" {
	  output_file="output_result"
	  ids=[]
	  project_ids=[]
	  scaling_configuration_name= "test"
	}

```
*/
package ksyun

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunScalingConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ScalingConfiguration IDs, all the ScalingConfiguration resources belong to this region will be retrieved if the ID is `\"\"`.",
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
				Description: "Total number of ScalingConfiguration resources that satisfy the condition.",
			},
			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Project id that the desired ScalingConfiguration belongs to.",
			},
			"scaling_configuration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of ScalingConfiguration.",
			},
			"scaling_configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of ScalingConfiguration.",
						},

						"scaling_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Name of the desired ScalingConfiguration.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The CPU core size of the desired ScalingConfiguration.",
						},

						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Memory GB size of the desired ScalingConfiguration.",
						},

						"data_disk_gb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Local Volume GB size of the desired ScalingConfiguration.",
						},

						"gpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The GPU core size the desired ScalingConfiguration.",
						},

						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The System Image Id of the desired ScalingConfiguration.",
						},

						"need_monitor_agent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Monitor agent flag desired ScalingConfiguration.",
						},

						"need_security_agent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Security agent flag desired ScalingConfiguration.",
						},

						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge type.",
						},

						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The KEC instance type of the desired ScalingConfiguration.",
						},

						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The KEC instance name of the desired ScalingConfiguration.",
						},

						"instance_name_suffix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kec instance name suffix of the desired ScalingConfiguration.",
						},

						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Project Id of the desired ScalingConfiguration belong to.",
						},

						"keep_image_login": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The Flag with image login set of the desired ScalingConfiguration.",
						},

						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System disk type.",
						},

						"key_id": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The SSH key set of the desired ScalingConfiguration.",
						},

						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "System disk size.",
						},

						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "It is a nested type which documented below.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The EBS Data Disk Type of the desired data_disk.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The EBS Data Disk Size of the desired data_disk.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The Flag with delete EBS Data Disk when KEC Instance destroy.",
									},
								},
							},
						},

						"instance_name_time_suffix": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The kec instance name suffix of the desired ScalingConfiguration.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation of ScalingGroup, formatted in RFC3339 time string.",
						},

						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user data of the desired ScalingConfiguration.",
						},

						"address_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP band width.",
						},

						"band_width_share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the BWS.",
						},

						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the line.",
						},

						"address_project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project ID of the IP address.",
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunScalingConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var allScalingConfigurations []interface{}
	var err error
	r := dataSourceKsyunScalingConfigurations()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	allScalingConfigurations = []interface{}{}
	result = []map[string]interface{}{}
	readScalingConfiguration := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"ids":                        {mapping: "ScalingConfigurationId", Type: TransformWithN},
		"project_ids":                {mapping: "ProjectId", Type: TransformWithN},
		"scaling_configuration_name": {},
	}

	readScalingConfiguration, err = SdkRequestAutoMapping(d, r, false, only, nil)

	if err != nil {
		return fmt.Errorf("error on reading ScalingConfiguration list, %s", err)
	}

	for {
		readScalingConfiguration["MaxResults"] = limit
		readScalingConfiguration["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingConfiguration", readScalingConfiguration)
		resp, err := client.kecconn.DescribeScalingConfiguration(&readScalingConfiguration)
		if err != nil {
			return fmt.Errorf("error on reading ScalingConfiguration list req(%v):%v", readScalingConfiguration, err)
		}
		l := (*resp)["ScalingConfigurationSet"].([]interface{})
		allScalingConfigurations = append(allScalingConfigurations, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range allScalingConfigurations {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["ScalingConfigurationName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, allScalingConfigurations)
	}

	err = dataSourceKsyunScalingConfigurationsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingConfigurationName list, %s", err)
	}

	return nil
}

func scalingConfigurationSpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["StorageSize"] = SdkResponseMapping{Field: "data_disk_gb"}
	specialMapping["DataDiskEbsDetail"] = SdkResponseMapping{
		Field: "data_disks",
		FieldRespFunc: func(i interface{}) interface{} {
			var result []map[string]interface{}
			result = []map[string]interface{}{}
			v := i.(string)
			var dat []interface{}
			if err := json.Unmarshal([]byte(v), &dat); err == nil {
				for _, v := range dat {
					d := v.(map[string]interface{})
					r := make(map[string]interface{})
					r["delete_with_instance"] = d["deleteWithInstance"]
					r["disk_size"] = d["size"]
					r["disk_type"] = d["type"]
					result = append(result, r)
				}
			}
			return result
		},
	}
	specialMapping["InstanceTypeSet"] = SdkResponseMapping{
		Field: "instance_type",
		FieldRespFunc: func(i interface{}) interface{} {
			var result string
			for _, v := range i.([]interface{}) {
				result = result + v.(string) + ","
			}
			if len(result) > 0 {
				result = result[0 : len(result)-1]
			}
			return result
		},
	}
	specialMapping["KeyIds"] = SdkResponseMapping{
		Field: "key_id",
		FieldRespFunc: func(i interface{}) interface{} {
			var (
				value []string
			)
			_ = json.Unmarshal([]byte(i.(string)), &value)
			return value

		},
	}
	return specialMapping
}

func dataSourceKsyunScalingConfigurationsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingConfigurations()
	targetName := "scaling_configurations"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingConfigurationId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item["ScalingConfigurationId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			delete(item, "InstanceType")
			return SdkResponseAutoMapping(resource, targetName, item, nil, scalingConfigurationSpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}
