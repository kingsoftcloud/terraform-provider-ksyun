package helper

import (
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// GetSchemaListHeadMap returns string key map if argument is MaxItem: 1 List Type
func GetSchemaListHeadMap(d *schema.ResourceData, key string) (result map[string]interface{}, ok bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return
	}
	interfaces, ok := v.([]interface{})
	if !ok || len(interfaces) == 0 {
		ok = false
		return
	}
	head := interfaces[0]
	result, ok = head.(map[string]interface{})
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
		}
		rm[Underline2Hump(ck)] = cv
		// delete(m, ck)
	}
	return rm
}
func ConvertMapKey2Underline(m map[string]interface{}) map[string]interface{} {
	rm := make(map[string]interface{}, len(m))
	for ck, cv := range m {
		rm[Hump2Underline(ck)] = cv
	}
	return rm
}

func IsEmpty(v interface{}) bool {
	value := reflect.ValueOf(v)
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
