package ksyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type IamGroupService struct {
	client *KsyunClient
}

func (s *IamGroupService) CreateIAMGroupCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateGroup(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				id, err = getSdkValue("CreateGroupResult.Group.GroupName", *resp)
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

func (s *IamGroupService) CreateIamGroup(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createGroupCall, err := s.CreateIamGroupCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createGroupCall)
	return apiProcess.Run()
}

func (s *IamGroupService) CreateIamGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateIAMGroupCommonCall(req, true)
}

func (s *IamGroupService) ReadAndSetIamGroup(d *schema.ResourceData, r *schema.Resource) (err error) {

	params := map[string]interface{}{}
	params["GroupName"] = d.Get("group_name")

	var data []interface{}
	data, err = s.ReadGroup(params)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *IamGroupService) ReadGroup(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.iamconn
	action := "GetGroup"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.GetGroup(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.GetGroup(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("GetGroupResult.Group", *resp)
	if err != nil {
		return data, err
	}
	data = append(data, results)
	return data, err
}

func (s *IamGroupService) ReadGroups(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	condition["MaxItems"] = 1000
	conn := s.client.iamconn
	action := "ListGroups"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.ListGroups(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.ListGroups(&condition)
		if err != nil {
			return data, err
		}
	}
	results, err = getSdkValue("ListGroupsResult.Groups.member", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *IamGroupService) DeleteIamGroup(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteIamGroupCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)

	return apiProcess.Run()
}

func (s *IamGroupService) DeleteIamGroupCall(d *schema.ResourceData) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}
	params["GroupName"] = d.Get("group_name")

	callback = ApiCall{
		param:  &params,
		action: "DeleteGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteGroup(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				if notFoundError(baseErr) {
					return nil
				}

				// it cannot be deleted if this is still using
				if isExpectError(baseErr, []string{
					"GroupNoSuchEntity",
					"GroupUserDeleteConflict",
					"GroupPolicyDeleteConflict",
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

func (s *IamGroupService) ReadAndSetIamGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
	req, err := mergeDataSourcesReq(d, r, nil)
	logger.Debug(logger.ReqFormat, "ListGroups", req)
	if err != nil {
		return err
	}
	data, err := s.ReadGroups(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "groups",
		extra:       map[string]SdkResponseMapping{},
	})
}
