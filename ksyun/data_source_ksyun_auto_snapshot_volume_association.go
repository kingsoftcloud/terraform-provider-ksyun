/*
Provides a list of Redis security groups in the current region.

# Example Usage

query asp and volume associations by volume id

```hcl
	data "ksyun_auto_snapshot_volume_association" "foo" {
		output_file = "output_result_volume_id"
		attach_volume_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}
```

query all associations.

```hcl
	data "ksyun_auto_snapshot_volume_association" "foo1" {
		output_file = "output_result_null"
	}
```

query asp and volume associations by auto_snapshot_policy_id
```hcl
	data "ksyun_auto_snapshot_volume_association" "foo2" {
		output_file = "output_result_policy_id"
		auto_snapshot_policy_id = "auto_snapshot_policy_id"
	}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunAutoSnapshotVolumeAssociation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAutoSnapshotVolumeAssociationRead,
		Schema: map[string]*schema.Schema{
			// parameter
			"attach_volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the volume.",
			},
			"auto_snapshot_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the auto snapshot policy.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of snapshot policies resources that satisfy the condition.",
			},

			"volume_asp_associations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of the associations of volumes and auto snapshot policy. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// return values by data source query
						"attach_volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the volume.",
						},

						"auto_snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy id.",
						},
					},
				},
			},
		},
	}
}

// dataSourceKsyunAutoSnapshotVolumeAssociationRead will read data source from ksyun
func dataSourceKsyunAutoSnapshotVolumeAssociationRead(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := NewAutoSnapshotSrv(meta.(*KsyunClient))

	var (
		sdkResponse []interface{}
		err         error
		policyId    string
	)

	r := dataSourceKsyunAutoSnapshotVolumeAssociation()

	if valIf, ok := d.GetOk("attach_volume_id"); ok {
		val := valIf.(string)
		sdkResponse, err = snapshotSrv.readAutoSnapshotPolicyVolumeAssociationById(val)
	} else {
		sdkResponse, err = snapshotSrv.readAutoSnapshotPolicyVolumeAssociationAll()
	}

	if err != nil {
		return err
	}

	if valIf, ok := d.GetOk("auto_snapshot_policy_id"); ok {
		policyId = valIf.(string)
	}

	associations, err := snapshotSrv.filterVolumesAutoSnapshotPolicyAssociations(sdkResponse, policyId)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  associations,
		idFiled:     "VolumeId",
		targetField: "volume_asp_associations",
		extra: map[string]SdkResponseMapping{
			"VolumeId": {
				Field: "attach_volume_id",
			},
		},
	})
}
