/*
Provides an KcrsVpcAttachment resource.

Example Usage

```hcl
# Create a KcrsVpcAttachment
resource "ksyun_KcrsVpcAttachment" "default" {
  link_type = "DDoS_BGP"
  ip_count = 10
  band = 30
  max_band = 30
  idc_band = 100
  duration = 1
  KcrsVpcAttachment_name = "ksc_KcrsService"
  bill_type = 1
  service_id = "KcrsVpcAttachment_30G"
  project_id="0"
}
```

Import

KcrsVpcAttachment can be imported using the id, e.g.

```
$ terraform import ksyun_KcrsVpcAttachment.default KcrsService67b91d3c-c363-4f57-b0cd-xxxxxxxxxxxx
```
*/

package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKcrsVpcAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKcrsVpcAttachmentCreate,
		Read:   resourceKsyunKcrsVpcAttachmentRead,
		Update: resourceKsyunKcrsVpcAttachmentUpdate,
		Delete: resourceKsyunKcrsVpcAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: commonImport(2, "instance_id", "vpc_id"),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the project.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reserve_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enable_vpc_domain_dns": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// compute values
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_parse_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eni_lb_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_endpoint_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunKcrsVpcAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		req        = make(map[string]interface{})
		instanceId = d.Get("instance_id").(string)
		vpcId      = d.Get("vpc_id").(string)
	)

	kcrsVpcAttachmentService := KcrsService{meta.(*KsyunClient)}

	req["InstanceId"] = instanceId
	req["VpcId"] = vpcId
	req["ReserveSubnetId"] = d.Get("reserve_subnet_id")

	conn := kcrsVpcAttachmentService.client.kcrsconn
	if _, err := conn.CreateInternalEndpoint(&req); err != nil {
		return fmt.Errorf("an error caused when creating endpoint %s", err)
	}
	id := AssembleIds(instanceId, vpcId)
	d.SetId(id)

	if d.Get("enable_vpc_domain_dns").(bool) {
		var (
			params  = make(map[string]interface{})
			eniLBIp string

			endpointDns = "PrivateDomain"
		)

		internalEndpoints, err := kcrsVpcAttachmentService.ReadInternalEndpoint(d, instanceId)
		if err != nil {
			return fmt.Errorf("an error caused when getting access ip of instance %q, %s", d.Get("instance_id"), err)
		}

		if v, ok := internalEndpoints["EniLBIp"]; ok {
			eniLBIp = v.(string)
		}
		if eniLBIp == "" {
			return fmt.Errorf("an error caused when getting access ip of instance %q", d.Get("instance_id"))
		}

		params["InstanceId"] = instanceId
		params["VpcId"] = vpcId
		params["InternalEndpointDns"] = endpointDns
		params["EniLBIp"] = eniLBIp

		if _, err := conn.CreateInternalEndpointDns(&params); err != nil {
			return fmt.Errorf("an error caused when creating internal endpoint dns of instance %q, %s", d.Get("instance_id"), err)
		}

		_ = d.Set("internal_endpoint_dns", endpointDns)
		_ = d.Set("eni_lb_ip", eniLBIp)
	}

	return resourceKsyunKcrsVpcAttachmentRead(d, meta)
}

func resourceKsyunKcrsVpcAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	kcrsVpcAttachmentService := KcrsService{meta.(*KsyunClient)}
	err = kcrsVpcAttachmentService.ReadAndSetInternalEndpoint(d, resourceKsyunKcrsVpcAttachment())
	if err != nil {
		return fmt.Errorf("error on reading kcrs VpcAttachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKcrsVpcAttachmentUpdate(d *schema.ResourceData, meta interface{}) (err error) {

	var (
		client  = meta.(*KsyunClient)
		kcrsSrv = KcrsService{client: client}
	)

	if d.HasChange("enable_vpc_domain_dns") {
		if err := kcrsSrv.ModifyInternalEndpointDns(d); err != nil {
			return fmt.Errorf("an error caused when modifying internal vpc domain dns %q, %s", d.Id(), err)
		}
	}

	return resourceKsyunKcrsVpcAttachmentRead(d, meta)
}

func resourceKsyunKcrsVpcAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {

	var (
		req                      = make(map[string]interface{})
		kcrsVpcAttachmentService = KcrsService{meta.(*KsyunClient)}
	)
	ids := DisassembleIds(d.Id())
	req["InstanceId"] = ids[0]
	req["VpcId"] = ids[1]
	req["EniLBIp"] = d.Get("eni_lb_ip")

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		conn := kcrsVpcAttachmentService.client.kcrsconn
		_, err := conn.DeleteInternalEndpoint(&req)
		if err != nil {
			if _, readErr := kcrsVpcAttachmentService.ReadInternalEndpoint(d, ids[0]); err != nil && notFoundError(readErr) {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error on deleting kcrs internal endpoint %q, %s", d.Id(), err)
	}
	return err

}
