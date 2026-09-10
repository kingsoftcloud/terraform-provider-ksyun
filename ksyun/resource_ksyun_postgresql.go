/*
Provides a PostgreSQL instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

# Example Usage

```hcl

# Create a PostgreSQL instance

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

	resource "ksyun_postgresql_security_group" "postgresql_sec_group_14" {
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

	resource "ksyun_postgresql" "my_postgresql_xx"{
	  db_instance_class= "db.ram.2|db.disk.21"
	  db_instance_name = "houbin_terraform_1-n"
	  db_instance_type = "HRDS_PG"
	  engine = "postgresql"
	  engine_version = "15"
	  master_user_name = "admin"
	  master_user_password = "123qweASD123"
	  vpc_id = "${ksyun_vpc.default.id}"
	  subnet_id = "${ksyun_subnet.foo.id}"
	  bill_type = "DAY"
	  security_group_id = "${ksyun_postgresql_security_group.postgresql_sec_group_14.id}"
	  preferred_backup_time = "01:00-02:00"
	  availability_zone_1 = "cn-shanghai-3a"
	  availability_zone_2 = "cn-shanghai-3b"
	  port=5432
	}

# Create a PostgreSQL instance with specific parameters

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

	resource "ksyun_postgresql_security_group" "postgresql_sec_group_14" {
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

	resource "ksyun_postgresql" "my_postgresql_xx"{
	  output_file = "output_file"
	  db_instance_class= "db.ram.2|db.disk.21"
	  db_instance_name = "houbin_terraform_1-n"
	  db_instance_type = "HRDS_PG"
	  engine = "postgresql"
	  engine_version = "15"
	  master_user_name = "admin"
	  master_user_password = "123qweASD123"
	  vpc_id = "${ksyun_vpc.default.id}"
	  subnet_id = "${ksyun_subnet.foo.id}"
	  bill_type = "DAY"
	  security_group_id = "${ksyun_postgresql_security_group.postgresql_sec_group_14.id}"
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

PostgreSQL can be imported using the id, e.g.

```
$ terraform import ksyun_postgresql.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunPostgresql() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunPostgresqlCreate,
		Update: resourceKsyunPostgresqlUpdate,
		Read:   resourceKsyunPostgresqlRead,
		Delete: resourceKsyunPostgresqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: postgresqlInstanceCustomizeDiff(),
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
					"db.ram is postgresql random access memory size, db.disk is disk size.",
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
				ForceNew:    true,
				Description: "instance type, valid values: HRDS_PG (highly available), RR_PG (read-only).",
				ValidateFunc: validation.StringInSlice([]string{
					"HRDS_PG",
					"RR_PG",
				}, false),
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "engine is db type, postgresql.",
				ForceNew:    true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10",
					"11",
					"12.5",
					"13",
					"14",
					"15",
					"17",
				}, false),
				Description: "db engine version, valid values: 10|11|12.5|13|14|15|17.",
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
				Description: "proprietary security group id for postgresql.",
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
				Description: "the id of parameter template that for the postgresql being created.",
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

			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The number of vCPUs for the DB instance. If not specified, defaults to half of the memory size.",
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceKsyunPostgresqlCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createPostgresqlInstance(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on creating instance , error is %e", err)
	}
	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunPostgresql(), "postgresql-instance", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
	}
	return resourceKsyunPostgresqlRead(d, meta)
}

func resourceKsyunPostgresqlRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetPostgresqlInstance(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on reading instance , error is %s", err)
	}
	if d.Id() == "" {
		return nil
	}
	err = readAndSetPostgresqlInstanceParameters(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance , error is %s", err)
	}
	return err
}

func resourceKsyunPostgresqlUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyPostgresqlInstance(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	client := meta.(*KsyunClient)
	if d.HasChange("tags") {
		tagService := TagService{client}
		tagCall, err := tagService.ReplaceResourcesTagsWithResourceCall(d, resourceKsyunPostgresql(), "postgresql-instance", false, true)
		if err != nil {
			return err
		}
		if err = tagCall.RightNow(d, client, false); err != nil {
			return fmt.Errorf("touching tags error: %s", err)
		}
	}

	err = resourceKsyunPostgresqlRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	return err
}

func resourceKsyunPostgresqlDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removePostgresqlInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting instance , error is %e", err)
	}
	return err
}
