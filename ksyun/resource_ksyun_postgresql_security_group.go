/*
Provide PostgreSQL security group function

# Example Usage

```hcl

	resource "ksyun_postgresql_security_group" "default" {
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

PostgreSQL security group can be imported using the id, e.g.

```
$ terraform import ksyun_postgresql_security_group.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceKsyunPostgresqlSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunPostgresqlSecurityGroupCreate,
		Update: resourceKsyunPostgresqlSecurityGroupUpdate,
		Read:   resourceKsyunPostgresqlSecurityGroupRead,
		Delete: resourceKsyunPostgresqlSecurityGroupDelete,
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
				Default:     "tf_postgresql_security_group",
				Description: "the name of the security group.",
			},
			"security_group_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tf_postgresql_security_group_desc",
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

func resourceKsyunPostgresqlSecurityGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createPostgresqlSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating sg , error is %e", err)
	}
	return resourceKsyunPostgresqlSecurityGroupRead(d, meta)
}

func resourceKsyunPostgresqlSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyPostgresqlSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating sg , error is %e", err)
	}
	return resourceKsyunPostgresqlSecurityGroupRead(d, meta)
}

func resourceKsyunPostgresqlSecurityGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readPostgresqlAndSetSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading sg , error is %e", err)
	}
	return err
}

func resourceKsyunPostgresqlSecurityGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removePostgresqlSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting sg , error is %e", err)
	}
	return err
}
