/*
This data source provides a list of Rabbitmq resources according to their name, Instance ID, Subnet ID, VPC ID and the Project ID they belong to.

# Example Usage

```hcl

	data "ksyun_rabbitmqs" "default" {
	  output_file = "output_result"
	  project_id = ""
	  instance_id = ""
	  instance_name = ""
	  subnet_id = ""
	  vpc_id = ""
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
func dataSourceKsyunRabbitmqs() *schema.Resource {
	return &schema.Resource{
		// Instance List Query Function
		Read: dataSourceRabbitmqInstancesRead,
		// Define input and output parameters
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Total number of RabbitMQs that satisfy the condition.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "One or more project IDs.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of Rabbitmq, all the Rabbitmqs belong to this region will be retrieved if the instance_id is `\"\"`.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of RabbitMQ.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the subnet.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the VPC.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of RabbitMQ.",
			},
			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vip of RabbitMQs.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of RabbitMQ instances. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the RabbitMQ instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the RabbitMQ instance.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project instance belongs to.",
						},
						"instance_password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The administrator password of instance.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of VPC linked to the instance.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of subnet linked to the instance.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of instance engine.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project name of instance belong.",
						},
						"bill_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance charge type.",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration of instance use.",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mode of instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The class of instance cpu and memory.",
						},
						"ssd_disk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of instance disk, measured in GB (GigaByte).",
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the number of instance node.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone where instance is located.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the engine of the instance.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of user.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"status_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the status name of the instance.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the vip of the instance.",
						},
						"web_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the web vip of the instance.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the protocol of the instance.",
						},
						"security_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "the id of the security group.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the network type of the instance.",
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the product id of the instance.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time of creation.",
						},
						"expiration_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time of expiration.",
						},
						"product_what": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "whether the instance is trial or not.",
						},
						"mode_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mode name of the instance.",
						},
						"eip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EIP.",
						},
						"web_eip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Web EIP.",
						},
						"eip_egress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "egress of the EIP.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "port number.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the status of the instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceRabbitmqInstancesRead(d *schema.ResourceData, meta interface{}) error {

	var (
		allInstances []interface{}
		limit        = 100
		nextToken    string
	)
	r := dataSourceKsyunRabbitmqs()

	only := map[string]SdkReqTransform{
		"instance_id": {Type: TransformDefault},
		"project_id":  {Type: TransformDefault},
		"vpc_id":      {Type: TransformDefault},
		"vip":         {Type: TransformDefault},
		"subnet_id":   {Type: TransformDefault},
	}

	req, err := SdkRequestAutoMapping(d, r, false, only, nil)
	if err != nil {
		return fmt.Errorf("error on reading Instance list, %s", err)
	}
	req["limit"] = fmt.Sprintf("%v", limit)

	conn := meta.(*KsyunClient).rabbitmqconn

	for {
		if nextToken != "" {
			req["offset"] = nextToken
		}
		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstances", req)

		resp, err := conn.DescribeInstances(&req)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", req, err)
		}
		logger.Debug(logger.RespFormat, "DescribeRabbitmqInstances", req, *resp)

		result, ok := (*resp)["Data"]
		if !ok {
			break
		}
		item, ok := result.(map[string]interface{})
		if !ok {
			break
		}
		items, ok := item["Instances"].([]interface{})
		if !ok {
			break
		}
		if len(items) < 1 {
			break
		}
		allInstances = append(allInstances, items...)
		if len(items) < limit {
			break
		}
		nextToken = strconv.Itoa(int(item["limit"].(float64)) + int(item["Offset"].(float64)))
	}

	values := GetSubSliceDByRep(allInstances, rabbitmqInstanceKeys)

	if err := dataSourceKscSave(d, "instances", []string{}, values); err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}

	return nil
}
