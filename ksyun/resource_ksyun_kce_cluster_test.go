package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/structor/v1/kce"
)

const testAccKceClusterConfig = `
data "ksyun_kce_instance_images" "test" {
}

variable az {
  default = "cn-beijing-6c"
}


variable "suffix" {
  default = "-kce-complete"
}

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf_kce_cluster"
  cluster_desc        = "description..."
  cluster_manage_mode = "DedicatedCluster"
  vpc_id              = "87e64f91-08eb-405c-8b22-75e075f89aca"
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.252.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.23.17"
  reserve_subnet_id   = "d12b6191-1c4f-433e-b760-419ef79673a3"

  master_config {
    role          = "Master_Etcd"
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "c771027a-fafd-4b3b-a6b9-daeab9d0c13a"
    security_group_id = ["59a87036-dc27-41cf-98ab-24a387501195"]
    charge_type       = "Daily"

	advanced_setting {
	  container_runtime = "containerd"
	  docker_path       = "/data/docker_new"
	  user_script       = "abc"
	}
  }

  worker_config {
    count         = 1
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "c771027a-fafd-4b3b-a6b9-daeab9d0c13a"
    security_group_id = ["59a87036-dc27-41cf-98ab-24a387501195"]
    charge_type       = "Daily"
	advanced_setting {
	  container_runtime = "containerd"
	  pre_user_script   = "def"
      label {
		key  = "tf_assembly_kce"
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
    instance_type = "S6.2A"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "c771027a-fafd-4b3b-a6b9-daeab9d0c13a"
    security_group_id = ["59a87036-dc27-41cf-98ab-24a387501195"]
    charge_type       = "Daily"
	advanced_setting {
	  container_runtime = "containerd"
	  pre_user_script   = "def"
      label {
		key  = "tf_assembly_kce"
		value = "advanced_setting"
      }
	  taints {
		key    = "key3"
		value  = "value3"
		effect = "NoSchedule"
      }
	}
  }
}
`

func TestAccKsyunKceCluster_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kce_cluster.default",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckKceWorkerDestroy,
		Steps: []resource.TestStep{
			{
				// PlanOnly: true,
				Config: testAccKceClusterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKceWorkerExists("ksyun_kce_cluster.default", &val),
					testAccCheckKceWorkerAttributes(&val),
				),
			},
			// {
			//	ResourceName:      "ksyun_kce_worker.foo",
			//	ImportStateId:     "bedfb5d0-bb8f-40dd-9f1d-8966fd1ace87:3a66ad5a-313b-41c0-8b72-6749e438ea17",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"lang"},
			// },
		},
	})
}

func TestAccKsyunKceVKKci_basic(t *testing.T) {
	// var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_kce_cluster.default",
		Providers:     testAccProviders,
		// CheckDestroy:  testAccCheckKceWorkerDestroy,
		Steps: []resource.TestStep{
			{
				// PlanOnly: true,
				Config: testAccKceClusterVKConfig,
				// Check: resource.ComposeTestCheckFunc(
				// 	testAccCheckKceWorkerExists("ksyun_kce_cluster.default", &val),
				// 	testAccCheckKceWorkerAttributes(&val),
				// ),
			},
			// {
			//	ResourceName:      "ksyun_kce_worker.foo",
			//	ImportStateId:     "bedfb5d0-bb8f-40dd-9f1d-8966fd1ace87:3a66ad5a-313b-41c0-8b72-6749e438ea17",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	//ImportStateVerifyIgnore: []string{"lang"},
			// },
		},
	})
}

const testAccKceManagementClusterConfig = `
data "ksyun_kce_instance_images" "test" {
}

variable az {
  default = "cn-beijing-6c"
}


variable "suffix" {
  default = "-kce-complete"
}

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf_kce_cluster"
  cluster_desc        = "description..."
  cluster_manage_mode = "ManagedCluster"
  vpc_id              = "87e64f91-08eb-405c-8b22-75e075f89aca"
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.254.0.0/16"
  network_type        = "Flannel"
  k8s_version         = "v1.25.3"
  reserve_subnet_id   = "d12b6191-1c4f-433e-b760-419ef79673a3"
  managed_cluster_multi_master {
	subnet_id = "4448cd56-17e2-4e7f-b04a-9d067792857d"
	security_group_id = "59a87036-dc27-41cf-98ab-24a387501195"
  }

  worker_config {
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set.0.image_id
    instance_type = "S6.4B"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = "4448cd56-17e2-4e7f-b04a-9d067792857d"
    security_group_id = ["59a87036-dc27-41cf-98ab-24a387501195"]
    charge_type       = "Daily"
  }
}
`

