/*
Query instance auto snapshot policies information

# Example Usage

## query auto snapshot policy with name or id
```hcl

		data "ksyun_auto_snapshot_policy" "foo" {
			name = "your auto snapshot policy name"
			auto_snapshot_policy_ids = ["auto snapshot policy id"]
			output_file = "output_result_snapshot"
		}

		output "ksyun_auto_snapshot_policy" {
			value = data.ksyun_auto_snapshot_policy.foo
		}

```

## query all auto snapshot policy

```hcl
		data "ksyun_auto_snapshot_policy" "foo" {
			output_file = "output_result_snapshot"
		}
```

*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func dataSourceKsyunAutoSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunAutoSnapshotPolicyRead,
		Schema: map[string]*schema.Schema{
			// parameter
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the name of auto snapshot policy.",
			},
			// query snapshot policy
			"auto_snapshot_policy_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The id of auto snapshot policy.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of auto snapshot policies resources that satisfy the condition.",
			},

			"snapshots": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of auto snapshot policy. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// return values by data source query
						"auto_snapshot_time": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The snapshot policy will be created in these hours.",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy creation date.",
						},
						"auto_snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy id.",
						},
						"auto_snapshot_date": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The snapshot policy will be triggered in these dates per month.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"auto_snapshot_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The snapshot policy name.",
						},
						"attach_local_volume_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The volume number that is attached to this policy.",
						},
						"attach_ebs_volume_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The snapshot retention period (unit: day).",
						},
					},
				},
			},
		},
	}
}

// dataSourceKsyunAutoSnapshotPolicyRead will read data source from ksyun
func dataSourceKsyunAutoSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	snapshotSrv := AutoSnapshotSrv{
		client: meta.(*KsyunClient),
	}
	r := dataSourceKsyunAutoSnapshotPolicy()

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

	data, err := snapshotSrv.querySnapshotPolicyByID(reqParameters)
	if err != nil {
		return err
	}

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
