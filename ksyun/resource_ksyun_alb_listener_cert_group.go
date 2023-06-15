/*
Provides a ALB Listener cert group resource.

# Example Usage

```hcl

	resource "ksyun_alb_listener_cert_group" "default" {
	}

```

# Import

ALB Listener Cert Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener_cert_group.example vserver-abcdefg
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunAlbListenerCertGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAlbListenerCertGroupCreate,
		Read:   resourceKsyunAlbListenerCertGroupRead,
		Update: resourceKsyunAlbListenerCertGroupUpdate,
		Delete: resourceKsyunAlbListenerCertGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alb_listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the ALB Listener.",
			},
			"alb_listener_cert_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the ALB Listener Cert Group.",
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				//Computed:    true,
				Description: "The certificate included in the cert group.",
				//DiffSuppressFunc: func(k, oldV, newV string, d *schema.ResourceData) bool {
				//	logger.Debug(logger.ReqFormat, "certificateDiffSuppressFunc k", k)
				//	logger.Debug(logger.ReqFormat, "certificateDiffSuppressFunc old", oldV)
				//	logger.Debug(logger.ReqFormat, "certificateDiffSuppressFunc new", newV)
				//	if k == "certificate.#" {
				//		if oldV == "0" && newV == "0" {
				//			return true
				//		}
				//	}
				//	return false
				//},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the certificate.",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the certificate.",
						},
						"cert_authority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate authority.",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The common name on the certificate.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expire time of the certificate.",
						},
					},
				},
			},
		},
	}
}

func resourceKsyunAlbListenerCertGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.CreateCertGroup(d, resourceKsyunAlbListenerCertGroup())
	if err != nil {
		return fmt.Errorf("error on creating ALB listener cert group %q, %s", d.Id(), err)
	}
	err = resourceKsyunAlbListenerCertGroupRead(d, meta)
	return
}
func resourceKsyunAlbListenerCertGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.ReadAndSetCertGroup(d, resourceKsyunAlbListenerCertGroup())
	if err != nil {
		return fmt.Errorf("error on reading ALB listener cert group %q, %s", d.Id(), err)
	}
	return
}
func resourceKsyunAlbListenerCertGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.ModifyCertGroup(d, resourceKsyunAlbListenerCertGroup())
	if err != nil {
		return fmt.Errorf("error on updating listener cert group %q, %s", d.Id(), err)
	}
	err = resourceKsyunAlbListenerCertGroupRead(d, meta)
	return
}
func resourceKsyunAlbListenerCertGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.RemoveCertGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener cert group %q, %s", d.Id(), err)
	}
	return
}
