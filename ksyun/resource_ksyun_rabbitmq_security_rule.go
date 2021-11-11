package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
	"time"
)

func resourceKsyunRabbitmqSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRabbitmqSecurityRuleCreate,
		Read:   resourceRabbitmqSecurityRuleRead,
		Update: resourceRabbitmqSecurityRuleUpdate,
		Delete: resourceRabbitmqSecurityRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return conflictResourceImport("instance_id", "cidr", "cidrs", d)
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cidrs"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) (r bool) {
					return conflictResourceDiffSuppressForSingle("cidrs", old, new, d)
				},
				ForceNew: true,
			},
			"cidrs": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) (r bool) {
					r = conflictResourceDiffSuppressForMultiple("cidr", "cidrs", d)
					if r {
						return r
					}
					return stringSplitDiffSuppressFunc(",")(k, old, new, d)
				},
				ValidateFunc:  stringSplitSchemaValidateFunc(","),
				Deprecated:    "`cidrs` is deprecated use resourceKsyunRabbitmq.cidrs instead ",
				ConflictsWith: []string{"cidr"},
			},
		},
	}
}

func resourceRabbitmqSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		err error
		add string
		del string
	)
	err, add, del = validModifyRabbitmqInstanceRules(d, resourceKsyunRabbitmqSecurityRule(), meta, d.Get("instance_id").(string), true)
	if err != nil {
		return fmt.Errorf("error on update rabbit Instance sg rule: %s", err)
	}
	err = addRabbitmqRules(d, meta, d.Get("instance_id").(string), add)
	if err != nil {
		return fmt.Errorf("error on update rabbit Instance sg rule: %s", err)
	}
	_, err = deleteRabbitmqRules(d, meta, d.Get("instance_id").(string), del)
	if err != nil {
		return fmt.Errorf("error on update rabbit Instance sg rule: %s", err)
	}
	return resourceRabbitmqSecurityRuleRead(d, meta)
}

func resourceRabbitmqSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		use string
		err error
		add string
	)
	use, err = checkConflictOnCreate("cidr", "cidrs", d)
	if err != nil {
		return err
	}
	err, add, _ = validModifyRabbitmqInstanceRules(d, resourceKsyunRabbitmqSecurityRule(), meta, d.Get("instance_id").(string), false)
	if err != nil {
		return fmt.Errorf("error on create rabbit Instance sg rule: %s", err)
	}
	err = addRabbitmqRules(d, meta, d.Get("instance_id").(string), add)
	if err != nil {
		return fmt.Errorf("error on create rabbit Instance sg rule: %s", err)
	}
	conflictResourceSetId(use, "instance_id", "cidr", "cidrs", d)
	return resourceRabbitmqSecurityRuleRead(d, meta)
}

func resourceRabbitmqSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		err error
		del string
	)
	if checkMultipleExist("cidrs", d) {
		del = d.Get("cidrs").(string)
	} else {
		del = d.Get("cidr").(string)
	}

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		_, err = deleteRabbitmqRules(d, meta, d.Get("instance_id").(string), del)
		if err == nil {
			return nil
		}
		if err != nil && inUseError(err) {
			return resource.RetryableError(err)
		}
		return nil
	})
}

func resourceRabbitmqSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {

	var (
		data   []interface{}
		err    error
		result string
	)
	data, err = readRabbitmqInstanceRules(d, meta, d.Get("instance_id").(string))
	if err != nil {
		return err
	}
	if checkMultipleExist("cidrs", d) {
		for _, v := range data {
			group := v.(map[string]interface{})
			if strings.Contains(d.Get("cidrs").(string), group["Cidr"].(string)) {
				result = result + group["Cidr"].(string) + ","
			}
		}
		if result != "" {
			result = result[0 : len(result)-1]
		}
		err = d.Set("cidrs", result)
	} else {
		if !checkValueInSliceMap(data, "Cidr", d.Get("cidr")) {
			return fmt.Errorf("can not read cidr [%s] from rabbitmq instance [%s]", d.Get("cidr"),
				d.Get("instance_id").(string))
		}
		result = d.Get("cidr").(string)
		err = d.Set("cidr", result)
	}

	return err
}
