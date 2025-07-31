package ksyun

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func (s *VpcService) CreateDirectConnectInterface(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(
		context.TODO(), d, s.client, false)

	createCall, err := s.createDirectConnectInterfaceCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(createCall)
	return apiProcess.Run()
}

func (s *VpcService) createDirectConnectInterfaceCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["EnableIpv6"] = d.Get("enable_ipv6")
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
				tunnelId, err := getSdkValue("DirectConnectInterface.DirectConnectInterfaceId", *resp)
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

func (s *VpcService) ReadAndSetDirectConnectInterface(d *schema.ResourceData, r *schema.Resource) (err error) {
	resp, err := s.readDirectConnectInterface(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, resp, nil)
	return nil
}

func (s *VpcService) readDirectConnectInterface(d *schema.ResourceData, dInterfaceId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if dInterfaceId == "" {
		dInterfaceId = d.Id()
	}
	req := map[string]interface{}{
		"DirectConnectInterfaceId.1": dInterfaceId,
	}
	results, err = s.readDirectConnectInterfaces(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("direct connect tunnel %s not exist ", dInterfaceId)
	}
	return data, err
}

func (s *VpcService) readDirectConnectInterfaces(condition map[string]interface{}) (data []interface{}, err error) {
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

func (s *VpcService) ModifyDirectConnectInterface(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(
		context.TODO(), d, s.client, false)

	call, err := s.modifyDirectConnectInterfaceCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	bfdCall, err := s.modifyDirectConnectInterfaceBfdCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(bfdCall)

	return apiProcess.Run()
}

func (s *VpcService) modifyDirectConnectInterfaceBfdCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	if !d.HasChange("bfd_config_id") {
		return callback, err
	}
	req := map[string]interface{}{
		"DirectConnectInterfaceId": d.Id(),
		"ReliabilityMethod":        d.Get("reliability_method"),
		"BfdConfigId":              d.Get("bfd_config_id"),
	}
	callback = ApiCall{
		param:  &req,
		action: "ModifyDirectConnectInterfaceBfd",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyDirectConnectInterfaceBfd(call.param)
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

func (s *VpcService) modifyDirectConnectInterfaceCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	if !d.HasChange("direct_connect_interface_name") {
		return callback, nil
	}

	req := map[string]interface{}{
		"DirectConnectInterfaceId":   d.Id(),
		"DirectConnectInterfaceName": d.Get("direct_connect_interface_name"),
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

func (s *VpcService) RemoveDirectConnectInterface(d *schema.ResourceData) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.removeDirectConnectInterfaceCall(d)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) removeDirectConnectInterfaceCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectInterfaceId": d.Id(),
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
				_, callErr := s.readDirectConnectInterface(d, "")
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

func (s *VpcService) attachDirectConnectGatewayCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectInterfaceId": d.Get("direct_connect_interface_id"),
		"DirectConnectGatewayId":   d.Get("direct_connect_gateway_id"),
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

func (s *VpcService) detachDirectConnectGatewayCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectInterfaceId": d.Get("direct_connect_interface_id"),
		"DirectConnectGatewayId":   d.Get("direct_connect_gateway_id"),
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

func (s *VpcService) CreateDirectConnectGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.createDirectConnectGatewayCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	// err = s.attachOrDettachInterfaceOnGateway(d, r, &apiProcess)
	// if err != nil {
	// 	return err
	// }

	return apiProcess.Run()
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
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				gatewayId, err := getSdkValue("DirectConnectGateway.DirectConnectGatewayId", *resp)
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

func (s *VpcService) RemoveDirectConnectGateway(d *schema.ResourceData) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.deleteDirectConnectGatewayCall(d)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) deleteDirectConnectGatewayCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DeleteDirectConnectGateway",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteDirectConnectGateway(call.param)
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

func (s *VpcService) ModifyDirectConnectGateway(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.modifyDirectConnectGatewayCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	if d.HasChange("vpc_id") {
		old, cur := d.GetChange("vpc_id")
		if old != "" {
			detachCall, err := s.detachDirectConnectGatewayWithVpcCall(d, r, old.(string))
			if err != nil {
				return err
			}
			apiProcess.PutCalls(detachCall)
		}
		if cur != "" {
			attachCall, err := s.attachDirectConnectGatewayWithVpcCall(d, r, cur.(string))
			if err != nil {
				return err
			}
			apiProcess.PutCalls(attachCall)
		}
	}

	// err = s.attachOrDettachInterfaceOnGateway(d, r, &apiProcess)
	// if err != nil {
	// 	return err
	// }

	return apiProcess.Run()
}

//
// func (s *VpcService) attachOrDettachInterfaceOnGateway(d *schema.ResourceData, r *schema.Resource, apiProcess *ApiProcess) (err error) {
// 	if !d.HasChange("direct_connect_interface_id") {
// 		return nil
// 	}
//
// 	old, cur := d.GetChange("direct_connect_interface_id")
// 	if old != "" {
// 		detachCall, err := s.detachDirectConnectGatewayCall(d, r, old.(string))
// 		if err != nil {
// 			return err
// 		}
// 		apiProcess.PutCalls(detachCall)
// 	}
// 	if cur != "" {
// 		attachCall, err := s.attachDirectConnectGatewayCall(d, r, cur.(string))
// 		if err != nil {
// 			return err
// 		}
// 		apiProcess.PutCalls(attachCall)
// 	}
//
// 	return nil
// }

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
	req := map[string]interface{}{
		"VpcId":                  vpcid,
		"DirectConnectGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "AttachDirectConnectGatewayWithVpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AttachDirectConnectGatewayWithVpc(call.param)
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
	req := map[string]interface{}{
		"VpcId":                  interId,
		"DirectConnectGatewayId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DetachDirectConnectGatewayWithVpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DetachDirectConnectGatewayWithVpc(call.param)
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

func (s *VpcService) CreateDirectConnectGatewayRoute(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.createDirectConnectGatewayRouteCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) RemoveDirectConnectGatewayRoute(d *schema.ResourceData) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.deleteDirectConnectGatewayRouteCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) createDirectConnectGatewayRouteCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "CreateDirectConnectGatewayRoute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.CreateDirectConnectGatewayRoute(call.param)
				if err != nil {
					return resp, err
				}

				routeId, err := getSdkValue("DirectConnectGatewayRoute.DirectConnectGatewayRouteId", *resp)
				if err != nil {
					return nil, err
				}
				d.SetId(routeId.(string))
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

func (s *VpcService) deleteDirectConnectGatewayRouteCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectGatewayRouteId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DetachDirectConnectGatewayWithVpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteDirectConnectGatewayRoute(call.param)
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

func (s *VpcService) ReadAndSetDirectConnectGatewayRoute(d *schema.ResourceData, r *schema.Resource) (err error) {
	resp, err := s.readDirectConnectGatewayRoute(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, resp, nil)
	return nil
}

func (s *VpcService) readDirectConnectGatewayRoute(d *schema.ResourceData, dcb string) (data map[string]interface{}, err error) {
	var results []interface{}
	if dcb == "" {
		dcb = d.Get("destination_cidr_block").(string)
	}
	req := map[string]interface{}{
		"DirectConnectGatewayId": d.Get("direct_connect_gateway_id"),
		"Filter.1.Name":          "cidr-block",
		"Filter.1.Value.1":       dcb,
	}
	results, err = s.readDirectConnectGatewayRoutes(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		vv := v.(map[string]interface{})
		if vv["DirectConnectGatewayRouteId"] == d.Id() {
			data = vv
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("direct connect gateway %s not exist ", dcb)
	}
	return data, err
}

func (s *VpcService) readDirectConnectGatewayRoutes(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeDirectConnectGatewayRoutes"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeDirectConnectGatewayRoute(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeDirectConnectGatewayRoute(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("DirectConnectGatewayRouteSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) PublishDirectConnectRoute(d *schema.ResourceData) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	if d.HasChange("bgp_status") {
		if d.Get("bgp_status").(string) == "Unpublished" {
			unpublishCall, err := s.unpublishDirectConnectRouteCall(d)
			if err != nil {
				return err
			}
			apiProcess.PutCalls(unpublishCall)
		}
		if d.Get("bgp_status").(string) == "Published" {
			publishCall, err := s.publishDirectConnectRouteCall(d)
			if err != nil {
				return err
			}
			apiProcess.PutCalls(publishCall)
		}
	}

	return apiProcess.Run()
}

func (s *VpcService) publishDirectConnectRouteCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectGatewayRouteId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "PublishDirectConnectRoute",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.PublishDirectConnectRoute(call.param)
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

func (s *VpcService) unpublishDirectConnectRouteCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DirectConnectGatewayRouteId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "UnpublishDirectConnectRoute",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.UnpublishDirectConnectRoute(call.param)
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

func (s *VpcService) ReadAndSetDirectConnects(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "DirectConnectId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadDirectConnects(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "DirectConnectName",
		idFiled:     "DirectConnectId",
		targetField: "direct_connects",
		extra: map[string]SdkResponseMapping{
			"DirectConnectName": {
				Field:    "name",
				KeepAuto: true,
			},
			"DirectConnectId": {
				Field:    "id",
				KeepAuto: true,
			},
			"VpcNOCId": {
				Field: "vpc_noc_id",
			},
		},
	})
}

