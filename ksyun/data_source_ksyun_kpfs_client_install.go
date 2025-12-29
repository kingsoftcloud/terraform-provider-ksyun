/*
Query KPFS client installation package information and mount IP by file system ID.

# Example Usage

```hcl

	data "ksyun_kpfs_client_install" "default" {
		output_file="output_result"
		id         = "b7449ea1d57f428595f7c68a1fbeeafd"
	}

```
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunKpfsClientInstall() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKpfsClientInstallRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File system ID.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of KPFS clusters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_data_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data IP of the KPFS cluster.",
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The download URL for the KPFS client.",
						},
						"os_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The OS version supported by the KPFS client.",
						},
						"kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The kernel version supported by the KPFS client.",
						},
						"nic_driver": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The NIC driver version supported by the KPFS client.",
						},
					},
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
		},
	}
}

func dataSourceKsyunKpfsClientInstallRead(d *schema.ResourceData, meta interface{}) error {
	kpfsService := KpfsService{meta.(*KsyunClient)}
	return kpfsService.ReadKpfsClientInstallInfo(d, dataSourceKsyunKpfsClientInstall())
}
