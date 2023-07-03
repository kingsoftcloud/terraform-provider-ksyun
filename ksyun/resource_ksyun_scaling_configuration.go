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
		Create: resourceKsyunScalingConfigurationCreate,
		Read:   resourceKsyunScalingConfigurationRead,
		Delete: resourceKsyunScalingConfigurationDelete,
		Update: resourceKsyunScalingConfigurationUpdate,
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
				Description:  "The system disk type of the desired ScalingConfiguration.Valid Values:'Local_SSD', 'SSD3.0', 'EHDD'.",
			},

			"system_disk_size": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskSize,
				Description:  "The system disk size of the desired ScalingConfiguration.",
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
							Description:  "The EBS Data Disk Type of the desired data_disk.Valid Values: 'SSD3.0', 'EHDD'.",
						},
						"disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateKecDataDiskSize,
							Description:  "The EBS Data Disk Size of the desired data_disk.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "The Flag with delete EBS Data Disk when KEC Instance destroy.",
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
				Computed:    true,
				Description: "Charge type.",
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
		"instance_type":    {Type: TransformWithN},
		"key_id":           {Type: TransformWithN},
		"system_disk_type": {mapping: "SystemDisk.DiskType"},
		"system_disk_size": {mapping: "SystemDisk.DiskSize"},
		"data_disks": {mappings: map[string]string{
			"data_disks": "DataDisk",
			"disk_size":  "Size",
			"disk_type":  "Type",
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
	projectErr := addProjectInfo(d, &readScalingConfiguration, client)
	if projectErr != nil {
		return projectErr
	}
	action := "DescribeScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, readScalingConfiguration)
	resp, err := conn.DescribeScalingConfiguration(&readScalingConfiguration)
	if err != nil {
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
