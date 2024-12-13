/*
Provides an replica set MongoDB resource.

# Example Usage

```hcl

	resource "ksyun_mongodb_instance" "default" {
	  name = "InstanceName"
	  instance_account = "root"
	  instance_password = "admin"
	  instance_class = "1C2G"
	  storage = 5
	  node_num = 3
	  vpc_id = "VpcId"
	  vnet_id = "VnetId"
	  db_version = "3.6"
	  pay_type = "byDay"
	  iam_project_id = "0"
	  availability_zone = "cn-shanghai-3b"
	}

```

# Import

MongoDB can be imported using the id, e.g.

```
$ terraform import ksyun_mongodb_instance.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunMongodbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbInstanceCreate,
		Delete: resourceMongodbInstanceDelete,
		Update: resourceMongodbInstanceUpdate,
		Read:   resourceMongodbInstanceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of instance, which contains 6-64 characters and only support Chinese, English, numbers, '-', '_'.",
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// set stable value before api support query
				ValidateFunc: validation.StringInSlice([]string{
					"3.2",
					"3.6",
					"4.0",
				}, false),
				Default:     "3.2",
				Description: "The version of instance engine, and support `3.2`, `3.6`, `4.0`, default is `3.2`.",
			},
			"instance_class": {
				Type:     schema.TypeString,
				Optional: true,
				// set stable value before api support query
				ValidateFunc: validation.StringInSlice([]string{
					"1C2G",
					"2C4G",
					"4C8G",
					"8C16G",
					"8C32G",
					"16C64G",
					"16C128G",
				}, false),
				Default:     "1C2G",
				Description: "The class of instance cpu and memory.",
			},
			"storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
				Default:      5,
				Description:  "The size of instance disk, measured in GB (GigaByte).",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of VPC linked to the instance.",
			},
			"vnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of subnet linked to the instance.",
			},
			"instance_account": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "root",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"root",
				}, false),
				Description: "The administrator name of instance, if not defined `instance_account`, the instance will use `root`.",
			},
			"instance_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The administrator password of instance.",
			},
			"pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "hourlyInstantSettlement",
				ValidateFunc: validation.StringInSlice([]string{
					"byMonth",
					"byDay",
					"hourlyInstantSettlement",
				}, false),
				ForceNew:    true,
				Description: "Instance charge type, if not defined `pay_type`, the instance will use `byMonth`.",
			},
			"duration": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: durationSchemaDiffSuppressFunc("pay_type", "byMonth"),
				ForceNew:         true,
				Description:      "The duration of instance use, if `pay_type` is `byMonth`, the duration is required.",
			},
			"iam_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0",
				Description: "The project id of instance belong, if not defined `iam_project_id`, the instance will use `0`.",
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
				ValidateFunc: validation.IntInSlice([]int{
					3, 5, 7,
				}),
				Description: "The num of instance node.",
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
				Description:      "Availability zone where instance is located.",
			},
			"cidrs": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
				Description:      "network cidr.",
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VPC",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VPC",
				}, false),
				Description: "the type of network.",
			},

			"tags": tagsSchema(),

			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User ID.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of instance.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the status of instance.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of instance.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of security group.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "port number.",
			},
			"timing_switch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "timing switch for backup.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "timezone of backup.",
			},
			"time_cycle": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "time cycle of backup.",
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the product.",
			},
			"product_what": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "whether the instance is trial or not.",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of the MongoDB.",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration date of the MongoDB.",
			},
			"iam_project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the project.",
			},
			"mongos_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "number of mongos.",
			},
			"shard_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "number of shards.",
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MongoDB cluster mode.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance specification.",
			},
		},
	}
}

func resourceMongodbInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createMongodbInstanceCommon(d, meta, resourceKsyunMongodbInstance())
	if err != nil {
		return err
	}
	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunKrds(), "mongodb-instance", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
	}
	return resourceMongodbInstanceRead(d, meta)
}

func resourceMongodbInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		err = removeMongodbInstance(d, meta)
		if err == nil {
			return nil
		} else {
			_, err = readMongodbInstance(d, meta, "")
			if err != nil && canNotFoundMongodbError(err) {
				return nil
			}
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceMongodbInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyMongodbInstanceCommon(d, meta, resourceKsyunMongodbInstance())
	if err != nil {
		return err
	}
	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunKrds(), "mongodb-instance", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
	}
	return resourceMongodbInstanceRead(d, meta)
}

func resourceMongodbInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	return readMongodbInstanceCommon(d, meta, resourceKsyunMongodbInstance())
}
