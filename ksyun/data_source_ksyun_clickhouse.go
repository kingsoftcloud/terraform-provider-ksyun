/*
Query Clickhouse instance information

# Example Usage

```hcl

	data "ksyun_clickhouse" "default"{

  		instance_id = "instance_id"

  		product_type = "product_type"
  		project_ids = "project_ids"
  		tag_id = "tag_id"

  		fuzzy_search = "fuzzy_search"

  		offset = 0
  		limit = 10
	}

```
*/
package ksyun

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// 定义字段映射配置
var clickhouseFieldMap = map[string]string{
	// API字段名 -> Terraform字段名
	"InstanceId":      "instance_id",
	"InstanceName":    "instance_name",
	"InstacneConfig":  "instance_config",
	"AdminUser":       "admin_user",
	"StatusName":      "status_name",
	"Status":          "status",
	"NetworkType":     "network_type",
	"VpcId":           "vpc_id",
	"SubnetId":        "subnet_id",
	"Vip":             "vip",
	"Engine":          "engine",
	"EngineVersion":   "engine_version",
	"ProjectId":       "project_id",
	"ProjectName":     "project_name",
	"BillType":        "bill_type",
	"EbsSize":         "ebs_size",
	"EbsType":         "ebs_type",
	"Mem":             "mem",
	"Cpu":             "cpu",
	"TcpPort":         "tcp_port",
	"HttpPort":        "http_port",
	"NodeNum":         "node_num",
	"ProductId":       "product_id",
	"ProductType":     "product_type",
	"ProductTypeName": "product_type_name",
	"CreateDate":      "create_date",
	"UpdateDate":      "update_date",
	"Region":          "region",
	"Az":              "az",
	"UserId":          "user_id",
	"SecurityGroupId": "security_group_id",
	"ProductWhat":     "product_what",
	"ServiceEndTime":  "service_end_time",
	"Replicas":        "replicas",

	//实例详情差异字段
	"StorageSize":          "storage_size",
	"UsedStorageSize":      "used_storage_size",
	"StorageType":          "storage_type",
	"SecurityGroupName":    "security_group_name",
	"SecurityGroupDesc":    "security_group_desc",
	"DirectConnectionVips": "direct_connection_vips",
	"HotAndCold":           "hot_and_cold",
	"MultiAz":              "multiaz",
	//Area 特殊处理
	//ShardList 特殊处理
}

// 数值类型字段列表
var clickhouseIntFields = map[string]bool{
	"bill_type":    true,
	"ebs_size":     true,
	"mem":          true,
	"cpu":          true,
	"tcp_port":     true,
	"http_port":    true,
	"node_num":     true,
	"product_type": true,
	"replicas":     true,
	"product_what": true,
}

