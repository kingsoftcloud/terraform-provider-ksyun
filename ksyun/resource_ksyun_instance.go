/*
Provides a KEC instance resource.

**Note**  At present, 'Monthly' instance cannot be deleted and must wait it to be outdated and released automatically.

Example Usage

```hcl

## get images list
data "ksyun_images" "centos-8_0" {
  platform = "centos-8.0"
}

data "ksyun_availability_zones" "default" {
}

## vpc settings of creating instance
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

resource "ksyun_instance" "foo" {
  image_id      = data.ksyun_images.centos-8_0.images[0].image_id
  instance_type = "N3.2B"

  subnet_id         = ksyun_subnet.default.id
  instance_password = "Xuan663222"
  keep_image_login  = false
  charge_type       = "Daily"
  purchase_time     = 1
  security_group_id = [ksyun_security_group.default.id]
  instance_name     = "ksyun-kec-tf-demotion"
  sriov_net_support = "false"
  data_disks {
    disk_type            = "SSD3.0"
    disk_size            = 40
    delete_with_instance = true
  }
  key_id          = []
  auto_create_ebs = true
}
```

Import

Instance can be imported using the id, e.g.

```
$ terraform import ksyun_instance.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func instanceConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID for the image to use for the instance.",
		},
		"instance_status": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				"active",
				"stopped",
			}, false),
			Description: "The state of instance.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The type of instance to start. <br> - NOTE: it's may trigger this instance to power off, if instance type will be demotion.",
		},
		"system_disk": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "System disk parameters.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"SSD3.0",
							"EHDD",
							"Local_SSD",
							"ESSD_SYSTEM_PL0",
							"ESSD_SYSTEM_PL1",
							"ESSD_SYSTEM_PL2",
						}, false),
						Description: "System disk type. `Local_SSD`, Local SSD disk. `SSD3.0`, The SSD cloud disk. `EHDD`, The EHDD cloud disk, `ESSD_SYSTEM_PL0`, The x7 machine type ESSD disk, `ESSD_SYSTEM_PL1`, The x7 machine type ESSD disk, `ESSD_SYSTEM_PL2`, The x7 machine type ESSD disk.",
					},
					"disk_size": {
						Type:         schema.TypeInt,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IntBetween(20, 500),
						Description:  "The size of the data disk. value range: [20, 500].",
					},
				},
			},
		},
		"data_disk_gb": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 16000),
			Description:  "The size of the local SSD disk.",
		},
		"data_disks": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    8,
			Computed:    true,
			Description: "The list of data disks created with instance.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"SSD3.0",
							"EHDD",
							"Local_SSD",
							"ESSD_PL0",
							"ESSD_PL1",
							"ESSD_PL2",
							"ESSD_PL3",
						}, false),
						Description: "Data disk type.",
					},
					"disk_size": {
						Type:         schema.TypeInt,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(10, 16000),
						Description:  "Data disk size. value range: [10, 16000].",
					},
					// 快照建盘（API不返回这个值，所以diff时忽略这个值）
					"disk_snapshot_id": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						DiffSuppressFunc: kecDiskSnapshotIdDiffSuppress,
						Description:      "When the cloud disk opens, the snapshot id is entered.",
					},
					"delete_with_instance": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						ForceNew:    true,
						Description: "Delete this data disk when the instance is destroyed. It only works on EBS disk.",
					},
					"disk_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the disk.",
					},
				},
			},
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of subnet. the instance will use the subnet in the current region.",
		},
		"extension_network_interface": {
			Type:        schema.TypeSet,
			Computed:    true,
			Description: "extension network interface information.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_interface_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the extension network interface.",
					},
				},
			},
			Set: kecNetworkInterfaceHash,
		},

		"local_volume_snapshot_id": {
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			DiffSuppressFunc: kecImportDiffSuppress,
			Description:      "When the local data disk opens, the snapshot id is entered.",
		},

		"instance_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Sensitive:   true,
			Description: "Password to an instance is a string of 8 to 32 characters.",
		},
		"keep_image_login": {
			Type:             schema.TypeBool,
			Optional:         true,
			DiffSuppressFunc: kecImportDiffSuppress,
			Description:      "Keep the initial settings of the custom image.",
		},

		"key_id": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Set:         schema.HashString,
			Description: "The certificate id of the instance.",
		},

		"charge_type": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Daily",
				"HourlyInstantSettlement",
			}, false),
			Description: "charge type of the instance.",
		},
		"purchase_time": {
			Type:             schema.TypeInt,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
			ValidateFunc:     validation.IntBetween(0, 36),
			Description:      "The duration that you will buy the resource.",
		},
		"security_group_id": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Required:    true,
			Set:         schema.HashString,
			MinItems:    1,
			Description: "Security Group to associate with.",
		},
		"private_ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Instance private IP address can be specified when you creating new instance.",
		},
		// eip和主机的绑定关系，放在绑定的resource里描述，不在vm的结构里提供这个字段
		// 否则后绑定，资源创建完成时这个字段为空
		// "public_ip": {
		//	Type:     schema.TypeString,
		//	Computed: true,
		// },
		"instance_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The name of instance, which contains 2-64 characters and only support Chinese, English, numbers.",
		},
		"sriov_net_support": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"true",
				"false",
			}, false),
			Description: "whether support networking enhancement.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The project instance belongs to.",
		},
		"data_guard_id": {
			Type:     schema.TypeString,
			Optional: true,
			// ForceNew:    true,
			Description: "Add instance being created to a disaster tolerance group. It will be quit the disaster tolerance group, if this field change to null.",
		},
		"host_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The hostname of the instance. only effective when image support cloud-init.",
		},
		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
			// ForceNew:         true,
			DiffSuppressFunc: kecImportDiffSuppress,
			Description:      "The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB. only effective when image support cloud-init.",
		},
		"iam_role_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "name of iam role.",
		},
		"force_reinstall_system": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicate whether to reinstall system.",
		},
		"dns1": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "DNS1 of the primary network interface.",
		},
		"dns2": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "DNS2 of the primary network interface.",
		},
		"tags": tagsSchema(),
		// "has_init_info": {
		//	Type:     schema.TypeBool,
		//	Computed: true,
		// },
		// some control
		"has_modify_system_disk": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "whether the system disk has modified.",
		},
		"has_modify_password": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "whether the password has modified.",
		},
		"has_modify_keys": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "whether the certificate key has modified.",
		},

		"force_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Deprecated:  "this field is Deprecated and no effect for change",
			Description: "Indicate whether to delete instance directly or not.",
		},

		"sync_tag": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Indicate whether to sync tags to instance.",
		},

		"network_interface_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the network interface.",
		},

		"instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the instance.",
		},

		// AutoCreateEbs: 是否自动创建数据盘；
		// 针对整机镜像，如果是true会自动使用镜像关联的快照建盘，这种情况用tf管理盘会比较困难，因此tf默认设置为false
		"auto_create_ebs": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
			// ForceNew:         true,
			DiffSuppressFunc: kecImportDiffSuppress,
			Description:      "Whether to create EBS volumes from snapshots in the custom image, default is false.",
		},
	}
}

func resourceKsyunInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstanceCreate,
		Update: resourceKsyunInstanceUpdate,
		Read:   resourceKsyunInstanceRead,
		Delete: resourceKsyunInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: instanceConfig(),
	}
}

func resourceKsyunInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.createKecInstance(d, resourceKsyunInstance())
	if err != nil {
		return fmt.Errorf("error on creating Instance: %s", err)
	}
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.readAndSetKecInstance(d, resourceKsyunInstance(), false)
	if err != nil {
		return fmt.Errorf("error on reading Instance: %s", err)
	}
	return err
}

func resourceKsyunInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	if d.HasChange("sync_tag") {
		return fmt.Errorf("changing sync_tag is not supported")
	}
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.modifyKecInstance(d, resourceKsyunInstance())
	if err != nil {
		return fmt.Errorf("error on updating Instance: %s", err)
	}
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.removeKecInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting Instance: %s", err)
	}
	return err
}
