package ksyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

//const skipCreate = true

type KceService struct {
	client *KsyunClient
}

// 获取kce列表
func (s *KceService) readKceClusters(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp           *map[string]interface{}
		clusterResults interface{}
	)

	return pageQuery(condition, "MaxResults", "Marker", 10, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kceconn

		if condition == nil {
			resp, err = conn.DescribeCluster(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeCluster(&condition)
			if err != nil {
				return data, err
			}
		}
		//logger.Debug("resp", "DescribeCluster", resp)
		clusterResults, err = getSdkValue("ClusterSet", *resp)
		if err != nil {
			return data, err
		}
		data = clusterResults.([]interface{})
		//logger.Debug("kce list", "123", data)
		return data, err
	})

}

func (s *KceService) ReadAndSetKceClusters(d *schema.ResourceData, r *schema.Resource) (err error) {

	transform := map[string]SdkReqTransform{
		"cluster_id": {
			mapping: "ClusterId",
			Type:    TransformDefault,
		},
		"search": {
			mapping: "Search",
			Type:    TransformDefault,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readKceClusters(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "ClusterId",
		nameField:   "ClusterName",
		targetField: "cluster_set",
		extra: map[string]SdkResponseMapping{
			"EnableKMSE": {
				Field: "enable_kmse",
			},
		},
	})
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case string:
		return v.(string) == ""
	}
	return false
}

func formatKceInstancePara(nodeConfig map[string]interface{}) (para string) {

	// todo: 网卡管理是缺失的

	ignoreFields := []string{
		// tag这个忽略的设置有点问题，kec的terraform是单独调了tag接口，但实际上主机的接口是支持tag的
		"instance_status", "force_delete", "force_reinstall_system",
		"extension_network_interface",
		"tags",
		"role",
	}

	paraMap := map[string]interface{}{}
	for k, v := range nodeConfig {
		// 忽略部分字段
		if stringSliceContains(ignoreFields, k) {
			continue
		}
		k = Downline2Hump(k)
		//if v == nil {
		//	continue
		//}
		if isEmpty(v) {
			continue
		}
		switch k {
		case "Count":
			paraMap["MinCount"] = v
			paraMap["MaxCount"] = v
			break
		case "KeyId":
			keyIdList := v.(*schema.Set).List()
			for keyIdx, keyId := range keyIdList {
				if keyId != nil {
					paraMap[fmt.Sprintf("KeyId.%d", keyIdx+1)] = keyId
				}
			}
			break
		case "SystemDisk":
			for _, diskSrc := range v.([]interface{}) {
				disk := diskSrc.(map[string]interface{})
				if disk["disk_type"] != nil {
					paraMap["System.DiskType"] = disk["disk_type"]
				}
				if disk["disk_size"] != nil {
					paraMap["System.DiskSize"] = disk["disk_size"]
				}
			}
			break
		// todo: 容器和主机不一致，不支持.N, 但是返回值是数组，所以传参保持list格式，取值先取第一个
		case "SecurityGroupId":
			sgIdList := v.(*schema.Set).List()
			if len(sgIdList) > 0 {
				paraMap["SecurityGroupId"] = sgIdList[0]
			}
			//for sgIdx, sgId := range sgIdList {
			//	if sgId != nil {
			//		paraMap[fmt.Sprintf("SecurityGroupId.%d", sgIdx+1)] = sgId
			//	}
			//}
			break
		case "DataDisks":
			for diskIdx, diskSrc := range v.([]interface{}) {
				disk := diskSrc.(map[string]interface{})
				if disk["disk_type"] != nil {
					paraMap[fmt.Sprintf("DataDisk.%d.Type", diskIdx+1)] = disk["disk_type"]
				}
				if disk["disk_size"] != nil {
					paraMap[fmt.Sprintf("DataDisk.%d.Size", diskIdx+1)] = disk["disk_size"]
				}
				if disk["delete_with_instance"] != nil {
					paraMap[fmt.Sprintf("DataDisk.%d.DeleteWithInstance", diskIdx+1)] = disk["delete_with_instance"]
				}
				if disk["disk_snapshot_id"] != nil {
					paraMap[fmt.Sprintf("DataDisk.%d.SnapshotId", diskIdx+1)] = disk["disk_snapshot_id"]
				}
			}
			break
		default:
			paraMap[k] = v
		}
	}
	b, err := json.Marshal(paraMap)
	if err != nil {
		//logger.Debug("[%s] %+v", "error", err)
		return
	}
	//logger.Debug("[%s] %+v %s", "paraMap", paraMap, b)
	return string(b)
}

