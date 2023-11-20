package ksyun

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type DnsService struct {
	client *KsyunClient
}

// CreatePrivateDnsZone creates private dns-2.0
func (s *DnsService) CreatePrivateDnsZone(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	createPdnsCall, err := s.createPdnsZoneWithCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createPdnsCall)

	return apiProcess.Run()
}

func (s *DnsService) ModifyPrivateDnsZone(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	modifyPdnsCall, err := s.modifyPdnsZoneWithCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(modifyPdnsCall)

	return apiProcess.Run()
}

func (s *DnsService) DeletePrivateDnsZone(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	removeCall, err := s.deletePdnsZoneWithCall(d)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(removeCall)

	return apiProcess.Run()
}

func (s *DnsService) ReadAndSetPrivateDnsZones(d *schema.ResourceData, r *schema.Resource) (err error) {
	var (
		data []interface{}
		req  map[string]interface{}
	)

	req = make(map[string]interface{})
	if zoneIdsIface, ok := d.GetOk("zone_ids"); ok {
		zoneIds := zoneIdsIface.(*schema.Set)
		for idx, zoneId := range zoneIds.List() {
			k := "Filter." + strconv.Itoa(idx+1)
			req[k] = zoneId
		}
	}

	data, err = s.ReadPrivateDnsZonePlural(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "ZoneId",
		targetField: "zones",
		extra:       nil,
	})
}

func (s *DnsService) ReadAndSetPrivateDnsZone(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadPrivateDnsZone(d, d.Id())
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading PrivateDnsZone %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, nil)
			return nil
		}
	})
}

