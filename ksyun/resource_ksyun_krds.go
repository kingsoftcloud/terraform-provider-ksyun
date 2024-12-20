/*
Provides an RDS instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

# Example Usage

```hcl

# Create a RDS MySQL instance

	provider "ksyun"{
	  region = "cn-shanghai-3"
	  access_key = ""
	  secret_key = ""
	}

	variable "available_zone" {
	  default = "cn-shanghai-3a"
	}

	resource "ksyun_vpc" "default" {
	  vpc_name   = "ksyun-vpc-tf"
	  cidr_block = "10.7.0.0/21"
	}

	resource "ksyun_subnet" "foo" {
	  subnet_name      = "ksyun-subnet-tf"
	  cidr_block = "10.7.0.0/21"
	  subnet_type = "Reserve"
	  dhcp_ip_from = "10.7.0.2"
	  dhcp_ip_to = "10.7.0.253"
	  vpc_id  = "${ksyun_vpc.default.id}"
	  gateway_ip = "10.7.0.1"
	  dns1 = "198.18.254.41"
	  dns2 = "198.18.254.40"
	  availability_zone = "${var.available_zone}"
	}

	resource "ksyun_krds_security_group" "krds_sec_group_14" {
	  security_group_name = "terraform_security_group_14"
	  security_group_description = "terraform-security-group-14"
	  security_group_rule{
	    security_group_rule_protocol = "182.133.0.0/16"
	    security_group_rule_name = "asdf"
	  }
	  security_group_rule{
	    security_group_rule_protocol = "182.134.0.0/16"
	    security_group_rule_name = "asdf2"
	  }
	}

	resource "ksyun_krds" "my_rds_xx"{
	  db_instance_class= "db.ram.2|db.disk.21"
	  db_instance_name = "houbin_terraform_1-n"
	  db_instance_type = "HRDS"
	  engine = "mysql"
	  engine_version = "5.7"
	  master_user_name = "admin"
	  master_user_password = "123qweASD123"
	  vpc_id = "${ksyun_vpc.default.id}"
	  subnet_id = "${ksyun_subnet.foo.id}"
	  bill_type = "DAY"
	  security_group_id = "${ksyun_krds_security_group.krds_sec_group_14.id}"
	  preferred_backup_time = "01:00-02:00"
	  availability_zone_1 = "cn-shanghai-3a"
	  availability_zone_2 = "cn-shanghai-3b"
	  port=3306
	}

# Create a RDS MySQL instance with specific parameters

	provider "ksyun"{
	  region = "cn-shanghai-3"
	  access_key = ""
	  secret_key = ""
	}

	variable "available_zone" {
	  default = "cn-shanghai-3a"
	}

	resource "ksyun_vpc" "default" {
	  vpc_name   = "ksyun-vpc-tf"
	  cidr_block = "10.7.0.0/21"
	}

	resource "ksyun_subnet" "foo" {
	  subnet_name      = "ksyun-subnet-tf"
	  cidr_block = "10.7.0.0/21"
	  subnet_type = "Reserve"
	  dhcp_ip_from = "10.7.0.2"
	  dhcp_ip_to = "10.7.0.253"
	  vpc_id  = "${ksyun_vpc.default.id}"
	  gateway_ip = "10.7.0.1"
	  dns1 = "198.18.254.41"
	  dns2 = "198.18.254.40"
	  availability_zone = "${var.available_zone}"
	}

	resource "ksyun_krds_security_group" "krds_sec_group_14" {
	  output_file = "output_file"
	  security_group_name = "terraform_security_group_14"
	  security_group_description = "terraform-security-group-14"
	  security_group_rule{
	    security_group_rule_protocol = "182.133.0.0/16"
	    security_group_rule_name = "asdf"
	  }
	  security_group_rule{
	    security_group_rule_protocol = "182.134.0.0/16"
	    security_group_rule_name = "asdf2"
	  }
	}

	resource "ksyun_krds" "my_rds_xx"{
	  output_file = "output_file"
	  db_instance_class= "db.ram.2|db.disk.21"
	  db_instance_name = "houbin_terraform_1-n"
	  db_instance_type = "HRDS"
	  engine = "mysql"
	  engine_version = "5.7"
	  master_user_name = "admin"
	  master_user_password = "123qweASD123"
	  vpc_id = "${ksyun_vpc.default.id}"
	  subnet_id = "${ksyun_subnet.foo.id}"
	  bill_type = "DAY"
	  security_group_id = "${ksyun_krds_security_group.krds_sec_group_14.id}"
	  preferred_backup_time = "01:00-02:00"
	  parameters {
	    name = "auto_increment_increment"
	    value = "8"
	  }

	  parameters {
	    name = "binlog_format"
	    value = "ROW"
	  }

	  parameters {
	    name = "delayed_insert_limit"
	    value = "108"
	  }
	  parameters {
	    name = "auto_increment_offset"
	    value= "2"
	  }
	  availability_zone_1 = "cn-shanghai-3a"
	  availability_zone_2 = "cn-shanghai-3b"
	  instance_has_eip = true
	}

```

# Import

KRDS can be imported using the id, e.g.

```
$ terraform import ksyun_krds.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Description: "instance ID.",
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				Description: "this value regex db.ram.d{1,9}|db.disk.d{1,9}, " +
					"db.ram is rds random access memory size, db.disk is disk size.",
				ValidateFunc: validDbInstanceClass(),
			},
			"db_instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance name.",
			},
			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance type, valid values: HRDS, TRDS, ERDS, SINGLERDS.",
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
				Description: "engine is db type, only support mysql|percona.",
				ForceNew:    true,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db engine version only support 5.5|5.6|5.7|8.0.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "region code.",
			},
			"master_user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "database primary account name.",
			},
			"master_user_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "master account password.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of th VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the subnet.",
			},
			"bill_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DAY",
				ForceNew: true,
				//				//ValidateFunc: validation.StringInSlice([]string{
				//	"DAY",
				//	"YEAR_MONTH",
				// }, false),
				Description: "bill type, valid values: DAY, YEAR_MONTH, HourlyInstantSettlement. Default is DAY.",
			},
			"duration": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: durationSchemaDiffSuppressFunc("bill_type", "YEAR_MONTH"),
				Description:      "purchase duration in months.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "proprietary security group id for krds.",
			},
			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "virtual IP.",
			},

			// 参数组不能手动指定
			// 创建时如果有parameters，会创建一个临时的参数组，创建实例时传入，实例创建完毕删除
			// ！！！如果有指定的需求，需要注意改动清理临时参数组的逻辑，避免误删
			"db_parameter_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the parameter group.",
			},
			"preferred_backup_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "backup time.",
			},
			"availability_zone_1": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "zone 1.",
			},
			"availability_zone_2": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "zone 2.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "project ID.",
			},
			"db_parameter_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "the id of parameter template that for the krds being created.",
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of the parameter.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value of the parameter.",
						},
					},
				},
				Set:         parameterToHash,
				Optional:    true,
				Computed:    true,
				Description: "database parameters.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "port number.",
			},
			"instance_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance create time.",
			},
			"instance_has_eip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "attach eip for instance.",
			},
			"eip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EIP address.",
			},
			"eip_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "EIP port.",
			},
			"force_restart": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set it to true to make some parameter efficient when modifying them. Default to false.",
			},

			"tags": tagsSchema(),
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
	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunKrds(), "krds", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
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
	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunKrds(), "krds", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
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
