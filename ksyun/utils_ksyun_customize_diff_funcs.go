package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
)

func kecNetworkInterfaceCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {
	if d.Id() != "" && (d.HasChange("private_ip_address") || d.HasChange("subnet_id") || d.HasChange("security_group_ids")) {
		var data []interface{}
		vpcService := VpcService{meta.(*KsyunClient)}
		condition := map[string]interface{}{
			"NetworkInterfaceId.1": d.Id(),
		}
		data, err = vpcService.ReadNetworkInterfaces(condition)
		if err != nil {
			return err
		}
		if len(data) != 1 {
			return fmt.Errorf("NetworkInterface %s not exist ", d.Id())
		}
		if _, ok := data[0].(map[string]interface{})["InstanceId"]; ok {
			return err
		}
		if d.HasChange("private_ip_address") {
			err = d.ForceNew("private_ip_address")
			if err != nil {
				return err
			}
		}
		if d.HasChange("subnet_id") {
			err = d.ForceNew("subnet_id")
			if err != nil {
				return err
			}
		}
		if d.HasChange("security_group_ids") {
			err = d.ForceNew("security_group_ids")
			if err != nil {
				return err
			}
		}
	}
	return err
}

func scalingConfigurationCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {
	for _, v := range d.Get("data_disks").([]interface{}) {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		diskType, _ := m["disk_type"].(string)
		size, _ := m["disk_size"].(int)
		if size == 0 {
			continue
		}
		r, ok := kecDataDiskSizeRangeMap[diskType]
		if !ok {
			// disk_type not set (computed): fall back to the global envelope of all disk types.
			if size < 10 || size > 65536 {
				return fmt.Errorf(
					"the disk_size of data_disks must be in range [10, 65536] when disk_type is not set, got: %d",
					size)
			}
			continue
		}
		if size < r[0] || size > r[1] {
			return fmt.Errorf(
				"the disk_size of data_disks with disk_type %q must be in range [%d, %d], got: %d",
				diskType, r[0], r[1], size)
		}
	}
	return err
}

func networkAclEntryCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {
	if d.HasChange("network_acl_entries") {
		m := make(map[string]interface{})
		for _, v := range d.Get("network_acl_entries").(*schema.Set).List() {
			num := v.(map[string]interface{})["rule_number"].(int)
			direction := v.(map[string]interface{})["direction"].(string)
			if _, ok := m[strconv.Itoa(num)+"-"+direction]; ok {
				return fmt.Errorf("rule_number must unique")
			} else {
				m[strconv.Itoa(num)] = ""
			}
		}
	}
	return err
}

func bareMetalCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {

	if d.Get("network_interface_mode") != "dual" {
		if d.Get("extension_subnet_id") != "" ||
			d.Get("extension_private_ip_address") != "" ||
			len(d.Get("extension_security_group_ids").(*schema.Set).List()) > 0 ||
			d.Get("extension_dns1") != "" ||
			d.Get("extension_dns2") != "" {
			return fmt.Errorf("extension network must set empty when network_interface_mode is bond4 or single")
		}
	}

	if d.Id() != "" && d.HasChange("network_interface_mode") {
		o, n := d.GetChange("network_interface_mode")
		var flag bool
		if o == "bond4" && n == "single" {
			flag = true
		}
		if o == "single" && n == "bond4" {
			flag = true
		}
		if o == "dual" && n == "single" {
			flag = true
			if d.Get("extension_subnet_id") != "" ||
				d.Get("extension_private_ip_address") != "" ||
				len(d.Get("extension_security_group_ids").(*schema.Set).List()) > 0 ||
				d.Get("extension_dns1") != "" ||
				d.Get("extension_dns2") != "" {
				return fmt.Errorf("extension network must set empty when network_interface_mode is bond4 or single")
			}
		}
		if o == "dual" && n == "bond4" {
			flag = true
			if d.Get("extension_subnet_id") != "" ||
				d.Get("extension_private_ip_address") != "" ||
				len(d.Get("extension_security_group_ids").(*schema.Set).List()) > 0 ||
				d.Get("extension_dns1") != "" ||
				d.Get("extension_dns2") != "" {
				return fmt.Errorf("extension network must set empty when network_interface_mode is boun4 or single")
			}
		}
		if !flag {
			return d.ForceNew("network_interface_mode")
		}
	}
	return err
}
