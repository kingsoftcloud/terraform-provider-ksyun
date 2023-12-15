/*
Provides a Tagv2 resource.

# Example Usage

```hcl

resource "ksyun_tag_v2" "tag" {
  key = "test_tag_key"
  value = "test_tag_value"
}

```

# Import

Tagv2 can be imported using the `key&value`, e.g.

```
$ terraform import ksyun_tag_v2.tag ${tagv2_key}:${tagv2_value}
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunTagv2() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunTagv2Create,
		Read:   resourceKsyunTagv2Read,
		Update: resourceKsyunTagv2Update,
		Delete: resourceKsyunTagv2Delete,
		Importer: &schema.ResourceImporter{
			State: importTagResource,
		},
		CustomizeDiff: resourceKsyunTagv2Diff(),
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
		},
	}
}

func resourceKsyunTagv2Diff() schema.CustomizeDiffFunc {
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

func resourceKsyunTagv2Create(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagService{meta.(*KsyunClient)}
	err = tagService.CreateTag(d, resourceKsyunTagv2())
	if err != nil {
		return fmt.Errorf("error on creating Tagv2 %q, %s", d.Id(), err)
	}
	return resourceKsyunTagv2Read(d, meta)
}

func resourceKsyunTagv2Update(d *schema.ResourceData, meta interface{}) (err error) {
	return
}

func resourceKsyunTagv2Read(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagService{meta.(*KsyunClient)}
	err = tagService.ReadAndSetTag(d, resourceKsyunTagv2())
	if err != nil {
		return fmt.Errorf("error on reading Tagv2, %s", err)
	}
	return
}

func resourceKsyunTagv2Delete(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagService{meta.(*KsyunClient)}
	err = tagService.DeleteTag(d)
	if err != nil {
		return fmt.Errorf("error on deleting Tagv2 %q, %s", d.Id(), err)
	}
	return
}
