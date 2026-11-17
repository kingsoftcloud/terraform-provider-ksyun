/*
Provides a ScalingConfiguration resource.

# Example Usage

```hcl

	resource "ksyun_scaling_configuration" "foo" {
	  scaling_configuration_name = "tf-xym-test-1"
	  image_id = "IMG-5465174a-6d71-4770-b8e1-917a0dd92466"
	  instance_type = "N3.1B"
	  password = "Aa123456"
	}

```

# Import

scalingConfiguration can be imported using the `id`, e.g.

```
$ terraform import ksyun_scaling_configuration.example scaling-configuration-abc123456
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create:        resourceKsyunScalingConfigurationCreate,
		Read:          resourceKsyunScalingConfigurationRead,
		Delete:        resourceKsyunScalingConfigurationDelete,
		Update:        resourceKsyunScalingConfigurationUpdate,
		CustomizeDiff: scalingConfigurationCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"scaling_configuration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tf-scaling-config",
				ForceNew:    false,
				Description: "The Name of the desired ScalingConfiguration.",
			},

			"image_id": {
				Type:        schema.TypeString,
				ForceNew:    false,
				Required:    true,
				Description: "The System Image Id of the desired ScalingConfiguration.",
			},

			"instance_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "I1.1A",
				ForceNew:         false,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
				Description:      "The KEC instance type of the desired ScalingConfiguration.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password.",
			},

			"system_disk_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskType,
				Description:  "The system disk type of the desired ScalingConfiguration.Valid Values:'Local_SSD', 'SSD3.0', 'EHDD', 'ESSD_SYSTEM_PL0', 'ESSD_SYSTEM_PL1', 'ESSD_SYSTEM_PL2'.",
			},

			"system_disk_size": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskSize,
				Description:  "The system disk size of the desired ScalingConfiguration.",
			},

			"system_disk_resize_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
				Description:  "The resize type of the desired system disk when scaling.Valid Values: 'offline', 'online'.",
			},

			"data_disk_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The Local Volume GB size of the desired ScalingConfiguration.",
			},

			"data_disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "A list of data disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateKecDataDiskType,
							Description:  "The EBS Data Disk Type of the desired data_disk.Valid Values: 'SSD3.0', 'EHDD', 'ESSD_PL0', 'ESSD_PL1', 'ESSD_PL2'.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The EBS Data Disk Size of the desired data_disk. Valid range by disk_type: 'SSD3.0'/'EHDD': [10, 65536]; 'ESSD_PL0'/'ESSD_PL1': [20, 65536]; 'ESSD_PL2': [461, 65536].",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "The Flag with delete EBS Data Disk when KEC Instance destroy.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The Snapshot ID of the desired data disk. If set, the data disk will be created from this snapshot. Note: ModifyScalingConfiguration does not support changing snapshot_id, so any change forces a new resource.",
						},
					},
				},
			},

			"key_id": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The SSH key set of the desired ScalingConfiguration.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The Project Id of the desired ScalingConfiguration belong to.",
			},

			"keep_image_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The Flag with image login set of the desired ScalingConfiguration.",
			},

			"tag": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "A list of instance tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of the tag.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the tag.",
						},
					},
				},
			},

			"sync_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to sync the instance tags to EBS data disks.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The KEC instance name of the desired ScalingConfiguration.",
			},

			"instance_name_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The kec instance name suffix of the desired ScalingConfiguration.",
			},

			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user data of the desired ScalingConfiguration.",
			},

			"instance_name_time_suffix": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The kec instance name time suffix of the desired ScalingConfiguration.",
			},

			"need_monitor_agent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKecInstanceAgent,
				Description:  "The Monitor agent flag desired ScalingConfiguration.",
			},

			"need_security_agent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKecInstanceAgent,
				Description:  "The Security agent flag desired ScalingConfiguration.",
			},

			"address_band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The EIP BandWidth.",
			},

			"band_width_share_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of BandWidthShare.",
			},

			"line_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Line ID Of EIP.",
			},

			"address_project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The Project ID of EIP.",
			},

			"charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The charge type of the instance, e.g. 'HourlyInstantSettlement' (pay-as-you-go).",
			},

			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CPU.",
			},

			"gpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "GPU.",
			},

			"mem": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
		},
	}
}

func resourceKsyunScalingConfigurationExtra(d *schema.ResourceData, forceGet bool) map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	var r map[string]SdkReqTransform

	r = map[string]SdkReqTransform{
		"instance_type":           {Type: TransformWithN},
		"key_id":                  {Type: TransformWithN},
		"system_disk_type":        {mapping: "SystemDisk.DiskType"},
		"system_disk_size":        {mapping: "SystemDisk.DiskSize"},
		"system_disk_resize_type": {mapping: "SystemDisk.ResizeType"},
		"data_disks": {mappings: map[string]string{
			"data_disks":           "DataDisk",
			"disk_size":            "Size",
			"disk_type":            "Type",
			"delete_with_instance": "DeleteWithInstance",
			"snapshot_id":          "SnapshotId",
		}, Type: TransformListN},
		"tag": {mappings: map[string]string{
			"tag":   "Tag",
			"key":   "Key",
			"value": "Value",
		}, Type: TransformListN},
	}
	extra = SdkRequestAutoExtra(r, d, forceGet)
	return extra
}

func resourceKsyunScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	scalingConfiguration := resourceKsyunScalingConfiguration()

	var resp *map[string]interface{}
	var err error

	createScalingConfiguration, err := SdkRequestAutoMapping(d, scalingConfiguration, false, nil,
		resourceKsyunScalingConfigurationExtra(d, false))
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfiguration, %s", err)
	}

	action := "CreateScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, createScalingConfiguration)
	resp, err = conn.CreateScalingConfiguration(&createScalingConfiguration)
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfiguration, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ScalingConfigurationId"].(string))
	}
	return resourceKsyunScalingConfigurationRead(d, meta)
}

func resourceKsyunScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	scalingConfiguration := resourceKsyunScalingConfiguration()

	var err error

	modifyScalingConfiguration, err := SdkRequestAutoMapping(d, scalingConfiguration, true, nil, resourceKsyunScalingConfigurationExtra(d, true))
	if err != nil {
		return fmt.Errorf("error on modifying ScalingConfiguration, %s", err)
	}

	if len(modifyScalingConfiguration) > 0 {
		modifyScalingConfiguration["ScalingConfigurationId"] = d.Id()
		action := "ModifyScalingConfiguration"
		logger.Debug(logger.ReqFormat, action, modifyScalingConfiguration)
		_, err = conn.ModifyScalingConfiguration(&modifyScalingConfiguration)
		if err != nil {
			return fmt.Errorf("error on modifying ScalingConfiguration, %s", err)
		}
	}
	return resourceKsyunScalingConfigurationRead(d, meta)
}

func resourceKsyunScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	readScalingConfiguration := make(map[string]interface{})
	readScalingConfiguration["ScalingConfigurationId.1"] = d.Id()
	// 不注入 ProjectId 过滤：按 ScalingConfigurationId 查询本身即可命中，
	// 注入全账号项目列表会触发服务端 ProjectId.N 参数配额上限（ParameterQuotaOutProjectId），
	// 且配置落在默认项目(0)时会被非 0 项目过滤查空，导致 Read 误报资源不存在。
	action := "DescribeScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, readScalingConfiguration)
	resp, err := conn.DescribeScalingConfiguration(&readScalingConfiguration)
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading ScalingConfiguration %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingConfigurationSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		delete(items[0].(map[string]interface{}), "InstanceType")
		SdkResponseAutoResourceData(d, resourceKsyunScalingConfiguration(), items[0], scalingConfigurationSpecialMapping())
	}
	return nil
}

func resourceKsyunScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	deleteScalingConfiguration := make(map[string]interface{})
	deleteScalingConfiguration["ScalingConfigurationId.1"] = d.Id()
	action := "DeleteScalingConfiguration"
	otherErrorRetry := 10

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteScalingConfiguration)
		resp, err1 := conn.DeleteScalingConfiguration(&deleteScalingConfiguration)
		logger.Debug(logger.AllFormat, action, deleteScalingConfiguration, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return OtherErrorProcess(&otherErrorRetry, fmt.Errorf("error on  deleting ScalingConfiguration %q, %s", d.Id(), err1))
		}
	})

}
