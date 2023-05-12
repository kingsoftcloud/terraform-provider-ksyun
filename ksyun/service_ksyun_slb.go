package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"strings"
	"time"
)

type SlbService struct {
	client *KsyunClient
}

//start slb

func (s *SlbService) ReadLoadBalancers(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeLoadBalancers"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeLoadBalancers(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeLoadBalancers(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("LoadBalancerDescriptions", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})

	// merge attrs
	//describeLoadBalancerAttributes
	s.describeLoadBalancersAttributes(data)

	return data, err
}

func (s *SlbService) ReadLoadBalancer(d *schema.ResourceData, loadBalancerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if loadBalancerId == "" {
		loadBalancerId = d.Id()
	}
	req := map[string]interface{}{
		"LoadBalancerId.1": loadBalancerId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}
	results, err = s.ReadLoadBalancers(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("LoadBalancer %s not exist ", loadBalancerId)
	}
	return data, err
}

// 读取attributes，增加日志配置字段
func (s *SlbService) describeLoadBalancersAttributes(lbs []interface{}) {

	apiNotSuppotRegion := false
	for _, lb := range lbs {
		lbItem := lb.(map[string]interface{})

		// 当前region如果已经报过错，不再查接口，直接给日志设为false
		if apiNotSuppotRegion {
			lbItem["AccessLogsEnabled"] = false
			lbItem["AccessLogsS3Bucket"] = nil
			continue
		}
		err := s.describeLoadBalancerAttributes(&lbItem)

		// 不支持的region给日志设为false
		if err != nil && strings.Contains(err.Error(), "ApiNotSupportRegion") {
			apiNotSuppotRegion = true
			lbItem["AccessLogsEnabled"] = false
			lbItem["AccessLogsS3Bucket"] = nil //"ApiNotSupportRegion"
		}
	}

}

// 读取单个lb的attributes
func (s *SlbService) describeLoadBalancerAttributes(lb *map[string]interface{}) (err error) {
	defer func() {
		e := recover()
		logger.Debug("recover err: %s %s", "describeLoadBalancerAttributes", e)
	}()
	conn := s.client.slbconn

	params := map[string]interface{}{
		"LoadBalancerId": (*lb)["LoadBalancerId"],
	}
	logger.Debug("%s, %s", "describeLoadBalancerAttributes", params)
	var resp *map[string]interface{}
	resp, err = conn.DescribeLoadBalancerAttributes(&params)
	logger.Debug("%s, %s %s", "describeLoadBalancerAttributes resp", resp, err)

	if err != nil {
		return err
	}

	if v, ok := (*resp)["LoadBalancerAttributeSet"]; ok {
		attributes := v.([]interface{})
		logger.Debug("%s, %s", "describeLoadBalancerAttributes ok", attributes)
		for _, attr := range attributes {

			item := attr.(map[string]interface{})

			logger.Debug("%s, %s", "describeLoadBalancerAttributes k:v", attr)
			if item["Key"].(string) == "access_logs.s3.enabled" {
				if item["Value"].(string) == "false" {
					(*lb)["AccessLogsEnabled"] = false
				}
				if item["Value"].(string) == "true" {
					(*lb)["AccessLogsEnabled"] = true
				}
			}
			if item["Key"].(string) == "access_logs.s3.bucket" {
				//d.Set("log_bucket", item["Value"])
				(*lb)["AccessLogsS3Bucket"] = item["Value"]
			}
		}
	} else {
		logger.Debug("%s, %s", "describeLoadBalancerAttributes not ok", "")
	}
	return
}

