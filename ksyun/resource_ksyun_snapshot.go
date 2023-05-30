// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ksyun

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunSnapshot() *schema.Resource {
	return &schema.Resource{
		Read:   resourceKsyunSnapshotRead,
		Create: resourceKsyunSnapshotCreate,
		Delete: resourceKsyunSnapshotDelete,
		Update: resourceKsyunSnapshotUpdate,

		Schema: map[string]*schema.Schema{
			// parameter
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of KEC snapshot policy",
			},
			"auto_snapshot_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of auto snapshot policy",
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
				Required: true,
				// Default:     []int{0},
				Description: "Setting the snapshot time in a day, its scope is between 0 and 23",
			},
			"auto_snapshot_date": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(1, 7),
				},
				Required:    true,
				Description: "Setting the snapshot date in a week, its scope is between 1 and 7",
			},

			"retention_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 9999),
				Description:  "the snapshot will be reserved for when, the cap is 9999",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The snapshot policy creation date",
			},
		},
	}
}

func resourceKsyunSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}
	r := resourceKsyunSnapshot()

	reqTransform := map[string]SdkReqTransform{
		"name":                    {mapping: "AutoSnapshotPolicyName"},
		"auto_snapshot_policy_id": {mapping: "AutoSnapshotPolicyId", Type: TransformWithN},
	}

	reqParameters, err := SdkRequestAutoMapping(d, r, false, reqTransform, nil)
	if err != nil {
		return err
	}
	// call query function
	action := "DescribeAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdkResponse, err := snapshotSrv.querySnapshotPolicyByID(reqParameters)
	if err != nil || len(sdkResponse) < 1 {
		return fmt.Errorf("while query snapshot policy have encountered an error detail: %s", err)
	}
	policySet, err := getSdkValue("AutoSnapshotPolicySet", sdkResponse)
	if err != nil || policySet == nil {
		return fmt.Errorf("the snapshot doesn't exsit from ksyun sdk. auto_snapshot_policy_id: %s, name: %s "+
			"\n This resource has been deleted on ksyun. you can delete local resource in ./.terraform.tfstate",
			d.Get("auto_snapshot_policy_id"), d.Get("name"))
	}

	data := policySet.([]interface{})

	result := data[0].(map[string]interface{})

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
	}
	SdkResponseAutoResourceData(d, r, result, extra)

	return nil
}

func resourceKsyunSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}
	r := resourceKsyunSnapshot()

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

	sdkResponse, err := snapshotSrv.createAutoSnapshotPolicy(reqParameters)
	if err != nil {
		return err
	}
	results, err := getSdkValue("AutoSnapshotPolicyId", sdkResponse)
	if err != nil {
		return err
	}
	policyId := results.(string)
	_ = d.Set("auto_snapshot_policy_id", policyId)
	d.SetId(policyId)
	err = resourceKsyunSnapshotRead(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func resourceKsyunSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}
	r := resourceKsyunSnapshot()

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

	return resourceKsyunSnapshotRead(d, meta)
}

func resourceKsyunSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}

	removeMap := map[string]interface{}{
		"AutoSnapshotPolicyId.1": d.Id(),
	}

	action := "DeleteAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, removeMap)

	resp, err := snapshotSrv.deleteAutoSnapshotPolicy(removeMap)
	if err != nil {
		return err
	}

	policySet, err := getSdkValue("AutoSnapshotPolicySet", resp)
	if err != nil || policySet == nil {
		return fmt.Errorf("fail to delete snapshot from ksyun sdk. auto_snapshot_policy_id: %s", d.Get("auto_snapshot_policy_id"))
	}

	return nil
}
