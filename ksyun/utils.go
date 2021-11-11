package ksyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func SchemaSetToInstanceMap(s interface{}, prefix string, input *map[string]interface{}) {
	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)[prefix+"."+strconv.Itoa(count)] = v
	}
}

func SchemaSetToFilterMap(s interface{}, prefix string, index int, input *map[string]interface{}) {
	(*input)["Filter."+strconv.Itoa(index)+".Name"] = prefix
	count := int(0)
	for _, v := range s.(*schema.Set).List() {
		count = count + 1
		(*input)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(count)] = v
	}
}

func SchemaSetsToFilterMap(d *schema.ResourceData, filters []string, req *map[string]interface{}) *map[string]interface{} {
	index := 0
	for _, v := range filters {
		var idsString []string
		if ids, ok := d.GetOk(v); ok {
			idsString = SchemaSetToStringSlice(ids)
		}
		if len(idsString) > 0 {
			index++
			(*req)[fmt.Sprintf("Filter.%v.Name", index)] = strings.Replace(v, "_", "-", -1)
		}
		for k1, v1 := range idsString {
			if v1 == "" {
				continue
			}
			(*req)[fmt.Sprintf("Filter.%v.Value.%d", index, k1+1)] = v1
		}
	}
	return req
}
func hashStringArray(arr []string) string {
	var buf bytes.Buffer

	for _, s := range arr {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func writeToFile(filePath string, data interface{}) error {
	absPath, err := getAbsPath(filePath)
	if err != nil {
		return err
	}
	os.Remove(absPath)
	var bs []byte
	switch data := data.(type) {
	case string:
		bs = []byte(data)
	default:
		bs, err = json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v and got an error: %#v", data, err)
		}
	}

	return ioutil.WriteFile(absPath, bs, 0422)
}

func getAbsPath(filePath string) (string, error) {
	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("get current user got an error: %#v", err)
		}

		if usr.HomeDir != "" {
			filePath = strings.Replace(filePath, "~", usr.HomeDir, 1)
		}
	}
	return filepath.Abs(filePath)
}

func merageResultDirect(result *[]map[string]interface{}, source []interface{}) {
	for _, v := range source {
		*result = append(*result, v.(map[string]interface{}))
	}
}

// schemaSetToStringSlice used for converting terraform schema set to a string slice
func SchemaSetToStringSlice(s interface{}) []string {
	vL := []string{}

	for _, v := range s.(*schema.Set).List() {
		vL = append(vL, v.(string))
	}

	return vL
}

func getSdkValue(keyPattern string, obj interface{}) (interface{}, error) {
	keys := strings.Split(keyPattern, ".")
	root := obj
	for index, k := range keys {
		if reflect.ValueOf(root).Kind() == reflect.Map {
			root = root.(map[string]interface{})[k]
			if root == nil {
				return root, nil
			}

		} else if reflect.ValueOf(root).Kind() == reflect.Slice {
			i, err := strconv.Atoi(k)
			if err != nil {
				return nil, fmt.Errorf("keyPattern %s index %d must number", keyPattern, index)
			}
			if len(root.([]interface{})) < i {
				return nil, nil
			}
			root = root.([]interface{})[i]
		}
	}
	return root, nil
}

type SliceMappingFunc func(map[string]interface{}) map[string]interface{}

type IdMappingFunc func(string, map[string]interface{}) string

type SdkSliceData struct {
	IdField          string
	IdMappingFunc    IdMappingFunc
	SliceMappingFunc SliceMappingFunc
	TargetName       string
}

func sliceMapping(ids []string, data []map[string]interface{}, sdkSliceData SdkSliceData, item interface{}) ([]string, []map[string]interface{}) {
	if mm, ok := item.(map[string]interface{}); ok {
		if sdkSliceData.IdMappingFunc != nil && sdkSliceData.IdField != "" {
			ids = append(ids, sdkSliceData.IdMappingFunc(sdkSliceData.IdField, mm))
		}
		if sdkSliceData.SliceMappingFunc != nil {
			data = append(data, sdkSliceData.SliceMappingFunc(mm))
		}
	}
	return ids, data
}

func mapMapping(sdkSliceData SdkSliceData, item interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	if mm, ok := item.(map[string]interface{}); ok {
		if sdkSliceData.SliceMappingFunc != nil {
			data = sdkSliceData.SliceMappingFunc(mm)
		}
	}
	return data
}

