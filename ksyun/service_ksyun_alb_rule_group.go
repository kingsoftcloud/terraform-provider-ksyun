package ksyun

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type AlbRuleGroup struct {
	client *KsyunClient
}

var (
	albRuleGroupSessionNecessary = []string{"cookie_type", "session_persistence_period"}
	albRuleGroupHealthNecessary  = []string{"interval", "timeout", "healthy_threshold", "unhealthy_threshold", "url_path", "host_name"}
)

const (
	albRuleTypeForwardGroup  = "ForwardGroup"
	albRuleTypeRedirect      = "Redirect"
	albRuleTypeFixedResponse = "FixedResponse"
	albRuleTypeRewrite       = "Rewrite"
)

var (
	fixedResponseConfigResourceElem = func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"content": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The content of response.",
				},
				"http_code": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The response http code. Valid Values: 2xx|4xx|5xx. e.g. 503.",
				},
				"content_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The type of content. Valid Values: `text/plain`|`text/css`|`text/html`|`application/javascript`|`application/json`.",
				},
			},
		}
	}
)

var albRuleTypeMappingFields = map[string]string{
	"domain":   "alb_rule_value",
	"url":      "alb_rule_value",
	"header":   "header_value",
	"cookie":   "cookie_value",
	"query":    "query_value",
	"method":   "method_value",
	"sourceIp": "source_ip_value",
}

func (s *AlbRuleGroup) createRuleGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		// "alb_rule_set": {mappings: map[string]string{
		//	"alb_rule_type":  "AlbRuleType",
		//	"alb_rule_value": "AlbRuleValue",
		// }, Type: TransformListN},
		"fixed_response_config": {
			Ignore: false,
		},
		"rewrite_config": {
			Ignore: false,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	var albRuleSet []map[string]interface{}
	for _, item := range req["AlbRuleSet"].([]interface{}) {
		_item := item.(map[string]interface{})
		// _m := make(map[string]interface{}, len(_item))
		// if err := recursiveMapToTransformListN(_item, SdkReqTransform{}, &_m, ""); err != nil {
		// 	return callback, fmt.Errorf("error on recursiveMapToTransformListN, %s", err)
		// }
		_m := helper.ConvertMapKey2Title(_item, true)
		albRuleSet = append(albRuleSet, _m)
	}
	req["AlbRuleSet"] = albRuleSet

	if _, ok := req["HostName"]; !ok {
		req["HostName"] = d.Get("host_name")
	}

	if _, ok := d.GetOk("fixed_response_config"); ok {
		if mm, mOk := helper.GetSchemaListHeadMap(d, "fixed_response_config"); mOk {
			req["FixedResponseConfig"] = helper.ConvertMapKey2Title(mm, true)
		}
	}

	if _, ok := d.GetOk("rewrite_config"); ok {
		if mm, mOk := helper.GetSchemaListHeadMap(d, "rewrite_config"); mOk {
			req["RewriteConfig"] = helper.ConvertMapKey2Title(mm, true)
		}
	}

	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateAlbRuleGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateAlbRuleGroup(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("AlbRuleGroup.AlbRuleGroupId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("alb_rule_group_id", d.Id())
		},
	}
	return
}

func (s *AlbRuleGroup) CreateAlbRuleGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.createRuleGroupCall(d, r)
	if err != nil {
		return err
	}
	err = ksyunApiCallNew([]ApiCall{
		call,
	}, d, s.client, false)
	return
}

func (s *AlbRuleGroup) readRuleGroups(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.slbconn
		action := "DescribeAlbRuleGroups"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeAlbRuleGroups(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeAlbRuleGroups(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("AlbRuleGroupSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *AlbRuleGroup) readRuleGroup(d *schema.ResourceData, ruleGroupId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if ruleGroupId == "" {
		ruleGroupId = d.Id()
	}
	req := map[string]interface{}{
		"AlbRuleGroupId.1": ruleGroupId,
	}
	results, err = s.readRuleGroups(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("ALb rule group %s not exist ", ruleGroupId)
	}
	return
}

func (s *AlbRuleGroup) ReadAndSetRuleGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "AlbRuleGroupId",
			Type:    TransformWithN,
		},
		"alb_listener_id": {
			mapping: "alblistener-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readRuleGroups(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "AlbRuleGroupName",
		idFiled:     "AlbRuleGroupId",
		targetField: "alb_rule_groups",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *AlbRuleGroup) ReadAndSetRuleGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.readRuleGroup(d, "")
	if err != nil {
		return err
	}

	// get defines set
	albRuleSet := d.Get("alb_rule_set").([]interface{})

	// get response set
	retRuleSet, ok := data["AlbRuleSet"].([]interface{})
	if !ok {
		return errors.New("parse alb rule group response error")
	}

	var storeRuleSet []map[string]interface{}
	for _, albRuleIf := range albRuleSet {
		albRule, ok := albRuleIf.(map[string]interface{})
		if !ok {
			return errors.New("parse alb rule group response error")
		}
		for t, v := range albRule {
			if t != "alb_rule_type" {
				continue
			}
			for _, retItem := range retRuleSet {
				retM := retItem.(map[string]interface{})
				// 1. get rule type
				// 2. get type field
				if retM["AlbRuleType"] == v {
					hashFunc := helper.HashFuncWithKeys("alb_rule_type", albRuleTypeMappingFields[v.(string)])
					convertMap := helper.ConvertMapKey2Underline(retM)
					convertAlbMap := helper.ConvertMapKey2Underline(albRule)
					if hashFunc(convertMap) == hashFunc(convertAlbMap) {
						storeRuleSet = append(storeRuleSet, convertMap)
						goto breakDouble
					}
				}

			}
		}
	breakDouble:
	}

	if err := d.Set("alb_rule_set", storeRuleSet); err != nil {
		return err
	}
	// alb rule set by manual set
	delete(data, "AlbRuleSet")

	if v, ok := data["CookieType"]; ok {
		if v != "RewriteCookie" {
			delete(data, "CookieName")
		}
	}

	if v, ok := data[fixedResponseConfig]; ok {
		vm := v.(map[string]interface{})
		if len(vm) > 0 {
			responseConfig := helper.ConvertMapKey2Underline(vm)
			data[fixedResponseConfig] = []interface{}{responseConfig}
		}
	}

	extra := map[string]SdkResponseMapping{}
	SdkResponseAutoResourceData(d, r, data, extra)
	return
}

