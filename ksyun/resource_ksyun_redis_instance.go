/*
Provides an redis instance resource.

# Example Usage

```hcl

	variable "available_zone" {
	  default = "cn-beijing-6a"
	}

	variable "subnet_name" {
	  default = "ksyun_subnet_tf"
	}

	variable "vpc_name" {
	  default = "ksyun_vpc_tf"
	}

	variable "vpc_cidr" {
	  default = "10.1.0.0/21"
	}

	variable "protocol" {
	  default = "4.0"
	}

	resource "ksyun_vpc" "default" {
	  vpc_name   = "${var.vpc_name}"
	  cidr_block = "${var.vpc_cidr}"
	}

	resource "ksyun_subnet" "default" {
	  subnet_name      = "${var.subnet_name}"
	  cidr_block = "10.1.0.0/21"
	  subnet_type = "Normal"
	  dhcp_ip_from = "10.1.0.2"
	  dhcp_ip_to = "10.1.0.253"
	  vpc_id  = "${ksyun_vpc.default.id}"
	  gateway_ip = "10.1.0.1"
	  dns1 = "198.18.254.41"
	  dns2 = "198.18.254.40"
	  available_zone = "${var.available_zone}"
	}

	resource "ksyun_redis_sec_group" "default" {
	  available_zone = "${var.available_zone}"
	  name = "testTerraform777"
	  description = "testTerraform777"
	}

	resource "ksyun_redis_instance" "default" {
	  available_zone        = "${var.available_zone}"
	  name                  = "MyRedisInstance1101"
	  mode                  = 2
	  capacity              = 1
	  slave_num              = 2
	  net_type              = 2
	  vnet_id               = "${ksyun_subnet.default.id}"
	  vpc_id                = "${ksyun_vpc.default.id}"
	  security_group_id     = "${ksyun_redis_sec_group.default.id}"
	  bill_type             = 5
	  duration              = ""
	  duration_unit         = ""
	  pass_word             = "Shiwo1101"
	  iam_project_id        = "0"
	  protocol              = "${var.protocol}"
	  reset_all_parameters  = false
	  timing_switch         = "On"
	  timezone              = "07:00-08:00"
	  available_zone        = "cn-beijing-6a"
	  prepare_az_name       = "cn-beijing-6b"
	  rr_az_name            = "cn-beijing-6a"
	  parameters = {
	    "appendonly"                  = "no",
	    "appendfsync"                 = "everysec",
	    "maxmemory-policy"            = "volatile-lru",
	    "hash-max-ziplist-entries"    = "513",
	    "zset-max-ziplist-entries"    = "129",
	    "list-max-ziplist-size"       = "-2",
	    "hash-max-ziplist-value"      = "64",
	    "notify-keyspace-events"      = "",
	    "zset-max-ziplist-value"      = "64",
	    "maxmemory-samples"           = "5",
	    "set-max-intset-entries"      = "512",
	    "timeout"                     = "600",
	  }
	}

```

# Import

redis  can be imported using the `id`, e.g.

```
$ terraform import ksyun_redis_instance.example xxxxxxxxx
```
*/
package ksyun

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The Zone to launch the DB instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of DB instance.",
			},
			"mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 4),
				Description:  "The KVStore instance system architecture required by the user. Valid values:  1(cluster),2(single),3(SelfDefineCluster).",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "mem of redis, if mode is selfDefineCluster(mode=3) then capacity = shard_size * shard_num.",
			},
			"slave_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 8),
				Description:  "The readonly node num required by the user. Valid values: {0-7}.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Used to retrieve instances belong to specified VPC.",
			},
			"vnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of subnet. the instance will use the subnet in the current region.",
			},
			"bill_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 5, 87}),
				Description:  "Valid values are 1 (Monthly), 5(Daily), 87(HourlyInstantSettlement).",
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
				Description: "Only meaningful if bill_type is 1, Valid values:{1~36}.",
			},
			"pass_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password of the  instance.The password is a string of 8 to 30 characters and must contain uppercase letters, lowercase letters, and numbers.",
			},
			"iam_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project instance belongs to.",
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
				Description: "Engine version. Supported values: 2.8, 4.0 and 5.0.",
			},
			"backup_time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Auto backup time zone. Example: \"03:00-04:00\".",
			},
			"security_group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
				Description:      "The id of security group.",
			},
			"reset_all_parameters": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether reset all parameters.",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        schema.TypeString,
				Description: "Set of parameters needs to be set after instance was launched. Available parameters can refer to the  docs https://docs.ksyun.com/documents/1018.",
			},
			// "security_group_ids": {
			//	Type:     schema.TypeSet,
			//	Computed: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	Set: schema.HashString,
			// },
			"net_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 2),
				Description:  "The network type. Valid values: 2(vpc).",
			},

			"cache_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of cache.",
			},
			"az": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "availability zone.",
			},
			"timing_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "auto backup On or Off.",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "Auto backup time zone. Example: \"03:00-04:00\".",
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
			"delete_directly": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Default is `false`, deleted instance will remain in the recycle bin. Setting the value to `true`, instance is permanently deleted without being recycled.",
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "engine.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "size.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "port.",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vip.",
			},
			"slave_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "slave vip.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "status.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time.",
			},
			"used_memory": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "used memory.",
			},
			"sub_order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "sub order ID.",
			},
			"product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "project id.",
			},
			"order_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "order type.",
			},
			"order_use": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "order use.",
			},
			"source": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "source.",
			},
			"service_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "service status.",
			},
			"service_begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "service begin time.",
			},
			"service_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "service end time.",
			},
			"iam_project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "project name.",
			},

			"tags": tagsSchema(),
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
	// d.Set("delete_directly", d.Get("delete_directly"))

	// logger.Debug(logger.ReqFormat, "delete_directly", d.Get("delete_directly"))
	createParam, err := resourceRedisInstanceParameterCheckAndPrepare(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	r := resourceRedisInstance()
	transform := map[string]SdkReqTransform{
		"reset_all_parameters": {Ignore: true},
		"delete_directly":      {Ignore: true},
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
	// generate req
	createReq, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	// create redis instance
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
	// AllocateSecurityGroup
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

	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunKrds(), "redis-instance", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
	}
	return resourceRedisInstanceRead(d, meta)
}

func resourceRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	deleteReq := make(map[string]interface{})
	deleteReq["CacheId"] = d.Id()
	deleteReq["AvailableZone"] = d.Get("az")
	deleteReq["DeleteDirectly"] = d.Get("delete_directly")

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
	// sg
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

	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunKrds(), "redis-instance", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
	}
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

	if _, ok := d.GetOk("tags"); ok {
		err = mergeTagsData(d, &item, meta.(*KsyunClient), "redis-instance")
		if err != nil {
			return fmt.Errorf("reading tags error: %s", err)
		}
	}
	SdkResponseAutoResourceData(d, resourceRedisInstance(), item, extra)

	// merge parameters
	err = resourceRedisInstanceParamRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}

	// merge securityGroupIds
	err = resourceRedisInstanceSgRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}

	return d.Set("reset_all_parameters", d.Get("reset_all_parameters"))
}
