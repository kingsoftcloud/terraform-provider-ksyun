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

type emptySlice []interface{}

func (s *KcrsService) CreateKcrsInstance(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.createKcrsInstanceCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)

	if err := apiProcess.Run(); err != nil {
		return err
	}

	if d.Get("open_public_operation").(bool) {
		openCall, err := s.modifyExternalEndpointWithCall(d, true)
		if err != nil {
			return err
		}
		apiProcess.PutCalls(openCall)
	}

	if d.HasChange("external_policy") {
		ePolicyCall, err := s.modifyExternalEndpointPolicyWithCall(d)
		if err != nil {
			return err
		}
		apiProcess.PutCalls(ePolicyCall...)
	}

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

	apiProcess.PutCalls(IdleApiCall(10*time.Second), call)

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
	apiProcess.PutCalls(IdleApiCall(10*time.Second), call)

	return apiProcess.Run()
}
func (s *KcrsService) ModifyKcrsInstanceEoIEndpoint(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	if d.HasChange("open_public_operation") {
		isOpen := d.Get("open_public_operation").(bool)
		publicCall, err := s.modifyExternalEndpointWithCall(d, isOpen)
		if err != nil {
			return err
		}
		apiProcess.PutCalls(publicCall)
	}

	if d.HasChange("external_policy") {
		externalPolicyCall, err := s.modifyExternalEndpointPolicyWithCall(d)
		if err != nil {
			return err
		}
		apiProcess.PutCalls(externalPolicyCall...)
	}

	return apiProcess.Run()
}

func (s *KcrsService) createKcrsInstanceCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	trans := map[string]SdkReqTransform{
		"open_public_operation": {
			Ignore: true,
		},
		"delete_bucket": {
			Ignore: true,
		},
	}
	params, err := SdkRequestAutoMapping(d, r, false, trans, nil, SdkReqParameter{onlyTransform: false})
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

func (s *KcrsService) modifyExternalEndpointWithCall(d *schema.ResourceData, open bool) (callback ApiCall, err error) {
	params := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	callback = ApiCall{
		param:  &params,
		action: "OpenExternalEndpoint",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.OpenExternalEndpoint(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	if !open {
		callback.executeCall = func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CloseExternalEndpoint(call.param)
			return resp, err
		}
		callback.action = "CloseExternalEndpoint"
		callback.afterCall = func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if !open {
				// clean the external policy set when open_public_operation is false
				_ = d.Set("external_policy", emptySlice{})
			}
			return err
		}
	}

	return callback, err
}

func (s *KcrsService) modifyExternalEndpointPolicyWithCall(d *schema.ResourceData) (callbacks []ApiCall, err error) {
	createCall := func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
		conn := client.kcrsconn
		logger.Debug(logger.RespFormat, call.action, *(call.param))
		resp, err = conn.CreateExternalEndpointPolicy(call.param)
		return resp, err
	}
	removeCall := func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
		conn := client.kcrsconn
		logger.Debug(logger.RespFormat, call.action, *(call.param))
		resp, err = conn.DeleteExternalEndpointPolicy(call.param)
		return resp, err
	}

	rawIf, currIf := d.GetChange("external_policy")

	raw := rawIf.(*schema.Set)
	curr := currIf.(*schema.Set)

	removePolicy := raw.Difference(curr)
	addPolicy := curr.Difference(raw)

	callback := ApiCall{
		// param:  &params,
		action:      "",
		executeCall: nil,
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	policyCallHandler := func(set *schema.Set, isCreate bool) {

		for _, policy := range set.List() {
			var req = make(map[string]interface{})
			call := callback.Copy()

			req["InstanceId"] = d.Id()
			p := policy.(map[string]interface{})
			req["Entry"] = p["entry"]
			if isCreate {
				call.action = "CreateExternalEndpointPolicy"
				call.executeCall = createCall
				if v, ok := p["desc"]; ok {
					req["Desc"] = v
				}
			} else {
				call.action = "DeleteExternalEndpointPolicy"
				call.executeCall = removeCall
			}
			call.param = &req
			callbacks = append(callbacks, call)
		}
	}

	policyCallHandler(removePolicy, false)
	policyCallHandler(addPolicy, true)

	return callbacks, err
}

func (s *KcrsService) ModifyInternalEndpointDns(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	mCall, err := s.modifyInternalEndpointDnsWithCall(d, d.Get("enable_vpc_domain_dns").(bool))
	if err != nil {
		return err
	}
	apiProcess.PutCalls(mCall)

	return apiProcess.Run()
}

