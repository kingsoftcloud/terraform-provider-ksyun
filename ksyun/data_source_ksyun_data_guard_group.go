/*
Query instance auto snapshot policies information

# Example Usage

```hcl

		data "ksyun_data_guard_group" "foo" {
		  output_file = "output_result"
		}
		data "ksyun_data_guard_group" "foo1" {
			data_guard_name = "Data Guard Name"
		}

		data "ksyun_data_guard_group" "foo2" {
			data_guard_id = "Data Guard Id"
		}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunDataGuardGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunDataGuardGroupRead,
		Schema: map[string]*schema.Schema{
			// parameter
			"data_guard_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of data guard group.",
			},
			// query data guard
			"data_guard_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of data guard group.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of snapshot policies resources that satisfy the condition.",
			},

			"data_guard_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of data guard groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// return values by data source query
						"data_guard_instances_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Description: "The data guard group includes instances.",
						},
						"data_guard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data guard group id.",
						},

						"data_guard_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data guard group name.",
						},
						"data_guard_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data guard group display type.",
						},
						"data_guard_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data guard group level, if the value is Host represent machine level, and the tol represent the domain of disaster tolerance.",
						},
						"data_guard_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity of data guard group.",
						},
						"data_guard_used_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "This data guard group includes the amount of instances.",
						},
					},
				},
			},
		},
	}
}

// dataSourceKsyunDataGuardGroupRead will read data source from ksyun
func dataSourceKsyunDataGuardGroupRead(d *schema.ResourceData, meta interface{}) error {
	dataGuardSrv := DataGuardSrv{
		client: meta.(*KsyunClient),
	}
	r := dataSourceKsyunDataGuardGroup()

	reqTransform := map[string]SdkReqTransform{
		"data_guard_name": {},
		"data_guard_id":   {},
	}

	reqParameters, err := mergeDataSourcesReq(d, r, reqTransform)
	if err != nil {
		return err
	}
	// call query function
	action := "DescribeDataGuardGroup"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdKData, err := dataGuardSrv.describeDataGuardGroup(reqParameters)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  sdKData,
		idFiled:     "DataGuardId",
		targetField: "data_guard_groups",
	})
}
