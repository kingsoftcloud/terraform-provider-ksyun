/*
This data source provides a list of Bare Metal resources according to their Bare Metal ID.

# Example Usage

```hcl
# Get  bare metals

	data "ksyun_bare_metals" "default" {
	  output_file="output_result"
	  ids = []
	  vpc_id = ["bfec0f43-9e5a-4f06-b7a1-df4768c1cd6f"]
	  project_id = []
	  host_name = []
	  subnet_id = []
	  cabinet_id = []
	  epc_host_status = []
	  os_name = []
	  product_type = []
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunBareMetals() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBareMetalsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of Bare Metal IDs, all the Bare Metals belong to this region will be retrieved if the ID is `\"\"`.",
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
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by Bare Metal name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Bare Metals that satisfy the condition.",
			},
			"host_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more Bare Metal host names.",
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more vpc IDs.",
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more subnet IDs.",
			},
			"cabinet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more Bare Metal cabinet IDs.",
			},
			"host_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more Bare Metal host types.",
			},
			"epc_host_status": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more Bare Metal status.",
			},
			"os_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "One or more Bare Metal operating system names.",
			},
			"product_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"lease",
						"customer",
						"lending",
					}, false),
				},
				Set:         schema.HashString,
				Description: "One or more Bare Metal product types. valid values: 'lease', 'customer', 'lending'.",
			},
			"bare_metals": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for Bare Metal.",
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Bare Metal.",
						},
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Bare Metal.",
						},
						"sn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SN of the Bare Metal.",
						},
						"cabinet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Cabinet.",
						},
						"raid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Raid type of the Bare Metal.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Image.",
						},
						"product_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "product type of the Bare metal.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the OS.",
						},
						"memory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the memory of the Bare Metal.",
						},
						"cpu": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "cpu specification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "model of the cpu.",
									},
									"frequence": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "frequence of the cpu.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of CPUs.",
									},
									"core_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of CPU cores.",
									},
								},
							},
						},
						"disk_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "a list of disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the disk.",
									},
									"raid": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "raid type of the disk.",
									},
									"space": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "space of the disk.",
									},
								},
							},
						},
						"network_interface_attribute_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "a list of network interfaces.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the VPC.",
									},
									"network_interface_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the Id of the network interface.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the ID of the subnet.",
									},
									"private_ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The private IP address assigned to the network interface.",
									},
									"dns1": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS1 of the network instance.",
									},
									"dns2": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS2 of the network instance.",
									},
									"mac": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "MAC of the network instance.",
									},
									"security_group_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "a list of security groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"security_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the security group.",
												},
											},
										},
									},
									"network_interface_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the network interface.",
									},
								},
							},
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability zone name.",
						},
						"cloud_monitor_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud monitor agent of the Bare Metal.",
						},
						"host_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status of the Bare Metal.",
						},
						"host_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of the Bare Metal.",
						},
						"network_interface_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "mode of the network interface.",
						},
						"security_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security agent of the Bare Metal.",
						},
						"enable_bond": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable bond.",
						},
						"data_disk_catalogue_suffix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data disk catalogue suffix.",
						},
						"data_file_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data file type.",
						},
						"system_volume_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System Volume Size.",
						},
						"nvme_data_disk_catalogue": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Nvme data disk catalogue.",
						},
						"enable_container": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether enable container.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"hyper_threading": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hyper Threading.",
						},
						"cabinet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cabinet Name.",
						},
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key ID.",
						},
						"allow_modify_hyper_threading": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Allow Modify Hyper Threading.",
						},
						"releasable_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Releasable Time.",
						},
						"rack_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rack Name.",
						},
						"kmr_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "KMR Agent.",
						},
						"nvme_data_disk_catalogue_suffix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Nvme data disk catalogue suffix.",
						},
						"support_ebs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support EBS.",
						},
						"kpl_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "KPL Agent.",
						},
						"service_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Servcie end time.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge type.",
						},
						"gpu": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
							Description: "Gpu specification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "model of the cpu.",
									},
									"frequence": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "frequence of the cpu.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of CPUs.",
									},
									"core_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of CPU cores.",
									},
									"gpu_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of GPU cores.",
									},
								},
							},
						},
						"roces": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Roces.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ip of the roce network.",
									},
									"mask": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mask of the roce network.",
									},
									"gate_way": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of roce.",
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

func dataSourceKsyunBareMetalsRead(d *schema.ResourceData, meta interface{}) error {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	return bareMetalService.ReadAndSetBareMetals(d, dataSourceKsyunBareMetals())
}
