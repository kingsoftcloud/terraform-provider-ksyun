/*
Provides a KCE attachment resource that attach a new instance to a cluster.

# Example Usage

```hcl
data "ksyun_kce_instance_images" "test" {
}

resource "ksyun_kce_cluster_attachment" "foo" {
  cluster_id = "67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx"

  worker_config {
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.2A"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "subnet-xxxxxx"
    security_group_id = ["sg-xxxxxx"]
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
$ terraform import ksyun_kce_cluster_attachment.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "worker_config.0.instance_password" {
						// Suppress the diff for instance_type as it is not modifiable
						return true
					}
					return false
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

	return nil
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
