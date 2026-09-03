package ksyun

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type InstanceModelService struct {
	client *KsyunClient
}

func (s *InstanceModelService) ReadModels(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	return pageQuery(condition, "MaxResults", "Marker", 50, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kecconn
		action := "DescribeModels"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeModels(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeModels(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("ModelSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *InstanceModelService) ReadModel(d *schema.ResourceData, modelId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if modelId == "" {
		modelId = d.Id()
	}
	req := map[string]interface{}{
		"ModelId.1": modelId,
	}

	results, err = s.ReadModels(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("InstanceModel %s not exist ", modelId)
	}
	return data, err
}

func (s *InstanceModelService) ReadAndSetModel(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadModel(d, "")
		if callErr != nil {
			if notFoundError(callErr) {
				d.SetId("")
				return nil
			}
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			return resource.NonRetryableError(fmt.Errorf("error on reading InstanceModel %q, %s", d.Id(), callErr))
		} else {
			extra := map[string]SdkResponseMapping{
				"ModelId":           {Field: "model_id"},
				"ModelName":         {Field: "model_name"},
				"CreateTime":        {Field: "create_time"},
				"ImageId":           {Field: "image_id"},
				"InstanceType":      {Field: "instance_type"},
				"SubnetId":          {Field: "subnet_id"},
				"ChargeType":        {Field: "charge_type"},
				"PurchaseTime":      {Field: "purchase_time"},
				"InstanceName":      {Field: "instance_name"},
				"DataDiskGb":        {Field: "data_disk_gb"},
				"ProjectId":         {Field: "project_id"},
				"IsDistributeIpv6":  {Field: "is_distribute_ipv6"},
				"SriovNetSupport":   {Field: "sriov_net_support"},
				"KeepImageLogin":    {Field: "keep_image_login"},
				"FailureAutoDelete": {Field: "failure_auto_delete"},
				"SyncTag":           {Field: "sync_tag"},
				"SecurityGroupId": {
					Field: "security_group_id",
					FieldRespFunc: func(i interface{}) interface{} {
						if i == nil {
							return []interface{}{}
						}
						if str, ok := i.(string); ok {
							if str != "" {
								return []interface{}{str}
							}
							return []interface{}{}
						}
						if list, ok := i.([]interface{}); ok {
							return list
						}
						return []interface{}{}
					},
				},
				"SystemDisk": {
					Field: "system_disk",
					FieldRespFunc: func(i interface{}) interface{} {
						if i == nil {
							return []interface{}{}
						}
						if m, ok := i.(map[string]interface{}); ok {
							item := make(map[string]interface{})
							if diskType, ok := m["DiskType"]; ok {
								item["disk_type"] = diskType
							}
							if diskSize, ok := m["DiskSize"]; ok {
								item["disk_size"] = diskSize
							}
							return []interface{}{item}
						}
						return []interface{}{}
					},
				},
				"DataDiskSet": {
					Field: "data_disks",
					FieldRespFunc: func(i interface{}) interface{} {
						var result []interface{}
						if i == nil {
							return result
						}
						if list, ok := i.([]interface{}); ok {
							for _, v := range list {
								if m, ok := v.(map[string]interface{}); ok {
									item := make(map[string]interface{})
									if diskType, ok := m["Type"]; ok {
										item["disk_type"] = diskType
									}
									if diskSize, ok := m["Size"]; ok {
										item["disk_size"] = diskSize
									}
									if delWithInst, ok := m["DeleteWithInstance"]; ok {
										item["delete_with_instance"] = delWithInst
									}
									if snapshotId, ok := m["SnapshotId"]; ok {
										item["disk_snapshot_id"] = snapshotId
									}
									if snapshotName, ok := m["SnapshotName"]; ok {
										item["snapshot_name"] = snapshotName
									}
									result = append(result, item)
								}
							}
						}
						return result
					},
				},
				"NetworkInterfaceSet": {
					Field: "network_interface",
					FieldRespFunc: func(i interface{}) interface{} {
						var result []interface{}
						if i == nil {
							return result
						}
						if list, ok := i.([]interface{}); ok {
							for _, v := range list {
								if m, ok := v.(map[string]interface{}); ok {
									item := make(map[string]interface{})
									if subnetId, ok := m["SubnetId"]; ok {
										item["subnet_id"] = subnetId
									}
									if sgIds, ok := m["SecurityGroupIdSet"]; ok {
										var sgList []interface{}
										if sgListRaw, ok := sgIds.([]interface{}); ok {
											for _, sg := range sgListRaw {
												if sgMap, ok := sg.(map[string]interface{}); ok {
													if sgId, ok := sgMap["SecurityGroupId"]; ok {
														sgList = append(sgList, sgId)
													}
												}
											}
										}
										item["security_group_id"] = sgList
									}
									if privateIp, ok := m["PrivateIpAddress"]; ok {
										item["private_ip_address"] = privateIp
									}
									result = append(result, item)
								}
							}
						}
						return result
					},
				},
				"TagSet": {
					Field: "tags",
					FieldRespFunc: func(i interface{}) interface{} {
						var result []interface{}
						if i == nil {
							return result
						}
						if list, ok := i.([]interface{}); ok {
							for _, v := range list {
								if m, ok := v.(map[string]interface{}); ok {
									item := make(map[string]interface{})
									if key, ok := m["Key"]; ok {
										item["key"] = key
									}
									if value, ok := m["Value"]; ok {
										item["value"] = value
									}
									if id, ok := m["Id"]; ok {
										item["id"] = id
									}
									result = append(result, item)
								}
							}
						}
						return result
					},
				},
			}

			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (s *InstanceModelService) CreateModelCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"security_group_id": {
			mapping: "SecurityGroupId",
			Type:    TransformWithN,
		},
		"system_disk": {
			Type: TransformListUnique,
		},
		"data_disks": {
			mappings: map[string]string{
				"data_disks":           "DataDisk",
				"disk_size":            "Size",
				"disk_type":            "Type",
				"delete_with_instance": "DeleteWithInstance",
				"disk_snapshot_id":     "SnapshotId",
				"snapshot_name":        "SnapshotName",
			},
			Type: TransformListN,
		},
		"network_interface": {
			mappings: map[string]string{
				"network_interface":  "NetworkInterface",
				"subnet_id":          "SubnetId",
				"security_group_id":  "SecurityGroupId",
				"private_ip_address": "PrivateIpAddress",
			},
			Type: TransformListN,
		},
		"address_bandwidth": {
			mapping: "AddressBandWidth",
		},
		"is_distribute_ipv6": {
			mapping: "IsDistributeIpv6",
		},
		"sync_tag": {
			mapping: "SyncTag",
		},
		"model_id":    {Ignore: true},
		"create_time": {Ignore: true},
		"tags":        {Ignore: true},
	}

	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}

	// Convert address_charge_type to API value
	addressChargeTypeMap := map[string]string{
		"Monthly":          "PrePaidByMonth",
		"Peak":             "PostPaidByPeak",
		"Daily":            "PostPaidByDay",
		"TrafficMonthly":   "PostPaidByTransfer",
		"HourlySettlement": "PostPaidByHour",
		"RegionPeak":       "PostPaidByRegionPeak",
	}
	if v, ok := req["AddressChargeType"]; ok {
		if s, ok := v.(string); ok {
			if mapped, ok := addressChargeTypeMap[s]; ok {
				req["AddressChargeType"] = mapped
			}
		}
	}

	// Handle tags
	if tags, ok := d.GetOk("tags"); ok {
		tagsList := tags.([]interface{})
		for idx, tag := range tagsList {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				req["Tag."+strconv.Itoa(idx+1)+".Key"] = tagMap["key"]
				req["Tag."+strconv.Itoa(idx+1)+".Value"] = tagMap["value"]
				if tagId, ok := tagMap["id"]; ok && tagId.(string) != "" {
					req["Tag."+strconv.Itoa(idx+1)+".Id"] = tagId
				}
			}
		}
	}

	// Handle sync_tag - GetOk returns false for bool zero value, so set it manually
	req["SyncTag"] = d.Get("sync_tag")

	callback = ApiCall{
		param:  &req,
		action: "CreateModel",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.ReqFormat, call.action, *(call.param))
			resp, err = conn.CreateModel(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			modelId, err := getSdkValue("ModelId", *resp)
			if err != nil {
				return err
			}
			d.SetId(modelId.(string))
			return err
		},
	}
	return callback, err
}

