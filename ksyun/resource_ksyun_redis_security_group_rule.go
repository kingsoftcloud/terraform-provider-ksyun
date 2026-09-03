package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// redis security group rule
func resourceRedisSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupRuleCreate,
		Delete: resourceRedisSecurityGroupRuleDelete,
		Update: resourceRedisSecurityGroupRuleUpdate,
		Read:   resourceRedisSecurityGroupRuleRead,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return conflictResourceImport("security_group_id", "rule", "rules", d)
			},
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"rule": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"rules"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return conflictResourceDiffSuppressForSingle("rules", old, new, d)
				},
			},

			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				Deprecated:    "Use resourceRedisSecurityGroup().rules instead",
				ConflictsWith: []string{"rule"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return conflictResourceDiffSuppressForMultiple("rule", "rules", d)
				},
			},
		},
	}
}

func resourceRedisSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		use string
		err error
	)
	use, err = checkConflictOnCreate("rule", "rules", d)
	if err != nil {
		return err
	}
	transform := map[string]SdkReqTransform{
		"rules": {
			mapping: "Cidrs",
			Type:    TransformWithN,
		},
		"rule": {
			mapping: "Cidrs",
			Type:    TransformSingleN,
		},
	}
	err = processRedisSecurityGroupRule(d, meta, transform, false, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	conflictResourceSetId(use, "security_group_id", "rule", "rules", d)
	return resourceRedisSecurityGroupRuleRead(d, meta)
}

func resourceRedisSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
		del  []interface{}
	)
	resp, err = readRedisSecurityGroup(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		if validateRedisSgExists(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	current := make(map[string]string)
	data := (*resp)["Data"].(map[string]interface{})

	//get rule id for del
	if rules, ok := data["rules"]; ok {
		for _, r := range rules.([]interface{}) {
			rule := r.(map[string]interface{})
			del = append(del, rule["id"])
			current[rule["cidr"].(string)] = rule["id"].(string)
		}
	}

	transformDel := map[string]SdkReqTransform{
		"rules": {
			mapping: "SecurityGroupRuleId",
			Type:    TransformWithN,
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				if len(del) > 0 && checkMultipleExist("rules", d) {
					return del, true
				}
				return nil, true
			},
		},
		"rule": {
			mapping: "SecurityGroupRuleId",
			Type:    TransformSingleN,
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				if !checkMultipleExist("rules", d) {
					if len(current) > 0 {
						if id, ok := current[d.Get("rule").(string)]; ok {
							return id, true
						}
					}
				}
				return nil, true
			},
		},
	}
	err = processRedisSecurityGroupRule(d, meta, transformDel, true, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	return nil
}

func resourceRedisSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	err := updateRedisSecurityGroupRules(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	return resourceRedisSecurityGroupRuleRead(d, meta)
}

func resourceRedisSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = readRedisSecurityGroup(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		if validateRedisSgExists(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	data := (*resp)["Data"].(map[string]interface{})
	extra := make(map[string]SdkResponseMapping)
	if checkMultipleExist("rules", d) {
		extra["rules"] = SdkResponseMapping{
			Field: "rules",
			FieldRespFunc: func(i interface{}) interface{} {
				return setRedisSgCidrs(i.([]interface{}), d)
			},
		}
	} else {
		if !checkValueInSliceMap(data["rules"].([]interface{}), "cidr", d.Get("rule")) {
			d.SetId("")
			return nil
		}
		extra["rules"] = SdkResponseMapping{
			Field: "rule",
			FieldRespFunc: func(i interface{}) interface{} {
				var cidr string
				for _, v := range i.([]interface{}) {
					if v.(map[string]interface{})["cidr"].(string) == d.Get("rule").(string) {
						cidr = v.(map[string]interface{})["cidr"].(string)
						break
					}
				}
				return cidr
			},
		}
	}

	SdkResponseAutoResourceData(d, resourceRedisSecurityGroupRule(), data, extra)
	return nil
}
