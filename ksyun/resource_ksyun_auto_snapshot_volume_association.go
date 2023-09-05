/*
Provides an auto snapshot policy associate to volume.

# Example Usage

```hcl

	resource "ksyun_auto_snapshot_volume_association" "foo" {
	  attach_volume_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  auto_snapshot_policy_id = "auto_snapshot_policy_id"
	}
```

# Import

ksyun_auto_snapshot_volume_association can be imported using the `id`, e.g.

```
$ terraform import ksyun_auto_snapshot_volume_association.foo ${auto_snapshot_policy_id}:${attach_volume_id}
```
*/

package ksyun

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunAutoSnapshotVolumeAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAutoSnapshotVolumeAssociationCreate,
		Read:   resourceKsyunAutoSnapshotVolumeAssociationRead,
		Delete: resourceKsyunAutoSnapshotVolumeAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: importAutoSnapshotPolicyAssociate,
		},

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
	SnapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

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
	associationSet, err := SnapshotSrv.associatedAutoSnapshotPolicy(reqParameters)
	if err != nil {
		return err
	}
	logger.Debug(logger.RespFormat, action, reqParameters, associationSet)

	combineIds := []string{d.Get("auto_snapshot_policy_id").(string), d.Get("attach_volume_id").(string)}
	d.SetId(strings.Join(combineIds, ":"))
	return nil
}

func resourceKsyunAutoSnapshotVolumeAssociationRead(d *schema.ResourceData, meta interface{}) error {
	// DescribeVolumes query ebs volumes
	SnapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

	r := resourceKsyunNatAssociation()

	volumeId := d.Get("attach_volume_id").(string)

	sdkResponse, err := SnapshotSrv.readAutoSnapshotPolicyVolumeAssociationById(volumeId)
	if err != nil {
		return err
	}

	if sdkResponse != nil && len(sdkResponse) > 0 {
		SdkResponseAutoResourceData(d, r, sdkResponse[0], map[string]SdkResponseMapping{
			"VolumeId": {
				Field: "attach_volume_id",
			},
		})
	}

	return nil
}

func resourceKsyunAutoSnapshotVolumeAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	SnapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

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
	sdkResponse, err := SnapshotSrv.unassociatedAutoSnapshotPolicy(reqParameters)
	if err != nil {
		return err
	}

	logger.Debug(logger.RespFormat, action, reqParameters, sdkResponse)

	return nil
}
