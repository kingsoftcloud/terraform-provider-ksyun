package ksyun

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type AlbRuleGroup struct {
	client *KsyunClient
}

func (s *AlbRuleGroup) createRuleGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		// "alb_rule_set": {mappings: map[string]string{
		//	"alb_rule_type":  "AlbRuleType",
		//	"alb_rule_value": "AlbRuleValue",
		// }, Type: TransformListN},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	var albRuleSet []map[string]interface{}
	for _, item := range req["AlbRuleSet"].([]interface{}) {
		albRuleSet = append(albRuleSet, map[string]interface{}{
			"AlbRuleType":  item.(map[string]interface{})["alb_rule_type"],
			"AlbRuleValue": item.(map[string]interface{})["alb_rule_value"],
		})
	}
	req["AlbRuleSet"] = albRuleSet

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
			m := make(map[string]interface{}, 2)
			for _, retItem := range retRuleSet {
				retM := retItem.(map[string]interface{})
				for rt, rv := range retM {
					if reflect.DeepEqual(rt, "AlbRuleType") && reflect.DeepEqual(v, rv) {
						m["alb_rule_value"] = retM["AlbRuleValue"]
						m["alb_rule_type"] = rv
						storeRuleSet = append(storeRuleSet, m)
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
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	logger.Debug(logger.ReqFormat, "modifyRuleGroupCall", req)
	if err != nil {
		return callback, err
	}

	if albRuleSetParams, ok := req["AlbRuleSet"]; ok {
		var albRuleSet []map[string]interface{}
		for _, item := range albRuleSetParams.([]interface{}) {
			albRuleSet = append(albRuleSet, map[string]interface{}{
				"AlbRuleType":  item.(map[string]interface{})["alb_rule_type"],
				"AlbRuleValue": item.(map[string]interface{})["alb_rule_value"],
			})
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
