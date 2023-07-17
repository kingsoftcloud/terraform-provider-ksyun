package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
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
	var (
		results []interface{}
	)
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

	//tagService := TagService{s.client}
	//tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, r, "eip", true, false)
	//if err != nil {
	//	return err
	//}
	return ksyunApiCallNew(calls, d, alb.client, true)
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
			SdkResponseAutoResourceData(d, r, data, chargeExtraForVpc(data))
			return nil
		}
	})
	//return
}
func (alb *AlbService) stateRefreshFunc(d *schema.ResourceData, albId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
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

	return ksyunApiCallNew(calls, d, alb.client, true)
}
