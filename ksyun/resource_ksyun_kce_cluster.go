package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func nodeAdvancedSetting() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data_disk": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"auto_format_and_mount": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"file_system": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"mount_target": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"container_runtime": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"docker", "containerd",
			}, false),
		},
		"docker_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"container_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_script": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pre_user_script": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"schedulable": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"label": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"extra_arg": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"container_log_max_size": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"container_log_max_files": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"taint": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type: schema.TypeString,
					},
					"value": {
						Type: schema.TypeString,
					},
					"effect": {
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

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
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_manage_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ManagedCluster", // 是否可以先不创建worker？
					"DedicatedCluster",
				}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pod_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"service_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},
			"max_pod_per_node": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{16, 32, 64, 128, 256}),
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Flannel", "Canal"}, false),
			},
			"k8s_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"v1.17.6", "v1.19.3", "v1.21.3"}, false),
			},
			"reserve_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			"public_api_server": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"master_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: instanceForNode(),
				},
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