func (s *KcrsService) modifyInternalEndpointDnsWithCall(d *schema.ResourceData, enable bool) (callback ApiCall, err error) {
	ids := DisassembleIds(d.Id())

	params := map[string]interface{}{
		"InstanceId":          ids[0],
		"VpcId":               ids[1],
		"EniLBIp":             d.Get("eni_lb_ip"),
		"InternalEndpointDns": d.Get("internal_endpoint_dns"),
	}

	callback = ApiCall{
		param:  &params,
		action: "CreateInternalEndpointDns",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateInternalEndpointDns(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	if !enable {
		callback.executeCall = func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcrsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteInternalEndpointDns(call.param)
			return resp, err
		}
		callback.action = "DeleteInternalEndpointDns"
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
			keyPrefix := orderKey + "Key"
			params[keyPrefix] = header["key"]

			value := header["value"]
			valuePrefix := orderKey + "Value.1"
			params[valuePrefix] = value
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
			d.SetId(Float64ToString(id.(float64)))

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
	params["InstanceId"] = d.Get("instance_id")
	params["Namespace"] = d.Get("namespace")

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
				return helper.StringBoolean(d.Get("enable").(bool)), true
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
		"InstanceId":   d.Id(),
		"DeleteBucket": "False",
	}

	if d.Get("delete_bucket").(bool) {
		removeReq["DeleteBucket"] = "True"
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

func (s *KcrsService) ReadAndSetInternalEndpoint(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		var (
			instanceId = d.Get("instance_id").(string)
		)

		data, callErr := s.ReadInternalEndpoint(d, instanceId)
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
			extra := map[string]SdkResponseMapping{
				"EniLBIp": {
					Field: "eni_lb_ip",
				},
				"SubnetId": {
					Field: "reserve_subnet_id",
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)

			if d.Get("enable_vpc_domain_dns").(bool) {
				var (
					params      = make(map[string]interface{})
					endpointDns = "PrivateDomain"
					resp        *map[string]interface{}
					err         error
				)
				params["InstanceId"] = instanceId
				params["VpcId"] = d.Get("vpc_id")
				params["InternalEndpointDns"] = endpointDns
				params["EniLBIp"] = d.Get("eni_lb_ip")

				conn := s.client.kcrsconn
				action := "DescribeInternalEndpoint"
				logger.Debug(logger.ReqFormat, action, params)
				resp, err = conn.DescribeInternalEndpointDns(&params)
				if err != nil {
					if notFoundError(err) {
						return resource.NonRetryableError(err)
					}
					return resource.RetryableError(err)
				}
				status, err := getSdkValue("InternalEndpointDnsSet.0.Status", *resp)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				stat, _ := If2String(status)
				if err := d.Set("dns_parse_status", stat); err != nil {
					return resource.NonRetryableError(err)
				}
				_ = d.Set("internal_endpoint_dns", endpointDns)
			} else {
				if err := d.Set("dns_parse_status", "Closed"); err != nil {
					return resource.NonRetryableError(err)
				}
			}

			return nil
		}
	})
}

func (s *KcrsService) ReadAndSetKcrsTokens(d *schema.ResourceData, r *schema.Resource) (err error) {
	req := map[string]interface{}{
		"InstanceId": d.Get("instance_id"),
	}

	data, err := s.readKcrsInstanceTokens(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection: data,
		// nameField:   "InstanceName",
		idFiled:     "TokenId",
		targetField: "tokens",
		extra:       map[string]SdkResponseMapping{},
	})
}
func (s *KcrsService) ReadAndSetKcrsNamespaces(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"namespace": {
			mapping: "Namespace",
			Type:    TransformDefault,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}

	req["InstanceId"] = d.Get("instance_id")

	data, err := s.ReadKcrsNamespaces(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection: data,
		// nameField:   "InstanceName",
		idFiled:     "Namespace",
		targetField: "namespace_items",
		extra:       map[string]SdkResponseMapping{},
	})
}
func (s *KcrsService) ReadAndSetKcrsWebhookTriggers(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"trigger_id": {
			mapping: "TriggerId",
			Type:    TransformDefault,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}

	req["InstanceId"] = d.Get("instance_id")
	req["Namespace"] = d.Get("namespace")

	data, err := s.ReadKcrsWebhookTriggers(req)
	if err != nil {
		return err
	}

	for _, dd := range data {
		switch dd.(type) {
		case map[string]interface{}:
			dm := dd.(map[string]interface{})
			if v, ok := dm["TriggerId"]; ok {
				if vv, vOk := v.(float64); vOk {
					dm["TriggerId"] = Float64ToString(vv)
				}
			}
		}
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection: data,
		// nameField:   "InstanceName",
		idFiled:     "TriggerId",
		targetField: "triggers",
		extra:       map[string]SdkResponseMapping{},
	})
}
func (s *KcrsService) ReadAndSetKcrsInstances(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "InstanceId",
			Type:    TransformWithN,
		},
		"project_ids": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}

	if _, ok := req["ProjectId.1"]; !ok {
		if err := addProjectInfo(d, &req, s.client); err != nil {
			return err
		}
	}

	data, err := s.ReadKcrsInstances(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "InstanceName",
		idFiled:     "InstanceId",
		targetField: "instances",
		extra:       map[string]SdkResponseMapping{},
	})
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

			descExternalEndpoint := map[string]interface{}{
				"InstanceId": d.Id(),
			}

			conn := s.client.kcrsconn

			resp, err := conn.DescribeExternalEndpoint(&descExternalEndpoint)
			if err != nil {
				return resource.NonRetryableError(err)
			}

			results, err := getSdkValue("Status", *resp)
			if err != nil {
				return resource.NonRetryableError(err)
			}

			status := results.(string)
			if status == "Closed" {
				_ = d.Set("open_public_operation", false)
			} else {
				_ = d.Set("open_public_operation", true)
			}

			policySetIf, err := getSdkValue("PolicySet", *resp)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if policySetIf != nil {
				policySet := policySetIf.([]interface{})

				ps := make([]interface{}, 0, len(policySet))
				for _, pp := range policySet {
					switch pp.(type) {
					case map[string]interface{}:
						p := pp.(map[string]interface{})
						m := make(map[string]interface{}, 2)
						for k, v := range p {
							m[Hump2Downline(k)] = v
						}
						ps = append(ps, m)
					}
				}
				if len(ps) > 0 {
					_ = d.Set("external_policy", ps)
				}
			} else {
				_ = d.Set("external_policy", emptySlice{})
			}

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

