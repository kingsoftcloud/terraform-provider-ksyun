/*
Provides a PostgreSQL Read Only instance resource. A DB read only instance is an isolated database environment in the cloud.

# Example Usage

```hcl

	resource "ksyun_postgresql_rr" "my_postgresql_rr"{
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

PostgreSQL Read Only instance resource can be imported using the id, e.g.

```
$ terraform import ksyun_postgresql_rr.my_postgresql_rr 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var postgresqlRrNotSupport = []string{
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

func resourceKsyunPostgresqlRr() *schema.Resource {
	rrSchema := resourceKsyunPostgresql().Schema
	for key := range rrSchema {
		for _, n := range postgresqlRrNotSupport {
			if key == n {
				delete(rrSchema, key)
			}
		}
	}
	rrSchema["db_instance_identifier"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "passes in the instance ID of the PostgreSQL highly available instance. A PostgreSQL highly available instance can have at most three read-only instances.",
		ForceNew:    true,
	}
	rrSchema["db_instance_class"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		// ForceNew:     true,
		ValidateFunc: validDbInstanceClass(),
		Description:  "this value regex db.ram.d{1,3}|db.disk.d{1,5}, db.ram is PostgreSQL random access memory size, db.disk is disk size.",
	}
	rrSchema["db_instance_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "instance type, valid values: HRDS_PG (highly available), RR_PG (read-only).",
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
		Description: "engine is db type, postgresql.",
	}
	rrSchema["engine_version"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "db engine version, valid values: 10|11|12.5|13|14|15|17.",
	}

	return &schema.Resource{
		Create: resourceKsyunPostgresqlRrCreate,
		Update: resourceKsyunPostgresqlRrUpdate,
		Read:   resourceKsyunPostgresqlRrRead,
		Delete: resourceKsyunPostgresqlRrDelete,
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

func resourceKsyunPostgresqlRrCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createPostgresqlInstance(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on creating rr instance , error is %e", err)
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

	return resourceKsyunPostgresqlRrRead(d, meta)
}

func resourceKsyunPostgresqlRrUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyPostgresqlInstance(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on updating rr instance , error is %e", err)
	}
	err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error on updating rr instance , error is %e", err)
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

	err = resourceKsyunPostgresqlRrRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating rr instance , error is %e", err)
	}
	return err
}
func resourceKsyunPostgresqlRrRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetPostgresqlInstance(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on reading rr instance , error is %s", err)
	}
	if d.Id() == "" {
		return nil
	}
	err = readAndSetPostgresqlInstanceParameters(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading rr instance , error is %s", err)
	}
	return err
}
func resourceKsyunPostgresqlRrDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removePostgresqlInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting rr instance , error is %e", err)
	}
	return err
}
