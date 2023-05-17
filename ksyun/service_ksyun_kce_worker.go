package ksyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"strings"
	"time"
)

type KceWorkerService struct {
	client *KsyunClient
}

type KceExistedInstance struct {
	Available         bool   `json:"Available,omitempty"`
	InstanceId        string `json:"InstanceId,omitempty"`
	UnavailableReason string `json:"UnavailableReason,omitempty"`
}
type KceAddExistedInstance struct {
	InstanceId string `json:"InstanceId,omitempty"`
	Reason     string `json:"Reason,omitempty"`
	Return     bool   `json:"Return,omitempty"`
}

func (s *KceWorkerService) addInstanceStateRefreshFunc(d *schema.ResourceData, instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err  error
			data *map[string]interface{}
		)

		data, err = s.client.kceconn.DescribeClusterInstance(&map[string]interface{}{
			"ClusterId":        d.Get("cluster_id"),
			"Filter.1.Name":    "instance-id",
			"Filter.1.Value.1": instanceId,
		})

		if err != nil {
			return nil, "", err
		}

		var status interface{}
		status, err = getSdkValue("InstanceSet.0.InstanceStatus", *data)
		//logger.Debug("test", "addInstanceStateRefreshFunc", status)
		if err != nil {
			return nil, "", err
		}

		for _, v := range failStates {
			if v == status.(string) {
				return nil, "", fmt.Errorf("instance status  error, status:%v", status)
			}
		}
		return data, status.(string), nil
	}
}

