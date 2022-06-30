package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

type KecService struct {
	client *KsyunClient
}

func (s *KecService) readAndSetKecInstance(d *schema.ResourceData, r *schema.Resource, disableSetTag ...bool) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.readKecInstance(d, "", false)
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading instane %q, %s", d.Id(), callErr))
			}
		} else {
			//InstanceConfigure
			SdkResponseAutoResourceData(d, r, data["InstanceConfigure"], nil)
			//InstanceState
			stateExtra := map[string]SdkResponseMapping{
				"Name": {
					Field: "instance_status",
				},
			}
			SdkResponseAutoResourceData(d, r, data["InstanceState"], stateExtra)
			extra := map[string]SdkResponseMapping{}
			if data["NetworkInterfaceSet"] != nil {
				//Primary network_interface
				for _, vif := range data["NetworkInterfaceSet"].([]interface{}) {
					if vif.(map[string]interface{})["NetworkInterfaceType"] == "primary" {
						extra := map[string]SdkResponseMapping{
							"SecurityGroupSet": {
								Field: "security_group_id",
								FieldRespFunc: func(i interface{}) interface{} {
									var result []interface{}
									for _, v := range i.([]interface{}) {
										result = append(result, v.(map[string]interface{})["SecurityGroupId"])
									}
									return result
								},
							},
						}
						SdkResponseAutoResourceData(d, r, vif, extra)
						//read dns info
						var networkInterface map[string]interface{}
						networkInterface, err = s.readKecNetworkInterface(d.Get("network_interface_id").(string))
						if err != nil {
							return resource.NonRetryableError(err)
						}
						for k := range networkInterface {
							if k == "DNS1" || k == "DNS2" {
								continue
							}
							delete(networkInterface, k)
						}
						extra = map[string]SdkResponseMapping{
							"DNS1": {
								Field: "dns1",
							},
							"DNS2": {
								Field: "dns2",
							},
						}
						SdkResponseAutoResourceData(d, r, networkInterface, extra)
						break
					}
				}

				//extension_network_interface
				extra["NetworkInterfaceSet"] = SdkResponseMapping{
					Field: "extension_network_interface",
					FieldRespFunc: func(i interface{}) interface{} {
						var result []interface{}
						for _, v := range i.([]interface{}) {
							if v.(map[string]interface{})["NetworkInterfaceType"] != "primary" {
								result = append(result, v)
							}
						}
						return result
					},
				}
			}
			extra["KeySet"] = SdkResponseMapping{
				Field: "key_id",
			}

			//tag
			if len(disableSetTag) == 0 || !disableSetTag[0] {
				err = mergeTagsData(d, &data, s.client, "instance")
				if err != nil {
					return resource.NonRetryableError(err)
				}
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			if v, ok := d.GetOk("force_reinstall_system"); ok {
				err = d.Set("force_reinstall_system", v)
			} else {
				err = d.Set("force_reinstall_system", false)
			}
			//control
			_ = d.Set("has_modify_system_disk", false)
			_ = d.Set("has_modify_password", false)
			_ = d.Set("has_modify_keys", false)
			return resource.NonRetryableError(err)
		}
	})
}

func (s *KecService) ReadAndSetKecInstances(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "InstanceId",
			Type:    TransformWithN,
		},
		"project_id": {
			Type: TransformWithN,
		},
		"vpc_id": {
			Type: TransformWithFilter,
		},
		"subnet_id": {
			Type: TransformWithFilter,
		},
		"network_interface": {
			Type: TransformListFilter,
		},
		"instance_state": {
			Type: TransformListFilter,
		},
		"availability_zone": {
			mappings: map[string]string{
				"availability_zone.name": "availability-zone-name",
			},
			Type: TransformListFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readKecInstances(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "InstanceId",
		nameField:   "InstanceName",
		targetField: "instances",
		extra: map[string]SdkResponseMapping{
			"KeySet": {
				Field: "key_id",
			},
		},
	})
}

