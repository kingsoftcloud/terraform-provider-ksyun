/*
Provides an attachment for pinning tag upon resource.

> Note: supported all of resource_type
> The tag will be created if it is not existed.

# Example Usage

```hcl
resource "ksyun_tag_v2" "tagv2" {
  key   = "test_tag_key"
  value = "test_tag_value"
}


resource "ksyun_tag_v2_attachment" "tag" {
  key           = "test_tag_key"
  value         = "test_tag_value"
  resource_type = "redis-instance"
  resource_id   = "1f4e8c22-xxxx-xxxx-xxxx-cc6345011af4"
}

```

# Import

Tagv2Attachment can be imported using the `id`, e.g.

```
$ terraform import ksyun_tag_v2_attachment.tag ${tag_key}:${tag_value},${resource_type}:${resource_id}
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunTagv2Attachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunTagv2AttachmentCreate,
		Read:   resourceKsyunTagv2AttachmentRead,
		Update: resourceKsyunTagv2AttachmentUpdate,
		Delete: resourceKsyunTagv2AttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importTagV1Resource,
		},
		CustomizeDiff: resourceKsyunTagv2AttachmentDiff(),
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag key.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag value.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource type. [supported type](https://docs.ksyun.com/documents/43391).",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource ID.",
			},

			"tag_id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Tag id.",
			},
		},
	}
}

func resourceKsyunTagv2AttachmentDiff() schema.CustomizeDiffFunc {
	return func(diff *schema.ResourceDiff, i interface{}) (err error) {
		keys := []string{"key", "value", "resource_type", "resource_id"}

		for _, k := range keys {
			if diff.HasChange(k) {
				err = diff.ForceNew(k)
				if err != nil {
					return
				}
			}
		}

		return
	}
}

func resourceKsyunTagv2AttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagService{meta.(*KsyunClient)}
	err = tagService.CreateTagResourceAttachment(d, resourceKsyunTagv2Attachment())
	if err != nil {
		return fmt.Errorf("error on creating Tagv2Attachment %q, %s", d.Id(), err)
	}
	return resourceKsyunTagv2AttachmentRead(d, meta)
}

func resourceKsyunTagv2AttachmentUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunTagv2AttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagService{meta.(*KsyunClient)}
	err = tagService.ReadAndSetTagAttachment(d, resourceKsyunTagv2Attachment())
	if err != nil {
		return fmt.Errorf("error on reading Tagv2Attachment, %s", err)
	}
	return
}

func resourceKsyunTagv2AttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagService{meta.(*KsyunClient)}
	err = tagService.DetachResourceTags(d)
	if err != nil {
		return fmt.Errorf("error on deleting Tagv2Attachment %q, %s", d.Id(), err)
	}
	return
}
