/*
Provide RDS security group function

# Example Usage

```hcl

	resource "ksyun_krds_security_group" "default" {
	  security_group_name = "terraform_security_group_13"
	  security_group_description = "terraform-security-group-13"
	  security_group_rule{
	    security_group_rule_protocol = "182.133.0.0/16"
	    security_group_rule_name = "asdf"
	  }
	  security_group_rule{
	    security_group_rule_protocol = "182.134.0.0/16"
	    security_group_rule_name = "asdf2"
	  }
	}

```

# Import

RDS security group can be imported using the id, e.g.

```
$ terraform import ksyun_krds_security_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceKsyunKrdsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKrdsSecurityGroupCreate,
		Update: resourceKsyunKrdsSecurityGroupUpdate,
		Read:   resourceKsyunKrdsSecurityGroupRead,
		Delete: resourceKsyunKrdsSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"security_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tf_krds_security_group",
				Description: "the name of the security group.",
			},
			"security_group_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tf_krds_security_group_desc",
				Description: "description of security group.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security group ID.",
			},
			"security_group_rule": {
				Type:        schema.TypeSet,
				Set:         secParameterToHash,
				Optional:    true,
				Description: "the rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the security group rule.",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the creation time of the rule.",
						},
						"security_group_rule_protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "the protocol of the rule.",
						},
						"security_group_rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "the name of the rule.",
						},
					},
				},
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation time of the resource.",
			},
		},
	}
}

func secParameterToHash(ruleMap interface{}) int {
	rule := ruleMap.(map[string]interface{})
	return hashcode.String(rule["security_group_rule_protocol"].(string) + "|" + rule["security_group_rule_name"].(string))
}

func resourceKsyunKrdsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createKrdsSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating sg , error is %e", err)
	}
	return resourceKsyunKrdsSecurityGroupRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyKrdsSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating sg , error is %e", err)
	}
	return resourceKsyunKrdsSecurityGroupRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readKrdsAndSetSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading sg , error is %e", err)
	}
	return err
}

func resourceKsyunKrdsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removeKrdsSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting sg , error is %e", err)
	}
	return err
}