func (s *KecService) readKecNetworkInterface(networkInterfaceId string) (data map[string]interface{}, err error) {
	var (
		networkInterfaces []interface{}
	)
	vpcService := VpcService{s.client}
	req := map[string]interface{}{
		"NetworkInterfaceId.1": networkInterfaceId,
		"Filter.1.Name":        "instance-type",
		"Filter.1.Value.1":     "kec",
	}
	networkInterfaces, err = vpcService.ReadNetworkInterfaces(req)
	if err != nil {
		return data, err
	}
	for _, v := range networkInterfaces {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Kec network interface %s not exist ", networkInterfaceId)
	}
	return data, err
}

func (s *KecService) readKecInstance(d *schema.ResourceData, instanceId string, allProject bool) (data map[string]interface{}, err error) {
	var (
		kecInstanceResults []interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"InstanceId.1": instanceId,
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

	kecInstanceResults, err = s.readKecInstances(req)
	if err != nil {
		return data, err
	}
	for _, v := range kecInstanceResults {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Kec instance %s not exist ", instanceId)
	}
	return data, err
}

func (s *KecService) readKecInstances(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	return pageQuery(condition, "MaxResults", "Marker", 200, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kecconn
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeInstances(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeInstances(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("InstancesSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func readKecNetworkInterfaces(d *schema.ResourceData, meta interface{}, condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp                    *map[string]interface{}
		networkInterfaceResults interface{}
	)
	conn := meta.(*KsyunClient).vpcconn
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

func (s *KecService) createKecInstance(d *schema.ResourceData, resource *schema.Resource) (err error) {
	var (
		callbacks []ApiCall
	)
	createCall, err := s.createKecInstanceCommon(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, createCall)
	tagService := TagService{s.client}
	tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resource, "instance", false, true)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, tagCall)
	dnsCall, err := s.initKecInstanceNetwork(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, dnsCall)
	// dryRun
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *KecService) modifyKecInstance(d *schema.ResourceData, resource *schema.Resource) (err error) {
	var (
		callbacks []ApiCall
	)
	//project
	projectCall, err := s.modifyKecInstanceProject(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, projectCall)
	//tag
	tagService := TagService{s.client}
	tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resource, "instance", true, false)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, tagCall)
	//name
	nameCall, err := s.modifyKecInstanceName(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, nameCall)
	//role
	roleCall, err := s.modifyKecInstanceIamRole(d)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, roleCall)
	//network update
	networkCall, err := s.modifyKecInstanceNetwork(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, networkCall)
	//force stop or start
	stateCall, err := s.stopOrStartKecInstance(d)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, stateCall)
	//need to stop
	//image
	imageCall, err := s.modifyKecInstanceImage(d, resource)
	if err != nil {
		return err
	}
	//password
	passCall, err := s.modifyKecInstancePassword(d, resource)
	if err != nil {
		return err
	}
	//key
	addCall, removeCall, err := s.modifyKecInstanceKeys(d)
	if err != nil {
		return err
	}
	if passCall.executeCall != nil || imageCall.executeCall != nil || addCall.executeCall != nil || removeCall.executeCall != nil {
		stopCall, err := s.stopKecInstance(d)
		if err != nil {
			return err
		}
		callbacks = append(callbacks, stopCall)
	}
	callbacks = append(callbacks, passCall)
	callbacks = append(callbacks, imageCall)
	callbacks = append(callbacks, addCall)
	callbacks = append(callbacks, removeCall)
	if passCall.executeCall != nil || imageCall.executeCall != nil || addCall.executeCall != nil || removeCall.executeCall != nil {
		startCall, err := s.startKecInstance(d)
		if err != nil {
			return err
		}
		callbacks = append(callbacks, startCall)
	}
	//need to restart
	specCall, err := s.modifyKecInstanceType(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, specCall)
	hostNameCall, err := s.modifyKecInstanceHostName(d, resource)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, hostNameCall)

	//if hostNameCall.executeCall != nil {
	//	stopCall, err := s.stopKecInstance(d)
	//	if err != nil {
	//		return err
	//	}
	//	callbacks = append(callbacks, stopCall)
	//	startCall, err := s.startKecInstance(d)
	//	if err != nil {
	//		return err
	//	}
	//	callbacks = append(callbacks, startCall)
	//}

	// 2022-03-17 [更配重启问题记录] by ydx
	// 先stop再start，有时候stop执行后机器没有关闭（默认不使用强制重启，避免客户在不知情的情况下影响服务）；
	// 如果卡在stop，用户使用其他方式重启了机器，stop就会一直retry
	// 这里改用reboot，然后等待active，如果没有成功，用户从控制台重启后也会变成active状态，retry到此状态就可以正常退出了
	if specCall.executeCall != nil || hostNameCall.executeCall != nil {
		rebootCall, err := s.rebootKecInstance(d)
		if err != nil {
			return err
		}
		callbacks = append(callbacks, rebootCall)
	}
	return ksyunApiCallNew(callbacks, d, s.client, true)
}