func getSchemeElem(resource *schema.Resource, keys []string) *schema.Resource {
	r := resource
	if r == nil {
		return nil
	}
	for _, v := range keys {
		if elem, o := r.Schema[v].Elem.(*schema.Resource); o {
			r = elem
		}
	}
	return r
}

func transformFieldReqFunc(v interface{}, k string, t SdkReqTransform, index int, req *map[string]interface{}) (bool, int, error) {
	if t.FieldReqFunc != nil {
		i, err := t.FieldReqFunc(v, t.mapping, t.mappings, index, k, req)
		if err != nil {
			return false, i, err
		}
		return true, i, nil
	}
	return false, index, nil
}

func transformDefault(v interface{}, k string, t SdkReqTransform, req *map[string]interface{}) error {
	if ok, _, err := transformFieldReqFunc(v, k, t, 0, req); ok {
		return err
	}
	if strings.TrimSpace(t.mapping) == "" {
		(*req)[Downline2Hump(k)] = v
	} else {
		(*req)[t.mapping] = v
	}
	return nil
}

func transformSingleN(v interface{}, k string, t SdkReqTransform, req *map[string]interface{}) error {
	if ok, _, err := transformFieldReqFunc(v, k, t, 0, req); ok {
		if err != nil {
			return fmt.Errorf("error on transformSingleN with transformFieldReqFunc %s", err)
		}
		return nil
	}
	if strings.TrimSpace(t.mapping) == "" {
		(*req)[Downline2Hump(k)+".1"] = v
	} else {
		(*req)[t.mapping+".1"] = v
	}
	return nil
}

func transformWithN(v interface{}, k string, t SdkReqTransform, req *map[string]interface{}) error {
	if x, ok := v.(*schema.Set); ok {
		if ok, _, err := transformFieldReqFunc(v, k, t, 0, req); ok {
			if err != nil {
				return fmt.Errorf("error on transformWithN with transformFieldReqFunc %s", err)
			}
			return nil
		}
		for i, value := range (*x).List() {
			if strings.TrimSpace(t.mapping) == "" {
				(*req)[Downline2Hump(k)+"."+strconv.Itoa(i+1)] = value
			} else {
				(*req)[t.mapping+"."+strconv.Itoa(i+1)] = value
			}
		}
	}
	if x, ok := v.([]interface{}); ok {
		for i, value := range x {
			if strings.TrimSpace(t.mapping) == "" {
				(*req)[Downline2Hump(k)+"."+strconv.Itoa(i+1)] = value
			} else {
				(*req)[t.mapping+"."+strconv.Itoa(i+1)] = value
			}
		}
	}
	if x, ok := v.(string); ok {
		for i, value := range strings.Split(x, ",") {
			if strings.TrimSpace(t.mapping) == "" {
				(*req)[Downline2Hump(k)+"."+strconv.Itoa(i+1)] = value
			} else {
				(*req)[t.mapping+"."+strconv.Itoa(i+1)] = value
			}
		}
	}

	return nil
}

func transformListN(v interface{}, k string, t SdkReqTransform, req *map[string]interface{}) error {
	if list, ok := v.([]interface{}); ok {
		if ok, _, err := transformFieldReqFunc(v, k, t, 0, req); ok {
			if err != nil {
				return fmt.Errorf("error on transformListN with transformFieldReqFunc %s", err)
			}
			return nil
		}
		for index, v1 := range list {
			if m1, ok := v1.(map[string]interface{}); ok {
				for k2, v2 := range m1 {
					k3 := getFinalKey(t, k) + "." + strconv.Itoa(index+1) + "." + getFinalKey(t, k2)
					(*req)[k3] = v2
				}
			}
		}
	}
	return nil
}

