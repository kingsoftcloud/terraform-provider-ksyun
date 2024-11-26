/*
Provides a EBS resource.

# Example Usage

```hcl

		resource "ksyun_volume" "default" {
		  volume_name       = "test"
		  volume_type       = "SSD3.0"
		  size              = 15
		  charge_type       = "Daily"
		  availability_zone = "cn-shanghai-3a"
		  volume_desc       = "test"

		  ## 传入快照ID，用快照创建EBS盘
		  ## 注意：
	      ##   如果使用的整机镜像创建主机，API默认会自动根据镜像中包含的快照创建数据盘，不需在tf配置中定义数据盘
		  # snapshot_id = "snapshot_id"
		}

```

# Import

Instance can be imported using the `id`, e.g.

```
$ terraform import ksyun_volume.default xxxxxx
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeCreate,
		Update: resourceKsyunVolumeUpdate,
		Read:   resourceKsyunVolumeRead,
		Delete: resourceKsyunVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the EBS volume.",
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					// "SSD2.0",
					// "SSD3.0",
					// "EHDD",
					// "SATA2.0",
					"SSD3.0",
					"EHDD",
					"ESSD_PL0",
					"ESSD_PL1",
					"ESSD_PL2",
					"ESSD_PL3",
				}, false),
				Default:     "SSD3.0",
				Description: "The type of the EBS volume. Valid values:ESSD_PL0/ESSD_PL1/ESSD_PL2/ESSD_PL3/SSD3.0/EHDD, default is `SSD3.0`.",
			},
			"volume_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the EBS volume.",
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 32000),
				Default:      10,
				Description:  "The capacity of the EBS volume, in GB. Value range: [10, 32000], Default is 10.",
			},

			"online_resize": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: volumeDiffSuppressFunc,
				Description:      "Specifies whether to expand the capacity of the EBS volume online, default is true.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The availability zone in which the EBS volume resides.",
			},

			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HourlyInstantSettlement",
					"Daily",
				}, false),
				Description: "The billing mode of the EBS volume. Valid values: 'HourlyInstantSettlement', 'Daily'.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The ID of the project.",
			},
			"volume_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EBS volume.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the EBS volume was created.",
			},
			"volume_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The category to which the EBS volume belongs. Valid values: 'system' and 'data'.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the KEC instance to which the EBS volume is to be attached.",
			},

			// 快照建盘（API不返回这个值，所以diff时忽略这个值）
			"snapshot_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: kecDiskSnapshotIdDiffSuppress,
				Description:      "When the cloud disk snapshot opens, the snapshot id is entered.",
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceKsyunVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.CreateVolume(d, resourceKsyunVolume())
	if err != nil {
		return fmt.Errorf("error on creating volume %q, %s", d.Id(), err)
	}
	err = resourceKsyunVolumeRead(d, meta)
	return err
}

func resourceKsyunVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ReadAndSetVolume(d, resourceKsyunVolume())
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVolumeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ModifyVolume(d, resourceKsyunVolume())
	if err != nil {
		return fmt.Errorf("error on updating volume %q, %s", d.Id(), err)
	}
	err = resourceKsyunVolumeRead(d, meta)
	return err
}

func resourceKsyunVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.RemoveVolume(d)
	if err != nil {
		return fmt.Errorf("error on deleting volume %q, %s", d.Id(), err)
	}
	return err
}