func formatKceClusterReq(d *schema.ResourceData, resource *schema.Resource) (createReq map[string]interface{}, err error) {

	transform := map[string]SdkReqTransform{
		//"node_config": {Ignore: true},
	}
	createReq, err = SdkRequestAutoMapping(d, resource, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if nodeConfigs, ok := createReq["MasterConfig"]; ok {

		for idx, nodeConfigSrc := range nodeConfigs.([]interface{}) {
			nodeConfig := nodeConfigSrc.(map[string]interface{})

			//logger.Debug("[%s] %d:%+v", "test", idx, nodeConfig)
			createReq[fmt.Sprintf("InstanceForNode.%d.NodeRole", idx+1)] = nodeConfig["role"]
			createReq[fmt.Sprintf("InstanceForNode.%d.NodeConfig.1.Para", idx+1)] = formatKceInstancePara(nodeConfig)
		}
		//logger.Debug("[%s] %+v", "test", createReq)

	} else {
		err = fmt.Errorf("node_config is required")
	}

	delete(createReq, "MasterConfig")
	//for k, v := range createReq {
	//	logger.Debug("[%s] %s:%v", "createReq", k, v)
	//}

	return
}

func (s *KceService) kceClusterStateRefreshFunc(clusterId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		//return
		data, err := s.client.kceconn.DescribeCluster(&map[string]interface{}{
			"ClusterId": clusterId,
		})
		//logger.Debug("[%s] %+v %+v", "DescribeCluster", data, err)
		if err != nil {
			return nil, "", err
		}
		status, err := getSdkValue("ClusterSet.0.Status", *data)
		logger.Debug("[%s] %+v %+v %+v", "DescribeCluster", data, err, status)
		if err != nil {
			return nil, "", err
		}
		if stringSliceContains(failStates, status.(string)) {
			return nil, "", fmt.Errorf("instance status  error, status:%v", status)
		}

		return data, status.(string), nil
	}
}

func (s *KceService) checkClusterState(clusterId string, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{},
		Target:       target,
		Refresh:      s.kceClusterStateRefreshFunc(clusterId, []string{"error", "NotReady"}),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
		MinTimeout:   1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *KceService) CreateCluster(d *schema.ResourceData, resource *schema.Resource) (err error) {
	var createReq map[string]interface{}
	createReq, err = formatKceClusterReq(d, resource)

	if err != nil {
		return
	}
	//logger.Debug("", "test", createReq)
	callback := ApiCall{
		param:  &createReq,
		action: "CreateCluster",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			b, e := json.Marshal(call.param)
			logger.Debug("[%s] err: %+v", call.action, e)
			logger.Debug(logger.RespFormat, call.action, string(b))

			// XXX: create
			conn := client.kcev2conn
			resp, err = conn.CreateCluster(call.param)

			// XXX debug: skip create
			//resp = &map[string]interface{}{
			//	"ClusterId": "xxxx",
			//}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var clusterId interface{}
			if resp != nil {
				clusterId, err = getSdkValue("ClusterId", *resp)
				if err != nil {
					return
				}
				d.SetId(clusterId.(string))
			}
			err = s.checkClusterState(clusterId.(string), []string{"running"}, d.Timeout(schema.TimeoutCreate))
			// checkClusterState的err可以忽略(无论是否异常，都要加载一次集群数据用于生成结果)
			err = s.ReadAndSetKceCluster(d, resource)
			return
		},
	}

	callbacks := []ApiCall{callback}
	err = ksyunApiCallNew(callbacks, d, s.client, true)

	return
}

func (s *KceService) getAllMasters(clusterId string) ([]interface{}, error) {

	var (
		resp                   *map[string]interface{}
		clusterInstanceResults interface{}
	)
	condition := map[string]interface{}{
		"ClusterId":        clusterId,
		"Filter.1.Name":    "instance-role",
		"Filter.1.Value.1": "Master_Etcd",
		"Filter.1.Value.2": "Master",
		"Filter.1.Value.3": "Master_Etcd",
	}
	return pageQuery(condition, "MaxResults", "Marker", 10, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kceconn
		var list []interface{}
		var err error
		if condition == nil {
			resp, err = conn.DescribeClusterInstance(nil)
			if err != nil {
				return list, err
			}
		} else {
			resp, err = conn.DescribeClusterInstance(&condition)
			if err != nil {
				return list, err
			}
		}
		clusterInstanceResults, err = getSdkValue("InstanceSet", *resp)
		if err != nil {
			return list, err
		}
		list = clusterInstanceResults.([]interface{})
		return list, err
	})

}

