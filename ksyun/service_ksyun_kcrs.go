package ksyun

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type KcrsService struct {
	client *KsyunClient
}

func (s *KcrsService) CreateKcrsInstance(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.createKcrsInstanceCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}
func (s *KcrsService) CreateKcrsWebhookTrigger(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.createKcrsWebhookTriggerWithCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}

func (s *KcrsService) CreateKcrsToken(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.createKcrsTokenWithCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)

	if d.Get("enable").(bool) {
		modifyCall, err := s.modifyInstanceTokenStatusWithCall(d, r)
		if err != nil {
			return err
		}
		apiProcess.PutCalls(modifyCall)
	}

	return apiProcess.Run()
}

func (s *KcrsService) CreateKcrsNamespace(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.createKcrsNamespaceWithCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}

func (s *KcrsService) createKcrsInstanceCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	params, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}

	callback = ApiCall{
		param:  &params,
		action: "CreateInstance",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateInstance(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("InstanceId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			if err := s.checkKcrsInstanceState(d, d.Id(), []string{"Running"}, d.Timeout(schema.TimeoutCreate)); err != nil {
				return fmt.Errorf("waiting for kcrs instance caused an error: %s", err)
			}
			return err
		},
	}
	return callback, err
}

func prepareWebhookTriggerParameters(d *schema.ResourceData, r *schema.Resource, isUpdate bool) (map[string]interface{}, error) {
	t := map[string]SdkReqTransform{
		"trigger": {
			Type: TransformListUnique,
		},
	}

	params, err := SdkRequestAutoMapping(d, r, isUpdate, t, nil, SdkReqParameter{onlyTransform: false})
	if err != nil {
		return nil, err
	}

	// deal with event.types
	if _, ok := params["Trigger.EventTypes"]; ok {
		delete(params, "Trigger.EventTypes")

		for idx, event := range d.Get("trigger.0.event_types").(*schema.Set).List() {
			params["Trigger.EventType."+strconv.Itoa(idx+1)] = event
		}
	}

	if _, ok := params["Trigger.Headers"]; ok {
		delete(params, "Trigger.Headers")

		dMap, ok := helper.GetSchemaListHeadMap(d, "trigger")
		if !ok {
			return nil, errors.New("trigger attribute is not blank")
		}

		headers := dMap["headers"]

		prefix := "Trigger.Header."
		for idx, headerIf := range headers.([]interface{}) {
			header := headerIf.(map[string]interface{})
			orderKey := prefix + strconv.Itoa(idx+1) + "."

			// header keys
			params[orderKey+"Key"] = header["key"]

			values := header["values"]
			secondPrefix := orderKey + "Value"
			for j, value := range values.([]interface{}) {
				headerValueKey := secondPrefix + "." + strconv.Itoa(j+1)

				params[headerValueKey] = value
			}
		}

	}

	return params, nil
}

func (s *KcrsService) createKcrsWebhookTriggerWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	params, err := prepareWebhookTriggerParameters(d, r, false)
	if err != nil {
		return callback, err
	}

	callback = ApiCall{
		param:  &params,
		action: "CreateWebhookTrigger",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateWebhookTrigger(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("triggerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))

			return err
		},
	}
	return callback, err
}
func (s *KcrsService) modifyKcrsWebhookTriggerWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	params, err := prepareWebhookTriggerParameters(d, r, true)
	if err != nil {
		return callback, err
	}

	if len(params) < 1 {
		return callback, err
	}

	params["Trigger.TriggerId"] = d.Id()

	callback = ApiCall{
		param:  &params,
		action: "ModifyWebhookTrigger",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyWebhookTrigger(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *KcrsService) createKcrsTokenWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	params, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}

	callback = ApiCall{
		param:  &params,
		action: "CreateInstanceToken",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateInstanceToken(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("tokenId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))

			return err
		},
	}
	return callback, err
}

func (s *KcrsService) createKcrsNamespaceWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	trans := map[string]SdkReqTransform{
		"public": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				if d.Get("public").(bool) {
					return "True", true
				}
				return "False", true
			},
		},
	}

	params, err := SdkRequestAutoMapping(d, r, false, trans, nil, SdkReqParameter{onlyTransform: false})
	if err != nil {
		return callback, err
	}

	callback = ApiCall{
		param:  &params,
		action: "CreateNamespace",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateNamespace(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)

			instanceId := d.Get("instance_id").(string)
			namespace := d.Get("namespace").(string)
			id := AssembleIds(instanceId, namespace)
			d.SetId(id)

			return err
		},
	}
	return callback, err
}

