package ksyun

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func (s *VpcService) CreateDirectConnectTunnel(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(
		context.TODO(), d, s.client, false)

	createCall, err := s.createDirectConnectTunnelCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(createCall)
	return apiProcess.Run()
}

func (s *VpcService) createDirectConnectTunnelCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "CreateDirectConnectInterface",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.CreateDirectConnectInterface(call.param)
				if err != nil {
					return resp, err
				}
				tunnelId, err := getSdkValue("DirectConnectInterface.0.DirectConnectInterfaceId", resp)
				if err != nil {
					return resp, err
				}
				d.SetId(tunnelId.(string))

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

func (s *VpcService) readAndSetDirectConnectTunnel(d *schema.ResourceData, r *schema.Resource) (err error) {
	resp, err := s.readDirectConnectTunnel(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, resp, nil)
	return nil
}

func (s *VpcService) readDirectConnectTunnel(d *schema.ResourceData, dTunnelId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if dTunnelId == "" {
		dTunnelId = d.Id()
	}
	req := map[string]interface{}{
		"DirectConnectInterfaceId.1": dTunnelId,
	}
	results, err = s.ReadVpnTunnels(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("direct connect tunnel %s not exist ", dTunnelId)
	}
	return data, err
}

func (s *VpcService) readDirectConnectTunnels(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeDirectConnectInterfaces"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeDirectConnectInterfaces(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeDirectConnectInterfaces(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("DirectConnectInterfaceSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) ModifyDirectConnectTunnel(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(
		context.TODO(), d, s.client, false)

	call, err := s.modifyDirectConnectTunnelCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) modifyDirectConnectTunnelCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	if !d.HasChange("direct_connect_interface_name") {
		return callback, nil
	}

	req := map[string]interface{}{
		"DirectConnectInterfaceId":   d.Id(),
		"DirectConnectInterfaceName": d.Get("direct_connect_interface_name"),
	}
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "ModifyDirectConnectInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyDirectConnectInterface(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	return callback, err
}

func (s *VpcService) removeDirectConnectTunnelCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectInterfaceId": d.Id(),
	}
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "DeleteDirectConnectInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteDirectConnectInterface(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.readDirectConnectTunnel(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading direct connect tunnel when delete %q, %s", d.Id(), callErr))
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

func (s *VpcService) attachDirectConnectGatewayCall(d *schema.ResourceData, r *schema.Resource, interId string) (callback ApiCall, err error) {
	if d.Get("direct_connect_interface_id") == "" {
		return callback, fmt.Errorf("direct_connect_gateway_id is required for attach direct connect gateway")
	}

	req := map[string]interface{}{
		"DirectConnectInterfaceId": interId,
		"DirectConnectGatewayId":   d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "AttachDirectConnectGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AttachDirectConnectGateway(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	return callback, err
}

func (s *VpcService) detachDirectConnectGatewayCall(d *schema.ResourceData, r *schema.Resource, interId string) (callback ApiCall, err error) {
	if d.Get("direct_connect_interface_id") == "" {
		return callback, fmt.Errorf("direct_connect_gateway_id is required for attach direct connect gateway")
	}

	req := map[string]interface{}{
		"DirectConnectInterfaceId": interId,
		"DirectConnectGatewayId":   d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DetachDirectConnectGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DetachDirectConnectGateway(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	return callback, err
}

func (s *VpcService) createDirectConnectGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"vpc_id":                      {},
		"direct_connect_gateway_name": {},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{onlyTransform: true})
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "CreateDirectConnectGateway",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.CreateDirectConnectGateway(call.param)
				if err != nil {
					return resp, err
				}

				gatewayId, err := getSdkValue("DirectConnectGateway.DirectConnectGatewayId", resp)
				if err != nil {
					return resp, err
				}
				d.SetId(gatewayId.(string))

				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return
}

func (s *VpcService) deleteDirectConnectGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DeleteDirectConnectGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteCustomerGateway(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
}

func (s *VpcService) modifyDirectConnectGatewayCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	if !d.HasChange("direct_connect_gateway_name") {
		return callback, nil
	}
	req := map[string]interface{}{
		"DirectConnectGatewayId":   d.Id(),
		"DirectConnectGatewayName": d.Get("direct_connect_gateway_name"),
	}
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "ModifyDirectConnectGateway",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyDirectConnectGateway(call.param)
				if err != nil {
					return resp, err
				}
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return
}

func (s *VpcService) ReadAndSetDirectConnectGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	resp, err := s.readDirectConnectGateway(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, resp, nil)
	return nil
}

func (s *VpcService) readDirectConnectGateway(d *schema.ResourceData, gatewayId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if gatewayId == "" {
		gatewayId = d.Id()
	}
	req := map[string]interface{}{
		"DirectConnectGatewayId.1": gatewayId,
	}
	results, err = s.readDirectConnectGateways(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("direct connect gateway %s not exist ", gatewayId)
	}
	return data, err
}

func (s *VpcService) readDirectConnectGateways(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeDirectConnectGateways"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeDirectConnectGateways(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeDirectConnectGateways(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("DirectConnectGatewaySet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) attachDirectConnectGatewayWithVpcCall(d *schema.ResourceData, r *schema.Resource, vpcid string) (callback ApiCall, err error) {
	if d.Get("vpc_id") == "" {
		return callback, fmt.Errorf("vpc_id is required for attach direct connect gateway")
	}

	req := map[string]interface{}{
		"VpcId":                  vpcid,
		"DirectConnectGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "AttachDirectConnectGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AttachDirectConnectGateway(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	return callback, err
}

func (s *VpcService) detachDirectConnectGatewayWithVpcCall(d *schema.ResourceData, r *schema.Resource, interId string) (callback ApiCall, err error) {
	if d.Get("vpc_id") == "" {
		return callback, fmt.Errorf("vpc_id is required for attach direct connect gateway")
	}

	req := map[string]interface{}{
		"VpcId":                  interId,
		"DirectConnectGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DetachDirectConnectGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DetachDirectConnectGateway(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	return callback, err
}