func (s *KceService) DeleteKceCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
	req := make(map[string]interface{})
	req["ClusterId"] = d.Id()
	_, err = s.client.kceconn.DeleteCluster(&req)
	if err != nil {
		return
	}
	var data []interface{}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, "DeleteCluster", req)

		data, err = s.readKceClusters(req)
		if len(data) == 0 {
			return nil
		}

		status, err := getSdkValue("Status", data[0])
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if status.(string) != "deleting" {
			err = errors.New("cluster status not available")
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(errors.New("deleting"))

	})
}

func (s *KceService) UpdateCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
	params := map[string]interface{}{
		"ClusterId": d.Get("cluster_id"),
	}
	if d.HasChange("cluster_name") {
		params["ClusterName"] = d.Get("cluster_name")
	}
	if d.HasChange("cluster_desc") {
		params["ClusterDesc"] = d.Get("cluster_desc")
	}
	_, err = s.client.kceconn.ModifyClusterInfo(&params)
	return
}

func (s *KceService) ReadAndSetKceCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
	//fmt.Println(d, resource)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		// 获取集群信息
		data, err := s.client.kceconn.DescribeCluster(&map[string]interface{}{
			"ClusterId": d.Id(),
		})
		//logger.Debug("[%s] %+v, %+v", "DescribeCluster", data, err)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if data == nil {
			return resource.NonRetryableError(fmt.Errorf("cluster not found"))
		}
		clusterSet := (*data)["ClusterSet"].([]interface{})

		if len(clusterSet) <= 0 {
			return resource.NonRetryableError(fmt.Errorf("cluster not found"))
		}
		clusterInfo := clusterSet[0].(map[string]interface{})

		extra := map[string]SdkResponseMapping{}
		SdkResponseAutoResourceData(d, r, clusterInfo, extra)
		//
		err = s.readAndSetMasters(d, r)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func (s *KceService) readKceInstanceImages(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	//return pageQuery(condition, "MaxResults", "Marker", 10, 0, func(condition map[string]interface{}) ([]interface{}, error) {
	conn := s.client.kceconn

	if condition == nil {
		resp, err = conn.DescribeInstanceImage(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeInstanceImage(&condition)
		if err != nil {
			return data, err
		}
	}
	//logger.Debug("kce instance images: %v %v", "DescribeInstanceImage", resp, err)
	results, err = getSdkValue("ImageSet", *resp)
	logger.Debug("kce instance images: %v %v", "results", results, err)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
	//})

}

func (s *KceService) ReadAndSetKceInstanceImages(d *schema.ResourceData, r *schema.Resource) (err error) {
	var data []interface{}
	data, err = s.readKceInstanceImages(nil)
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "ImageId",
		nameField:   "ImageName",
		targetField: "image_set",
	})
}

// 把tf里master的配置，和master机器列表对应起来
func (s *KceService) readAndSetMasters(d *schema.ResourceData, r *schema.Resource) (err error) {

	//masterConfigSrc := d.Get("master_config")
	//logger.Debug("readAndSetMasters", "master_config", masterConfigSrc)
	//if masterConfigList, ok := masterConfigSrc.([]interface{}); ok {
	//	for _, configItem := range masterConfigList {
	//		logger.Debug("readAndSetMasters", "master_config_item", configItem.(map[string]interface{}))
	//	}
	//}
	//
	//var instances []interface{}
	//instances, err = s.getAllMasters(d.Id())
	//
	//logger.Debug("readAndSetMasters", "instances", instances)
	//for _, itemSrc := range instances {
	//	if item, ok := itemSrc.(map[string]interface{}); ok {
	//		instanceRole := item["InstanceRole"]
	//	}
	//}
	//d.Get("master_config")

	// todo:
	// 把机器列表格式化一组字符串，然后将master_config也格式化成一组字符串，
	// 然后把机器串匹配master_config，能匹配上就累加数字，如果最终有差异，就成为diff

	return
}

//// 将master机器列表
//func getMasterUniqKeyFromInstances() {}
//func getMasterUniqKeyFromConfigs()   {}
