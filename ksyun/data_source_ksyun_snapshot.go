/*
Provides a list of Redis security groups in the current region.

# Example Usage

```hcl

		data "ksyun_snapshot" "foo" {
			name = "your auto snapshot policy name"
			auto_snapshot_policy_ids = ["auto snapshot policy id"] // a list of auto snapshot policy id that can be null
			output_file = "output_result_snapshot"
		}

		output "ksyun_snapshot" {
			value = data.ksyun_snapshot.foo
		}

		output "ksyun_snapshots_total_count" {
			value = data.ksyun_snapshot.foo.total_count
		}
```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSnapshotRead,
		Schema: map[string]*schema.Schema{
			// parameter
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of KEC snapshot policy",
			},
			// query snapshot policy
			"auto_snapshot_policy_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The id of auto snapshot policy",
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

			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// return values by data source query
						"auto_snapshot_time": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The snapshot policy will be created in these hours",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy creation date",
						},
						"auto_snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy id",
						},
						"auto_snapshot_date": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The snapshot policy will be triggered in these dates per month",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"auto_snapshot_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy name",
						},
						"attach_local_volume_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The volume number that is attached to this policy",
						},
						"attach_ebs_volume_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The snapshot retention period (unit: day)",
						},
					},
				},
			},
		},
	}
}

// dataSourceKsyunSnapshotRead will read data source from ksyun
func dataSourceKsyunSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := SnapshotSrv{
		client: meta.(*KsyunClient),
	}
	r := dataSourceKsyunSnapshot()

	reqTransform := map[string]SdkReqTransform{
		"name":                     {mapping: "AutoSnapshotPolicyName"},
		"auto_snapshot_policy_ids": {mapping: "AutoSnapshotPolicyId", Type: TransformWithN},
	}

	reqParameters, err := mergeDataSourcesReq(d, r, reqTransform)
	if err != nil {
		return err
	}
	// call query function
	action := "DescribeAutoSnapshotPolicy"
	logger.Debug(logger.ReqFormat, action, reqParameters)

	sdKResponse, err := snapshotSrv.querySnapshotPolicyByID(reqParameters)
	if err != nil {
		return err
	}
	results, err := getSdkValue("AutoSnapshotPolicySet", sdKResponse)
	if err != nil {
		return err
	}
	data := results.([]interface{})

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "AutoSnapshotPolicyId",
		targetField: "snapshots",
		extra: map[string]SdkResponseMapping{
			"AttachEBSVolumeNum": {
				Field: "attach_ebs_volume_num",
			},
		},
	})
}