const testAccKceClusterVKConfig = `
data "ksyun_kce_instance_images" "test" {
}

variable "region" {
  default = "cn-beijing-6"
  type    = string

}

variable "az" {
  default = "cn-beijing-6e"
  type    = string
}


variable "suffix" {
  default = "-agora-demo"
  type    = string
}

resource "ksyun_vpc" "test" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

# resource "ksyun_vpc" "test-2" {
#   vpc_name   = "tf-alb-vpc"
#   cidr_block = "10.0.0.0/16"
# }

resource "ksyun_subnet" "normal" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}


resource "ksyun_subnet" "reserve" {
  subnet_name       = "tf-alb-subnet-reserve"
  cidr_block        = "10.0.101.0/24"
  subnet_type       = "Reserve"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.test.id
}

resource "ksyun_security_group" "test" {
  vpc_id              = ksyun_vpc.test.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "test" {
  security_group_id = ksyun_security_group.test.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}

resource "ksyun_kce_cluster" "default" {
  cluster_name        = "tf-managed-cluster${var.suffix}"
  cluster_desc        = "description..."
  cluster_manage_mode = "ManagedCluster"
  vpc_id              = ksyun_vpc.test.id
  pod_cidr            = "172.16.0.0/16"
  service_cidr        = "10.252.0.0/16"
  network_type        = "Calico"
  k8s_version         = "v1.23.17"
  reserve_subnet_id   = ksyun_subnet.reserve.id

  managed_cluster_multi_master {
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = ksyun_security_group.test.id
  }
  worker_config {
    count         = 3
    image_id      = data.ksyun_kce_instance_images.test.image_set[1].image_id
    instance_type = "S6.4B"
    instance_name = "tf_kce_worker${var.suffix}"
    system_disk {
      disk_size = 20
      disk_type = "SSD3.0"
    }
    subnet_id         = ksyun_subnet.normal.id
    security_group_id = [ksyun_security_group.test.id]
    charge_type       = "Daily"
  }

  component {
    name         = "virtual-kubelet"
    namespace    = "kube-system"
    release_name = "vk"
    config {
      virtual_kubelet {
        nodename               = "tf-vk"
        kcilet_kubeconfig_path = "/root/.kube/config"
        kci_pod_deletion_cost  = 0
        anonymous_auth         = false
        instance_settings = {
          "k8s.ksyun.com/kci-security-group-id" : ksyun_security_group.test.id,
          "k8s.ksyun.com/kci-host-aliases" : "[ ]",
          "k8s.ksyun.com/kci-kube-proxy-enabled" : "false",
          "k8s.ksyun.com/kci-subnet-id" : ksyun_subnet.normal.id,
        }

        openapi {
          aksk_config_map = {
            "name" : "user-temp-aksk",
            "namespace" : "kube-system",
          }
          region = var.region
        }

      }
    }
    resources {
      requests {
        cpu    = "0.1"
        memory = "256MiB"
      }
      limits {
        cpu    = "2"
        memory = "6144MiB"
      }
    }
  }
}
`

func TestHandleAdvancedSetting(t *testing.T) {

	advanced := kce.AdvancedSetting{
		DataDisk: &kce.DataDisk{
			AutoFormatAndMount: true,
			FileSystem:         "ext4",
			MountTarget:        "/data",
		},
		PreUserScript: "adds",
		UserScript:    "sdasad",
		Schedulable:   false,
		Label:         []kce.Label{{Key: "tf_assembly_kce", Value: "advanced_setting"}},
		ExtraArg: &kce.ExtraArg{
			Kubelet: []string{"custom_arg"},
		},
		ContainerRuntime:     "containerd",
		ContainerPath:        "/data/containerd",
		DockerPath:           "/data/docker",
		ContainerLogMaxFiles: 200,
		ContainerLogMaxSize:  2000,
	}
	m := handleAdvancedSetting2Map(advanced)
	t.Logf("%+v", m)
}