func (s *KecService) createKecInstanceCommon(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"key_id": {
			Type: TransformWithN,
		},
		"system_disk": {
			Type: TransformListUnique,
		},
		"security_group_id": {
			Type: TransformWithN,
		},
		"data_disks": {
			mappings: map[string]string{
				"data_disks": "DataDisk",
				"disk_size":  "Size",
				"disk_type":  "Type",
			}, Type: TransformListN,
		},
		"instance_status":        {Ignore: true},
		"force_delete":           {Ignore: true},
		"force_reinstall_system": {Ignore: true},
		"tags":                   {Ignore: true},
	}
	createReq, err := SdkRequestAutoMapping(d, resource, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	createReq["MaxCount"] = "1"
	createReq["MinCount"] = "1"

	callback = ApiCall{
		param:  &createReq,
		action: "RunInstances",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.RunInstances(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var (
				instanceId interface{}
			)
			if resp != nil {
				instanceId, err = getSdkValue("InstancesSet.0.InstanceId", *resp)
				if err != nil {
					return err
				}
				d.SetId(instanceId.(string))
			}
			err = s.checkKecInstanceState(d, "", []string{"active"}, d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return err
			}
			return s.readAndSetKecInstance(d, resource, true)
		},
	}
	return callback, err
}

func (s *KecService) modifyKecInstanceType(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"instance_type": {},
		"data_disk_gb":  {},
	}
	if d.HasChange("system_disk") && !d.Get("has_modify_system_disk").(bool) {
		transform["system_disk.0.disk_size"] = SdkReqTransform{
			mapping: "SystemDisk.DiskSize",
		}
		transform["system_disk.0.disk_type"] = SdkReqTransform{
			mapping:          "SystemDisk.DiskType",
			forceUpdateParam: true,
		}
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		updateReq["InstanceId"] = d.Id()
		callback = ApiCall{
			param:  &updateReq,
			action: "ModifyInstanceType",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyInstanceType(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = s.checkKecInstanceState(d, "", []string{"resize_success_local", "migrating_success", "migrating_success_off_line"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				return err
			},
		}
	}
	return callback, err
}
func (s *KecService) modifyKecInstanceIamRole(d *schema.ResourceData) (callback ApiCall, err error) {
	if d.HasChange("iam_role_name") {
		_, nr := d.GetChange("iam_role_name")
		if nr == "" {
			//unbind
			updateReq := map[string]interface{}{
				"InstanceId.1": d.Id(),
			}
			callback = ApiCall{
				param:  &updateReq,
				action: "DetachInstancesIamRole",
				executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
					conn := client.kecconn
					logger.Debug(logger.RespFormat, call.action, *(call.param))
					resp, err = conn.DetachInstancesIamRole(call.param)
					return resp, err
				},
				afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
					logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
					return err
				},
			}
		} else {
			//change
			updateReq := map[string]interface{}{
				"InstanceId.1": d.Id(),
				"IamRoleName":  nr,
			}
			callback = ApiCall{
				param:  &updateReq,
				action: "AttachInstancesIamRole",
				executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
					conn := client.kecconn
					logger.Debug(logger.RespFormat, call.action, *(call.param))
					resp, err = conn.AttachInstancesIamRole(call.param)
					return resp, err
				},
				afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
					logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
					return err
				},
			}
		}
	}
	return callback, err
}

