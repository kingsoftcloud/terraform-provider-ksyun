/*
Provide RDS security group rule

# Example Usage

```hcl

	resource "ksyun_krds_security_group_rule" "default" {
	  security_group_rule_protocol = "182.133.0.0/16"
	  security_group_id = "62540"
	}

```

# Import

RDS security group rule can be imported using the id, e.g.

```
$ terraform import ksyun_krds_security_group_rule.default ${security_group_id}:${security_group_rule_protocol}
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKrdsSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceKsyunKrdsSecurityGroupRuleCreate,
		Read:     resourceKsyunKrdsSecurityGroupRuleRead,
		Delete:   resourceKsyunKrdsSecurityGroupRuleDelete,
		Importer: importKrdsSecurityGroupRule(),
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security group id.",
			},
			"security_group_rule_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security group rule protocol.",
			},
			"security_group_rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "security group rule name.",
			},
			"security_group_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group rule id.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation time.",
			},
		},
	}
}

func resourceKsyunKrdsSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createKrdsSecurityGroupRule(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating krds sg rule , error is %e", err)
	}
	return resourceKsyunKrdsSecurityGroupRuleRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return removeKrdsSecurityGroupRule(d, meta)
}

func resourceKsyunKrdsSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetKrdsSecurityGroupRule(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading krds sg rule , error is %e", err)
	}
	return err
}
