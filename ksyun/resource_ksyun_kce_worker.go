package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunKceWorker() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKceWorkerCreate,
		Update: resourceKsyunKceWorkerUpdate,
		Read:   resourceKsyunKceWorkerRead,
		Delete: resourceKsyunKceWorkerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the kce cluster",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the kec instance. The instance will be shut down while being added to the kce cluster.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the image which support KCE.",
			},
			"instance_password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Computed:    true,
				Description: "The password of the instance.",
			},
			//"key_id": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "The certificate id of the instance.",
			//},

			// todo:
			//"force_stop": {
			//	Type:        schema.TypeBool,
			//	Default:     true,
			//	Optional:    true,
			//	Description: "Whether to force shutdown before the instance joins the cluster, with the default value of true. If set to false, a normal shutdown will be performed.",
			//},

			// 创建接口有advancedSetting，但是考虑到部分字段在创建后就不再更新，比如驱逐状态、label等，会使用节点上的字段返回
			// 因此这里将advancedSetting中的字段拉平，并在逻辑中处理字段值的更新
			"data_disk": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Data Disk config.",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_format_and_mount": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Whether to format and mount the data disk, with the default value of true. If set to false, the FileSystem and MountTarget fields will not take effect.",
						},
						"file_system": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The file system of the data disk, with optional values of ext3, ext4, and xfs. The default value is ext4. If the disk already has a file system, no processing will be performed. If there is no file system, it will be formatted according to the user's definition, only taking effect on the first disk.",
						},
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The mounting point of the data disk, which will be mounted and only take effect on the first disk.",
						},
					},
				},
			},

			"container_runtime": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"docker",
					"containerd",
				}, false),
				Description: "container runtime instruction.",
			},

			"docker_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The storage path of the container. If not specified, the default is /data/docker.",
			},
			"container_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The storage path of the container. If not specified, the default is /data/container. Note: when this parameter is passed, the DockerPath parameter is invalid.",
			},
			"user_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The user script in base64 encoding. This script will be executed on the node after the k8s component runs. Users need to ensure the re-entry and retry logic of the script. The script and the generated log file can be found in the /usr/local/ksyun/kce/userscript directory.",
			},
			"pre_user_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The user script in base64 encoding. This script will be executed on the node before the k8s component runs. Users need to ensure the re-entry and retry logic of the script. The script and the generated log file can be found in the /usr/local/ksyun/kce/pre_userscript directory.",
			},
			"schedulable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				//Computed:    true,
				Description: "Whether the node can be normally scheduled after being added to the cluster. The default is true.",
			},
			// todo: label部分的openapi不完整，需要按照console的功能实现一套
			//"label": {
			//	Type:        schema.TypeList,
			//	Optional:    true,
			//	Description: "The labels that are pre-set when the node is added to the cluster.",
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"key": {
			//				Type:        schema.TypeString,
			//				Required:    true,
			//				Description: "label key.",
			//			},
			//			"value": {
			//				Type:        schema.TypeString,
			//				Required:    true,
			//				Description: "label value.",
			//			},
			//		},
			//	},
			//},
			"extra_arg": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Custom parameters for k8s components on the node.",
			},
			"container_log_max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     100,
				Description: "The maximum size of a container log file. When the size of a container log file reaches this limit, a new container log file is generated for data writing. The default value is 100 MB.",
			},
			"container_log_max_files": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     10,
				Description: "Specify custom data to configure a node, namely, specify the script to run after you deploy the node. You must ensure the reentrancy and retry logic of the script. You can view the script and its log files in the /usr/local/ksyun/kce/userscript directory of the node.",
			},
		},
	}
}

func resourceKsyunKceWorkerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KceWorkerService{
		meta.(*KsyunClient),
	}
	err = s.AddWorker(d, resourceKsyunKceWorker())
	return
}
func resourceKsyunKceWorkerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := KceWorkerService{
		meta.(*KsyunClient),
	}
	err = s.UpdateWorker(d, resourceKsyunKceWorker())
	return
}
func resourceKsyunKceWorkerRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceWorkerService{meta.(*KsyunClient)}
	err = srv.ReadAndSetWorker(d, resourceKsyunKceWorker())
	if err != nil {
		return fmt.Errorf("error on create kce worker: %s", err)
	}
	return
}
func resourceKsyunKceWorkerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := KceWorkerService{meta.(*KsyunClient)}
	err = srv.DeleteKceWorker(d, resourceKsyunKceWorker())
	if err != nil {
		return fmt.Errorf("error on delete kce cluster: %s", err)
	}
	return
}
