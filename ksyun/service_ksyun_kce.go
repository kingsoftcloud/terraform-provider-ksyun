package ksyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/structor/v1/kce"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/structor/v1/kec"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// const skipCreate = true

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
		// logger.Debug("resp", "DescribeCluster", resp)
		clusterResults, err = getSdkValue("ClusterSet", *resp)
		if err != nil {
			return data, err
		}
		data = clusterResults.([]interface{})
		// logger.Debug("kce list", "123", data)
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
		"advanced_setting",
	}

	paraMap := map[string]interface{}{}
	for k, v := range nodeConfig {
		// 忽略部分字段
		if stringSliceContains(ignoreFields, k) {
			continue
		}
		k = Downline2Hump(k)
		// if v == nil {
		//	continue
		// }
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
			// for sgIdx, sgId := range sgIdList {
			//	if sgId != nil {
			//		paraMap[fmt.Sprintf("SecurityGroupId.%d", sgIdx+1)] = sgId
			//	}
			// }
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
		// logger.Debug("[%s] %+v", "error", err)
		return
	}
	// logger.Debug("[%s] %+v %s", "paraMap", paraMap, b)
	return string(b)
}

func formatKceClusterReq(d *schema.ResourceData, resource *schema.Resource) (createReq map[string]interface{}, err error) {
	transform := map[string]SdkReqTransform{
		// "node_config": {Ignore: true},
		"managed_cluster_multi_master": {
			Type: TransformListN,
		},
	}
	createReq, err = SdkRequestAutoMapping(d, resource, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	handleKecPara := func(createParams *map[string]interface{}, nodeConfigs []interface{}, index int, topKey string) int {

		for idx, nodeConfigSrc := range nodeConfigs {
			nodeConfig := nodeConfigSrc.(map[string]interface{})

			// logger.Debug("[%s] %d:%+v", "test", idx, nodeConfig)
			index++
			_idx := index
			(*createParams)[fmt.Sprintf("InstanceForNode.%d.NodeRole", _idx)] = nodeConfig["role"]
			(*createParams)[fmt.Sprintf("InstanceForNode.%d.NodeConfig.1.Para", _idx)] = formatKceInstancePara(nodeConfig)

			advancedSettingParams := map[string]interface{}{}

			if as, ok := nodeConfig["advanced_setting"]; ok {
				asSet := as.([]interface{})
				if len(asSet) == 0 {
					continue
				}

				advancedSetting := asSet[0].(map[string]interface{})

				for k, v := range advancedSetting {
					if _, ok := d.GetOk(fmt.Sprintf("%s.%d.advanced_setting.0.%s", topKey, idx, k)); !ok {
						continue
					}

					hump := Downline2Hump(k)
					formatAdvancedSettingParams(&advancedSettingParams, hump, v, true)
					logger.Debug("advanced_setting", "advanced_setting", advancedSettingParams)
				}

				handleAdvancedConfigWithPrefix(createParams, []interface{}{advancedSettingParams}, fmt.Sprintf("InstanceForNode.%d.NodeConfig", _idx), 0)
			}
		}

		return index
	}

	var (
		instanceIdx int
	)

	if nodeConfigs, ok := createReq["MasterConfig"]; ok {
		logger.Debug("[%s] %+v", "test", createReq)

		instanceIdx = handleKecPara(&createReq, nodeConfigs.([]interface{}), instanceIdx, "master_config")
	}
	delete(createReq, "MasterConfig")

	if workerCfg, ok := createReq["WorkerConfig"]; ok {
		logger.Debug("[%s] handles worker config %+v", "test", createReq)
		instanceIdx = handleKecPara(&createReq, workerCfg.([]interface{}), instanceIdx, "worker_config")
	}
	delete(createReq, "WorkerConfig")

	return
}

func (s *KceService) kceClusterStateRefreshFunc(clusterId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// return
		data, err := s.client.kceconn.DescribeCluster(&map[string]interface{}{
			"ClusterId": clusterId,
		})
		// logger.Debug("[%s] %+v %+v", "DescribeCluster", data, err)
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

func (s *KceService) CreateCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
	var createReq map[string]interface{}
	createReq, err = formatKceClusterReq(d, r)

	if err != nil {
		return
	}
	// logger.Debug("", "test", createReq)
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
			// resp = &map[string]interface{}{
			//	"ClusterId": "xxxx",
			// }
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
			if err != nil {
				return
			}
			// checkClusterState的err可以忽略(无论是否异常，都要加载一次集群数据用于生成结果)
			// 统一拿到最外层的create方法处理
			// err = s.ReadAndSetKceCluster(d, resource)

			// TODO: query all nodes and set to tf
			var (
				nodes []interface{}
			)

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				var (
					queryErr error
				)

				nodes, queryErr = s.getAllNodeWithFilter(clusterId.(string), nil)
				if queryErr != nil {
					return resource.RetryableError(queryErr)
				}
				return nil
			})
			if err != nil {
				return
			}

			var (
				workerNodeIds = make([]string, 0)
				masterNodeIds = make([]string, 0)
			)

			for _, nodeIf := range nodes {
				if node, ok := nodeIf.(map[string]interface{}); ok {
					kceInstance := kce.InstanceSet{}

					err := helper.MapstructureFiller(node, &kceInstance, "")
					if err != nil {
						return fmt.Errorf("convert master instance failed: %s", err)
					}
					nodeMap := map[string]interface{}{
						"instance_type": kceInstance.KecInstancePara.InstanceType,
						"image_id":      kceInstance.KecInstancePara.ImageID,
						"subnet_id":     kceInstance.KecInstancePara.SubnetID,
						"role":          kceInstance.InstanceRole,
					}

					networkInterface := kec.NetworkInterfaceSet{}
					for _, networkInterfaceIf := range kceInstance.KecInstancePara.NetworkInterfaceSet {
						if networkInterfaceIf.SubnetId == kceInstance.KecInstancePara.SubnetID {
							networkInterface = networkInterfaceIf
						}
					}

					securityIds := make([]string, len(networkInterface.SecurityGroupSet))
					for _, sg := range networkInterface.SecurityGroupSet {
						securityIds = append(securityIds, sg.SecurityGroupId)
					}
					nodeMap["security_group_id"] = securityIds

					hashcode := kceInstanceNodeHashFunc()(nodeMap)

					nodeHashAndId := AssembleIds(kceInstance.InstanceID, strconv.Itoa(hashcode))

					switch kceInstance.InstanceRole {
					case "Worker":
						workerNodeIds = append(workerNodeIds, nodeHashAndId)
					default:
						masterNodeIds = append(masterNodeIds, nodeHashAndId)
					}
				}

			}

			_ = d.Set("master_id_list", masterNodeIds)
			_ = d.Set("worker_id_list", workerNodeIds)
			return
		},
	}

	callbacks := []ApiCall{callback}
	err = ksyunApiCallNew(callbacks, d, s.client, true)

	return
}

