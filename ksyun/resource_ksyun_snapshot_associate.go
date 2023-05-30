// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ksyun

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func ResourceKsyunSnapshotAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSnapshotAssociationCreate,
		Read:   resourceKsyunSnapshotAssociationRead,
		Delete: resourceKsyunSnapshotAssociationDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: importNatAssociate,
		// },

		Schema: map[string]*schema.Schema{
			"attach_volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the volume.",
			},
			"auto_snapshot_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the auto_snapshot_policy_id.",
			},
		},
	}
}

func resourceKsyunSnapshotAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	SnapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}

	r := resourceKsyunNatAssociation()

	reqTransform := map[string]SdkReqTransform{
		"attach_volume_id":        {Type: TransformWithN},
		"auto_snapshot_policy_id": {},
	}

	reqParameters, err := SdkRequestAutoMapping(d, r, false, reqTransform, nil)
	if err != nil {
		return err
	}

	action := "ApplyAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)
	sdkResponse, err := SnapshotSrv.associatedAutoSnapshotPolicy(reqParameters)
	if err != nil {
		return err
	}
	associationSet, err := getSdkValue("ReturnSet", sdkResponse)
	if err != nil || associationSet == nil {
		return fmt.Errorf("the snapshot policy is fail to associate attach volume, \n auto_snapshot_policy_id: %s \n attach_volume_id: %s", d.Get("auto_snapshot_policy_id"), d.Get("attach_volume_id"))
	}

	combineIds := []string{d.Get("auto_snapshot_policy_id").(string), d.Get("attach_volume_id").(string)}
	d.SetId(strings.Join(combineIds, ":"))
	return nil
}

func resourceKsyunSnapshotAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKsyunSnapshotAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	SnapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}

	r := resourceKsyunNatAssociation()

	reqTransform := map[string]SdkReqTransform{
		"attach_volume_id": {
			Type: TransformWithN,
		},
		"auto_snapshot_policy_id": {},
	}

	reqParameters, err := SdkRequestAutoMapping(d, r, false, reqTransform, nil)
	if err != nil {
		return err
	}

	action := "CancelAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdkResponse, err := SnapshotSrv.unassociatedAutoSnapshotPolicy(reqParameters)
	if err != nil {
		return err
	}
	associationSet, err := getSdkValue("ReturnSet", sdkResponse)
	if err != nil || associationSet == nil {
		return fmt.Errorf("the snapshot policy is fail to unassociate attach volume, \n auto_snapshot_policy_id: %s \n attach_volume_id: %s", d.Get("auto_snapshot_policy_id"), d.Get("attach_volume_id"))
	}
	return nil
}
