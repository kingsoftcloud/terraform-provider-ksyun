/*
Provides a kpfs acl rule resource.

# Example Usage

```hcl

	resource "ksyun_kpfs_acl" "default" {
	  epc_id = "c6c683f8-5bb4-4747-8516-9a61f01c4bce"
	  kpfs_acl_id = "4a42284d7f354e6a9b2d7a3454e0b495"
	}

```

# Import

KPFS ACL rules can be imported using the id, e.g.

```
$ terraform import ksyun_kpfs_acl.example ${epc_id}_${kpfs_acl_id}
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunKpfsAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunPerformanceOnePosixAclCreate,
		Read:   resourceKsyunPerformanceOnePosixAclRead,
		Delete: resourceKsyunPerformanceOnePosixAclDelete,
		Importer: &schema.ResourceImporter{
			State: importKpfsAcl,
		},
		Schema: map[string]*schema.Schema{
			"epc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The epc instance id.",
			},
			"kpfs_acl_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The posix acl rule id.",
			},
		},
	}
}

func resourceKsyunPerformanceOnePosixAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	err = kpfsService.readPerformanceOnePosixAcl(d, resourceKsyunKpfsAcl())
	if err != nil {
		return fmt.Errorf("error on read posix acl %s", err)
	}
	logger.Debug(logger.RespFormat, "readPerformanceOnePosixAcl", d)
	return err
}

func resourceKsyunPerformanceOnePosixAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	r := resourceKsyunKpfsAcl()
	err = kpfsService.updatePerformanceOnePosixAcl(d, r)
	if err != nil {
		return fmt.Errorf("error on add posix acl %s", err)
	}
	err = resourceKsyunPerformanceOnePosixAclRead(d, meta)
	transform := map[string]SdkReqTransform{
		"epc_id":      {mapping: "epcId"},
		"kpfs_acl_id": {mapping: "aclId"},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return fmt.Errorf("error on add posix acl %s", err)
	}
	d.SetId(fmt.Sprintf("%s-%s", req["epcId"], req["aclId"]))
	return err
}

func resourceKsyunPerformanceOnePosixAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	err = kpfsService.deletePerformanceOnePosixAclIp(d, resourceKsyunKpfsAcl())
	if err != nil {
		return fmt.Errorf("error on delete posix acl %s", err)
	}
	return
}