func dataSourceKsyunClickhouse() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunClickhouseRead,
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
				Description: "File name where to save data source results (after running terraform plan).",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of resources that satisfy the condition.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The ClickHouse instance ID. When provided, returns detailed information for that specific instance; otherwise returns a list of all instances.",
			},
			"product_type": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The product type of the instance. Valid values: 'ClickHouse_Single' (single replica) or 'ClickHouse' (high availability).",
			},
			"tag_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Filter instances by tag ID.",
			},
			"project_ids": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Comma-separated list of project IDs to filter instances.",
			},

			"fuzzy_search": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Fuzzy search filter that matches instance name, VIP, or instance ID.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "The starting offset for pagination. Default is 0 (first page).",
			},
			"limit": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "The maximum number of records to return per page. Default is 10.",
			},
			// 输出字段
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ClickHouse instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"instance_config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance configuration.",
						},
						"admin_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Admin user name.",
						},
						"status_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the instance status (Chinese).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network type of the instance. Currently supports 'VPC'.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP address.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database engine.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine version.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project name.",
						},
						"bill_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The billing type of the instance.",
						},
						"ebs_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size in GB.",
						},
						"ebs_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk type.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size in GB.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU cores.",
						},
						"tcp_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP port.",
						},
						"http_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP port.",
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes.",
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product ID.",
						},
						"product_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Product type.",
						},
						"product_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product type name.",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"az": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User ID.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group ID.",
						},
						"product_what": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The product category identifier.",
						},
						"service_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service end time.",
						},
						"replicas": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of replicas.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag ID.",
									},
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"area": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Instance area configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Master area.",
									},
									"standby": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Standby area.",
									},
								},
							},
						},
						"shard_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance shard list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Shard ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Shard name.",
									},
								},
							},
						},
						"direct_connection_vips": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Direct connection VIPs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage size.",
						},
						"used_storage_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Used storage size.",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage type.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group name.",
						},
						"security_group_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group description.",
						},
						"hot_and_cold": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hot and cold configuration.",
						},
						"multiaz": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Multi-Availability Zone deployment configuration flag.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunClickhouseRead(d *schema.ResourceData, meta interface{}) error {
	connClient, ok := meta.(*KsyunClient)
	if !ok || connClient == nil || connClient.clickhouseconn == nil {
		return fmt.Errorf("invalid ClickHouse client")
	}
	conn := connClient.clickhouseconn

	// 构建请求参数
	req := make(map[string]interface{})

	// 手动映射查询参数（只在存在时添加）
	if v, ok := d.GetOk("instance_id"); ok {
		req["InstanceId"] = v
	}
	if v, ok := d.GetOk("project_ids"); ok {
		req["ProjectIds"] = v
	}
	if v, ok := d.GetOk("tag_id"); ok {
		req["TagId"] = v
	}
	if v, ok := d.GetOk("fuzzy_search"); ok {
		req["FuzzySearch"] = v
	}
	if v, ok := d.GetOk("offset"); ok {
		req["Offset"] = v
	}
	if v, ok := d.GetOk("limit"); ok {
		req["Limit"] = v
	}
	if v, ok := d.GetOk("product_type"); ok {
		req["ProductType"] = v
	}

	action := "ListInstance"
	var resp *map[string]interface{}
	var err error

	// 若指定了 InstanceId，则调用另一个接口获取单个实例并将返回的单条数据规范化为 Data 列表
	if v, ok := req["InstanceId"]; ok && fmt.Sprintf("%v", v) != "" {
		action = "DescribeInstance"
		logger.Debug(logger.ReqFormat, action, req)
		singleResp, errGet := conn.DescribeInstance(&req)
		logger.Debug(logger.AllFormat, action, req, singleResp, errGet)
		if errGet != nil {
			return fmt.Errorf("error reading ClickHouse instance: %s", errGet)
		}

		// 规范化为与 ListInstance 相同的响应结构：Data 为列表
		newResp := make(map[string]interface{})
		if singleResp != nil {
			for k, val := range *singleResp {
				newResp[k] = val
			}
		}
		if data, exists := newResp["Data"]; exists && data != nil {
			switch data.(type) {
			case []interface{}:
				// 已经是列表，无需处理
			default:
				newResp["Data"] = []interface{}{data}
			}
		} else {
			newResp["Data"] = []interface{}{}
		}
		resp = &newResp
	} else {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.ListInstance(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)
		if err != nil {
			return fmt.Errorf("error reading ClickHouse instances: %s", err)
		}
	}

	body, ok := (*resp)["Data"]
	var instances []interface{}
	var instanceMaps []map[string]interface{}
	var instanceIds []string
	var totalCount int // 定义总数量变量

	if body == nil {
		logger.Debug(logger.ReqFormat, action, "no Data in response, returning empty list")
		// 初始化空数组
		instances = make([]interface{}, 0)
		instanceMaps = make([]map[string]interface{}, 0)
		instanceIds = make([]string, 0)
		totalCount = 0 // 设置数量为0
	} else {
		// 处理有数据的情况
		instances, ok = body.([]interface{})
		if !ok {
			// 如果类型断言失败，也当作空数据处理
			instances = make([]interface{}, 0)
		}

		// 使用映射函数处理数据
		instanceIds = make([]string, 0, len(instances))
		instanceMaps = make([]map[string]interface{}, 0, len(instances))
		for _, instance := range instances {
			if instanceInfo, ok := instance.(map[string]interface{}); ok && instanceInfo != nil {
				mappedInstance := mapClickhouseInstance(instanceInfo)
				if mappedInstance == nil {
					// 跳过无效映射，继续处理其余项
					continue
				}
				instanceMaps = append(instanceMaps, mappedInstance)

				if v, exists := mappedInstance["instance_id"]; exists && v != nil {
					instanceIds = append(instanceIds, fmt.Sprintf("%v", v))
				} else {
					// 与 instanceMaps 保持一一对应，使用空字符串占位
					instanceIds = append(instanceIds, "")
				}
			}
		}
		totalCount = len(instanceMaps) // 设置实际数量
	}

	// 将结果安全地写回 Terraform schema
	if err := d.Set("instances", instanceMaps); err != nil {
		return fmt.Errorf("failed to set instances: %w", err)
	}
	if err := d.Set("total_count", totalCount); err != nil { // 使用统一的totalCount变量
		return fmt.Errorf("failed to set total_count: %w", err)
	}

	// 设置 resource id（保证非空）
	if len(instanceIds) > 0 && instanceIds[0] != "" {
		d.SetId(instanceIds[0])
	} else {
		d.SetId(fmt.Sprintf("ksyun_clickhouse_%d", time.Now().UnixNano()))
	}

	// 尝试持久化到本地 DB，但不阻塞返回（记录错误）
	if err := dataDbSave(d, "instances", instanceIds, instanceMaps); err != nil {
		logger.Debug(logger.ReqFormat, "dataDbSave error", err)
	}

	return nil
}

