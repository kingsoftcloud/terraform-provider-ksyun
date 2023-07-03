/*
This data source provides a list of MongoDB resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to .

# Example Usage

```hcl

	data "ksyun_mongodbs" "default" {
	  output_file = "output_result"
	  iam_project_id = ""
	  instance_id = ""
	  vnet_id = ""
	  vpc_id = ""
	  name = ""
	  vip = ""
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

// instance List
func dataSourceKsyunMongodbs() *schema.Resource {
	return &schema.Resource{
		// Instance List Query Function
		Read: dataSourceMongodbInstancesRead,
		// Define input and output parameters
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of MongoDB, all the MongoDBs belong to this region will be retrieved if the instance_id is `\"\"`.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of MongoDB, all the MongoDBs belong to this region will be retrieved if the name is `\"\"`.",
			},
			"vnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of subnet. the instance will use the subnet in the current region.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of VPC. the instance will use the VPC in the current region.",
			},
			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vip of instances.",
			},

			"iam_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project instance belongs to.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of MongoDBs that satisfy the condition.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the user.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the MongoDB.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the MongoDB.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the MongoDB.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A list of MongoDB node IPs.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the MongoDB.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the MongoDB.",
						},
						"instance_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance specification.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "storage size of the instance disk.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the security group.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of the MongoDB.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of network.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of VPC linked to the instance.",
						},
						"vnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of subnet linked to the instance.",
						},
						"timing_switch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "timing switch for backup.",
						},
						"timezone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "timezone of backup.",
						},
						"time_cycle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "time cycle of backup.",
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the product.",
						},
						"pay_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of pay.",
						},
						"product_what": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "whether the instance is trial or not.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the MongoDB.",
						},
						"expiration_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expiration date of the MongoDB.",
						},
						"iam_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the project.",
						},
						"iam_project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the project.",
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of nodes.",
						},
						"mongos_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of mongos.",
						},
						"shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of shards.",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MongoDB cluster mode.",
						},
						"config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance specification.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMongodbInstancesRead(d *schema.ResourceData, meta interface{}) error {

	var (
		allInstances []interface{}
		limit        = 100
		nextToken    string
	)

	readReq := make(map[string]interface{})
	filters := []string{"iam_project_id", "instance_id", "vnet_id", "vpc_id", "name", "vip"}
	for _, v := range filters {
		if value, ok := d.GetOk(v); ok {
			readReq[Downline2Hump(v)] = fmt.Sprintf("%v", value)
		}
	}
	readReq["Limit"] = fmt.Sprintf("%v", limit)

	conn := meta.(*KsyunClient).mongodbconn

	for {
		if nextToken != "" {
			readReq["Offset"] = nextToken
		}
		logger.Debug(logger.ReqFormat, "DescribeMongoDBInstances", readReq)

		resp, err := conn.DescribeMongoDBInstances(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", readReq, err)
		}
		logger.Debug(logger.RespFormat, "DescribeMongoDBInstances", readReq, *resp)

		itemSet, ok := (*resp)["MongoDBInstancesResult"]
		if !ok {
			break
		}
		items, ok := itemSet.([]interface{})
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
		nextToken = strconv.Itoa(int((*resp)["limit"].(float64)) + int((*resp)["offset"].(float64)))
	}

	values := GetSubSliceDByRep(allInstances, mongodbInstanceKeys)
	for _, v := range values {
		v["ip"] = v["i_p"]
		delete(v, "i_p")
	}
	if err := dataSourceKscSave(d, "instances", []string{}, values); err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}

	return nil
}