func (s *KecService) modifyKecInstanceProject(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
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

func (s *KecService) modifyKecInstanceName(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"instance_name": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		updateReq["InstanceId"] = d.Id()
		callback = ApiCall{
			param:  &updateReq,
			action: "ModifyInstanceAttribute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyInstanceAttribute(call.param)
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

func (s *KecService) modifyKecInstanceHostName(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"host_name": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		updateReq["InstanceId"] = d.Id()
		callback = ApiCall{
			param:  &updateReq,
			action: "ModifyInstanceAttribute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyInstanceAttribute(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = s.checkKecInstanceState(d, "", []string{"active"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				return err
			},
		}
	}
	return callback, err
}

func (s *KecService) modifyKecInstancePassword(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"instance_password": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 && !d.Get("has_modify_password").(bool) {
		updateReq["InstanceId"] = d.Id()
		callback = ApiCall{
			param:  &updateReq,
			action: "ModifyInstanceAttribute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyInstanceAttribute(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return s.checkKecInstanceState(d, "", []string{"stopped"}, d.Timeout(schema.TimeoutUpdate))
			},
		}
	}
	return callback, err
}

func (s *KecService) updateKecInstanceNetwork(updateReq map[string]interface{}, resource *schema.Resource, init bool) (callback ApiCall, err error) {
	if len(updateReq) > 0 {
		callback = ApiCall{
			param:  &updateReq,
			action: "ModifyNetworkInterfaceAttribute",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				(*call.param)["InstanceId"] = d.Id()
				(*call.param)["NetworkInterfaceId"] = d.Get("network_interface_id")
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyNetworkInterfaceAttribute(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = s.checkKecInstanceState(d, "", []string{"active"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				if init {
					return s.readAndSetKecInstance(d, resource)
				}
				return err
			},
		}
	}
	return callback, err
}

func (s *KecService) modifyKecInstanceNetwork(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"security_group_id": {
			forceUpdateParam: true,
			Type:             TransformWithN,
		},
		"subnet_id":       {},
		"private_address": {},
		"dns1": {
			mapping: "DNS1",
		},
		"dns2": {
			mapping: "DNS2",
		},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	_, updateSubnet := updateReq["SubnetId"]
	_, updateIp := updateReq["PrivateAddress"]
	_, updateDns1 := updateReq["DNS1"]
	_, updateDns2 := updateReq["DNS2"]
	if updateSubnet || updateIp || updateDns1 || updateDns2 {
		return s.updateKecInstanceNetwork(updateReq, resource, false)
	}
	return callback, err
}

func (s *KecService) initKecInstanceNetwork(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"security_group_id": {
			Type: TransformWithN,
		},
		"dns1": {
			mapping: "DNS1",
		},
		"dns2": {
			mapping: "DNS2",
		},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, false, transform, nil)
	if err != nil {
		return callback, err
	}
	var init bool
	if _, ok := updateReq["DNS1"]; ok {
		init = true
	}
	if _, ok := updateReq["DNS2"]; ok {
		init = true
	}
	if init {
		callback, err = s.updateKecInstanceNetwork(updateReq, resource, true)
		callback.disableDryRun = true
		return callback, err
	}

	return callback, err
}

func (s *KecService) modifyKecInstanceImage(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"system_disk": {
			forceUpdateParam: true,
			Type:             TransformListUnique,
		},
		"key_id": {
			forceUpdateParam: true,
			Type:             TransformWithN,
		},
		"keep_image_login": {
			forceUpdateParam: true,
		},
		"instance_password": {
			forceUpdateParam: true,
		},
	}
	if d.HasChange("force_reinstall_system") && d.Get("force_reinstall_system").(bool) {
		transform["image_id"] = SdkReqTransform{forceUpdateParam: true}
	} else {
		transform["image_id"] = SdkReqTransform{}
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if _, ok := updateReq["ImageId"]; ok {
		updateReq["InstanceId"] = d.Id()
		err = d.Set("has_modify_system_disk", true)
		if err != nil {
			return callback, err
		}
		err = d.Set("has_modify_password", true)
		if err != nil {
			return callback, err
		}
		err = d.Set("has_modify_keys", true)
		if err != nil {
			return callback, err
		}
		callback = ApiCall{
			param:  &updateReq,
			action: "ModifyInstanceImage",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyInstanceImage(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				err = s.checkKecInstanceState(d, "", []string{"active"}, d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
				return err
			},
		}
	}
	return callback, err
}

func (s *KecService) processKeysChange(d *schema.ResourceData, keys []interface{}, isAdd bool) (callback ApiCall, err error) {
	if len(keys) > 0 {
		updateReq := map[string]interface{}{
			"InstanceId.1": d.Id(),
		}
		count := 1
		for _, key := range keys {
			updateReq["KeyId."+strconv.Itoa(count)] = key
			count++
		}
		var action string
		if isAdd {
			action = "AttachKey"
		} else {
			action = "DetachKey"
		}
		callback = ApiCall{
			param:  &updateReq,
			action: action,
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kecconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				if call.action == "AttachKey" {
					resp, err = conn.AttachKey(call.param)
				} else {
					resp, err = conn.DetachKey(call.param)
				}
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return s.checkKecInstanceState(d, "", []string{"stopped"}, d.Timeout(schema.TimeoutUpdate))
			},
		}
	}
	return callback, err
}

func (s *KecService) modifyKecInstanceKeys(d *schema.ResourceData) (add ApiCall, remove ApiCall, err error) {
	if d.HasChange("key_id") && !d.Get("has_modify_keys").(bool) {
		oldK, newK := d.GetChange("key_id")
		removeKeys := oldK.(*schema.Set).Difference(newK.(*schema.Set)).List()
		newKeys := newK.(*schema.Set).Difference(oldK.(*schema.Set)).List()
		remove, err = s.processKeysChange(d, removeKeys, false)
		if err != nil {
			return add, remove, err
		}
		add, err = s.processKeysChange(d, newKeys, true)
		if err != nil {
			return add, remove, err
		}
	}
	return add, remove, err
}

func (s *KecService) stopOrStartKecInstance(d *schema.ResourceData) (callback ApiCall, err error) {
	if d.HasChange("instance_status") {
		if d.Get("instance_status") == "active" {
			return s.startKecInstance(d)
		} else {
			return s.stopKecInstance(d)
		}
	}
	return callback, err
}

func (s *KecService) rebootKecInstance(d *schema.ResourceData) (callback ApiCall, err error) {
	updateReq := map[string]interface{}{
		"InstanceId.1": d.Id(),
	}
	callback = ApiCall{
		param:  &updateReq,
		action: "RebootInstances",
		beforeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (doExecute bool, err error) {
			data, err := s.readKecInstance(d, "", false)
			if err != nil {
				return doExecute, err
			}
			status, err := getSdkValue("InstanceState.Name", data)
			if err != nil {
				return doExecute, err
			}
			if status.(string) == "stopped" {
				doExecute = false
			} else {
				doExecute = true
			}
			return doExecute, err
		},
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.RebootInstances(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			err = s.checkKecInstanceState(d, "", []string{"active"}, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			return err
		},
	}
	return callback, err
}

func (s *KecService) stopKecInstance(d *schema.ResourceData) (callback ApiCall, err error) {
	updateReq := map[string]interface{}{
		"InstanceId.1": d.Id(),
	}
	callback = ApiCall{
		param:  &updateReq,
		action: "StopInstances",
		beforeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (doExecute bool, err error) {
			data, err := s.readKecInstance(d, "", false)
			if err != nil {
				return doExecute, err
			}
			status, err := getSdkValue("InstanceState.Name", data)
			if err != nil {
				return doExecute, err
			}
			if status.(string) == "stopped" {
				doExecute = false
			} else {
				doExecute = true
			}
			return doExecute, err
		},
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.StopInstances(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			err = s.checkKecInstanceState(d, "", []string{"stopped"}, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			return err
		},
	}
	return callback, err
}

func (s *KecService) startKecInstance(d *schema.ResourceData) (callback ApiCall, err error) {
	updateReq := map[string]interface{}{
		"InstanceId.1": d.Id(),
	}
	callback = ApiCall{
		param:  &updateReq,
		action: "StartInstances",
		beforeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (doExecute bool, err error) {
			data, err := s.readKecInstance(d, "", false)
			if err != nil {
				return doExecute, err
			}
			status, err := getSdkValue("InstanceState.Name", data)
			if err != nil {
				return doExecute, err
			}
			if status.(string) == "active" {
				doExecute = false
			} else {
				doExecute = true
			}
			return doExecute, err
		},
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.StartInstances(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			err = s.checkKecInstanceState(d, "", []string{"active"}, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			return err
		},
	}
	return callback, err
}

func (s *KecService) removeKecInstance(d *schema.ResourceData, meta interface{}) (err error) {
	conn := meta.(*KsyunClient).kecconn
	req := make(map[string]interface{})
	req["InstanceId.1"] = d.Id()
	req["ForceDelete"] = true
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		action := "TerminateInstances"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.TerminateInstances(&req)
		if err == nil {
			return nil
		}
		_, err = s.readKecInstance(d, "", false)
		if err != nil {
			if notFoundError(err) {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading instance when delete %q, %s", d.Id(), err))
			}
		}
		return nil
	})
}

func (s *KecService) checkKecInstanceState(d *schema.ResourceData, instanceId string, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{},
		Target:       target,
		Refresh:      s.kecInstanceStateRefreshFunc(d, instanceId, []string{"error"}),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
		MinTimeout:   1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *KecService) kecInstanceStateRefreshFunc(d *schema.ResourceData, instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := s.readKecInstance(d, instanceId, true)
		if err != nil {
			return nil, "", err
		}

		status, err := getSdkValue("InstanceState.Name", data)
		if err != nil {
			return nil, "", err
		}

		for _, v := range failStates {
			if v == status.(string) {
				return nil, "", fmt.Errorf("instance status  error, status:%v", status)
			}
		}
		return data, status.(string), nil
	}
}

func (s *KecService) createNetworkInterfaceCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	vpcService := VpcService{s.client}
	data, err := vpcService.ReadSubnet(d, d.Get("subnet_id").(string))
	if err != nil {
		return callback, err
	}
	if data["SubnetType"] != "Normal" {
		return callback, fmt.Errorf("Subnet type %s not support for kec network interface ", data["SubnetType"].(string))
	}
	transform := map[string]SdkReqTransform{
		"security_group_ids": {
			mapping: "SecurityGroupId",
			Type:    TransformWithN,
		},
	}
	createReq, err := SdkRequestAutoMapping(d, resource, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}

	return vpcService.CreateNetworkInterfaceCall(&createReq)
}

