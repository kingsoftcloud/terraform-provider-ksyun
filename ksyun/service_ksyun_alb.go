package ksyun

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type AlbService struct {
	client *KsyunClient
}

func (alb *AlbService) readAlbs(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := alb.client.slbconn
		action := "DescribeAlbs"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeAlbs(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeAlbs(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("ApplicationLoadBalancerSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})

		for _, itemInterface := range data {
			if item, ok := itemInterface.(map[string]interface{}); ok {
				if klogInfoInterface, ok := item["KlogInfo"]; ok {
					klogInfo := klogInfoInterface.(map[string]interface{})
					if v, ok := klogInfo["LogpoolName"]; ok {
						klogInfo["LogPoolName"] = v
					}
				}
			}
		}

		return data, err
	})
}

func (alb *AlbService) readAlb(d *schema.ResourceData, albId string, allProject bool) (data map[string]interface{}, err error) {
	var results []interface{}
	if albId == "" {
		albId = d.Id()
	}
	req := map[string]interface{}{
		"AlbId.1": albId,
	}
	if allProject {
		err = addProjectInfoAll(d, &req, alb.client)
		if err != nil {
			return data, err
		}
	} else {
		err = addProjectInfo(d, &req, alb.client)
		if err != nil {
			return data, err
		}
	}

	if _, ok := d.GetOk("tags"); ok {
		req["IsContainTag"] = true
	}

	results, err = alb.readAlbs(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("ALB %s not exist ", albId)
	}
	return
}

func (alb *AlbService) modifyProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		callback = ApiCall{
			param: &updateReq,
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				return resp, ModifyProjectInstanceNew(d.Id(), call.param, client)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				return err
			},
		}
	}
	return callback, err
}

func (alb *AlbService) modifyStateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}
	req["AlbId"] = d.Id()
	req["State"] = d.Get("state")
	callback = ApiCall{
		param:  &req,
		action: "SetAlbStatus",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetAlbStatus(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) modifyAlbNameCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}
	req["AlbId"] = d.Id()
	req["AlbName"] = d.Get("alb_name")
	callback = ApiCall{
		param:  &req,
		action: "SetAlbName",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetAlbName(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) modifyAccessLogCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}

	callback = ApiCall{
		param:  &req,
		action: "SetEnableAlbAccessLog",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			req["AlbId"] = d.Id()
			req["EnabledLog"] = d.Get("enabled_log")
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetEnableAlbAccessLog(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) setAlbDeleteProtectionCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}

	callback = ApiCall{
		param:  &req,
		action: "SetAlbDeleteProtection",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			req["AlbId"] = d.Id()
			req["DeleteProtection"] = d.Get("delete_protection")
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetAlbDeleteProtection(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) setAlbModificationProtectionCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}

	callback = ApiCall{
		param:  &req,
		action: "SetAlbModificationProtection",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			req["AlbId"] = d.Id()
			req["ModificationProtection"] = d.Get("modification_protection")
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetAlbModificationProtection(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) modifyAccessLogInfo(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}

	callback = ApiCall{
		param:  &req,
		action: "SetAlbAccessLog",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			req["AlbId"] = d.Id()
			req["ProjectName"] = d.Get("klog_info.0.project_name")
			req["LogPoolName"] = d.Get("klog_info.0.log_pool_name")
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetAlbAccessLog(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) ModifyAlb(d *schema.ResourceData, r *schema.Resource) (err error) {
	calls := []ApiCall{}

	if d.HasChange("project_id") {
		modifyProjectCall, err := alb.modifyProjectCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyProjectCall)
	}
	if d.HasChange("alb_name") {
		modifyNameCall, err := alb.modifyAlbNameCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyNameCall)
	}
	if d.HasChange("state") {
		modifyStateCall, err := alb.modifyStateCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyStateCall)
	}
	if d.HasChange("klog_info") {
		modifyAccessLogInfoCall, err := alb.modifyAccessLogInfo(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyAccessLogInfoCall)
	}
	if d.HasChange("enabled_log") {
		modifyEnabledLogCall, err := alb.modifyAccessLogCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyEnabledLogCall)
	}

	if d.HasChange("delete_protection") {
		modifyDeleteProtectionCall, err := alb.setAlbDeleteProtectionCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyDeleteProtectionCall)
	}

	if d.HasChange("modification_protection") {
		modifyModificationProtectionCall, err := alb.setAlbModificationProtectionCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyModificationProtectionCall)
	}

	if d.HasChanges("alb_version", "enable_hpa") {
		modifyAlbCall, err := alb.modifyAlbCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyAlbCall)
	}

	if d.HasChange("tags") {
		tagService := TagService{alb.client}
		tagsCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, r, "loadbalancer", true, false)
		if err != nil {
			return err
		}
		calls = append(calls, tagsCall)
	}

	if d.HasChange("protocol_layers") {
		modifyProtocolLayersCall, err := alb.modifyProtocolLayersCall(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, modifyProtocolLayersCall)
	}
	return ksyunApiCallNew(calls, d, alb.client, true)
}