func transformListUnique(v interface{}, k string, t SdkReqTransform, req *map[string]interface{}, d *schema.ResourceData, forceGet bool) error {
	if list, ok := v.([]interface{}); ok {
		if ok, _, err := transformFieldReqFunc(v, k, t, 0, req); ok {
			if err != nil {
				return fmt.Errorf("error on transformListN with transformFieldReqFunc %s", err)
			}
			return nil
		}
		for i, v1 := range list {
			if m1, ok := v1.(map[string]interface{}); ok {
				for k2, v2 := range m1 {
					flag := false
					schemaKey := fmt.Sprintf("%s.%d.%s", k, i, k2)
					if forceGet {
						if t.forceUpdateParam || (d.HasChange(schemaKey) && !d.IsNewResource()) {
							flag = true
						}
					} else {
						if _, ok := d.GetOk(schemaKey); ok {
							flag = true
						}
					}
					if flag {
						k3 := k + "." + k2
						if target, ok := t.mappings[k3]; ok {
							(*req)[target] = v2
						} else {
							(*req)[fmt.Sprintf("%s.%s", Downline2Hump(k), Downline2Hump(k2))] = v2
						}
					}
				}
				break
			}
		}
	}
	return nil
}

func transformWithFilter(v interface{}, k string, t SdkReqTransform, index int, req *map[string]interface{}) (int, error) {
	if x, ok := v.([]interface{}); ok {
		v = schema.NewSet(schema.HashString, x)
	}
	if x, ok := v.(string); ok {
		if ok, j, err := transformFieldReqFunc(v, k, t, index, req); ok {
			if err != nil {
				return index, fmt.Errorf("error on transformWithFilter with transformFieldReqFunc %s", err)
			}
			return j, nil
		}
		if strings.TrimSpace(t.mapping) == "" {
			(*req)["Filter."+strconv.Itoa(1)+".Name"] = Downline2Filter(k)
		} else {
			(*req)["Filter."+strconv.Itoa(1)+".Name"] = t.mapping
		}
		(*req)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(1)] = x
	}
	if x, ok := v.(*schema.Set); ok {
		if ok, j, err := transformFieldReqFunc(v, k, t, index, req); ok {
			if err != nil {
				return index, fmt.Errorf("error on transformWithFilter with transformFieldReqFunc %s", err)
			}
			return j, nil
		}
		for i, value := range (*x).List() {
			if i == 0 {
				if strings.TrimSpace(t.mapping) == "" {
					(*req)["Filter."+strconv.Itoa(index)+".Name"] = Downline2Filter(k)
				} else {
					(*req)["Filter."+strconv.Itoa(index)+".Name"] = t.mapping
				}

			}
			(*req)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(i+1)] = value
		}
		index = index + 1
	}
	return index, nil
}

func transformListFilter(v interface{}, k string, t SdkReqTransform, index int, req *map[string]interface{}) (int, error) {
	var err error
	if list, ok := v.([]interface{}); ok {
		if ok, j, err := transformFieldReqFunc(v, k, t, index, req); ok {
			if err != nil {
				return index, fmt.Errorf("error on transformListFilter with transformFieldReqFunc %s", err)
			}
			return j, nil
		}
		for _, v1 := range list {
			if m1, ok := v1.(map[string]interface{}); ok {
				for k2, v2 := range m1 {
					if v3, ok := v2.(*schema.Set); ok {
						for i, v4 := range (*v3).List() {
							if i == 0 {
								k = k + "." + k2
								if target, ok := t.mappings[k]; ok {
									(*req)["Filter."+strconv.Itoa(index)+".Name"] = target
								} else {
									(*req)["Filter."+strconv.Itoa(index)+".Name"] = Downline2Filter(k)
								}

							}
							(*req)["Filter."+strconv.Itoa(index)+".Value."+strconv.Itoa(i+1)] = v4
						}
						index = index + 1
					} else {
						index, err = transformListFilter(v2, k, t, index, req)
						if err != nil {
							return index, err
						}
					}
				}
			}
		}
	}
	return index, nil
}

func getFinalKey(t SdkReqTransform, k string) string {
	if target, ok := t.mappings[k]; ok {
		return target
	} else {
		return Downline2Hump(k)
	}
}

func SdkRequestAutoExtra(r map[string]SdkReqTransform, d *schema.ResourceData, forceGet bool) map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = make(map[string]SdkRequestMapping)
	index := 1
	for k := range r {
		extra[k] = SdkRequestMapping{
			FieldReqFunc: func(item interface{}, field string, source string, m *map[string]interface{}) error {
				var err error
				switch r[source].Type {
				case TransformListUnique:
					err = transformListUnique(item, source, r[source], m, d, forceGet)
				case TransformWithN:
					err = transformWithN(item, source, r[source], m)
				case TransformListN:
					err = transformListN(item, source, r[source], m)
				case TransformDefault:
					err = transformDefault(item, source, r[source], m)
				case TransformWithFilter:
					index, err = transformWithFilter(item, source, r[source], index, m)
				case TransformListFilter:
					index, err = transformListFilter(item, source, r[source], index, m)
				}
				return err
			},
		}
	}
	return extra
}