func (s *InstanceModelService) CreateModel(d *schema.ResourceData, r *schema.Resource) (err error) {
	// Validate purchase_time when charge_type is Monthly
	if d.Get("charge_type").(string) == "Monthly" {
		if _, ok := d.GetOk("purchase_time"); !ok {
			return fmt.Errorf("purchase_time is required when charge_type is Monthly")
		}
	}
	// Validate address_purchase_time when address_charge_type is Monthly
	if d.Get("address_charge_type").(string) == "Monthly" {
		if _, ok := d.GetOk("address_purchase_time"); !ok {
			return fmt.Errorf("address_purchase_time is required when address_charge_type is Monthly")
		}
	}
	// Validate keep_image_login and key_id are mutually exclusive
	if d.Get("keep_image_login").(bool) {
		if _, ok := d.GetOk("key_id"); ok {
			return fmt.Errorf("keep_image_login and key_id are mutually exclusive, key_id must be empty when keep_image_login is true")
		}
	}
	if _, ok := d.GetOk("key_id"); ok {
		if d.Get("keep_image_login").(bool) {
			return fmt.Errorf("keep_image_login and key_id are mutually exclusive, keep_image_login must be false when key_id is set")
		}
	}

	call, err := s.CreateModelCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *InstanceModelService) RemoveModelCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"ModelId.1": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "TerminateModels",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.ReqFormat, call.action, *(call.param))
			resp, err = conn.TerminateModels(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(1*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadModel(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading InstanceModel when delete %q, %s", d.Id(), callErr))
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
			return nil
		},
	}
	return callback, err
}

func (s *InstanceModelService) RemoveModel(d *schema.ResourceData) (err error) {
	call, err := s.RemoveModelCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