func (alb *AlbService) modifyProtocolLayersCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{}
	req["AlbId"] = d.Id()
	if protocolLayers, ok := d.GetOk("protocol_layers"); ok {
		req["ProtocolLayers"] = protocolLayers
	} else {
		req["ProtocolLayers"] = []string{"http", "https"}
	}
	callback = ApiCall{
		param:  &req,
		action: "ModifyAlbProtocolLayers",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.SetLbProtocolLayers(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return nil
		},
	}
	return callback, err
}

func (alb *AlbService) removeAlbCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"AlbId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteAlb",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteAlb(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				data, callErr := alb.readAlb(d, "", true)
				logger.Debug(logger.RespFormat, call.action, data, callErr)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading ALB when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				data, callErr := alb.readAlb(d, "", true)
				logger.Debug(logger.RespFormat, call.action, data, callErr)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading ALB when delete %q, %s", d.Id(), callErr))
					}
				}
				status, callErr := getSdkValue("Status", data)
				logger.Debug(logger.RespFormat, "data.status", status)
				if callErr == nil {
					if statusValue, ok := status.(string); ok {
						if statusValue == "deleting" {
							return resource.RetryableError(errors.New("deleting..."))
						}
					}
				}
				return resource.RetryableError(callErr)
			})
		},
	}
	return
}

func (alb *AlbService) RemoveAlb(d *schema.ResourceData) (err error) {
	call, err := alb.removeAlbCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, alb.client, true)
}

func (alb *AlbService) ReadAndSetAlbs(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "AlbId",
			Type:    TransformWithN,
		},
		"project_id": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"state": {
			mapping: "state",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := alb.readAlbs(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "AlbId",
		targetField: "albs",
		extra: map[string]SdkResponseMapping{
			"AlbId": {
				Field:    "id",
				KeepAuto: true,
			},
		},
	})
}

func (alb *AlbService) ReadAndSetAlb(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := alb.readAlb(d, "", true)
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading ALB %q, %s", d.Id(), callErr))
			}
		} else {
			extra := chargeExtraForVpc(data)
			extra["ModifyProtection"] = SdkResponseMapping{
				Field: "modification_protection",
			}
			extra["TagSet"] = SdkResponseMapping{
				Field: "tags",
				FieldRespFunc: func(i interface{}) interface{} {
					tags := i.([]interface{})
					tagMap := make(map[string]interface{})
					for _, tag := range tags {
						_m := tag.(map[string]interface{})
						tagMap[_m["TagKey"].(string)] = _m["TagValue"].(string)
					}
					return tagMap
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
	// return
}

func (alb *AlbService) stateRefreshFunc(d *schema.ResourceData, albId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var err error
		data, err := alb.readAlb(d, albId, true)
		if err != nil {
			return nil, "", err
		}

		status, err := getSdkValue("Status", data)
		if err != nil {
			return nil, "", err
		}

		for _, v := range failStates {
			if v == status.(string) {
				return nil, "", fmt.Errorf("ALB status  error, status:%v", status)
			}
		}
		return data, status.(string), nil
	}
}

func (alb *AlbService) checkState(d *schema.ResourceData, albId string, target []string, timeout time.Duration) (state interface{}, err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    alb.stateRefreshFunc(d, albId, []string{"error"}),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	return stateConf.WaitForState()
}

func (alb *AlbService) createAlbCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateAlb",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))

			resp, err = conn.CreateAlb(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("ApplicationLoadBalancer.AlbId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))

			_, err = alb.checkState(d, "", []string{"active"}, d.Timeout(schema.TimeoutUpdate))
			return
		},
	}
	return
}