func (s *SlbService) ReadAndSetLoadBalancer(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadLoadBalancer(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading lb %q, %s", d.Id(), callErr))
			}
		} else {
			err = mergeTagsData(d, &data, s.client, "loadbalancer")
			if err != nil {
				return resource.NonRetryableError(err)
			}

			extra := map[string]SdkResponseMapping{
				"ProjectId": {
					Field: "project_id",
					FieldRespFunc: func(i interface{}) interface{} {
						v, _ := strconv.Atoi(i.(string))
						return v
					},
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (s *SlbService) ReadAndSetLoadBalancers(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "LoadBalancerId",
			Type:    TransformWithN,
		},
		"project_id": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadLoadBalancers(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "LoadBalancerName",
		idFiled:     "LoadBalancerId",
		targetField: "lbs",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateLoadBalancerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"load_balancer_state": {
			mapping: "AdminStateUp",
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				if data.Get("load_balancer_state") == "start" {
					return true, true
				} else {
					return false, true
				}
			},
		},
		"tags": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	if req["Type"] == "internal" {
		if _, ok := req["SubnetId"]; !ok {
			return callback, fmt.Errorf(" Lb type is internal,must set SubnetId")
		}
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateLoadBalancer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateLoadBalancer(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("LoadBalancerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateLoadBalancer(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateLoadBalancerCall(d, r)
	if err != nil {
		return err
	}
	tagService := TagService{s.client}
	tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, r, "loadbalancer", false, true)
	if err != nil {
		return err
	}

	attributesCall, err := s.modifyLoadBalancerAttributesCall(d, r, false)
	if err != nil {
		return err
	}

	return ksyunApiCallNew([]ApiCall{call, tagCall, attributesCall}, d, s.client, true)
}

func (s *SlbService) ModifyLoadBalancerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id":            {Ignore: true},
		"tags":                  {Ignore: true},
		"access_logs_enabled":   {Ignore: true},
		"access_logs_s3_bucket": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["LoadBalancerId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyLoadBalancer",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyLoadBalancer(call.param)
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

func (s *SlbService) ModifyLoadBalancerProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		callback = ApiCall{
			param: &updateReq,
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				return resp, ModifyProjectInstanceNew(d.Id(), call.param, client)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				return err
			},
		}
	}
	return callback, err
}

func (s *SlbService) modifyLoadBalancerAttributesCall(d *schema.ResourceData, resource *schema.Resource, isUpdate bool) (callback ApiCall, err error) {
	if isUpdate && !d.HasChange("access_logs_enabled") && !d.HasChange("access_logs_s3_bucket") {
		return callback, err
	}

	if !isUpdate && !d.Get("access_logs_enabled").(bool) {
		return callback, err
	}

	params := map[string]interface{}{
		"LoadBalancerId":            d.Id(),
		"Attributes.member.1.Key":   "access_logs.s3.enabled",
		"Attributes.member.1.Value": d.Get("access_logs_enabled"),
		"Attributes.member.2.Key":   "access_logs.s3.bucket",
		"Attributes.member.2.Value": d.Get("access_logs_s3_bucket"),
	}

	callback = ApiCall{
		param: &params,
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			// 新建的资源，在这里才能拿到刚建好的实例信息，id在这里加到参数里
			if !isUpdate && d.Id() != "" {
				(*call.param)["LoadBalancerId"] = d.Id()
			}
			_, err = client.slbconn.ModifyLoadBalancerAttributes(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			return err
		},
	}
	return callback, err
}

func (s *SlbService) ModifyLoadBalancer(d *schema.ResourceData, r *schema.Resource) (err error) {
	projectCall, err := s.ModifyLoadBalancerProjectCall(d, r)
	if err != nil {
		return err
	}
	call, err := s.ModifyLoadBalancerCall(d, r)
	if err != nil {
		return err
	}
	tagService := TagService{s.client}
	tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, r, "loadbalancer", true, false)
	if err != nil {
		return err
	}

	attributesCall, err := s.modifyLoadBalancerAttributesCall(d, r, true)
	if err != nil {
		return err
	}

	return ksyunApiCallNew([]ApiCall{projectCall, call, tagCall, attributesCall}, d, s.client, true)
}

func (s *SlbService) RemoveLoadBalancerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteLoadBalancer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteLoadBalancer(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLoadBalancer(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveLoadBalancer(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLoadBalancerCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

//start listener

func (s *SlbService) ReadListeners(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeListeners"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeListeners(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeListeners(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("ListenerSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadListener(d *schema.ResourceData, listenerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if listenerId == "" {
		listenerId = d.Id()
	}
	req := map[string]interface{}{
		"ListenerId.1": listenerId,
	}
	results, err = s.ReadListeners(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Lb listener %s not exist ", listenerId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadListener(d, "")
	if err != nil {
		return err
	}

	// 如果没有返回acl信息，手动加上这个字段（否则在控制台解绑后，tf无法正常获取到绑定状态的更改）
	if _, ok := data["LoadBalancerAclId"]; !ok {
		data["LoadBalancerAclId"] = nil
	}

	//logger.Debug("test", "ReadAndSetListener", data)
	extra := map[string]SdkResponseMapping{
		"Session": {
			Field: "session",
			FieldRespFunc: func(i interface{}) interface{} {
				return []interface{}{
					i,
				}
			},
		},
		"HealthCheck": {
			Field: "health_check",
			FieldRespFunc: func(i interface{}) interface{} {
				if len(i.(map[string]interface{})) > 0 {
					return []interface{}{
						i,
					}
				}
				return nil
			},
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return err
}

func (s *SlbService) ReadAndSetListeners(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "ListenerId",
			Type:    TransformWithN,
		},
		"load_balancer_id": {
			mapping: "load-balancer-id",
			Type:    TransformWithFilter,
		},
		"certificate_id": {
			mapping: "certificate-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadListeners(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "ListenerName",
		idFiled:     "ListenerId",
		targetField: "listeners",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"enable_http2": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return data.Get("enable_http2"), true
			},
		},
		"session": {
			Type: TransformListUnique,
		},
		"health_check": {
			Ignore: true,
		},
		"load_balancer_acl_id": {
			Ignore: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	if req["listener_protocol"] != "HTTPS" {
		delete(req, "EnableHttp2")
	}
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
	}
	// if session is zero need set default SessionState stop
	if _, ok := req["SessionState"]; !ok {
		req["SessionState"] = "stop"
	}

	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateListeners",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateListeners(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("ListenerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("listener_id", d.Id())
		},
	}
	return callback, err
}

func (s *SlbService) CreateHealthCheckWithListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return callback, err
	}
	for k, v := range req {
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	if d.Get("listener_protocol") != "HTTP" && d.Get("listener_protocol") != "HTTPS" {
		delete(req, "UrlPath")
		delete(req, "HostName")
		delete(req, "IsDefaultHostName")
	}
	if len(req) > 0 {
		return s.CreateHealthCheckCommonCall(req, false)
	}
	return callback, err
}

func (s *SlbService) CreateListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateListenerCall(d, r)
	if err != nil {
		return err
	}
	health, err := s.CreateHealthCheckWithListenerCall(d, r)
	if err != nil {
		return err
	}
	var acl ApiCall
	if v, ok := d.GetOk("load_balancer_acl_id"); ok {
		acl, err = s.CreateLoadBalancerAclAssociateWithListenerCall(d, r, v.(string))
	}
	return ksyunApiCallNew([]ApiCall{call, health, acl}, d, s.client, false)
}

func (s *SlbService) ModifyHealthCheckWithListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	for k, v := range req {
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	//special
	req["HealthCheckState"] = d.Get("health_check.0.health_check_state")
	if d.Get("listener_protocol") == "HTTP" || d.Get("listener_protocol") == "HTTPS" {
		req["UrlPath"] = d.Get("health_check.0.url_path")
	}
	return s.ModifyHealthCheckCommonCall(req, d.Get("health_check.0.health_check_id").(string))
}

func (s *SlbService) ModifyListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"session": {
			Type: TransformListUnique,
		},
		"health_check": {
			Ignore: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	//特殊处理下"Session."
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
	}
	if len(req) > 0 {
		req["ListenerId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyListeners",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyListeners(call.param)
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

func (s *SlbService) ModifyListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyListenerCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	var healthCheckCall ApiCall
	if len(d.Get("health_check").([]interface{})) > 0 && d.Get("health_check.0.health_check_id") != "" {
		healthCheckCall, err = s.ModifyHealthCheckWithListenerCall(d, r)
	} else {
		healthCheckCall, err = s.CreateHealthCheckWithListenerCall(d, r)
	}
	if err != nil {
		return err
	}
	callbacks = append(callbacks, healthCheckCall)
	if d.HasChange("load_balancer_acl_id") {
		o, n := d.GetChange("load_balancer_acl_id")
		if d.Get("load_balancer_acl_id") == "" {
			var aclRemoveCall ApiCall
			aclRemoveCall, err = s.RemoveLoadBalancerAclAssociateCommonCall(d.Get("listener_id").(string), o.(string))
			if err != nil {
				return err
			}
			callbacks = append(callbacks, aclRemoveCall)
		} else {
			if o == "" {
				var aclAddCall ApiCall
				aclAddCall, err = s.CreateLoadBalancerAclAssociateWithListenerCall(d, r, n.(string))
				if err != nil {
					return err
				}
				callbacks = append(callbacks, aclAddCall)
			} else {
				var aclAddCall ApiCall
				var aclRemoveCall ApiCall
				aclRemoveCall, err = s.RemoveLoadBalancerAclAssociateCommonCall(d.Get("listener_id").(string), o.(string))
				if err != nil {
					return err
				}
				callbacks = append(callbacks, aclRemoveCall)
				aclAddCall, err = s.CreateLoadBalancerAclAssociateWithListenerCall(d, r, n.(string))
				if err != nil {
					return err
				}
				callbacks = append(callbacks, aclAddCall)
			}
		}
	}
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) RemoveListenerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"ListenerId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteListeners",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteListeners(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadListener(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading health check when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveListener(d *schema.ResourceData) (err error) {
	call, err := s.RemoveListenerCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

//start healthCheck

func (s *SlbService) ReadLoadHealthChecks(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeHealthChecks"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeHealthChecks(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeHealthChecks(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("HealthCheckSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadLoadHealthCheck(d *schema.ResourceData, healthCheckId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if healthCheckId == "" {
		healthCheckId = d.Id()
	}
	req := map[string]interface{}{
		"HealthCheckId.1": healthCheckId,
	}
	results, err = s.ReadLoadHealthChecks(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("HealthCheck %s not exist ", healthCheckId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetHealCheck(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLoadHealthCheck(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	listener, err := s.ReadListener(nil, data["ListenerId"].(string))
	if err != nil {
		return err
	}
	if listener["ListenerProtocol"] == "HTTP" || listener["ListenerProtocol"] == "HTTPS" {
		if data["HostName"] == nil {
			err = d.Set("is_default_host_name", true)
		} else {
			err = d.Set("is_default_host_name", false)
		}
		if err != nil {
			return err
		}
	}
	return d.Set("listener_protocol", listener["ListenerProtocol"])
}

func (s *SlbService) ReadAndSetHealthChecks(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "HealthCheckId",
			Type:    TransformWithN,
		},
		"listener_id": {
			mapping: "listener-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadLoadHealthChecks(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "HealthCheckId",
		targetField: "health_checks",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateHealthCheckCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "ConfigureHealthCheck",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			if _, ok := (*(call.param))["ListenerId"]; !ok {
				(*(call.param))["ListenerId"] = d.Get("listener_id")
			}
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ConfigureHealthCheck(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if isSetId {
				var id interface{}
				id, err = getSdkValue("HealthCheckId", *resp)
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
func (s *SlbService) CreateHealthCheckCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateHealthCheckCommonCall(req, true)
}

func (s *SlbService) CreateHealthCheck(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateHealthCheckCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) ModifyHealthCheckCommonCall(req map[string]interface{}, healthCheckId string) (callback ApiCall, err error) {
	if len(req) > 0 {
		req["HealthCheckId"] = healthCheckId
		callback = ApiCall{
			param:  &req,
			action: "ModifyHealthCheck",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyHealthCheck(call.param)
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

func (s *SlbService) ModifyHealthCheckCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check_state": {
			forceUpdateParam: true,
		},
		"url_path": {
			forceUpdateParam: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if d.Get("listener_protocol") != "HTTPS" && d.Get("listener_protocol") != "HTTP" {
		delete(req, "UrlPath")
	}
	if err != nil {
		return callback, err
	}
	return s.ModifyHealthCheckCommonCall(req, d.Id())
}

func (s *SlbService) ModifyHealthCheck(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyHealthCheckCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) RemoveHealthCheckCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"HealthCheckId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteHealthCheck",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteHealthCheck(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLoadHealthCheck(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading health check when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveHealthCheck(d *schema.ResourceData) (err error) {
	call, err := s.RemoveHealthCheckCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// LbRule start

func (s *SlbService) ReadLbRules(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeRules"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeRules(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeRules(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("RuleSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadLbRule(d *schema.ResourceData, lbRuleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if lbRuleId == "" {
		lbRuleId = d.Id()
	}
	req := map[string]interface{}{
		"RuleId.1": lbRuleId,
	}
	results, err = s.ReadLbRules(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Lb Rule %s not exist ", lbRuleId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetLbRule(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLbRule(d, "")
	if err != nil {
		return err
	}
	hostname, _ := getSdkValue("HealthCheck.HostName", data)
	if hostname == nil {
		data["HealthCheck"].(map[string]interface{})["IsDefaultHostName"] = true
	} else {
		data["HealthCheck"].(map[string]interface{})["IsDefaultHostName"] = false
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	if err != nil {
		return err
	}
	return err
}

func (s *SlbService) ReadAndSetLbRules(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "RuleId",
			Type:    TransformWithN,
		},
		"host_header_id": {
			mapping: "host-header-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadLbRules(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "RuleId",
		targetField: "lb_rules",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateLbRuleCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"session": {
			Type: TransformListUnique,
		},
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	//特殊处理下"Session." 和 "HealthCheck."
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	//非同步模式下 需要检查一下
	if req["ListenerSync"] == "off" && req["CookieType"] == "RewriteCookie" {
		if _, ok := req["CookieName"]; !ok {
			return callback, fmt.Errorf("Session.CookieType  is RewriteCookie,must set CookieName")
		}
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateSlbRule",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateSlbRule(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			id, err = getSdkValue("Rule.RuleId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateLbRule(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateLbRuleCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) ModifyLbRuleCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"session": {
			Type: TransformListUnique,
		},
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	//特殊处理下"Session." 和 "HealthCheck."
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	if len(req) > 0 {
		req["RuleId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifySlbRule",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifySlbRule(call.param)
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

func (s *SlbService) ModifyLbRule(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyLbRuleCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) RemoveLbRuleCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"RuleId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteRule",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteRule(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLbRule(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb rule when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveLbRule(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLbRuleCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// start host header

func (s *SlbService) ReadHostHeaders(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeHostHeaders"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeHostHeaders(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeHostHeaders(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("HostHeaderSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadHostHeader(d *schema.ResourceData, hostHeaderId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if hostHeaderId == "" {
		hostHeaderId = d.Id()
	}
	req := map[string]interface{}{
		"HostHeaderId.1": hostHeaderId,
	}
	results, err = s.ReadHostHeaders(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Host header %s not exist ", hostHeaderId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetHostHeader(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadHostHeader(d, "")
	if err != nil {
		return err
	}
	listenerId := data["ListenerId"].(string)
	listener, err := s.ReadListener(nil, listenerId)
	if err != nil {
		return err
	}
	data["ListenerProtocol"] = listener["ListenerProtocol"]
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetHostHeaders(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "HostHeaderId",
			Type:    TransformWithN,
		},
		"listener_id": {
			mapping: "listener-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadHostHeaders(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "HostHeaderId",
		targetField: "host_headers",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateHostHeaderCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	//获取一次监听器信息 1判定协议 2 校验是需要证书
	vip, err := s.ReadListener(nil, req["ListenerId"].(string))
	if err != nil {
		return callback, err
	}
	if vip["ListenerProtocol"] != "HTTP" && vip["ListenerProtocol"] != "HTTPS" {
		return callback, fmt.Errorf("Listener Protocol must HTTP or HTTPS when create a HostHeader ")
	}
	if _, ok := req["CertificateId"]; !ok && vip["ListenerProtocol"] == "HTTPS" {
		return callback, fmt.Errorf("CertificateId must set because Listener Protocol is HTTPS when create a HostHeader ")
	}
	if vip["ListenerProtocol"] == "HTTP" {
		delete(req, "CertificateId")
	}

	callback = ApiCall{
		param:  &req,
		action: "CreateHostHeader",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateHostHeader(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("HostHeader.HostHeaderId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateHostHeader(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateHostHeaderCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) ModifyHostHeaderCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["HostHeaderId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyHostHeader",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyHostHeader(call.param)
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

func (s *SlbService) ModifyHostHeader(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyHostHeaderCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) RemoveHostHeaderCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"HostHeaderId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteHostHeader",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteHostHeader(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadHostHeader(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading host header when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveHostHeader(d *schema.ResourceData) (err error) {
	call, err := s.RemoveHostHeaderCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// start LoadBalancerAcl

func (s *SlbService) ReadLoadBalancerAcls(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 5, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.slbconn
		action := "DescribeLoadBalancerAcls"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeLoadBalancerAcls(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeLoadBalancerAcls(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("LoadBalancerAclSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *SlbService) ReadLoadBalancerAcl(d *schema.ResourceData, loadBalancerAclId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if loadBalancerAclId == "" {
		loadBalancerAclId = d.Id()
	}
	req := map[string]interface{}{
		"LoadBalancerAclId.1": loadBalancerAclId,
	}
	results, err = s.ReadLoadBalancerAcls(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("LoadBalancerAcls %s not exist ", loadBalancerAclId)
	}
	return data, err
}

func (s *SlbService) ReadLoadBalancerAclEntry(d *schema.ResourceData, loadBalancerAclId string) (data map[string]interface{}, err error) {
	acl, err := s.ReadLoadBalancerAcl(d, loadBalancerAclId)
	if err != nil {
		return data, err
	}
	num := int64(d.Get("rule_number").(int))
	cidr := d.Get("cidr_block").(string)
	found := false
	for _, entry := range acl["LoadBalancerAclEntrySet"].([]interface{}) {
		m := entry.(map[string]interface{})
		if num == int64(m["RuleNumber"].(float64)) && cidr == m["CidrBlock"] {
			found = true
			data = m
			break
		}
	}
	if !found {
		return data, fmt.Errorf("LoadBalancerAclEntry not exist")
	}
	return data, err
}

func (s *SlbService) ReadLoadBalancerAclAssociate(listenerId string, loadBalancerAclId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	req := map[string]interface{}{
		"ListenerId.1": listenerId,
	}
	results, err = s.ReadListeners(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf(" LoadBalancerAclAssociate listener_id [%s] and load_balancer_acl_id [%s] not exist ",
			listenerId, loadBalancerAclId)
	}
	if _, ok := data["LoadBalancerAclId"]; !ok || data["LoadBalancerAclId"] != loadBalancerAclId {
		return data, fmt.Errorf(" LoadBalancerAclAssociate listener_id [%s] and load_balancer_acl_id [%s] not exist ",
			listenerId, loadBalancerAclId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetLoadBalancerAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLoadBalancerAcl(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetLoadBalancerAclEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLoadBalancerAclEntry(d, d.Get("load_balancer_acl_id").(string))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetLoadBalancerAclAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLoadBalancerAclAssociate(d.Get("listener_id").(string), d.Get("load_balancer_acl_id").(string))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetLoadBalancerAcls(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "LoadBalancerAclId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadLoadBalancerAcls(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "LoadBalancerAclId",
		nameField:   "LoadBalancerAclName",
		targetField: "lb_acls",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateLoadBalancerAclCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"load_balancer_acl_entry_set": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateLoadBalancerAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateLoadBalancerAcl(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("LoadBalancerAcl.LoadBalancerAclId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("load_balancer_acl_id", d.Id())
		},
	}
	return callback, err
}

func (s *SlbService) CreateLoadBalancerAclEntryCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateLoadBalancerAclEntryCommonCall(req, true)
}

func (s *SlbService) CreateLoadBalancerAclEntryWithAclCall(d *schema.ResourceData, r *schema.Resource) (callbacks []ApiCall, err error) {
	if entries, ok := d.GetOk("load_balancer_acl_entry_set"); ok {
		//check
		if len(schema.NewSet(loadBalancerAclEntryNumberHash, entries.(*schema.Set).List()).List()) != len(entries.(*schema.Set).List()) {
			return callbacks, fmt.Errorf("RuleNumber must unique ")
		}
		if len(schema.NewSet(loadBalancerAclEntryCidrHash, entries.(*schema.Set).List()).List()) != len(entries.(*schema.Set).List()) {
			return callbacks, fmt.Errorf("CidrBlock must unique ")
		}
		for _, entry := range entries.(*schema.Set).List() {
			var (
				req      map[string]interface{}
				callback ApiCall
			)
			transform := make(map[string]SdkReqTransform)
			key := strconv.Itoa(loadBalancerAclEntryHash(entry))
			for k := range entry.(map[string]interface{}) {
				key := "load_balancer_acl_entry_set." + key + "." + k
				transform[key] = SdkReqTransform{mapping: Downline2Hump(k)}
			}
			logger.Debug(logger.RespFormat, "Demo", d.Get("load_balancer_acl_entry_set"))
			req, err = SdkRequestAutoMapping(d, r, false, transform, nil)
			if err != nil {
				return callbacks, err
			}
			logger.Debug(logger.RespFormat, "Demo", req)
			callback, err = s.CreateLoadBalancerAclEntryCommonCall(req, false)
			if err != nil {
				return callbacks, err
			}
			callbacks = append(callbacks, callback)
		}
	}
	return callbacks, err
}

func (s *SlbService) CreateLoadBalancerAclEntryCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateLoadBalancerAclEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			(*(call.param))["LoadBalancerAclId"] = d.Get("load_balancer_acl_id")
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateLoadBalancerAclEntry(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if isSetId {
				_, err = s.ReadLoadBalancerAcl(d, (*(call.param))["LoadBalancerAclId"].(string))
				if err != nil {
					return err
				}
				d.SetId((*(call.param))["LoadBalancerAclId"].(string) + ":" + strconv.Itoa(d.Get("rule_number").(int)) + ":" + d.Get("cidr_block").(string))
			}
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateLoadBalancerAclAssociateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	return s.CreateLoadBalancerAclAssociateCommonCall(req, true)
}

func (s *SlbService) CreateLoadBalancerAclAssociateWithListenerCall(d *schema.ResourceData, r *schema.Resource, loadBalancerAclId string) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"LoadBalancerAclId": loadBalancerAclId,
	}
	return s.CreateLoadBalancerAclAssociateCommonCall(req, false)
}

func (s *SlbService) CreateLoadBalancerAclAssociateCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "AssociateLoadBalancerAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			(*(call.param))["ListenerId"] = d.Get("listener_id")
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AssociateLoadBalancerAcl(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if isSetId {
				d.SetId(d.Get("listener_id").(string) + ":" + d.Get("load_balancer_acl_id").(string))
			}
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateLoadBalancerAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateLoadBalancerAclCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	entries, err := s.CreateLoadBalancerAclEntryWithAclCall(d, r)
	if err != nil {
		return err
	}
	for _, entryCall := range entries {
		callbacks = append(callbacks, entryCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, false)
}

func (s *SlbService) CreateLoadBalancerAclEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateLoadBalancerAclEntryCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) CreateLoadBalancerAclAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateLoadBalancerAclAssociateCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) ModifyLoadBalancerAclCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"load_balancer_acl_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["LoadBalancerAclId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyLoadBalancerAcl",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyLoadBalancerAcl(call.param)
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

func (s *SlbService) ModifyLoadBalancerAclEntryCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["LoadBalancerAclEntryId"] = d.Get("load_balancer_acl_entry_id")
		return s.ModifyLoadBalancerAclEntryCommonCall(req)
	}
	return callback, err
}

func (s *SlbService) ModifyLoadBalancerAclEntryWithAclCall(d *schema.ResourceData, r *schema.Resource) (callbacks []ApiCall, err error) {
	if d.HasChange("load_balancer_acl_entry_set") {
		o, n := d.GetChange("load_balancer_acl_entry_set")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		//check change is valid
		if len(schema.NewSet(loadBalancerAclEntryNumberHash, ns.Difference(os).List()).List()) != len(ns.Difference(os).List()) {
			return callbacks, fmt.Errorf("RuleNumber must unique ")
		}
		if len(schema.NewSet(loadBalancerAclEntryCidrHash, ns.Difference(os).List()).List()) != len(ns.Difference(os).List()) {
			return callbacks, fmt.Errorf("CidrBlock must unique ")
		}
		//generate new hashcode without can modify field
		mayAdd := schema.NewSet(loadBalancerAclEntrySimpleHash, ns.Difference(os).List())
		mayRemove := schema.NewSet(loadBalancerAclEntrySimpleHash, os.Difference(ns).List())
		addCache := make(map[int]interface{})
		for _, entry := range mayAdd.List() {
			index := loadBalancerAclEntrySimpleHash(entry)
			addCache[index] = entry
		}
		//compare hashcode  can modify field
		//need add entries
		add := mayAdd.Difference(mayRemove)
		//need remove entries
		remove := mayRemove.Difference(mayAdd)
		//need modify entries
		modify := mayRemove.Difference(remove)
		//process modify
		if len(modify.List()) > 0 {
			for _, entry := range modify.List() {
				var (
					callback ApiCall
				)
				index := loadBalancerAclEntrySimpleHash(entry)
				req := make(map[string]interface{})
				req["RuleNumber"] = addCache[index].(map[string]interface{})["rule_number"]
				req["RuleAction"] = addCache[index].(map[string]interface{})["rule_action"]
				req["LoadBalancerAclEntryId"] = entry.(map[string]interface{})["load_balancer_acl_entry_id"]
				logger.Debug(logger.ReqFormat, "DemoModify", req)
				callback, err = s.ModifyLoadBalancerAclEntryCommonCall(req)
				if err != nil {
					return callbacks, err
				}
				callbacks = append(callbacks, callback)
			}
		}
		//process remove
		if len(remove.List()) > 0 {
			for _, entry := range remove.List() {
				var (
					callback ApiCall
				)
				callback, err = s.RemoveLoadBalancerAclEntryCommonCall(d.Id(), entry.(map[string]interface{})["load_balancer_acl_entry_id"].(string))
				if err != nil {
					return callbacks, err
				}
				callbacks = append(callbacks, callback)
			}
		}
		//process add
		if len(add.List()) > 0 {
			for _, entry := range add.List() {
				var (
					req      map[string]interface{}
					callback ApiCall
				)
				index := loadBalancerAclEntryHash(entry)
				transform := make(map[string]SdkReqTransform)
				for k := range entry.(map[string]interface{}) {
					key := "load_balancer_acl_entry_set." + strconv.Itoa(index) + "." + k
					transform[key] = SdkReqTransform{mapping: Downline2Hump(k)}
				}
				req, err = SdkRequestAutoMapping(d, r, false, transform, nil)
				if err != nil {
					return callbacks, err
				}
				callback, err = s.CreateLoadBalancerAclEntryCommonCall(req, false)
				if err != nil {
					return callbacks, err
				}
				callbacks = append(callbacks, callback)
			}
		}

	}
	return callbacks, err
}

func (s *SlbService) ModifyLoadBalancerAclEntryCommonCall(req map[string]interface{}) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "ModifyLoadBalancerAclEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyLoadBalancerAclEntry(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *SlbService) ModifyLoadBalancerAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyLoadBalancerAclCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	entries, err := s.ModifyLoadBalancerAclEntryWithAclCall(d, r)
	if err != nil {
		return err
	}
	for _, entryCall := range entries {
		callbacks = append(callbacks, entryCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) ModifyLoadBalancerAclEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyLoadBalancerAclEntryCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) RemoveLoadBalancerAclCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"LoadBalancerAclId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteLoadBalancerAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteLoadBalancerAcl(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLoadBalancerAcl(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb acl when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveLoadBalancerAclEntryCommonCall(aclId string, entryId string) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"LoadBalancerAclId":      aclId,
		"LoadBalancerAclEntryId": entryId,
	}
	callback = ApiCall{
		param:  &req,
		action: "DeleteLoadBalancerAclEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteLoadBalancerAclEntry(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				data, callErr := s.ReadLoadBalancerAcl(d, aclId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb rulr entry when delete %q, %s", d.Id(), callErr))
					}
				}
				if len(data["LoadBalancerAclEntrySet"].([]interface{})) == 0 {
					return nil
				} else {
					found := false
					for _, item := range data["LoadBalancerAclEntrySet"].([]interface{}) {
						if item.(map[string]interface{})["LoadBalancerAclEntryId"] == entryId {
							found = true
							break
						}
					}
					if !found {
						return nil
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

func (s *SlbService) RemoveLoadBalancerAclAssociateCommonCall(listenerId string, loadBalancerAclId string) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"ListenerId": listenerId,
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DisassociateLoadBalancerAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DisassociateLoadBalancerAcl(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLoadBalancerAclAssociate(listenerId, loadBalancerAclId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb acl when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveLoadBalancerAcl(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLoadBalancerAclCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) RemoveLoadBalancerAclEntry(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLoadBalancerAclEntryCommonCall(d.Get("load_balancer_acl_id").(string), d.Get("load_balancer_acl_entry_id").(string))
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) RemoveLoadBalancerAclAssociate(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLoadBalancerAclAssociateCommonCall(d.Get("listener_id").(string), d.Get("load_balancer_acl_id").(string))
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// start RealServer

func (s *SlbService) ReadRealServers(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeInstancesWithListener"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeInstancesWithListener(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeInstancesWithListener(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("RealServerSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadRealServer(d *schema.ResourceData, registerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if registerId == "" {
		registerId = d.Id()
	}
	req := map[string]interface{}{
		"RegisterId.1": registerId,
	}
	results, err = s.ReadRealServers(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Real Server %s not exist ", registerId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetRealServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadRealServer(d, "")
	if err != nil {
		return err
	}
	listenerId := data["ListenerId"].(string)
	listener, err := s.ReadListener(nil, listenerId)
	if err != nil {
		return err
	}
	data["ListenerMethod"] = listener["Method"]
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetRealServers(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "RegisterId",
			Type:    TransformWithN,
		},
		"listener_id": {
			mapping: "listener-id",
			Type:    TransformWithFilter,
		},
		"real_server_ip": {
			mapping: "real-server-ip",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadRealServers(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "RegisterId",
		targetField: "servers",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateRealServerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	//获取一次监听器信息 判定method类型
	vip, err := s.ReadListener(nil, req["ListenerId"].(string))
	if err != nil {
		return callback, err
	}
	if vip["Method"] != "MasterSlave" {
		delete(req, "MasterSlaveType")
	}
	callback = ApiCall{
		param:  &req,
		action: "RegisterInstancesWithListener",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.RegisterInstancesWithListener(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("RegisterId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateRealServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateRealServerCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) ModifyRealServerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["RegisterId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyInstancesWithListener",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyInstancesWithListener(call.param)
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

func (s *SlbService) ModifyRealServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyRealServerCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) RemoveRealServerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"RegisterId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeregisterInstancesFromListener",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeregisterInstancesFromListener(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadRealServer(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading real server when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveRealServer(d *schema.ResourceData) (err error) {
	call, err := s.RemoveRealServerCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// start BackendServerGroup

func (s *SlbService) ReadBackendServerGroups(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeBackendServerGroups"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeBackendServerGroups(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeBackendServerGroups(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("BackendServerGroupSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadBackendServerGroup(d *schema.ResourceData, backendServerGroupId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if backendServerGroupId == "" {
		backendServerGroupId = d.Id()
	}
	req := map[string]interface{}{
		"BackendServerGroupId.1": backendServerGroupId,
	}
	results, err = s.ReadBackendServerGroups(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("BackendServerGroup %s not exist ", backendServerGroupId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetBackendServerGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadBackendServerGroup(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetBackendServerGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "BackendServerGroupId",
			Type:    TransformWithN,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"backend_server_group_type": {
			mapping: "backend-server-group-type",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadBackendServerGroups(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "BackendServerGroupId",
		targetField: "backend_server_groups",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateBackendServerGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	//特殊处理下"HealthCheck."
	for k, v := range req {
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	if _, ok := req["UrlPath"]; !ok && req["BackendServerGroupType"] == "Mirror" {
		return callback, fmt.Errorf("BackendServerGroupType is Mirror must set HealthCheck")
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateBackendServerGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateBackendServerGroup(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("BackendServerGroup.BackendServerGroupId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateBackendServerGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateBackendServerGroupCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) ModifyBackendServerGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"backend_server_group_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["BackendServerGroupId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyBackendServerGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyBackendServerGroup(call.param)
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

func (s *SlbService) ModifyBackendServerGroupHealthCheckCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	//特殊处理下"HealthCheck."
	for k, v := range req {
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	if len(req) > 0 {
		req["BackendServerGroupId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyBackendServerGroupHealthCheck",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyBackendServerGroupHealthCheck(call.param)
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

func (s *SlbService) ModifyBackendServerGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyBackendServerGroupCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	healthCheckCall, err := s.ModifyBackendServerGroupHealthCheckCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, healthCheckCall)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) RemoveBackendServerGroupCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"BackendServerGroupId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteBackendServerGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteBackendServerGroup(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadBackendServerGroup(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading backend server group when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveBackendServerGroup(d *schema.ResourceData) (err error) {
	call, err := s.RemoveBackendServerGroupCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// start backend server group server

func (s *SlbService) ReadBackendServers(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeBackendServers"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeBackendServers(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeBackendServers(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("BackendServerSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadBackendServer(d *schema.ResourceData, registerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if registerId == "" {
		registerId = d.Id()
	}
	req := map[string]interface{}{
		"RegisterId.1": registerId,
	}
	results, err = s.ReadBackendServers(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("BackendServer %s not exist ", registerId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetBackendServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadBackendServer(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetBackendServers(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "RegisterId",
			Type:    TransformWithN,
		},
		"backend_server_group_id": {
			mapping: "backend-server-group-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadBackendServers(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "RegisterId",
		targetField: "register_backend_servers",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateBackendServerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "RegisterBackendServer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.RegisterBackendServer(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("BackendServer.RegisterId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateBackendServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateBackendServerCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) ModifyBackendServerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["RegisterId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyBackendServer",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyBackendServer(call.param)
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

func (s *SlbService) ModifyBackendServer(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyBackendServerCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *SlbService) RemoveBackendServerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"RegisterId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeregisterBackendServer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeregisterBackendServer(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadBackendServer(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading backend server when delete %q, %s", d.Id(), callErr))
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

func (s *SlbService) RemoveBackendServer(d *schema.ResourceData) (err error) {
	call, err := s.RemoveBackendServerCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
