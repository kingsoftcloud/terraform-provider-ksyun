/*
Provides a KCE cluster resource.

# Example Usage

```hcl

data "ksyun_kce_instance_images" "test" {
  output_file = "output_result"
}

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf_test_cluster"
  cluster_desc        = "description..."
  cluster_manage_mode = "DedicatedCluster"
  vpc_id              = ksyun_vpc.tf_test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.254.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.19.3"
  reserve_subnet_id   = ksyun_subnet.tf_test_reserve_subnet.id

  master_config {
    role          = "Master_Etcd"
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.tf_test_subnet.id
    security_group_id = [ksyun_security_group.default.id]
    charge_type       = "Daily"
  }
}
```

# Import

KCE cluster can be imported using the id, e.g.

```
$ terraform import ksyun_kce_cluster.default 67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

//func nodeAdvancedSetting() map[string]*schema.Schema {
//	return map[string]*schema.Schema{
//		"data_disk": {
//			Type:     schema.TypeList,
//			MaxItems: 1,
//			Elem: &schema.Resource{
//				Schema: map[string]*schema.Schema{
//					"auto_format_and_mount": {
//						Type:     schema.TypeBool,
//						Optional: true,
//					},
//					"file_system": {
//						Type:     schema.TypeString,
//						Optional: true,
//					},
//					"mount_target": {
//						Type:     schema.TypeString,
//						Optional: true,
//					},
//				},
//			},
//		},
//		"container_runtime": {
//			Type:     schema.TypeString,
//			Optional: true,
//			ValidateFunc: validation.StringInSlice([]string{
//				"docker", "containerd",
//			}, false),
//		},
//		"docker_path": {
//			Type:     schema.TypeString,
//			Optional: true,
//		},
//		"container_path": {
//			Type:     schema.TypeString,
//			Optional: true,
//		},
//		"user_script": {
//			Type:     schema.TypeString,
//			Optional: true,
//		},
//		"pre_user_script": {
//			Type:     schema.TypeString,
//			Optional: true,
//		},
//		"schedulable": {
//			Type:     schema.TypeBool,
//			Optional: true,
//		},
//		"label": {
//			Type:     schema.TypeList,
//			Optional: true,
//			Elem: &schema.Resource{
//				Schema: map[string]*schema.Schema{
//					"key": {
//						Type:     schema.TypeString,
//						Required: true,
//					},
//					"value": {
//						Type:     schema.TypeString,
//						Required: true,
//					},
//				},
//			},
//		},
//		"extra_arg": {
//			Type:     schema.TypeList,
//			Optional: true,
//			Elem: &schema.Schema{
//				Type: schema.TypeString,
//			},
//		},
//		"container_log_max_size": {
//			Type:     schema.TypeInt,
//			Optional: true,
//		},
//		"container_log_max_files": {
//			Type:     schema.TypeInt,
//			Optional: true,
//		},
//		"taint": {
//			Type:     schema.TypeList,
//			Optional: true,
//			Elem: &schema.Resource{
//				Schema: map[string]*schema.Schema{
//					"key": {
//						Type: schema.TypeString,
//					},
//					"value": {
//						Type: schema.TypeString,
//					},
//					"effect": {
//						Type: schema.TypeString,
//					},
//				},
//			},
//		},
//	}
//}

func instanceForNode() map[string]*schema.Schema {
	m := instanceConfig()

	m["key_id"].Computed = true
	m["tags"].Computed = true

	m["count"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	m["role"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			//"Worker",
			"Master_Etcd", "Master", "Etcd",
		}, false),
	}
	//m["security_group_id"] = &schema.Schema{
	//	Type:     schema.TypeString,
	//	Required: true,
	//}

	return m
	//	//"advanced_setting": {
	//	//	Type: schema.TypeSet,
	//	//	//MinItems: 1,
	//	//	MaxItems: 1,
	//	//	Elem: &schema.Resource{
	//	//		Schema: nodeAdvancedSetting(),
	//	//	},
	//	//	//Elem:     nodeAdvancedSetting(),
	//	//},
	//}

}

// 独立集群和托管集群分开管理？？？
func resourceKsyunKceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKceClusterCreate,
		Update: resourceKsyunKceClusterUpdate,
		Read:   resourceKsyunKceClusterRead,
		Delete: resourceKsyunKceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cluster.",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cluster.",
			},
			"cluster_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the cluster.",
			},
			"cluster_manage_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Default:  "DedicatedCluster",
				ValidateFunc: validation.StringInSlice([]string{
					// "ManagedCluster", // 是否可以先不创建worker？
					"DedicatedCluster",
				}, false),
				Description: "The management mode of the master node.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC.",
			},
			"pod_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "The pod CIDR block.",
			},
			"service_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
				Description:  "The service CIDR block.",
			},
			"max_pod_per_node": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{16, 32, 64, 128, 256}),
				Description:  "The maximum number of pods that can be run on each node. valid values: 16, 32, 64, 128, 256.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Flannel", "Canal"}, false),
				Description:  "The network type of the cluster. valid values: 'Flannel', 'Canal'.",
			},
			"k8s_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"v1.17.6", "v1.19.3", "v1.21.3"}, false),
				Description:  "kubernetes version, valid values:\"v1.17.6\", \"v1.19.3\", \"v1.21.3\".",
			},
			"reserve_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the reserve subnet.",
			},
			// todo
			//"managed_cluster_multi_master": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"subnet_id": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"security_group_id": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//		},
			//	},
			//},
			"master_etcd_separate": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "The deployment method for the Master and Etcd components of the cluster. " +
					"if set to True, Deploy the Master and Etcd components on dedicated nodes. " +
					"if set to false, Deploy the Master and Etcd components on shared nodes.",
			},
			"public_api_server": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "Whether to expose the apiserver to the public network. " +
					"If not needed, do not fill in this option. " +
					"If selected, a public SLB and EIP will be created to enable public access to the cluster's API server. " +
					"Users need to pass the Elastic IP creation pass-through parameter, which should be a JSON-formatted string.",
			},
			"master_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: instanceForNode(),
				},
				Description: "The configuration for the master nodes.",
			},
		},
	}
}

func resourceKsyunKceClusterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceService{meta.(*KsyunClient)}
	err = srv.CreateCluster(d, resourceKsyunKceCluster())
	if err != nil {
		return fmt.Errorf("error on create kce cluster: %s", err)
	}
	return
}
func resourceKsyunKceClusterUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return
}
func resourceKsyunKceClusterRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceService{meta.(*KsyunClient)}
	err = srv.ReadAndSetKceCluster(d, resourceKsyunKceCluster())
	if err != nil {
		return fmt.Errorf("error on create kce cluster: %s", err)
	}
	return
}
func resourceKsyunKceClusterDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceService{meta.(*KsyunClient)}
	err = srv.DeleteKceCluster(d, resourceKsyunKceCluster())
	if err != nil {
		return fmt.Errorf("error on delete kce cluster: %s", err)
	}
	return
}
