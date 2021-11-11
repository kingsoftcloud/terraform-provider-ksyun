package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
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
		Description: "instance identifier",
		ForceNew:    true,
	}
	rrSchema["db_instance_class"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validDbInstanceClass(),
	}
	rrSchema["db_instance_type"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	rrSchema["availability_zone_1"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
	}

	rrSchema["engine"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	rrSchema["engine_version"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
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
