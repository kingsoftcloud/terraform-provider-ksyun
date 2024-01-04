/*
Provides an internal access attachment resource with a vpc under kcrs repository instance.

Example Usage

```hcl
# repository instance
resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfunittest"
	instance_type = "basic"
}

# To attach a vpc for an instance
resource "ksyun_kcrs_vpc_attachment" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	vpc_id = "vpc_id"
	reserve_subnet_id = "subnet_id"
	enable_vpc_domain_dns = true
}
```

Import

KcrsVpcAttachment can be imported using `instance_id:vpc_id`, e.g.

```
$ terraform import ksyun_kcrs_vpc_attachment.foo ${instance_id}:${vpc_id}
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
				Description: "Instance id of repository.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Vpc id.",
			},
			"reserve_subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of subnet type is '**Reserve**'.",
			},

			"enable_vpc_domain_dns": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable vpc domain dns. Default value is `false`.",
			},

			// compute values
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the internal access.",
			},
			"dns_parse_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the DNS parsed.",
			},
			"eni_lb_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the internal access.",
			},
			"internal_endpoint_dns": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint Domain of the internal access.",
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
