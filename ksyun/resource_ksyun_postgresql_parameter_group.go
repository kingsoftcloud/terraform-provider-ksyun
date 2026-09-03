/*
Provides a PostgreSQL parameter template group.

# Example Usage

```hcl

	resource "ksyun_postgresql_parameter_group" "dpg" {
		name = "tf_dpg_on_hcl"
		description = "tf_configuration_test"
		engine = "postgresql"
		engine_version = "15"
		parameters {
	    	name = "auto_increment_increment"
	    	value = "8"
		}
		parameters {
			name = "binlog_format"
			value = "ROW"
		}
		parameters {
			name = "delayed_insert_limit"
			value = "108"
		}
		parameters {
			name = "auto_increment_offset"
			value= "2"
		}
	}

```

# Import

Tag can be imported using the `id`, e.g.

```
$ terraform import ksyun_postgresql_parameter_group.foo "id"
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunPostgresqlParameterGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceKsyunPostgresqlParameterGroupRead,
		Create: resourceKsyunPostgresqlParameterGroupCreate,
		Delete: resourceKsyunPostgresqlParameterGroupDelete,
		Update: resourceKsyunPostgresqlParameterGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// parameter
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of postgresql parameter group.",
			},
			"db_parameter_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of postgresql parameter group.",
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of the parameter.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value of the parameter.",
						},
					},
				},
				Set:         parameterToHash,
				Optional:    true,
				Computed:    true,
				Description: "database parameters.",
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "postgresql database engine.",
			},

			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10",
					"11",
					"12.5",
					"13",
					"14",
					"15",
					"17",
				}, false),
				Description: "postgresql database version, valid values: 10|11|12.5|13|14|15|17.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				// Default:     "",
				ValidateFunc: validateName,
				Description:  "The description of this db parameter group.",
			},

			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "identify this resource.",
			},
		},
	}
}

func resourceKsyunPostgresqlParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	postgresqlParameterSrv := NewPostgresqlParameterSrv(meta.(*KsyunClient))
	r := resourceKsyunPostgresqlParameterGroup()

	reqParameters := make(map[string]interface{})

	// call query function
	reqParameters["DBParameterGroupId"] = d.Id()
	action := "DescribeDBParameterGroupById"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdkResponse, err := postgresqlParameterSrv.describeDBParameterGroupById(reqParameters)
	if err != nil || len(sdkResponse) < 1 {
		if err != nil {
			if notFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("while query db parameter group have encountered an error, detail: %s", err)
		}
		d.SetId("")
		return nil
	}
	if err := TransformMapValue2StringWithKey("Parameters", sdkResponse); err != nil {
		return err
	}
	data := sdkResponse[0]

	extra := map[string]SdkResponseMapping{
		"DBParameterGroupName": {
			Field: "name",
		},
		"Parameters": {
			FieldRespFunc: func(i interface{}) interface{} {
				if parameter, ok := i.(map[string]interface{}); ok {
					var remote = make([]map[string]interface{}, 0, len(parameter))

					for k, v := range parameter {
						m := make(map[string]interface{})
						m["name"] = k
						m["value"] = fmt.Sprintf("%v", v)
						remote = append(remote, m)
					}
					return remote
				}
				return nil
			},
		},
	}

	SdkResponseAutoResourceData(d, r, data, extra)

	return nil
}

func resourceKsyunPostgresqlParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	postgresqlParameterSrv := NewPostgresqlParameterSrv(meta.(*KsyunClient))

	postgresqlEngine := d.Get("engine").(string)
	postgresqlEngineVersion := d.Get("engine_version").(string)

	reqParameters, _, err := checkAndProcessPostgresqlParameters(d, meta)
	if err != nil {
		return err
	}

	reqParameters["DBParameterGroupName"] = d.Get("name")
	reqParameters["Engine"] = postgresqlEngine
	reqParameters["EngineVersion"] = postgresqlEngineVersion
	reqParameters["Description"] = d.Get("description")

	action := "CreateAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	dbParameterId, err := postgresqlParameterSrv.createDBParameterGroup(reqParameters)
	if err != nil {
		return err
	}

	d.SetId(dbParameterId)
	_ = d.Set("resource_name", ResourcePostgresqlParameterGroup)

	return d.Set("db_parameter_group_id", dbParameterId)
}

func resourceKsyunPostgresqlParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	postgresqlParameterSrv := NewPostgresqlParameterSrv(meta.(*KsyunClient))

	reqParameters, _, err := checkAndProcessPostgresqlParameters(d, meta)
	if err != nil {
		return err
	}

	if d.HasChange("name") {
		reqParameters["DBParameterGroupName"] = d.Get("name")
	}
	if d.HasChange("description") {
		reqParameters["Description"] = d.Get("description")
	}
	if len(reqParameters) < 1 {
		return fmt.Errorf("db parameter group only modify the parameter of name|description|parameters")
	}
	reqParameters["DBParameterGroupId"] = d.Id()

	action := "ModifyDBParameterGroup"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	if _, err := postgresqlParameterSrv.modifyDBParameterGroup(reqParameters); err != nil {
		return err
	}

	return resourceKsyunPostgresqlParameterGroupRead(d, meta)
}

func resourceKsyunPostgresqlParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	postgresqlParameterSrv := NewPostgresqlParameterSrv(meta.(*KsyunClient))

	removeMap := map[string]interface{}{
		"DBParameterGroupId": d.Id(),
	}

	action := "DeleteDBParameterGroup"
	logger.Debug(logger.ReqFormat, action, removeMap)

	return postgresqlParameterSrv.deleteDBParameterGroup(removeMap)
}