func (s *DnsService) ReadPrivateDnsZonePlural(condition map[string]interface{}) (data []interface{}, err error) {

	data, err = pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		var (
			resp    *map[string]interface{}
			results interface{}
		)

		conn := s.client.pdnsconn
		action := "DescribePdnsZones"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil || len(condition) == 0 {
			resp, err = conn.DescribePdnsZones(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribePdnsZones(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("ZoneSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})

	return data, err
}

func (s *DnsService) ReadPrivateDnsZone(d *schema.ResourceData, zoneId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if zoneId == "" {
		zoneId = d.Id()
	}
	req := map[string]interface{}{
		"Filter.1": zoneId,
	}
	// err = addProjectInfo(d, &req, s.client)
	// if err != nil {
	// 	return data, err
	// }
	results, err = s.ReadPrivateDnsZonePlural(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("PrivateDnsZone %s not exist ", zoneId)
	}
	return data, err
}

func (s *DnsService) createPdnsZoneWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	params, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}

	callback = ApiCall{
		param:  &params,
		action: "CreatePdnsZone",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreatePdnsZone(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("ZoneVpc.ZoneId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *DnsService) modifyPdnsZoneWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	trans := map[string]SdkReqTransform{
		"zone_ttl": {},
	}
	params, err := SdkRequestAutoMapping(d, r, true, trans, nil)
	if err != nil {
		return callback, err
	}

	if params == nil || len(params) < 1 {
		return callback, err
	}

	params["ZoneId"] = d.Id()

	callback = ApiCall{
		param:  &params,
		action: "ModifyPdnsZone",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyPdnsZone(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *DnsService) deletePdnsZoneWithCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"ZoneId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeletePdnsZone",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeletePdnsZone(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadPrivateDnsZone(d, d.Id())
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading ReadPdnsZone when delete %q, %s", d.Id(), callErr))
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

func (s *DnsService) createZoneRecordWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	params, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}

	callback = ApiCall{
		param:  &params,
		action: "CreateZoneRecord",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateZoneRecord(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Record.RecordId", *resp)
			if err != nil {
				return err
			}
			zoneId := d.Get("zone_id").(string)
			d.SetId(AssembleIds(zoneId, id.(string)))
			return err
		},
	}
	return callback, err
}

func (s *DnsService) modifyZoneRecordWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	trans := map[string]SdkReqTransform{
		"record_ttl": {},
	}
	params, err := SdkRequestAutoMapping(d, r, true, trans, nil)
	if err != nil {
		return callback, err
	}

	if params == nil || len(params) < 1 {
		return callback, err
	}
	fullId := DisassembleIds(d.Id())
	params["RecordId"] = fullId[1]
	params["ZoneId"] = fullId[0]

	callback = ApiCall{
		param:  &params,
		action: "ModifyZoneRecord",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyZoneRecord(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *DnsService) deleteZoneRecordWithCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	fullId := DisassembleIds(d.Id())
	recordId := fullId[1]
	zoneId := fullId[0]
	req := map[string]interface{}{}

	trans := map[string]SdkReqTransform{
		"record_value": {
			forceUpdateParam: true,
		},
		"priority": {
			forceUpdateParam: true,
		},
		"weight": {
			forceUpdateParam: true,
		},
		"port": {
			forceUpdateParam: true,
		},
	}

	req, err = SdkRequestAutoMapping(d, r, false, trans, nil)
	if err != nil {
		return callback, err
	}
	req["RecordId"] = recordId
	req["ZoneId"] = zoneId

	callback = ApiCall{
		param:  &req,
		action: "DeleteZoneRecord",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteZoneRecord(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(3*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadPrivateDnsRecord(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading ReadPdnsRecord when delete %q, %s", d.Id(), callErr))
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

func (s *DnsService) ReadAndSetPrivateDnsRecords(d *schema.ResourceData, r *schema.Resource) error {
	var (
		data []interface{}
		err  error
		req  = make(map[string]interface{})
	)
	trans := map[string]SdkReqTransform{
		"record_ids": {
			Type:    TransformWithN,
			mapping: "RecordId",
		},

		"region_name": {
			Type:    TransformWithN,
			mapping: "RegionName",
		},
	}

	req, err = mergeDataSourcesReq(d, r, trans)
	if err != nil {
		return err
	}

	data, err = s.ReadPrivateDnsRecords(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "RecordId",
		targetField: "records",
	})

}

func (s *DnsService) ReadPrivateDnsRecord(d *schema.ResourceData, recordFullId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	zoneId := DisassembleIds(d.Id())[0]
	recordId := DisassembleIds(d.Id())[1]
	if recordFullId != "" {
		zoneId = DisassembleIds(recordFullId)[0]
		recordId = DisassembleIds(recordFullId)[1]
	}
	req := map[string]interface{}{
		"RecordId.1": recordId,
		"ZoneId":     zoneId,
	}
	// err = addProjectInfo(d, &req, s.client)
	// if err != nil {
	// 	return data, err
	// }
	results, err = s.ReadPrivateDnsRecords(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("PrivateDnsRecord %s not exist ", zoneId)
	}
	return data, err
}

func (s *DnsService) ReadPrivateDnsRecords(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.pdnsconn
	action := "DescribeZoneRecord"
	logger.Debug(logger.ReqFormat, action, condition)

	resp, err = conn.DescribeZoneRecord(&condition)
	if err != nil {
		return data, err
	}

	results, err = getSdkValue("RecordSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err

}

func (s *DnsService) CreatePrivateDnsRecord(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	createRecordCall, err := s.createZoneRecordWithCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createRecordCall)

	return apiProcess.Run()
}

func (s *DnsService) ModifyPrivateDnsRecord(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	modifyRecordCall, err := s.modifyZoneRecordWithCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(modifyRecordCall)

	return apiProcess.Run()
}

func (s *DnsService) DeletePrivateDnsRecord(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	removeCall, err := s.deleteZoneRecordWithCall(d, r)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(removeCall)

	return apiProcess.Run()
}

func (s *DnsService) ReadAndSetPrivateDnsRecord(d *schema.ResourceData, r *schema.Resource) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadPrivateDnsRecord(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.NonRetryableError(callErr)
			} else {
				return resource.RetryableError(fmt.Errorf("error on reading PrivateDnsRecord %q, %s", d.Id(), callErr))
			}
		} else {
			recordValueFunc := func(i interface{}) interface{} {
				var r map[string]interface{}
				switch i.(type) {
				case map[string]interface{}:
					r = i.(map[string]interface{})
				}
				v, _ := getSdkValue("RecordDataSet.0.RecordValue", r)
				if v == nil {
					return ""
				}
				return v
			}
			SdkResponseAutoResourceData(d, r, data, nil)

			rValue := recordValueFunc(data)
			_ = d.Set("record_value", rValue)

			return nil
		}
	})
}

func (s *DnsService) bindOrUnbindZoneVpcWithCall(d *schema.ResourceData, r *schema.Resource, isBind bool) (callback ApiCall, err error) {
	var (
		zoneId string
		vpcId  string
	)

	params := make(map[string]interface{})

	zoneId = d.Get("zone_id").(string)
	if v, ok := d.GetOk("vpc_set"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["vpc_id"]; ok {
				vpcId = v.(string)
				params["Vpcs.1.VpcId.1"] = vpcId
			}
			if v, ok := dMap["region_name"]; ok {
				params["Vpcs.1.RegionName"] = v
			}
		}
	}
	if len(params) < 1 {
		err = errors.New("vpc_set is empty, which should be spcify")
		return
	}

	params["ZoneId"] = zoneId
	targetStatus := "active"
	callback = ApiCall{
		param:  &params,
		action: "BindZoneVpc",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.pdnsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))

			if isBind {
				resp, err = conn.BindZoneVpc(call.param)
			} else {
				targetStatus = "empty"
				resp, err = conn.UnbindZoneVpc(call.param)
			}
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			if isBind {
				d.SetId(AssembleIds(zoneId, vpcId))
			}
			err = s.checkPdnsBindState(d, []string{targetStatus}, d.Timeout(schema.TimeoutCreate))
			return err
		},
	}
	if isBind {
		callback.action = "UnbindZoneVpc"
	}
	return callback, err
}

func (s *DnsService) BindZoneVpc(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	bindCall, err := s.bindOrUnbindZoneVpcWithCall(d, r, true)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(bindCall)

	return apiProcess.Run()
}

func (s *DnsService) UnbindZoneVpc(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	unbindCall, err := s.bindOrUnbindZoneVpcWithCall(d, r, false)
	if err != nil {
		return err
	}

	apiProcess.PutCalls(unbindCall)

	return apiProcess.Run()
}

func (s *DnsService) ReadAndSetZoneVpcAttachment(d *schema.ResourceData, r *schema.Resource) error {

	vpcSet, err := s.ReadZoneVpcAttachment(d)
	if err != nil {
		return err
	}

	vpcList := make([]interface{}, 0, len(vpcSet))

	for _, v := range vpcSet {
		delete(v, "status")
		vpcList = append(vpcList, v)
	}
	err = d.Set("vpc_set", vpcList)
	if err != nil {
		return err
	}
	return nil
}

func (s *DnsService) ReadZoneVpcAttachment(d *schema.ResourceData) (attachments []map[string]interface{}, err error) {
	data, err := s.ReadPrivateDnsZone(d, d.Get("zone_id").(string))
	if err != nil {
		return attachments, err
	}

	fullId := DisassembleIds(d.Id())
	vpcId := fullId[1]

	defaultErr := errors.New(fmt.Sprintf("the bind relationship is not exist, %s", d.Id()))
	bindSetIf, ok := data["BindVpcSet"]
	if !ok {
		return attachments, defaultErr
	}
	vpcSet := make([]map[string]interface{}, 0, 1)

	switch bindSetIf.(type) {
	case []interface{}:
		bindSet := bindSetIf.([]interface{})

		for _, bindVpc := range bindSet {
			bMap := make(map[string]interface{})
			bindVpcMap, ok := bindVpc.(map[string]interface{})
			if !ok {
				return attachments, defaultErr
			}

			if v, ok := bindVpcMap["VpcId"]; ok {
				if !reflect.DeepEqual(vpcId, v) {
					continue
				}
				bMap["vpc_id"] = v
				bMap["status"] = bindVpcMap["Status"]

				if v, ok := bindVpcMap["RegionName"]; ok {
					bMap["region_name"] = v
				}
				vpcSet = append(vpcSet, bMap)
			}
		}
	}
	attachments = vpcSet
	if len(attachments) < 1 {
		err = defaultErr
	}
	return
}

func (s *DnsService) checkPdnsBindState(d *schema.ResourceData, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{},
		Target:       target,
		Refresh:      s.pdnsBindStateRefreshFunc(d, []string{"error"}),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
		MinTimeout:   1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *DnsService) pdnsBindStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err    error
			status string
		)
		data, err := s.ReadZoneVpcAttachment(d)
		if err != nil {
			if notFoundError(err) {
				isNotExist := "empty"
				return data, isNotExist, nil
			}
			return nil, "", err
		}
		for _, statusRemote := range data {
			if v, ok := statusRemote["status"]; ok {
				status = v.(string)
			}
		}

		for _, v := range failStates {
			if v == status {
				return nil, "", fmt.Errorf("pdns zone bind vpc status error, status:%v", status)
			}
		}
		return data, status, nil
	}
}
