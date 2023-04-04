/*
This data source provides a list of EIP resources (Elastic IP address) according to their EIP ID.

# Example Usage

```hcl

	data "ksyun_eips" "default" {
	  output_file="output_result"

	  ids=[]
	  project_id=[]
	  instance_type=[]
	  network_interface_id=[]
	  internet_gateway_id=[]
	  band_width_share_id=[]
	  line_id=[]
	  public_ip=[]
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunEipsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Elastic IP IDs, all the EIPs belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more project IDs.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Elastic IPs that satisfy the condition.",
			},
			"network_interface_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of NetworkInterface IDs.",
			},
			"internet_gateway_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of InternetGateway IDs.",
			},
			"instance_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Instance Type.",
			},
			"band_width_share_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of BandWidthShare IDs.",
			},
			"line_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Line IDs.",
			},
			"public_ip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of EIP address.",
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ipv4",
					"ipv6",
					"all",
				}, false),
				Description: "IP Version.",
			},
			"eips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of EIP. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP Version.",
						},
						"internet_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "InternetGateway ID.",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NetworkInterface ID.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the EIP.",
						},
						"allocation_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the EIP.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "project ID.",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID.",
						},
						"band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "bandwidth of the EIP.",
						},

						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance id to bind with the EIP.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state of the EIP.",
						},

						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance type to bind with the EIP.",
						},

						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EIP address.",
						},

						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the EIP.",
						},
						"band_width_share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the ID of the BWS which the EIP associated.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunEipsRead(d *schema.ResourceData, meta interface{}) error {
	eipService := EipService{meta.(*KsyunClient)}
	return eipService.ReadAndSetAddresses(d, dataSourceKsyunEips())
}
