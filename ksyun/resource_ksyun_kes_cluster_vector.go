/*
Create a KES vector cluster (KES8, 8.x/9.x) via the shared kmr/2021-09-02
LaunchCluster backend switched into KES mode by ClusterType="KES".

KES8 clusters run Elasticsearch 8.x / 9.x and require Xpack credentials. The
authoritative creation path (confirmed by a Ksyun console capture + a production
probe) is the SHARED kmr/2021-09-02 LaunchCluster action, with the body field
ClusterType="KES" switching the backend from KMR mode into KES8 mode. The SDK
ships no typed client for that version, so this resource builds a custom request
struct and dispatches it through the KES client's generic signing channel (see
service_ksyun_kes_vector.go).

On the kmr/2021-09-02 backend, placement (AvailabilityZone, VpcSubnetId, VpcId)
and the per-group index/ChargeType live INSIDE each InstanceGroups entry, and
VpcSubnetId is also sent at top level. PurchaseTime and VpcEpcSubnetId are not
part of the KES8 shape and are therefore not modelled here.

# Example Usage

```hcl

	resource "ksyun_kes_cluster_vector" "example" {
	  cluster_name      = "my-kes-vector-cluster"
	  main_version      = "8.19.13"
	  distribution      = "ElasticSearch"
	  scenario          = "General"
	  charge_type       = "Monthly"
	  project_id        = 0
	  vpc_domain_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  vpc_subnet_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	  security_group_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

	  # Required for any 8.x / 9.x cluster.
	  xpack_username = "elastic"
	  xpack_password = "your-xpack-password"

	  # Optional KES8 switches (defaults shown).
	  xpack_enabled             = true
	  enable_best_delivery_mode = false
	  enable_eip                = false
	  instance_login_type       = "later"
	  services                  = ["elasticsearch"]
	  plugins                   = ["repository-ks3"]

	  instance_groups {
	    instance_group_type     = "MASTER"
	    instance_type           = "ES.I3.2C4G"
	    instance_type_code      = "ES.I3.2C4G"
	    resource_type           = "KEC"
	    instance_count          = 3
	    instance_group_index    = 1
	    availability_zone       = "cn-beijing-6a"
	    availability_zone_index = 1
	    vpc_id                  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	    vpc_subnet_id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	    charge_type             = "Monthly"
	    volume_type             = "SSD3.0"
	    volume_size             = 100
	    volume_count            = 1
	    system_disk_type        = "SSD3.0"
	    system_disk_size        = 50
	  }

	  instance_groups {
	    instance_group_type     = "DATA"
	    instance_type           = "ES.D4.8C16G"
	    instance_type_code      = "ES.D4.8C16G"
	    resource_type           = "KEC"
	    instance_count          = 3
	    instance_group_index    = 2
	    availability_zone       = "cn-beijing-6a"
	    availability_zone_index = 1
	    vpc_id                  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	    vpc_subnet_id           = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	    charge_type             = "Monthly"
	    volume_type             = "SSD3.0"
	    volume_size             = 500
	    volume_count            = 1
	    system_disk_type        = "SSD3.0"
	    system_disk_size        = 50
	  }
	}

```

# Import

KES vector cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kes_cluster_vector.example ${cluster_id}
```
*/
package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunKesClusterVector() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKesClusterVectorCreate,
		Read:   resourceKsyunKesClusterVectorRead,
		Delete: resourceKsyunKesClusterVectorDelete,
		Importer: &schema.ResourceImporter{
			State: importKesClusterVector,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the KES vector cluster.",
			},
			"main_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The main version of the KES vector cluster, e.g. 8.19.13/9.x. Xpack credentials are required for any 8.x or 9.x cluster.",
			},
			"distribution": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "ElasticSearch",
				Description: "The distribution of the KES vector cluster. Defaults to ElasticSearch.",
			},
			"scenario": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The scenario of the KES vector cluster, e.g. General.",
			},
			"charge_type": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The charge type of the KES vector cluster, e.g. Monthly/Minutely.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The project id of the KES vector cluster.",
			},
			"vpc_domain_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The VPC domain id of the KES vector cluster.",
			},
			"vpc_subnet_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The VPC subnet id of the KES vector cluster. Also sent per instance group on the kmr/2021-09-02 backend.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The security group id of the KES vector cluster.",
			},
			"xpack_username": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Sensitive:   true,
				Description: "The Xpack username. Required for KES8 (8.x/9.x) clusters.",
			},
			"xpack_password": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Sensitive:   true,
				Description: "The Xpack password (at least 6 characters). Required for KES8 (8.x/9.x) clusters.",
			},
			"xpack_enabled": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
				Description: "Whether Xpack is enabled. Defaults to true for KES8 clusters.",
			},
			"enable_best_delivery_mode": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: "Whether best delivery mode is enabled.",
			},
			"enable_eip": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Whether EIP is enabled for the cluster.",
			},
			"instance_login_type": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The instance login type, e.g. 'later' (no password) or 'password'.",
			},
			"instance_password": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "The instance login password. Only sent when instance_login_type is 'password'.",
			},
			"services": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The services to deploy. Defaults to [\"elasticsearch\"] when omitted.",
			},
			"plugins": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The plugins to install, e.g. repository-ks3. Only sent when set.",
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
						"instance_type_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance type code of the instance group, e.g. ES.X6.4C8G.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type of the instance group, e.g. KEC.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The instance count of the instance group.",
						},
						"instance_group_index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The index of the instance group.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The availability zone of the instance group, e.g. cn-beijing-6a.",
						},
						"availability_zone_index": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The availability zone index of the instance group.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The VPC id of the instance group.",
						},
						"vpc_subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The VPC subnet id of the instance group.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The charge type of the instance group, e.g. Monthly/Minutely.",
						},
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The volume type of the instance group, e.g. SSD3.0.",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The volume size (GB) of the instance group.",
						},
						"volume_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The volume count of the instance group.",
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
						"system_disk_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The system disk type of the instance group, e.g. SSD3.0.",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The system disk size (GB) of the instance group.",
						},
					},
				},
				Description: "The instance groups of the KES vector cluster. Placement (availability_zone/vpc_subnet_id/vpc_id) and the per-group index/charge_type are set per group.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the KES vector cluster.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the KES vector cluster (KES).",
			},
			"cluster_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the KES vector cluster.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the KES vector cluster.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the KES vector cluster.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the KES vector cluster.",
			},
			"serving_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The serving minutes of the KES vector cluster.",
			},
			"proxy_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The proxy port of the KES vector cluster.",
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
				Description: "The tags of the KES vector cluster.",
			},
		},
	}
}

func resourceKsyunKesClusterVectorCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kesVectorService := KesVectorService{meta.(*KsyunClient)}
	err = kesVectorService.CreateVectorCluster(d, resourceKsyunKesClusterVector())
	if err != nil {
		return fmt.Errorf("error on creating KES vector cluster %s", err)
	}
	return err
}

func resourceKsyunKesClusterVectorRead(d *schema.ResourceData, meta interface{}) (err error) {
	kesVectorService := KesVectorService{meta.(*KsyunClient)}
	err = kesVectorService.ReadAndSetVectorClusterResource(d, resourceKsyunKesClusterVector())
	if err != nil {
		return fmt.Errorf("error on reading KES vector cluster %s", err)
	}
	logger.Debug(logger.RespFormat, "readKesClusterVector", d)
	return err
}

func resourceKsyunKesClusterVectorDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kesVectorService := KesVectorService{meta.(*KsyunClient)}
	err = kesVectorService.DeleteVectorCluster(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting KES vector cluster %s", err)
	}
	return
}

func importKesClusterVector(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
