package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
)

// instance List
func dataSourceRedisInstances() *schema.Resource {
	return &schema.Resource{
		// Instance List Query Function
		Read: dataSourceRedisInstancesRead,
		// Define input and output parameters
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fuzzy_search": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bill_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"order_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iam_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iam_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timing_switch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shard_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"shard_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_ro": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"readonly_node": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"create_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"proxy": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRedisInstancesRead(d *schema.ResourceData, meta interface{}) error {
	var (
		allInstances []interface{}
		//az           map[string]string
		item      interface{}
		resp      *map[string]interface{}
		ok        bool
		limit     = 100
		nextToken string
		err       error
	)

	action := "DescribeCacheClusters"
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	if v, ok := d.GetOk("iam_project_id"); ok {
		readReq["IamProjectId"] = v
	}
	if v, ok := d.GetOk("fuzzy_search"); ok {
		readReq["FuzzySearch"] = v
	}
	if v, ok := d.GetOk("cache_id"); ok {
		readReq["CacheId"] = v
	}
	if v, ok := d.GetOk("vnet_id"); ok {
		readReq["VnetId"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		readReq["VpcId"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		readReq["Name"] = v
	}
	if v, ok := d.GetOk("vip"); ok {
		readReq["Vip"] = v
	}
	//if az, err = queryAz(conn); err != nil {
	//	return fmt.Errorf("error on reading instances, because there is no available area in the region")
	//}
	for {
		readReq["Limit"] = fmt.Sprintf("%v", limit)
		if nextToken != "" {
			readReq["Offset"] = nextToken
		}
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err := conn.DescribeCacheClusters(&readReq)
		if err != nil {
			return fmt.Errorf("error on reading instance list req(%v):%s", readReq, err)
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

	readOnlyAction := "DescribeCacheReadonlyNode"
	readOnlyConn := meta.(*KsyunClient).kcsv2conn
	readOnlyReq := make(map[string]interface{})

	paramAction := "DescribeCacheParameters"
	paramConn := meta.(*KsyunClient).kcsv1conn
	readParamReq := make(map[string]interface{})
	for _, v := range allInstances {
		instance := v.(map[string]interface{})

		// query instance parameter
		readParamReq["CacheId"] = instance["cacheId"]
		if instance["az"] != nil {
			readParamReq["AvailableZone"] = instance["az"]
		}
		logger.Debug(logger.ReqFormat, paramAction, readParamReq)
		if resp, err = paramConn.DescribeCacheParameters(&readParamReq); err != nil {
			return fmt.Errorf("error on reading instance parameter %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, paramAction, readParamReq, *resp)
		paramData := (*resp)["Data"].([]interface{})
		if len(paramData) > 0 {
			params := make(map[string]interface{})
			for _, d := range paramData {
				param := d.(map[string]interface{})
				params[param["name"].(string)] = fmt.Sprintf("%v", param["currentValue"])
			}
			instance["parameters"] = params
		}

		// query instance node
		if int(instance["mode"].(float64)) != 2 {
			continue
		}
		readOnlyReq["CacheId"] = instance["cacheId"]
		if instance["az"] != nil {
			readOnlyReq["AvailableZone"] = instance["az"]
		}
		logger.Debug(logger.ReqFormat, readOnlyAction, readOnlyReq)
		if resp, err = readOnlyConn.DescribeCacheReadonlyNode(&readOnlyReq); err != nil {
			fmt.Printf("error on reading instance node %q, %s", d.Id(), err)
			continue
		}
		logger.Debug(logger.RespFormat, readOnlyAction, readOnlyReq, *resp)
		if item, ok = (*resp)["Data"]; !ok {
			continue
		}
		items, ok := item.([]interface{})
		if !ok || len(items) == 0 {
			continue
		}
		result := make(map[string]interface{})
		var data []interface{}
		for _, v := range items {
			vMap := v.(map[string]interface{})
			result["instance_id"] = vMap["instanceId"]
			result["name"] = vMap["name"]
			result["port"] = fmt.Sprintf("%v", vMap["port"])
			result["ip"] = vMap["ip"]
			result["status"] = vMap["status"]
			result["create_time"] = vMap["createTime"]
			result["proxy"] = vMap["proxy"]
			data = append(data, result)
		}
		instance["readonlyNode"] = data
	}
	values := GetSubSliceDByRep(allInstances, redisInstanceKeys)
	if err := dataSourceKscSave(d, "instances", []string{}, values); err != nil {
		return fmt.Errorf("error on save instance list, %s", err)
	}
	return nil
}

func queryAz(conn *kcsv1.Kcsv1) (map[string]string, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	result := make(map[string]string)
	az := []string{"1", "2"}
	action := "DescribeAvailabilityZones"
	readAz := make(map[string]interface{})
	for _, v := range az {
		readAz["Mode"] = v
		logger.Debug(logger.ReqFormat, action, readAz)
		if resp, err = conn.DescribeAvailabilityZones(&readAz); err != nil {
			return nil, fmt.Errorf("error on reading az")
		}
		logger.Debug(logger.RespFormat, action, readAz, *resp)
		set := (*resp)["AvailabilityZoneSet"].([]interface{})
		if len(set) == 0 {
			return result, nil
		}
		for _, v := range set {
			vv := v.(map[string]interface{})
			if vv["Region"].(string) == *conn.Config.Region {
				result[vv["AvailabilityZone"].(string)] = ""
			}
		}
	}
	logger.Info("region:", *conn.Config.Region, "az:", result)
	return result, nil
}
