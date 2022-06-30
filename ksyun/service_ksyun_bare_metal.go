package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

type BareMetalService struct {
	client *KsyunClient
}

func (s *BareMetalService) ReadBareMetals(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.epcconn
		action := "DescribeEpcs"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeEpcs(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeEpcs(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("HostSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *BareMetalService) ReadBareMetal(d *schema.ResourceData, hostId string, allProject bool) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if hostId == "" {
		hostId = d.Id()
	}
	req := map[string]interface{}{
		"HostId.1": hostId,
	}
	if allProject {
		err = addProjectInfoAll(d, &req, s.client)
		if err != nil {
			return data, err
		}
	} else {
		err = addProjectInfo(d, &req, s.client)
		if err != nil {
			return data, err
		}
	}

	results, err = s.ReadBareMetals(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("BareMetal %s not exist ", hostId)
	}
	return data, err
}

func (s *BareMetalService) ReadAndSetBareMetal(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadBareMetal(d, "", false)
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading address %q, %s", d.Id(), callErr))
			}
		} else {
			if vifs, ok := data["NetworkInterfaceAttributeSet"]; ok {
				for _, vif := range vifs.([]interface{}) {
					networkInterfaceType := vif.(map[string]interface{})["NetworkInterfaceType"]
					for k, v := range vif.(map[string]interface{}) {
						if networkInterfaceType == "primary" {
							data[k] = v
						} else {
							data["Extension"+k] = v
						}
					}
				}
				delete(data, "NetworkInterfaceAttributeSet")
			}
			extra := map[string]SdkResponseMapping{
				"RaidTemplateId": {
					Field: "raid_id",
				},
				"DNS1": {
					Field: "dns1",
				},
				"DNS2": {
					Field: "dns2",
				},
				"SecurityGroupSet": {
					Field: "security_group_ids",
					FieldRespFunc: func(i interface{}) interface{} {
						var value []interface{}
						for _, v := range i.([]interface{}) {
							value = append(value, v.(map[string]interface{})["SecurityGroupId"])
						}
						return value
					},
				},
				"ExtensionDNS1": {
					Field: "extension_dns1",
				},
				"ExtensionDNS2": {
					Field: "extension_dns1",
				},
				"ExtensionSecurityGroupSet": {
					Field: "extension_security_group_ids",
					FieldRespFunc: func(i interface{}) interface{} {
						var value []interface{}
						for _, v := range i.([]interface{}) {
							value = append(value, v.(map[string]interface{})["SecurityGroupId"])
						}
						return value
					},
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (s *BareMetalService) ReadAndSetBareMetals(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "HostId",
			Type:    TransformWithN,
		},
		"project_id": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"host_name": {
			mapping: "host-name",
			Type:    TransformWithFilter,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"subnet_id": {
			mapping: "subnet-id",
			Type:    TransformWithFilter,
		},
		"cabinet_id": {
			mapping: "cabinet-id",
			Type:    TransformWithFilter,
		},
		"host_type": {
			mapping: "host-type",
			Type:    TransformWithFilter,
		},
		"epc_host_status": {
			mapping: "epc-host-status",
			Type:    TransformWithFilter,
		},
		"os_name": {
			mapping: "os-name",
			Type:    TransformWithFilter,
		},
		"product_type": {
			mapping: "product-type",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadBareMetals(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "HostId",
		targetField: "bare_metals",
		extra: map[string]SdkResponseMapping{
			"DNS1": {
				Field: "dns1",
			},
			"DNS2": {
				Field: "dns2",
			},
		},
	})
}

func (s *BareMetalService) BareMetalStateRefreshFunc(d *schema.ResourceData, hostId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := s.ReadBareMetal(d, hostId, true)
		if err != nil {
			return nil, "", err
		}

		status, err := getSdkValue("HostStatus", data)
		if err != nil {
			return nil, "", err
		}

		for _, v := range failStates {
			if v == status.(string) {
				return nil, "", fmt.Errorf("baremetal status  error, status:%v", status)
			}
		}
		return data, status.(string), nil
	}
}

func (s *BareMetalService) CheckBareMetalState(d *schema.ResourceData, hostId string, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    s.BareMetalStateRefreshFunc(d, hostId, []string{"failed"}),
		Timeout:    timeout,
		Delay:      1 * time.Minute,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *BareMetalService) CreateBareMetalCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"security_group_ids": {
			mapping: "SecurityGroupId",
			Type:    TransformWithN,
		},
		"extension_security_group_ids": {
			mapping: "ExtensionSecurityGroupId",
			Type:    TransformWithN,
		},
		"dns1": {
			mapping: "DNS1",
		},
		"dns2": {
			mapping: "DNS2",
		},
		"extension_dns1": {
			mapping: "ExtensionDNS1",
		},
		"extension_dns2": {
			mapping: "ExtensionDNS2",
		},
		"force_re_install": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, resource, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	req["ChargeType"] = "Daily"

	callback = ApiCall{
		param:  &req,
		action: "CreateEpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.epcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateEpc(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var (
				hostId interface{}
			)
			if resp != nil {
				hostId, err = getSdkValue("Host.HostId", *resp)
				if err != nil {
					return err
				}
				d.SetId(hostId.(string))
			}
			err = s.CheckBareMetalState(d, "", []string{"Running"}, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return err
			}
			return s.ReadAndSetBareMetal(d, resource)
		},
	}
	return callback, err
}

