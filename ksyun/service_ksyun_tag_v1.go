package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type TagV1Service struct {
	client *KsyunClient
}

func (s *TagV1Service) ReadTags(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.tagv1conn
		action := "DescribeTags"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeTags(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeTags(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("TagSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *TagV1Service) ReadAndSetTag(d *schema.ResourceData, r *schema.Resource) (err error) {

	params := map[string]interface{}{}
	params["Filter.1.Name"] = "key"
	params["Filter.1.Value.1"] = d.Get("key")
	params["Filter.2.Name"] = "value"
	params["Filter.2.Value.1"] = d.Get("value")
	params["Filter.3.Name"] = "resource-type"
	params["Filter.3.Value.1"] = d.Get("resource_type")
	params["Filter.4.Name"] = "resource-id"
	params["Filter.4.Value.1"] = d.Get("resource_id")

	var data []interface{}
	data, err = s.ReadTags(params)
	SdkResponseAutoResourceData(d, r, data, nil)

	return
}

func (s *TagV1Service) CreateTagCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}
	params["Tag.0.Key"] = d.Get("key")
	params["Tag.0.Value"] = d.Get("value")
	params["Resource.0.Type"] = d.Get("resource_type")
	params["Resource.0.Id"] = d.Get("resource_id")
	callback = ApiCall{
		param:  &params,
		action: "CreateTags",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.tagv1conn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateTags(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			t_key := fmt.Sprintf("%s", d.Get("key"))
			t_value := fmt.Sprintf("%s", d.Get("value"))
			r_type := fmt.Sprintf("%s", d.Get("resource_type"))
			r_id := fmt.Sprintf("%s", d.Get("resource_id"))
			d.SetId(t_key + ":" + t_value + "," + r_type + ":" + r_id)
			return err
		},
	}
	return callback, err
}

func (s *TagV1Service) DeleteTagCall(d *schema.ResourceData) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}
	params["Tag.0.Key"] = d.Get("key")
	params["Tag.0.Value"] = d.Get("value")
	params["Resource.0.Type"] = d.Get("resource_type")
	params["Resource.0.Id"] = d.Get("resource_id")
	callback = ApiCall{
		param:  &params,
		action: "DeleteTags",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.tagv1conn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteTags(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
}

func (s *TagV1Service) ModifyTag(d *schema.ResourceData, r *schema.Resource) (err error) {
	return
}

func (s *TagV1Service) CreateTag(d *schema.ResourceData, r *schema.Resource) (err error) {

	call, err := s.CreateTagCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *TagV1Service) DeleteTag(d *schema.ResourceData) (err error) {
	call, err := s.DeleteTagCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *TagV1Service) ReadAndSetTags(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"keys": {
			mapping: "key",
			Type:    TransformWithFilter,
		},
		"values": {
			mapping: "value",
			Type:    TransformWithFilter,
		},
		"resource_types": {
			mapping: "resource-type",
			Type:    TransformWithFilter,
		},
		"resource_ids": {
			mapping: "resource-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	logger.Debug(logger.ReqFormat, "Demo", req)
	if err != nil {
		return err
	}
	data, err := s.ReadTags(req)
	if err != nil {
		return err
	}

	//for _, item := range data {
	//	if itemData, ok := item.(map[string]interface{}); ok {
	//		itemData["id"] = fmt.Sprintf("%s:%s,%s:%s", itemData["Key"], itemData["Value"], itemData["ResourceType"], itemData["ResourceId"])
	//	} else {
	//		itemData["id"] = ""
	//	}
	//}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "tags",
		extra:       map[string]SdkResponseMapping{},
	})
}
