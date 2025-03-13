package helper

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
)

// GetSchemaListHeadMap returns string key map if argument is MaxItem: 1 List Type
func GetSchemaListHeadMap(d *schema.ResourceData, key string) (result map[string]interface{}, ok bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return
	}
	var interfaces []interface{}
	switch v.(type) {
	case []interface{}:
		interfaces = v.([]interface{})
	case *schema.Set:
		interfaces = v.(*schema.Set).List()
	}

	if !ok || len(interfaces) == 0 {
		ok = false
		return
	}
	head := interfaces[0]
	result, ok = head.(map[string]interface{})
	return
}

func GetSchemaSingerMapWithKey(d *schema.ResourceData, key string) (result map[string]interface{}, ok bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return
	}
	result, ok = v.(map[string]interface{})
	return
}

func GetSchemaMapListWithKey(d *schema.ResourceData, key string) (result []map[string]interface{}, ok bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return
	}
	var interfaces []interface{}
	switch v.(type) {
	case []interface{}:
		interfaces = v.([]interface{})
	case *schema.Set:
		interfaces = v.(*schema.Set).List()
	}

	if !ok || len(interfaces) == 0 {
		ok = false
		return
	}
	for _, inter := range interfaces {
		if v, ok := inter.(map[string]interface{}); !ok {
			return nil, ok
		} else {
			result = append(result, v)
		}
	}
	return
}

func StringBoolean(s bool) string {
	if s {
		return "True"
	}
	return "False"
}

func GetDWithBool(d *schema.ResourceData, key string) (b bool, ok bool) {
	if v, dOk := d.GetOk(key); dOk {
		b = v.(bool)
		return b, dOk
	} else {
		return false, false
	}
}

// Underline2Hump The underline is converted to a hump simply.
func Underline2Hump(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return strings.Title(s)
	}
	var s1 []string
	ss := strings.Split(s, "_")
	for _, v := range ss {
		vv := strings.ToUpper(v[:1]) + v[1:]
		s1 = append(s1, vv)
	}
	return strings.Join(s1, "")
}

// Hump2Underline The hump is converted to an underline simply, and no special treatment is required for even uppercase letters.
//
//ex:aDDCC ->a_d_d_c_c
func Hump2Underline(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return strings.ToLower(s)
	}
	var s1 string
	if len(s) == 1 {
		s1 = strings.ToLower(s[:1])
		return s1
	}
	for k, v := range s {
		if k == 0 {
			s1 = strings.ToLower(s[0:1])
			continue
		}
		if v >= 65 && v <= 90 {
			v1 := "_" + strings.ToLower(s[k:k+1])
			s1 = s1 + v1
		} else {
			s1 = s1 + s[k:k+1]
		}
	}
	return s1
}

func ConvertMapKey2Title(m map[string]interface{}, hitBlankStr bool) map[string]interface{} {
	rm := make(map[string]interface{}, len(m))
	for ck, cv := range m {
		switch cv.(type) {
		case string:
			if hitBlankStr && cv == "" {
				continue
			}
		case map[string]interface{}:
			cv = ConvertMapKey2Title(cv.(map[string]interface{}), hitBlankStr)
		case []interface{}:
			for i, v := range cv.([]interface{}) {
				switch v.(type) {
				case map[string]interface{}:
					cv.([]interface{})[i] = ConvertMapKey2Title(v.(map[string]interface{}), hitBlankStr)
				}
			}
		}
		rm[Underline2Hump(ck)] = cv
		// delete(m, ck)
	}
	return rm
}

func ConvertMapKey2Underline(m map[string]interface{}) map[string]interface{} {
	rm := make(map[string]interface{}, len(m))
	for ck, cv := range m {
		switch cv.(type) {
		case map[string]interface{}:
			cv = ConvertMapKey2Underline(cv.(map[string]interface{}))
		case []interface{}:
			for i, v := range cv.([]interface{}) {
				switch v.(type) {
				case map[string]interface{}:
					cv.([]interface{})[i] = ConvertMapKey2Underline(v.(map[string]interface{}))
				}
			}
		}
		rm[Hump2Underline(ck)] = cv
	}
	return rm
}

func IsEmpty(v interface{}) bool {
	value := reflect.ValueOf(v)
	if !value.IsValid() {
		return true
	}
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return true
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Slice, reflect.Map:
		return value.Len() == 0
	}
	zeroValue := reflect.Zero(value.Type())
	return reflect.DeepEqual(v, zeroValue.Interface())
}