func (s *KcrsService) modifyInstanceTokenStatus(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.modifyInstanceTokenStatusWithCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *KcrsService) ModifyWebhookTrigger(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.modifyKcrsWebhookTriggerWithCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *KcrsService) modifyInstanceTokenStatusWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	trans := map[string]SdkReqTransform{
		"enable": {
			ValueFunc: func(d *schema.ResourceData) (interface{}, bool) {
				if d.Get("enable").(bool) {
					return "True", true
				}
				return "False", true
			},
		},
	}

	params, err := SdkRequestAutoMapping(d, r, true, trans, nil)
	if err != nil {
		return callback, err
	}

	if len(params) < 1 {
		return
	}

	params["InstanceId"] = d.Get("instance_id")
	params["TokenId"] = d.Id()

	callback = ApiCall{
		param:  &params,
		action: "ModifyInstanceTokenStatus",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyInstanceTokenStatus(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)

			return err
		},
	}
	return callback, err
}

func (s *KcrsService) deleteKcrsInstanceWithCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteInstance",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteInstance(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadKcrsInstance(d, d.Id())
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading kcrs instance when delete %q, %s", d.Id(), callErr))
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
	return callback, err
}

func (s *KcrsService) RemoveKcrsWebhookTrigger(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.deleteKcrsWebhookTriggerWithCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}

func (s *KcrsService) deleteKcrsWebhookTriggerWithCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"InstanceId": d.Get("instance_id"),
		"Namespace":  d.Get("namespace"),
		"TriggerId":  d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteWebhookTrigger",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteWebhookTrigger(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadKcrsInstance(d, d.Id())
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading kcrs webhook trigger when delete %q, %s", d.Id(), callErr))
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
	return callback, err
}

func (s *KcrsService) deleteKcrsNamespaceWithCall(d *schema.ResourceData) (callback ApiCall, err error) {
	insIds := DisassembleIds(d.Id())

	removeReq := map[string]interface{}{
		"InstanceId": insIds[0],
		"Namespace":  insIds[1],
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteNamespace",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteNamespace(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadKcrsNamespace(d, insIds[0], insIds[1])
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading kcrs namespace when delete %q, %s", d.Id(), callErr))
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
	return callback, err
}

func (s *KcrsService) ReadAndSetKcrsInstance(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadKcrsInstance(d, d.Id())
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading kcrs instance %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, nil)
			return nil
		}
	})
}
func (s *KcrsService) ReadAndSetKcrsInstanceToken(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadKcrsInstanceToken(d, d.Get("instance_id").(string))
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading kcrs instance %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, nil)
			return nil
		}
	})
}

func (s *KcrsService) ReadAndSetKcrsNamespace(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		insIds := DisassembleIds(d.Id())
		data, callErr := s.ReadKcrsNamespace(d, insIds[0], insIds[1])
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading kcrs namespace %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, nil)
			return nil
		}
	})
}

func (s *KcrsService) RemoveKcrsInstance(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.deleteKcrsInstanceWithCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}

func (s *KcrsService) RemoveKcrsNamespace(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.deleteKcrsNamespaceWithCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}

func (s *KcrsService) ReadKcrsInstance(d *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"InstanceId.1": instanceId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}
	results, err = s.ReadKcrsInstances(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kcrs instance %s not exist ", instanceId)
	}
	return data, err
}

func (s *KcrsService) ReadAndSetWebhookTrigger(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadKcrsWebhookTrigger(d, "", "", "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading kcrs namespace %q, %s", d.Id(), callErr))
			}
		} else {

			trigger := make(map[string]interface{}, len(data))

			// parse headers
			headersM := make(map[string][]string)
			if receivedHeaderIf, ok := data["Header"]; ok {
				if receivedHeader, cOk := receivedHeaderIf.(string); cOk {
					kvs := strings.Split(receivedHeader, ";")
					for _, kv := range kvs {
						kvSlice := strings.Split(kv, ":")
						if len(kvSlice) < 2 {
							continue
						}
						k := kvSlice[0]
						if _, ok := headersM[k]; !ok {
							headersM[k] = make([]string, 0)
						}

						headersM[k] = append(headersM[k], kvSlice[1])

					}
				}
			}

			headers := make([]map[string]interface{}, 0, len(headersM))
			for hk, hvs := range headersM {
				m := make(map[string]interface{}, 2)
				m["key"] = hk
				m["values"] = hvs
				headers = append(headers, m)
			}

			trigger["headers"] = headers

			// set trigger
			for k, v := range data {
				switch k {
				case "EventType":
					trigger["event_types"] = v
				case "TriggerName", "TriggerUrl", "Enabled":
					trigger[Hump2Downline(k)] = v
				}
			}
			if err := d.Set("trigger", trigger); err != nil {
				return resource.NonRetryableError(err)
			}
			return nil
		}
	})
}

