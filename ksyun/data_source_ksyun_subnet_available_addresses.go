/*
This data source provides a list of subnet available IPs.

# Example Usage

```hcl

	data "ksyun_subnet_available_addresses" "default" {
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

func dataSourceKsyunSubnetAvailableAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSubnetAvailableAddressesRead,
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
				Description: "Total number of available IPs that satisfy the condition.",
			},

			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of subnet IDs.",
			},
			"subnet_available_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of available IPs.",
			},
		},
	}
}

func dataSourceKsyunSubnetAvailableAddressesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).vpcconn
	req := make(map[string]interface{})
	var SubnetAvailableAddresseIds []string

	if ids, ok := d.GetOk("ids"); ok {
		SubnetAvailableAddresseIds = SchemaSetToStringSlice(ids)
	}
	if len(SubnetAvailableAddresseIds) == 0 {
		if ids, ok := d.GetOk("subnet_id"); ok {
			SubnetAvailableAddresseIds = SchemaSetToStringSlice(ids)
		}
	}
	if len(SubnetAvailableAddresseIds) > 0 {
		req["Filter.1.Name"] = "subnet-id"
	}

	for k, v := range SubnetAvailableAddresseIds {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("Filter.1.Value.%d", k+1)] = v
	}
	resp, err := conn.DescribeSubnetAvailableAddresses(&req)
	if err != nil {
		return fmt.Errorf("error on reading SubnetAvailableAddresse list(%v) %s", req, err)
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

	var datas []string
	for _, v := range items {
		datas = append(datas, v.(string))
	}
	err = dataSourceKscSaveSlice(d, "subnet_available_addresses", SubnetAvailableAddresseIds, datas)
	if err != nil {
		return fmt.Errorf("error on save SubnetAvailableAddresse list, %s", err)
	}
	return nil
}