func (s *VpcService) ReadDirectConnects(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	return pageQueryWithNextToken(condition, "MaxResults", "NextToken", 100, func(condition map[string]interface{}) ([]interface{}, string, error) {
		conn := s.client.vpcconn
		action := "DescribeDirectConnects"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeDirectConnects(nil)
			if err != nil {
				return data, "", err
			}
		} else {
			resp, err = conn.DescribeDirectConnects(&condition)
			if err != nil {
				return data, "", err
			}
		}
		nextToken := (*resp)["NextToken"]
		results, err = getSdkValue("DirectConnectSet", *resp)
		if err != nil {
			return data, "", err
		}
		data = results.([]interface{})
		return data, indirectString(nextToken), err
	})
}

func (s *VpcService) createDirectConnectBfdConfigCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateBfdConfig",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateBfdConfig(call.param)
			if err != nil {
				return resp, err
			}

			bfdConfigId, _ := getSdkValue("BfdConfig.BfdConfigId", *resp)
			d.SetId(indirectString(bfdConfigId))

			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}

	return callback, err
}

func (s *VpcService) ReadAndSetDirectConnectBfdConfig(d *schema.ResourceData, r *schema.Resource) (err error) {
	resp, err := s.readDirectConnectBfdConfig(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, resp, nil)
	return nil
}

