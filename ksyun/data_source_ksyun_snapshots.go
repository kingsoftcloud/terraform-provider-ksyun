/*
This data source provides a list of EBS snapshots.

# Example Usage

```hcl

	data "ksyun_snapshots" "default" {
	  output_file="output_result"
	}

```
*/

package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the volume.",
			},
			"volume_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "data"}, false),
				Description:  "The category of the volume.",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Id of the snapshot.",
			},
			"snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the snapshot.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "availability zone.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},
			"snapshots": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of EBS snapshot. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of the snapshot.",
						},
						"snapshot_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the snapshot.",
						},
						"volume_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Id of the volume.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "snapshot size, unit: GB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
						},
						"snapshot_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "snapshot status.",
						},
						"volume_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The category of the volume.",
						},
						"progress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Snapshot progress. Example value: 100%.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability zone.",
						},
						"volume_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Volume status.",
						},
						"snapshot_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "snapshot type.",
						},
						//"auto_snapshot": {
						//	Type:     schema.TypeBool,
						//	Computed: true,
						//},
						//"image_related": {
						//	Type:     schema.TypeBool,
						//	Computed: true,
						//},
						//"copy_from": {
						//	Type:     schema.TypeBool,
						//	Computed: true,
						//},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	s := SnapshotService{meta.(*KsyunClient)}
	return s.ReadAndSetSnapshots(d, dataSourceKsyunSnapshots())
}
