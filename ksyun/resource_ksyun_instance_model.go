/*
Provides a KEC instance model resource.

Instance Model is a launch template that allows you to quickly create instances with predefined configurations.

**Note**  Instance model cannot be modified after creation. If you need to change any parameter, you must delete and recreate the model.

Example Usage

```hcl
# get images list
data "ksyun_images" "centos-8_0" {
  platform = "centos-8.0"
}

data "ksyun_availability_zones" "default" {
}

# vpc settings of creating instance
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}

resource "ksyun_subnet" "default" {
  subnet_name       = "ksyun-subnet-tf"
  cidr_block        = "10.7.0.0/21"
  subnet_type       = "Normal"
  vpc_id            = ksyun_vpc.default.id
  availability_zone = data.ksyun_availability_zones.default.availability_zones[0].availability_zone_name
}

resource "ksyun_security_group" "default" {
  vpc_id              = ksyun_vpc.default.id
  security_group_name = "ksyun-security-group"
}

resource "ksyun_instance_model" "default" {
  model_name          = "web-server-model"
  image_id            = data.ksyun_images.centos-8_0.images[0].image_id
  instance_type       = "I2.8B"
  charge_type         = "Monthly"
  purchase_time       = 12
  security_group_ids  = [ksyun_security_group.default.id]
  subnet_id           = ksyun_subnet.default.id
  key_id              = "key-12345678"
  instance_name       = "db-server"
  instance_name_suffix = "1"
  project_id          = 1001
  allocate_address    = true
  address_bandwidth   = 5
  address_charge_type = "Peak"

  data_disks {
    type = "SSD3.0"
    size = 500
  }

  tags {
    key   = "env"
    value = "production"
  }
}
```

Import

Instance Model can be imported using the model id, e.g.

```
$ terraform import ksyun_instance_model.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func instanceModelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"model_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the instance model.",
		},
		"model_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the instance model. Must be globally unique.",
		},
		"image_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID for the image to use for the instance.",
		},
		"charge_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Monthly",
				"Daily",
				"HourlyInstantSettlement",
				"Spot",
			}, false),
			Description: "Charge type of the instance. Valid values: Monthly, Daily, HourlyInstantSettlement, Spot.",
		},
		"security_group_id": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Security Group IDs to associate with. Currently only supports 1 security group.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The type of instance to start.",
		},
		"purchase_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 36),
			Description:  "The duration that you will buy the resource. Required when charge_type is Monthly, value range: [1, 36].",
		},
		"data_disk_gb": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 16000),
			Description:  "The size of the local SSD disk. Not effective for general purpose instances. Value range: [0, 16000].",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of subnet. The instance will use the subnet in the current region.",
		},
		"keep_image_login": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Keep the initial settings of the custom image. Mutually exclusive with password/key_id.",
		},
		"key_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The certificate id of the instance. Mutually exclusive with password. Not supported for other-linux.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.",
		},
		"instance_name_suffix": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The suffix of the instance name. Range: 0-9999, effective for batch creation.",
		},
		"sriov_net_support": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Description: "Whether to support networking enhancement. Valid for I1/C1/I2(8C+) with standard images.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "The project instance belongs to. 0 is the default project.",
		},
		"data_guard_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Add instance being created to a disaster tolerance group.",
		},
		"allocate_address": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Whether to allocate EIP to the instance. EIP parameters take effect when true.",
		},
		"address_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "The bandwidth of the EIP.",
		},
		"line_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The line id of the EIP.",
		},
		"address_charge_type": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Monthly",
				"Peak",
				"Daily",
				"TrafficMonthly",
				"DailyPaidByTransfer",
				"HourlyInstantSettlement",
				"RegionPeak",
				"HourlySettlement",
				"PrepaidByTime",
				"PostpaidByTime",
				"PostPaidByAdvanced95Peak",
				"DailyPaidByTransfer",
			}, false),
			Description: "The charge type of the EIP.",
		},
		"address_purchase_time": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 36),
			Description:  "The purchase time of the EIP. Required when address_charge_type is Monthly.",
		},
		"address_project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "The project ID of the EIP.",
		},
		"failure_auto_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Whether to automatically delete the instance when creation fails.",
		},
		"host_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The hostname of the instance. OS internal computer name.",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The user data to be specified into this instance. Must be encrypted in base64 format.",
		},
		"is_distribute_ipv6": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Whether to distribute IPv6 address.",
		},
		"mem": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The memory size of the instance.",
		},
		"cpu": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The CPU count of the instance.",
		},
		"iam_role_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The name of IAM role.",
		},
		"assembled_image_data_disk_type": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Data disk type for assembled image.",
		},
		"local_volume_snapshot_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Local volume snapshot ID.",
		},
		"sync_tag": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Whether to sync EBS tags.",
		},
		"system_disk": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Description: "System disk parameters.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"SSD3.0",
							"SSD2.0",
							"EHDD",
							"SATA2.0",
							"Local_SSD",
							"ESSD_PL0",
							"ESSD_PL1",
							"ESSD_PL2",
							"ESSD_Entry",
							"ESSD_AutoPL",
							"ESSD_SYSTEM_PL0",
							"ESSD_SYSTEM_PL1",
							"ESSD_SYSTEM_PL2",
							"ESSD_SYSTEM_Entry",
							"ESSD_SYSTEM_AutoPL",
						}, false),
						Description: "System disk type.",
					},
					"disk_size": {
						Type:         schema.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(20, 500),
						Description:  "System disk size. Value range: [20, 500].",
					},
				},
			},
		},
		"data_disks": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    8,
			Description: "The list of data disks created with instance.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"SSD3.0",
							"EHDD",
							"Local_SSD",
							"ESSD_PL1",
							"ESSD_PL2",
							"ESSD_PL0",
							"ESSD_Entry",
							"ESSD_AutoPL",
						}, false),
						Description: "Data disk type.",
					},
					"disk_size": {
						Type:         schema.TypeInt,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(1, 65536),
						Description:  "Data disk size. Value range: [1, 65536].",
					},
					"delete_with_instance": {
						Type:        schema.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
						Description: "Whether to delete the data disk when the instance is deleted.",
					},
					"disk_snapshot_id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Snapshot ID for creating data disk.",
					},
					"snapshot_name": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Snapshot name for creating data disk.",
					},
				},
			},
		},
		"network_interface": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MinItems:    1,
			Description: "Network interface configurations.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"subnet_id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Subnet ID for the network interface.",
					},
					"security_group_id": {
						Type:        schema.TypeList,
						Optional:    true,
						ForceNew:    true,
						Description: "Security group IDs for the network interface.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"private_ip_address": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Private IP address for the network interface.",
					},
				},
			},
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "Tags to bind to the instance.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Tag key.",
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Tag value.",
					},
					"id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Tag ID.",
					},
				},
			},
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the instance model.",
		},
	}
}

func resourceKsyunInstanceModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstanceModelCreate,
		Read:   resourceKsyunInstanceModelRead,
		Delete: resourceKsyunInstanceModelDelete,
		// No Update - model template cannot be modified
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: instanceModelSchema(),
	}
}

func resourceKsyunInstanceModelCreate(d *schema.ResourceData, meta interface{}) (err error) {
	instanceModelService := InstanceModelService{meta.(*KsyunClient)}
	err = instanceModelService.CreateModel(d, resourceKsyunInstanceModel())
	if err != nil {
		return fmt.Errorf("error on creating InstanceModel %q, %s", d.Id(), err)
	}
	return resourceKsyunInstanceModelRead(d, meta)
}

func resourceKsyunInstanceModelRead(d *schema.ResourceData, meta interface{}) (err error) {
	instanceModelService := InstanceModelService{meta.(*KsyunClient)}
	err = instanceModelService.ReadAndSetModel(d, resourceKsyunInstanceModel())
	if err != nil {
		return fmt.Errorf("error on reading InstanceModel %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunInstanceModelDelete(d *schema.ResourceData, meta interface{}) (err error) {
	instanceModelService := InstanceModelService{meta.(*KsyunClient)}
	err = instanceModelService.RemoveModel(d)
	if err != nil {
		return fmt.Errorf("error on deleting InstanceModel %q, %s", d.Id(), err)
	}
	return err
}
