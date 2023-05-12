package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
)

//重要提示
//拦截器 返回false代表不拦截变更 反正则拦截变更
// 不要使用isNewResource() 判断是否是新建资源还是更新资源，因为这是内置机制，新建资源时候isNewResource()也是false 需要用d.id()替代

func purchaseTimeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("charge_type"); ok && (v.(string) == "Monthly" || v.(string) == "PrePaidByMonth") {
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
	//logger.Debug("test", "test", d.Id(), k, strings.Contains(k, "disk_snapshot_id"))
	if d.Id() != "" {
		if strings.Contains(k, "disk_snapshot_id") {
			//logger.Debug("test1", "test", 123)
			return true
		}
	}
	return false
}

func kecImportDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	//由于一些字段暂时无法支持从查询中返回 所以现在设立做特殊处理拦截变更 用来适配导入的场景 后续支持后在对导入场景做优化
	if d.Id() != "" {
		if k == "local_volume_snapshot_id" {
			return true
		}
		if k == "user_data" {
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
	if d.Id() == "" && (k == "host_status" || k == "force_re_install") {
		return true
	}
	return false
}
