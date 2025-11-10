package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

type CenService struct {
	client *KsyunClient
}

func (s *CenService) ReadCens(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.cenconn
	action := "DescribeCens"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeCens(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeCens(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("CenSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *CenService) ReadCen(d *schema.ResourceData, cenId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if cenId == "" {
		cenId = d.Id()
	}
	req := map[string]interface{}{
		"CenId.1": cenId,
	}
	results, err = s.ReadCens(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Cen %s not exist ", cenId)
	}
	return data, err
}

func (s *CenService) ReadAndSetCen(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadCen(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *CenService) ReadAndSetCens(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "CenId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadCens(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "CenName",
		idFiled:     "CenId",
		targetField: "cens",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *CenService) CreateCenCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateCen",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.cenconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateCen(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Cen.CenId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *CenService) CreateCen(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateCenCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *CenService) ModifyCenCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["CenId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyCen",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.cenconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyCen(call.param)
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

func (s *CenService) ModifyCen(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyCenCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *CenService) RemoveCenCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"CenId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteCen",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.cenconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteCen(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadCen(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading cen when delete %q, %s", d.Id(), callErr))
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

func (s *CenService) RemoveCen(d *schema.ResourceData) (err error) {
	call, err := s.RemoveCenCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
