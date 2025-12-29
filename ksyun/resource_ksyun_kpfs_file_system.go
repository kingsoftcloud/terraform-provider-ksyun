/*
Create a kpfs file system.

# Example Usage

```hcl

	resource "ksyun_kpfs_file_system" "example" {
	  file_system_name = "examplefs9"
	  region           = "cn-qingyangtest-1"
	  avail_zone       = "cn-qingyangtest-1a"
	  charge_type      = "dailySettlement"
	  project_id       = "0"
	  purchase_time    = "-1"
	  store_class      = "KPFS-P-S01"
	  capacity         = 102
	  chunk_size_type  = "Balanced"
	  cluster_code     = "7c2cd53b-8eec-440d-a450-18b7db2b0040-5"
	}

```

# Import

KPFS filesystem can be imported using the id, e.g.

```
$ terraform import ksyun_kpfs_file_system.example ${filesystem_id}
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunKpfsFilesystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKpfsFileSystemCreate,
		Read:   resourceKsyunKpfsFileSystemRead,
		Update: resourceKsyunKpfsFileSystemUpdate,
		Delete: resourceKsyunKpfsFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: importKpfsFileSystem,
		},
		Schema: map[string]*schema.Schema{
			"file_system_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the file system.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region of the file system.",
			},
			"avail_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The availability zone of the file system.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The charge type of the file system.",
			},
			"purchase_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The purchase time of the file system.",
			},
			"store_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The store class of the file system.",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The capacity of the file system.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The project id of the file system.",
			},
			"chunk_size_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The chunk size type of the file system.",
			},
			"cluster_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster code of the file system.",
			},
		},
	}
}

func resourceKsyunKpfsFileSystemRead(d *schema.ResourceData, meta interface{}) (err error) {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	err = kpfsService.ReadKpfsFileSystemOne(d, resourceKsyunKpfsFilesystem())
	if err != nil {
		return fmt.Errorf("error on read file system %s", err)
	}
	logger.Debug(logger.RespFormat, "readKpfsFile", d)
	return err
}

func resourceKsyunKpfsFileSystemCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	err = kpfsService.CreateFileSystem(d, resourceKsyunKpfsFilesystem())
	if err != nil {
		return fmt.Errorf("error on creating filesystem %s", err)
	}
	return err
}

func resourceKsyunKpfsFileSystemDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	err = kpfsService.DeleteFileSystem(d, meta)
	if err != nil {
		return fmt.Errorf("error on delete file sysetem %s", err)
	}
	return
}

func importKpfsFileSystem(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

func resourceKsyunKpfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return fmt.Errorf("failed to update file system: %s", err)
}