// Auto Transform Terraform Resource to SDK Request Parameter
// d : Transform schema.ResourceData Ptr
// resource: Transform schema.Resource Ptr
// onlyTransform : map[string]TransformType ,If set this field,Transform will with this array instead of Transform schema.Resource
// extraMapping : map[string]SdkRequestMapping , if set this field, the key in map will instead of Transform schema.Resource key or only key
func SdkRequestAutoMapping(d *schema.ResourceData, resource *schema.Resource, isUpdate bool, transform map[string]SdkReqTransform, extraMapping map[string]SdkRequestMapping, params ...SdkReqParameter) (map[string]interface{}, error) {
	var req map[string]interface{}
	var err error
	req = make(map[string]interface{})
	count := 1
	var onlyMode bool
	if params != nil && len(params) == 1 {
		onlyMode = params[0].onlyTransform
	} else {
		onlyMode = true
	}

	if transform != nil && onlyMode {
		for k, v := range transform {
			if isUpdate {
				count, err = requestUpdateMapping(d, k, v, count, nil, &req)
			} else {
				count, err = requestCreateMapping(d, k, v, count, nil, &req, false)
			}
		}
	} else {
		for k := range resource.Schema {
			if v, ok := transform[k]; ok {
				if isUpdate {
					count, err = requestUpdateMapping(d, k, v, count, nil, &req)
				} else {
					count, err = requestCreateMapping(d, k, v, count, nil, &req, false)
				}
			} else {
				if isUpdate {
					count, err = requestUpdateMapping(d, k, SdkReqTransform{}, count, extraMapping, &req)
				} else {
					count, err = requestCreateMapping(d, k, SdkReqTransform{}, count, extraMapping, &req, false)
				}
			}
		}
	}

	return req, err
}

func requestCreateMapping(d *schema.ResourceData, k string, t SdkReqTransform, index int, extraMapping map[string]SdkRequestMapping, req *map[string]interface{}, forceGet bool) (int, error) {
	var err error
	var ok bool
	var v interface{}

	if t.Ignore {
		return index, err
	}

	if t.ValueFunc != nil {
		v, ok = t.ValueFunc(d)
	} else {
		if forceGet {
			v = d.Get(k)
			ok = true

		} else {
			v, ok = d.GetOk(k)
		}
	}

	if extraMapping == nil {
		extraMapping = make(map[string]SdkRequestMapping)
	}
	if ok {
		if _, ok := extraMapping[k]; !ok {
			switch t.Type {
			case TransformDefault:
				err = transformDefault(v, k, t, req)
			case TransformSingleN:
				err = transformSingleN(v, k, t, req)
			case TransformWithN:
				err = transformWithN(v, k, t, req)
			case TransformListN:
				err = transformListN(v, k, t, req)
			case TransformListUnique:
				err = transformListUnique(v, k, t, req, d, forceGet)
			case TransformWithFilter:
				index, err = transformWithFilter(v, k, t, index, req)
			case TransformListFilter:
				index, err = transformListFilter(v, k, t, index, req)
			}

		} else {
			m := extraMapping[k]
			if m.FieldReqFunc == nil {
				(*req)[m.Field] = v
			} else {
				err = m.FieldReqFunc(v, m.Field, k, req)
			}
		}
	}
	return index, err
}

func requestUpdateMapping(d *schema.ResourceData, k string, t SdkReqTransform, index int, extraMapping map[string]SdkRequestMapping, req *map[string]interface{}) (int, error) {
	var err error

	if t.forceUpdateParam || (d.HasChange(k) && !d.IsNewResource()) {
		index, err = requestCreateMapping(d, k, t, index, extraMapping, req, true)
	}
	return index, err
}

