/*
This data source providers a list of instance resources according to their availability zone, instance ID.

# Example Usage

```hcl
# Get  instances

	data "ksyun_instances" "default" {
	  output_file = "output_result"
	  ids = []
	  search = ""
	  project_id = []
	  network_interface {
	  	network_interface_id = []
	  	subnet_id = []
	  	group_id = []
	  }
	  instance_state {
	  	name =  []
	  }
	  availability_zone {
	  	name =  []
	  }
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of instance IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by instance name.",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex string to filter results by instance name or privateIpAddress.",
			},

			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Description: "Total number of instance that satisfy the condition.",
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID of subnet linked to the instance.",
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID of VPC linked to the instance.",
			},
			"network_interface": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "a list of network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The ID of VPC linked to the network interface.",
						},
						"network_interface_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "the ID of the network interface.",
						},
						"group_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The ID of security group linked to the network interface.",
						},
					},
				},
			},
			"instance_state": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The state of instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "name of the state.",
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Description: "the availability zone that the instance locates at.",
			},

			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "a list of instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the ID of the instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the instance.",
						},
						"instance_configure": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Description: "the configure of the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"v_c_p_u": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the number of the vcpu.",
									},
									"g_p_u": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the number of the gpu.",
									},
									"memory_gb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "memory capacity.",
									},
									"data_disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the data disk.",
									},
									"data_disk_gb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "size of the data disk.",
									},
									"root_disk_gb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "size of the root disk.",
									},
								},
							},
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the image.",
						},

						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of the instance.",
						},
						"instance_state": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Description: "state of the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the state.",
									},
								},
							},
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of subnet linked to the instance.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of VPC linked to the instance.",
						},
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance private IP address.",
						},
						"monitoring": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Description: "state of the monitoring.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the state.",
									},
								},
							},
						},

						"sriov_net_support": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "whether support networking enhancement.",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of creation for instance.",
						},
						"network_interface_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "a list of network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_interface_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the network interface.",
									},
									"network_interface_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the network interface.",
									},
									"mac_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "MAC address.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the subnet.",
									},
									"private_ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private ip address of the network interface.",
									},
									"public_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "public ip address of the network interface.",
									},
									"security_group_set": {
										Type:        schema.TypeList,
										Computed:    true,
										MaxItems:    1,
										Description: "a list of the security group.",
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
									"group_set": {
										Type:        schema.TypeList,
										Computed:    true,
										MaxItems:    1,
										Description: "a list of the security group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the security group.",
												},
											},
										},
									},
									"d_n_s1": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The dns1 of the network interface.",
									},
									"d_n_s2": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The dns2 of the network interface.",
									},
								},
							},
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project instance belongs to.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance charge type.",
						},

						"system_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Description: "System disk information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the system disk.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "size of the system disk.",
									},
								},
							},
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "count of the instance.",
						},
						"stopped_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "stopped mode.",
						},
						"availability_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone name.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone name.",
						},
						"product_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "product type of the instance.",
						},
						"product_what": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "whether the instance is trial or not.",
						},
						"auto_scaling_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of auto scaling.",
						},
						"is_show_sriov_net_support": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether support networking enhancement.",
						},
						"key_id": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The certificate id of the instance.",
						},
						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "a list of the data disks.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the data disk.",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "type of the data disk.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "size of the data disk.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Decides whether the disk is deleted with instance.",
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

func dataSourceKsyunInstancesRead(d *schema.ResourceData, meta interface{}) error {
	kecService := KecService{meta.(*KsyunClient)}
	return kecService.ReadAndSetKecInstances(d, dataSourceKsyunInstances())
}
