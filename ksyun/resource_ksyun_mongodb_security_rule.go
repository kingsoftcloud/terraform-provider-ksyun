package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceKsyunMongodbSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbSecurityRuleCreate,
		Delete: resourceMongodbSecurityRuleDelete,
		Update: resourceMongodbSecurityRuleUpdate,
		Read:   resourceMongodbSecurityRuleRead,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return conflictResourceImport("instance_id", "cidr", "cidrs", d)
			},
		},
		CustomizeDiff: conflictResourceCustomizeDiffFunc("cidr", "cidrs"),
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cidrs"},
				ForceNew:      true,
			},
			"cidrs": {
				Type:             schema.TypeString,
				Optional:         true,
				Deprecated:       "`cidrs` is deprecated use resourceKsyunMongodbInstance.cidrs or resourceKsyunMongodbShardInstance instead ",
				ConflictsWith:    []string{"cidr"},
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
			},
		},
	}
}

func resourceMongodbSecurityRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		use   string
		addV4 string
		addV6 string
	)
	use, err = checkConflictOnCreate("cidr", "cidrs", d)
	if err != nil {
		return err
	}
	err, addV4, _, addV6, _ = checkMongodbSecurityGroupRulesChange(d, meta, use, d.Get("instance_id").(string))
	if err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	err = addMongodbSecurityGroupRules(d, meta, d.Get("instance_id").(string), addV4, addV6)
	if err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	conflictResourceSetId(use, "instance_id", "cidr", "cidrs", d)
	return resourceMongodbSecurityRuleRead(d, meta)
}

func resourceMongodbSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		err   error
		delV4 string
		delV6 string
	)
	if checkMultipleExist("cidrs", d) {
		err, delV4, delV6 = checkMongodbSecurityGroupRulesDel(d, meta, d.Get("instance_id").(string), d.Get("cidrs").(string))
	} else {
		err, delV4, delV6 = checkMongodbSecurityGroupRulesDel(d, meta, d.Get("instance_id").(string), d.Get("cidr").(string))
	}

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		err = delMongodbSecurityGroupRules(d, meta, d.Get("instance_id").(string), delV4, delV6)
		if err == nil {
			return nil
		}
		if err != nil && inUseError(err) {
			return resource.RetryableError(err)
		}
		return nil
	})
}

func resourceMongodbSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		addV4 string
		delV4 string
		addV6 string
		delV6 string
	)
	err, addV4, delV4, addV6, delV6 = checkMongodbSecurityGroupRulesChange(d, meta, "cidrs", d.Get("instance_id").(string))
	if err != nil {
		return fmt.Errorf("error on update instance security rules: %s", err)
	}
	err = addMongodbSecurityGroupRules(d, meta, d.Get("instance_id").(string), addV4, addV6)
	if err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	err = delMongodbSecurityGroupRules(d, meta, d.Get("instance_id").(string), delV4, delV6)
	if err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	return resourceMongodbSecurityRuleRead(d, meta)
}

func resourceMongodbSecurityRuleRead(d *schema.ResourceData, meta interface{}) (err error) {

	var (
		cidrs string
	)
	if checkMultipleExist("cidrs", d) {
		cidrs, err = readMongodbSecurityGroupCidrs(d, meta, "cidrs", d.Get("instance_id").(string))
		if err != nil {
			return fmt.Errorf("error on read instance security rule: %s", err)
		}
		if cidrs != "" {
			err = d.Set("cidrs", cidrs)
		}
	} else {
		cidrs, err = readMongodbSecurityGroupCidrs(d, meta, "cidr", d.Get("instance_id").(string))
		if err != nil {
			return fmt.Errorf("error on read instance security rule: %s", err)
		}
		if cidrs == "" {
			return fmt.Errorf("can not read cidr [%s] from mongodb instance [%s]", d.Get("cidr"),
				d.Get("instance_id").(string))
		}
		err = d.Set("cidr", cidrs)
	}
	return err
}
