/*
Query ksyun postgresql security group information

# Example Usage

```hcl
# Get  postgresql_security_groups

	data "ksyun_postgresql_security_groups" "security_groups"{
	  output_file = "output_file"
	  security_group_id = 123
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func dataSourceKsyunPostgresqlSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunPostgresqlSecurityGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The filename of the content store will be returned.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of resources that satisfy the condition.",
			},
			"security_group_id": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "Security group ID.",
			},

			// 与存入数据一致datakey
			"security_groups": {
				Type: schema.TypeList,
				//Optional:    true,
				Computed:    true,
				Description: "An information list of PostgreSQL security groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"security_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Security group ID.",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Security group name.",
						},
						"security_group_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Security group description.",
						},
						"created": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The time of creation.",
						},
						"instances": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "corresponding instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_instance_identifier": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "instance ID.",
									},
									"db_instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "instance name.",
									},
									"vip": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "instance virtual IP.",
									},
									"created": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The time of creation.",
									},
									"db_instance_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "instance type.",
									},
								},
							},
						},
						"security_group_rules": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "security group rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_rule_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "rule ID.",
									},
									"security_group_rule_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "rule name.",
									},
									"security_group_rule_protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "rule protocol.",
									},
									"created": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The time of creation.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunPostgresqlSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).postgresqlconn
	descReq := make(map[string]interface{})
	if v, ok := d.GetOk("security_group_id"); ok {
		descReq["SecurityGroupId"] = fmt.Sprintf("%v", v)
	}
	action := "DescribeSecurityGroup"
	logger.DebugInfo("+-+-+-+-+  %+v  ------  %T", descReq)
	logger.Debug(logger.ReqFormat, action, descReq)
	resp, err := conn.DescribeSecurityGroup(&descReq)
	logger.Debug(logger.AllFormat, action, descReq, *resp, err)

	if err != nil {
		return fmt.Errorf("error on request instance. security group id %q, %s", d.Id(), err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading response body, security group id %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["SecurityGroups"].([]interface{})

	postgresqlIds := make([]string, len(instances))
	postgresqlMapList := make([]map[string]interface{}, len(instances))
	for num, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		postgresqlMap := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "Instances" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[Camel2Hungarian(j)] = q
						}
						wtf[num] = rrmap
					}
					postgresqlMap["instances"] = wtf
				}
			} else if k == "SecurityGroupRules" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[Camel2Hungarian(j)] = q
						}
						wtf[num] = rrmap
					}
					postgresqlMap["security_group_rules"] = wtf
				}
			} else {
				postgresqlMap[Camel2Hungarian(k)] = v
			}
		}
		logger.DebugInfo(" converted ---- %+v ", postgresqlMap)

		postgresqlIds[num] = postgresqlMap["security_group_id"].(string)
		logger.DebugInfo("postgresqlIds fuck : %v", postgresqlIds)
		postgresqlMapList[num] = postgresqlMap
	}
	//d.Set("security_group_id",postgresqlIds[0])
	logger.DebugInfo(" converted ---- %+v ", postgresqlMapList)
	_ = dataDbSave(d, "security_groups", postgresqlIds, postgresqlMapList)

	return nil
}
