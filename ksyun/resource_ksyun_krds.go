package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceKsyunKrds() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunKrdsCreate,
		Update: resourceKsyunKrdsUpdate,
		Read:   resourceKsyunKrdsRead,
		Delete: resourceKsyunKrdsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: krdsInstanceCustomizeDiff(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Minute),
			Update: schema.DefaultTimeout(300 * time.Minute),
			Delete: schema.DefaultTimeout(300 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_identifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "source instance identifier",
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				Description: "this value regex db.ram.d{1,3}|db.disk.d{1,5} , " +
					"db.ram is rds random access memory size, db.disk is disk size",
				ValidateFunc: validDbInstanceClass(),
			},
			"db_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HRDS",
				ValidateFunc: validation.StringInSlice([]string{
					"HRDS",
					"TRDS",
					"ERDS",
					"SINGLERDS",
				}, false),
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "engine is db type, only support mysql|percona",
				ForceNew:    true,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db engine version only support 5.5|5.6|5.7|8.0",
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"master_user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
			"bill_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DAY",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DAY",
					"YEAR_MONTH",
				}, false),
			},
			"duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: durationSchemaDiffSuppressFunc("bill_type", "YEAR_MONTH"),
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "proprietary security group id for krds",
			},
			"vip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_parameter_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_has_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func parameterToHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceKsyunKrdsCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createKrdsInstance(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on creating instance , error is %e", err)
	}
	return resourceKsyunKrdsRead(d, meta)
}

func resourceKsyunKrdsRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetKrdsInstance(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on reading instance , error is %s", err)
	}
	err = readAndSetKrdsInstanceParameters(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance , error is %s", err)
	}
	return err
}

func resourceKsyunKrdsUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyKrdsInstance(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	err = checkKrdsInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	err = resourceKsyunKrdsRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	return err
}

func resourceKsyunKrdsDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removeKrdsInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting instance , error is %e", err)
	}
	return err
}
