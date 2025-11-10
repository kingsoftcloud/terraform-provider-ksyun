/*
Provides a kce_auth_attachment resource.

Example Usage

```hcl
resource "ksyun_kce_auth_attachment" "auth" {
  sub_user_id = "38435"
  permissions {
	  cluster_id = "4cf5b24b-de39-4f55-b0ce-fd7b28cb964c"
	  cluster_role = "kce:dev"
	  namespace = ""
  }
}
```

Import

kce can be imported using the id, e.g.

```
$ terraform import ksyun_kce_auth_attachment.auth ${sub_user_id}
```

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKceAuthAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKceAuthAttachmentCreate,
		Read:   resourceKsyunKceAuthAttachmentRead,
		Update: resourceKsyunKceAuthAttachmentUpdate,
		Delete: resourceKsyunKceAuthAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sub_user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the id of the sub user.",
			},
			"permissions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "the permissions of the sub user.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the kce cluster.",
						},
						"cluster_role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "the role for the sub user in the cluster. Valid Values: kce:admin, kce:dev, kce:ops, kce:restricted, kce:ns:dev, kce:ns:restricted.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "the namespace of the cluster role, if it's empty, this authorization will apply in all of namespace.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "the region of the kce cluster.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the kce cluster.",
						},
					},
				},
			},
		},
	}
}

func resourceKsyunKceAuthAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kceService := KceService{meta.(*KsyunClient)}
	err = kceService.AddAuthorization(d, resourceKsyunKceAuthAttachment())
	if err != nil {
		return fmt.Errorf("error on creating ksyun_kce_auth_attachment %q, %s", d.Id(), err)
	}
	return resourceKsyunKceAuthAttachmentRead(d, meta)
}

func resourceKsyunKceAuthAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	kceService := KceService{meta.(*KsyunClient)}
	err = kceService.ReadAndSetKceAuthAttachment(d, resourceKsyunKceAuthAttachment())
	if err != nil {
		return fmt.Errorf("error on reading ksyun_kce_auth_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKceAuthAttachmentUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kceService := KceService{meta.(*KsyunClient)}
	err = kceService.ModifyAuthorization(d, resourceKsyunKceAuthAttachment())
	if err != nil {
		return fmt.Errorf("error on updating ksyun_kce_auth_attachment %q, %s", d.Id(), err)
	}
	return resourceKsyunKceAuthAttachmentRead(d, meta)
}

func resourceKsyunKceAuthAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kceService := KceService{meta.(*KsyunClient)}
	_ = d.Set("permissions", []interface{}{})
	err = kceService.ModifyAuthorization(d, resourceKsyunKceAuthAttachment())
	if err != nil {
		return fmt.Errorf("error on removing ksyun_kce_auth_attachment %q, %s", d.Id(), err)
	}
	return err
}