func (s *KecService) createNetworkInterface(d *schema.ResourceData, resource *schema.Resource) (err error) {
	call, err := s.createNetworkInterfaceCall(d, resource)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KecService) readAndSetNetworkInterface(d *schema.ResourceData, resource *schema.Resource) (err error) {
	vpcService := VpcService{s.client}
	data, err := vpcService.ReadNetworkInterface(d, "")
	if err != nil {
		return err
	}
	if data["InstanceType"] != "kec" {
		return fmt.Errorf("Network interface type %s not support for kec ", data["InstanceType"].(string))
	}
	extra := map[string]SdkResponseMapping{
		"SecurityGroupSet": {
			Field: "security_group_ids",
			FieldRespFunc: func(i interface{}) interface{} {
				var sgIds []string
				for _, v := range i.([]interface{}) {
					sgIds = append(sgIds, v.(map[string]interface{})["SecurityGroupId"].(string))
				}
				return sgIds
			},
		},
	}
	SdkResponseAutoResourceData(d, resource, data, extra)
	return err
}

func (s *KecService) modifyNetworkInterfaceAttrCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	if d.HasChange("subnet_id") || d.HasChange("private_ip_address") || d.HasChange("security_group_ids") {
		transform := map[string]SdkReqTransform{
			"subnet_id": {
				forceUpdateParam: true,
			},
			"security_group_ids": {
				forceUpdateParam: true,
				mapping:          "SecurityGroupId",
				Type:             TransformWithN,
			},
			"private_ip_address": {},
		}
		updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
		if err != nil {
			return callback, err
		}
		if len(updateReq) > 0 {
			updateReq["NetworkInterfaceId"] = d.Id()
			updateReq["InstanceId"] = d.Get("instance_id")
			callback = ApiCall{
				param:  &updateReq,
				action: "ModifyNetworkInterfaceAttribute",
				executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
					conn := client.kecconn
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
	}
	return callback, err
}

func (s *KecService) modifyNetworkInterfaceCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"network_interface_name": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		vpcService := VpcService{s.client}
		updateReq["NetworkInterfaceId"] = d.Id()
		return vpcService.ModifyNetworkInterfaceCall(&updateReq)
	}
	return callback, err
}

