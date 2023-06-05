/*
Provides a tag resource.

# Example Usage

```hcl

	resource "ksyun_snapshot_volume_association" "foo" {
	  attach_volume_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  auto_snapshot_policy_id = "auto_snapshot_policy_id"
	}
```
*/

package ksyun

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunAutoSnapshotVolumeAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAutoSnapshotVolumeAssociationCreate,
		Read:   resourceKsyunAutoSnapshotVolumeAssociationRead,
		Delete: resourceKsyunAutoSnapshotVolumeAssociationDelete,
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

func resourceKsyunAutoSnapshotVolumeAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	SnapshotSrv := AutoSnapshotSrv{
		client: meta.(*KsyunClient),
	}

	r := resourceKsyunAutoSnapshotVolumeAssociation()

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

func resourceKsyunAutoSnapshotVolumeAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceKsyunAutoSnapshotVolumeAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	SnapshotSrv := AutoSnapshotSrv{
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
