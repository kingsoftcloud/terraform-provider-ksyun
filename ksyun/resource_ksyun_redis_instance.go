package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// instance
func resourceRedisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisInstanceCreate,
		Delete: resourceRedisInstanceDelete,
		Update: resourceRedisInstanceUpdate,
		Read:   resourceRedisInstanceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 4),
				Description:  "1-quickSelect;2-HA;3-selfDefineCluster",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "mem of redis, if mode is selfDefineCluster(mode=3) then capacity = shard_size * shard_num",
			},
			"slave_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 8),
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
			"bill_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 5, 87}),
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("bill_type"); ok && v == 1 {
						return false
					}
					return true
				},
			},
			"pass_word": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"4.0",
					"5.0",
					"6.0",
				}, false),
			},
			"backup_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
			},
			"reset_all_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			//"security_group_ids": {
			//	Type:     schema.TypeSet,
			//	Computed: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Set: schema.HashString,
			//},
			"net_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 2),
			},

			"cache_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timing_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "auto backup On or Off.",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"shard_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1024),
				Description:  "each shard mem size GB.",
			},
			"shard_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1024),
				Description:  "shard number.",
			},
			"prepare_az_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "assign standby instance area.",
			},
			"rr_az_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "assign read only instance area.",
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"sub_order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_use": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kcsv1conn
	var (
		resp *map[string]interface{}
		err  error
	)
	// valid parameters ...
	createParam, err := resourceRedisInstanceParameterCheckAndPrepare(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	r := resourceRedisInstance()
	transform := map[string]SdkReqTransform{
		"reset_all_parameters": {Ignore: true},
		"parameters":           {Ignore: true},
		"security_group_id":    {Ignore: true},
		"protocol": {ValueFunc: func(d *schema.ResourceData) (interface{}, bool) {
			v, ok := d.GetOk("protocol")
			if ok {
				return v, ok
			} else {
				mode := d.Get("mode").(int)
				switch mode {
				case 1:
					_ = d.Set("protocol", "4.0")
					return "4.0", true
				case 2:
					_ = d.Set("protocol", "4.0")
					return "4.0", true
				default:
					_ = d.Set("protocol", "4.0")
					return "4.0", true
				}
			}
		}},
	}
	//generate req
	createReq, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	//create redis instance
	action := "CreateCacheCluster"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.CreateCacheCluster(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["CacheId"].(string))
	}
	err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutCreate), "")
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	//AllocateSecurityGroup
	err = modifyRedisInstanceSg(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	if len(*createParam) > 0 {
		err = setResourceRedisInstanceParameter(d, meta, createParam)
		if err != nil {
			return fmt.Errorf("error on create Instance: %s", err)
		}
	}
	if timingSwitch, ok := d.GetOk("timing_switch"); ok && strings.ToLower(timingSwitch.(string)) == "on" {
		autoBackupReq := make(map[string]interface{})
		autoBackupReq["CacheId"] = d.Id()
		autoBackupReq["TimingSwitch"] = "On"
		if timezone, ok := d.GetOk("timezone"); ok {
			autoBackupReq["Timezone"] = timezone
		} else {
			return fmt.Errorf("error, timing_switch=on, but timezone is null. Timezone: %s", timezone)
		}
		backupAction := "SetTimingSnapshot"
		logger.Debug(logger.ReqFormat, backupAction, autoBackupReq)
		resp, err = conn.SetTimingSnapshot(&autoBackupReq)
		if err != nil {
			return fmt.Errorf("error on creating instance: %s", err)
		}
		logger.Debug(logger.RespFormat, action, autoBackupReq, *resp)
	}

	return resourceRedisInstanceRead(d, meta)
}

func resourceRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	deleteReq := make(map[string]interface{})
	deleteReq["CacheId"] = d.Id()
	deleteReq["AvailableZone"] = d.Get("az")

	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		var (
			resp *map[string]interface{}
			err  error
		)
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &deleteReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.DeleteCacheCluster(&deleteReq)
			},
		}
		action := "DeleteCacheCluster"
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err == nil {
			return nil
		}
		logger.Debug(logger.RespFormat, action, deleteReq, resp)
		_, err = describeRedisInstance(d, meta, "")
		if err != nil {
			if validateExists(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer func(d *schema.ResourceData, meta interface{}) {
		_err := resourceRedisInstanceRead(d, meta)
		if err == nil {
			err = _err
		} else {
			if _err != nil {
				err = fmt.Errorf(err.Error()+" %s", _err)
			}
		}

	}(d, meta)

	// valid parameters ...
	createParam, err := resourceRedisInstanceParameterCheckAndPrepare(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}

	// rename
	err = modifyRedisInstanceNameAndProject(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	// update password
	err = modifyRedisInstancePassword(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	//sg
	err = modifyRedisInstanceSg(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	// resize mem
	err = modifyRedisInstanceSpec(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	// auto backup time
	err = modifyRedisInstanceAutoBackup(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}

	// update parameter
	if len(*createParam) > 0 {
		err = setResourceRedisInstanceParameter(d, meta, createParam)
		if err != nil {
			return fmt.Errorf("error on create Instance: %s", err)
		}
	}
	err = d.Set("reset_all_parameters", d.Get("reset_all_parameters"))
	return err
}

func resourceRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		item map[string]interface{}
		resp *map[string]interface{}
		ok   bool
		err  error
	)
	resp, err = describeRedisInstance(d, meta, "")
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
		return nil
	}
	// merge some field
	add := make(map[string]interface{})
	for k, v := range item {
		if k == "az" {
			add["availableZone"] = v
		}
	}
	if _, ok = item["slaveNum"]; !ok {
		item["slaveNum"] = 0
	}
	for k, v := range add {
		item[k] = v
	}
	extra := make(map[string]SdkResponseMapping)
	extra["protocol"] = SdkResponseMapping{
		Field: "protocol",
		FieldRespFunc: func(i interface{}) interface{} {
			return strings.Replace(i.(string), "redis ", "", -1)
		},
	}
	extra["size"] = SdkResponseMapping{
		Field: "capacity",
	}
	SdkResponseAutoResourceData(d, resourceRedisInstance(), item, extra)

	//merge parameters
	err = resourceRedisInstanceParamRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}

	//merge securityGroupIds
	err = resourceRedisInstanceSgRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}

	return d.Set("reset_all_parameters", d.Get("reset_all_parameters"))
}
