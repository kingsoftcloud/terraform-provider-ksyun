/*
This data source provides a list of kec local snapshots in the current region.

# Example Usage

```hcl

	data "ksyun_local_snapshots" "default" {
	  output_file=""
	}

```
*/

package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func dataSourceKsyunLocalSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLocalSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"local_volume_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the volume.",
			},
			"source_local_volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the volume.",
			},
			"local_volume_snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the snapshot.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of local snapshots that satisfy the condition.",
			},
			"local_snapshot_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of KEC local snapshots. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_volume_snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_volume_snapshot_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_volume_snapshot_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_image": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"copy_from_remote": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"source_local_volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_local_volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_local_volume_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_local_volume_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"snapshot_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLocalSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	s := LocalVolumeService{meta.(*KsyunClient)}
	return s.ReadAndSetLocalSnapshots(d, dataSourceKsyunLocalSnapshots())
}
