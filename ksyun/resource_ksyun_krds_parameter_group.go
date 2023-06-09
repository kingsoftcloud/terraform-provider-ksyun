/*
Provides a tag resource.

# Example Usage

```hcl

	resource "ksyun_krds_parameter_group" "dpg1" {
		name = "tf_dpg_on_hcl"
		description = "tf configuration test"
		engine = "mysql"
		engine_version = "5.7"
		parameters = {
			connect_timeout = 20
			innodb_stats_on_metadata = "OFF"
			table_open_cache_instances = 1
			group_concat_max_len = 102
			max_connect_errors = 2000
			max_prepared_stmt_count = 65535
			max_user_connections = 65535
		}
	}

	data "ksyun_krds_parameter_group" "foo" {
		output_file = "output_result"
		db_parameter_group_id = ksyun_krds_parameter_group.dpg1.id
	}

	output "dpg_out" {
		value = data.ksyun_krds_parameter_group.foo
	}

```

# Import

Tag can be imported using the `id`, e.g.

```
$ terraform import ksyun_krds_parameter_group.foo "id"
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunKrdsParameterGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceKsyunKrdsParameterGroupRead,
		Create: resourceKsyunKrdsParameterGroupCreate,
		Delete: resourceKsyunKrdsParameterGroupDelete,
		Update: resourceKsyunKrdsParameterGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// parameter
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of krds parameter group.",
			},
			"db_parameter_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of krds parameter group.",
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The custom parameters.",
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKrdsEngine,
				Description:  "krds database type. Value options: mysql|percona|consistent_mysql|ebs_mysql.",
			},

			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "krds database version. Value options:<br> - Mysql: [ 5.5, 5.6, 5.7, 8.0 ] <br> - Percona: [ 5.6 ] <br> - Consistent_mysql: [ 5.7 ] <br> - Ebs_mysql: [ 5.6, 5.7 ].",
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

func resourceKsyunKrdsParameterGroupRead(d *schema.ResourceData, meta interface{}) error {
	krdsParameterSrv := NewKrdsParameterSrv(meta.(*KsyunClient))
	r := resourceKsyunKrdsParameterGroup()

	reqParameters := make(map[string]interface{})

	// call query function
	reqParameters["DBParameterGroupId"] = d.Id()
	action := "DescribeDBParameterGroupById"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdkResponse, err := krdsParameterSrv.describeDBParameterGroupById(reqParameters)
	if err != nil || len(sdkResponse) < 1 {
		return fmt.Errorf("while query db parameter group have encountered an error, detail: %s", err)
	}
	if err := TransformMapValue2StringWithKey("Parameters", sdkResponse); err != nil {
		return err
	}
	data := sdkResponse[0]

	extra := map[string]SdkResponseMapping{
		"DBParameterGroupName": {
			Field: "name",
		},
	}

	SdkResponseAutoResourceData(d, r, data, extra)

	return nil
}

func resourceKsyunKrdsParameterGroupCreate(d *schema.ResourceData, meta interface{}) error {
	krdsParameterSrv := NewKrdsParameterSrv(meta.(*KsyunClient))

	// TODO: to validate the relational mappings between engine and engine_version
	krdsEngine := d.Get("engine").(string)
	krdsEngineVersion := d.Get("engine_version").(string)

	if !validateKrdsEngineVersionWithEngine(krdsEngine, krdsEngineVersion) {
		return fmt.Errorf("the engine_version cannot match engine \n engine: %s need: %s", krdsEngine, KrdsEnginVersionMap[krdsEngine])
	}

	reqParameters, _, err := checkAndProcessKrdsParameters(d, meta)
	if err != nil {
		return err
	}

	reqParameters["DBParameterGroupName"] = d.Get("name")
	reqParameters["Engine"] = krdsEngine
	reqParameters["EngineVersion"] = krdsEngineVersion
	reqParameters["Description"] = d.Get("description")

	action := "CreateAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	dbParameterId, err := krdsParameterSrv.createDBParameterGroup(reqParameters)
	if err != nil {
		return err
	}

	d.SetId(dbParameterId)
	_ = d.Set("resource_name", ResourceKrdsParameterGroup)

	return d.Set("db_parameter_group_id", dbParameterId)
}

func resourceKsyunKrdsParameterGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	krdsParameterSrv := NewKrdsParameterSrv(meta.(*KsyunClient))

	reqParameters, _, err := checkAndProcessKrdsParameters(d, meta)
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

	if _, err := krdsParameterSrv.modifyDBParameterGroup(reqParameters); err != nil {
		return err
	}

	return resourceKsyunKrdsParameterGroupRead(d, meta)
}

func resourceKsyunKrdsParameterGroupDelete(d *schema.ResourceData, meta interface{}) error {
	krdsParameterSrv := NewKrdsParameterSrv(meta.(*KsyunClient))

	removeMap := map[string]interface{}{
		"DBParameterGroupId": d.Id(),
	}

	action := "DeleteDBParameterGroup"
	logger.Debug(logger.ReqFormat, action, removeMap)

	return krdsParameterSrv.deleteDBParameterGroup(removeMap)
}
