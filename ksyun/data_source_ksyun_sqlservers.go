/*
Query HRDS-ss instance information

# Example Usage

```hcl

	data "ksyun_sqlservers" "search-sqlservers"{
	  output_file = "output_file"
	  db_instance_identifier = "***"
	  db_instance_type = "HRDS-SS"
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
	"strings"
	"time"
)

func dataSourceKsyunSqlServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSqlServerRead,
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
				Description: "Total number of instance that satisfy the condition.",
			},
			"db_instance_identifier": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "source instance identifier.",
			},

			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "HRDS hrds (highly available), RR (read-only), trds (temporary).",
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
				Description: "defaults to all projects.",
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
			"sqlservers": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "a list of instance.",
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
										Description: "the number of the vcpu.",
									},
									"disk": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "hard disk size..",
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
							Optional:    true,
							Computed:    true,
							Description: "instance ID.",
						},
						"db_instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "instance name.",
						},
						"db_instance_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "instance status.",
						},
						"db_instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "instance type.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "group ID.",
						},
						"vip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "virtual IP.",
						},
						"engine": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Database Engine.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "database engine version.",
						},
						"instance_create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "instance creation time.",
						},
						"master_user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "primary account user name.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "virtual private network ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"publicly_accessible": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "publicly accessible.",
						},
						"read_replica_db_instance_identifiers": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "read only instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "ID.",
									},
									"vip": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "VIP.",
									},
									"read_replica_db_instance_identifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "ID.",
									},
								},
							},
						},
						"bill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Bill type.",
						},
						"order_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Order type.",
						},
						"order_source": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Order source.",
						},
						"master_availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Master AZ.",
						},
						"slave_availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Slave AZ.",
						},
						"multi_availability_zone": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Multi availability zone.",
						},
						"product_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Product ID.",
						},
						"order_use": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Order Use.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Project ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Project name.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Region.",
						},
						"bill_type_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Bill Type ID.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Port number.",
						},
						"db_parameter_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "parameter group ID.",
						},
						"datastore_version_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Data store version ID.",
						},
						"disk_used": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "hard disk usage.",
						},
						"preferred_backup_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "preferred backup time.",
						},
						"product_what": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Product what.",
						},
						"service_start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Service start time.",
						},
						"order_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Order ID.",
						},
						"sub_order_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Sub order ID.",
						},
						"audit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Audit.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Security group ID.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "AZ.",
						},
						"db_source": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "DB source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_instance_identifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "DB instance Identifier.",
									},
									"db_instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "DB instance name.",
									},
									"db_instance_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "DB instance Type.",
									},
									"point_in_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Point in time.",
									},
								},
							},
						},
						"service_end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Service end time.",
						},
						"eip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "EIP address.",
						},
						"eip_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "EIP Port number.",
						},
						"rip": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "rip.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	desReq := make(map[string]interface{})
	des := []string{
		"DBInstanceStatus",
		"DBInstanceType",
		"DBInstanceIdentifier",
		"Keyword",
		"ExpiryDateLessThan",
		"Marker",
		"MaxRecords",
	}
	for _, v := range des {
		if v1, ok := d.GetOk(strings.ToLower(v)); ok {
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
		return fmt.Errorf("error on reading Instance(sqlserver) body %+v", (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})
	if len(instances) == 0 {
		return fmt.Errorf("empty on reading Instance(sqlserver) body %+v", *resp)
	}

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
			} else if k == "ReadReplicaDBInstanceIdentifiers" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[Camel2Hungarian(j)] = q
						}
						wtf[num] = rrmap
					}
					krdsMap["read_replica_db_instance_identifiers"] = wtf
				}
			} else if k == "DBSource" {
				dbsource := v.(map[string]interface{})
				dbsourcemap := make(map[string]interface{})
				for j, q := range dbsource {
					dbsourcemap[Camel2Hungarian(j)] = q
				}
				wtf := make([]interface{}, 1)
				wtf[0] = dbsourcemap
				krdsMap["db_source"] = wtf
			} else {
				dk := Camel2Hungarian(k)
				if _, ok := dataSourceKsyunSqlServer().Schema["sqlservers"].Elem.(*schema.Resource).Schema[dk]; ok {
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
	_ = dataDbSave(d, "sqlservers", krdsIds, krdsMapList)

	return nil
}
