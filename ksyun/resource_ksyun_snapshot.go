/*
Provides a EBS snapshot resource.

# Example Usage

```hcl

		resource "ksyun_snapshot" "default" {
	  		snapshot_name = "test_tf_snapshot"
			snapshot_desc = "test descrition"
			volume_id = "xxxxxxxxx"
		}

```

# Import

Instance can be imported using the `id`, e.g.

```
$ terraform import ksyun_snapshot.default xxxxxx
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceKsyunSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSnapshotCreate,
		Update: resourceKsyunSnapshotUpdate,
		Read:   resourceKsyunSnapshotRead,
		Delete: resourceKsyunSnapShotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the volume. Snapshot requires the Volume to be in \"in-use\" or \"available\" status.When the Volume status is \"in-use\", the kec instance status can be either \"running\" or \"stopped\".",
			},
			"snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the snapshot.",
			},
			"snapshot_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the snapshot.",
			},
			"snapshot_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The type of the snapshot, valid values: 'LocalSnapShot', 'CommonSnapShot'. Default is 'CommonSnapShot'.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Availability zone.",
			},
			"volume_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The category of the volume, 'data' or 'system'.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the snapshot, unit is 'GB'.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"snapshot_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "snapshot status.",
			},
			"progress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Snapshot progress. Example value: 100%.",
			},
			"volume_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume status.",
			},
		},
	}
}

func resourceKsyunSnapshotCreate(d *schema.ResourceData, meta interface{}) (err error) {
	snapshotService := SnapshotService{meta.(*KsyunClient)}
	err = snapshotService.CreateSnapshot(d, resourceKsyunSnapshot())
	if err != nil {
		return fmt.Errorf("error on creating stnapshot %q, %s", d.Id(), err)
	}
	return err
}
func resourceKsyunSnapshotUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	snapshotService := SnapshotService{meta.(*KsyunClient)}
	err = snapshotService.ModifySnapshot(d, resourceKsyunSnapshot())
	if err != nil {
		return
	}
	err = resourceKsyunSnapshotRead(d, meta)
	return
}
func resourceKsyunSnapshotRead(d *schema.ResourceData, meta interface{}) (err error) {
	snapshotService := SnapshotService{meta.(*KsyunClient)}
	err = snapshotService.ReadAndSetSnapshot(d, resourceKsyunSnapshot())
	if err != nil {
		return fmt.Errorf("error on reading snapshot %q, %s", d.Id(), err)
	}
	return err
}
func resourceKsyunSnapShotDelete(d *schema.ResourceData, meta interface{}) (err error) {
	snapshotService := SnapshotService{meta.(*KsyunClient)}
	err = snapshotService.Remove(d)
	return err
}
