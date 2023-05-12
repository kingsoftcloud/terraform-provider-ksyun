/*
This data source provides a list of allocated IPs.

# Example Usage

```hcl

	data "ksyun_subnet_allocated_ip_addresses" "default" {
	  output_file="output_result"
	  ids=["494c3a64-eff9-4438-aa7c-694b7baxxxxx"]
	  subnet_id=["494c3a64-eff9-4438-aa7c-694b7baxxxxx"]
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunSubnetAllocatedIpAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSubnetAllocatedIpAddressessRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of subnet IDs.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Allocated IPs that satisfy the condition.",
			},

			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID of the subnet.",
			},
			"subnet_allocated_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of allocated IPs.",
			},
		},
	}
}

func dataSourceKsyunSubnetAllocatedIpAddressessRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).vpcconn
	req := make(map[string]interface{})
	var SubnetAllocatedIpAddressesIds []string
	if ids, ok := d.GetOk("ids"); ok {
		SubnetAllocatedIpAddressesIds = SchemaSetToStringSlice(ids)
	} else {
		return fmt.Errorf("subnet_id can not be empty")
	}
	if len(SubnetAllocatedIpAddressesIds) < 1 {
		return fmt.Errorf("subnet_id can not be empty")
	}
	req["Filter.1.Name"] = "subnet-id"
	for k, v := range SubnetAllocatedIpAddressesIds {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("Filter.1.Value.%d", k+1)] = v
	}

	resp, err := conn.DescribeSubnetAllocatedIpAddresses(&req)
	if err != nil {
		return fmt.Errorf("error on reading SubnetAllocatedIpAddresses list, %s", req)
	}
	itemSet, ok := (*resp)["AvailableIpAddress"]
	if !ok {
		return fmt.Errorf("error on reading SubnetAvailableAddresse set")
	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	//d.Set("subnet_allocated_ip_addresses",items)
	var datas []string
	for _, v := range items {
		datas = append(datas, v.(string))
	}
	err = dataSourceKscSaveSlice(d, "subnet_allocated_ip_addresses", SubnetAllocatedIpAddressesIds, datas)
	if err != nil {
		return fmt.Errorf("error on save SubnetAllocatedIpAddresses list, %s", err)
	}
	return nil
}
