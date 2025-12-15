/*
This data source provides a list of kpfs fileSystem resources according to their fileSystem ID, name.

# Example Usage

```hcl

	data "ksyun_kpfs_file_systems" "default" {
	  output_file="output_result"
	  id="fileSystemId"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKpfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKpfsFileSystemsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File system ID.",
			},
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page number for pagination.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page size for pagination.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of matching file systems.",
			},
			"data": {
				Description: "It is a nested type which documented below.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current state: creating/using/upgrading/renewing/shutdown.",
						},
						"charge_info": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "File system billing information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "monthly (prepaid) or dailySettlement (postpaid).",
									},
									"expired_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expiration time for prepaid (e.g., 2024-12-04T23:59:59).",
									},
									"is_trial": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether it is a trial file system.",
									},
								},
							},
						},
						"file_system_info": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_system_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System name (e.g., train-fs).",
									},
									"capacity": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Allocated capacity in TiB.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region code (e.g., cn-beijing-6).",
									},
									"region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region description (e.g., 华北1（北京）).",
									},
									"avail_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone (e.g., cn-beijing-6e).",
									},
									"file_system_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique identifier.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation timestamp.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modification time.",
									},
									"store_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Storage type: KPFS-capacity/KPFS-capacity2/KPFS-standard/KPFS-P-S01/KPFS-P-S02.",
									},
									"client_mount_command": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Complete mount instructions (Capacity/Standard only).",
									},
									"check_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Check size.",
									},
									"check_size_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Check size type.",
									},
									"throughput_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ThroughputLimit.",
									},
								},
							},
						},
						"volume_info": {
							Description: "It is a nested type which documented below.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inodes": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of files stored.",
									},
									"use_capacity": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Actual used space in bytes.",
									},
								},
							},
						},
					},
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
		},
	}
}

func dataSourceKsyunKpfsFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	return kpfsService.ReadKpfsFileSystemList(d, dataSourceKsyunKpfsFileSystems())
}
