package helper

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
