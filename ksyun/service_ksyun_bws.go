package ksyun

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type BwsService struct {
	client *KsyunClient
}

func (s *BwsService) ReadBandWidthShares(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.bwsconn
	action := "DescribeBandWidthShares"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeBandWidthShares(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeBandWidthShares(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("BandWidthShareSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *BwsService) ReadBandWidthShare(d *schema.ResourceData, bwsId string) (data map[string]interface{}, err error) {
	var results []interface{}
	if bwsId == "" {
		bwsId = d.Id()
	}
	req := map[string]interface{}{
		"BandWidthShareId.1": bwsId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}

	if _, ok := d.GetOk("tags"); ok {
		req["IsContainTag"] = true
	}

	results, err = s.ReadBandWidthShares(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("BandWidthShare %s not exist ", bwsId)
	}
	return data, err
}

func (s *BwsService) ReadAndSetBandWidthShare(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadBandWidthShare(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading bandWidthShare %q, %s", d.Id(), callErr))
			}
		} else {
			extra := chargeExtraForVpc(data)
			extra["TagSet"] = SdkResponseMapping{
				Field: "tags",
				FieldRespFunc: func(i interface{}) interface{} {
					tags := i.([]interface{})
					tagMap := make(map[string]interface{})
					for _, tag := range tags {
						_m := tag.(map[string]interface{})
						tagMap[_m["TagKey"].(string)] = _m["TagValue"].(string)
					}
					return tagMap
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (s *BwsService) ReadAndSetBandWidthShares(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "BandWidthShareId",
			Type:    TransformWithN,
		},
		"allocation_ids": {
			mapping: "allocation-id",
			Type:    TransformWithFilter,
		},
		"project_ids": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadBandWidthShares(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "BandWidthShareName",
		idFiled:     "BandWidthShareId",
		targetField: "band_width_shares",
		extra: map[string]SdkResponseMapping{
			"BandWidthShareName": {
				Field:    "name",
				KeepAuto: true,
			},
			"BandWidthShareId": {
				Field:    "id",
				KeepAuto: true,
			},
			"AssociateBandWidthShareInfoSet": {
				Field: "allocation_ids",
				FieldRespFunc: func(i interface{}) interface{} {
					var result []interface{}
					for _, v := range i.([]interface{}) {
						result = append(result, v.(map[string]interface{})["AllocationId"])
					}
					return result
				},
			},
		},
	})
}

func (s *BwsService) CreateBandWidthShareCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateBandWidthShare",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.bwsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateBandWidthShare(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("BandWidthShareId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *BwsService) CreateBandWidthShare(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateBandWidthShareCall(d, r)
	if err != nil {
		return err
	}
	tagsService := TagService{client: s.client}
	tagsCall, err := tagsService.ReplaceResourcesTagsWithResourceCall(d, r, "bws", false, false)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call, tagsCall}, d, s.client, true)
}

func (s *BwsService) ModifyBandWidthShareProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
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

func (s *BwsService) ModifyBandWidthShareCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
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
		req["BandWidthShareId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyBandWidthShare",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.bwsconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyBandWidthShare(call.param)
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

func (s *BwsService) ModifyBandWidthShare(d *schema.ResourceData, r *schema.Resource) (err error) {
	var calls []ApiCall
	projectCall, err := s.ModifyBandWidthShareProjectCall(d, r)
	if err != nil {
		return err
	}
	calls = append(calls, projectCall)
	call, err := s.ModifyBandWidthShareCall(d, r)
	if err != nil {
		return err
	}
	calls = append(calls, call)
	if d.HasChange("tags") {
		tagsService := TagService{client: s.client}
		tagsCall, err := tagsService.ReplaceResourcesTagsWithResourceCall(d, r, "bws", true, false)
		if err != nil {
			return err
		}
		calls = append(calls, tagsCall)
	}

	return ksyunApiCallNew(calls, d, s.client, true)
}

func (s *BwsService) RemoveBandWidthShareCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"BandWidthShareId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteBandWidthShare",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.bwsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteBandWidthShare(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadBandWidthShare(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading bandWidthShare when delete %q, %s", d.Id(), callErr))
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

func (s *BwsService) RemoveBandWidthShare(d *schema.ResourceData) (err error) {
	call, err := s.RemoveBandWidthShareCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *BwsService) ReadBandWidthShareAssociate(d *schema.ResourceData, bwsId string, allocationId string) (result map[string]interface{}, err error) {
	data, err := s.ReadBandWidthShare(d, bwsId)
	if err != nil {
		return nil, err
	}
	result = make(map[string]interface{})
	if len(data["AssociateBandWidthShareInfoSet"].([]interface{})) == 0 {
		return data, fmt.Errorf("AllocationId %s not associate in BandWidthShare %s ", allocationId, bwsId)
	}

	isFound := false
	for _, v := range data["AssociateBandWidthShareInfoSet"].([]interface{}) {
		if v1, ok := v.(map[string]interface{})["AllocationId"]; ok && v1 == allocationId {
			isFound = true
			result["AllocationId"] = v1
			break
		}
	}

	if !isFound {
		return data, fmt.Errorf("AllocationId %s not associate in BandWidthShare %s ", allocationId, bwsId)
	}

	result["BandWidthShareId"] = data["BandWidthShareId"]
	return result, err
}

func (s *BwsService) ReadAndSetAssociateBandWidthShare(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, err := s.ReadBandWidthShareAssociate(d, d.Get("band_width_share_id").(string), d.Get("allocation_id").(string))
		if err != nil {
			if strings.Contains(err.Error(), "not associate") {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		SdkResponseAutoResourceData(d, r, data, nil)
		return nil
	})
}

func (s *BwsService) AssociateBandWidthShareCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	// read eip
	eipService := EipService{s.client}
	eipData, err := eipService.ReadAddress(d, d.Get("allocation_id").(string))
	if err != nil {
		return callback, err
	}
	bandWidth := eipData["BandWidth"]
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "AssociateBandWidthShare",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.bwsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.AssociateBandWidthShare(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			d.SetId(d.Get("band_width_share_id").(string) + ":" + d.Get("allocation_id").(string))
			return d.Set("band_width", bandWidth)
		},
	}
	return callback, err
}

func (s *BwsService) AssociateBandWidthShare(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.AssociateBandWidthShareCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *BwsService) DisassociateBandWidthShareCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"BandWidthShareId": d.Get("band_width_share_id"),
		"AllocationId":     d.Get("allocation_id"),
	}
	if d.Get("band_width") == 0 {
		removeReq["BandWidth"] = 1
	} else {
		removeReq["BandWidth"] = d.Get("band_width")
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DisassociateBandWidthShare",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.bwsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DisassociateBandWidthShare(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadBandWidthShare(d, d.Get("band_width_share_id").(string))
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading bandWidthShare associate when disassociate %q, %s", d.Id(), callErr))
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

func (s *BwsService) DisassociateBandWidthShare(d *schema.ResourceData) (err error) {
	call, err := s.DisassociateBandWidthShareCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
