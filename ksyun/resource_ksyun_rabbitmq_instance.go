package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunRabbitmq() *schema.Resource {
	return &schema.Resource{
		Create: resourceRabbitmqInstanceCreate,
		Read:   resourceRabbitmqInstanceRead,
		Update: resourceRabbitmqInstanceUpdate,
		Delete: resourceRabbitmqInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"1",
				}, false),
			},
			"ssd_disk": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"bill_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("bill_type"); ok && v == 1 {
						return false
					}
					return true
				},
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(3, 3),
				ForceNew:     true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_plugins": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			//"band_width": {
			//	Type:     schema.TypeInt,
			//	Optional: true,
			//	Computed: true,
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		if v, ok := d.GetOk("enable_eip"); ok && v.(bool) {
			//			return false
			//		}
			//		return true
			//	},
			//},
			"enable_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cidrs": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
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
			"product_what": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_egress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func resourceRabbitmqInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp    *map[string]interface{}
		err     error
		addCidr string
	)

	transform := map[string]SdkReqTransform{
		"force_restart": {Ignore: true},
		"cidrs":         {Ignore: true},
	}

	conn := meta.(*KsyunClient).rabbitmqconn
	r := resourceKsyunRabbitmq()
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		false,
	})
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	err = checkRabbitmqAvailabilityZone(d, meta, req)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	err, _ = checkRabbitmqPlugins(d, meta, req)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	action := "CreateInstance"
	logger.Debug(logger.ReqFormat, action, req)
	if resp, err = conn.CreateInstance(&req); err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["InstanceId"].(string))
	}
	err = allocateRabbitmqInstanceEip(d, meta)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	err, addCidr, _ = validModifyRabbitmqInstanceRules(d, resourceKsyunRabbitmq(), meta, "", false)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	err = addRabbitmqRules(d, meta, "", addCidr)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	return resourceRabbitmqInstanceRead(d, meta)
}

func resourceRabbitmqInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).rabbitmqconn

	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId"] = d.Id()

	logger.Debug(logger.ReqFormat, "DeleteRabbitmqInstance", deleteReq)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteInstance(&deleteReq)
		if err != nil {
			return resource.RetryableError(errors.New(""))
		} else {
			return nil
		}
	})

	if err != nil {
		return fmt.Errorf("error on deleting instance %q, %s", d.Id(), err)
	}

	return resource.Retry(20*time.Minute, func() *resource.RetryError {

		queryReq := make(map[string]interface{})
		queryReq["InstanceId"] = d.Id()

		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstance", queryReq)
		resp, err := conn.DescribeInstance(&queryReq)
		logger.Debug(logger.RespFormat, "DescribeRabbitmqInstance", queryReq, resp)

		if err != nil {
			if strings.Contains(err.Error(), "InstanceNotFound") {
				return nil
			} else {
				return resource.NonRetryableError(err)
			}
		}

		_, ok := (*resp)["Data"].(map[string]interface{})

		if !ok {
			return nil
		}

		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceRabbitmqInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		enable  string
		disable string
		addCidr string
		delCidr string
	)
	//validModifyRabbitmqInstancePlugins before update resource
	err, enable, disable = validModifyRabbitmqInstancePlugins(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance plugins %q, %s", d.Id(), err)
	}

	err, addCidr, delCidr = validModifyRabbitmqInstanceRules(d, resourceKsyunRabbitmq(), meta, "", true)
	if err != nil {
		return fmt.Errorf("error on update instance cidrs %q, %s", d.Id(), err)
	}

	err = modifyRabbitmqInstanceNameAndProject(d, meta)

	if err != nil {
		return fmt.Errorf("error on update instance plugins %q, %s", d.Id(), err)
	}

	err = modifyRabbitmqInstancePassword(d, meta)

	if err != nil {
		return fmt.Errorf("error on update instance plugins %q, %s", d.Id(), err)
	}

	err = allocateRabbitmqInstanceEip(d, meta)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	err = deallocateRabbitmqInstanceEip(d, meta)
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}

	err = modifyRabbitmqInstancePlugins(d, meta, enable, disable)

	if err != nil {
		return fmt.Errorf("error on update instance plugins %q, %s", d.Id(), err)
	}

	err = restartRabbitmqInstance(d, meta)
	if err != nil {
		return err
	}

	err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	err = addRabbitmqRules(d, meta, "", addCidr)
	if err != nil {
		return fmt.Errorf("error on update instance cidrs %q, %s", d.Id(), err)
	}

	_, err = deleteRabbitmqRules(d, meta, "", delCidr)
	if err != nil {
		return fmt.Errorf("error on update instance cidrs %q, %s", d.Id(), err)
	}
	return resourceRabbitmqInstanceRead(d, meta)
}

func resourceRabbitmqInstanceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		item           map[string]interface{}
		ok             bool
		err            error
		plugins        []interface{}
		pluginStr      string
		rules          []interface{}
		ruleStr        string
		currentPlugins []string
		currentRules   []string
	)

	item, err = readRabbitmqInstance(d, meta, "")
	if err != nil {
		return err
	}
	if _, ok = item["AvailabilityZone"]; ok {
		delete(item, "AvailabilityZone")
	}

	plugins, err = readRabbitmqInstancePlugins(d, meta, "")
	if err != nil {
		return err
	}
	for _, plugin := range plugins {
		if int64(plugin.(map[string]interface{})["PluginStatus"].(float64)) == 1 {
			currentPlugins = append(currentPlugins, plugin.(map[string]interface{})["PluginName"].(string))
		}
	}
	pluginStr = stringSplitRead(",", "enable_plugins", currentPlugins, d)
	if pluginStr != "" {
		item["EnablePlugins"] = pluginStr
	}

	extra := make(map[string]SdkResponseMapping)
	extra["AvailabilityZoneEn"] = SdkResponseMapping{
		Field: "availability_zone",
	}

	extra["Eip"] = SdkResponseMapping{
		Field:    "eip",
		KeepAuto: true,
		FieldRespFunc: func(i interface{}) interface{} {
			if i.(string) != "" {
				_ = d.Set("enable_eip", true)
			} else {
				_ = d.Set("enable_eip", false)
			}
			return i
		},
	}

	rules, err = readRabbitmqInstanceRules(d, meta, "")
	if err != nil {
		return err
	}
	for _, rule := range rules {
		r := rule.(map[string]interface{})["Cidr"].(string)
		currentRules = append(currentRules, r)
	}
	ruleStr = stringSplitRead(",", "cidrs", currentRules, d)
	if ruleStr != "" {
		item["Cidrs"] = ruleStr
	}
	SdkResponseAutoResourceData(d, resourceKsyunRabbitmq(), item, extra)
	return d.Set("force_restart", d.Get("force_restart"))
}
