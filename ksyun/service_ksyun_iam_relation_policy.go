package ksyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
	"time"
)

type IamRelationPolicyService struct {
	client *KsyunClient
}

type ListAttachedUserPoliciesResult struct {
	ListAttachedUserPoliciesResult struct {
		AttachedPolicies struct {
			Member []ListAttachedUserOrRolePoliciesMembersResult `json:"member,omitempty"`
		} `json:"AttachedPolicies,omitempty"`
		IsTruncated bool `json:"IsTruncated,omitempty"`
		Total       int  `json:"Total,omitempty"`
	} `json:"ListAttachedUserPoliciesResult,omitempty"`
	RequestId string `json:"RequestId,omitempty"`
}

type ListAttachedRolePoliciesResult struct {
	ListAttachedRolePoliciesResult struct {
		AttachedPolicies struct {
			Member []ListAttachedUserOrRolePoliciesMembersResult `json:"member,omitempty"`
		} `json:"AttachedPolicies,omitempty"`
		IsTruncated bool   `json:"IsTruncated,omitempty"`
		Marker      string `json:"Marker,omitempty"`
	} `json:"ListAttachedRolePoliciesResult,omitempty"`
	RequestId string `json:"RequestId,omitempty"`
}

type ListAttachedUserOrRolePoliciesMembersResult struct {
	PolicyKrn     string `json:"PolicyKrn,omitempty"`
	PolicyName    string `json:"PolicyName,omitempty"`
	CreateTime    string `json:"CreateTime,omitempty"`
	Description   string `json:"Description,omitempty"`
	DescriptionEn string `json:"Description_en,omitempty"`
	Type          int    `json:"Type,omitempty"`
}

func (s *IamRelationPolicyService) CreateIAMRelationPolicyCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	if req["RelationType"] == 1 {
		sendParams := map[string]interface{}{}
		sendParams["UserName"] = req["Name"]

		conn := s.client.iamconn
		action := "GetUser"
		logger.Debug(logger.ReqFormat, action, sendParams)
		resp, err := conn.GetUser(&sendParams)
		if err != nil {
			return callback, err
		}
		user := (*resp)["GetUserResult"].(map[string]interface{})["User"].(map[string]interface{})
		re := regexp.MustCompile(`::(\d+):`)
		match := re.FindStringSubmatch(user["Krn"].(string))
		accountId := match[1]
		if req["PolicyType"] == "system" {
			sendParams["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::ksc:policy/%s", req["PolicyName"])
		} else {
			sendParams["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::%s:policy/%s", accountId, req["PolicyName"])
		}

		callback = ApiCall{
			param:  &sendParams,
			action: "AttachUserPolicy",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.iamconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.AttachUserPolicy(call.param)
				if err == nil {
					d.SetId(fmt.Sprintf("%s", accountId))
				}
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = d.Set("relation_type", 1)
				if err != nil {
					return err
				}
				err = d.Set("name", sendParams["UserName"])
				if err != nil {
					return err
				}
				err = d.Set("policy_name", req["PolicyName"])
				if err != nil {
					return err
				}
				err = d.Set("policy_type", req["PolicyType"])
				if err != nil {
					return err
				}
				return err
			},
		}
		return callback, err
	} else {
		sendParams := map[string]interface{}{}
		sendParams["RoleName"] = req["Name"]

		conn := s.client.iamconn
		action := "GetRole"
		logger.Debug(logger.ReqFormat, action, sendParams)
		resp, err := conn.GetRole(&sendParams)
		if err != nil {
			return callback, err
		}
		user := (*resp)["GetRoleResult"].(map[string]interface{})["Role"].(map[string]interface{})
		re := regexp.MustCompile(`::(\d+):`)
		match := re.FindStringSubmatch(user["Krn"].(string))
		accountId := match[1]
		if req["PolicyType"] == "system" {
			sendParams["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::ksc:policy/%s", req["PolicyName"])
		} else {
			sendParams["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::%s:policy/%s", accountId, req["PolicyName"])
		}

		callback = ApiCall{
			param:  &sendParams,
			action: "AttachRolePolicy",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.iamconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.AttachRolePolicy(call.param)
				if err == nil {
					d.SetId(fmt.Sprintf("%s", accountId))
				}
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = d.Set("relation_type", 2)
				if err != nil {
					return err
				}
				err = d.Set("name", sendParams["RoleName"])
				if err != nil {
					return err
				}
				err = d.Set("policy_name", req["PolicyName"])
				if err != nil {
					return err
				}
				return err
			},
		}
		return callback, err
	}
}

func (s *IamRelationPolicyService) CreateIamRelationPolicy(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createPolicyCall, err := s.CreateIamRelationPolicyCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createPolicyCall)
	return apiProcess.Run()
}

func (s *IamRelationPolicyService) CreateIamRelationPolicyCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateIAMRelationPolicyCommonCall(req, false)
}

func (s *IamRelationPolicyService) DeleteIamRelationPolicy(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteIamRelationPolicyCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)

	return apiProcess.Run()
}

