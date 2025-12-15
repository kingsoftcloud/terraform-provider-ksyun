package ksyun

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// 重要提示
// 拦截器 返回false代表不拦截变更 反正则拦截变更
// 不要使用isNewResource() 判断是否是新建资源还是更新资源，因为这是内置机制，新建资源时候isNewResource()也是false 需要用d.id()替代

func purchaseTimeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("charge_type"); ok && (v.(string) == "Monthly" || v.(string) == "PrePaidByMonth") {
		return false
	}
	return true
}

func purchaseTimeTrialDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("trial").(bool) {
		return false
	}
	return true
}

func chargeSchemaDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	mappings := map[string]string{
		"PostPaidByPeak":     "Peak",
		"PostPaidByDay":      "Daily",
		"PostPaidByTransfer": "TrafficMonthly",
		"PrePaidByMonth":     "Monthly",
		"Peak":               "PostPaidByPeak",
		"Daily":              "PostPaidByDay",
		"TrafficMonthly":     "PostPaidByTransfer",
		"Monthly":            "PrePaidByMonth",
	}
	if old == new {
		return true
	}
	if v, ok := mappings[old]; ok && v == new {
		return true
	}
	return false
}

func kecDiskSnapshotIdDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// logger.Debug("test", "test", d.Id(), k, strings.Contains(k, "disk_snapshot_id"))
	if d.Id() != "" {
		if strings.Contains(k, "disk_snapshot_id") {
			// logger.Debug("test1", "test", 123)
			return true
		}
	}
	return false
}

func kecImportDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// 由于一些字段暂时无法支持从查询中返回 所以现在设立做特殊处理拦截变更 用来适配导入的场景 后续支持后在对导入场景做优化
	if d.Id() != "" {
		if k == "local_volume_snapshot_id" {
			return true
		}
		if k == "user_data" {
			if d.HasChange("image_id") {
				return false
			}
			return true
		}
		if k == "auto_create_ebs" {
			return true
		}
	}
	if (k == "keep_image_login" || k == "key_id") && d.Id() != "" && !d.HasChange("image_id") {
		return true
	}

	return false
}

func kcsParameterDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if old != "" && new == "" {
		return true
	}
	if k == "parameters.notify-keyspace-events" && old == "" && new == "" {
		return true
	}
	return false
}

func rdsParameterDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if k == "parameters.#" {
		logger.Debug(logger.RespFormat, "DemoTest", d.ConnInfo())
		logger.Debug(logger.RespFormat, "DemoTest", d.Get("parameters"))
		return false
	}
	return true
}

func kcsSecurityGroupDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("security_group_ids") != nil && old == "" && new != "" {
		if sgs, ok := d.Get("security_group_ids").(*schema.Set); ok {
			if (*sgs).Contains(new) {
				err := d.Set("security_group_id", new)
				if err == nil {
					return true
				}
			}
		}
	}
	return false
}

func networkAclEntryDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("protocol") != "icmp" && (k == "icmp_type" || k == "icmp_code") {
		return true
	}
	if d.Get("protocol") != "tcp" && d.Get("protocol") != "udp" && (k == "port_range_from" || k == "port_range_to") {
		return true
	}
	return false
}

func securityGroupEntryDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("protocol") != "icmp" && (k == "icmp_type" || k == "icmp_code") {
		return true
	}
	if d.Get("protocol") != "tcp" && d.Get("protocol") != "udp" && (k == "port_range_from" || k == "port_range_to") {
		return true
	}
	return false
}

func loadBalancerDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("type") != "internal" && (k == "subnet_id" || k == "private_ip_address") {
		return true
	}
	return false
}

func securityGroupEntryLiteDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if k == "cidr_block.#" {
		return false
	}
	oldBlockIf, _ := d.GetChange("cidr_block")
	oldBlock := oldBlockIf.([]interface{})
	newBlock, _ := helper.GetSchemaListWithString(d, "cidr_block")
	if len(oldBlock) != len(newBlock) {
		return false
	}

	for _, cidrBlock := range oldBlock {
		if !stringSliceContains(newBlock, cidrBlock.(string)) {
			return false
		}
	}

	return true
}

func albInternalDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("alb_type") != "internal" && (k == "subnet_id" || k == "private_ip_address") {
		return true
	}
	return false
}

func AlbListenerDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// if k == "certificate_id" || k == "tls_cipher_policy" || k == "enable_http2" {
	//	return true
	// }
	// if d.Get("listener_protocol") != "HTTP" && k == "redirect_listener_id" {
	//	return true
	// }
	// if d.Get("listener_protocol") != "HTTPS" && d.Get("listener_protocol") != "HTTP" &&
	//	(k == "http_protocol" ||
	//		k == "health_check.0.host_name" ||
	//		k == "health_check.0.url_path" ||
	//		k == "health_check.0.is_default_host_name" ||
	//		k == "session.0.cookie_type" ||
	//		k == "session.0.cookie_name") {
	//	return true
	// }
	if k == "session.0.cookie_name" && d.Get("session.0.cookie_type") != "RewriteCookie" {
		return true
	}
	// if k == "health_check.0.host_name" && d.Get("health_check.0.is_default_host_name").(bool) {
	//	return true
	// }
	return false
}

// AlbRuleGroupSyncOffDiffSuppressFunc find field difference when `listener_sync` is off
func AlbRuleGroupSyncOffDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("listener_sync").(string) == "on" {
		return true
	}
	switch k {
	case "cookie_type", "session_persistence_period":
		if n := d.Get("session_state"); n == "start" {
			return false
		}

	// health check
	case "interval", "timeout", "healthy_threshold", "unhealthy_threshold",
		"health_protocol", "health_port", "http_method", "url_path", "host_name":
		if n := d.Get("health_check_state"); n == "start" {
			switch k {
			case "http_method", "url_path", "host_name":
				if d.Get("health_protocol") == "HTTP" {
					return false
				}
				return true
			}
			return false
		}
	case "cookie_name":
		if n := d.Get("session_state"); n == "start" {
			if d.Get("cookie_type") == "RewriteCookie" {
				return false
			}
		}

	case "session_state", "health_check_state":
		return false

	}

	return true
}

func albRuleGroupDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	var (
		fieldKey    string
		resourceKey string
	)
	if strings.Contains(k, ".") {
		keys := strings.Split(k, ".")
		fieldKey = keys[2]
		resourceKey = strings.Join(keys[:2], ".") + ".alb_rule_type"
	}
	switch d.Get(resourceKey) {
	case "url", "domain":
		if fieldKey == "alb_rule_value" {
			return false
		}
	case "header":
		if fieldKey == "header_value" {
			return false
		}
	case "method":
		if fieldKey == "method_value" {
			return false
		}
	case "sourceIp":
		if fieldKey == "source_ip_value" {
			return false
		}
	case "query":
		if fieldKey == "query_value" {
			return false
		}
	case "cookie":
		if fieldKey == "cookie_value" {
			return false
		}
	}
	return true
}

func albRuleGroupTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	var (
		fieldKey    string
		resourceKey string
	)
	if strings.Contains(k, ".") {
		keys := strings.Split(k, ".")
		fieldKey = keys[0]
		resourceKey = "type"
	} else {
		resourceKey = "type"
		fieldKey = k
	}

	switch d.Get(resourceKey) {
	case albRuleTypeForwardGroup:
		if fieldKey == "backend_server_group_id" {
			return false
		}
		return true
	case albRuleTypeRewrite:
		switch fieldKey {
		case "backend_server_group_id", "rewrite_config":
			return false
		}
		return true

	case albRuleTypeFixedResponse:
		if fieldKey == "fixed_response_config" {
			return false
		}
		return true

	case albRuleTypeRedirect:
		switch fieldKey {
		case "redirect_alb_listener_id", "redirect_http_code":
			return false
		}
		return true
	}
	return false
}

func lbListenerDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("listener_protocol") != "HTTPS" && (k == "certificate_id" || k == "tls_cipher_policy" || k == "enable_http2") {
		return true
	}
	if d.Get("listener_protocol") != "HTTP" && k == "redirect_listener_id" {
		return true
	}
	if d.Get("listener_protocol") != "HTTPS" && d.Get("listener_protocol") != "HTTP" &&
		(k == "http_protocol" ||
			k == "health_check.0.host_name" ||
			k == "health_check.0.url_path" ||
			k == "health_check.0.is_default_host_name" ||
			k == "session.0.cookie_type" ||
			k == "session.0.cookie_name") {
		return true
	}
	if k == "session.0.cookie_name" && d.Get("session.0.cookie_type") != "RewriteCookie" {
		return true
	}
	if k == "health_check.0.host_name" && d.Get("health_check.0.is_default_host_name").(bool) {
		return true
	}

	// deal with health_check excepted health_check_state
	// if health_check_state is stop, all health_check fields are suppressed
	if (strings.HasPrefix(k, "health_check") && !strings.HasSuffix(k, "health_check_state")) &&
		d.Get("health_check.0.health_check_state") == "stop" {
		return true
	}

	return false
}

func lbHealthCheckDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("listener_protocol") != "" && d.Get("listener_protocol") != "HTTP" && d.Get("listener_protocol") != "HTTPS" &&
		(k == "url_path" || k == "host_name" || k == "is_default_host_name") {
		return true
	}
	if d.Get("host_name") != "" && k == "is_default_host_name" {
		return true
	}
	if k == "host_name" && d.Get("is_default_host_name").(bool) {
		return true
	}
	return false
}

func lbRuleDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("listener_sync") != "off" && (k == "method" || strings.HasPrefix(k, "session.") || strings.HasPrefix(k, "health_check.")) {
		return true
	}
	if k == "session.0.cookie_name" && d.Get("session.0.cookie_type") != "RewriteCookie" {
		return true
	}
	if k == "health_check.0.host_name" && d.Get("health_check.0.is_default_host_name").(bool) {
		return true
	}
	return false
}

func hostHeaderDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("listener_protocol") != "" && d.Get("listener_protocol") != "HTTPS" && k == "certificate_id" {
		return true
	}
	return false
}

func lbRealServerDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("real_server_type") != "host" && k == "instance_id" {
		return true
	}
	if d.Get("listener_method") != "" && d.Get("listener_method") != "MasterSlave" && k == "master_slave_type" {
		return true
	}
	return false
}

func lbBackendServerDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("backend_server_group_type") != "Mirror" && strings.HasPrefix(k, "health_check.") {
		return true
	}
	return false
}

func lbBackendServerHealthCheckDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasPrefix(k, "health_check.") && d.Get("health_check.0.health_protocol") == "TCP" {
		keyNames := strings.Split(k, ".")
		keyName := keyNames[len(keyNames)-1]
		switch keyName {
		case "health_code", "host_name", "url_path", "http_method":
			return true
		}
	}
	return false
}

func volumeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() != "" && d.HasChange("size") && k == "online_resize" {
		return false
	}
	return true
}

func bareMetalDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("network_interface_mode") != "dual" && strings.HasPrefix(k, "extension_") {
		return true
	}
	if d.Get("network_interface_mode") != "bond4" && k == "bond_attribute" {
		return true
	}
	if (d.Id() == "" || d.Get("host_type") != "COLO") && (k == "server_ip" || k == "path") {
		return true
	}
	// if d.Id() == "" && (k == "host_status" || k == "force_re_install") {
	// 	return true
	// }
	return false
}

func bareMetalReinstallDiffSuppressFunc(k, oldV, newV string, d *schema.ResourceData) bool {
	if d.Id() == "" || (d.HasChange("force_re_install") && d.Get("force_re_install").(bool)) {
		if helper.StringInSlice(k, []string{"gpu_image_driver_id", "overclocking_attribute", "kmr_agent", "kes_agent", "computer_name", "container_agent", "password_inherit", "data_disk_mount"}) {
			return false
		}
	}

	return true
}

func bareMetalRoceNetwork(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() != "" {
		if k == "roce_network" {
			return true
		}
	}
	return false
}

func bareMetalCreateDiff(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() != "" {
		switch k {
		case "roce_network", "trial":
			return true
		}
	}
	return false
}

func activateHotStandbyDSF(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("host_status") == "HotStandbyToBeActivated" && d.Get(k).(bool) {
		return false
	}
	return true
}

func vpnV2ParamsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	isV2 := d.Get("vpn_gateway_version") == "2.0"
	// is vpn version 1
	if stringSliceContains(vpnV1Attribute, k) {
		if isV2 {
			return true
		}
	} else if stringSliceContains(vpnV2Attribute, k) {
		if !isV2 {
			return true
		}
		switch k {
		case "local_peer_ip", "customer_peer_ip":
			if d.Get("type") != "RouteIpsec" {
				return true
			}
		}
	}
	return false
}

func pdnsZoneRecordDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	t := d.Get("type")
	if !(t == "SRV" || t == "MX") {
		return true
	}

	switch k {
	case "weight", "port":
		if t != "SRV" {
			return true
		}
	}
	return false
}
