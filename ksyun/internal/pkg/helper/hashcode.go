package helper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// generate a hash code that's value if from 'backend_server_group_id'
// and its type is  schema.SchemaSetFunc
func HashFuncWithKeys(keys ...string) schema.SchemaSetFunc {
	return func(v interface{}) int {
		if v == nil {
			return hashcode.String("")
		}
		m := v.(map[string]interface{})
		var buf bytes.Buffer

		for _, key := range keys {
			if vv, ok := m[key]; ok {
				switch vv.(type) {
				case string:
					if vv == "" {
						break
					}
					buf.WriteString(fmt.Sprintf("%s", strings.ToLower(vv.(string))))
				case int:
					buf.WriteString(fmt.Sprintf("%d", vv.(int)))
				case float64:
					buf.WriteString(fmt.Sprintf("%d", int(vv.(float64))))
				case bool:
					if vv.(bool) {
						buf.WriteString("1")
					} else {
						buf.WriteString("0")
					}
				case []interface{}:
					vvs := vv.([]interface{})
					if len(vvs) < 1 {
						buf.WriteString("[]")
					} else {
						buf.WriteString("[")
						for _, vvv := range vvs {
							buf.WriteString(fmt.Sprintf("%v", vvv))
						}
						buf.WriteString("]")
					}
				case map[string]interface{}:
					vvs := vv.([]interface{})
					if len(vvs) < 1 {
						buf.WriteString("{}")
					} else {
						buf.WriteString("{")
						for _, vvv := range vvs {
							buf.WriteString(fmt.Sprintf("%v", vvv))
						}
						buf.WriteString("}")
					}
				}
			}
		}
		return hashcode.String(buf.String())
	}
}