func (s *KceWorkerService) checkAddInstanceProgress(d *schema.ResourceData, instanceId string, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{},
		Target:       target,
		Refresh:      s.addInstanceStateRefreshFunc(d, instanceId, []string{"error"}),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
		MinTimeout:   1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func formatAdvancedSettingParams(params *map[string]interface{}, currentParamKey string, currentParamValue interface{}, isNew bool) {
	switch currentParamKey {
	case "DataDisk":
		value := currentParamValue.(*schema.Set)
		if value.Len() == 0 {
			return
		}
		dataDiskItem := value.List()[0].(map[string]interface{})
		for k, v := range dataDiskItem {
			(*params)["DataDisk."+Downline2Hump(k)] = v
		}
		return
	case "Label":
		value := currentParamValue.([]interface{})
		if len(value) == 0 {
			return
		}
		for idx, item := range value {
			(*params)[fmt.Sprintf("Label.%d.Key", idx+1)] = item.(map[string]interface{})["key"]
			(*params)[fmt.Sprintf("Label.%d.Value", idx+1)] = item.(map[string]interface{})["value"]
		}
		return
	case "ExtraArg":
		value := currentParamValue.([]interface{})
		if len(value) == 0 {
			return
		}
		for idx, item := range value {
			if isNew || item.(string) != "" {
				(*params)[fmt.Sprintf("ExtraArg.%d", idx+1)] = item
			}
		}
		return
	default:
		if isNew {
			if valueStr, ok := currentParamValue.(string); ok {
				if valueStr != "" {
					(*params)[currentParamKey] = currentParamValue
				}
			} else {
				(*params)[currentParamKey] = currentParamValue
			}
		} else {
			(*params)[currentParamKey] = currentParamValue
		}
		return
	}
	//logger.Debug("AddWorker", "AddWorker", currentParamKey, currentParamValue)
}

func (s *KceWorkerService) AddWorker(d *schema.ResourceData, r *schema.Resource) (err error) {
	clusterId := d.Get("cluster_id")
	instanceId := d.Get("instance_id")
	imageId := d.Get("image_id")

	//as, ok := d.GetOk("advanced_setting")
	//logger.Debug("AddWorker", "advanced_setting", as, ok)
	//dd, ok := d.GetOk("data_disk")
	//logger.Debug("AddWorker", "data_disk", dd, ok)
	//lb, ok := d.GetOk("data_disk.0.auto_format_and_mount")
	//logger.Debug("AddWorker", "label", lb, ok)
	//return

	// 查询是否可以移入
	var resp *map[string]interface{}
	resp, err = s.client.kceconn.DescribeExistedInstances(&map[string]interface{}{
		"ClusterId":    clusterId,
		"InstanceId.1": instanceId,
	})
	if err != nil {
		return
	}

	var nodes []KceExistedInstance
	err = transInterfaceToStruct((*resp)["InstanceSet"], &nodes)
	if err != nil {
		return
	}
	if len(nodes) == 0 {
		err = errors.New("instance not exists")
		return
	}
	if !nodes[0].Available {
		if !strings.Contains(nodes[0].UnavailableReason, "The instance is not in the stopped state") {
			err = errors.New(nodes[0].UnavailableReason)
			return
		}
		var callbacks []ApiCall
		kecService := KecService{s.client}
		dst := &schema.ResourceData{}
		dst.SetId(instanceId.(string))
		var stopFunc ApiCall
		stopFunc, err = kecService.stopKecInstance(dst)
		if err != nil {
			return
		}
		callbacks = append(callbacks, stopFunc)
		err = ksyunApiCallNew(callbacks, dst, s.client, true)
		if err != nil {
			return
		}
	}

	//[map[InstanceId:0b0f6f62-25ef-478f-9576-20a93c11e5dc Reason:The instance modify image fail Return:false]

	// inti params
	params := map[string]interface{}{
		"ClusterId":                        clusterId, //"cd19855c-ed77-447a-9d4f-0fb6f7707df6",
		"ExistedInstanceKecSet.1.NodeRole": "worker",
	}

	// 整理 KecPara 参数
	kecPara := map[string]interface{}{
		"InstanceId": instanceId,
		"ImageId":    imageId,
	}
	if instancePassword, ok := d.GetOk("instance_password"); ok {
		kecPara["InstancePassword"] = instancePassword
	}
	if keyId, ok := d.GetOk("key_id"); ok {
		kecPara["KeyId.1"] = keyId
	}
	kecParaBytes, _ := json.Marshal(&kecPara)
	params["ExistedInstanceKecSet.1.KecPara.1"] = string(kecParaBytes)

	// 整理 AdvanceSetting 参数

	advancedFields := []string{
		"data_disk",
		"container_runtime",
		"docker_path",
		"container_path",
		"user_script",
		"pre_user_script",
		"schedulable",
		"label",
		"extra_arg",
		"container_log_max_size",
		"container_log_max_files",
	}
	advancedSettingParams := map[string]interface{}{}
	for _, f := range advancedFields {
		v, ok := d.GetOk(f)
		logger.Debug("AddWorker", f, v, ok)
		if ok {
			f = Downline2Hump(f)
			formatAdvancedSettingParams(&advancedSettingParams, f, v, true)
			logger.Debug("advanced_setting", "advanced_setting", advancedSettingParams)
		} else {
			logger.Debug("AddWorker", "no advanced_setting", ok)
		}
	}
	for k, v := range advancedSettingParams {
		params[fmt.Sprintf("ExistedInstanceKecSet.1.AdvancedSetting.%s", k)] = v
	}

	paramsBytes, parseJsonErr := json.Marshal(params)
	logger.Debug("AddWorker", "params json", string(paramsBytes), parseJsonErr)

	resp, err = s.client.kceconn.AddExistedInstances(&params)
	if err != nil {
		return
	}
	var addExistedInstances []KceAddExistedInstance
	err = transInterfaceToStruct((*resp)["InstanceSet"], &addExistedInstances)
	if err != nil {
		return
	}
	if len(addExistedInstances) == 0 {
		err = errors.New("addExistedInstances response error")
		return
	}
	if !addExistedInstances[0].Return {
		err = errors.New(addExistedInstances[0].Reason)
		return
	}

	// InstanceStatus:normal
	err = s.checkAddInstanceProgress(d, instanceId.(string), []string{"normal"}, d.Timeout(schema.TimeoutUpdate))

	d.SetId(clusterId.(string) + ":" + instanceId.(string))
	return
}

func (s *KceWorkerService) DeleteKceWorker(d *schema.ResourceData, r *schema.Resource) (err error) {
	var resp *map[string]interface{}
	resp, err = s.client.kceconn.DeleteClusterInstances(&map[string]interface{}{
		"ClusterId":          d.Get("cluster_id"),
		"InstanceDeleteMode": "Remove",
		"InstanceId.1":       d.Get("instance_id"),
	})
	if err != nil {
		return
	}

	var deleteKceInstances []KceAddExistedInstance
	err = transInterfaceToStruct((*resp)["InstanceSet"], &deleteKceInstances)
	if err != nil {
		return
	}
	if len(deleteKceInstances) == 0 {
		err = errors.New("unknown error")
		return
	}
	if !deleteKceInstances[0].Return {
		err = errors.New(deleteKceInstances[0].Reason)
		return
	}
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		var data *map[string]interface{}
		data, err = s.client.kceconn.DescribeClusterInstance(&map[string]interface{}{
			"ClusterId":        d.Get("cluster_id"),
			"Filter.1.Name":    "instance-id",
			"Filter.1.Value.1": d.Get("instance_id"),
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		instanceSetSrc := (*data)["InstanceSet"]
		if instanceSetSrc == nil {
			return nil
		}
		instanceSet := instanceSetSrc.([]interface{})
		if len(instanceSet) == 0 {
			return nil
		}
		var status interface{}
		status, err = getSdkValue("InstanceStatus", instanceSet[0])
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if status.(string) != "deleting" {
			return resource.NonRetryableError(errors.New("instance status not available"))
		}
		return resource.RetryableError(errors.New("deleting"))
	})

}

func updateResourceData(d *schema.ResourceData, fieldName string, fieldValue interface{}) {
	if fieldValue == nil {
		return
	}
	fieldNameDownline := Hump2Downline(fieldName)
	switch fieldName {
	case "DataDisk":
		item := fieldValue.(map[string]interface{})
		itemValueFormatted := map[string]interface{}{}
		for itemKey, itemValue := range item {
			itemValueFormatted[Hump2Downline(itemKey)] = itemValue
		}
		logger.Debug("DataDisk", "DataDisk", fieldNameDownline, itemValueFormatted)
		d.Set(fieldNameDownline, []interface{}{itemValueFormatted})
		return
	case "Label":
		items := fieldValue.([]interface{})
		formattedValue := []map[string]interface{}{}
		for _, item := range items {
			itemValueFormatted := map[string]interface{}{}
			for itemKey, itemValue := range item.(map[string]interface{}) {
				itemValueFormatted[Hump2Downline(itemKey)] = itemValue
			}
			formattedValue = append(formattedValue, itemValueFormatted)
		}
		d.Set(fieldNameDownline, formattedValue)
		return
	case "ExtraArg":
		d.Set(fieldNameDownline, fieldValue)
		return
	case "ContainerLogMaxFiles", "ContainerLogMaxSize":
		if vStr, ok := fieldValue.(string); ok {
			vInt, err := strconv.Atoi(vStr)
			if err == nil {
				d.Set(fieldNameDownline, vInt)
			}
		}
		return
	default:
		d.Set(fieldNameDownline, fieldValue)
		return
	}
}

// todo
// label和taint有openapi的读接口，但是写入操作是通过更新node的yaml实现的
// 具体接口包括：
// DeleteVirtualNode
// patchResourceYaml
func (s *KceWorkerService) readAndSetLabels(d *schema.ResourceData) (err error) {
	//s.client.kceconn.DescribeNodeLabels()
	// 获取label列表
	return
}

func (s *KceWorkerService) ReadAndSetWorker(d *schema.ResourceData, r *schema.Resource) (err error) {
	var data *map[string]interface{}
	id := d.Id()
	ids := strings.Split(id, ":")
	clusterId := ids[0]
	instanceId := ids[1]
	data, err = s.client.kceconn.DescribeClusterInstance(&map[string]interface{}{
		"ClusterId":        clusterId, //d.Get("cluster_id"),
		"Filter.1.Name":    "instance-id",
		"Filter.1.Value.1": instanceId, //d.Get("instance_id"),
	})
	if err != nil {
		return
	}
	instanceSetSrc := (*data)["InstanceSet"]
	if instanceSetSrc == nil {
		d.SetId("")
		return
	}
	instanceSet := instanceSetSrc.([]interface{})
	instanceInfo := instanceSet[0].(map[string]interface{})
	logger.Debug("ReadAndSetWorker", "ReadAndSetWorker", instanceInfo)

	d.Set("cluster_id", clusterId)
	d.Set("instance_id", instanceId)

	imageId, _ := getSdkValue("KecInstancePara.ImageId", instanceInfo)
	d.Set("image_id", imageId)

	if advancedSetting, ok := instanceInfo["AdvancedSetting"].(map[string]interface{}); ok {
		for k, v := range advancedSetting {
			updateResourceData(d, k, v)
		}
	}

	// todo：由于封锁接口未开放，暂不更新这个字段
	// 创建后，advanceSetting里的schedulable就不更新了，驱逐状态由节点上的UnSchedulable字段返回
	// 由于tf只暴露了schedulable，所以在这里将字段做个映射
	//unSchedulableInterface, _ := getSdkValue("UnSchedulable", instanceInfo)
	//logger.Debug("unSchedulableInterface", "unSchedulableInterface", unSchedulableInterface)
	//if unSchedulable, ok := unSchedulableInterface.(bool); ok {
	//	logger.Debug("unSchedulableInterface", "unSchedulableInterface", unSchedulable)
	//	updateResourceData(d, "schedulable", !unSchedulable)
	//}

	logger.Debug("ReadAndSetWorker", "ReadAndSetWorker", d)

	return
}

func (s *KceWorkerService) updateSchedulable(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {

	updateReq := map[string]interface{}{
		"ClusterId":    d.Get("cluster_id"),
		"InstanceId.1": d.Get("instance_id"),
		"IsCordonNode": !d.Get("schedulable").(bool),
	}
	callback = ApiCall{
		param: &updateReq,
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			//return resp, ModifyProjectInstanceNew(d.Id(), call.param, client)
			return
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			return err
		},
	}
	return callback, nil
}

func (s *KceWorkerService) UpdateWorker(d *schema.ResourceData, r *schema.Resource) (err error) {
	// todo 暂时没有可以update的字段
	return
	logger.Debug("UpdateWorker", "", d.HasChange("schedulable"),
		d.HasChange("label"))

	var apiCalls []ApiCall
	if d.HasChange("schedulable") {
		var apiCall ApiCall
		apiCall, err = s.updateSchedulable(d, r)
		apiCalls = append(apiCalls, apiCall)
	}
	return ksyunApiCallNew(apiCalls, d, s.client, true)
}
