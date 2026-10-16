/*
Create a KES cluster (standard / KES7, LaunchCluster v2020-12-15).

# Example Usage

```hcl

	resource "ksyun_kes_cluster" "example" {
	  cluster_name      = "my-kes-cluster"
	  main_version      = "7.10.0"
	  scenario          = "general"
	  availability_zone = "cn-beijing-6a"
	  charge_type       = "Monthly"
	  purchase_time     = 1
	  project_id        = 0
	  vpc_domain_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  vpc_subnet_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  security_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

	  instance_groups {
	    instance_group_type = "MASTER"
	    instance_type       = "ES.I3.2C4G"
	    resource_type       = "KEC"
	    instance_count      = 3
	    volume_type         = "SSD3.0"
	    volume_size         = 100
	    volume_count        = 1
	    system_disk_type    = "SSD3.0"
	    system_disk_size    = 50
	  }

	  instance_groups {
	    instance_group_type = "DATA"
	    instance_type       = "ES.D4.8C16G"
	    resource_type       = "KEC"
	    instance_count      = 3
	    volume_type         = "SSD3.0"
	    volume_size         = 500
	    volume_count        = 1
	    system_disk_type    = "SSD3.0"
	    system_disk_size    = 50
	  }
	}

```

# Import

KES cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kes_cluster.example ${cluster_id}
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunKesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKesClusterCreate,
		Read:   resourceKsyunKesClusterRead,
		Delete: resourceKsyunKesClusterDelete,
		Importer: &schema.ResourceImporter{
			State: importKesCluster,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the KES cluster.",
			},
			"main_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The main version of the KES cluster, e.g. 7.10.0/7.4.2/6.8.4/5.6.16.",
			},
			"scenario": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The scenario of the KES cluster.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The availability zone of the KES cluster.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The charge type of the KES cluster.",
			},
			"purchase_time": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The purchase time of the KES cluster.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The project id of the KES cluster.",
			},
			"vpc_domain_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The VPC domain id of the KES cluster.",
			},
			"vpc_subnet_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The VPC subnet id of the KES cluster.",
			},
			"vpc_epc_subnet_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The VPC EPC subnet id of the KES cluster.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The security group id of the KES cluster.",
			},
			"instance_groups": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_group_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the instance group, e.g. MASTER/DATA.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The instance type of the instance group.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type of the instance group.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The instance count of the instance group.",
						},
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The volume type of the instance group.",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The volume size (GB) of the instance group.",
						},
						"multi_instance_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The multi instance count of the instance group.",
						},
						"raid_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The raid type of the instance group.",
						},
						"volume_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The volume count of the instance group.",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The system disk type of the instance group.",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The system disk size (GB) of the instance group.",
						},
					},
				},
				Description: "The instance groups of the KES cluster.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the KES cluster.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the KES cluster.",
			},
			"cluster_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the KES cluster.",
			},
			"enable_eip": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether EIP is enabled.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the KES cluster.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the KES cluster.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the KES cluster.",
			},
			"serving_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The serving minutes of the KES cluster.",
			},
			"proxy_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The proxy port of the KES cluster.",
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "The tags of the KES cluster.",
			},
		},
	}
}

func resourceKsyunKesClusterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kesService := KesService{meta.(*KsyunClient)}
	err = kesService.CreateCluster(d, resourceKsyunKesCluster())
	if err != nil {
		return fmt.Errorf("error on creating KES cluster %s", err)
	}
	return err
}

func resourceKsyunKesClusterRead(d *schema.ResourceData, meta interface{}) (err error) {
	kesService := KesService{meta.(*KsyunClient)}
	err = kesService.ReadAndSetClusterResource(d, resourceKsyunKesCluster())
	if err != nil {
		return fmt.Errorf("error on reading KES cluster %s", err)
	}
	logger.Debug(logger.RespFormat, "readKesCluster", d)
	return err
}

func resourceKsyunKesClusterDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kesService := KesService{meta.(*KsyunClient)}
	err = kesService.DeleteCluster(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting KES cluster %s", err)
	}
	return
}

func importKesCluster(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
