package ksyun

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func kecNetworkInterfaceHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	return hashcode.String(m["network_interface_id"].(string))
}

func networkAclEntryHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	buf := networkAclEntryHashBase(m)
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["description"].(string))))
	return hashcode.String(buf.String())
}

func networkAclEntrySimpleHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	buf := networkAclEntryHashBase(m)
	return hashcode.String(buf.String())
}

func networkAclEntryHashBase(m map[string]interface{}) (buf bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%d-", m["rule_number"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["cidr_block"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["direction"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["rule_action"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["protocol"].(string))))
	buf.WriteString(fmt.Sprintf("%d-", m["icmp_type"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["icmp_code"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["port_range_from"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["port_range_to"].(int)))
	return buf
}

func securityGroupEntryHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	buf := securityGroupEntryHashBase(m, false)
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["description"].(string))))
	return hashcode.String(buf.String())
}

func securityGroupEntrySimpleHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	buf := securityGroupEntryHashBase(v, false)
	return hashcode.String(buf.String())
}

func securityGroupEntrySimpleHashWithHump(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	buf := securityGroupEntryHashBase(v, true)
	return hashcode.String(buf.String())
}

func securityGroupEntryHashBase(v interface{}, isHump bool) (buf bytes.Buffer) {
	strField := []string{
		"protocol",
		"direction",
		"cidr_block",
	}
	logger.Debug(logger.RespFormat, "Demo", v)
	protocol := ""
	if m, ok1 := v.(map[string]interface{}); ok1 {
		for _, s := range strField {
			if !isHump {
				if _, ok := m[s]; ok {
					buf.WriteString(fmt.Sprintf("%s:", strings.ToLower(m[s].(string))))
				}
				protocol = strings.ToLower(m["protocol"].(string))
			} else {
				if _, ok := m[Downline2Hump(s)]; ok {
					buf.WriteString(fmt.Sprintf("%s:", strings.ToLower(m[Downline2Hump(s)].(string))))
				}
				protocol = strings.ToLower(m["Protocol"].(string))
			}
		}
		intField := generateEntryField(protocol)
		for _, s := range intField {
			if !isHump {
				if _, ok := m[s]; ok {
					buf.WriteString(fmt.Sprintf("%d:", int64(m[s].(int))))
				}
			} else {
				if _, ok := m[Downline2Hump(s)]; ok {
					buf.WriteString(fmt.Sprintf("%d:", int64(m[Downline2Hump(s)].(float64))))
				}
			}
		}
	} else if d, ok2 := v.(*schema.ResourceData); ok2 {
		for _, s := range strField {
			if _, ok := d.GetOk(s); ok {
				buf.WriteString(fmt.Sprintf("%s:", strings.ToLower(d.Get(s).(string))))
			}
			protocol = strings.ToLower(d.Get("protocol").(string))
		}
		intField := generateEntryField(protocol)
		for _, s := range intField {
			buf.WriteString(fmt.Sprintf("%d:", int64(d.Get(s).(int))))
		}
	}
	return buf
}

func generateEntryField(protocol string) (fields []string) {
	if protocol == "icmp" {
		fields = []string{
			"icmp_type",
			"icmp_code",
		}
	} else if protocol == "tcp" || protocol == "udp" {
		fields = []string{
			"port_range_from",
			"port_range_to",
		}
	}
	return fields
}

func loadBalancerAclEntryHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	buf := loadBalancerAclEntryHashBase(m)
	buf.WriteString(fmt.Sprintf("%d-", m["rule_number"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["rule_action"].(string))))
	return hashcode.String(buf.String())
}

func loadBalancerAclEntrySimpleHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	buf := loadBalancerAclEntryHashBase(m)
	return hashcode.String(buf.String())
}

func loadBalancerAclEntryNumberHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d-", m["rule_number"].(int)))
	return hashcode.String(buf.String())
}

func loadBalancerAclEntryCidrHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["cidr_block"].(string))))
	return hashcode.String(buf.String())
}

func loadBalancerAclEntryHashBase(m map[string]interface{}) (buf bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["cidr_block"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["protocol"].(string))))
	return buf
}

func kceInstanceNodeHashFunc() schema.SchemaSetFunc {
	return func(v interface{}) int {
		if v == nil {
			return hashcode.String("")
		}
		m := v.(map[string]interface{})

		var buf bytes.Buffer

		hashKeys := []string{
			"instance_type",
			"image_id",
			"security_group_id",
			"subnet_id",
			"role",
		}
		for _, key := range hashKeys {
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
							buf.WriteString(fmt.Sprintf("%s", strings.ToLower(vvv.(string))))
						}
						buf.WriteString("]")
					}
				}
			}
		}

		return hashcode.String(buf.String())
	}
}