func GetSchemaListWithString(d *schema.ResourceData, k string) (l []string, ok bool) {
	v, ok := d.GetOk(k)
	if !ok {
		return
	}
	interfaces, ok := v.([]interface{})
	if !ok || len(interfaces) == 0 {
		ok = false
		return
	}
	l = make([]string, len(interfaces))
	for i, inter := range interfaces {
		if v, ok := inter.(string); !ok {
			return nil, ok
		} else {
			l[i] = v
		}
	}
	return
}

func MapCopy(m map[string]interface{}) (map[string]interface{}, bool) {
	if m == nil {
		return nil, false
	}
	n := make(map[string]interface{}, len(m))
	for k, v := range m {
		n[k] = v
	}
	return n, true
}

func GetDiffMap(base map[string]interface{}, targets ...map[string]interface{}) (diff map[string]interface{}) {
	if len(targets) == 0 {
		return base
	}
	var isDiff bool

	diff = make(map[string]interface{}, len(base))
	for _, target := range targets {
		diff, isDiff = diffMap(base, target)
		if isDiff {
			return diff
		}
	}
	return
}

func diffMap(base map[string]interface{}, target map[string]interface{}) (diff map[string]interface{}, isDiff bool) {
	diff = make(map[string]interface{})
	for k, v := range base {
		if tv, ok := target[k]; ok {
			if !reflect.DeepEqual(v, tv) {
				isDiff = true
			}

			diff[k] = tv
		} else {
			diff[k] = v
		}
	}
	return
}

func MapstructureFiller(i interface{}, o interface{}, tag string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("MapstructureFiller panic recovered: ", r)
			err = fmt.Errorf("MapstructureFiller panic recovered: %v", r)
		}
	}()

	oVal := reflect.ValueOf(o)
	if oVal.Kind() != reflect.Ptr {
		return fmt.Errorf("o must be a pointer")
	}
	// it's basic convert
	convertConfig := &mapstructure.DecoderConfig{
		IgnoreUntaggedFields: true,
		ZeroFields:           false,
		WeaklyTypedInput:     true,
		Metadata:             nil,
		Result:               o,
	}

	// convertConfig.DecodeHook = decodeHookFunc()

	if tag != "" {
		convertConfig.TagName = tag
	}

	decoder, err := mapstructure.NewDecoder(convertConfig)
	if err != nil {
		return err
	}
	if err := decoder.Decode(i); err != nil {
		return err
	}

	if tag == "" {
		return nil
	}

	sVal := reflect.ValueOf(i)
	switch sVal.Kind() {
	case reflect.Struct:
		oMap := o.(*map[string]interface{})
		// specify to custom type such as list, filter
		for j := 0; j < sVal.NumField(); j++ {
			fieldType := sVal.Type().Field(j)
			fieldVal := sVal.Field(j)

			mapstructureTag := fieldType.Tag.Get(tag)
			tagList := strings.Split(mapstructureTag, ",")
			var (
				tagKind string
				tagName string
			)
			tagName = tagList[0]
			if tagName == "" {
				continue
			}

			if isZeroOfUnderlyingType(fieldVal.Interface()) {
				delete(*oMap, tagName)
				continue
			}

			switch fieldType.Type.Kind() {
			case reflect.Slice:
				if fieldVal.Len() == 0 {
					break
				}
				tempList := make([]interface{}, 0, fieldVal.Len())
				for k := 0; k < fieldVal.Len(); k++ {
					if fieldVal.Index(k).Kind() == reflect.Struct {
						tempMap := make(map[string]interface{})

						err = MapstructureFiller(fieldVal.Index(k).Interface(), &tempMap, tag)
						if err != nil {
							return err
						}
						tempList = append(tempList, tempMap)
					}
				}
				if len(tempList) > 0 {
					(*oMap)[tagName] = tempList
				}
			}

			if len(tagList) > 1 {
				tagKind = tagList[1]
			}

			switch tagKind {
			case "tf-list":
				if IsEmpty((*oMap)[tagName]) || reflect.ValueOf((*oMap)[tagName]).Kind() == reflect.Slice {
					continue
				}
				(*oMap)[tagName] = []interface{}{(*oMap)[tagName]}
			}
		}
	}

	return nil
}

func decodeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() == reflect.String && t.Kind() == reflect.Bool {
			return data, nil
		}

		result := make(map[string]interface{})

		// 遍历struct的所有字段
		val := reflect.ValueOf(data)
		for i := 0; i < val.NumField(); i++ {
			valueField := val.Field(i)
			typeField := val.Type().Field(i)

			// 检查字段是否为空值（这里只检查了基本类型和字符串，你可能需要扩展以处理其他类型）
			if isZeroOfUnderlyingType(valueField.Interface()) {
				continue // 忽略空值
			}

			// 使用字段名作为map的键，字段值作为map的值
			result[typeField.Name] = valueField.Interface()
		}

		return result, nil
	}
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