func (alb *AlbService) CreateAlb(d *schema.ResourceData, r *schema.Resource) (err error) {
	calls := []ApiCall{}
	call, err := alb.createAlbCall(d, r)
	if err != nil {
		return err
	}
	calls = append(calls, call)

	if _, ok := d.GetOk("klog_info"); ok {
		callLogInfo, err := alb.modifyAccessLogInfo(d, r)
		if err != nil {
			return err
		}
		calls = append(calls, callLogInfo)
	}
	if v, ok := d.GetOk("enabled_log"); ok {
		if v.(bool) {
			callEnableLog, err := alb.modifyAccessLogCall(d, r)
			if err != nil {
				return err
			}
			calls = append(calls, callEnableLog)
		}
	}

	if d.HasChange("tags") {
		tagsService := TagService{client: alb.client}
		tagsCall, err := tagsService.ReplaceResourcesTagsWithResourceCall(d, r, "loadbalancer", false, false)
		if err != nil {
			return err
		}
		calls = append(calls, tagsCall)
	}

	return ksyunApiCallNew(calls, d, alb.client, true)
}

func (alb *AlbService) modifyAlbCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"alb_version": {},
		"enable_hpa":  {},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: true,
	})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["AlbId"] = d.Id()

		callback = ApiCall{
			param:  &req,
			action: "ModifyAlb",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))

				resp, err = conn.ModifyAlb(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return
			},
		}
	}
	return
}

func (alb *AlbService) CreateAlbBackendServerGroup(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, alb.client, true)

	createCall, err := alb.createAlbBackendServerGroupCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(createCall)
	return apiProcess.Run()
}

func (alb *AlbService) ModifyAlbBackendServerGroup(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, alb.client, true)

	modifyCall, err := alb.modifyAlbBackendServerGroupCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(modifyCall)

	return apiProcess.Run()
}

func (alb *AlbService) ReadAndSetAlbBackendServerGroup(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := alb.readAlbBackendServerGroup(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading ALB %q, %s", d.Id(), callErr))
			}
		} else {

			extra := map[string]SdkResponseMapping{
				"Session": {
					Field: "session",
					FieldRespFunc: func(i interface{}) interface{} {
						return []interface{}{
							i,
						}
					},
				},
				"HealthCheck": {
					Field: "health_check",
					FieldRespFunc: func(i interface{}) interface{} {
						return []interface{}{
							i,
						}
					},
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (alb *AlbService) RemoveAlbBackendServerGroup(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, alb.client, true)

	removeCall, err := alb.removeAlbBackendServerGroupCall(d)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(removeCall)

	return apiProcess.Run()
}

func (alb *AlbService) createAlbBackendServerGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check": {
			Ignore: true,
		},
		"session": {
			Ignore: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}

	if d.HasChanges("health_check") {
		if healthCheck, ok := helper.GetSchemaListHeadMap(d, "health_check"); ok {
			for k, v := range healthCheck {
				switch v.(type) {
				case int, string, bool:
					if !helper.IsEmpty(v) {
						req[helper.Underline2Hump(k)] = v
					}
				}
			}
		}
	}

	if d.HasChanges("session") {
		if sessionMap, ok := helper.GetSchemaListHeadMap(d, "session"); ok {
			for k, v := range sessionMap {
				switch v.(type) {
				case int, string, bool:
					if !helper.IsEmpty(v) {
						req[helper.Underline2Hump(k)] = v
					}
				}
			}
		}
	}

	callback = ApiCall{
		param:  &req,
		action: "CreateAlbBackendServerGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))

			resp, err = conn.CreateAlbBackendServerGroup(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("BackendServerGroup.BackendServerGroupId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))

			return
		},
	}
	return callback, err
}

func (alb *AlbService) modifyAlbBackendServerGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"name":               {},
		"upstream_keepalive": {},
		"method":             {},

		"health_check": {
			Ignore: true,
		},
		"session": {
			Ignore: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}

	if d.HasChanges("health_check") {
		if healthCheck, ok := helper.GetSchemaListHeadMap(d, "health_check"); ok {
			for k, v := range healthCheck {
				switch v.(type) {
				case int, string, bool:
					if !helper.IsEmpty(v) {
						req[helper.Underline2Hump(k)] = v
					}
				}
			}
		}
	}

	if d.HasChanges("session") {
		if sessionMap, ok := helper.GetSchemaListHeadMap(d, "session"); ok {
			for k, v := range sessionMap {
				switch v.(type) {
				case int, string, bool:
					if !helper.IsEmpty(v) {
						req[helper.Underline2Hump(k)] = v
					}
				}
			}
		}
	}

	if len(req) > 0 {
		req["BackendServerGroupId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyAlbBackendServerGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))

				resp, err = conn.ModifyAlbBackendServerGroup(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return
			},
		}
	}

	return callback, err
}

