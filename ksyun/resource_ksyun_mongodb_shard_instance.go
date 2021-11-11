package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceKsyunMongodbShardInstance() *schema.Resource {
	instanceSchema := resourceKsyunMongodbInstance().Schema
	subSchema := map[string]*schema.Schema{
		"shard_num": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          3,
			ValidateFunc:     validation.IntBetween(2, 32),
			DiffSuppressFunc: mongodbShardInstanceSchemaDiffSuppressFunc(),
		},
		"mongos_num": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          2,
			ValidateFunc:     validation.IntBetween(2, 32),
			DiffSuppressFunc: mongodbShardInstanceSchemaDiffSuppressFunc(),
		},
		"shard_class": {
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
			Default:          "1C2G",
			DiffSuppressFunc: mongodbShardInstanceSchemaDiffSuppressFunc(),
		},
		"mongos_class": {
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
			Default:          "1C2G",
			DiffSuppressFunc: mongodbShardInstanceSchemaDiffSuppressFunc(),
		},
		"storage": {
			Type:             schema.TypeInt,
			Optional:         true,
			ValidateFunc:     validation.IntBetween(5, 1000),
			Default:          50,
			DiffSuppressFunc: mongodbShardInstanceSchemaDiffSuppressFunc(),
		},
		"db_version": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"3.2",
				"3.6",
			}, false),
			Default:  "3.2",
			ForceNew: true,
		},
		"node_num": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"total_storage": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"instance_class": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	for k, v := range instanceSchema {
		if _, ok := subSchema[k]; ok {
			continue
		}
		subSchema[k] = v
	}
	return &schema.Resource{
		Create: resourceMongodbShardInstanceCreate,
		Delete: resourceMongodbInstanceDelete,
		Update: resourceMongodbShardInstanceUpdate,
		Read:   resourceMongodbShardInstanceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		CustomizeDiff: mongodbShardInstanceCustomizeDiffFunc(),
		Schema:        subSchema,
	}
}

func resourceMongodbShardInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createMongodbInstanceCommon(d, meta, resourceKsyunMongodbShardInstance())
	if err != nil {
		return err
	}
	return resourceMongodbShardInstanceRead(d, meta)
}

func resourceMongodbShardInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyMongodbInstanceCommon(d, meta, resourceKsyunMongodbShardInstance())
	if err != nil {
		return err
	}
	return resourceMongodbShardInstanceRead(d, meta)
}

func resourceMongodbShardInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	return readMongodbInstanceCommon(d, meta, resourceKsyunMongodbShardInstance())
}
