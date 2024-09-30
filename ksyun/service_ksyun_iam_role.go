package ksyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type IamRoleService struct {
	client *KsyunClient
}

func (s *IamRoleService) CreateIAMRoleCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateRole",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateRole(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				id, err = getSdkValue("CreateRoleResult.Role.RoleName", *resp)
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

func (s *IamRoleService) CreateIamRole(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createRoleCall, err := s.CreateIamRoleCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createRoleCall)
	return apiProcess.Run()
}

func (s *IamRoleService) CreateIamRoleCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateIAMRoleCommonCall(req, true)
}

func (s *IamRoleService) ReadAndSetIamRole(d *schema.ResourceData, r *schema.Resource) (err error) {

	params := map[string]interface{}{}
	params["RoleName"] = d.Get("role_name")

	var data []interface{}
	data, err = s.ReadRole(params)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *IamRoleService) ReadRole(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.iamconn
	action := "GetRole"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.GetRole(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.GetRole(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("GetRoleResult.Role", *resp)
	if err != nil {
		return data, err
	}
	data = append(data, results)
	return data, err
}

func (s *IamRoleService) ReadRoles(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	condition["MaxItems"] = 1000
	conn := s.client.iamconn
	action := "ListRoles"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.ListRoles(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.ListRoles(&condition)
		if err != nil {
			return data, err
		}
	}
	results, err = getSdkValue("ListRolesResult.Roles.member", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *IamRoleService) DeleteIamRole(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteIamRoleCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)

	return apiProcess.Run()
}

func (s *IamRoleService) DeleteIamRoleCall(d *schema.ResourceData) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}
	params["RoleName"] = d.Get("role_name")

	callback = ApiCall{
		param:  &params,
		action: "DeleteRole",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteRole(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				if notFoundError(baseErr) {
					return nil
				}

				// it cannot be deleted if this is still using
				if isExpectError(baseErr, []string{
					"RoleNoSuchEntity",
					"DeleteConflict",
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

func (s *IamRoleService) ReadAndSetIamRoles(d *schema.ResourceData, r *schema.Resource) (err error) {
	req, err := mergeDataSourcesReq(d, r, nil)
	logger.Debug(logger.ReqFormat, "ListRoles", req)
	if err != nil {
		return err
	}
	data, err := s.ReadRoles(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "roles",
		extra:       map[string]SdkResponseMapping{},
	})
}