// mapClickhouseInstance 映射单个实例数据
func mapClickhouseInstance(instanceInfo map[string]interface{}) map[string]interface{} {
	mapped := make(map[string]interface{})

	for apiField, tfField := range clickhouseFieldMap {
		if value, ok := instanceInfo[apiField]; ok {
			// 处理数值类型转换
			if clickhouseIntFields[tfField] {
				switch v := value.(type) {
				case float64:
					mapped[tfField] = int(v)
				case string:
					// 字符串转整数
					if intVal, err := strconv.Atoi(v); err == nil {
						mapped[tfField] = intVal
					} else {
						// 转换失败时保持原值或设置默认值
						log.Printf("[WARN] Failed to convert string to int for field %s: %v", tfField, err)
						mapped[tfField] = 0 // 或者保持原值: mapped[tfField] = value
					}
				case int:
					mapped[tfField] = v
				default:
					// 其他类型，记录警告并使用默认值
					log.Printf("[WARN] Unexpected type for field %s: %T, value: %v", tfField, value, value)
					mapped[tfField] = 0
				}
			} else {
				// 处理空值
				if value == nil {
					mapped[tfField] = ""
				} else {
					mapped[tfField] = value
				}
			}
		}
	}

	// 特殊处理 Tags
	mapped["tags"] = mapClickhouseTags(instanceInfo["Tags"])
	// 特殊处理 ShardList
	mapped["shard_list"] = mapClickhouseShardList(instanceInfo["ShardList"])
	// 特殊处理 Area
	mapped["area"] = mapClickhouseArea(instanceInfo["Area"])

	return mapped
}

// mapClickhouseTags 映射标签数据
func mapClickhouseTags(tags interface{}) []map[string]interface{} {
	tagList, ok := tags.([]interface{})
	if !ok {
		return make([]map[string]interface{}, 0)
	}

	mappedTags := make([]map[string]interface{}, 0, len(tagList))
	for _, tag := range tagList {
		if tagMap, ok := tag.(map[string]interface{}); ok {
			mappedTag := make(map[string]interface{})
			if v, ok := tagMap["tagId"].(string); ok {
				mappedTag["tag_id"] = v
			}
			if v, ok := tagMap["tagKey"].(string); ok {
				mappedTag["tag_key"] = v
			}
			if v, ok := tagMap["tagValue"].(string); ok {
				mappedTag["tag_value"] = v
			}
			mappedTags = append(mappedTags, mappedTag)
		}
	}
	return mappedTags
}

// mapClickhouseShardList 映射分片列表数据
func mapClickhouseShardList(shards interface{}) []map[string]interface{} {
	shardList, ok := shards.([]interface{})
	if !ok {
		return make([]map[string]interface{}, 0)
	}

	mappedShards := make([]map[string]interface{}, 0, len(shardList))
	for _, s := range shardList {
		if shardMap, ok := s.(map[string]interface{}); ok {
			mapped := make(map[string]interface{})
			if v, ok := shardMap["id"].(string); ok {
				mapped["id"] = v
			}
			if v, ok := shardMap["name"].(string); ok {
				mapped["name"] = v
			}
			mappedShards = append(mappedShards, mapped)
		}
	}
	return mappedShards
}

// mapClickhouseArea 映射 Area 字段，返回包含 master 与 standby 的映射
func mapClickhouseArea(area interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"master":  "",
		"standby": "",
	}

	areaMap, ok := area.(map[string]interface{})
	if !ok || areaMap == nil {
		return result
	}

	// 支持可能的大小写差异
	if v, ok := areaMap["Master"].(string); ok && v != "" {
		result["master"] = v
	} else if v, ok := areaMap["master"].(string); ok && v != "" {
		result["master"] = v
	}

	if v, ok := areaMap["Standby"].(string); ok && v != "" {
		result["standby"] = v
	} else if v, ok := areaMap["standby"].(string); ok && v != "" {
		result["standby"] = v
	}

	return result
}
