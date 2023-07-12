/*
This data source provides a list of Network Interface resources according to their Network Interface ID.

# Example Usage

```hcl

	data "ksyun_network_interfaces" "default" {
	  output_file="output_result"
	  ids=[]
	  securitygroup_id=[]
	  instance_type=[]
	  instance_id=[]
	  private_ip_address=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunNetworkInterfacesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Network Interface IDs, all the Network Interfaces belong to this region will be retrieved if the ID is `\"\"`.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of network interface resources that satisfy the condition.",
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of VPC IDs.",
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of subnet IDs.",
			},
			"securitygroup_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of security group IDs.",
			},
			"instance_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of instance types.",
			},
			"instance_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of VPC instance IDs.",
			},
			"private_ip_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of private IPs.",
			},
			"network_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of network interfaces. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the network interface.",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the network interface.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network interface.",
						},
						"network_interface_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network interface.",
						},
						"network_interface_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the network interface.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mac address of the network interface.",
						},
						"security_group_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of security groups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the security group.",
									},
									"security_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the security group.",
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the instance.",
						},
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "private IP.",
						},
						"d_n_s1": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 1.",
						},
						"d_n_s2": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 2.",
						},
						"assigned_private_ip_address_set": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_ip_address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Secondary Private IP.",
									},
								},
							},
							Description: "Assign secondary private ips to the network interface.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunNetworkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.ReadAndSetNetworkInterfaces(d, dataSourceKsyunNetworkInterfaces())
}
