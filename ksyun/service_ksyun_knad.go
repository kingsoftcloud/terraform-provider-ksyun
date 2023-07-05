package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

type KnadService struct {
	client *KsyunClient
}

func (s *KnadService) ReadAndSetKnads(d *schema.ResourceData, r *schema.Resource) (err error) {
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

	data, err := s.ReadKnads(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "knads",
		idFiled:     "KnadId",
		extra: map[string]SdkResponseMapping{
			//类型转换 float转string
			"ProjectId": {
				Field: "project_id",
				FieldRespFunc: func(element interface{}) interface{} {
					projectIdStr := strconv.FormatFloat(element.(float64), 'f', -1, 64)
					return projectIdStr
				},
			},
		},
	})
}

func (s *KnadService) ReadKnads(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.knadconn
	action := "DescribeKnad"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeKnad(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeKnad(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("KnadSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

// resource_查 有重试    跟data_resource里的查不一样
func (s *KnadService) ReadAndSetKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadKnad(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
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

func (s *KnadService) ReadKnad(d *schema.ResourceData, knadId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if knadId == "" {
		knadId = d.Id()
	}
	req := map[string]interface{}{
		"KnadId.1": knadId,
	}
	results, err = s.ReadKnads(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("knad %s not exist ", knadId)
	}
	return data, err
}

func (s *KnadService) CreateKnadCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateKnad",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.knadconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateKnad(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Kid", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *KnadService) CreateKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateKnadCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// 更新入口 ModifyKnad
func (s *KnadService) ModifyKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	/*projectCall, err := s.ModifyKnadProjectCall(d, r)
	if err != nil {
		return err
	}*/
	call, err := s.ModifyKnadCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
	//return ksyunApiCallNew([]ApiCall{projectCall, call}, d, s.client, true)
}

func (s *KnadService) ModifyKnadProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil, SdkReqParameter{
		false,
	})

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

func (s *KnadService) ModifyKnadCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"service_id": {forceUpdateParam: true},
		"project_id": {Ignore: true},
		"ip_count":   {forceUpdateParam: true},
		"band":       {forceUpdateParam: true},
		"max_band":   {forceUpdateParam: true},
		"idc_band":   {forceUpdateParam: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["KnadId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyKnad",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.knadconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyKnad(call.param)
				if err != nil {
					return nil, err
				}
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				return err
			},
		}
	}
	return callback, err
}

func (s *KnadService) RemoveKnadCall(d *schema.ResourceData) (callback ApiCall, err error) {
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
			return resource.Retry(1*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadKnad(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading Knad when delete %q, %s", d.Id(), callErr))
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

// 删除knad入口
func (s *KnadService) RemoveKnad(d *schema.ResourceData) (err error) {
	call, err := s.RemoveKnadCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KnadService) ReadAndSetAssociateKnad(d *schema.ResourceData, r *schema.Resource) (err error) {

	data, err := s.ReadKnadAssociate(d, d.Get("knad_id").(string), d.Get("ip").(*schema.Set))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

// 判断subset是否是superset的子集
func is_Subset(subset []string, superset []string) bool {
	checkset := make(map[string]bool)
	for _, element := range subset {
		checkset[element] = true
	}
	for _, value := range superset {
		if checkset[value] {
			delete(checkset, value)
		}
	}
	return len(checkset) == 0
}

func (s *KnadService) ReadKnadAssociate(d *schema.ResourceData, knadId string, ip *schema.Set) (result map[string]interface{}, err error) {
	data, err := s.ReadKnadIpList(d, knadId)
	result = make(map[string]interface{})
	//var emptyMap = make(map[string]interface{})
	var associateIps []string //db里读出来的ips
	//if len(data) == 0 {
	//	return emptyMap, fmt.Errorf("instance has not band ips")
	//}

	for _, v := range data {
		associateIps = append(associateIps, v.(map[string]interface{})["Ip"].(string))
	}
	/*ipSlice := make([]string, 0) //d里获取的ip
	for _, _ip := range ip.List() {
		ipSlice = append(ipSlice, _ip.(string))
	}
	if len(ipSlice) != len(associateIps) {
		return emptyMap, fmt.Errorf("ip not associate knad")
	}*/
	//isFound := is_Subset(ipSlice, associateIps)
	//
	//if !isFound {
	//	return emptyMap, fmt.Errorf("ip not associate knad") //todo
	//}

	result["KnadId"] = knadId
	result["Ip"] = associateIps
	return result, err
}

// 已绑ip列表
func (s *KnadService) ReadKnadIpList(d *schema.ResourceData, knadId string) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	if knadId == "" {
		knadId = d.Id()
	}
	req := map[string]interface{}{
		"KnadId": knadId,
	}
	/*
		err = addProjectInfo(d, &req, s.client)
		if err != nil {
			return data, err
		}
	*/
	conn := s.client.knadconn
	resp, err = conn.IpList(&req)
	if err != nil {
		return data, err
	}
	results, err = getSdkValue("EipSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *KnadService) DisassociateKnad(d *schema.ResourceData) (err error) {
	call, err := s.DisassociateKnadCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KnadService) DisassociateKnadCall(d *schema.ResourceData) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"ip": {Type: TransformWithN},
	}

	req, err := SdkRequestAutoMapping(d, resourceKsyunKnadAssociate(), false, transform, nil)
	callback = ApiCall{
		param:  &req,
		action: "DisassociateIp",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.knadconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DisassociateIp(call.param)
			return resp, err
		},
		/**
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadKnad(d, d.Get("knad_id").(string))
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading knad associate when disassociate %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		*/
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *KnadService) AssociateKnad(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.AssociateKnadCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
func (s *KnadService) AssociateKnadCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"knad_id": {},
		"ip":      {Type: TransformWithN},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)

	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AssociateIp",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.knadconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AssociateIp(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("knad_id").(string))
			return err
		},
	}
	return callback, err
}