func (s *KcrsService) ReadKcrsWebhookTrigger(d *schema.ResourceData, instanceId, triggerId, namespace string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if instanceId == "" {
		instanceId = d.Get("instance_id").(string)
	}

	if triggerId == "" {
		triggerId = d.Id()
	}

	if namespace == "" {
		namespace = d.Get("namespace").(string)
	}

	req := map[string]interface{}{
		"InstanceId": instanceId,
		"Namespace":  namespace,
		"TriggerId":  triggerId,
	}

	results, err = s.ReadKcrsWebhookTriggers(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kcrs webhook trigger %s not exist ", instanceId)
	}
	return data, err
}

func (s *KcrsService) ReadKcrsInstanceToken(d *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if instanceId == "" {
		instanceId = d.Get("instance_id").(string)
	}
	req := map[string]interface{}{
		"InstanceId": instanceId,
	}

	results, err = s.readKcrsInstanceTokens(req)
	if err != nil {
		return data, err
	}

	for _, v := range results {
		temp := v.(map[string]interface{})
		if v, ok := temp["TokenId"]; ok && v == d.Id() {
			data = temp
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kcrs instance %s not exist ", instanceId)
	}
	return data, err
}

func (s *KcrsService) ReadKcrsNamespace(d *schema.ResourceData, instanceId, namespace string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if instanceId == "" {
		instanceId = d.Get("instance_id").(string)
	}

	if namespace == "" {
		namespace = d.Get("namespace").(string)
	}

	req := map[string]interface{}{
		"InstanceId": instanceId,
		"Namespace":  namespace,
	}

	results, err = s.ReadKcrsNamespaces(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kcrs namespace %s not exist ", instanceId)
	}
	return data, err
}

func (s *KcrsService) ReadKcrsInstances(condition map[string]interface{}) (data []interface{}, err error) {

	var (
		resp    *map[string]interface{}
		results interface{}
	)
	if condition == nil {
		condition = make(map[string]interface{})
	}
	conn := s.client.kcrsconn
	action := "DescribeInstance"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err = conn.DescribeInstance(&condition)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("InstanceSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})

	return data, err
}

func (s *KcrsService) ReadKcrsWebhookTriggers(condition map[string]interface{}) (data []interface{}, err error) {

	var (
		resp    *map[string]interface{}
		results interface{}
	)
	data, err = pageQuery(condition, "MaxResults", "Maker", 10, 0, func(condition map[string]interface{}) ([]interface{}, error) {

		if condition == nil {
			condition = make(map[string]interface{})
		}
		conn := s.client.kcrsconn
		action := "DescribeWebhookTrigger"
		logger.Debug(logger.ReqFormat, action, condition)
		resp, err = conn.DescribeWebhookTrigger(&condition)
		if err != nil {
			return data, err
		}

		results, err = getSdkValue("TriggerSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})

		return data, err
	})

	return data, err
}

func (s *KcrsService) readKcrsInstanceTokens(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "Marker", 99, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kcrsconn
		action := "DescribeInstanceToken"
		logger.Debug(logger.ReqFormat, action, condition)

		resp, err = conn.DescribeInstanceToken(&condition)
		if err != nil {
			return data, err
		}

		results, err = getSdkValue("TokenSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})

		return data, err
	})
}

func (s *KcrsService) ReadKcrsNamespaces(condition map[string]interface{}) (data []interface{}, err error) {

	var (
		resp    *map[string]interface{}
		results interface{}
	)
	if condition == nil {
		condition = make(map[string]interface{})
	}
	conn := s.client.kcrsconn
	action := "DescribeNamespace"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err = conn.DescribeNamespace(&condition)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("NamespaceSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err

}

func (s *KcrsService) checkKcrsInstanceState(d *schema.ResourceData, instanceId string, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{},
		Target:       target,
		Refresh:      s.kcrsInstanceStateRefreshFunc(d, instanceId),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
		MinTimeout:   1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *KcrsService) kcrsInstanceStateRefreshFunc(d *schema.ResourceData, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := s.ReadKcrsInstance(d, instanceId)
		if err != nil {
			return nil, "", err
		}

		status, err := getSdkValue("InstanceStatus", data)
		if err != nil {
			return nil, "", err
		}

		return data, status.(string), nil
	}
}