func (s *IamRelationPolicyService) DeleteIamRelationPolicyCall(d *schema.ResourceData) (callback ApiCall, err error) {
	relationType := d.Get("relation_type").(int)
	if relationType == 1 {
		// 构成参数
		params := map[string]interface{}{}
		params["UserName"] = d.Get("name").(string)

		policyType := d.Get("policy_type").(string)
		if policyType == "system" {
			params["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::ksc:policy/%s", d.Get("policy_name").(string))
		} else {
			params["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::%s:policy/%s", d.Id(), d.Get("policy_name").(string))
		}
		callback = ApiCall{
			param:  &params,
			action: "DetachUserPolicy",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.iamconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.DetachUserPolicy(call.param)
				return resp, err
			},
			callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					if notFoundError(baseErr) {
						return nil
					}

					// it cannot be deleted if this is still using
					if isExpectError(baseErr, []string{
						"PolicyNoSuchEntity",
						"UserNoSuchEntity",
						"UserPolicyNoSuchEntity",
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
	} else {
		// 构成参数
		params := map[string]interface{}{}
		params["RoleName"] = d.Get("name").(string)
		policyType := d.Get("policy_type").(string)
		if policyType == "system" {
			params["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::ksc:policy/%s", d.Get("policy_name").(string))
		} else {
			params["PolicyKrn"] = fmt.Sprintf("krn:ksc:iam::%s:policy/%s", d.Id(), d.Get("policy_name").(string))
		}

		callback = ApiCall{
			param:  &params,
			action: "DetachRolePolicy",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.iamconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.DetachRolePolicy(call.param)
				return resp, err
			},
			callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					if notFoundError(baseErr) {
						return nil
					}

					// it cannot be deleted if this is still using
					if isExpectError(baseErr, []string{
						"PolicyNoSuchEntity",
						"RoleNoSuchEntity",
						"RolePolicyNoSuchEntity",
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
	}
	return
}

func (s *IamRelationPolicyService) ReadAndSetIamRelationPolicy(d *schema.ResourceData, r *schema.Resource) (err error) {
	relationType := d.Get("relation_type").(int)
	name := d.Get("name").(string)
	policyName := d.Get("policy_name").(string)

	var data []interface{}
	data, err = s.ReadRelationPolicy(relationType, name, policyName)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *IamRelationPolicyService) ReadRelationPolicy(relationType int, name string, policyName string) (data []interface{}, err error) {
	var resp *map[string]interface{}

	var policyMemberResult ListAttachedUserOrRolePoliciesMembersResult
	if relationType == 1 {
		var condition = map[string]interface{}{}
		page, MaxItems := 1, 1
		for {
			condition["UserName"] = name
			condition["Page"] = page
			condition["MaxItems"] = MaxItems
			conn := s.client.iamconn
			action := "ListAttachedUserPolicies"
			logger.Debug(logger.ReqFormat, action, condition)
			resp, err = conn.ListAttachedUserPolicies(&condition)
			if err != nil {
				return data, err
			}
			respJson, err := json.Marshal(resp)
			if err != nil {
				return data, err
			}
			var result ListAttachedUserPoliciesResult
			err = json.Unmarshal(respJson, &result)
			if err != nil {
				return data, err
			}
			if len(result.ListAttachedUserPoliciesResult.AttachedPolicies.Member) <= 0 {
				break
			}
			for _, item := range result.ListAttachedUserPoliciesResult.AttachedPolicies.Member {
				if item.PolicyName == policyName {
					policyMemberResult = item
					break
				}
			}
			if !result.ListAttachedUserPoliciesResult.IsTruncated {
				break
			}
			page++
		}
	} else {
		var condition = map[string]interface{}{}
		marker, MaxItems := "", 1
		for {
			condition["RoleName"] = name
			condition["Marker"] = marker
			condition["MaxItems"] = MaxItems
			conn := s.client.iamconn
			action := "ListAttachedRolePolicies"
			logger.Debug(logger.ReqFormat, action, condition)
			resp, err = conn.ListAttachedRolePolicies(&condition)
			if err != nil {
				return data, err
			}
			respJson, err := json.Marshal(resp)
			if err != nil {
				return data, err
			}
			var result ListAttachedRolePoliciesResult
			err = json.Unmarshal(respJson, &result)
			if err != nil {
				return data, err
			}
			if len(result.ListAttachedRolePoliciesResult.AttachedPolicies.Member) <= 0 {
				break
			}
			for _, item := range result.ListAttachedRolePoliciesResult.AttachedPolicies.Member {
				if item.PolicyName == policyName {
					policyMemberResult = item
					break
				}
			}
			if !result.ListAttachedRolePoliciesResult.IsTruncated {
				break
			}
			marker = result.ListAttachedRolePoliciesResult.Marker
		}
	}

	if policyMemberResult.PolicyName != "" {
		data = append(data, map[string]interface{}{
			"PolicyKrn":     policyMemberResult.PolicyKrn,
			"PolicyName":    policyMemberResult.PolicyName,
			"CreateTime":    policyMemberResult.CreateTime,
			"Description":   policyMemberResult.Description,
			"DescriptionEn": policyMemberResult.DescriptionEn,
			"Type":          policyMemberResult.Type,
		})
	}

	return data, err
}