func (s *KceService) getAllNodeWithFilter(clusterId string, filter map[string]interface{}) ([]interface{}, error) {

	var (
		resp                   *map[string]interface{}
		clusterInstanceResults interface{}
	)
	condition := map[string]interface{}{
		"ClusterId": clusterId,
		// "Filter.1.Name":    "instance-role",
		// "Filter.1.Value.1": "Master_Etcd",
		// "Filter.1.Value.2": "Master",
		// "Filter.1.Value.3": "Master_Etcd",
	}

	if filter != nil {
		idx := 0
		for k, v := range filter {
			idx++
			condition["Filter."+strconv.Itoa(idx)+".Name"] = k
			reflectType := reflect.TypeOf(v)
			switch reflectType.Kind() {
			case reflect.Slice:
				for i, vv := range v.([]interface{}) {
					condition["Filter."+strconv.Itoa(idx)+".Value."+strconv.Itoa(i+1)] = vv
				}
			case reflect.String, reflect.Int:
				condition["Filter."+strconv.Itoa(idx)+".Value."+strconv.Itoa(1)] = v
			default:
				continue
			}
		}
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
	// fmt.Println(d, resource)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		// 获取集群信息
		data, err := s.client.kceconn.DescribeCluster(&map[string]interface{}{
			"ClusterId": d.Id(),
		})
		// logger.Debug("[%s] %+v, %+v", "DescribeCluster", data, err)
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

		// read node
		err = s.readAndSetInstance(d, r)
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
	// return pageQuery(condition, "MaxResults", "Marker", 10, 0, func(condition map[string]interface{}) ([]interface{}, error) {
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
	// logger.Debug("kce instance images: %v %v", "DescribeInstanceImage", resp, err)
	results, err = getSdkValue("ImageSet", *resp)
	logger.Debug("kce instance images: %v %v", "results", results, err)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
	// })

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

// readAndSetInstance read the instance information from ksyun remote
// and make connection with the local config block with hashcode of machine types.
// there are eight steps to finish the function:
// 1. get the local instance id list.
// 2. calculate the hashcode of the machine type sequence in local tfstate.
// 3. according to local instance id, querying the kec information and role in kce cluster.
// 4. calculate the hashcode of the node information from ksyun remote.
// 5. deal with the node information from ksyun remote with sequence index and save the results to resources map.
// 6. deal with the rest of the node information and append it to tail of the resources map, if machine type was changed on the ksyun remote.
// 7. make the rest of nodes to relate with the local config block with hashcode of machine types.
// 8. save the resources map to the local tfstate.

// this function is so complicate, that we split hardly it to some small functions.
func (s *KceService) readAndSetInstance(d *schema.ResourceData, r *schema.Resource) (err error) {
	type nodeInfo struct {
		nodeMap  map[string]interface{}
		advanced map[string]interface{}
		role     string
	}
	type nodeInfoList struct {
		deal bool
		list []nodeInfo
	}

	var (
		mIDs []string
		wIDs []string

		masterInstances []interface{}
		workerInstances []interface{}

		nodeInfoMap = make(map[int]nodeInfoList, 0)

		instanceHashMap = make(map[int]string)

		masterInstanceList = make([]string, 0)
		workerInstanceList = make([]string, 0)
	)

	// get the local instance id list.
	mIDs, _ = helper.GetSchemaListWithString(d, "master_id_list")
	wIDs, _ = helper.GetSchemaListWithString(d, "worker_id_list")

	kecClient := KecService{client: s.client}

	// the function that query the kec instances from kec OpenAPI actions.
	queryKec := func(queryIds []string) ([]interface{}, error) {
		var (
			kecQuery        = map[string]interface{}{}
			retry           int
			infraErr        error
			masterInstances []interface{}
		)

		for idx, queryId := range queryIds {
			kecQuery[fmt.Sprintf("InstanceId.%d", idx+1)] = queryId
		}
	again:
		masterInstances, infraErr = kecClient.readKecInstances(kecQuery)
		if infraErr != nil && retry < 3 {
			retry++
			time.Sleep(2 * time.Second)
			goto again
		}
		return masterInstances, infraErr
	}

	// the function that query the role of the instance from kce cluster.
	queryRole := func(instanceID string) (*kce.InstanceSet, error) {
		filter := map[string]interface{}{
			"instance-id": instanceID,
		}
		nodes, infraErr := s.getAllNodeWithFilter(d.Id(), filter)
		if infraErr != nil {
			return nil, infraErr
		}
		var instance = &kce.InstanceSet{}
		_ = helper.MapstructureFiller(nodes[0], instance, "")
		return instance, nil
	}

	// the function that calculate the hashcode of the node information from ksyun remote.
	calculateNodeHashAndSave := func(node nodeInfo) {
		hashcode := kceInstanceNodeHashFunc()(node.nodeMap)

		if nodeList, ok := nodeInfoMap[hashcode]; !ok {
			nodeInfoMap[hashcode] = nodeInfoList{
				deal: false,
				list: []nodeInfo{node},
			}
		} else {
			nodeList.list = append(nodeList.list, node)
			nodeInfoMap[hashcode] = nodeList
		}

		var saveInstanceList []string
		switch node.role {
		case "Worker":
			saveInstanceList = workerInstanceList
		default:
			saveInstanceList = masterInstanceList
		}

		instanceHashId := generateInstanceIdWithHashcode(node.nodeMap)
		saveInstanceList = append(saveInstanceList, instanceHashId)
	}

	localMachineTypeSequenceM := map[int]int{}
	localMachineTypeSequenceW := map[int]int{}

	// the function that calculate the hashcode of the machine type.
	// and record the sequence's index in local config block with the hashcode.
	calculateMachineHash := func(blockKey string) {
		var (
			sm  = localMachineTypeSequenceM
			mtl []interface{}
		)

		switch blockKey {
		case "worker_config":
			sm = localMachineTypeSequenceW
		}

		mtlI, ok := d.GetOk(blockKey)
		if ok {
			mtl = mtlI.([]interface{})
		} else {
			return
		}

		for idx, m := range mtl {
			hashcode := kceInstanceNodeHashFunc()(m)
			sm[hashcode] = idx
			sm[idx] = hashcode
		}
	}

	// calculate the hashcode of the machine type sequence in local tfstate.
	// avoid the machine hashcode is empty.
	for _, blockKey := range []string{"master_config", "worker_config"} {
		calculateMachineHash(blockKey)
	}

	// according to local instance id, querying the kec information from kec OpenAPI actions.
	// and save the information to the local machine type sequence.
	if mIDs != nil && len(mIDs) > 0 {
		queryIds := make([]string, 0, len(mIDs))
		for _, mID := range mIDs {
			instanceIds := DisassembleIds(mID)
			hashcode, err := strconv.Atoi(instanceIds[1])
			if err != nil {
				return fmt.Errorf("convert hashcode failed: %s", err)
			}
			logger.Debug("read master instances", "hashcode", hashcode, "instance_id", instanceIds[0])
			instanceHashMap[hashcode] = instanceIds[0]
			queryIds = append(queryIds, instanceIds[0])
		}
		masterInstances, err = queryKec(queryIds)
		if err != nil {
			return fmt.Errorf("read master instances failed: %s", err)
		}
	}

	// according to local instance id, querying the kec information from kec OpenAPI actions.
	// and save the information to the local machine type sequence.
	if wIDs != nil && len(wIDs) > 0 {
		queryIds := make([]string, 0, len(wIDs))
		for _, wID := range wIDs {
			instanceIds := DisassembleIds(wID)
			hashcode, _ := strconv.Atoi(instanceIds[1])
			logger.Debug("read master instances", "hashcode", instanceIds[1], "instance_id", instanceIds[0])
			instanceHashMap[hashcode] = instanceIds[0]
			queryIds = append(queryIds, instanceIds[0])
		}
		workerInstances, err = queryKec(queryIds)
		if err != nil {
			return fmt.Errorf("read worker instances failed: %s", err)
		}
	}

	var (
		resourceMap = make(map[string]interface{})
	)

	// query the role and some information of master instances in kce cluster
	// and convert attributes of the master instances to tf schema map.
	for _, masterInstanceIf := range masterInstances {
		// handles master instance and set to data resources
		masterInstance := masterInstanceIf.(map[string]interface{})
		masterSaveMap, err := convertInstanceToMapForSchema(masterInstance)
		if err != nil {
			return fmt.Errorf("convert master instance failed: %s", err)
		}
		instanceId := masterInstance["InstanceId"].(string)
		// masterSaveMap["count"] = 1
		role, queryErr := queryRole(instanceId)
		if queryErr != nil {
			return fmt.Errorf("query %s role failed: %s", instanceId, err)
		}
		masterSaveMap["role"] = role.InstanceRole

		advanced := handleAdvancedSetting2Map(*role.AdvancedSetting)
		node := nodeInfo{
			nodeMap:  masterSaveMap,
			advanced: advanced,
			role:     role.InstanceRole,
		}
		calculateNodeHashAndSave(node)
	}

	// query the worker instances and convert attributes of the worker instances to tf schema map.
	for _, workerInstanceIf := range workerInstances {
		// handles worker instance and set to data resources
		workerInstance := workerInstanceIf.(map[string]interface{})
		workerSaveMap, err := convertInstanceToMapForSchema(workerInstance)
		if err != nil {
			return fmt.Errorf("convert master instance failed: %s", err)
		}
		// workerSaveMap["count"] = 1
		instanceId := workerInstance["InstanceId"].(string)
		role, queryErr := queryRole(instanceId)
		if queryErr != nil {
			return fmt.Errorf("query %s role failed: %s", instanceId, err)
		}
		workerSaveMap["role"] = role.InstanceRole

		advanced := handleAdvancedSetting2Map(*role.AdvancedSetting)
		node := nodeInfo{
			nodeMap:  workerSaveMap,
			advanced: advanced,
			role:     role.InstanceRole,
		}
		calculateNodeHashAndSave(node)
	}

	// deal with the rank of nodes information from ksyun remote with hashcode about machine type in local tfstate.
	// if the hashcode between the local machine type sequence and the remote machine type sequence is equal,
	// set the remote information to the local machine type sequence with the same index.
	if nodeInfoMap != nil && len(nodeInfoMap) > 0 {
		blocks := []string{"master_config", "worker_config"}
		for _, block := range blocks {
			localBlock, ok := helper.GetSchemaMapListWithKey(d, block)
			if !ok {
				continue
			}
			results := make([]interface{}, len(localBlock), len(localBlock))
			for idx, nodeLocalInfo := range localBlock {
				schemeAdvancedKey := fmt.Sprintf("%s.%d.advanced_setting", block, idx)
				hashcode := kceInstanceNodeHashFunc()(nodeLocalInfo)

				var _nodeInfoList []nodeInfo
				if nl, ok := nodeInfoMap[hashcode]; !ok {
					continue
				} else {
					// if the information of nodes from remote has the identical hashcode with local machine type sequence,
					// mark the nodeInfoList as dealt.
					_nodeInfoList = nl.list
					nl.deal = true
					nodeInfoMap[hashcode] = nl
				}

				nodeMapList := make([]map[string]interface{}, 0, len(_nodeInfoList))
				nodeAdvancedMapList := make([]map[string]interface{}, 0, len(_nodeInfoList))
				for _, _nodeInfo := range _nodeInfoList {
					nodeMapList = append(nodeMapList, _nodeInfo.nodeMap)
					nodeAdvancedMapList = append(nodeAdvancedMapList, _nodeInfo.advanced)
				}
				diffMap := helper.GetDiffMap(nodeLocalInfo, nodeMapList...)

				localAdvanced, ok := helper.GetSchemaListHeadMap(d, schemeAdvancedKey)
				if ok {
					advancedDiff := helper.GetDiffMap(localAdvanced, nodeAdvancedMapList...)
					diffMap["advanced_setting"] = []interface{}{advancedDiff}
				}
				diffMap["hashcode"] = hashcode
				diffMap["count"] = len(_nodeInfoList)
				results[idx] = diffMap
			}
			resourceMap[block] = results
		}

		// deal with the rest nodes of the nodeInfoMap
		// the rest nodes are the nodes that are not in the local machine type sequence.
		// so, we just append the rest nodes to the resourceMap.
		for hashcode, _nodeInfoList := range nodeInfoMap {
			if _nodeInfoList.deal {
				continue
			}

			var (
				block string
			)
			theFirstNode := _nodeInfoList.list[0]
			theFirstMap := theFirstNode.nodeMap
			switch theFirstNode.role {
			case "Worker":
				block = "worker_config"
			default:
				block = "master_config"
			}
			theFirstMap["hashcode"] = hashcode
			theFirstMap["count"] = len(_nodeInfoList.list)
			if theFirstNode.advanced != nil && len(theFirstNode.advanced) > 0 {
				theFirstMap["advanced_setting"] = []interface{}{theFirstNode.advanced}
			}

			resourceMap[block] = append(resourceMap[block].([]interface{}), theFirstNode.nodeMap)
		}

		// omit the element that is for occupied.
		// compact the resourceMap
		findCandidateFunc := func(instanceId string) (int, bool) {
			for hashcode, _nodeInfoList := range nodeInfoMap {
				for _, _nodeInfo := range _nodeInfoList.list {
					if _nodeInfo.nodeMap["instance_id"] == instanceId {
						return hashcode, true
					}
				}
			}
			return 0, false
		}

		// find the candidate block with the hashcode
		findSaveMapWithHashcodeFunc := func(hashcode int, candidates []interface{}) (map[string]interface{}, int) {
			for idx, candidate := range candidates {
				if helper.IsEmpty(candidate) {
					continue
				}
				candidateMap := candidate.(map[string]interface{})
				if candidateMap["hashcode"] == hashcode {
					return candidateMap, idx
				}
			}
			return nil, 0
		}

		// find a candidate instance block to replace the old block for occupy the position.
		// if old block is empty, other word, the original machine type was changed and that it not exists in the remote.
		// so we need to find a candidate block to replace the old block. In order to keep the position for make a diff by terraform.
		for k, ss := range resourceMap {
			results := make([]interface{}, 0, len(ss.([]interface{})))
			// stateSlice := d.Get(k).([]interface{})

			sm := localMachineTypeSequenceM
			if k == "worker_config" {
				sm = localMachineTypeSequenceW
			}

			sss := ss.([]interface{})
			stateLength := len(sss)
			for idx := 0; idx < stateLength; idx++ {
				if idx >= stateLength {
					break
				}

				block := sss[idx]
				if helper.IsEmpty(block) {
					oldHashcode := sm[idx]
					candidateInstanceId := instanceHashMap[oldHashcode]
					hashcode, ok := findCandidateFunc(candidateInstanceId)
					logger.Debug("readAndSetInstance", "hashcode", hashcode, "candidateInstanceId", candidateInstanceId)
					if !ok {
						continue
					}
					candidate, dealIdx := findSaveMapWithHashcodeFunc(hashcode, sss)
					logger.Debug("readAndSetInstance", "candidate", candidate, "dealIdx", dealIdx)

					if candidate != nil {
						// handle the taints from instance block with the old hashcode
						// because the new candidate block has not the taints from the remote.
						// so we need to keep the taints from the old block.

						oldBlockSchemaKey := fmt.Sprintf("%s.%d.advanced_setting.0.taints", k, idx)
						oldTaintsBlock, blockExist := d.GetOk(oldBlockSchemaKey)
						logger.Debug("readAndSetInstance", "oldTaintsBlock", oldTaintsBlock, "blockExist", blockExist)
						logger.Debug("readAndSetInstance", "oldBlockSchemaKey", oldBlockSchemaKey)
						if blockExist {
							blockASs := candidate["advanced_setting"].([]interface{})
							blockAS := blockASs[0].(map[string]interface{})
							blockAS["taints"] = oldTaintsBlock
							candidate["advanced_setting"] = []interface{}{blockAS}
						}
						results = append(results, candidate)
					}

					// remove the element that is occupied
					for i := dealIdx; i < len(sss)-1; i++ {
						sss[i] = sss[i+1]
					}
					stateLength--
				} else {
					results = append(results, block)
				}
			}
			resourceMap[k] = results
		}
	}

	if len(masterInstanceList) > 0 {
		_ = d.Set("master_id_list", masterInstanceList)
	}
	if len(workerInstanceList) > 0 {
		_ = d.Set("worker_id_list", workerInstanceList)
	}
	SdkResponseAutoResourceData(d, r, resourceMap, nil)

	// todo: Done
	// 把机器列表格式化一组字符串，然后将master_config也格式化成一组字符串，
	// 然后把机器串匹配master_config，能匹配上就累加数字，如果最终有差异，就成为diff

	return
}

// // 将master机器列表
// func getMasterUniqKeyFromInstances() {}
// func getMasterUniqKeyFromConfigs()   {}

func convertInstanceToMapForSchema(insResp map[string]interface{}) (map[string]interface{}, error) {
	// logger.Debug("convertInstanceToMapForSchema", "ins", ins)
	insMap := map[string]interface{}{}

	schemaMap := instanceConfig()

	ignoreFields := []string{"advanced_setting"}

	// handle the top level fields
	for k, v := range insResp {
		// logger.Debug("convertInstanceToMapForSchema", "k", k, "v", v)
		underK := helper.Hump2Underline(k)
		if checkValueInSlice(ignoreFields, underK) {
			continue
		}
		if schemaMap[underK] == nil || schemaMap[underK].Elem != nil {
			continue
		}
		insMap[underK] = v
		logger.Debug("convertInstanceToMapForSchema", "k", k, "v", v)
	}

	// handle the special fields
	kecInstance := &kec.Instance{}
	err := helper.MapstructureFiller(insResp, kecInstance, "")
	if err != nil {
		return nil, fmt.Errorf("convert worker instance failed: %s", err)
	}

	// InstanceConfigure
	insMap["data_disk_gb"] = kecInstance.InstanceConfigure.DataDiskGb

	// InstanceState
	insMap["instance_status"] = kecInstance.InstanceState.Name

	// NetworkInterfaceSet
	for _, ni := range kecInstance.NetworkInterfaceSet {
		if ni.NetworkInterfaceType == "primary" {
			insMap["private_ip_address"] = ni.PrivateIpAddress
			insMap["subnet_id"] = ni.SubnetId
			insMap["network_interface_id"] = ni.NetworkInterfaceId
			insMap["vpc_id"] = ni.VpcId

			sgIds := make([]string, 0, len(ni.SecurityGroupSet))
			sgIds = append(sgIds, ni.SecurityGroupSet[0].SecurityGroupId)
			insMap["security_group_id"] = sgIds
		}

	}

	systemDiskMap := make(map[string]interface{}, 2)
	_ = helper.MapstructureFiller(kecInstance.SystemDisk, &systemDiskMap, "tf-schema")
	insMap["system_disk"] = []interface{}{systemDiskMap}

	return insMap, nil
}

func handleAdvancedSetting2Map(advancedSetting kce.AdvancedSetting) map[string]interface{} {
	advanced := make(map[string]interface{})
	_ = helper.MapstructureFiller(advancedSetting, &advanced, "tf-schema")
	if !helper.IsEmpty(advancedSetting.ExtraArg) {
		advanced["extra_arg"] = advancedSetting.ExtraArg.Kubelet
	}
	return advanced
}

func localNodeCfgHashcode(d *schema.ResourceData, field string) int {
	localConfig, _ := helper.GetSchemaListHeadMap(d, field)
	return kceInstanceNodeHashFunc()(localConfig)
}

func generateInstanceIdWithHashcode(instance map[string]interface{}) string {
	hashcode := kceInstanceNodeHashFunc()(instance)
	instanceId := instance["instance_id"].(string)
	return AssembleIds(instanceId, strconv.Itoa(hashcode))
}
