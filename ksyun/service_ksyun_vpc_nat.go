package ksyun

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type Filter interface {
	Filler() map[string]string
	GetFilterParams() (map[string]interface{}, error)
}

type RouteFilter struct {
	// VpcId the id of vpc
	VpcId string `mapstructure:"vpc-id"`

	// InstanceId, the next route of the instance(kec)
	InstanceId string `mapstructure:"instance-id"`

	// DestCidrBlock destination cidr blocks
	DestCidrBlock string `mapstructure:"destination-cidr-block"`
}

type DescribeRoutesParam struct {
	RouteId ParamsList `mapstructure:"RouteId" type:"list"`

	Filter RouteFilter `mapstructure:"Filter" type:"filter"`
}

type BandwidthLimitFilter struct {
	// PrivateIp of the instance ip
	PrivateIpAddress string `mapstructure:"private-ip-address"`

	// NetworkInterfaceId, the network interface belong to instance
	NetworkInterfaceId string `mapstructure:"network-interface-id"`

	// instanceType destination cidr blocks
	InstanceType string `mapstructure:"instance-type"`
}

type DescribeNatRateLimitParam struct {
	NatId string `mapstructure:"NatId" type:"string"`

	// slice type must be []interface{}
	Filter BandwidthLimitFilter `mapstructure:"Filter" type:"filter"`

	MaxResults int    `mapstructure:"MaxResults" type:"string"`
	NextToken  string `mapstructure:"NextToken" type:"string"`
}

type CreateDnatParams struct {
	DnatName string `mapstructure:"DnatName"`

	NatId string `mapstructure:"NatId"`

	NatIp            string `mapstructure:"NatIp,omitempty"`
	PublicPort       string `mapstructure:"PublicPort,omitempty"`
	PrivateIpAddress string `mapstructure:"PrivateIpAddress,omitempty"`
	IpProtocol       string `mapstructure:"IpProtocol,omitempty"`
	PrivatePort      string `mapstructure:"PrivatePort,omitempty"`
	Description      string
}

type ModifyDnatParams struct {
	DnatName string `mapstructure:"DnatName"`

	NatId string `mapstructure:"NatId"`

	NatIp            string `mapstructure:"NatIp"`
	PublicPort       string `mapstructure:"PublicPort"`
	PrivateIpAddress string `mapstructure:"PrivateIpAddress"`
	IpProtocol       string `mapstructure:"IpProtocol"`
	PrivatePort      string `mapstructure:"PrivatePort"`
	Description      string `mapstructure:"Description"`
}

type DescribeDnatsParams struct {
	Filter  DnatsFilter   `mapstructure:"Filter,omitempty" type:"filter"`
	DnatIds []interface{} `mapstructure:"DnatId" type:"list"`
}

type DnatsFilter struct {
	NatIp              string `mapstructure:"nat-ip,omitempty" type:"string"`
	PublicPort         string `mapstructure:"public-port,omitempty" type:"string"`
	PrivateIpAddress   string `mapstructure:"private-ip-address,omitempty" type:"string"`
	IpProtocol         string `mapstructure:"ip-protocol,omitempty" type:"string"`
	NatId              string `mapstructure:"nat-id,omitempty" type:"string"`
	DnatName           string `mapstructure:"dnat-name,omitempty" type:"string"`
	NetworkInterfaceId string `mapstructure:"network-interface-id,omitempty" type:"string"`
}

type Dnat struct {
	NatId            string `mapstructure:"NatId"`
	DnatId           string `mapstructure:"DnatId"`
	DnatName         string `mapstructure:"DnatName"`
	IpProtocol       string `mapstructure:"IpProtocol"`
	NatIp            string `mapstructure:"NatIp"`
	PublicPort       string `mapstructure:"PublicPort"`
	PrivateIpAddress string `mapstructure:"PrivateIpAddress"`
	PrivatePort      string `mapstructure:"PrivatePort"`
	Description      string `mapstructure:"Description"`
	CreateTime       string `mapstructure:"CreateTime"`
}

type DescribeDnatsResponse struct {
	RequestId string `mapstructure:"RequestId"`
	DnatSet   []Dnat `mapstructure:"DnatSet"`
}

func NewDescribeDnatsResponse() DescribeDnatsResponse {
	return DescribeDnatsResponse{}
}

func NewDescribeDnatsParams() DescribeDnatsParams {
	return DescribeDnatsParams{
		DnatIds: []interface{}{},
	}
}

func (s *VpcService) CreateDnat(d *schema.ResourceData, r *schema.Resource) error {
	var apiProcess = NewApiProcess(context.Background(), d, s.client, true)
	call, err := s.CreateDnatCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) CreateDnatCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "CreateDnat",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				return conn.CreateDnat(call.param)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				id, err := getSdkValue("Dnat.DnatId", *resp)
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

func (s *VpcService) ModifyDnat(d *schema.ResourceData, r *schema.Resource) error {
	var apiProcess = NewApiProcess(context.Background(), d, s.client, true)
	call, err := s.ModifyDnatCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *VpcService) ModifyDnatCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["DnatId"] = d.Id()
		req["NatId"] = d.Get("nat_id")
		callback = ApiCall{
			param:  &req,
			action: "ModifyDnat",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				return conn.ModifyDnat(call.param)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return nil
			},
		}
	}
	return callback, err
}

func (s *VpcService) RemoveDnat(d *schema.ResourceData, r *schema.Resource) error {
	var apiProcess = NewApiProcess(context.Background(), d, s.client, true)
	call, err := s.DeleteDnatCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)

	return apiProcess.Run()
}

func (s *VpcService) DeleteDnatCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req := map[string]interface{}{
		"DnatId": d.Id(),
	}
	if len(req) > 0 {
		callback = ApiCall{
			param:  &req,
			action: "DeleteDnat",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error) {
				conn := client.vpcconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				return conn.DeleteDnat(call.param)
			},
			callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
				return resource.Retry(RetryTimeoutMinute, func() *resource.RetryError {
					params := NewDescribeDnatsParams()
					params.DnatIds = append(params.DnatIds, d.Id())
					_, err := s.DescribeDnats(params)
					if err != nil {
						if notFoundError(err) {
							return nil
						}
						return retryError(err)
					}

					_, err = call.executeCall(d, client, call)

					return retryError(err)
				})

			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return nil
			},
		}
	}
	return callback, err
}

func (s *VpcService) ReadAndSetDnats(d *schema.ResourceData, r *schema.Resource) error {
	param := NewDescribeDnatsParams()
	param.DnatIds = []interface{}{d.Id()}
	dnats, err := s.DescribeDnats(param)
	if err != nil {
		return err
	}
	dnat := dnats[0]

	dnatMap := dnat.(map[string]interface{})
	SdkResponseAutoResourceData(d, r, dnatMap, nil)
	return nil

}

func (s *VpcService) DescribeDnats(param DescribeDnatsParams) (dnats []interface{}, err error) {
	req := make(map[string]interface{})
	err = StructureConverter(param, &req)
	if err != nil {
		return nil, err
	}
	return pageQuery(req, "MaxResults", "NextToken", 200, 1, func(m map[string]interface{}) ([]interface{}, error) {
		conn := s.client.vpcconn
		action := "DescribeDnats"
		logger.Debug(logger.RespFormat, action, req)
		resp, err := conn.DescribeDnats(&req)
		if err != nil {
			return nil, err
		}
		dataRaw, err := getSdkValue("DnatSet", *resp)
		if err != nil {
			return nil, fmt.Errorf("not found any dnat: request detail: +%v", req)
		}
		dnats = dataRaw.([]interface{})
		return dnats, nil
	})
}
