/*
Query ksyun dnats information

# Example Usage

```hcl

data "ksyun_dnats" "default" {
  private_ip_address = "10.7.x.xxx"
  nat_id             = "5c7b7925-xxxx-xxxx-xxxx-434fc8042329"
  dnat_ids           = ["5cxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx"]
  output_file        = "output_result"
}

```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunDnats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunDnatsRead,
		Schema: map[string]*schema.Schema{
			"dnat_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The id list of dnats.",
			},
			"nat_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The nat ip.",
			},
			"public_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The public port.",
			}, // query data guard
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private ip address.",
			},
			"ip_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol of dnat rule.",
			},
			"nat_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The nat id of dnat associated.",
			},
			"dnat_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of dnat.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network interface id of dnat rule associated.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of snapshot policies resources that satisfy the condition.",
			},

			"dnats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of krds db parameter groups. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// return values by data source query
						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The nat id.",
						},
						"dnat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of dnat.",
						},
						"dnat_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of dnat.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of dnat.",
						},
						"ip_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ip protocol of dnat.",
						},
						"nat_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The nat ip of nat associated.",
						},
						"public_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public port.",
						},
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private ip address.",
						},
						"private_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The private port.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time created.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunDnatsRead(d *schema.ResourceData, meta interface{}) error {
	params := NewDescribeDnatsParams()
	if dnatId, ok := d.GetOk("dnat_ids"); ok {
		params.DnatIds = dnatId.([]interface{})
	}

	if dnatName, ok := d.GetOk("dnat_name"); ok {
		params.Filter.DnatName = dnatName.(string)
	}

	if publicPort, ok := d.GetOk("public_port"); ok {
		params.Filter.PublicPort = publicPort.(string)
	}

	if infraIp, ok := d.GetOk("private_ip_address"); ok {
		params.Filter.PrivateIpAddress = infraIp.(string)
	}

	if ipProtocol, ok := d.GetOk("ip_protocol"); ok {
		params.Filter.IpProtocol = ipProtocol.(string)
	}
	if natId, ok := d.GetOk("nat_id"); ok {
		params.Filter.NatId = natId.(string)
	}
	if kniId, ok := d.GetOk("network_interface_id"); ok {
		params.Filter.NetworkInterfaceId = kniId.(string)
	}
	if kNatIp, ok := d.GetOk("nat_ip"); ok {
		params.Filter.NatIp = kNatIp.(string)
	}

	vpcSrv := VpcService{
		client: meta.(*KsyunClient),
	}
	sdkResponse, err := vpcSrv.DescribeDnats(params)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, dataSourceKsyunDnats(), ksyunDataSource{
		collection:  sdkResponse,
		idFiled:     "DnatId",
		targetField: "dnats",
		extra:       nil,
	})
}
