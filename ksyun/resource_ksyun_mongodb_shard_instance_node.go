package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceKsyunMongodbShardInstanceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbShardInstanceNodeCreate,
		Delete: resourceMongodbShardInstanceNodeDelete,
		Update: resourceMongodbShardInstanceNodeUpdate,
		Read:   resourceMongodbShardInstanceNodeRead,
		Importer: &schema.ResourceImporter{
			State: mongodbShardInstanceNodeImportStateFunc(),
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		CustomizeDiff: mongodbShardInstanceCustomizeDiffFunc(),
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Optional: true,
				// set stable value before api support query
				ValidateFunc: validation.StringInSlice([]string{
					"shard",
					"mongos",
				}, false),
				Default:  "shard",
				ForceNew: true,
			},
			"node_class": {
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
				}, false),
				Default: "1C2G",
			},
			"node_storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 1000),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("node_type").(string) == "mongos" {
						return true
					}
					return false
				},
				Default: 50,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongodbShardInstanceNodeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createMongodbShardInstanceNode(d, meta)
	if err != nil {
		return fmt.Errorf("create shard instance node error %s ", err)
	}
	return resourceMongodbShardInstanceNodeRead(d, meta)
}

func resourceMongodbShardInstanceNodeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyMongodbShardInstanceNode(d, meta)
	if err != nil {
		return fmt.Errorf("update shard instance node error %s ", err)
	}
	return resourceMongodbShardInstanceNodeRead(d, meta)
}

func resourceMongodbShardInstanceNodeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	if d.Get("node_type").(string) == "shard" {
		return fmt.Errorf("can not support remove shard node from instance")
	}
	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		err = delMongodbShardInstanceNode(d, meta)
		if err == nil {
			return nil
		} else {
			_, _, err = readMongodbShardInstanceNode(d, meta)
			if err != nil && canNotFoundMongodbError(err) {
				return nil
			}
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceMongodbShardInstanceNodeRead(d *schema.ResourceData, meta interface{}) (err error) {
	v, extra, err := readMongodbShardInstanceNode(d, meta)
	if err != nil {
		return fmt.Errorf("read shard instance node error %s ", err)
	}
	SdkResponseAutoResourceData(d, resourceKsyunMongodbShardInstanceNode(), v, extra)
	return err
}
