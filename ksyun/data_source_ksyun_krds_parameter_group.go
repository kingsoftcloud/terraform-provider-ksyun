/*
Query ksyun krds parameter group information

# Example Usage

```hcl

		provider "ksyun" {
			region = "cn-beijing-6"
		}


		data "ksyun_krds_parameter_group" "foo" {
			output_file = "output_result"
			// if you give db_parameter_group_id will return the single krds parameter group
			// if you don't give this value, it will return a list of krds parameter groups
			db_parameter_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

			// keyword is a filter value that can query the results by name of description
			keyword = "name or description"
		}
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunKrdsParameterGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKrdsParameterGroupRead,
		Schema: map[string]*schema.Schema{
			// parameter
			"db_parameter_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of db parameter group.",
			},
			// query data guard
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The keyword uses to filter parameter group.",
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

			"db_parameter_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of krds db parameter groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// return values by data source query
						"parameters": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The custom parameters.",
						},
						"db_parameter_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The krds db parameter group id.",
						},
						"db_parameter_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The krds db parameter group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of this db parameter group.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The db parameter group adapts to what krds engine.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of engine.",
						},
					},
				},
			},
		},
	}
}

// dataSourceKsyunKrdsParameterGroupRead will read data source from ksyun
func dataSourceKsyunKrdsParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	krdsParameterSrv := NewKrdsParameterSrv(meta.(*KsyunClient))

	var (
		r           = dataSourceKsyunKrdsParameterGroup()
		sdkResponse []interface{}
		isById      bool
	)

	reqTransform := map[string]SdkReqTransform{
		"db_parameter_group_id": {mapping: "DBParameterGroupId"},
		"keyword":               {},
	}

	reqParameters, err := mergeDataSourcesReq(d, r, reqTransform)
	if err != nil {
		return err
	}
	// call query function
	action := "DescribeDBParameterGroup"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	if _, ok := reqParameters["DBParameterGroupId"]; ok {
		sdkResponse, err = krdsParameterSrv.describeDBParameterGroupById(reqParameters)
		if err != nil {
			return err
		}
		isById = true
	} else {
		sdkResponse, err = krdsParameterSrv.describeDBParameterGroupAll(reqParameters)
		if err != nil {
			return err
		}
	}
	if len(sdkResponse) == 0 {
		return fmt.Errorf("the current available zone not exsits any krds parameter group")
	}

	if isById {
		if err := TransformMapValue2StringWithKey("Parameters", sdkResponse); err != nil {
			return err
		}
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  sdkResponse,
		idFiled:     "DBParameterGroupId",
		targetField: "db_parameter_groups",
		extra: map[string]SdkResponseMapping{
			"DBParameterGroupId": {
				Field: "db_parameter_group_id",
			},
			"DBParameterGroupName": {
				Field: "db_parameter_group_name",
			},
		},
	})
}
