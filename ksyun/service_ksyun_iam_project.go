package ksyun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type IamProjectService struct {
	client *KsyunClient
}

func (s *IamProjectService) CreateIAMProjectCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateProject",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.iamconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateProject(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				id, err = getSdkValue("Result", *resp)
				if err != nil {
					return err
				}
				d.SetId(Float64ToString(id.(float64)))
			}
			return err
		},
	}
	return callback, err
}

func (s *IamProjectService) CreateIamProject(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createProjectCall, err := s.CreateIamProjectCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createProjectCall)
	return apiProcess.Run()
}

func (s *IamProjectService) CreateIamProjectCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateIAMProjectCommonCall(req, true)
}

func (s *IamProjectService) ReadAndSetIamProject(d *schema.ResourceData, r *schema.Resource) (err error) {

	params := map[string]interface{}{}
	params["ProjectName"] = d.Get("project_name")

	var data []interface{}
	data, err = s.ReadProject(params)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *IamProjectService) ReadProject(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.iamconn
	action := "GetAccountAllProjectList"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.GetAccountAllProjectList(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.GetAccountAllProjectList(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("ListProjectResult.ProjectList", *resp)
	if err != nil {
		return data, err
	}
	projects := results.([]interface{})
	projectInfo := make(map[string]interface{})
	for _, item := range projects {
		value := item.(map[string]interface{})
		if value["ProjectName"] == condition["ProjectName"] {
			projectInfo = value
			break
		}
	}
	data = append(data, projectInfo)
	return data, err
}

func (s *IamProjectService) ReadProjects(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.iamconn
	action := "GetAccountAllProjectList"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.GetAccountAllProjectList(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.GetAccountAllProjectList(&condition)
		if err != nil {
			return data, err
		}
	}
	results, err = getSdkValue("ListProjectResult.ProjectList", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *IamProjectService) ReadAndSetIamProjects(d *schema.ResourceData, r *schema.Resource) (err error) {
	req, err := mergeDataSourcesReq(d, r, nil)
	logger.Debug(logger.ReqFormat, "GetAccountAllProjectList", req)
	if err != nil {
		return err
	}
	data, err := s.ReadProjects(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "projects",
		extra:       map[string]SdkResponseMapping{},
	})
}
