/*
This data source provides a list of SSH Key resources according to their SSH Key ID.

# Example Usage

```hcl

	data "ksyun_ssh_keys" "default" {
	  output_file="output_result"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunSSHKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSSHKeysRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of SSH Key IDs, all the SSH Key belong to this region will be retrieved if the ID is `\"\"`.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to filter results by kay name.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of keys that satisfy the condition.",
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_names"},
				Description:   "ssh key name.",
			},
			"key_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"key_name"},
				Description:   "a list of ssh key name.",
			},

			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of ssh key. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the key.",
						},
						"key_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of the key.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "public key.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time of the key.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunSSHKeysRead(d *schema.ResourceData, meta interface{}) error {
	sksService := SksService{meta.(*KsyunClient)}
	return sksService.ReadAndSetKeys(d, dataSourceKsyunSSHKeys())
}
