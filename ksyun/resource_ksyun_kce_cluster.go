/*
Provides a KCE cluster resource.

~> **NOTE:** We recommend that uses the `ksyun_kce_cluster` resource to create a cluster with few `worker_config`.
If you want to manage more worker instances in this cluster, to use the `ksyun_kce_cluster_attach_existence` or `ksyun_kce_cluster_attachment` resource to attach the worker instances to the cluster. The reason is that the `worker_config` is unchangeable and may cause the cluster to be re-created because it is marked *ForceNew*.

# Example Usage

## basic dependency resources

```hcl

data "ksyun_kce_instance_images" "test" {
  output_file = "output_result"
}

data "ksyun_kce_instance_images" "test" {
}

variable "az" {
  default = "cn-beijing-6e"
}


variable "suffix" {
  default = "-kce-complete"
}
```

## create a ManagementCluster

```hcl

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf-modification${var.suffix}"
  cluster_desc        = "description...modification"
  cluster_manage_mode = "ManagedCluster"
  vpc_id              = ksyun_vpc.test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.252.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.23.17"
  reserve_subnet_id   = ksyun_subnet.reserve.id

  managed_cluster_multi_master {
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = ksyun_security_group.test.id
  }

  worker_config {
    count         = 2
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_worker"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }
}
```

## create a DedicatedCluster

```hcl

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf-modification${var.suffix}"
  cluster_desc        = "description...modification"
  cluster_manage_mode = "DedicateCluster"
  vpc_id              = ksyun_vpc.test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.252.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.23.17"
  reserve_subnet_id   = ksyun_subnet.reserve.id

  managed_cluster_multi_master {
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = ksyun_security_group.test.id
  }

  master_config {
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_master"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }

  worker_config {
    count         = 2
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_worker"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
  }

## config a different machine type
  worker_config {
    count         = 2
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4C"
    instance_name = "tf_kce_worker"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
    advanced_setting {
      container_runtime = "containerd"
      label {
        key   = "tf_assembly_kce"
        value = "advanced_setting"
      }
      taints {
        key    = "key1"
        value  = "value1"
        effect = "NoSchedule"

      }
    }
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
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/helper"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

const (
	kceManagedModeManaged   = "ManagedCluster"
	kceManagedModeDedicated = "DedicatedCluster"
)

var instanceNodeForceNewField = []string{"image_id", "instance_name", "subnet_id", "security_group_id", "charge_type", "instance_type"}

func nodeAdvancedSetting() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"data_disk": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			ForceNew:    true,
			Description: "The mount setting of data disk. **Notes:** Only impact on the first data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"auto_format_and_mount": {
						Type:     schema.TypeBool,
						Optional: true,
						ForceNew: true,
						Description: "Whether to format and mount the data disk, default value: true." +
							" If this field is filled with false, then the file_system and mount_target fields will not take effect.",
					},
					"file_system": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						Description: "The file system of the data disk. The default value is ext4." +
							"Valid values: ext3, ext4, xfs.",
					},
					"mount_target": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "The mount target of the data disk.",
					},
				},
			},
		},
		"container_runtime": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
			Description: "Container Runtime. Valid Values: `docker`, `containerd`.",
			// ValidateFunc: validation.StringInSlice([]string{
			// 	"docker", "containerd",
			// }, false),
		},
		"docker_path": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The storage path of the container. The default value is /data/docker.",
		},
		"container_path": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The storage path of the container. The default value is /data/container. **Notes:** If this path is specified, the docker_path field will be ignored.",
		},
		"user_script": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Description: "A user script encoded in base64, which will be executed on the node **after** the Kubernetes components run. " +
				"Users need to ensure the script's re-entrant and retry logic. The script and its generated logs can be found in the directory /usr/local/ksyun/kce/pre_userscript.",
		},
		"pre_user_script": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Description: "A user script encoded in base64, which will be executed on the node **before** the Kubernetes components run. " +
				"Users need to ensure the script's re-entrant and retry logic. The script and its generated logs can be found in the directory /usr/local/ksyun/kce/pre_userscript.",
		},
		// "schedulable": {
		// 	Type:     schema.TypeBool,
		// 	Optional: true,
		// 	ForceNew: true,
		// 	Computed: true,
		// },
		"label": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The key of label.",
						ForceNew:    true,
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The value of label.",
						ForceNew:    true,
					},
				},
			},
		},
		"extra_arg": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "The extra arguments for the kubelet. The format is key=value. For example, --kubelet-extra-args=\"key1=value1,key2=value2\".",
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
			},
		},
		"container_log_max_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Description:  "Customize the maximum size of the log file. The default value is 100m.",
			ValidateFunc: validation.IntBetween(1, 9999),
		},
		"container_log_max_files": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Description:  "Customize the number of log files. The default value is 10.",
			ValidateFunc: validation.IntBetween(1, 9999),
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "Taints.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "The key of the taint.",
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "The value of the taint.",
					},
					"effect": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "The effect of the taint. Valid values: NoSchedule, PreferNoSchedule, NoExecute.",
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

	m["instance_type"].Computed = false
	m["instance_type"].Required = true
	m["instance_type"].Optional = false

	m["security_group_id"].MaxItems = 1

	m["role"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "Master_Etcd",
		ValidateFunc: validation.StringInSlice([]string{
			"Worker",
			"Master_Etcd", // only one role for master node with tf create.
			"Master", "Etcd",
		}, false),
	}

	m["hashcode"] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		ForceNew:    true,
		Description: "",
	}

	m["advanced_setting"] = &schema.Schema{
		Type: schema.TypeList,
		// MinItems: 1,
		MaxItems:    1,
		Optional:    true,
		ForceNew:    true,
		Description: "Advanced settings.",
		Elem: &schema.Resource{
			Schema: nodeAdvancedSetting(),
		},
	}

	for _, field := range instanceNodeForceNewField {
		m[field].ForceNew = true
	}
	// m["security_group_id"] = &schema.Schema{
	//	Type:     schema.TypeString,
	//	Required: true,
	// }

	return m
}

func instanceForWorkerNode() map[string]*schema.Schema {
	m := instanceForNode()

	m["role"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "Worker",
		ValidateFunc: validation.StringInSlice([]string{
			"Worker",
		}, false),

		Description: "The role of instance. Valid values: Worker.",
	}

	m["count"] = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		ForceNew:    true,
		Description: "The number of worker nodes.",
	}

	return m
}

func instanceForMasterNode() map[string]*schema.Schema {
	m := instanceForNode()

	m["count"] = &schema.Schema{
		Type:         schema.TypeInt,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.IntInSlice([]int{3, 5}),
		Description:  "The number of master nodes. The count of master nodes must be 3 or 5.",
	}

	return m
}

func vkConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The version of the virtual kubelet.",
		},
		"kubeconfig": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The config string of the virtual kubelet.",
		},
		"nodename": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the virtual kubelet node.",
		},
		"namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The namespace of the virtual kubelet.",
		},
		"anonymous_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to enable anonymous authentication.",
		},
		"enable_node_lease": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to enable node lease.",
		},
		"server": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The server of the virtual kubelet.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"listen_port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The port of the virtual kubelet server.",
					},
					"server_cert_file": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The server certificate file of the virtual kubelet.",
					},
					"server_key_file": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The server key file of the virtual kubelet.",
					},
				},
			},
		},
		"openapi": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The openapi of the virtual kubelet.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aksk_config_map": {
						Type:        schema.TypeMap,
						Optional:    true,
						Description: "The aksk config map of the virtual kubelet.",
					},
					"region": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The region of the virtual kubelet.",
					},
					"access_key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The access key of the virtual kubelet.",
					},
					"secret_key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The secret key of the virtual kubelet.",
					},
				},
			},
		},
		"metrics_server": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The metrics server of the virtual kubelet.",
		},
		"log_level": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The log level of the virtual kubelet.",
		},
		"startup_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The startup timeout of the virtual kubelet.",
		},
		"custom_labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The custom labels of the virtual kubelet.",
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Taints.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
						// ForceNew:    true,
						Description: "The key of the taint.",
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
						// ForceNew:    true,
						Description: "The value of the taint.",
					},
					"effect": {
						Type:     schema.TypeString,
						Required: true,
						// ForceNew:    true,
						Description: "The effect of the taint. Valid values: NoSchedule, PreferNoSchedule, NoExecute.",
					},
				},
			},
		},
		"disable_taint": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to disable taint.",
		},
		"capacity": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The capacity of the virtual kubelet.",
		},
		"allow_privileged": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Whether to allow privileged containers.",
		},
		"cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the cluster.",
		},
		"kcilet_heartbeat": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The heartbeat of the virtual kubelet.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"failuret_threshold": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The failure threshold of the virtual kubelet.",
					},
					"period": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "The period of the virtual kubelet.",
					},
				},
			},
		},
		"kcilet_kubeconfig_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The kubeconfig path of the virtual kubelet.",
		},
		"kci_pod_deletion_cost": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The pod deletion cost of the virtual kubelet.",
		},
		"batch_create_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to enable batch creation of pods in the virtual kubelet.",
		},
		"pod_sync_rate": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The rate of pod synchronization in the virtual kubelet.",
		},
		"pod_sync_workers": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The number of workers for pod synchronization in the virtual kubelet.",
		},
		"image_cache_enable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to enable image cache in the virtual kubelet.",
		},
		"pod_sync_bucket_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The bucket size of pod synchronization in the virtual kubelet.",
		},
		"summary_sync_interval": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The interval of summary synchronization in the virtual kubelet.",
		},
		"cluster_dns": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The DNS of the virtual kubelet.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cluster_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The domain of the virtual kubelet.",
		},
		"instance_settings": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The instance settings of the virtual kubelet.",
		},
		"advanced_setting": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The advanced settings of the virtual kubelet.",
		},
	}
}

func componentForNode() map[string]*schema.Schema {
	m := make(map[string]*schema.Schema)
	m["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the component.",
	}
	m["release_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the release.",
	}
	m["namespace"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The namespace of the component.",
	}

	m["version"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The version of the component.",
	}
	m["config"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"config_string": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The config string of the component.",
				},
				"virtual_kubelet": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The config of the virtual kubelet.",
					MaxItems:    1,
					Elem:        &schema.Resource{Schema: vkConfig()},
				},
			},
		},
		Description: "The config of the component.",
	}

	m["resources"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "The resources of the component.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"requests": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The requests of the component.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: resourcesForComponent(),
					},
				},
				"limits": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The limits of the component.",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: resourcesForComponent(),
					},
				},
			},
		},
	}

	return m
}

func resourcesForComponent() map[string]*schema.Schema {
	m := map[string]*schema.Schema{
		"cpu": {
			Type:     schema.TypeString,
			Optional: true,
			// ForceNew:    true,
			Description: "The cpu of the component.",
		},
		"memory": {
			Type:     schema.TypeString,
			Optional: true,
			// ForceNew:    true,
			Description: "The memory of the component.",
		},
	}
	return m
}

// 独立集群和托管集群分开管理？？？
func resourceKsyunKceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKceClusterCreate,
		Update: resourceKsyunKceClusterUpdate,
		Read:   resourceKsyunKceClusterRead,
		Delete: resourceKsyunKceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: importKceCluster,
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
				// Computed: true,
				Default: "DedicatedCluster",
				ValidateFunc: validation.StringInSlice([]string{
					kceManagedModeManaged, // 是否可以先不创建worker？
					kceManagedModeDedicated,
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
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				// ValidateFunc: validation.IntInSlice([]int{16, 32, 64, 128, 256}),
				Description: "The maximum number of pods that can be run on each node. valid values: 16, 32, 64, 128, 256.",
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validation.StringInSlice([]string{"Flannel", "Canal"}, false),
				Description: "The network type of the cluster. valid values: 'Flannel', 'Canal', 'Calico'.",
			},
			"k8s_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// ValidateFunc: validation.StringInSlice([]string{"v1.17.6", "v1.19.3", "v1.21.3"}, false),
				Description: "The latest three kubernetes version. Current valid values:\"v1.25.7\", \"v1.23.17\", \"v1.21.3\"." +
					" **Notes:** The version is updated in real time with the K8s official. Therefore, you can view the maintaining strategies in [Kingsoft Cloud K8s Version Strategies](https://docs.ksyun.com/documents/43229?type=3) and get the latest versions.",
			},
			"reserve_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the reserve subnet.",
			},
			// todo
			"managed_cluster_multi_master": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the subnet for the managed cluster masters.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the security group for the managed cluster masters.",
						},
					},
				},
				Description: "The configuration for the managed cluster multi master. If the cluster_manage_mode is ManagedCluster, this field is **required**.",
			},
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
			"expose_public_api_server": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				// ForceNew:    true,
				Description: "Whether to expose the apiserver to the public network, when cluster type is ManagedCluster, it works.",
			},
			"master_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				// Computed: true,
				// MaxItems: 1,
				Elem: &schema.Resource{
					Schema: instanceForMasterNode(),
				},
				Description: "The configuration for the master nodes. If the cluster_manage_mode is DedicatedCluster, this field is **required**." +
					" **Notes:** work_config block is identified by the **instance_type, subnet_id, security_group_id, role, image_id**. " +
					"If the unique identification is the same, the instance config block is conflict, and then **cause an error**." +
					"If the unique identification is changed, that leads to the cluster **re-creation**.",
			},
			"worker_config": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: instanceForWorkerNode(),
				},
				Description: "The configuration for the worker nodes. If the cluster_manage_mode is ManagedCluster, this field is **required**." +
					" **Notes:** work_config block is identified by the **instance_type, subnet_id, security_group_id, role, image_id**. " +
					"If the unique identification is the same, the instance config block is conflict, and then **cause an error**." +
					"If the unique identification is changed, that leads to the cluster **re-creation**.",
			},

			"component": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The component of the cluster.",
				Set:         helper.HashFuncWithKeys("name", "release_name", "namespace"),
				Elem: &schema.Resource{
					Schema: componentForNode(),
				},
			},

			"master_id_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID list of the master nodes.",
			},
			"worker_id_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID list of the worker nodes.",
			},
			"kube_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration for the kubernetes cluster.",
			},
			"kube_config_intranet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration for the private kubernetes cluster.",
			},
		},
	}
}

func resourceKsyunKceClusterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = kcePreinspection(d)
	if err != nil {
		return
	}

	srv := KceService{meta.(*KsyunClient)}
	err = srv.CreateCluster(d, resourceKsyunKceCluster())
	if err != nil {
		return fmt.Errorf("error on create kce cluster: %s", err)
	}

	if d.HasChange("component") {
		err = srv.InstallComponentInCluster(d, resourceKsyunKceCluster())
		if err != nil {
			return fmt.Errorf("error on install component in kce cluster: %s", err)
		}
	}

	if d.HasChange("expose_public_api_server") {
		err = srv.ModifyPublicApiServer(d, resourceKsyunKceCluster())
		if err != nil {
			logger.Info("error on expose public api server in kce cluster: %s", err)
		}
	}

	return resourceKsyunKceClusterRead(d, meta)
}

func resourceKsyunKceClusterUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = kcePreinspection(d)
	if err != nil {
		return
	}

	srv := KceService{meta.(*KsyunClient)}
	if d.HasChange("component") {
		err = srv.InstallComponentInCluster(d, resourceKsyunKceCluster())
		if err != nil {
			return fmt.Errorf("error on install component in kce cluster: %s", err)
		}
	}

	if d.HasChanges("cluster_name", "cluster_desc") {
		err = srv.UpdateCluster(d, resourceKsyunKceCluster())
		if err != nil {
			return fmt.Errorf("error on update kce cluster: %s", err)
		}
	}
	if d.HasChange("expose_public_api_server") {
		err = srv.ModifyPublicApiServer(d, resourceKsyunKceCluster())
		if err != nil {
			logger.Info("error on expose public api server in kce cluster: %s", err)
			vv := d.Get("expose_public_api_server")
			switch vv.(type) {
			case bool:
				vvv := vv.(bool)
				d.Set("expose_public_api_server", !vvv) // reset the value to false if error occurs
			}
		}
	}

	return resourceKsyunKceClusterRead(d, meta)
}

func resourceKsyunKceClusterRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceService{meta.(*KsyunClient)}
	err = srv.ReadAndSetKceCluster(d, resourceKsyunKceCluster())
	if err != nil {
		return fmt.Errorf("error on create kce cluster: %s", err)
	}

	if _, ok := d.GetOk("component"); ok {
		err = srv.ReadComponentInCluster(d, resourceKsyunKceCluster())
		if err != nil {
			return fmt.Errorf("error on read component in kce cluster: %s", err)
		}
	}

	if err = srv.ReadKubeConfigInCluster(d, resourceKsyunKceCluster()); err != nil {
		return fmt.Errorf("error on read kube_config in kce cluster: %s", err)
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

func kcePreinspection(d *schema.ResourceData) error {
	kceManagedMode := d.Get("cluster_manage_mode").(string)

	blockErrFmt := "%s is duplication machine type block, the unique identification of master_config block conbine with instance_type, subnet_id, security_group_id, role, image_id. " +
		"If the unique identification is the same, the instance config block is duplication. Please check the master_config block. " +
		"And for the details of instance config block, please refer to the document https://registry.terraform.io/providers/kingsoftcloud/ksyun/latest/docs/resources/kce_cluster"

	switch kceManagedMode {
	case kceManagedModeManaged:
		wc, ok := d.GetOk("worker_config")
		if !ok {
			return fmt.Errorf("worker_config is required when cluster_manage_mode is %s", kceManagedMode)
		}
		if helper.IsEmpty(wc) {
			return fmt.Errorf("worker_config is required when cluster_manage_mode is %s", kceManagedMode)
		}
		if isDuplicationMachineType(wc.([]interface{})) {
			return fmt.Errorf(blockErrFmt, "worker_config")
		}

		mc, ok := d.GetOk("master_config")
		if ok {
			if !helper.IsEmpty(mc) {
				return fmt.Errorf("you don't need define master_config when cluster_manage_mode is %s", kceManagedMode)
			}
		}

		if helper.IsEmpty(d.Get("managed_cluster_multi_master")) {
			return fmt.Errorf("managed_cluster_multi_master is required when cluster_manage_mode is %s", kceManagedMode)
		}

	case kceManagedModeDedicated:
		mc, ok := d.GetOk("master_config")
		if ok {
			if helper.IsEmpty(mc) {
				return fmt.Errorf("master_config is required when cluster_manage_mode is %s", kceManagedMode)
			}
			if isDuplicationMachineType(mc.([]interface{})) {
				return fmt.Errorf(blockErrFmt, "master_config")
			}
		} else {
			return fmt.Errorf("master_config is required when cluster_manage_mode is %s", kceManagedMode)
		}
	}
	return nil
}

func isDuplicationMachineType(machines []interface{}) bool {
	m := make(map[int]bool)
	for _, machine := range machines {
		hashcode := kceInstanceNodeHashFunc()(machine)
		if _, ok := m[hashcode]; ok {
			return true
		}
		m[hashcode] = true
	}
	return false
}
