package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceKsyunInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstanceCreate,
		Update: resourceKsyunInstanceUpdate,
		Read:   resourceKsyunInstanceRead,
		Delete: resourceKsyunInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"stopped",
				}, false),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SSD3.0",
								"EHDD",
								"Local_SSD",
							}, false),
						},
						"disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(20, 500),
						},
					},
				},
			},
			"data_disk_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16000),
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 8,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SSD3.0",
								"EHDD",
								"Local_SSD",
								"ESSD_PL1",
								"ESSD_PL2",
								"ESSD_PL3",
							}, false),
						},
						"disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(10, 16000),
						},
						"disk_snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extension_network_interface": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: kecNetworkInterfaceHash,
			},

			"local_volume_snapshot_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: kecImportDiffSuppress,
			},

			"instance_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"keep_image_login": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: kecImportDiffSuppress,
			},

			"key_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Daily",
					"HourlyInstantSettlement",
				}, false),
			},
			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
				ValidateFunc:     validation.IntBetween(0, 36),
			},
			"security_group_id": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sriov_net_support": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"data_guard_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: kecImportDiffSuppress,
			},
			"iam_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force_reinstall_system": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dns1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			//"has_init_info": {
			//	Type:     schema.TypeBool,
			//	Computed: true,
			//},
			//some control
			"has_modify_system_disk": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_modify_password": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_modify_keys": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"force_delete": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "this field is Deprecated and no effect for change",
			},

			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.createKecInstance(d, resourceKsyunInstance())
	if err != nil {
		return fmt.Errorf("error on creating Instance: %s", err)
	}
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.readAndSetKecInstance(d, resourceKsyunInstance())
	if err != nil {
		return fmt.Errorf("error on reading Instance: %s", err)
	}
	return err
}

func resourceKsyunInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.modifyKecInstance(d, resourceKsyunInstance())
	if err != nil {
		return fmt.Errorf("error on updating Instance: %s", err)
	}
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.removeKecInstance(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting Instance: %s", err)
	}
	return err
}