func (alb *AlbService) removeAlbBackendServerGroupCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"BackendServerGroupId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteAlbBackendServerGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteAlbBackendServerGroup(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(3*time.Minute, func() *resource.RetryError {
				data, callErr := alb.readAlbBackendServerGroup(d, "")
				logger.Debug(logger.RespFormat, call.action, data, callErr)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading AlbBackendServerGroup when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
}

func (alb *AlbService) readAlbBackendServerGroups(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := alb.client.slbconn
		action := "DescribeAlbBackendServerGroups"
		logger.Debug(logger.ReqFormat, action, condition)
		resp, err = conn.DescribeAlbBackendServerGroups(&condition)
		if err != nil {
			return data, err
		}
		results, err = getSdkValue("BackendServerGroupSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})

		return data, err
	})
}

func (alb *AlbService) readAlbBackendServerGroup(d *schema.ResourceData, backendId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if backendId == "" {
		backendId = d.Id()
	}
	req := map[string]interface{}{
		"BackendServerGroupId.1": backendId,
	}

	results, err = alb.readAlbBackendServerGroups(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("AlbBackendServerGroup %s is not exist ", backendId)
	}
	return
}

func (alb *AlbService) ReadAndSetAlbBackendServerGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "BackendServerGroupId",
			Type:    TransformWithN,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"backend_server_type": {
			mapping: "backend-server-type",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := alb.readAlbBackendServerGroups(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "BackendServerGroupId",
		targetField: "alb_backend_server_groups",
	})
}

func (alb *AlbService) createAlbBackendServerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "RegisterAlbBackendServer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.RegisterAlbBackendServer(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("BackendServer.BackendServerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) modifyAlbBackendServerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["BackendServerId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyAlbBackendServer",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyAlbBackendServer(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return callback, err
}

func (alb *AlbService) ReadAlbBackendServers(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := alb.client.slbconn
		action := "DescribeAlbBackendServers"
		logger.Debug(logger.ReqFormat, action, condition)

		resp, err = conn.DescribeAlbBackendServers(&condition)
		if err != nil {
			return data, err
		}

		results, err = getSdkValue("BackendServerSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (alb *AlbService) ReadAlbBackendServer(d *schema.ResourceData, backendServerId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if backendServerId == "" {
		backendServerId = d.Id()
	}
	req := map[string]interface{}{
		"BackendServerId.1": backendServerId,
	}
	results, err = alb.ReadAlbBackendServers(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("AlbBackendServer %s not exist ", backendServerId)
	}
	return data, err
}

func (alb *AlbService) ReadAndSetAlbBackendServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := alb.ReadAlbBackendServer(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (alb *AlbService) RemoveAlbBackendServerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"BackendServerId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeregisterAlbBackendServer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeregisterAlbBackendServer(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := alb.ReadAlbBackendServer(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading alb backend server when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (alb *AlbService) RemoveAlbBackendServer(d *schema.ResourceData) (err error) {
	apiProcess := NewApiProcess(context.Background(), d, alb.client, true)

	call, err := alb.RemoveAlbBackendServerCall(d)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (alb *AlbService) CreateAlbBackendServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.Background(), d, alb.client, true)

	call, err := alb.createAlbBackendServerCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (alb *AlbService) ModifyAlbBackendServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.Background(), d, alb.client, true)

	call, err := alb.modifyAlbBackendServerCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}