func (s *BareMetalService) CreateBareMetal(d *schema.ResourceData, resource *schema.Resource) (err error) {
	var (
		callbacks []ApiCall
	)
	createCall, err := s.CreateBareMetalCall(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, createCall)
	// dryRun
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *BareMetalService) ModifyBareMetalInfoCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"host_name": {},
	}
	req, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["HostId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyEpc",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.epcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyEpc(call.param)
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

func (s *BareMetalService) ModifyBareMetalNetworkCall(d *schema.ResourceData, resource *schema.Resource, isPrimary bool) (callback ApiCall, err error) {
	var transform map[string]SdkReqTransform
	if isPrimary {
		transform = map[string]SdkReqTransform{
			"subnet_id":          {},
			"private_ip_address": {mapping: "IpAddress"},
		}
	} else {
		transform = map[string]SdkReqTransform{
			"extension_subnet_id":          {},
			"extension_private_ip_address": {mapping: "IpAddress"},
		}
	}
	req, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["HostId"] = d.Id()
		if isPrimary {
			req["NetworkInterfaceId"] = d.Get("network_interface_id")
			if _, ok := req["SubnetId"]; !ok {
				req["SubnetId"] = d.Get("subnet_id")
			}
		} else {
			req["NetworkInterfaceId"] = d.Get("extension_network_interface_id")
			if _, ok := req["SubnetId"]; !ok {
				req["SubnetId"] = d.Get("extension_subnet_id")
			}
		}
		callback = ApiCall{
			param:  &req,
			action: "ModifyNetworkInterfaceAttribute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.epcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyNetworkInterfaceAttribute(call.param)
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

func (s *BareMetalService) ModifyBareMetalDnsCall(d *schema.ResourceData, resource *schema.Resource, isPrimary bool) (callback ApiCall, err error) {
	var transform map[string]SdkReqTransform
	if isPrimary {
		transform = map[string]SdkReqTransform{
			"dns1": {mapping: "DNS1"},
			"dns2": {mapping: "DNS2"},
		}
	} else {
		transform = map[string]SdkReqTransform{
			"extension_dns1": {mapping: "DNS1"},
			"extension_dns2": {mapping: "DNS2"},
		}
	}
	req, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["HostId"] = d.Id()
		if isPrimary {
			req["NetworkInterfaceId"] = d.Get("network_interface_id")
		} else {
			req["NetworkInterfaceId"] = d.Get("extension_network_interface_id")
		}
		callback = ApiCall{
			param:  &req,
			action: "ModifyDns",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.epcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyDns(call.param)
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

func (s *BareMetalService) ModifyBareMetalSecurityGroupCall(d *schema.ResourceData, resource *schema.Resource, isPrimary bool) (callback ApiCall, err error) {
	var transform map[string]SdkReqTransform
	if isPrimary {
		transform = map[string]SdkReqTransform{
			"security_group_ids": {
				mapping: "SecurityGroupId",
				Type:    TransformWithN,
			},
		}
	} else {
		transform = map[string]SdkReqTransform{
			"extension_security_group_ids": {
				mapping: "SecurityGroupId",
				Type:    TransformWithN,
			},
		}
	}
	req, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["HostId"] = d.Id()
		if isPrimary {
			req["NetworkInterfaceId"] = d.Get("network_interface_id")
		} else {
			req["NetworkInterfaceId"] = d.Get("extension_network_interface_id")
		}
		callback = ApiCall{
			param:  &req,
			action: "ModifySecurityGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.epcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifySecurityGroup(call.param)
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

func (s *BareMetalService) ReinstallBareMetalCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	if d.Get("host_type") == "COLO" {
		return callback, err
	}
	transform := map[string]SdkReqTransform{
		"host_name":                    {Ignore: true},
		"dns1":                         {Ignore: true},
		"dns2":                         {Ignore: true},
		"private_ip_address":           {Ignore: true},
		"security_group_ids":           {Ignore: true},
		"subnet_id":                    {Ignore: true},
		"extension_dns1":               {Ignore: true},
		"extension_dns2":               {Ignore: true},
		"extension_security_group_ids": {Ignore: true},
		"extension_private_ip_address": {Ignore: true},
		"extension_subnet_id":          {Ignore: true},
		"server_ip":                    {Ignore: true},
		"path":                         {Ignore: true},
		"force_re_install":             {Ignore: true},
	}
	if d.HasChange("force_re_install") && d.Get("force_re_install").(bool) {
		transform["image_id"] = SdkReqTransform{
			forceUpdateParam: true,
		}
	}
	req, err := SdkRequestAutoMapping(d, resource, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		_transform := make(map[string]SdkReqTransform)
		for k, v := range resource.Schema {
			if _, ok := transform[k]; !ok && !v.ForceNew {
				_transform[k] = SdkReqTransform{
					forceUpdateParam: true,
				}
			}
		}
		req, err = SdkRequestAutoMapping(d, resource, true, _transform, nil)
		if err != nil {
			return callback, err
		}
		for k, v := range req {
			if v == "" {
				delete(req, k)
			}
		}
		req["HostId"] = d.Id()
		if _, ok := req["ImageId"]; !ok {
			req["ImageId"] = d.Get("image_id")
		}
		callback = ApiCall{
			param:  &req,
			action: "ReinstallEpc",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.epcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ReinstallEpc(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = s.CheckBareMetalState(d, "", []string{"Running"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				return err
			},
			callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
				if d.HasChange("force_re_install") {
					_ = d.Set("force_re_install", !d.Get("force_re_install").(bool))
				}
				return baseErr
			},
		}
	}
	return callback, err
}

func (s *BareMetalService) ReinstallCustomerBareMetalCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	if d.Get("host_type") != "COLO" {
		return callback, err
	}
	transform := map[string]SdkReqTransform{
		"server_ip": {},
		"path":      {},
	}
	if d.HasChange("force_re_install") && d.Get("force_re_install").(bool) {
		transform["server_ip"] = SdkReqTransform{
			forceUpdateParam: true,
		}
		transform["path"] = SdkReqTransform{
			forceUpdateParam: true,
		}
	}
	req, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		if len(req) < 2 {
			if d.HasChange("force_re_install") {
				_ = d.Set("force_re_install", !d.Get("force_re_install").(bool))
			}
			return callback, fmt.Errorf("server_ip and path must both set")
		}
		if req["ServerIp"] == "" || req["Path"] == "" {
			if d.HasChange("force_re_install") {
				_ = d.Set("force_re_install", !d.Get("force_re_install").(bool))
			}
			return callback, fmt.Errorf("server_ip and path must both set")
		}
		req["HostId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ReinstallCustomerEpc",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.epcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ReinstallCustomerEpc(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = s.CheckBareMetalState(d, "", []string{"Running"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				return err
			},
			callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
				if d.HasChange("force_re_install") {
					_ = d.Set("force_re_install", !d.Get("force_re_install").(bool))
				}
				return baseErr
			},
		}
	}
	return callback, err
}

func (s *BareMetalService) ModifyBareMetalProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
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

func (s *BareMetalService) ModifyBareMetal(d *schema.ResourceData, resource *schema.Resource) (err error) {
	var (
		callbacks []ApiCall
	)
	reinstallCall, err := s.ReinstallBareMetalCall(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, reinstallCall)

	customerReinstallCall, err := s.ReinstallCustomerBareMetalCall(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, customerReinstallCall)

	infoCall, err := s.ModifyBareMetalInfoCall(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, infoCall)

	sgPrimaryCall, err := s.ModifyBareMetalSecurityGroupCall(d, resource, true)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, sgPrimaryCall)

	sgCall, err := s.ModifyBareMetalSecurityGroupCall(d, resource, false)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, sgCall)

	dnsPrimaryCall, err := s.ModifyBareMetalDnsCall(d, resource, true)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, dnsPrimaryCall)

	dnsCall, err := s.ModifyBareMetalDnsCall(d, resource, false)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, dnsCall)

	networkPrimaryCall, err := s.ModifyBareMetalNetworkCall(d, resource, true)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, networkPrimaryCall)

	networkCall, err := s.ModifyBareMetalNetworkCall(d, resource, false)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, networkCall)

	projectCall, err := s.ModifyBareMetalProjectCall(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, projectCall)
	// dryRun
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *BareMetalService) RemoveBareMetalCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"HostId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteEpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.epcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteEpc(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadBareMetal(d, "", false)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading bare metal when delete %q, %s", d.Id(), callErr))
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

func (s *BareMetalService) RemoveBareMetal(d *schema.ResourceData) (err error) {
	call, err := s.RemoveBareMetalCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *BareMetalService) ReadImages(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.epcconn
		action := "DescribeImages"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeImages(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeImages(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("ImageSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *BareMetalService) ReadAndSetImages(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "ImageId",
			Type:    TransformWithN,
		},
		"image_type": {
			mapping: "image-type",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadImages(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "ImageId",
		targetField: "images",
		nameField:   "ImageName",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *BareMetalService) ReadRaidAttributes(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.epcconn
		action := "DescribeEpcRaidAttributes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeEpcRaidAttributes(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeEpcRaidAttributes(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("EpcRaidAttributeSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *BareMetalService) ReadAndSetRaidAttributes(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"host_type": {
			mapping: "host-type",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadRaidAttributes(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "RaidId",
		targetField: "raid_attributes",
		nameField:   "TemplateName",
		extra:       map[string]SdkResponseMapping{},
	})
}
