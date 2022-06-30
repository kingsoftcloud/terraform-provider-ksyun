package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

type EipService struct {
	client *KsyunClient
}

func (s *EipService) ReadAddresses(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.eipconn
		action := "DescribeAddresses"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeAddresses(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeAddresses(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("AddressesSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *EipService) ReadAddress(d *schema.ResourceData, allocationId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if allocationId == "" {
		allocationId = d.Id()
	}
	req := map[string]interface{}{
		"AllocationId.1": allocationId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}
	results, err = s.ReadAddresses(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Address %s not exist ", allocationId)
	}
	return data, err
}

func (s *EipService) ReadAndSetAddress(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadAddress(d, "")
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
			err = mergeTagsData(d, &data, s.client, "eip")
			if err != nil {
				return resource.NonRetryableError(err)
			}
			SdkResponseAutoResourceData(d, r, data, chargeExtraForVpc(data))
			return nil
		}
	})
}

func (s *EipService) ReadAndSetAddresses(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "AllocationId",
			Type:    TransformWithN,
		},
		"project_id": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"network_interface_id": {
			mapping: "network-interface-id",
			Type:    TransformWithFilter,
		},
		"instance_type": {
			mapping: "instance-type",
			Type:    TransformWithFilter,
		},
		"internet_gateway_id": {
			mapping: "internet-gateway-id",
			Type:    TransformWithFilter,
		},
		"band_width_share_id": {
			mapping: "band-width-share-id",
			Type:    TransformWithFilter,
		},
		"line_id": {
			mapping: "line-id",
			Type:    TransformWithFilter,
		},
		"public_ip": {
			mapping: "public-ip",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadAddresses(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "AllocationId",
		targetField: "eips",
		extra: map[string]SdkResponseMapping{
			"AllocationId": {
				Field:    "id",
				KeepAuto: true,
			},
		},
	})
}

func (s *EipService) CreateAddressCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"tags": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AllocateAddress",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.eipconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AllocateAddress(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("AllocationId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *EipService) CreateAddress(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateAddressCall(d, r)
	if err != nil {
		return err
	}
	tagService := TagService{s.client}
	tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, r, "eip", false, true)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call, tagCall}, d, s.client, true)
}

func (s *EipService) ModifyAddressProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
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

func (s *EipService) ModifyAddressCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {Ignore: true},
		"tags":       {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		false,
	})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["AllocationId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyAddress",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.eipconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyAddress(call.param)
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

func (s *EipService) ModifyAddress(d *schema.ResourceData, r *schema.Resource) (err error) {
	projectCall, err := s.ModifyAddressProjectCall(d, r)
	if err != nil {
		return err
	}
	call, err := s.ModifyAddressCall(d, r)
	if err != nil {
		return err
	}
	tagService := TagService{s.client}
	tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, r, "eip", true, false)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{projectCall, call, tagCall}, d, s.client, true)
}

func (s *EipService) RemoveAddressCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"AllocationId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "ReleaseAddress",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.eipconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ReleaseAddress(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadAddress(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading address when delete %q, %s", d.Id(), callErr))
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

func (s *EipService) RemoveAddress(d *schema.ResourceData) (err error) {
	call, err := s.RemoveAddressCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *EipService) ReadAddressAssociate(d *schema.ResourceData, allocationId string, instanceId string, networkInterfaceId string) (data map[string]interface{}, err error) {
	data, err = s.ReadAddress(d, allocationId)
	if id, ok := data["InstanceId"]; ok {
		if id != instanceId {
			return data, fmt.Errorf("InstanceId %s not associate in Address %s ", instanceId, allocationId)
		}
	} else {
		return data, fmt.Errorf("InstanceId %s not associate in Address %s ", instanceId, allocationId)
	}
	if networkInterfaceId != "" {
		if vifId, ok := data["NetworkInterfaceId"]; ok {
			if vifId != networkInterfaceId {
				return data, fmt.Errorf("InstanceId %s not associate in Address %s ", instanceId, allocationId)
			}
		} else {
			return data, fmt.Errorf("InstanceId %s not associate in Address %s ", instanceId, allocationId)
		}
	}
	return data, err
}

func (s *EipService) ReadAndSetAddressAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadAddressAssociate(d, d.Get("allocation_id").(string), d.Get("instance_id").(string), d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *EipService) CreateAddressAssociateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AssociateAddress",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.eipconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AssociateAddress(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("allocation_id").(string) + ":" + d.Get("instance_id").(string) + ":" + d.Get("network_interface_id").(string))
			return err
		},
	}
	return callback, err
}

func (s *EipService) CreateAddressAssociate(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateAddressAssociateCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *EipService) RemoveAddressAssociateCall(d *schema.ResourceData, allocationId string, instanceId string, networkInterfaceId string) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"AllocationId": allocationId,
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DisassociateAddress",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.eipconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DisassociateAddress(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadAddressAssociate(d, allocationId, instanceId, networkInterfaceId)
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading address associate when delete %q, %s", d.Id(), callErr))
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

func (s *EipService) RemoveAddressAssociate(d *schema.ResourceData) (err error) {
	call, err := s.RemoveAddressAssociateCall(d, d.Get("allocation_id").(string), d.Get("instance_id").(string), d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *EipService) ReadLines() (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.eipconn
	action := "GetLines"
	logger.Debug(logger.ReqFormat, action, nil)
	resp, err = conn.GetLines(nil)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("LineSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *EipService) ReadAndSetLines(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLines()
	if err != nil {
		return err
	}
	var newData []interface{}
	if name, ok := d.GetOk("line_name"); ok {
		for _, line := range data {
			m := line.(map[string]interface{})
			if m["LineName"] == name {
				newData = append(newData, line)
			}
		}
	} else {
		newData = data
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  newData,
		nameField:   "LineName",
		idFiled:     "LineId",
		targetField: "lines",
	})
}
