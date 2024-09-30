package ksyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
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
