/*
Provides a KCE worker resource.

# Example Usage

```hcl

resource "ksyun_kce_cluster_attachment" "foo" {
  cluster_id = ksyun_kce_cluster.default.id

  worker_config {
    image_id      = "fbafd8cd-b570-47c4-a3db-ff9702108f17"
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "c771027a-fafd-4b3b-a6b9-daeab9d0c13a"
    security_group_id = ["59a87036-dc27-41cf-98ab-24a387501195"]
    charge_type       = "Daily"

  }
  advanced_setting {
    container_runtime = "containerd"
    pre_user_script   = "def"
    label {
      key   = "tf_assembly_kce"
      value = "on_configuration_files"
    }
    taints {
      key    = "key2"
      value  = "value3"
      effect = "NoSchedule"
    }
  }
}

```

# Import

KCE worker can be imported using the id, e.g.

```
$ terraform import ksyun_kce_worker.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKceClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKceClusterAttachmentCreate,
		Update: resourceKsyunKceClusterAttachmentUpdate,
		Read:   resourceKsyunKceClusterAttachmentRead,
		Delete: resourceKsyunKceClusterAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the kce cluster.",
			},

			"instance_delete_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Terminate",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Description: "The instance delete mode when the instance is removed from the cluster. The value can be 'Terminate' or 'Remove'.",
			},

			"worker_config": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The instance node configuration for attach on cluster.",
				Elem: &schema.Resource{
					Schema: func() map[string]*schema.Schema {
						m := instanceForWorkerNode()
						delete(m, "count")
						delete(m, "advanced_setting")
						return m
					}(),
				},
			},

			"advanced_setting": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The advanced settings of the worker node.",
				Elem: &schema.Resource{
					Schema: nodeAdvancedSetting(),
				},
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the kec instance. The instance will be shut down while being added to the kce cluster.",
			},
		},
	}
}

func resourceKsyunKceClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KceWorkerService{
		meta.(*KsyunClient),
	}
	err = s.AddNewInstances(d, resourceKsyunKceClusterAttachment())
	return resourceKsyunKceClusterAttachmentRead(d, meta)
}
func resourceKsyunKceClusterAttachmentUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	// s := KceWorkerService{
	// 	meta.(*KsyunClient),
	// }
	// err = s.UpdateWorker(d, resourceKsyunKceClusterAttachment())

	return errors.New("you can't change anything at now. If you need to change, please move to Ksyun Console")
}
func resourceKsyunKceClusterAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceWorkerService{meta.(*KsyunClient)}
	err = srv.readAndSetAttachment(d, resourceKsyunKceClusterAttachment())
	if err != nil {
		return fmt.Errorf("error on create kce worker: %s", err)
	}
	return
}
func resourceKsyunKceClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceWorkerService{meta.(*KsyunClient)}
	err = srv.DeleteKceWorker(d, resourceKsyunKceClusterAttachment())
	if err != nil {
		return fmt.Errorf("error on delete kce cluster: %s", err)
	}
	return
}
