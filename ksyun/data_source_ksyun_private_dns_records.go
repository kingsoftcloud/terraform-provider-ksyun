/*
This data source provides a list of Private Dns Record resources according to their Zone ID.

# Example Usage

```hcl

data "ksyun_private_dns_records" "default" {
  output_file = "pdns_records_output_result"
  zone_id = "a5ae6bf0-0ff4-xxxxxx-xxxxx-xxxxxxxxxx"
  region_name = ["cn-beijing-6"]
  record_ids = []
}

```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunPrivateDnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunPrivateDnsRecordsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Id of the private dns zone. Required.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"record_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.StringIsNotWhiteSpace,
					),
				},
				Set:         schema.HashString,
				Description: "A list of Record IDs, the Records belong to this private-dns-zone. The value of id is not be `\"\"`.",
			},

			"region_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						validation.StringIsNotWhiteSpace,
					),
				},
				Set:         schema.HashString,
				Description: "A list of the filter values that is region name. Such `cn-beijing-6`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Private Dns Record that satisfy the condition.",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of `Private Dns Records`. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the record.",
						},
						"record_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of record.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},

						"record_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The record ttl.",
						},
						"record_data_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The record value and other information like priority and weight etc.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"record_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Record value.",
									},
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The priority of record.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The weight of record.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The port of record.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunPrivateDnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	s := DnsService{meta.(*KsyunClient)}
	return s.ReadAndSetPrivateDnsRecords(d, dataSourceKsyunPrivateDnsRecords())
}
