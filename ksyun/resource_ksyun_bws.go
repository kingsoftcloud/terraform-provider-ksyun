/*
Provides a BandWidthShare resource.

Example Usage

```hcl
resource "ksyun_bws" "default" {
  line_id = "5fc2595f-1bfd-481b-bf64-2d08f116d800"
  charge_type = "PostPaidByPeak"
  band_width = 12
}
```

Import

BWS can be imported using the id, e.g.

```
$ terraform import ksyun_bws.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```

*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunBandWidthShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBandWidthShareCreate,
		Read:   resourceKsyunBandWidthShareRead,
		Update: resourceKsyunBandWidthShareUpdate,
		Delete: resourceKsyunBandWidthShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"line_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the line.",
			},
			"band_width_share_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "name of the BWS.",
			},
			"band_width": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 15000),
				Description:  "bandwidth value, value range: [1, 15000].",
			},
			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				// ValidateFunc: validation.StringInSlice([]string{
				// 	"PostPaidByPeak",
				// 	"PostPaidByDay",
				// 	"PostPaidByTransfer",
				// }, false),
				DiffSuppressFunc: chargeSchemaDiffSuppressFunc,
				Description:      "The charge type of the BWS. Valid values: PostPaidByPeak, PostPaidByDay, PostPaidByTransfer, DailyPaidByTransfer.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     0,
				Description: "ID of the project.",
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceKsyunBandWidthShareCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.CreateBandWidthShare(d, resourceKsyunBandWidthShare())
	if err != nil {
		return fmt.Errorf("error on creating bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunBandWidthShareRead(d, meta)
}

func resourceKsyunBandWidthShareRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetBandWidthShare(d, resourceKsyunBandWidthShare())
	if err != nil {
		return fmt.Errorf("error on reading bandWidthShare %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBandWidthShareUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ModifyBandWidthShare(d, resourceKsyunBandWidthShare())
	if err != nil {
		return fmt.Errorf("error on updating bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunBandWidthShareRead(d, meta)
}

func resourceKsyunBandWidthShareDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.RemoveBandWidthShare(d)
	if err != nil {
		return fmt.Errorf("error on deleting bandWidthShare %q, %s", d.Id(), err)
	}
	return err
}
