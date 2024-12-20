package ksyun

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type TagService struct {
	client *KsyunClient
}

type Tag struct {
	Id         int    `mapstructure:"Id"`
	Key        string `mapstructure:"Key"`
	Value      string `mapstructure:"Value"`
	CreateTime string `mapstructure:"CreateTime"`
	CanDelete  int    `mapstructure:"CanDelete"`
	IsBillTag  int    `mapstructure:"IsBillTag"`
}

type Tags []*Tag

func (t Tags) GetTagsParams(rsType, rsUuid string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for i, tag := range t {
		m["Tag_"+strconv.Itoa(i+1)+"_Key"] = tag.Key
		m["Tag_"+strconv.Itoa(i+1)+"_Value"] = tag.Value
	}
	m["ResourceType"] = rsType
	rpTagMap := map[string]interface{}{
		"ResourceUuids": rsUuid,
	}
	m["ReplaceTags"] = []interface{}{rpTagMap}
	return m, nil
}

func (tag Tag) GetTagParam(rsType, rsUuid string) map[string]interface{} {
	m := make(map[string]interface{})
	m["TagKey"] = tag.Key
	m["TagValue"] = tag.Value
	m["ResourceType"] = rsType
	m["ResourceUuid"] = rsUuid
	return m
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
	var results []interface{}
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
				d.SetId(Float64ToString(id.(float64)))
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

func (s *TagService) CreateTag(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	createTagCall, err := s.CreateTagCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createTagCall)
	return apiProcess.Run()
}

func (s *TagService) CreateTagResourceAttachment(d *schema.ResourceData, r *schema.Resource) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	createTagCall, err := s.tagAttachResourceWithCall(d, r)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(createTagCall)
	return apiProcess.Run()
}

func (s *TagService) tagAttachResourceWithCall(d *schema.ResourceData, r *schema.Resource) (ApiCall, error) {
	var (
		req    map[string]interface{}
		rsType = d.Get("resource_type").(string)
		rsId   = d.Get("resource_id").(string)
	)
	t := &Tag{
		Key:   d.Get("key").(string),
		Value: d.Get("value").(string),
	}
	tags := Tags{t}

	tagsMutex.Lock()
	// query existed tags
	results, err := s.ReadTagByResourceId(d, rsId, rsType)
	if err != nil {
		return ApiCall{}, fmt.Errorf("an error caused while merging tags, %s", err)
	}
	for _, res := range results {
		switch res.(type) {
		case map[string]interface{}:
			r := res.(map[string]interface{})
			k, ok := r["TagKey"]
			if !ok {
				continue
			}
			v, vOk := r["TagValue"]
			if !vOk {
				continue
			}
			tt := &Tag{
				Key:   k.(string),
				Value: v.(string),
			}
			tags = append(tags, tt)
		}
	}

	req, _ = tags.GetTagsParams(rsType, rsId)

	attachCall, err := s.ReplaceResourcesTagsCommonCall(req, true)
	if err != nil {
		return ApiCall{}, err
	}
	rawExecuteCall := attachCall.executeCall
	attachCall.executeCall = func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
		retryErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = rawExecuteCall(d, client, call)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			readErr := s.ReadAndSetTagAttachment(d, nil)
			if readErr != nil {
				if notFoundError(readErr) {
					return resource.RetryableError(readErr)
				}
				return resource.NonRetryableError(readErr)
			}
			return nil
		})
		if retryErr != nil {
			return nil, retryErr
		}
		return
	}
	rawErrorCall := attachCall.callError
	attachCall.callError = func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
		defer tagsMutex.Unlock()
		if rawErrorCall != nil {
			return rawErrorCall(d, client, call, baseErr)
		}
		return nil
	}
	rawAfterCall := attachCall.afterCall
	attachCall.afterCall = func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error {
		defer tagsMutex.Unlock()
		if rawAfterCall != nil {
			return rawAfterCall(d, client, resp, call)
		}
		return nil
	}
	return attachCall, nil
}

func (s *TagService) DetachResourceTags(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	call, err := s.DetachResourceTagsWithCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(call)
	return apiProcess.Run()
}

func (s *TagService) DetachResourceTagsWithCall(d *schema.ResourceData) (callback ApiCall, err error) {
	params := map[string]interface{}{}
	params["ResourceType"] = d.Get("resource_type")
	params["ResourceUuid"] = d.Get("resource_id")
	params["TagIds"] = d.Get("tag_id")

	callback = ApiCall{
		param:  &params,
		action: "DeleteTags",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.tagconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DetachResourceTags(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				if notFoundError(err) {
					return nil
				}
				return resource.RetryableError(err)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
}

func (s *TagService) ReadAndSetTag(d *schema.ResourceData, r *schema.Resource) error {
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	req := make(map[string]interface{})
	req["Key"] = key
	req["Value"] = value

	results, err := s.ReadTags(req)
	if err != nil {
		return err
	}

	var data map[string]interface{}

	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return fmt.Errorf("tag %s:%s is not exist ", key, value)
	}
	extra := map[string]SdkResponseMapping{
		"Id": {
			Field: "tag_id",
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return nil
}

func (s *TagService) ReadAndSetTagAttachment(d *schema.ResourceData, r *schema.Resource) error {
	var (
		rsId   = d.Get("resource_id").(string)
		rsType = d.Get("resource_type").(string)

		tagKey   = d.Get("key")
		tagValue = d.Get("value")

		tagId float64
		found bool
	)

	results, err := s.ReadTagByResourceId(d, rsId, rsType)
	if err != nil {
		return err
	}

	for _, result := range results {
		switch result.(type) {
		case map[string]interface{}:
			r := result.(map[string]interface{})
			if reflect.DeepEqual(tagKey, r["TagKey"]) && reflect.DeepEqual(tagValue, r["TagValue"]) {
				if v, ok := r["TagId"]; ok {
					tagId = v.(float64)
				}

				if err := d.Set("tag_id", tagId); err != nil {
					return err
				}

				id := AssembleIds(rsId, Float64ToString(tagId))
				d.SetId(id)
				found = true
				goto exit
			}
		}
	}

exit:
	if !found {
		return fmt.Errorf("the attachment between tag and resource %s is not exist", rsId)
	}
	return nil
}

func (s *TagService) DeleteTag(d *schema.ResourceData) error {
	apiProcess := NewApiProcess(context.Background(), d, s.client, true)

	deleteCall, err := s.DeleteTagCall(d)
	if err != nil {
		return err
	}
	apiProcess.PutCalls(deleteCall)

	return apiProcess.Run()
}

func (s *TagService) DeleteTagCall(d *schema.ResourceData) (callback ApiCall, err error) {
	// 构成参数
	params := map[string]interface{}{}

	tags := make([]map[string]interface{}, 0)

	tag := map[string]interface{}{
		"Key":   d.Get("key"),
		"Value": d.Get("value"),
	}
	params["Tags"] = append(tags, tag)

	callback = ApiCall{
		param:  &params,
		action: "DeleteTag",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.tagconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteTag(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(5*time.Minute, func() *resource.RetryError {
				if notFoundError(baseErr) {
					return nil
				}

				// it cannot be deleted if this is still using
				if isExpectError(baseErr, []string{"TagDeleteConflict"}) {
					return resource.NonRetryableError(baseErr)
				}
				return resource.RetryableError(baseErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
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
	if len(req) > 0 || d.HasChange("tags") {
		req["ResourceType"] = resourceType
		return s.ReplaceResourcesTagsCommonCall(req, disableDryRun)
	}
	return callback, err
}
