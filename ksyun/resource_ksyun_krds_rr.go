/*
Provides an RDS Read Only instance resource. A DB read only instance is an isolated database environment in the cloud.

# Example Usage

```hcl

	resource "ksyun_krds_rr" "my_rds_rr"{
	  db_instance_identifier= "******"
	  db_instance_class= "db.ram.2|db.disk.50"
	  db_instance_name = "houbin_terraform_888_rr_1"
	  bill_type = "DAY"
	  security_group_id = "******"

	  parameters {
	    name = "auto_increment_increment"
	    value = "7"
	  }

	  parameters {
	    name = "binlog_format"
	    value = "ROW"
	  }
	}

```

# Import

RDS Read Only instance resource can be imported using the id, e.g.

```
$ terraform import ksyun_krds_rr.my_rds_rr 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var krdsRrNotSupport = []string{
	"db_instance_identifier",
	"availability_zone_2",
	"master_user_name",
	"master_user_password",
	"engine",
	"engine_version",
	"db_instance_type",
	"vpc_id",
	"subnet_id",
	"preferred_backup_time",
	"availability_zone_1",
	"db_instance_class",
	"db_parameter_template_id",
}

func resourceKsyunKrdsRr() *schema.Resource {
	rrSchema := resourceKsyunKrds().Schema
	for key := range rrSchema {
		for _, n := range krdsRrNotSupport {
			if key == n {
				delete(rrSchema, key)
			}
		}
	}
	rrSchema["db_instance_identifier"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "passes in the instance ID of the RDS highly available instance. A RDS highly available instance can have at most three read-only instances.",
		ForceNew:    true,
	}
	rrSchema["db_instance_class"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		// ForceNew:     true,
		ValidateFunc: validDbInstanceClass(),
		Description:  "this value regex db.ram.d{1,3}|db.disk.d{1,5}, db.ram is rds random access memory size, db.disk is disk size.",
	}
	rrSchema["db_instance_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "instance type, valid values: HRDS, TRDS, ERDS, SINGLERDS.",
	}
	rrSchema["availability_zone_1"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "zone 1.",
	}

	rrSchema["engine"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "engine is db type, only support mysql|percona.",
	}
	rrSchema["engine_version"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "db engine version only support 5.5|5.6|5.7|8.0.",
	}

	return &schema.Resource{
		Create: resourceKsyunKrdsRrCreate,
		Update: resourceKsyunKrdsRrUpdate,
		Read:   resourceKsyunKrdsRrRead,
		Delete: resourceKsyunKrdsRrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: rrSchema,
	}
}

func resourceKsyunKrdsRrCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createKrdsInstance(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on creating rr instance , error is %e", err)
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

	return resourceKsyunKrdsRrRead(d, meta)
}

func resourceKsyunKrdsRrUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyKrdsInstance(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on updating rr instance , error is %e", err)
	}
	err = checkKrdsInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error on updating rr instance , error is %e", err)
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

	err = resourceKsyunKrdsRrRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating rr instance , error is %e", err)
	}
	return err
}
func resourceKsyunKrdsRrRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetKrdsInstance(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on reading rr instance , error is %s", err)
	}
	err = readAndSetKrdsInstanceParameters(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading rr instance , error is %s", err)
	}
	return err
}
func resourceKsyunKrdsRrDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removeKrdsInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting rr instance , error is %e", err)
	}
	return err
}
