package ksyun

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

var (
	// vpnV1Attribute the flowing fields are invalid when vpn1.0
	vpnV2Attribute = []string{"ha_mode", "open_health_check", "local_peer_ip", "customer_peer_ip", "ike_version"}

	// vpnV1Attribute the flowing fields are invalid when vpn2.0
	vpnV1Attribute = []string{"vpn_gre_ip", "ha_vpn_gre_ip", "customer_gre_ip", "ha_customer_gre_ip"}
)

type VpnSrv struct {
	client *KsyunClient
}

func NewVpnSrv(client *KsyunClient) VpnSrv {
	return VpnSrv{
		client: client,
	}
}

func (v *VpnSrv) CreateVpnGatewayRoute(d *schema.ResourceData, r *schema.Resource) error {
	ap := NewApiProcess(context.Background(), d, v.client, true)
	createCall, err := v.CreateVpnGatewayRouteCall(d, r)
	if err != nil {
		return err
	}
	ap.PutCalls(createCall)

	return ap.Run()
}

func (v *VpnSrv) CreateVpnGatewayRouteCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "CreateVpnGatewayRoute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				return conn.CreateVpnGatewayRoute(call.param)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				id, err := getSdkValue("RouteId", *resp)
				if err != nil {
					return err
				}
				d.SetId(id.(string))
				return nil
			},
		}
	}
	return callback, err
}

func (v *VpnSrv) ReadAndSetVpnGatewayRoute(d *schema.ResourceData, r *schema.Resource) error {
	data, err := v.ReadVpnGatewayRoute(d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	var currResult []interface{}
	for _, routeIf := range data {
		routeId, err := getSdkValue("VpnGatewayRouteId", routeIf)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(routeId, d.Id()) {
			currResult = append(currResult, routeIf)
		}
	}

	SdkResponseAutoResourceData(d, r, currResult, nil)
	return nil
}

func (v *VpnSrv) ReadVpnGatewayRoute(vpnGatewayId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)

	req := map[string]interface{}{
		"VpnGatewayId": vpnGatewayId,
	}
	results, err = v.DescribeVpnGatewayRoutes(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("vpn gateway route %s not exist ", vpnGatewayId)
	}
	return data, err
}

func (v *VpnSrv) DescribeVpnGatewayRoutes(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := v.client.vpcconn
		action := "DescribeVpnGatewayRoutes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeVpnGatewayRoutes(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeVpnGatewayRoutes(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("VpnGatewayRouteSet", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			return []interface{}{}, nil
		}
		data = results.([]interface{})
		return data, err
	})
}

func (v *VpnSrv) RemoveVpnGatewayRoute(d *schema.ResourceData) error {
	ap := NewApiProcess(context.Background(), d, v.client, true)

	removeCall, err := v.RemoveVpnGatewayRouteCall(d)
	if err != nil {
		return err
	}
	ap.PutCalls(removeCall)
	return ap.Run()
}

func (v *VpnSrv) RemoveVpnGatewayRouteCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeParams := map[string]interface{}{
		"VpnGatewayRouteId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeParams,
		action: "DeleteVpnGatewayRoute",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err := conn.DeleteVpnGatewayRoute(call.param)
			return resp, err
		},
	}

	return callback, err
}

func (v *VpnSrv) ReadAndSetVpnGatewayRoutes(d *schema.ResourceData, r *schema.Resource) error {
	transform := map[string]SdkReqTransform{
		"vpn_gateway_id": {
			mapping: "VpnGatewayId",
			Type:    TransformDefault,
		},
		"next_hop_types": {
			mapping: "nexthop",
			Type:    TransformWithFilter,
		},
		"cidr_blocks": {
			mapping: "cidr-block",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := v.DescribeVpnGatewayRoutes(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "VpnGatewayRouteId",
		targetField: "vpn_gateway_routes",
		extra:       nil,
	})
}
