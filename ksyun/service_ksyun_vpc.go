package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

type VpcService struct {
	client *KsyunClient
}

func (s *VpcService) ReadNetworkInterfaces(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp                    *map[string]interface{}
		networkInterfaceResults interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeNetworkInterfaces"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeNetworkInterfaces(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeNetworkInterfaces(&condition)
		if err != nil {
			return data, err
		}
	}

	networkInterfaceResults, err = getSdkValue("NetworkInterfaceSet", *resp)
	if err != nil {
		return data, err
	}
	data = networkInterfaceResults.([]interface{})
	return data, err
}

func (s *VpcService) ReadNetworkInterface(d *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		networkInterfaceResults []interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"NetworkInterfaceId.1": instanceId,
	}
	networkInterfaceResults, err = s.ReadNetworkInterfaces(req)
	if err != nil {
		return data, err
	}
	for _, v := range networkInterfaceResults {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("NetworkInterface %s not exist ", instanceId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetNetworkInterfaces(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "NetworkInterfaceId",
			Type:    TransformWithN,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"subnet_id": {
			mapping: "subnet-id",
			Type:    TransformWithFilter,
		},
		"securitygroup_id": {
			mapping: "securitygroup-id",
			Type:    TransformWithFilter,
		},
		"instance_type": {
			mapping: "instance-type",
			Type:    TransformWithFilter,
		},
		"instance_id": {
			mapping: "instance-id",
			Type:    TransformWithFilter,
		},
		"private_ip_address": {
			mapping: "private-ip-address",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadNetworkInterfaces(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "NetworkInterfaceName",
		idFiled:     "NetworkInterfaceId",
		targetField: "network_interfaces",
		extra: map[string]SdkResponseMapping{
			"NetworkInterfaceName": {
				Field:    "name",
				KeepAuto: true,
			},
			"NetworkInterfaceId": {
				Field:    "id",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) RemoveNetworkInterface(d *schema.ResourceData) (err error) {
	call, err := s.RemoveNetworkInterfaceCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveNetworkInterfaceCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"NetworkInterfaceId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteNetworkInterface(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadNetworkInterface(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading network interface when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) CreateNetworkInterfaceCall(createReq *map[string]interface{}) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  createReq,
		action: "CreateNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateNetworkInterface(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var (
				instanceId interface{}
			)
			if resp != nil {
				instanceId, err = getSdkValue("NetworkInterfaceId", *resp)
				if err != nil {
					return err
				}
				d.SetId(instanceId.(string))
			}
			return err
		},
	}
	return callback, err
}

func (s *VpcService) ModifyNetworkInterfaceCall(updateReq *map[string]interface{}) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  updateReq,
		action: "ModifyNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyNetworkInterface(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *VpcService) ReadVpcs(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeVpcs"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeVpcs(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeVpcs(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("VpcSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadVpc(d *schema.ResourceData, vpcId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if vpcId == "" {
		vpcId = d.Id()
	}
	req := map[string]interface{}{
		"VpcId.1": vpcId,
	}
	results, err = s.ReadVpcs(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Vpc %s not exist ", vpcId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetVpc(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadVpc(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *VpcService) ReadAndSetVpcs(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "VpcId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadVpcs(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "VpcName",
		idFiled:     "VpcId",
		targetField: "vpcs",
		extra: map[string]SdkResponseMapping{
			"VpcName": {
				Field:    "name",
				KeepAuto: true,
			},
			"VpcId": {
				Field:    "id",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) CreateVpcCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateVpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateVpc(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Vpc.VpcId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateVpc(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateVpcCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifyVpcCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["VpcId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyVpc",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyVpc(call.param)
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

func (s *VpcService) ModifyVpc(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyVpcCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveVpcCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"VpcId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteVpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteVpc(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadVpc(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading vpc when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveVpc(d *schema.ResourceData) (err error) {
	call, err := s.RemoveVpcCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadSubnets(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeSubnets"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeSubnets(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeSubnets(&condition)
		if err != nil {
			return data, err
		}
	}
	logger.Debug(logger.ReqFormat, action, *resp)
	results, err = getSdkValue("SubnetSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadSubnet(d *schema.ResourceData, subnetId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if subnetId == "" {
		subnetId = d.Id()
	}
	req := map[string]interface{}{
		"SubnetId.1": subnetId,
	}
	results, err = s.ReadSubnets(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})

		// 手动补AvailabilityZoneName字段，解决和配置字段不能对应引发的问题
		// 例如: import的时候读不到这个值，导致认为这个值是空，plan会发现值变化
		if v, ok := data["AvailabilityZoneName"]; ok {
			data["AvailabilityZone"] = v
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Subnet %s not exist ", subnetId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetSubnet(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadSubnet(d, "")
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"AvailableIpNumber": {Field: "available_ip_number"},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return err
}

func (s *VpcService) ReadAndSetSubnets(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "SubnetId",
			Type:    TransformWithN,
		},
		"vpc_ids": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"subnet_types": {
			mapping: "subnet-type",
			Type:    TransformWithFilter,
		},
		"nat_ids": {
			mapping: "nat-id",
			Type:    TransformWithFilter,
		},
		"network_acl_ids": {
			mapping: "network-acl-id",
			Type:    TransformWithFilter,
		},
		"availability_zone_names": {
			mapping: "availability-zone-name",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	logger.Debug(logger.ReqFormat, "Demo", req)
	if err != nil {
		return err
	}
	data, err := s.ReadSubnets(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "SubnetName",
		idFiled:     "SubnetId",
		targetField: "subnets",
		extra: map[string]SdkResponseMapping{
			"SubnetName": {
				Field: "name",
			},
			"SubnetId": {
				Field:    "id",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) SubnetAutoMatch(req *map[string]interface{}) {
	for k, v := range *req {
		if k == "SubnetType" && v.(string) != "Reserve" {
			gw, start, end := getCidrIpRange((*req)["CidrBlock"].(string))
			if _, ok := (*req)["GatewayIp"]; !ok {
				(*req)["GatewayIp"] = gw
			}
			if _, ok := (*req)["DhcpIpFrom"]; !ok {
				(*req)["DhcpIpFrom"] = start
			}
			if _, ok := (*req)["DhcpIpTo"]; !ok {
				(*req)["DhcpIpTo"] = end
			}
		}
	}
}

func (s *VpcService) CreateSubnetCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	s.SubnetAutoMatch(&req)
	if req["SubnetType"] != "Reserve" {
		if _, ok := req["GatewayIp"]; !ok {
			return callback, fmt.Errorf("GatewayIp must set when SubnetType is not Reserve ")
		}
		if _, ok := req["DhcpIpFrom"]; !ok {
			return callback, fmt.Errorf("DhcpIpFrom must set when SubnetType is not Reserve ")
		}
		if _, ok := req["DhcpIpTo"]; !ok {
			return callback, fmt.Errorf("DhcpIpTo must set when SubnetType is not Reserve ")
		}
	}

	callback = ApiCall{
		param:  &req,
		action: "CreateSubnet",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateSubnet(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Subnet.SubnetId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateSubnet(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateSubnetCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifySubnetCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["SubnetId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifySubnet",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifySubnet(call.param)
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

func (s *VpcService) ModifySubnet(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifySubnetCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveSubnetCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"SubnetId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteSubnet",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteSubnet(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadSubnet(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading subnet when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveSubnet(d *schema.ResourceData) (err error) {
	call, err := s.RemoveSubnetCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadRoutes(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeRoutes"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeRoutes(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeRoutes(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("RouteSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadRoute(d *schema.ResourceData, subnetId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if subnetId == "" {
		subnetId = d.Id()
	}
	req := map[string]interface{}{
		"RouteId.1": subnetId,
	}
	results, err = s.ReadRoutes(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Route %s not exist ", subnetId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetRoute(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadRoute(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *VpcService) ReadAndSetRoutes(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "RouteId",
			Type:    TransformWithN,
		},
		"vpc_ids": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
		"instance_ids": {
			mapping: "instance-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadRoutes(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "RouteId",
		targetField: "routes",
		extra: map[string]SdkResponseMapping{
			"RouteId": {
				Field:    "id",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) CreateRouteCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	switch req["RouteType"] {
	case "Tunnel":
		if _, ok := req["TunnelId"]; !ok {
			return callback, fmt.Errorf("TunnelId must set when RouteType is Tunnel")
		}
		break
	case "Host":
		if _, ok := req["InstanceId"]; !ok {
			return callback, fmt.Errorf("InstanceId must set when RouteType is Host")
		}
		break
	case "Peering":
		if _, ok := req["VpcPeeringConnectionId"]; !ok {
			return callback, fmt.Errorf("VpcPeeringConnectionId must set when RouteType is Peering")
		}
		break
	case "DirectConnect":
		if _, ok := req["DirectConnectGatewayId"]; !ok {
			return callback, fmt.Errorf("DirectConnectGatewayId must set when RouteType is DirectConnect")
		}
		break
	case "Vpn":
		if _, ok := req["VpnTunnelId"]; !ok {
			return callback, fmt.Errorf("VpnTunnelId must set when RouteType is Vpn")
		}
		break
	default:
		break
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateRoute",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateRoute(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("RouteId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateRoute(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateRouteCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveRouteCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"RouteId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteRoute",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteRoute(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadRoute(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading route when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveRoute(d *schema.ResourceData) (err error) {
	call, err := s.RemoveRouteCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadNats(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeNats"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeNats(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeNats(&condition)
		if err != nil {
			return data, err
		}
	}
	logger.Debug(logger.ReqFormat, action, *resp)
	results, err = getSdkValue("NatSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadNat(d *schema.ResourceData, natId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if natId == "" {
		natId = d.Id()
	}
	req := map[string]interface{}{
		"NatId.1": natId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}
	results, err = s.ReadNats(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Nat %s not exist ", natId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetNats(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "NatId",
			Type:    TransformWithN,
		},
		"project_ids": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"vpc_ids": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadNats(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "NatName",
		idFiled:     "NatId",
		targetField: "nats",
		extra: map[string]SdkResponseMapping{
			"NatId": {
				Field:    "id",
				KeepAuto: false,
			},
		},
	})
}

func (s *VpcService) ReadAndSetNat(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadNat(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading nat %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, chargeExtraForVpc(data))
			return nil
		}
	})
}

func (s *VpcService) CreateNatCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateNat",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateNat(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("NatId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateNat(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateNatCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifyNatCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{false})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["NatId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyNat",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyNat(call.param)
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

func (s *VpcService) modifyNatProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
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

func (s *VpcService) ModifyNat(d *schema.ResourceData, r *schema.Resource) (err error) {
	projectCall, err := s.modifyNatProjectCall(d, r)
	if err != nil {
		return err
	}
	call, err := s.ModifyNatCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{projectCall, call}, d, s.client, true)
}

func (s *VpcService) RemoveNatCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"NatId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteNat",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteNat(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadNat(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading nat when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveNat(d *schema.ResourceData) (err error) {
	call, err := s.RemoveNatCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadNatAssociate(d *schema.ResourceData, natId string, subnetId string) (data map[string]interface{}, err error) {
	data, err = s.ReadNat(d, natId)
	if items, ok := data["AssociateNatSet"]; ok {
		if len(items.([]interface{})) == 0 {
			return data, fmt.Errorf("Subnet %s not exist in Nat %s ", subnetId, natId)
		}
		found := false
		for _, item := range items.([]interface{}) {
			if item.(map[string]interface{})["SubnetId"] == subnetId {
				found = true
				break
			}
		}
		if !found {
			return data, fmt.Errorf("Subnet %s not exist in Nat %s ", subnetId, natId)
		}
	} else {
		return data, fmt.Errorf("Subnet %s not exist in Nat %s ", subnetId, natId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetNatAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadNatAssociate(d, d.Get("nat_id").(string), d.Get("subnet_id").(string))
	if err != nil {
		return err
	}
	data["SubnetId"] = d.Get("subnet_id")
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *VpcService) CreateNatAssociateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AssociateNat",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AssociateNat(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("nat_id").(string) + ":" + d.Get("subnet_id").(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateNatAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateNatAssociateCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveNatAssociateCall(d *schema.ResourceData, natId string, subnetId string) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"NatId":    natId,
		"SubnetId": subnetId,
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DisassociateNat",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DisassociateNat(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadNatAssociate(d, natId, subnetId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading nat associate when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveNatAssociate(d *schema.ResourceData) (err error) {
	call, err := s.RemoveNatAssociateCall(d, d.Get("nat_id").(string), d.Get("subnet_id").(string))
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadNetworkAcls(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeNetworkAcls"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeNetworkAcls(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeNetworkAcls(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("NetworkAclSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadNetworkAcl(d *schema.ResourceData, networkAclId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if networkAclId == "" {
		networkAclId = d.Id()
	}
	req := map[string]interface{}{
		"NetworkAclId.1": networkAclId,
	}
	results, err = s.ReadNetworkAcls(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Acl %s not exist ", networkAclId)
	}
	return data, err
}

func (s *VpcService) ReadNetworkAclEntry(d *schema.ResourceData, networkAclId string) (data map[string]interface{}, err error) {
	acl, err := s.ReadNetworkAcl(d, networkAclId)
	if err != nil {
		return data, err
	}
	num := int64(d.Get("rule_number").(int))
	direction := d.Get("direction")
	found := false
	for _, entry := range acl["NetworkAclEntrySet"].([]interface{}) {
		m := entry.(map[string]interface{})
		if num == int64(m["RuleNumber"].(float64)) && direction == m["Direction"] {
			found = true
			data = m
			break
		}
	}
	if !found {
		return data, fmt.Errorf("network acl not exist")
	}
	return data, err
}

func (s *VpcService) ReadNetworkAclAssociate(d *schema.ResourceData, networkAclId string, subnetId string) (data map[string]interface{}, err error) {
	data, err = s.ReadSubnet(d, subnetId)
	if err != nil {
		return data, err
	}
	if data["NetworkAclId"] != networkAclId {
		return data, fmt.Errorf("network acl %s not associate sunbet %s", networkAclId, subnetId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetNetworkAcls(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "NetworkAclId",
			Type:    TransformWithN,
		},
		"vpc_ids": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadNetworkAcls(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "NetworkAclName",
		idFiled:     "NetworkAclId",
		targetField: "network_acls",
		extra: map[string]SdkResponseMapping{
			"NetworkAclId": {
				Field:    "id",
				KeepAuto: true,
			},
			"NetworkAclName": {
				Field:    "name",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) ReadAndSetNetworkAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadNetworkAcl(d, "")
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"NetworkAclEntrySet": {
			Field: "network_acl_entries",
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return err
}

func (s *VpcService) ReadAndSetNetworkAclEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadNetworkAclEntry(d, d.Get("network_acl_id").(string))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *VpcService) ReadAndSetNetworkAclAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadNetworkAclAssociate(d, d.Get("network_acl_id").(string), d.Get("subnet_id").(string))
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading network acl associate %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, nil)
			return nil
		}
	})
}

func (s *VpcService) CreateNetworkAclCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"vpc_id":           {},
		"network_acl_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateNetworkAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateNetworkAcl(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("NetworkAclId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("network_acl_id", id)
		},
	}
	return callback, err
}

func (s *VpcService) CreateNetworkAclEntryWithAclCall(d *schema.ResourceData, r *schema.Resource) (callbacks []ApiCall, err error) {
	if entries, ok := d.GetOk("network_acl_entries"); ok {
		for _, entry := range entries.(*schema.Set).List() {
			var (
				req      map[string]interface{}
				callback ApiCall
			)
			key := networkAclEntryHash(entry)
			transform := make(map[string]SdkReqTransform)
			for k := range entry.(map[string]interface{}) {
				key := "network_acl_entries." + strconv.Itoa(key) + "." + k
				if k == "icmp_type" || k == "icmp_code" {
					transform[key] = SdkReqTransform{
						mapping: Downline2Hump(k),
						ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
							return d.Get(key), true
						},
					}
				} else {
					transform[key] = SdkReqTransform{mapping: Downline2Hump(k)}
				}

			}
			req, err = SdkRequestAutoMapping(d, r, false, transform, nil)
			if err != nil {
				return callbacks, err
			}
			callback, err = s.CreateNetworkAclEntryCommonCall(req, false)
			if err != nil {
				return callbacks, err
			}
			callbacks = append(callbacks, callback)
		}
	}
	return callbacks, err
}

func (s *VpcService) CreateNetworkAclEntryCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"icmp_type": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return d.Get("icmp_type"), true
			},
		},
		"icmp_code": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return d.Get("icmp_code"), true
			},
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		false,
	})
	if req["Protocol"] != "icmp" {
		delete(req, "IcmpType")
		delete(req, "IcmpCode")
	}

	if err != nil {
		return callback, err
	}
	return s.CreateNetworkAclEntryCommonCall(req, true)
}

func (s *VpcService) CreateNetworkAclAssociateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AssociateNetworkAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AssociateNetworkAcl(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("network_acl_id").(string) + ":" + d.Get("subnet_id").(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateNetworkAclEntryCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	//check
	if req["Protocol"] == "icmp" {
		if _, ok := req["IcmpType"]; !ok {
			return callback, fmt.Errorf("NetworkAcl Protocol is icmp,must set IcmpType")
		}
		if _, ok := req["IcmpCode"]; !ok {
			return callback, fmt.Errorf("NetworkAcl Protocol is icmp,must set IcmpCode")
		}
	}
	if req["Protocol"] == "udp" || req["Protocol"] == "tcp" {
		if _, ok := req["PortRangeFrom"]; !ok {
			return callback, fmt.Errorf("NetworkAcl Protocol is udp/tcp,must set PortRangeFrom")
		}
		if _, ok := req["PortRangeTo"]; !ok {
			return callback, fmt.Errorf("NetworkAcl Protocolt is udp/tcp,must set PortRangeTo")
		}
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateNetworkAclEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			(*(call.param))["NetworkAclId"] = d.Get("network_acl_id")
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateNetworkAclEntry(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if isSetId {
				_, err = s.ReadNetworkAclEntry(d, (*(call.param))["NetworkAclId"].(string))
				if err != nil {
					return err
				}
				d.SetId((*(call.param))["NetworkAclId"].(string) + ":" + strconv.Itoa(d.Get("rule_number").(int)) + ":" + d.Get("direction").(string))
			}
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateNetworkAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateNetworkAclCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	entries, err := s.CreateNetworkAclEntryWithAclCall(d, r)
	if err != nil {
		return err
	}
	for _, entryCall := range entries {
		callbacks = append(callbacks, entryCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, false)
}

func (s *VpcService) CreateNetworkAclEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateNetworkAclEntryCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) CreateNetworkAclAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateNetworkAclAssociateCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifyNetworkAclCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"network_acl_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["NetworkAclId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyNetworkAcl",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyNetworkAcl(call.param)
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

func (s *VpcService) ModifyNetworkAclEntryCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"description": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["NetworkAclEntryId"] = d.Get("network_acl_entry_id")
		return s.ModifyNetworkAclEntryCommonCall(req)
	}
	return callback, err
}

func (s *VpcService) ModifyNetworkAclEntryWithAclCall(d *schema.ResourceData, r *schema.Resource) (callbacks []ApiCall, err error) {
	if d.HasChange("network_acl_entries") {
		o, n := d.GetChange("network_acl_entries")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		//generate new hashcode without description
		mayAdd := schema.NewSet(networkAclEntrySimpleHash, ns.Difference(os).List())
		mayRemove := schema.NewSet(networkAclEntrySimpleHash, os.Difference(ns).List())
		addCache := make(map[int]interface{})
		for _, entry := range mayAdd.List() {
			index := networkAclEntrySimpleHash(entry)
			addCache[index] = entry
		}
		//compare hashcode without description
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
				index := networkAclEntrySimpleHash(entry)
				req := make(map[string]interface{})
				req["Description"] = addCache[index].(map[string]interface{})["description"]
				req["NetworkAclEntryId"] = entry.(map[string]interface{})["network_acl_entry_id"]
				callback, err = s.ModifyNetworkAclEntryCommonCall(req)
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
				callback, err = s.RemoveNetworkAclEntryCommonCall(d.Id(), entry.(map[string]interface{})["network_acl_entry_id"].(string))
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
				index := networkAclEntryHash(entry)
				transform := make(map[string]SdkReqTransform)
				for k := range entry.(map[string]interface{}) {
					key := "network_acl_entries." + strconv.Itoa(index) + "." + k
					transform[key] = SdkReqTransform{mapping: Downline2Hump(k)}
				}
				req, err = SdkRequestAutoMapping(d, r, false, transform, nil)
				if err != nil {
					return callbacks, err
				}
				callback, err = s.CreateNetworkAclEntryCommonCall(req, false)
				if err != nil {
					return callbacks, err
				}
				callbacks = append(callbacks, callback)
			}
		}
	}
	return callbacks, err
}

func (s *VpcService) ModifyNetworkAclEntryCommonCall(req map[string]interface{}) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "ModifyNetworkAclEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyNetworkAclEntry(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *VpcService) ModifyNetworkAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyNetworkAclCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	entries, err := s.ModifyNetworkAclEntryWithAclCall(d, r)
	if err != nil {
		return err
	}
	for _, entryCall := range entries {
		callbacks = append(callbacks, entryCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *VpcService) ModifyNetworkAclEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifyNetworkAclEntryCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *VpcService) RemoveNetworkAclCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"NetworkAclId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteNetworkAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteNetworkAcl(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadNetworkAcl(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading network acl when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveNetworkAclEntryCommonCall(aclId string, entryId string) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"NetworkAclId":      aclId,
		"NetworkAclEntryId": entryId,
	}
	callback = ApiCall{
		param:  &req,
		action: "DeleteNetworkAclEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteNetworkAclEntry(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				data, callErr := s.ReadNetworkAcl(d, aclId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading nat when delete %q, %s", d.Id(), callErr))
					}
				}
				if len(data["NetworkAclEntrySet"].([]interface{})) == 0 {
					return nil
				} else {
					found := false
					for _, item := range data["NetworkAclEntrySet"].([]interface{}) {
						if item.(map[string]interface{})["NetworkAclEntryId"] == entryId {
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

func (s *VpcService) RemoveNetworkAclAssociateCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"NetworkAclId": d.Get("network_acl_id"),
		"SubnetId":     d.Get("subnet_id"),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DisassociateNetworkAcl",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DisassociateNetworkAcl(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				aclId := (*(call.param))["NetworkAclId"].(string)
				subnetId := (*(call.param))["SubnetId"].(string)
				_, callErr := s.ReadNetworkAclAssociate(d, aclId, subnetId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading network acl associate when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveNetworkAcl(d *schema.ResourceData) (err error) {
	call, err := s.RemoveNetworkAclCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveNetworkAclEntry(d *schema.ResourceData) (err error) {
	call, err := s.RemoveNetworkAclEntryCommonCall(d.Get("network_acl_id").(string), d.Get("network_acl_entry_id").(string))
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveNetworkAclAssociate(d *schema.ResourceData) (err error) {
	call, err := s.RemoveNetworkAclAssociateCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadSecurityGroups(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeSecurityGroups"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeSecurityGroups(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeSecurityGroups(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("SecurityGroupSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadSecurityGroup(d *schema.ResourceData, securityGroupId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if securityGroupId == "" {
		securityGroupId = d.Id()
	}
	req := map[string]interface{}{
		"SecurityGroupId.1": securityGroupId,
	}
	results, err = s.ReadSecurityGroups(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("security group %s not exist ", securityGroupId)
	}
	return data, err
}

func (s *VpcService) ReadSecurityGroupEntry(d *schema.ResourceData, securityGroupId string) (data map[string]interface{}, err error) {
	sg, err := s.ReadSecurityGroup(d, securityGroupId)
	if err != nil {
		return data, err
	}
	found := false
	for _, entry := range sg["SecurityGroupEntrySet"].([]interface{}) {
		h1 := securityGroupEntrySimpleHashWithHump(entry)
		h2 := securityGroupEntrySimpleHash(d)
		if h1 == h2 {
			found = true
			data = entry.(map[string]interface{})
			break
		}
	}
	if !found {
		return data, fmt.Errorf("security group entry not exist")
	}
	return data, err
}

func (s *VpcService) ReadAndSetSecurityGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "SecurityGroupId",
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
	data, err := s.ReadSecurityGroups(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "SecurityGroupName",
		idFiled:     "SecurityGroupId",
		targetField: "security_groups",
		extra: map[string]SdkResponseMapping{
			"SecurityGroupId": {
				Field:    "id",
				KeepAuto: true,
			},
			"SecurityGroupName": {
				Field:    "name",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) ReadAndSetSecurityGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadSecurityGroup(d, "")
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"SecurityGroupEntrySet": {
			Field: "security_group_entries",
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return err
}

func (s *VpcService) ReadAndSetSecurityGroupEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadSecurityGroupEntry(d, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *VpcService) CreateSecurityGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"vpc_id":              {},
		"security_group_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateSecurityGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateSecurityGroup(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("SecurityGroup.SecurityGroupId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("security_group_id", d.Id())
		},
	}
	return callback, err
}

func (s *VpcService) CreateSecurityGroupEntryWithSgCall(d *schema.ResourceData, r *schema.Resource) (callbacks []ApiCall, err error) {
	if entries, ok := d.GetOk("security_group_entries"); ok {
		for _, entry := range entries.(*schema.Set).List() {
			var (
				req      map[string]interface{}
				callback ApiCall
			)
			index := securityGroupEntryHash(entry)
			transform := make(map[string]SdkReqTransform)
			for k := range entry.(map[string]interface{}) {
				key := "security_group_entries." + strconv.Itoa(index) + "." + k
				if k == "icmp_type" || k == "icmp_code" {
					transform[key] = SdkReqTransform{
						mapping: Downline2Hump(k),
						ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
							return d.Get(key), true
						},
					}
				} else {
					transform[key] = SdkReqTransform{mapping: Downline2Hump(k)}
				}

			}
			req, err = SdkRequestAutoMapping(d, r, false, transform, nil)
			if err != nil {
				return callbacks, err
			}
			callback, err = s.CreateSecurityGroupEntryCommonCall(req, false)
			if err != nil {
				return callbacks, err
			}
			callbacks = append(callbacks, callback)
		}
	}
	return callbacks, err
}

func (s *VpcService) CreateSecurityGroupEntryCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"icmp_type": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return d.Get("icmp_type"), true
			},
		},
		"icmp_code": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return d.Get("icmp_code"), true
			},
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		false,
	})
	if req["Protocol"] != "icmp" {
		delete(req, "IcmpType")
		delete(req, "IcmpCode")
	}
	if err != nil {
		return callback, err
	}
	return s.CreateSecurityGroupEntryCommonCall(req, true)
}

func (s *VpcService) CreateSecurityGroupEntryCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	//check
	if req["Protocol"] == "icmp" {
		if _, ok := req["IcmpType"]; !ok {
			return callback, fmt.Errorf("SecurityGroup entry Protocol is icmp,must set IcmpType")
		}
		if _, ok := req["IcmpCode"]; !ok {
			return callback, fmt.Errorf("SecurityGroup entry Protocol is icmp,must set IcmpCode")
		}
	}
	if req["Protocol"] == "udp" || req["Protocol"] == "tcp" {
		if _, ok := req["PortRangeFrom"]; !ok {
			return callback, fmt.Errorf("SecurityGroup entry Protocol is udp/tcp,must set PortRangeFrom")
		}
		if _, ok := req["PortRangeTo"]; !ok {
			return callback, fmt.Errorf("SecurityGroup entry Protocolt is udp/tcp,must set PortRangeTo")
		}
	}
	callback = ApiCall{
		param:  &req,
		action: "AuthorizeSecurityGroupEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			(*(call.param))["SecurityGroupId"] = d.Get("security_group_id")
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AuthorizeSecurityGroupEntry(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if isSetId {
				var data map[string]interface{}
				data, err = s.ReadSecurityGroupEntry(d, (*(call.param))["SecurityGroupId"].(string))
				if err != nil {
					return err
				}
				buf := securityGroupEntryHashBase(data, true)
				d.SetId((*(call.param))["SecurityGroupId"].(string) + buf.String())
			}
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateSecurityGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.CreateSecurityGroupCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	entries, err := s.CreateSecurityGroupEntryWithSgCall(d, r)
	if err != nil {
		return err
	}
	for _, entryCall := range entries {
		callbacks = append(callbacks, entryCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, false)
}

func (s *VpcService) CreateSecurityGroupEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateSecurityGroupEntryCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifySecurityGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"security_group_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["SecurityGroupId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifySecurityGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
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

func (s *VpcService) ModifySecurityGroupEntryCommonCall(req map[string]interface{}) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "ModifySecurityGroupEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifySecurityGroupEntry(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *VpcService) ModifySecurityGroupEntryWithSgCall(d *schema.ResourceData, r *schema.Resource) (callbacks []ApiCall, err error) {
	if d.HasChange("security_group_entries") {
		o, n := d.GetChange("security_group_entries")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		//generate new hashcode without description
		mayAdd := schema.NewSet(securityGroupEntrySimpleHash, ns.Difference(os).List())
		mayRemove := schema.NewSet(securityGroupEntrySimpleHash, os.Difference(ns).List())
		addCache := make(map[int]interface{})
		for _, entry := range mayAdd.List() {
			index := securityGroupEntrySimpleHash(entry)
			addCache[index] = entry
		}
		//compare hashcode without description
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
				index := securityGroupEntrySimpleHash(entry)
				req := make(map[string]interface{})
				req["Description"] = addCache[index].(map[string]interface{})["description"]
				req["SecurityGroupEntryId"] = entry.(map[string]interface{})["security_group_entry_id"]
				callback, err = s.ModifySecurityGroupEntryCommonCall(req)
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
				callback, err = s.RemoveSecurityGroupEntryCommonCall(d.Id(), entry.(map[string]interface{})["security_group_entry_id"].(string))
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
				index := securityGroupEntryHash(entry)
				transform := make(map[string]SdkReqTransform)
				for k := range entry.(map[string]interface{}) {
					key := "security_group_entries." + strconv.Itoa(index) + "." + k
					transform[key] = SdkReqTransform{mapping: Downline2Hump(k)}
				}
				req, err = SdkRequestAutoMapping(d, r, false, transform, nil)
				if err != nil {
					return callbacks, err
				}
				callback, err = s.CreateSecurityGroupEntryCommonCall(req, false)
				if err != nil {
					return callbacks, err
				}
				callbacks = append(callbacks, callback)
			}
		}
	}
	return callbacks, err
}

func (s *VpcService) ModifySecurityGroupEntryCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"description": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["SecurityGroupEntryId"] = d.Get("security_group_entry_id")
		return s.ModifySecurityGroupEntryCommonCall(req)
	}
	return callback, err
}

func (s *VpcService) ModifySecurityGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifySecurityGroupCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	entries, err := s.ModifySecurityGroupEntryWithSgCall(d, r)
	if err != nil {
		return err
	}
	for _, entryCall := range entries {
		callbacks = append(callbacks, entryCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *VpcService) ModifySecurityGroupEntry(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	call, err := s.ModifySecurityGroupEntryCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, call)
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *VpcService) RemoveSecurityGroupCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"SecurityGroupId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteSecurityGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteSecurityGroup(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadSecurityGroup(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading security group when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveSecurityGroupEntryCommonCall(sgId string, entryId string) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"SecurityGroupId":      sgId,
		"SecurityGroupEntryId": entryId,
	}
	callback = ApiCall{
		param:  &req,
		action: "RevokeSecurityGroupEntry",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.RevokeSecurityGroupEntry(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				data, callErr := s.ReadSecurityGroup(d, sgId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading security group entry when delete %q, %s", d.Id(), callErr))
					}
				}
				if len(data["SecurityGroupEntrySet"].([]interface{})) == 0 {
					return nil
				} else {
					found := false
					for _, item := range data["SecurityGroupEntrySet"].([]interface{}) {
						if item.(map[string]interface{})["SecurityGroupEntryId"] == entryId {
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

func (s *VpcService) RemoveSecurityGroup(d *schema.ResourceData) (err error) {
	call, err := s.RemoveSecurityGroupCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveSecurityGroupEntry(d *schema.ResourceData) (err error) {
	call, err := s.RemoveSecurityGroupEntryCommonCall(d.Get("security_group_id").(string), d.Get("security_group_entry_id").(string))
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadVpnGateways(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeVpnGateways"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeVpnGateways(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeVpnGateways(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("VpnGatewaySet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadVpnGateway(d *schema.ResourceData, vpnGatewayId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if vpnGatewayId == "" {
		vpnGatewayId = d.Id()
	}
	req := map[string]interface{}{
		"VpnGatewayId.1": vpnGatewayId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}
	results, err = s.ReadVpnGateways(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Vpn Gateway  %s not exist ", vpnGatewayId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetVpnGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadVpnGateway(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(callErr)
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, chargeExtraForVpc(data))
			return nil
		}
	})
}

func (s *VpcService) ReadAndSetVpnGateways(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "VpnGatewayId",
			Type:    TransformWithN,
		},
		"project_ids": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"vpc_ids": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadVpnGateways(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "VpnGatewayName",
		idFiled:     "VpnGatewayId",
		targetField: "vpn_gateways",
		extra: map[string]SdkResponseMapping{
			"VpnGatewayId": {
				Field:    "id",
				KeepAuto: true,
			},
			"VpnGatewayName": {
				Field:    "name",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) CreateVpnGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	if _, ok := req["PurchaseTime"]; !ok && req["ChargeType"] == "Monthly" {
		return callback, fmt.Errorf("ChargeType is Monthly must set PurchaseTime")
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateVpnGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateVpnGateway(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("VpnGateway.VpnGatewayId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateVpnGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateVpnGatewayCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifyVpnGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["VpnGatewayId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyVpnGateway",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyVpnGateway(call.param)
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

func (s *VpcService) ModifyVpnGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyVpnGatewayCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveVpnGatewayCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"VpnGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteVpnGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteVpnGateway(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadVpnGateway(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading vpn gateway when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveVpnGateway(d *schema.ResourceData) (err error) {
	call, err := s.RemoveVpnGatewayCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadVpnCustomerGateways(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeCustomerGateways"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeCustomerGateways(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeCustomerGateways(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("CustomerGatewaySet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadVpnCustomerGateway(d *schema.ResourceData, vpnCustomerGatewayId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if vpnCustomerGatewayId == "" {
		vpnCustomerGatewayId = d.Id()
	}
	req := map[string]interface{}{
		"CustomerGatewayId.1": vpnCustomerGatewayId,
	}
	results, err = s.ReadVpnCustomerGateways(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Customer gateway %s not exist ", vpnCustomerGatewayId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetVpnCustomerGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadVpnCustomerGateway(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *VpcService) ReadAndSetVpnCustomerGateways(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "CustomerGatewayId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadVpnCustomerGateways(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "CustomerGatewayName",
		idFiled:     "CustomerGatewayId",
		targetField: "customer_gateways",
		extra: map[string]SdkResponseMapping{
			"CustomerGatewayId": {
				Field:    "id",
				KeepAuto: true,
			},
			"CustomerGatewayName": {
				Field:    "name",
				KeepAuto: true,
			},
		},
	})
}

func (s *VpcService) CreateVpnCustomerGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateCustomerGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateCustomerGateway(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("CustomerGateway.CustomerGatewayId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateVpnCustomerGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateVpnCustomerGatewayCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifyVpnCustomerGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["CustomerGatewayId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyCustomerGateway",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyCustomerGateway(call.param)
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

func (s *VpcService) ModifyVpnCustomerGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyVpnCustomerGatewayCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveVpnCustomerGatewayCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"CustomerGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteCustomerGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteCustomerGateway(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadVpnCustomerGateway(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading customer gateway when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveVpnCustomerGateway(d *schema.ResourceData) (err error) {
	call, err := s.RemoveVpnCustomerGatewayCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadVpnTunnels(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeVpnTunnels"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeVpnTunnels(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeVpnTunnels(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("VpnTunnelSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadVpnTunnel(d *schema.ResourceData, vpnTunnelId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if vpnTunnelId == "" {
		vpnTunnelId = d.Id()
	}
	req := map[string]interface{}{
		"VpnTunnelId.1": vpnTunnelId,
	}
	results, err = s.ReadVpnTunnels(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Vpn tunnel %s not exist ", vpnTunnelId)
	}
	return data, err
}

func (s *VpcService) ReadAndSetVpnTunnel(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadVpnTunnel(d, "")
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"IkeDHGroup": {
			Field: "ike_dh_group",
			FieldRespFunc: func(i interface{}) interface{} {
				result, _ := strconv.Atoi(i.(string))
				return result
			},
		},
		"Type": {
			Field: "type",
			FieldRespFunc: func(i interface{}) interface{} {
				return Downline2Hump(i.(string))
			},
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return err
}

func (s *VpcService) ReadAndSetVpnTunnels(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "VpnTunnelId",
			Type:    TransformWithN,
		},
		"vpn_gateway_ids": {
			mapping: "vpn-gateway-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadVpnTunnels(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "VpnTunnelName",
		idFiled:     "VpnTunnelId",
		targetField: "vpn_tunnels",
		extra: map[string]SdkResponseMapping{
			"VpnTunnelId": {
				Field:    "id",
				KeepAuto: true,
			},
			"VpnTunnelName": {
				Field:    "name",
				KeepAuto: true,
			},
			"IkeDHGroup": {
				Field: "ike_dh_group",
				FieldRespFunc: func(i interface{}) interface{} {
					result, _ := strconv.Atoi(i.(string))
					return result
				},
			},
		},
	})
}

func (s *VpcService) CreateVpnTunnelCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"ike_dh_group": {
			mapping: "IkeDHGroup",
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		false,
	})
	if err != nil {
		return callback, err
	}
	//check
	if _, ok := req["VpnGreIp"]; !ok && req["Type"] == "GreOverIpsec" {
		return callback, fmt.Errorf("Vpn tunnel type is GreOverIpsec must set VpnGreIp ")
	}
	if _, ok := req["HaVpnGreIp"]; !ok && req["Type"] == "GreOverIpsec" {
		return callback, fmt.Errorf("Vpn tunnel type is GreOverIpsec must set HaVpnGreIp ")
	}
	if _, ok := req["CustomerGreIp"]; !ok && req["Type"] == "GreOverIpsec" {
		return callback, fmt.Errorf("Vpn tunnel type is GreOverIpsec must set CustomerGreIp ")
	}
	if _, ok := req["HaCustomerGreIp"]; !ok && req["Type"] == "GreOverIpsec" {
		return callback, fmt.Errorf("Vpn tunnel type is GreOverIpsec must set HaCustomerGreIp ")
	}
	callback = ApiCall{
		param:  &req,
		action: " CreateVpnTunnel",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateVpnTunnel(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("VpnTunnel.VpnTunnelId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *VpcService) CreateVpnTunnel(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateVpnTunnelCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ModifyVpnTunnelCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["VpnTunnelId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyVpnTunnel",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyVpnTunnel(call.param)
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

func (s *VpcService) ModifyVpnTunnel(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyVpnTunnelCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) RemoveVpnTunnelCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"VpnTunnelId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteVpnTunnel",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteVpnTunnel(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadVpnTunnel(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading vpn tunnel when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) RemoveVpnTunnel(d *schema.ResourceData) (err error) {
	call, err := s.RemoveVpnTunnelCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *VpcService) ReadAvailabilityZones(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeAvailabilityZones"
	logger.Debug(logger.ReqFormat, action, condition)
	logger.Debug(logger.ReqFormat, action, conn.Endpoint)
	resp, err = conn.DescribeAvailabilityZones(nil)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("AvailabilityZoneInfo", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ReadAndSetAvailabilityZones(d *schema.ResourceData, r *schema.Resource) (err error) {
	req, err := mergeDataSourcesReq(d, r, nil)
	if err != nil {
		return err
	}
	data, err := s.ReadAvailabilityZones(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "AvailabilityZoneName",
		idFiled:     "AvailabilityZoneName",
		targetField: "availability_zones",
	})
}
