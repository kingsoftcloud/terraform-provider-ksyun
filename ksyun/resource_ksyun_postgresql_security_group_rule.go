/*
Provide PostgreSQL security group rule

# Example Usage

```hcl

	resource "ksyun_postgresql_security_group_rule" "default" {
	  security_group_rule_protocol = "182.133.0.0/16"
	  security_group_id = "62540"
	}

```

# Import

PostgreSQL security group rule can be imported using the id, e.g.

```
$ terraform import ksyun_postgresql_security_group_rule.default ${security_group_id}:${security_group_rule_protocol}
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunPostgresqlSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceKsyunPostgresqlSecurityGroupRuleCreate,
		Read:     resourceKsyunPostgresqlSecurityGroupRuleRead,
		Delete:   resourceKsyunPostgresqlSecurityGroupRuleDelete,
		Importer: importPostgresqlSecurityGroupRule(),
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

func resourceKsyunPostgresqlSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createPostgresqlSecurityGroupRule(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating postgresql sg rule , error is %e", err)
	}
	return resourceKsyunPostgresqlSecurityGroupRuleRead(d, meta)
}

func resourceKsyunPostgresqlSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return removePostgresqlSecurityGroupRule(d, meta)
}

func resourceKsyunPostgresqlSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetPostgresqlSecurityGroupRule(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading postgresql sg rule , error is %e", err)
	}
	return err
}
