package ksyun

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type ksyunDataSource struct {
	collection  []interface{}
	nameField   string
	idFiled     string
	targetField string
	extra       map[string]SdkResponseMapping
	compute     map[string]interface{}
}

type matchPlugin func(*schema.ResourceData, map[string]interface{}) (map[string]interface{}, bool, error)

func mergeDataSourcesReq(d *schema.ResourceData, r *schema.Resource, transform map[string]SdkReqTransform) (req map[string]interface{}, err error) {
	if transform == nil {
		transform = make(map[string]SdkReqTransform)
	}
	if _, ok := transform["name_regex"]; !ok {
		transform["name_regex"] = SdkReqTransform{
			Ignore: true,
		}
	}
	if _, ok := transform["output_file"]; !ok {
		transform["output_file"] = SdkReqTransform{
			Ignore: true,
		}
	}
	return SdkRequestAutoMapping(d, r, false, transform, nil,
		SdkReqParameter{false})
}

func mergeDataSourcesResp(d *schema.ResourceData, r *schema.Resource, dataSource ksyunDataSource, plugIns ...matchPlugin) (err error) {
	var (
		result []map[string]interface{}
	)

	if plugIns != nil && len(plugIns) > 0 {
		for _, plugIn := range plugIns {
			var filter []interface{}
			for _, item := range dataSource.collection {
				var (
					temp map[string]interface{}
					flag bool
				)
				temp, flag, err = plugIn(d, item.(map[string]interface{}))
				if err != nil {
					return err
				}
				if flag {
					if temp != nil {
						filter = append(filter, temp)
					}
				} else {
					filter = append(filter, item.(map[string]interface{}))
				}
			}
			dataSource.collection = filter
		}
	}

	for _, item := range dataSource.collection {
		var (
			temp map[string]interface{}
			flag bool
		)
		temp, flag, err = mergeNameRegex(d, item.(map[string]interface{}), dataSource.nameField)
		if err != nil {
			return err
		}
		if flag {
			if temp != nil {
				result = append(result, temp)
			}
		} else {
			result = append(result, item.(map[string]interface{}))
		}
	}
	_, _, err = SdkSliceMapping(d, result, SdkSliceData{
		IdField: dataSource.idFiled,
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item[idField].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(r, dataSource.targetField, item, dataSource.compute, dataSource.extra)
		},
		TargetName: dataSource.targetField,
	})
	return err
}

func mergeNameRegex(d *schema.ResourceData, data map[string]interface{}, nameField string) (result map[string]interface{}, flag bool, err error) {
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		match := regexp.MustCompile(nameRegex.(string))
		if match.MatchString(data[nameField].(string)) {
			return data, true, err
		}
		return nil, true, err
	}
	return nil, false, err
}

func mergeDataSourcesRespIdInObj(d *schema.ResourceData, r *schema.Resource, dataSource ksyunDataSource, plugIns ...matchPlugin) (err error) {
	var (
		result []map[string]interface{}
	)

	if plugIns != nil && len(plugIns) > 0 {
		for _, plugIn := range plugIns {
			var filter []interface{}
			for _, item := range dataSource.collection {
				var (
					temp map[string]interface{}
					flag bool
				)
				temp, flag, err = plugIn(d, item.(map[string]interface{}))
				if err != nil {
					return err
				}
				if flag {
					if temp != nil {
						filter = append(filter, temp)
					}
				} else {
					filter = append(filter, item.(map[string]interface{}))
				}
			}
			dataSource.collection = filter
		}
	}

	for _, item := range dataSource.collection {
		var (
			temp map[string]interface{}
			flag bool
		)
		temp, flag, err = mergeNameRegex(d, item.(map[string]interface{}), dataSource.nameField)
		if err != nil {
			return err
		}
		if flag {
			if temp != nil {
				result = append(result, temp)
			}
		} else {
			result = append(result, item.(map[string]interface{}))
		}
	}
	_, _, err = SdkSliceMapping(d, result, SdkSliceData{
		IdField: dataSource.idFiled,
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			logger.DebugInfo("idField ---- %+v ", idField)
			logger.DebugInfo("item ---- %+v ", item)
			//idFiled 以.切割成数组
			idFields := strings.Split(idField, ".")
			//不修改原始item 只要值
			itemT := item
			for _, v := range idFields {
				if v == idFields[len(idFields)-1] {
					return itemT[v].(string)
				} else {
					itemT = itemT[v].(map[string]interface{})
				}
			}
			return ""
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			return SdkResponseAutoMapping(r, dataSource.targetField, item, dataSource.compute, dataSource.extra)
		},
		TargetName: dataSource.targetField,
	})
	return err
}
