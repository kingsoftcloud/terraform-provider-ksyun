/*
Provides an auto snapshot policy resource.

# Example Usage

```hcl

	resource "ksyun_auto_snapshot_policy" "foo" {
	  name   = "your auto snapshot policy name"
	  auto_snapshot_date = [1,3,4,5]
	  auto_snapshot_time = [1,3,4,5,9,22]
	}

```

# Import

`ksyun_auto_snapshot_policy` can be imported using the `id`, e.g.

```
$ terraform import ksyun_auto_snapshot_policy.foo "auto_snapshot_policy_id"
```
*/

package ksyun

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Read:   resourceKsyunAutoSnapshotPolicyRead,
		Create: resourceKsyunAutoSnapshotPolicyCreate,
		Delete: resourceKsyunAutoSnapshotPolicyDelete,
		Update: resourceKsyunAutoSnapshotPolicyUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// parameter
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of auto snapshot policy.",
			},
			"auto_snapshot_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of auto snapshot policy.",
			},
			"auto_snapshot_time": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(int)
						if v < 0 || v > 23 {
							errs = append(errs, fmt.Errorf("%q must be between 0 and 23 inclusive, got: %d", key, v))
						}
						return
					},
				},
				Required:    true,
				Description: "Setting the snapshot time in a day, its scope is between 0 and 23.",
			},
			"auto_snapshot_date": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(1, 7),
				},
				Required:    true,
				Description: "Setting the snapshot date in a week, its scope is between 1 and 7.",
			},

			"retention_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 9999),
				Description:  "the snapshot will be reserved for when, the cap is 9999.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The snapshot policy creation date.",
			},
		},
	}
}

func resourceKsyunAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

	r := resourceKsyunAutoSnapshotPolicy()

	reqParameters := make(map[string]interface{})

	if name, ok := d.GetOk("name"); ok {
		reqParameters["AutoSnapshotPolicyName"] = name
	}
	reqParameters["AutoSnapshotPolicyId.0"] = d.Id()
	// call query function
	action := "DescribeAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdkResponse, err := snapshotSrv.querySnapshotPolicyByID(reqParameters)
	if err != nil || len(sdkResponse) < 1 {
		return fmt.Errorf("while query snapshot policy have encountered an error detail: %s", err)
	}

	result := sdkResponse[0].(map[string]interface{})

	extra := map[string]SdkResponseMapping{
		"AutoSnapshotTime": {
			Field: "auto_snapshot_time",
			FieldRespFunc: func(values interface{}) interface{} {
				snapshotTimes := values.([]interface{})
				timeList := make([]int, 0, len(snapshotTimes))
				for _, v := range snapshotTimes {
					tv, _ := strconv.Atoi(v.(string))
					timeList = append(timeList, tv)
				}
				return timeList
			},
		},
		"AutoSnapshotDate": {
			Field: "auto_snapshot_date",
			FieldRespFunc: func(values interface{}) interface{} {

				snapshotDates := values.([]interface{})
				datesList := make([]int, 0, len(snapshotDates))
				for _, v := range snapshotDates {
					tv, _ := strconv.Atoi(v.(string))
					datesList = append(datesList, tv)
				}
				return datesList
			},
		},
		"AutoSnapshotPolicyName": {
			Field: "name",
		},
	}
	SdkResponseAutoResourceData(d, r, result, extra)

	return nil
}

func resourceKsyunAutoSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

	r := resourceKsyunAutoSnapshotPolicy()

	if err := checkTimesAndDatesEmpty(d); err != nil {
		return err
	}

	reqTransform := map[string]SdkReqTransform{
		"name":               {mapping: "AutoSnapshotPolicyName"},
		"auto_snapshot_time": {Type: TransformWithN},
		"auto_snapshot_date": {Type: TransformWithN},
		"retention_time":     {},
	}

	reqParameters, err := SdkRequestAutoMapping(d, r, false, reqTransform, nil)
	if err != nil {
		return err
	}

	action := "CreateAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	policyId, err := snapshotSrv.createAutoSnapshotPolicy(reqParameters)
	if err != nil {
		return err
	}

	_ = d.Set("auto_snapshot_policy_id", policyId)
	d.SetId(policyId)
	err = resourceKsyunAutoSnapshotPolicyRead(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func resourceKsyunAutoSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

	r := resourceKsyunAutoSnapshotPolicy()

	if err := checkTimesAndDatesEmpty(d); err != nil {
		return err
	}

	reqTransform := map[string]SdkReqTransform{
		"auto_snapshot_time": {Type: TransformWithN},
		"auto_snapshot_date": {Type: TransformWithN},

		"retention_time": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return d.GetOk("retention_time")
			},
			forceUpdateParam: true,
		},

		"name": {
			mapping:          "AutoSnapshotPolicyName",
			forceUpdateParam: true,
		},
		"auto_snapshot_policy_id": {
			forceUpdateParam: true,
		},
	}

	reqParameters, err := SdkRequestAutoMapping(d, r, true, reqTransform, nil)
	if err != nil {
		return err
	}

	action := "ModifyAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	if _, err := snapshotSrv.modifyAutoSnapshotPolicy(reqParameters); err != nil {
		return err
	}

	return resourceKsyunAutoSnapshotPolicyRead(d, meta)
}

func resourceKsyunAutoSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

	removeMap := map[string]interface{}{
		"AutoSnapshotPolicyId.1": d.Id(),
	}

	action := "DeleteAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, removeMap)

	_, err := snapshotSrv.deleteAutoSnapshotPolicy(removeMap)
	if err != nil {
		return err
	}

	return nil
}

func checkTimesAndDatesEmpty(d *schema.ResourceData) error {
	snapshotTimes := d.Get("auto_snapshot_time").(*schema.Set)
	snapshotDates := d.Get("auto_snapshot_date").(*schema.Set)
	if snapshotTimes.Len() < 1 {
		return fmt.Errorf("auto_snapshot_time is empty, please set your snpashot times")
	}
	if snapshotDates.Len() < 1 {
		return fmt.Errorf("auto_snapshot_date is empty, please set your snpashot times")
	}
	return nil
}