func (s *VpcService) readDirectConnectBfdConfig(d *schema.ResourceData, bfdId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if bfdId == "" {
		bfdId = d.Id()
	}
	req := map[string]interface{}{
		"BfdConfigId.1": bfdId,
	}
	results, err = s.readDirectConnectBfdConfigs(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("bfd config %s not exist ", bfdId)
	}
	return data, err
}

func (s *VpcService) readDirectConnectBfdConfigs(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeBfdConfig"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeBfdConfig(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeBfdConfig(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("BfdConfigSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) CreateDirectConnectBfdConfig(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.createDirectConnectBfdConfigCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) ModifyDirectConnectBfdConfig(d *schema.ResourceData, r *schema.Resource) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)

	call, err := s.modifyBfdConfigCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) modifyBfdConfigCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["BfdConfigId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyBfdConfig",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyBfdConfig(call.param)
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
	return callback, err
}

func (s *VpcService) RemoveDirectConnectBfdConfig(d *schema.ResourceData) (err error) {
	apiProcess := NewApiProcess(context.TODO(), d, s.client, false)
	call, err := s.removeDirectConnectBfdConfigCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) removeDirectConnectBfdConfigCall(d *schema.ResourceData) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"BfdConfigId": d.Id(),
	}
	callback = ApiCall{
		param:  &req,
		action: "DeleteBfdConfig",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteBfdConfig(call.param)
			if err != nil {
				return resp, err
			}
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.readDirectConnectBfdConfig(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading direct connect bfd config when delete %q, %s", d.Id(), callErr))
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