func SdkResponseDefault(p string, d interface{}, item *interface{}) {
	path := strings.Split(p, ".")

	if m, ok := (*item).(map[string]interface{}); ok {
		root := m
		for i, s := range path {
			if i < len(path)-1 {
				if v1, ok := root[s]; ok {
					if v2, ok := v1.(map[string]interface{}); ok {
						root = v2
					} else {
						break
					}
				} else {
					break
				}
			} else {
				if _, ok := root[s]; !ok {
					(root)[s] = d
				}
			}
		}
	}
}

func SdkResponseAutoResourceData(d *schema.ResourceData, resource *schema.Resource, item interface{}, extra map[string]SdkResponseMapping, start ...bool) interface{} {
	setFlag := false
	if start == nil || (len(start) > 0 && start[0]) {
		setFlag = true
	}
	if reflect.ValueOf(item).Kind() == reflect.Map {
		result := make(map[string]interface{})
		root := item.(map[string]interface{})
		for k, v := range root {
			var value interface{}
			var err error
			//if v == nil {
			//	continue
			//}
			//if str, ok := v.(string); ok {
			//	if str == "" {
			//		continue
			//	}
			//}
			m := SdkResponseMapping{}
			target := Hump2Downline(k)
			if _, ok := extra[k]; ok {
				m = extra[k]
				target = m.Field
			}
			if r, ok := resource.Schema[target]; ok {
				if r.Elem != nil {
					if elem, ok := r.Elem.(*schema.Resource); ok {
						if m.FieldRespFunc != nil {
							value = SdkResponseAutoResourceData(d, elem, m.FieldRespFunc(v), extra, false)
						} else {
							value = SdkResponseAutoResourceData(d, elem, v, extra, false)
						}
					} else if _, ok := r.Elem.(*schema.Schema); ok {
						if m.FieldRespFunc != nil {
							value = m.FieldRespFunc(v)
						} else {
							value = v
						}
					}
				} else {
					if m.FieldRespFunc != nil {
						value = m.FieldRespFunc(v)
					} else {
						value = v
					}
				}
			} else {
				continue
			}
			if setFlag {
				if (resource.Schema[target].Type == schema.TypeList ||
					resource.Schema[target].Type == schema.TypeSet) &&
					reflect.ValueOf(value).Kind() == reflect.Map {
					err = d.Set(target, []interface{}{value})
				} else {
					err = d.Set(target, value)
				}

				if err != nil {
					log.Println(err.Error())
					panic("ERROR: " + err.Error())
				}

			} else {
				result[target] = value
			}
		}
		if len(result) > 0 {
			return result
		}
	} else if reflect.ValueOf(item).Kind() == reflect.Slice {
		var result []interface{}
		result = []interface{}{}
		root := item.([]interface{})
		for _, v := range root {
			value := SdkResponseAutoResourceData(d, resource, v, extra, false)
			result = append(result, value)
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func SdkResponseAutoMapping(resource *schema.Resource, collectField string, item map[string]interface{}, computeItem map[string]interface{},
	extraMapping map[string]SdkResponseMapping) map[string]interface{} {
	var result map[string]interface{}
	result = make(map[string]interface{})
	keys := strings.Split(collectField, ".")
	var extra map[string]interface{}
	extra = make(map[string]interface{})
	if len(keys) == 0 {
		return result
	}

	if computeItem != nil {
		for k, v := range computeItem {
			item[k] = v
		}
	}

	if _, ok := resource.Schema[keys[0]]; ok {
		elem := getSchemeElem(resource, keys)
		for k, v := range item {
			needExtraMapping := false
			target := Hump2Downline(k)
			m := SdkResponseMapping{}
			if extraMapping != nil {
				if _, ok := extraMapping[k]; ok {
					m = extraMapping[k]
					target = m.Field
					needExtraMapping = true
				}
			}
			if targetValue, ok := elem.Schema[target]; ok {
				if !needExtraMapping && (targetValue.Type == schema.TypeList || targetValue.Type == schema.TypeSet) {
					if _, ok := targetValue.Elem.(*schema.Schema); ok {
						extra[target] = v
					} else {
						if _, ok := extra[target]; !ok {
							if l, ok := v.([]interface{}); ok {
								_, result, _ := SdkSliceMapping(nil, l, SdkSliceData{
									SliceMappingFunc: func(m1 map[string]interface{}) map[string]interface{} {
										return SdkResponseAutoMapping(resource, collectField+"."+target, m1, computeItem, extraMapping)
									},
								})
								extra[target] = result
							} else if m, ok := v.(map[string]interface{}); ok {
								result, _ := SdkMapMapping(m, SdkSliceData{
									SliceMappingFunc: func(m1 map[string]interface{}) map[string]interface{} {
										return SdkResponseAutoMapping(resource, collectField+"."+target, m1, computeItem, extraMapping)
									},
								})
								extra[target] = []map[string]interface{}{
									result,
								}

							}
						}
					}
				}
			} else {
				continue
			}

			needDefaultMapping := false
			if _, ok := extra[target]; !ok {
				needDefaultMapping = true
			}
			if needDefaultMapping {
				if needExtraMapping {
					if m.FieldRespFunc == nil {
						result[m.Field] = v
					} else {
						result[m.Field] = m.FieldRespFunc(v)
					}
					if m.KeepAuto {
						result[Hump2Downline(k)] = result[m.Field]
					}
				} else {
					result[target] = v
				}
			} else {
				result[target] = extra[target]
			}
		}
	}
	return result
}

func SdkMapMapping(result interface{}, sdkSliceData SdkSliceData) (map[string]interface{}, error) {
	var data map[string]interface{}
	if reflect.TypeOf(result).Kind() == reflect.Map {
		if v, ok := result.(map[string]interface{}); ok {
			data = mapMapping(sdkSliceData, v)
		}
	}
	return data, nil
}

func SdkSliceMapping(d *schema.ResourceData, result interface{}, sdkSliceData SdkSliceData) ([]string, []map[string]interface{}, error) {
	var err error
	var ids []string
	ids = []string{}
	var data []map[string]interface{}
	data = []map[string]interface{}{}

	if reflect.TypeOf(result).Kind() == reflect.Slice {
		var length = 0
		if v, ok := result.([]map[string]interface{}); ok {
			length = len(v)
			for _, v1 := range v {
				ids, data = sliceMapping(ids, data, sdkSliceData, v1)
			}
		} else {
			root := result.([]interface{})
			length = len(root)
			for _, v2 := range root {
				ids, data = sliceMapping(ids, data, sdkSliceData, v2)
			}
		}

		if d != nil && sdkSliceData.TargetName != "" {
			d.SetId(hashStringArray(ids))
			_ = d.Set("total_count", length)
			err = d.Set(sdkSliceData.TargetName, data)
			if err != nil {
				return nil, nil, err
			}
			if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
				err = writeToFile(outputFile.(string), data)
				if err != nil {
					return nil, nil, err
				}
			}
		}

	}
	return ids, data, nil
}

func GetSdkParam(d *schema.ResourceData, params []string) map[string]interface{} {
	sdkParam := make(map[string]interface{})
	for _, v := range params {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			sdkParam[vv] = fmt.Sprintf("%v", v1)
		}
	}
	return sdkParam
}

func OtherErrorProcess(remain *int, error error) *resource.RetryError {
	*remain--
	if *remain >= 0 {
		return resource.RetryableError(error)
	} else {
		return resource.NonRetryableError(error)
	}
}

func ModifyProjectInstanceNew(resourceId string, param *map[string]interface{}, client *KsyunClient) error {
	if projectId, ok := (*param)["ProjectId"]; ok {
		req := make(map[string]interface{})
		req["InstanceId"] = resourceId
		req["ProjectId"] = projectId
		conn := client.iamconn
		_, err := conn.UpdateInstanceProjectId(&req)
		if err != nil {
			return err
		}
		delete(*param, "ProjectId")
	}
	return nil
}

func ModifyProjectInstance(resourceId string, param *map[string]interface{}, meta interface{}) error {
	if projectId, ok := (*param)["ProjectId"]; ok {
		req := make(map[string]interface{})
		req["InstanceId"] = resourceId
		req["ProjectId"] = projectId
		client := meta.(*KsyunClient)
		conn := client.iamconn
		_, err := conn.UpdateInstanceProjectId(&req)
		if err != nil {
			return err
		}
		delete(*param, "ProjectId")
	}
	return nil
}

func checkValueInSliceMap(data []interface{}, key string, value interface{}) (c bool) {
	for _, obj := range data {
		if m, ok := obj.(map[string]interface{}); ok {
			if v, ok1 := m[key]; ok1 && v == value {
				c = true
				return c
			}
		}
	}
	return c
}
