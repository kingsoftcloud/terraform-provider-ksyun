package ksyun

import (
	"context"

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
	return nil
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
