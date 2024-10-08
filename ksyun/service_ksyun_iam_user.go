package ksyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type IamUserService struct {
	client *KsyunClient
}

func (s *IamUserService) CreateIAMUserCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateUser",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateUser(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				id, err = getSdkValue("CreateUserResult.User.UserName", *resp)
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

func (s *IamUserService) CreateIamUser(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createUserCall, err := s.CreateIamUserCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createUserCall)
	return apiProcess.Run()
}

func (s *IamUserService) CreateIamUserCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateIAMUserCommonCall(req, true)
}

func (s *IamUserService) ReadAndSetIamUser(d *schema.ResourceData, r *schema.Resource) (err error) {

	params := map[string]interface{}{}
	params["UserName"] = d.Get("user_name")

	var data []interface{}
	data, err = s.ReadUser(params)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *IamUserService) ReadUser(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.iamconn
	action := "GetUser"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.GetUser(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.GetUser(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("GetUserResult.User", *resp)
	if err != nil {
		return data, err
	}
	data = append(data, results)
	return data, err
}

func (s *IamUserService) ReadUsers(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	condition["MaxItems"] = 1000
	conn := s.client.iamconn
	action := "ListUsers"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.ListUsers(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.ListUsers(&condition)
		if err != nil {
			return data, err
		}
	}
	results, err = getSdkValue("ListUserResult.Users.member", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *IamUserService) DeleteIamUser(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteIamUserCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)

	return apiProcess.Run()
}

func (s *IamUserService) DeleteIamUserCall(d *schema.ResourceData) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}
	params["UserName"] = d.Get("user_name")

	callback = ApiCall{
		param:  &params,
		action: "DeleteUser",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteUser(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				if notFoundError(baseErr) {
					return nil
				}

				// it cannot be deleted if this is still using
				if isExpectError(baseErr, []string{
					"UserNoSuchEntity",
					"UserAkDeleteConflict",
					"UserPolicyDeleteConflict",
					"UserMfaDeleteConflict",
					"UserGroupRelationDeleteConflict",
					"ProjectMemberRelationDeleteConflict",
					"ProjectMemberRelationDeleteConflict",
					"userAlarmDeleteConflict",
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

func (s *IamUserService) ReadAndSetIamUsers(d *schema.ResourceData, r *schema.Resource) (err error) {
	req, err := mergeDataSourcesReq(d, r, nil)
	logger.Debug(logger.ReqFormat, "ListUsers", req)
	if err != nil {
		return err
	}
	data, err := s.ReadUsers(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "users",
		extra:       map[string]SdkResponseMapping{},
	})
}