func (s *AlbRuleGroup) removeAlbRuleGroupCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"AlbRuleGroupId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteAlbRuleGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteAlbRuleGroup(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.readRuleGroup(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading rule group when delete %q, %s", d.Id(), callErr))
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
	return
}
func (s *AlbRuleGroup) RemoveRuleGroup(d *schema.ResourceData) (err error) {
	call, err := s.removeAlbRuleGroupCall(d)
	if err != nil {
		return err
	}
	err = ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
	return
}

func (s *AlbRuleGroup) modifyRuleGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{}

	if _, n := d.GetChange("listener_sync"); n == "off" {
		transform["session_state"] = SdkReqTransform{
			forceUpdateParam: true,
		}
		transform["health_check_state"] = SdkReqTransform{
			forceUpdateParam: true,
		}
		transform["method"] = SdkReqTransform{
			forceUpdateParam: true,
		}
		if d.Get("session_state") == "start" {
			for _, s := range albRuleGroupSessionNecessary {
				transform[s] = SdkReqTransform{
					forceUpdateParam: true,
				}
			}
		}
		if d.Get("health_check_state") == "start" {
			for _, s := range albRuleGroupHealthNecessary {
				transform[s] = SdkReqTransform{
					forceUpdateParam: true,
				}
			}
		}

		if d.HasChange("health_protocol"); d.Get("health_protocol") == "HTTP" {
			transform["url_path"] = SdkReqTransform{
				forceUpdateParam: true,
			}

			transform["http_method"] = SdkReqTransform{
				forceUpdateParam: true,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					methodValue, ok := data.GetOk("http_method")
					if !ok {
						return "HEAD", true
					}
					return methodValue, true
				},
			}
		}
	}

	switch d.Get("type") {
	case albRuleTypeForwardGroup:
		transform["backend_server_group_id"] = SdkReqTransform{
			forceUpdateParam: true,
		}

		// ignore others' id by manual
		transform["redirect_alb_listener_id"] = SdkReqTransform{
			Ignore: true,
		}
	case albRuleTypeRedirect:
		transform["redirect_alb_listener_id"] = SdkReqTransform{
			forceUpdateParam: true,
		}

		transform["backend_server_group_id"] = SdkReqTransform{
			Ignore: true,
		}
	case albRuleTypeFixedResponse:
		transform["backend_server_group_id"] = SdkReqTransform{
			Ignore: true,
		}
		transform["redirect_alb_listener_id"] = SdkReqTransform{
			Ignore: true,
		}
	case albRuleTypeRewrite:
		transform["backend_server_group_id"] = SdkReqTransform{
			forceUpdateParam: true,
		}
		transform["rewrite_config"] = SdkReqTransform{
			forceUpdateParam: true,
		}
	default:
	}

	transform["alb_rule_set"] = SdkReqTransform{
		forceUpdateParam: true,
	}
	transform["fixed_response_config"] = SdkReqTransform{
		Ignore: true,
	}

	transform["rewrite_config"] = SdkReqTransform{
		Ignore: true,
	}

	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	logger.Debug(logger.ReqFormat, "modifyRuleGroupCall", req)
	if err != nil {
		return callback, err
	}

	if d.HasChange("fixed_response_config") {
		// if the number of fixed_response_config more than 0, always set it
		if mm, mOk := helper.GetSchemaListHeadMap(d, "fixed_response_config"); mOk {
			req["FixedResponseConfig"] = helper.ConvertMapKey2Title(mm, true)
		}
	}

	if d.HasChange("rewrite_config") {
		if mm, mOk := helper.GetSchemaListHeadMap(d, "rewrite_config"); mOk {
			req["RewriteConfig"] = helper.ConvertMapKey2Title(mm, true)
		}
	}

	if albRuleSetParams, ok := req["AlbRuleSet"]; ok {
		var albRuleSet []map[string]interface{}
		for _, item := range albRuleSetParams.([]interface{}) {
			_item := item.(map[string]interface{})
			// _m := make(map[string]interface{}, len(_item))
			// if err := recursiveMapToTransformListN(_item, SdkReqTransform{}, &_m, ""); err != nil {
			// 	return callback, fmt.Errorf("error on recursiveMapToTransformListN, %s", err)
			// }
			_m := helper.ConvertMapKey2Title(_item, true)
			albRuleSet = append(albRuleSet, _m)
		}

		if albRuleSet != nil && len(albRuleSet) > 0 {
			req["AlbRuleSet"] = albRuleSet
		}
	}

	if len(req) > 0 {
		req["AlbRuleGroupId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyAlbRuleGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyAlbRuleGroup(call.param)
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
func (s *AlbRuleGroup) ModifyRuleGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	callbacks := []ApiCall{}
	var call ApiCall
	call, err = s.modifyRuleGroupCall(d, r)
	if err != nil {
		return
	}
	callbacks = append(callbacks, call)
	err = ksyunApiCallNew(callbacks, d, s.client, true)
	return
}
