package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

type SksService struct {
	client *KsyunClient
}

func (s *SksService) ReadKeys(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.sksconn
		action := "DescribeKeys"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeKeys(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeKeys(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("KeySet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *SksService) ReadKey(d *schema.ResourceData, keyId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if keyId == "" {
		keyId = d.Id()
	}
	req := map[string]interface{}{
		"KeyId.1": keyId,
	}
	results, err = s.ReadKeys(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Key %s not exist ", keyId)
	}
	return data, err
}

func (s *SksService) ReadAndSetKey(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadKey(d, "")
	if err != nil {
		if notFoundError(err) {
			//资源已被外部删除，从 state 中移除
			d.SetId("")
			return nil
		}
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return nil

}

func (s *SksService) ReadAndSetKeys(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "KeyId",
			Type:    TransformWithN,
		},
		"key_name": {
			mapping: "key-name",
			Type:    TransformWithFilter,
		},
		"key_names": {
			mapping: "key-name",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadKeys(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "KeyId",
		targetField: "keys",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SksService) CreateKeyCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	if _, ok := d.GetOk("public_key"); ok {
		return callback, err
	}
	transform := map[string]SdkReqTransform{
		"key_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateKey",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.sksconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateKey(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Key.KeyId", *resp)
			if err != nil {
				return err
			}
			pk, err := getSdkValue("PrivateKey", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("private_key", pk)
		},
	}
	return callback, err
}

func (s *SksService) ImportKeyCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	if _, ok := d.GetOk("public_key"); !ok {
		return callback, err
	}
	transform := map[string]SdkReqTransform{
		"key_name":   {},
		"public_key": {},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "ImportKey",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.sksconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ImportKey(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Key.KeyId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SksService) CreateKey(d *schema.ResourceData, r *schema.Resource) (err error) {
	createCall, err := s.CreateKeyCall(d, r)
	if err != nil {
		return err
	}
	importCall, err := s.ImportKeyCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{createCall, importCall}, d, s.client, true)
}

func (s *SksService) ModifyKeyCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"key_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["KeyId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyKey",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.sksconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyKey(call.param)
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

func (s *SksService) ModifyKey(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyKeyCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SksService) RemoveKeyCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"KeyId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteKey",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.sksconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteKey(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadKey(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading key when delete %q, %s", d.Id(), callErr))
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

func (s *SksService) RemoveKey(d *schema.ResourceData) (err error) {
	call, err := s.RemoveKeyCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
