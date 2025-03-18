/*
Provides a Bare Metal resource.

# Example Usage

```hcl

	resource "ksyun_bare_metal" "default" {
	  host_name = "test"
	  host_type = "MI-I2"
	  image_id = "eb8c0428-476e-49af-8ccb-9fad2455a54c"
	  key_id = "9c45b560-e51d-4aee-9e99-0e292476692d"
	  network_interface_mode = "single"
	  raid = "Raid1"
	  availability_zone = "cn-beijing-6b"
	  security_agent = "classic"
	  cloud_monitor_agent = "classic"
	  subnet_id = "d2fdc1b5-0280-4ca7-920b-0bd0453c130c"
	  security_group_ids = ["7e2f45b5-e79d-4612-a7fc-fe74a50b639a"]
	  system_file_type = "EXT4"
	  container_agent = "supported"
	  force_re_install = false
	}

```

# Import

Bare Metal can be imported using the id, e.g.

```
$ terraform import ksyun_bera_metal.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunBareMetal() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBareMetalCreate,
		Read:   resourceKsyunBareMetalRead,
		Update: resourceKsyunBareMetalUpdate,
		Delete: resourceKsyunBareMetalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
		CustomizeDiff: bareMetalCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Availability Zone.",
			},
			"host_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Bare Metal Host Type (e.g. CAL-III).",
			},
			"hyper_threading": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Open",
					"Close",
					"NoChange",
				}, false),
				Default:     "NoChange",
				Description: "The HyperThread status of the Bare Metal. Valid Values:'Open','Close','NoChange'.Default is 'NoChange'.",
			},
			"roce_network": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Open",
					"Close",
				}, false),

				DiffSuppressFunc: bareMetalRoceNetwork,
				Description:      "The value of roce network that indicates acquiring whether an instance supplied roce network. Valid Options: `Open` and `Close`.",
			},
			"host_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Running",
					"Stopped",
				}, false),
				Description: "The status of Bare Metal instance. That can set your Bare Metal instance status, `Running` or `Stopped`, on ksyun. In detail, the instance will start, when `host_status` is `Running` but its status is `Stopped` on ksyun. Similarly, the instance will be power off, when `host_status` is `Stopped` but its status is `Running` on ksyun. <br> Value Options: `Running`, `Stopped`.",
			},
			"raid": {
				Type:     schema.TypeString,
				Optional: true,
				// ValidateFunc: validation.StringInSlice([]string{
				// 	"Raid0",
				// 	"Raid1",
				// 	"Raid5",
				// 	"Raid10",
				// 	"Raid50",
				// 	"SRaid0",
				// }, false),
				ConflictsWith: []string{"raid_id"},
				Description:   "The Raid type of the Bare Metal. Valid Values:'Raid0','Raid1','Raid5','Raid10','Raid50','SRaid0', 'Jbod'. Conflict raid_id. If you don't set raid_id,raid is Required.",
			},
			"raid_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"raid"},
				Description:   "The Raid template id of Bare Metal.Conflict raid. If you don't set raid,raid_id is Required. If you want to use raid_id,you must in user white list.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the image.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project id of the Bare Metal.Default is '0'.",
			},
			"network_interface_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bond4",
					"single",
					"dual",
				}, false),
				Default:     "bond4",
				Description: "The network interface mode of the Bare Metal. Valid Values:'bond4','single','dual'.Default is 'bond4'.When bond4->single,single->bond4,dual->single,dual->bond4 can modify,otherwise is ForceNew.",
			},
			"bond_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bond0",
					"bond1",
				}, false),
				Default:          "bond1",
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The bond attribute of the Bare Metal. Valid Values:'bond0','bond1'.Default is 'bond1'. Only effective when network_interface_mode is bond4.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet id of the Bare Metal primary network interface.",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private ip address of the Bare Metal primary network interface.",
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 3,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The security_group_id set of the Bare Metal primary network interface.",
			},
			"dns1": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The dns1 of the Bare Metal primary network interface.",
			},
			"dns2": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The dns2 of the Bare Metal primary network interface.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The certificate id of the Bare Metal.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ksc_epc",
				Description: "The name of the Bare Metal.Default is 'ksc_epc'.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of the Bare Metal.",
			},
			"security_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"classic",
					"no",
				}, false),
				Default:     "no",
				Description: "The security agent choice of the Bare Metal. Valid Values:'classic','no'. Default is 'no'.",
			},
			"cloud_monitor_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"classic",
					"no",
				}, false),
				Default:     "no",
				Description: "The cloud monitor agent choice of the Bare Metal.Valid Values:'classic','no'.Default is 'no'.",
			},
			"extension_subnet_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The subnet id of the Bare Metal primary extension interface.Only effective when network_interface_mode is dual and Required.",
			},
			"extension_private_ip_address": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The private ip address of the Bare Metal extension network interface.Only effective when network_interface_mode is dual.",
			},
			"extension_security_group_ids": {
				Type:     schema.TypeSet,
				MaxItems: 3,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:              schema.HashString,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The security_group_id set of the Bare Metal extension network interface.Max is 3.Only effective when network_interface_mode is dual and Required.",
			},
			"extension_dns1": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The dns1 of the Bare Metal extension network interface.Only effective when network_interface_mode is dual.",
			},
			"extension_dns2": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The dns2 of the Bare Metal extension network interface.Only effective when network_interface_mode is dual.",
			},
			"system_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
				Default:     "EXT4",
				Description: "The system disk file type of the Bare Metal.Valid Values:'EXT4','XFS'.Default is 'EXT4'.",
			},
			"data_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
				Default:     "XFS",
				Description: "The data disk file type of the Bare Metal.Valid Values:'EXT4','XFS'.Default is 'XFS'.",
			},
			"data_disk_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"/DATA/disk",
					"/data",
				}, false),
				Default:     "/DATA/disk",
				Description: "The data disk catalogue of the Bare Metal.Valid Values:'/DATA/disk','/data'.Default is '/DATA/disk'.",
			},
			"data_disk_catalogue_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NoSuffix",
					"NaturalNumber",
					"NaturalNumberFromZero",
				}, false),
				Default:     "NaturalNumber",
				Description: "The data disk catalogue suffix of the Bare Metal.Valid Values:'NoSuffix','NaturalNumber','NaturalNumberFromZero'.Default is 'NaturalNumber'.",
			},
			"nvme_data_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
				Description: "The nvme data file type of the Bare Metal.Valid Values:'EXT4','XFS'.",
			},
			"nvme_data_disk_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"/DATA/disk",
					"/data",
				}, false),
				Description: "The nvme data disk catalogue of the Bare Metal.Valid Values:'/DATA/disk','/data'.",
			},
			"nvme_data_disk_catalogue_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NoSuffix",
					"NaturalNumber",
					"NaturalNumberFromZero",
				}, false),
				Description: "The nvme data disk catalogue suffix of the Bare Metal.Valid Values:'NoSuffix','NaturalNumber','NaturalNumberFromZero'.",
			},
			"container_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"supported",
					"unsupported",
				}, false),
				Default:          "unsupported",
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "Whether to support KCE cluster, valid values: 'supported', 'unsupported'.",
			},
			"computer_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "The computer name of the Bare Metal.",
			},
			"server_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The pxe server ip of the Bare Metal.Only effective on modify and host type is COLO.",
			},
			"path": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "The path of the Bare Metal.Only effective on modify and host type is COLO.",
			},
			"force_re_install": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
				Description:      "Indicate whether to reinstall system.",
			},

			// 2025/03/10 20:02 added fields
			// "sn": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "Cloud physical host serial number.",
			// },
			//
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,

				ValidateFunc: validation.StringInSlice([]string{"Daily"}, false),
				Description:  "Charge Type. Valid Value: `Daily`.",
			},

			"trial": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: bareMetalCreateDiff,
				Description:      "Trial this epc instance.",
			},

			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: purchaseTimeTrialDiffSuppressFunc,
				ForceNew:         true,
				Description:      "Purchase time. If trial is true, it works. Default: 7.",
			},
			"address_band_width": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The band width of elastic ip.",
			},
			"line_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Line id.",
			},
			"band_width_share_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of share band width.",
			},
			"address_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The charge type of elastic ip.",
			},
			"address_purchase_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The purchase time.",
			},

			"address_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project id of elastic ip.",
			},

			"kes_agent": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "unsupport",
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "The KES Agent.",
			},
			"kmr_agent": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "unsupport",
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "The KMR Agent.",
			},
			"overclocking_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The overclocking attribute.",
			},
			"gpu_image_driver_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "The GPU version.",
			},
			"system_volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "System disk type of cloud disk.",
			},
			"system_volume_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "System disk size of cloud disk.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone id, when creating pdns, is required.",
			},
			"zone_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone type, when creating pdns, is required.",
			},
			"use_hot_standby": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether use hot standy. Valid Values: `support`, `unsupport` and `onlyHotStandby`.",
			},
			"timed_regularization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Trial timed conversion to regular status, when charge_type is `Trial`. Valid Values: `support`, `unsupported`.",
			},
			"tags": tagsSchema(),

			"hot_standby": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hot_stand_by_host_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of hot standby.",
						},
						"retain_instance_info": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether retain the instance info. Valid Values: `RetainPrivateIP` `Notretain`.",
						},
					},
				},
				Description: "Indicate the hot standby to instead the master Host.",
			},

			"activate_hot_standby": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: activateHotStandbyDSF,
				Description:      "Activate hot standby epc. it works, when this instance is standby.",
			},

			// 2025, Mar 18
			"password_inherit": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "Inherit password or keys from image, Valid Values: `support` `unsupport`.",
			},
			"data_disk_mount": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalReinstallDiffSuppressFunc,
				Description:      "Whether mount for data disk. Valid Values: `support` `unsupport`.",
			},
			"storage_roce_network_card_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of storage roce network card. Valid Values `eth8x_bond` `storage_bond`.",
			},

			"extension_network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the extension network interface.",
			},
			"network_interface_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the primary network interface.",
			},
		},
	}
}

func resourceKsyunBareMetalCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.CreateBareMetal(d, resourceKsyunBareMetal())
	if err != nil {
		return fmt.Errorf("error on creating bare metal %q, %s", d.Id(), err)
	}
	return resourceKsyunBareMetalRead(d, meta)
}

func resourceKsyunBareMetalRead(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.ReadAndSetBareMetal(d, resourceKsyunBareMetal())
	if err != nil {
		return fmt.Errorf("error on reading bare metal %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBareMetalUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	if d.HasChange("roce_network") {
		return fmt.Errorf("argument `roce_network` cannot be modified for now")
	}

	notAllowModifyFields := []string{"system_volume_size", "system_volume_type"}
	if d.HasChanges(notAllowModifyFields...) {
		return fmt.Errorf("argument %v cannot be modified for now", notAllowModifyFields)
	}
	err = bareMetalService.ModifyBareMetal(d, resourceKsyunBareMetal())
	if err != nil {
		return fmt.Errorf("error on updating bare metal %q, %s", d.Id(), err)
	}
	return resourceKsyunBareMetalRead(d, meta)
}

func resourceKsyunBareMetalDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.RemoveBareMetal(d)
	if err != nil {
		return fmt.Errorf("error on deleting bare metal %q, %s", d.Id(), err)
	}
	return err
}
