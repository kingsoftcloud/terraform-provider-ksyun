package ksyun

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
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
				Type:     schema.TypeString,
				Required: true,
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
				Default: "3.2",
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
				Default: "1C2G",
			},
			"storage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
				Default:      5,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_account": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "root",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"root",
				}, false),
			},
			"instance_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
				ForceNew: true,
			},
			"duration": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: durationSchemaDiffSuppressFunc("pay_type", "byMonth"),
				ForceNew:         true,
			},
			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
				ValidateFunc: validation.IntInSlice([]int{
					3, 5, 7,
				}),
			},
			"availability_zone": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
			},
			"cidrs": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "VPC",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VPC",
				}, false),
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"timing_switch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_cycle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_what": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mongos_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shard_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongodbInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createMongodbInstanceCommon(d, meta, resourceKsyunMongodbInstance())
	if err != nil {
		return err
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
	return resourceMongodbInstanceRead(d, meta)
}

func resourceMongodbInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	return readMongodbInstanceCommon(d, meta, resourceKsyunMongodbInstance())
}
