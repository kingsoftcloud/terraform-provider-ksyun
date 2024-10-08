package ksyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

type IamPolicyService struct {
	client *KsyunClient
}

func (s *IamPolicyService) CreateIAMPolicyCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreatePolicy",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreatePolicy(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				id, err = getSdkValue("CreatePolicyResult.Policy.PolicyName", *resp)
				if err != nil {
					return err
				}
				d.SetId(id.(string))
				krn, err := getSdkValue("CreatePolicyResult.Policy.Krn", *resp)
				if err != nil {
					return err
				}
				err = d.Set("policy_krn", krn)
				if err != nil {
					return err
				}
			}
			return err
		},
	}
	return callback, err
}

func (s *IamPolicyService) CreateIamPolicy(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createPolicyCall, err := s.CreateIamPolicyCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createPolicyCall)
	return apiProcess.Run()
}

func (s *IamPolicyService) CreateIamPolicyCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateIAMPolicyCommonCall(req, true)
}

func (s *IamPolicyService) DeleteIamPolicy(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteIamPolicyCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)

	return apiProcess.Run()
}

func (s *IamPolicyService) DeleteIamPolicyCall(d *schema.ResourceData) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}
	params["PolicyKrn"] = d.Get("policy_krn").(string)

	callback = ApiCall{
		param:  &params,
		action: "DeletePolicy",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeletePolicy(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				if notFoundError(baseErr) {
					return nil
				}

				// it cannot be deleted if this is still using
				if isExpectError(baseErr, []string{
					"PolicyDeleteConflict",
				}) {
					return resource.NonRetryableError(baseErr)
				}
				return resource.RetryableError(baseErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
}

func (s *IamPolicyService) ReadAndSetIamPolicy(d *schema.ResourceData, r *schema.Resource) (err error) {

	params := map[string]interface{}{}
	params["PolicyKrn"] = d.Get("policy_krn").(string)

	var data []interface{}
	data, err = s.ReadPolicy(params)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *IamPolicyService) ReadPolicy(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.iamconn
	action := "ListPolicyVersions"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.ListPolicyVersions(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.ListPolicyVersions(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("ListPolicyVersionsResult.Versions.member", *resp)
	if err != nil {
		return data, err
	}
	res := results.([]interface{})
	for _, item := range res {
		itemMap := item.(map[string]interface{})
		data = append(data, itemMap)
	}
	return data, err
}
