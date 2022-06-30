package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
)

type TagService struct {
	client *KsyunClient
}

func (s *TagService) ReadTags(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "PageSize", "Page", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.tagconn
		action := "ListTags"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.ListTags(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.ListTags(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("Tags", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *TagService) ReadTagKeys(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "PageSize", "Page", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.tagconn
		action := "ListTagKeys"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.ListTagKeys(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.ListTagKeys(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("TagKeys", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *TagService) ReadTagValues(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "PageSize", "Page", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.tagconn
		action := "ListTagValues"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			return data, fmt.Errorf("TagKey must set when ListTagValues")
		} else {
			if _, ok := condition["TagKey"]; !ok {
				return data, fmt.Errorf("TagKey must set when ListTagValues")
			}
			resp, err = conn.ListTagValues(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("TagValues", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *TagService) ReadTagsByResourceIds(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	conn := s.client.tagconn
	action := "ListTagsByResourceIds"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		return data, fmt.Errorf("ResourceType and ResourceUuids must set when ListTagsByResourceIds")
	} else {
		if _, ok := condition["ResourceType"]; !ok {
			return data, fmt.Errorf("ResourceType must set when ListTagsByResourceIds")
		}
		if _, ok := condition["ResourceUuids"]; !ok {
			return data, fmt.Errorf("ResourceUuids must set when ListTagsByResourceIds")
		}
		resp, err = conn.ListTagsByResourceIds(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("Tags", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *TagService) ReadTagByTagValue(d *schema.ResourceData, tagKey string, tagValue string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	req := map[string]interface{}{
		"TagKeys": tagKey,
	}
	results, err = s.ReadTagValues(req)
	if err != nil {
		return data, err
	}
	if len(results) == 0 {
		return data, fmt.Errorf("tagKey %s not exist ", tagKey)
	}
	var findValue bool
	for _, v := range results {
		data = v.(map[string]interface{})
		if data["Value"] == tagValue {
			findValue = true
			break
		}
	}
	if !findValue {
		return nil, fmt.Errorf("tagValue %s not exist ", tagKey)
	}
	return data, err
}

func (s *TagService) ReadTagByResourceId(d *schema.ResourceData, resourceId string, resourceType string) (data []interface{}, err error) {
	req := map[string]interface{}{
		"ResourceType":  resourceType,
		"ResourceUuids": resourceId,
	}
	data, err = s.ReadTagsByResourceIds(req)
	if err != nil {
		return data, err
	}
	return data, err
}

func (s *TagService) CreateTagCommonCall(req map[string]interface{}, isSetId bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:  &req,
		action: "CreateTag",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.tagconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateTag(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			if isSetId {
				id, err = getSdkValue("TagId", *resp)
				if err != nil {
					return err
				}
				d.SetId(id.(string))
			}
			return err
		},
	}
	return callback, err
}
func (s *TagService) ReplaceResourcesTagsCommonCall(req map[string]interface{}, disableDryRun bool) (callback ApiCall, err error) {
	callback = ApiCall{
		param:         &req,
		action:        "ReplaceResourcesTags",
		disableDryRun: disableDryRun,
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			if _, ok := (*call.param)["ReplaceTags"]; !ok {
				instanceId := d.Id()
				if instanceId == "" {
					instanceId = "tempId"
				}
				var tags []interface{}
				tags = append(tags, map[string]interface{}{
					"ResourceUuids": instanceId,
				})
				(*call.param)["ReplaceTags"] = tags
			}
			conn := client.tagconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ReplaceResourcesTags(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *TagService) CreateTagCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.CreateTagCommonCall(req, true)
}

func (s *TagService) ReplaceResourcesTagsCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	return s.ReplaceResourcesTagsCommonCall(req, false)
}

func (s *TagService) ReplaceResourcesTagsWithResourceCall(d *schema.ResourceData, r *schema.Resource, resourceType string, isUpdate bool, disableDryRun bool) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"tags": {
			FieldReqFunc: func(i interface{}, s string, m map[string]string, i2 int, s2 string, m2 *map[string]interface{}) (int, error) {
				if tagMap, ok := i.(map[string]interface{}); ok {
					for k, v := range tagMap {
						(*m2)["Tag_"+strconv.Itoa(i2)+"_Key"] = k
						(*m2)["Tag_"+strconv.Itoa(i2)+"_Value"] = v
						i2++
					}
				}
				return 0, nil
			},
		},
	}
	req, err := SdkRequestAutoMapping(d, r, isUpdate, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["ResourceType"] = resourceType
		return s.ReplaceResourcesTagsCommonCall(req, disableDryRun)
	}
	return callback, err
}
