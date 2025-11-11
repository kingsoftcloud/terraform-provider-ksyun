package ksyun

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// MonitorService
type MonitorService struct {
	client *KsyunClient
}

// CreateAlarmPolicyCommonCall creates a common API call for creating alarm policy
func (s *MonitorService) CreateAlarmPolicyCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateAlarmPolicy",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.monitorv4conn
			// Create custom operation for CreateAlarmPolicy
			op := &request.Operation{
				Name:       "CreateAlarmPolicy",
				HTTPMethod: "POST",
				HTTPPath:   "/",
			}

			params := make(map[string]interface{})
			for k, v := range *(call.param) {
				params[k] = v
			}

			output := &map[string]interface{}{}
			req := conn.NewRequest(op, &params, output)
			req.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

			logger.Debug(logger.RespFormat, call.action, params)
			err = req.Send()
			if err != nil {
				return nil, err
			}
			return output, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				// Add detailed logging to inspect the response
				responseJSON, _ := json.Marshal(*resp)
				fmt.Println("Full API response: ", string(responseJSON))

				// Try different possible paths for PolicyId
				id, err = getSdkValue("policyId", *resp)
				if err != nil || id == nil {
					id, err = getSdkValue("data.policyId", *resp)
				}
				if err != nil {
					return fmt.Errorf("error getting PolicyId from response: %s", err)
				}
				if id == nil {
					return fmt.Errorf("PolicyId not found in response")
				}
				// Handle different numeric types
				switch v := id.(type) {
				case float64:
					d.SetId(fmt.Sprintf("%.0f", v))
				case int:
					d.SetId(strconv.Itoa(v))
				case int64:
					d.SetId(strconv.FormatInt(v, 10))
				case string:
					d.SetId(v)
				default:
					d.SetId(fmt.Sprintf("%v", v))
				}
			}
			return err
		},
	}
	return callback, err
}

// CreateAlarmPolicy creates an alarm policy
func (s *MonitorService) CreateAlarmPolicy(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, false)

	createCall, err := s.CreateAlarmPolicyCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createCall)
	return apiProcess.Run()
}

// CreateAlarmPolicyCall prepares the API call for creating alarm policy
func (s *MonitorService) CreateAlarmPolicyCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := make(map[string]interface{})

	// Basic fields - convert integers to strings for SDK
	req["PolicyName"] = d.Get("policy_name").(string)
	req["ProductType"] = d.Get("product_type").(int)
	req["PolicyType"] = d.Get("policy_type").(int)

	// Optional resource binding fields
	if v, ok := d.GetOk("resource_bind_type"); ok {
		req["ResourceBindType"] = v.(int)
	}
	if v, ok := d.GetOk("project_id"); ok {
		req["ProjectId"] = v.(int)
	}
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIds := v.([]interface{})
		if len(instanceIds) > 0 {
			req["InstanceIds"] = instanceIds
		}
	}

	// TriggerRules - required, serialize to JSON string
	if v, ok := d.GetOk("trigger_rules"); ok {
		triggerRules := v.([]interface{})
		rules := make([]map[string]interface{}, 0, len(triggerRules))
		for _, rule := range triggerRules {
			ruleMap := rule.(map[string]interface{})
			ruleData := make(map[string]interface{})
			ruleData["ItemKey"] = ruleMap["item_key"].(string)
			ruleData["ItemName"] = ruleMap["item_name"].(string)
			ruleData["EffectBT"] = ruleMap["effect_bt"].(string)
			ruleData["EffectET"] = ruleMap["effect_et"].(string)
			ruleData["Period"] = ruleMap["period"].(string)
			ruleData["Method"] = ruleMap["method"].(string)
			ruleData["Compare"] = ruleMap["compare"].(string)
			ruleData["TriggerValue"] = ruleMap["trigger_value"].(string)
			ruleData["Units"] = ruleMap["units"].(string)
			ruleData["Interval"] = ruleMap["interval"].(int)
			ruleData["MaxCount"] = ruleMap["max_count"].(int)
			ruleData["Points"] = ruleMap["points"].(int)
			rules = append(rules, ruleData)
		}
		req["TriggerRules"] = rules
	}

	// UserNotice - optional, serialize to JSON string
	if v, ok := d.GetOk("user_notice"); ok {
		userNotices := v.([]interface{})
		notices := make([]map[string]interface{}, 0, len(userNotices))
		for _, notice := range userNotices {
			noticeMap := notice.(map[string]interface{})
			noticeData := make(map[string]interface{})
			noticeData["ContactWay"] = noticeMap["contact_way"].(int)
			noticeData["ContactFlag"] = noticeMap["contact_flag"].(int)
			noticeData["ContactId"] = noticeMap["contact_id"].(int)
			notices = append(notices, noticeData)
		}
		req["UserNotice"] = notices
	}

	// URLNotice - optional
	if v, ok := d.GetOk("url_notice"); ok {
		urlNotices := v.([]interface{})
		urls := make([]string, 0, len(urlNotices))
		for _, url := range urlNotices {
			urls = append(urls, url.(string))
		}
		req["URLNotice"] = urls
	}

	return s.CreateAlarmPolicyCommonCall(req, true)
}

// ReadAndSetAlarmPolicy reads and sets alarm policy data
func (s *MonitorService) ReadAndSetAlarmPolicy(d *schema.ResourceData, r *schema.Resource) (err error) {
	params := map[string]interface{}{}
	params["PolicyId"] = d.Id()

	var data []interface{}
	data, err = s.ReadAlarmPolicy(params)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		d.SetId("")
		return nil
	}

	SdkResponseAutoResourceData(d, r, data, nil)
	return nil
}

// ReadAlarmPolicy reads alarm policy details
func (s *MonitorService) ReadAlarmPolicy(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.monitorv4conn
	action := "DescribeAlarmPolicy"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err = conn.DescribeAlarmPolicy(&condition)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("data", *resp)
	if err != nil {
		return data, err
	}
	data = append(data, results)
	return data, err
}

// DeleteAlarmPolicy deletes an alarm policy
func (s *MonitorService) DeleteAlarmPolicy(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteAlarmPolicyCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)
	return apiProcess.Run()
}

// DeleteAlarmPolicyCall prepares the API call for deleting alarm policy
func (s *MonitorService) DeleteAlarmPolicyCall(d *schema.ResourceData) (callback ApiCall, err error) {
	params := map[string]interface{}{}

	policyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return callback, fmt.Errorf("invalid policy id: %v", err)
	}
	params["PolicyIds"] = []int{policyId}

	callback = ApiCall{
		param:  &params,
		action: "DeleteAlarmPolicy",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.monitorv4conn
			// Create custom operation for DeleteAlarmPolicy
			op := &request.Operation{
				Name:       "DeleteAlarmPolicy",
				HTTPMethod: "POST",
				HTTPPath:   "/",
			}

			params := make(map[string]interface{})
			for k, v := range *(call.param) {
				params[k] = v
			}

			requestJSON, _ := json.Marshal(params)
			fmt.Println("Request params: ", string(requestJSON))

			output := &map[string]interface{}{}
			req := conn.NewRequest(op, &params, output)
			req.HTTPRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

			logger.Debug(logger.RespFormat, call.action, params)
			err = req.Send()
			if err != nil {
				return nil, err
			}
			return output, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}
