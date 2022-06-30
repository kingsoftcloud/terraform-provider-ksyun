package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunTagCreate,
		Read:   resourceKsyunTagRead,
		Update: resourceKsyunTagUpdate,
		Delete: resourceKsyunTagDelete,
		Importer: &schema.ResourceImporter{
			State: importTagV1Resource,
		},
		CustomizeDiff: resourceKsyunTagDiff(),
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKsyunTagDiff() schema.CustomizeDiffFunc {
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

func resourceKsyunTagCreate(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagV1Service{meta.(*KsyunClient)}
	err = tagService.CreateTag(d, resourceKsyunTag())
	if err != nil {
		return fmt.Errorf("error on creating tag %q, %s", d.Id(), err)
	}
	return resourceKsyunTagRead(d, meta)
}

func resourceKsyunTagUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	//tagService := TagV1Service{meta.(*KsyunClient)}
	//err = tagService.ModifyTag(d, resourceKsyunTag())
	//if err != nil {
	//	return fmt.Errorf("error on updating tag %q, %s", d.Id(), err)
	//}
	//return resourceKsyunTagRead(d, meta)
	return
}

func resourceKsyunTagRead(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagV1Service{meta.(*KsyunClient)}
	err = tagService.ReadAndSetTag(d, resourceKsyunTag())
	if err != nil {
		return fmt.Errorf("error on reading tag, %s", err)
	}
	return
}

func resourceKsyunTagDelete(d *schema.ResourceData, meta interface{}) (err error) {
	tagService := TagV1Service{meta.(*KsyunClient)}
	err = tagService.DeleteTag(d)
	if err != nil {
		return fmt.Errorf("error on deleting tag %q, %s", d.Id(), err)
	}
	return
}
