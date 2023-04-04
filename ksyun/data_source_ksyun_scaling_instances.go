/*
This data source provides a list of ScalingInstance resources in a ScalingGroup.

# Example Usage

```hcl

	data "ksyun_scaling_instances" "default" {
	  output_file="output_result"
	  scaling_group_id = "246b37be-5213-49da-a971-xxxxxxxxxxxx"
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunScalingInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingInstancesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of ScalingInstance resources that satisfy the condition.",
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A scaling group id that the desired ScalingInstance belong to.",
			},

			"scaling_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of scaling group ids that the desired ScalingInstance belong to.",
			},

			"health_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the health status that desired scalingInstance belong to.",
			},

			"creation_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the creation type that desired scalingInstance belong to.",
			},

			"scaling_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The KEC Instance ID of the desired ScalingInstance.",
						},
						"scaling_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The KEC Instance Name of the desired ScalingInstance.",
						},
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Health Status of the desired ScalingInstance.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation of ScalingInstance, formatted in RFC3339 time string.",
						},
						"creation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Creation Type of the desired ScalingInstance.",
						},
						"protected_from_detach": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The KEC Instance Protected Model of the desired ScalingInstance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunScalingInstancesRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var all []interface{}
	var err error
	r := dataSourceKsyunScalingInstances()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	all = []interface{}{}
	result = []map[string]interface{}{}
	req := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"scaling_group_id": {},
		"health_status":    {},
		"creation_type":    {},
	}

	req, err = SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading ScalingInstance list, %s", err)
	}

	if ids, ok := d.GetOk("scaling_instance_ids"); ok {
		SchemaSetToInstanceMap(ids, "ScalingInstanceId", &req)
	}

	for {
		req["MaxResults"] = limit
		req["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingInstance", req)
		resp, err := client.kecconn.DescribeScalingInstance(&req)
		if err != nil {
			return fmt.Errorf("error on reading ScalingInstance list req(%v):%v", req, err)
		}
		l := (*resp)["ScalingInstanceSet"].([]interface{})
		all = append(all, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	merageResultDirect(&result, all)

	err = dataSourceKsyunScalingInstancesSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingInstance list, %s", err)
	}
	return nil
}

func dataSourceKsyunScalingInstancesSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingInstances()
	targetName := "scaling_instances"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "InstanceId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string) + ":" + d.Get("scaling_group_id").(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(resource, targetName, item, nil, scalingInstanceSpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}

func scalingInstanceSpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["InstanceId"] = SdkResponseMapping{Field: "scaling_instance_id"}
	specialMapping["InstanceName"] = SdkResponseMapping{Field: "scaling_instance_name"}
	specialMapping["ProtectedFromScaleIn"] = SdkResponseMapping{Field: "protected_from_detach"}
	return specialMapping
}
