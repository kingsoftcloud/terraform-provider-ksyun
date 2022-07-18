package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var instanceForNode = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Worker", "Master_Etcd", "Master", "Etcd",
				}, false),
			},
			"node_config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"para": {
							Type:     schema.TypeString,
							Required: true,
						},
						"advanced_setting": {},
					},
				},
			},
		},
	},
}

func resourceKsyunKceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKceClusterCreate,
		Update: resourceKsyunKceClusterUpdate,
		Read:   resourceKsyunKceClusterRead,
		Delete: resourceKsyunKceClustereDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
					"ManagedCluster",
					"DedicatedCluster",
				}, false),
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
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
			"managed_cluster_multi_master": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"master_etcd_separate": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"public_api_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_for_node": &instanceForNode,
		},
	}
}

func resourceKsyunKceClusterCreate(d *schema.ResourceData, meta interface{}) (err error) {

}
func resourceKsyunKceClusterUpdate(d *schema.ResourceData, meta interface{}) (err error) {

}
func resourceKsyunKceClusterRead(d *schema.ResourceData, meta interface{}) (err error) {

}
func resourceKsyunKceClustereDelete(d *schema.ResourceData, meta interface{}) (err error) {

}