func (s *KcrsService) ReadInternalEndpoint(d *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if instanceId == "" {
		return nil, fmt.Errorf("instance id cannot be blank")
	}
	req := map[string]interface{}{
		"InstanceId": instanceId,
	}

	results, err = s.ReadInternalEndpoints(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("kcrs internal endpoint %s not exist ", instanceId)
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

			headers := make([]map[string]interface{}, 0)
			// parse headers
			if receivedHeaderIf, ok := data["Header"]; ok {
				if receivedHeader, cOk := receivedHeaderIf.(string); cOk {
					kvs := strings.Split(receivedHeader, ";")
					for _, kv := range kvs {
						kvSlice := strings.Split(kv, ":")
						if len(kvSlice) < 2 {
							continue
						}
						k := kvSlice[0]
						v := kvSlice[1]
						m := make(map[string]interface{}, 2)
						m["key"] = k
						m["value"] = v
						headers = append(headers, m)
					}
				}
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

			triggers := []interface{}{trigger}
			if err := d.Set("trigger", triggers); err != nil {
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
	return pageQuery(condition, "MaxResults", "Marker", 99, 0, func(condition map[string]interface{}) ([]interface{}, error) {
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

	})
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
func (s *KcrsService) ReadInternalEndpoints(condition map[string]interface{}) (data []interface{}, err error) {

	var (
		resp    *map[string]interface{}
		results interface{}
	)

	if condition == nil {
		condition = make(map[string]interface{})
	}
	conn := s.client.kcrsconn
	action := "DescribeInternalEndpoint"
	logger.Debug(logger.ReqFormat, action, condition)
	resp, err = conn.DescribeInternalEndpoint(&condition)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("AccessVpcSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})

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
		Refresh:      s.kcrsInstanceStateRefreshFunc(d, instanceId, "Error"),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
		MinTimeout:   1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *KcrsService) kcrsInstanceStateRefreshFunc(d *schema.ResourceData, instanceId string, failStates ...string) resource.StateRefreshFunc {
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

		for _, v := range failStates {
			if v == status {
				return nil, "", fmt.Errorf("Kcrs Instance status  error, status:%v", status)
			}
		}

		return data, status.(string), nil
	}
}

func IdleApiCall(idle time.Duration) ApiCall {
	return ApiCall{
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error) {
			time.Sleep(idle)
			return nil, nil
		},
	}
}
