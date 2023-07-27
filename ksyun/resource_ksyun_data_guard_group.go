/*
Provides a data guard group resource.

# Example Usage

```hcl

	resource "ksyun_data_guard_group" "foo" {
	  data_guard_name = "your data guard name"
	  data_guard_type = "host"
	}
```

# Import

Data guard group can be imported using the `id`, e.g.

```
$ terraform import ksyun_data_guard_group.foo "data_guard_id"
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunDataGuardGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceKsyunDataGuardGroupRead,
		Create: resourceKsyunDataGuardGroupCreate,
		Delete: resourceKsyunDataGuardGroupDelete,
		Update: resourceKsyunDataGuardGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// parameter
			"data_guard_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of data guard group.",
			},
			// query data guard
			"data_guard_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of data guard group.",
			},
			"data_guard_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "host" && v != "domain" {
						errs = append(errs, fmt.Errorf("%q must be the host and domain type, got: %s", key, v))
					}
					return
				},
				Description: "The data guard group display type, Value options: [ host, domain ]. <br> - Host: represent machine level <br> - Domain: represent the domain of disaster tolerance <br> Note: some account doesn't support domain data guard.",
			},
			"data_guard_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data guard group level, if the value is Host represent machine level, and the tol represent the domain of disaster tolerance.",
			},
			"data_guard_capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity of data guard group.",
			},
			"data_guard_used_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This data guard group includes the amount of instances.",
			},
		},
	}
}

func resourceKsyunDataGuardGroupRead(d *schema.ResourceData, meta interface{}) error {
	dataGuardSrv := NewDataGuardSrv(meta.(*KsyunClient))
	r := resourceKsyunDataGuardGroup()

	reqParameters := map[string]interface{}{}

	if dataGuardName, ok := d.GetOk("data_guard_name"); ok {
		reqParameters["DataGuardName"] = dataGuardName
	}

	reqParameters["DataGuardId"] = d.Id()
	// call query function
	action := "DescribeDataGuardGroup"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdkResponse, err := dataGuardSrv.describeDataGuardGroup(reqParameters)
	if err != nil || len(sdkResponse) < 1 {
		return fmt.Errorf("while query snapshot policy have encountered an error detail: %s", err)
	}

	if len(sdkResponse) < 1 {
		return fmt.Errorf("the data guard group doesn't exsit from ksyun. data_guard_id: %s, data_guard_name: %s "+
			"\n This resource has been deleted on ksyun. you can delete local resource in ./.terraform.tfstate",
			d.Get("data_guard_id"), d.Get("data_guard_name"))
	}

	result := sdkResponse[0].(map[string]interface{})

	SdkResponseAutoResourceData(d, r, result, nil)

	return nil
}

func resourceKsyunDataGuardGroupCreate(d *schema.ResourceData, meta interface{}) error {
	dataGuardSrv := NewDataGuardSrv(meta.(*KsyunClient))

	r := resourceKsyunDataGuardGroup()

	reqTransform := map[string]SdkReqTransform{
		"data_guard_name": {},
		"data_guard_type": {},
	}

	reqParameters, err := SdkRequestAutoMapping(d, r, false, reqTransform, nil)
	if err != nil {
		return err
	}

	action := "CreateDataGuardGroup"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	guardId, err := dataGuardSrv.createDataGuardGroup(reqParameters)
	if err != nil {
		return err
	}

	_ = d.Set("data_guard_id", guardId)
	d.SetId(guardId)

	return resourceKsyunDataGuardGroupRead(d, meta)
}

func resourceKsyunDataGuardGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	dataGuardSrv := DataGuardSrv{
		client: meta.(*KsyunClient),
	}

	if ok := d.HasChange("data_guard_name"); !ok {
		return fmt.Errorf("this resource only supports to modify data_guard_name")
	}
	reqParameters := map[string]interface{}{
		"DataGuardName": d.Get("data_guard_name"),
	}

	reqParameters["DataGuardId"] = d.Id()
	action := "ModifyModifyDataGuardGroups"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	if _, err := dataGuardSrv.modifyModifyDataGuardGroups(reqParameters); err != nil {
		return err
	}

	return resourceKsyunDataGuardGroupRead(d, meta)
}

func resourceKsyunDataGuardGroupDelete(d *schema.ResourceData, meta interface{}) error {
	dataGuardSrv := NewDataGuardSrv(meta.(*KsyunClient))

	removeMap := map[string]interface{}{
		"DataGuardId.1": d.Id(),
	}

	action := "DeleteDataGuardGroups"
	logger.Debug(logger.ReqFormat, action, removeMap)

	return dataGuardSrv.deleteDataGuardGroup(removeMap)
}
