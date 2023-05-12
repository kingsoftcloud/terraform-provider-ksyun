/*
Query HRDS and RDS-rr instance information

# Example Usage

```hcl

	data "ksyun_krds" "search-krds"{
	  output_file = "output_file"
	  db_instance_identifier = "***"
	  db_instance_type = "HRDS,RR,TRDS"
	  keyword = ""
	  order = ""
	  project_id = ""
	  marker = ""
	  max_records = ""
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func dataSourceKsyunKrds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKrdsRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "will return the file name of the content store.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of resources that satisfy the condition.",
			},
			"db_instance_identifier": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "instance ID (passed in the instance ID to get the details of the instance, otherwise get the list).",
			},

			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "HRDS (highly available), RR (read-only), TRDS (temporary).",
			},
			"db_instance_status": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "status of the instance, ACTIVE or INVALID.",
			},
			"keyword": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "fuzzy filter by name / VIP.",
			},
			"order": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "case sensitive, value range: default (default sorting method), group (sorting by replication group, will rank read-only instances after their primary instances).",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "the default value is all projects.",
			},
			"marker": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "record start offset.",
			},
			"max_records": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "the maximum number of entries in the result of each page. Value range: 1-100.",
			},
			// 与存入数据一致datakey
			"krds": {
				Type: schema.TypeList,
				//Optional:    true,
				Computed:    true,
				Description: "An information list of KRDS.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_class": {
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "instance specification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "id of the DBInstanceClass.",
									},
									"vcpus": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "number of CPUs.",
									},
									"disk": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "hard disk size.",
									},
									"ram": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "memory size.",
									},
									"iops": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "IOPS.",
									},
									"max_conn": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "max connection.",
									},
									"mem": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "memory size.",
									},
								},
							},
						},
						"db_instance_identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance ID.",
						},
						"db_instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance name.",
						},
						"db_instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance status.",
						},
						"db_instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance type.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the parameter group.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "virtual IP.",
						},
						"engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database Engine.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "database engine version.",
						},
						"instance_create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance creation time.",
						},
						"master_user_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "primary account user name.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the subnet.",
						},
						"publicly_accessible": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "network type, true(VPC instance), false(classic instance).",
						},
						"bill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "bill type.",
						},
						"master_availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AZ that master belongs to.",
						},
						"slave_availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AZ that slave belongs to.",
						},
						"multi_availability_zone": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Multi availability zone.",
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product ID.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Project ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "port number.",
						},
						"db_parameter_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the parameter group.",
						},
						"disk_used": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "hard disk usage.",
						},
						"preferred_backup_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Preferred backup time.",
						},
						"service_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service start time.",
						},
						"audit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether audit is supported or not.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the security group.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Availability zone.",
						},
						"service_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "service end time.",
						},
						"eip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EIP address.",
						},
						"eip_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EIP port.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKrdsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	desReq := make(map[string]interface{})
	des := []string{
		"DBInstanceIdentifier",
		"DBInstanceType",
		"DBInstanceStatus",
		"Keyword",
		"Order",
		"ProjectId",
		"Marker",
		"MaxRecords",
	}
	for _, v := range des {
		if v1, ok := d.GetOk(Camel2Hungarian(v)); ok {
			desReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, desReq)
	resp, err := conn.DescribeDBInstances(&desReq)
	logger.Debug(logger.AllFormat, action, desReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Instance(sqlserver)  %s", err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading Instance(krds) body %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})

	krdsIds := make([]string, len(instances))
	krdsMapList := make([]map[string]interface{}, len(instances))
	for num, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		krdsMap := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "DBInstanceClass" {
				dbclass := v.(map[string]interface{})
				dbinstanceclass := make(map[string]interface{})
				for j, q := range dbclass {
					dbinstanceclass[Camel2Hungarian(j)] = q
				}
				// shit 这里不传list会出现各种报错，我日了
				wtf := make([]interface{}, 1)
				wtf[0] = dbinstanceclass
				krdsMap["db_instance_class"] = wtf
			} else {
				dk := Camel2Hungarian(k)
				if _, ok := dataSourceKsyunKrds().Schema["krds"].Elem.(*schema.Resource).Schema[dk]; ok {
					krdsMap[dk] = v
				}
			}

		}

		logger.DebugInfo(" converted ---- %+v ", krdsMap)

		krdsIds[num] = krdsMap["db_instance_identifier"].(string)
		logger.DebugInfo("krdsIds fuck : %v", krdsIds)
		krdsMapList[num] = krdsMap
	}

	logger.DebugInfo(" converted ---- %+v ", krdsMapList)
	_ = dataDbSave(d, "krds", krdsIds, krdsMapList)

	return nil
}
