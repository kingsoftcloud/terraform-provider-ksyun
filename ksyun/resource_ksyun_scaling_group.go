/*
Provides a ScalingGroup resource.

# Example Usage

```hcl

	resource "ksyun_scaling_group" "foo" {
	  subnet_id_set = [ksyun_subnet.foo.id]
	  security_group_id = ksyun_security_group.foo.id
	  scaling_configuration_id = ksyun_scaling_configuration.foo.id
	  min_size = 0
	  max_size = 2
	  desired_capacity = 0
	  status = "Active"
	  slb_config_set  {
	    slb_id = ksyun_lb.foo.id}
	    listener_id = ksyun_lb_listener.foo.id
	    server_port_set = [80]
	  }
	}

```

# Import

scalingGroup can be imported using the `id`, e.g.

```
$ terraform import ksyun_scaling_group.example scaling-group-abc123456
```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"strings"
	"time"
)

func resourceKsyunScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingGroupCreate,
		Read:   resourceKsyunScalingGroupRead,
		Delete: resourceKsyunScalingGroupDelete,
		Update: resourceKsyunScalingGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(diff *schema.ResourceDiff, i interface{}) (err error) {
			if diff.HasChange("security_group_id") {
				err = diff.SetNewComputed("security_group_id_set")
			}
			if diff.HasChange("security_group_id_set") {
				err = diff.SetNewComputed("security_group_id")
			}
			return err
		},
		Schema: map[string]*schema.Schema{

			"scaling_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tf-scaling-group",
				Description: "The Name of the desired ScalingGroup.",
			},
			"scaling_configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Scaling Configuration ID of the desired ScalingGroup set to.",
			},

			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecScalingGroupSize,
				Description:  "The Min KEC instance size of the desired ScalingGroup set to.Valid Value 0-1000.",
			},

			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecScalingGroupSize,
				Description:  "The Max KEC instance size of the desired ScalingGroup set to.Valid Value 0-1000.",
			},

			"desired_capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateKecScalingGroupDesiredCapacity,
				Description:  "The Desire Capacity KEC instance count of the desired ScalingGroup set to.Valid Value 0-1000.",
			},

			"remove_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "RemoveOldestInstance",
				ValidateFunc: validateKecScalingGroupRemovePolicy,
				Description:  "The KEC instance remove policy of the desired ScalingGroup set to.Valid Values:'RemoveOldestInstance', 'RemoveNewestInstance'.",
			},

			"subnet_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "balanced-distribution",
				ValidateFunc: validateKecScalingGroupSubnetStrategy,
				Description:  "The Subnet Strategy of the desired ScalingGroup set to.Valid Values:'balanced-distribution', 'choice-first'.",
			},

			"subnet_id_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "The Subnet ID Set of the desired ScalingGroup set to.",
			},

			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_id_set"},
				Description:   "The Security Group ID of the desired ScalingGroup set to.",
			},

			"security_group_id_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				ConflictsWith: []string{"security_group_id"},
				Description:   "The Security Group ID List of the desired ScalingGroup set to.",
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Active",
				ValidateFunc: validateKecScalingGroupStatus,
				Description:  "The Status of the desired ScalingGroup.Valid Values:'Active', 'UnActive'.",
			},

			"slb_config_set": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of slb configs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slb_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The SLB ID of the desired ScalingGroup set to.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The Listener ID of the desired ScalingGroup set to.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     20,
							Description: "The weight of the desired ScalingGroup set to.Valid Values 1-100.",
						},
						"health_check_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "slb",
							ValidateFunc: validation.StringInSlice([]string{
								"slb",
								"kec",
							}, false),
							Description: "Health check type, valid values:'slb','kec'.",
						},
						"server_port_set": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Set:         schema.HashInt,
							Description: "The Server Port Set of the desired ScalingGroup set to.Valid Values 1-65535.",
						},
					},
				},
			},

			"scaling_configuration_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Scaling Configuration Name of the desired ScalingGroup set to.",
			},

			"instance_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The KEC instance Number of the desired ScalingGroup set to.Valid Value 0-10.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC.",
			},
		},
	}
}

func resourceKsyunScalingGroupExtra(d *schema.ResourceData, forceGet bool) map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	var r map[string]SdkReqTransform

	r = map[string]SdkReqTransform{
		"subnet_id_set":         {mapping: "SubnetId", Type: TransformWithN},
		"security_group_id_set": {mapping: "SecurityGroupId", Type: TransformWithN},
	}
	extra = SdkRequestAutoExtra(r, d, forceGet)
	extra["slb_config_set"] = SdkRequestMapping{
		Field: "Slb.",
		FieldReqFunc: func(item interface{}, s string, source string, m *map[string]interface{}) error {
			if arr, ok := item.([]interface{}); ok {
				for i, value := range arr {
					if d, ok := value.(map[string]interface{}); ok {
						for k, v := range d {
							if k == "slb_id" {
								(*m)[s+strconv.Itoa(i+1)+".Id"] = v
							}
							if k == "listener_id" {
								(*m)[s+strconv.Itoa(i+1)+".ListenerId"] = v
							}
							if k == "weight" {
								(*m)[s+strconv.Itoa(i+1)+".Weight"] = v
							}
							if k == "health_check_type" {
								(*m)[s+strconv.Itoa(i+1)+".HealthCheckType"] = v
							}
							if k == "server_port_set" {
								if x, ok := v.(*schema.Set); ok {
									for j, v1 := range (*x).List() {
										(*m)[s+strconv.Itoa(i+1)+".ServerPort."+strconv.Itoa(j+1)] = v1
									}
								}
							}
						}
					}
				}
			}
			return nil
		},
	}
	return extra
}

func resourceKsyunScalingGroupReqModify(req *map[string]interface{}, update bool) error {
	//sg
	v1, sg := (*req)["SecurityGroupId"]
	v2, sgn := (*req)["SecurityGroupId.1"]

	if !sg && !sgn && !update {
		return fmt.Errorf("you must set security_group_id or security_group_id_set")
	} else if sg && sgn {
		if v1 != v2 {
			return fmt.Errorf("security_group_id must equal security_group_id_set#0")
		}
		delete(*req, "SecurityGroupId")
	}
	return nil
}

func resourceKsyunScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingGroup()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, resourceKsyunScalingGroupExtra(d, false))
	if err != nil {
		return fmt.Errorf("error on creating ScalingGroup, %s", err)
	}

	//zero process
	if _, ok := req["MinSize"]; !ok {
		req["MinSize"] = 0
	}
	if _, ok := req["MaxSize"]; !ok {
		req["MaxSize"] = 0
	}
	if _, ok := req["DesiredCapacity"]; !ok {
		req["DesiredCapacity"] = 0
	}

	err = resourceKsyunScalingGroupReqModify(&req, false)
	if err != nil {
		return fmt.Errorf("error on creating ScalingGroup, %s", err)
	}

	action := "CreateScalingGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateScalingGroup(&req)
	if err != nil {
		return fmt.Errorf("error on creating ScalingGroup, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ReturnSet"].(map[string]interface{})["ScalingGroupId"].(string))
	}
	//set status
	if v, ok := d.GetOk("status"); ok {
		if v == "UnActive" {
			req = make(map[string]interface{})
			req["ScalingGroupId"] = d.Id()
			_, err = conn.DisableScalingGroup(&req)
			if err != nil {
				return fmt.Errorf("error on creating ScalingGroup, %s", err)
			}
		}
	}
	return resourceKsyunScalingGroupRead(d, meta)
}

func resourceKsyunScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingGroup()
	var action string

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, resourceKsyunScalingGroupExtra(d, true))
	if err != nil {
		return fmt.Errorf("error on modifying ScalingGroup, %s", err)
	}

	err = resourceKsyunScalingGroupReqModify(&req, true)
	if err != nil {
		return fmt.Errorf("error on modifying ScalingGroup, %s", err)
	}

	// distinguish modify lb info or other info
	reqLb := make(map[string]interface{})
	reqLb["ScalingGroupId"] = d.Id()
	for k, v := range req {
		if strings.HasPrefix(k, "Slb.") {
			reqLb[k] = v
			delete(req, k)
		}
	}
	action = "ModifyScalingLoadBalancers"
	logger.Debug(logger.ReqFormat, action, reqLb)
	_, err = conn.ModifyScalingLoadBalancers(&reqLb)
	if err != nil {
		return fmt.Errorf("error on modifying ScalingGroup, %s", err)
	}

	if len(req) > 0 {
		req1 := make(map[string]interface{})
		req1["ScalingGroupId"] = d.Id()
		if v, ok := req["Status"]; ok {
			if v == "Active" {
				action = "EnableScalingGroup"
				logger.Debug(logger.ReqFormat, action, req)
				_, err = conn.EnableScalingGroup(&req1)
				if err != nil {
					return fmt.Errorf("error on modifying ScalingGroup, %s", err)
				}
			} else {
				action = "DisableScalingGroup"
				logger.Debug(logger.ReqFormat, action, req)
				_, err = conn.DisableScalingGroup(&req1)
				if err != nil {
					return fmt.Errorf("error on modifying ScalingGroup, %s", err)
				}
			}
			delete(req, "Status")
		}
		if len(req) > 0 {
			req["ScalingGroupId"] = d.Id()
			action = "ModifyScalingGroup"
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.ModifyScalingGroup(&req)
			if err != nil {
				return fmt.Errorf("error on modifying ScalingGroup, %s", err)
			}
		}

	}
	return resourceKsyunScalingGroupRead(d, meta)
}

func resourceKsyunScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	req := make(map[string]interface{})
	req["ScalingGroupId.1"] = d.Id()
	action := "DescribeScalingGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeScalingGroup(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroup %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingGroupSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}

		SdkResponseAutoResourceData(d, resourceKsyunScalingGroup(), items[0], nil)
	}
	return nil
}

func resourceKsyunScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	req := make(map[string]interface{})
	action := "DeleteScalingGroup"
	//before delete need set DesiredCapacity=0 to release instance
	req["ScalingGroupId"] = d.Id()
	req["DesiredCapacity"] = 0
	req["MinSize"] = 0
	_, err := conn.ModifyScalingGroup(&req)
	if err != nil {
		return fmt.Errorf("error on deleting ScalingGroup, %s", err)
	}
	for k := range req {
		delete(req, k)
	}
	req["ScalingGroupId.1"] = d.Id()

	return resource.Retry(60*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		_, err1 := conn.DeleteScalingGroup(&req)
		if err1 == nil {
			return nil
		} else if notFoundErrorNew(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting ScalingGroup %q, %s", d.Id(), err1))
		}
	})

}
