package ksyun

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// PerKnadService: 按次计费原生高防（PerPayKnad）
type PerKnadService struct {
	client *KsyunClient
}

func (s *PerKnadService) doKnadCustomAction(action string, input *map[string]interface{}) (resp *map[string]interface{}, err error) {
	conn := s.client.knadconn
	op := &request.Operation{
		Name:       action,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}
	if input == nil {
		input = &map[string]interface{}{}
	}
	output := &map[string]interface{}{}
	req := conn.NewRequest(op, input, output)
	err = req.Send()
	return output, err
}

func (s *PerKnadService) ReadPerKnads(condition map[string]interface{}) (data []interface{}, err error) {
	// 目前 SDK 中无 DescribePerPayKnad，按次实例仍通过 DescribeKnad 查询
	knadService := KnadService{s.client}
	return knadService.ReadKnads(condition)
}

func (s *PerKnadService) ReadAndSetPerKnads(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "KnadId",
			Type:    TransformWithN,
		},
		"project_id": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadPerKnads(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "perknads",
		idFiled:     "KnadId",
	})
}

// resource_查 有重试    跟data_resource里的查不一样
func (s *PerKnadService) ReadAndSetPerKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadPerKnad(d, "")
		if callErr != nil {
			if notFoundError(callErr) {
				d.SetId("")
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading knad %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, map[string]SdkResponseMapping{
				"ProjectId": {
					Field: "project_id",
					FieldRespFunc: func(element interface{}) interface{} {
						projectIdStr := strconv.FormatFloat(element.(float64), 'f', -1, 64)
						return projectIdStr
					},
				},
			})
			return nil
		}
	})
}

func (s *PerKnadService) ReadPerKnad(d *schema.ResourceData, knadId string) (data map[string]interface{}, err error) {
	knadService := KnadService{s.client}
	return knadService.ReadKnad(d, knadId)
}

func (s *PerKnadService) CreatePerKnadCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreatePerPayKnad",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			// SDK 中暂未生成 CreatePerPayKnad 方法，这里通过自定义 Operation 调用
			return s.doKnadCustomAction(call.action, call.param)
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Kid", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return nil
		},
	}
	return callback, nil
}

func (s *PerKnadService) CreatePerKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreatePerKnadCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *PerKnadService) ModifyPerKnadCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {Ignore: true},
		"knad_name":  {Ignore: true},
		"ip_count":   {forceUpdateParam: true},
		"max_band":   {forceUpdateParam: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) == 0 {
		return callback, nil
	}

	req["KnadId"] = d.Id()
	callback = ApiCall{
		param:  &req,
		action: "ModifyPerPayKnad",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			// SDK 中暂未生成 ModifyPerPayKnad 方法，这里通过自定义 Operation 调用
			return s.doKnadCustomAction(call.action, call.param)
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return nil
		},
	}
	return callback, nil
}

func (s *PerKnadService) ModifyPerKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyPerKnadCall(d, r)
	if err != nil {
		return err
	}
	if call.executeCall == nil {
		return nil
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *PerKnadService) RemovePerKnadCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"KnadId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteKnad",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.knadconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteKnad(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			// 复用 knad 删除重试逻辑：删除后需要确认资源不存在
			return resource.Retry(1*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadPerKnad(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					}
					return resource.NonRetryableError(fmt.Errorf("error on reading perknad when delete %q, %s", d.Id(), callErr))
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
			return nil
		},
	}
	return callback, nil
}

func (s *PerKnadService) RemovePerKnad(d *schema.ResourceData) (err error) {
	call, err := s.RemovePerKnadCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
