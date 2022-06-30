package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func mergeTagsData(d *schema.ResourceData, data *map[string]interface{}, client *KsyunClient, resourceType string) (err error) {
	var tags []interface{}
	tagService := TagService{client}
	tags, err = tagService.ReadTagByResourceId(d, d.Id(), resourceType)
	if err != nil {
		//此处暂时兼容如果没有更改tags可以忽略listTags的权限检查。做到最大兼容性
		if !d.HasChange("tags") {
			errMessage := strings.ToLower(err.Error())
			if strings.Contains(errMessage, "lack of policy") {
				return nil
			}
		}
		return err
	}
	tagMap := make(map[string]interface{})
	for _, tag := range tags {
		_m := tag.(map[string]interface{})
		tagMap[_m["TagKey"].(string)] = _m["TagValue"].(string)
	}
	if len(tagMap) > 0 {
		(*data)["Tags"] = tagMap
	}
	return err
}
