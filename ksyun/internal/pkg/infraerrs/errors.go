package infraerrs

import (
	"reflect"
	"strings"
)

// IsContains returns whether value is within object
func IsContains(obj interface{}, value interface{}) bool {
	vv := reflect.ValueOf(obj)
	if vv.Kind() == reflect.Ptr || vv.Kind() == reflect.Interface {
		if vv.IsNil() {
			return false
		}
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Slice:
		for i := 0; i < vv.Len(); i++ {
			if reflect.DeepEqual(value, vv.Index(i).Interface()) {
				return true
			}
		}
		return false
	case reflect.Map:
		s := vv.MapKeys()
		for i := 0; i < len(s); i++ {
			if reflect.DeepEqual(value, s[i].Interface()) {
				return true
			}
		}
		return false
	case reflect.String:
		ss := reflect.ValueOf(value)
		switch ss.Kind() {
		case reflect.String:
			return strings.Contains(vv.String(), ss.String())
		}
		return false
	default:
		return reflect.DeepEqual(obj, value)
	}
}
