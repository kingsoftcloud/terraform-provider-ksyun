/*
Provides a list of Redis security groups in the current region.

# Example Usage

```hcl

	data "ksyun_redis_security_groups" "default" {
	  output_file       = "output_result1"
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
)

// redis security group List
func dataSourceRedisSecurityGroups() *schema.Resource {
	return &schema.Resource{
		// redis security group List List Query Function
		Read: dataSourceRedisSecurityGroupsRead,
		// Define input and output parameters
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Redis security groups that satisfy the condition.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of Redis security groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group description.",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
						},
						"updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "updated time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceRedisSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	var (
		allInstances []interface{}
		az           map[string]string
		limit        = 100
		nextToken    string
		err          error
	)

	action := "DescribeSecurityGroups"
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	if az, err = queryAz(conn); err != nil {
		return fmt.Errorf("error on reading instances, because there is no available area in the region")
	}
	for k := range az {
		readReq["AvailableZone"] = k
		for {
			readReq["Limit"] = fmt.Sprintf("%v", limit)
			if nextToken != "" {
				readReq["Offset"] = nextToken
			}
			logger.Debug(logger.ReqFormat, action, readReq)
			resp, err := conn.DescribeSecurityGroups(&readReq)
			if err != nil {
				return fmt.Errorf("error on reading redis security group list req(%v):%s", readReq, err)
			}
			logger.Debug(logger.RespFormat, action, readReq, *resp)
			result, ok := (*resp)["Data"]
			if !ok {
				break
			}
			item, ok := result.(map[string]interface{})
			if !ok {
				break
			}
			items, ok := item["list"].([]interface{})
			if !ok {
				break
			}
			if items == nil || len(items) < 1 {
				break
			}
			allInstances = append(allInstances, items...)
			if len(items) < limit {
				break
			}
			nextToken = strconv.Itoa(int(item["limit"].(float64)) + int(item["offset"].(float64)))
		}
	}

	values := GetSubSliceDByRep(allInstances, redisSecKeys)
	if err := dataSourceKscSave(d, "instances", []string{}, values); err != nil {
		return fmt.Errorf("error on save redis security group list, %s", err)
	}
	return nil
}