func (s *KecService) modifyNetworkInterface(d *schema.ResourceData, resource *schema.Resource) (err error) {
	var calls []ApiCall
	call, err := s.modifyNetworkInterfaceCall(d, resource)
	if err != nil {
		return err
	}
	calls = append(calls, call)
	if d.Get("instance_id") != "" {
		var (
			attrCall ApiCall
		)
		attrCall, err = s.modifyNetworkInterfaceAttrCall(d, resource)
		if err != nil {
			return err
		}
		calls = append(calls, attrCall)
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KecService) readAndSetNetworkInterfaceAttachment(d *schema.ResourceData, resource *schema.Resource) (err error) {
	vpcService := VpcService{s.client}
	data, err := vpcService.ReadNetworkInterface(d, d.Get("network_interface_id").(string))
	if err != nil {
		return err
	}
	if data["InstanceType"] != "kec" {
		return fmt.Errorf("Network interface instance type %s not support for kec ", data["InstanceType"].(string))
	}
	if data["NetworkInterfaceType"] != "extension" {
		return fmt.Errorf("Network interface type %s not support for kec network interface attachment ", data["NetworkInterfaceType"].(string))
	}
	if id, ok := data["InstanceId"]; ok {
		if id != d.Get("instance_id") {
			return fmt.Errorf("Network interface attachmemt %s not exist ", d.Id())
		}
	} else {
		return fmt.Errorf("Network interface attachmemt %s not exist ", d.Id())
	}
	SdkResponseAutoResourceData(d, resource, data, nil)
	return err
}

func (s *KecService) createNetworkInterfaceAttachmentCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	createReq, err := SdkRequestAutoMapping(d, resource, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &createReq,
		action: "AttachNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AttachNetworkInterface(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("network_interface_id").(string) + ":" + d.Get("instance_id").(string))
			return s.checkKecInstanceState(d, d.Get("instance_id").(string), []string{"active", "stopped"}, d.Timeout(schema.TimeoutUpdate))
		},
	}
	return callback, err
}

func (s *KecService) createNetworkInterfaceAttachment(d *schema.ResourceData, resource *schema.Resource) (err error) {
	call, err := s.createNetworkInterfaceAttachmentCall(d, resource)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KecService) modifyNetworkInterfaceAttachmentCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	createReq, err := SdkRequestAutoMapping(d, resource, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &createReq,
		action: "DetachNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kecconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DetachNetworkInterface(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *KecService) modifyNetworkInterfaceAttachment(d *schema.ResourceData, resource *schema.Resource) (err error) {
	call, err := s.modifyNetworkInterfaceAttachmentCall(d, resource)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
