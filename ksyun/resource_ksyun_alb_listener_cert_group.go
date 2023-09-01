/*
Provides a ALB Listener cert group resource.

# Example Usage

```hcl
# network and security group configuration
resource "ksyun_vpc" "example" {
  vpc_name   = "tf-alb-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "ksyun_subnet" "example" {
  subnet_name       = "tf-alb-subnet"
  cidr_block        = "10.0.1.0/24"
  subnet_type       = "Normal"
  availability_zone = var.az
  vpc_id            = ksyun_vpc.example.id
}

resource "ksyun_security_group" "example" {
  vpc_id              = ksyun_vpc.example.id
  security_group_name = "tf_sg"
}

resource "ksyun_security_group_entry" "example" {
  security_group_id = ksyun_security_group.example.id
  cidr_block        = "0.0.0.0/0"
  protocol          = "ip"
  direction         = "in"
}

# ---------------------------------------------
# resource ksyun alb
resource "ksyun_alb" "example" {
  alb_name    = "tf-alb-example-1"
  alb_version = "standard"
  alb_type    = "public"
  state       = "start"
  charge_type = "PrePaidByHourUsage"
  vpc_id      = ksyun_vpc.example.id
  project_id  = 0
  enabled_log = false
  ip_version  = "ipv4"
}

# query your certificates on ksyun
data "ksyun_certificates" "listener_cert" {
  name_regex = "test"
}

resource "ksyun_alb_listener" "example" {
  alb_id             = ksyun_alb.example.id
  alb_listener_name  = "alb-example-listener"
  protocol           = "HTTPS"
  port               = 8099
  alb_listener_state = "start"
  certificate_id     = data.ksyun_certificates.listener_cert.certificates.0.certificate_id
  http_protocol      = "HTTP1.1"
  session {
    cookie_type                = "ImplantCookie"
    cookie_name                = "KLBRSIDdad"
    session_state              = "start"
    session_persistence_period = 3100
  }
}

resource "ksyun_alb_listener_cert_group" "default" {
  alb_listener_id = ksyun_alb_listener.example.id
  certificate {
    certificate_id = data.ksyun_certificates.listener_cert.certificates.0.certificate_id
  }
}
```

# Import

ALB Listener Cert Group can be imported using the `id`, e.g.

```
$ terraform import ksyun_alb_listener_cert_group.example vserver-abcdefg
```
*/

package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunAlbListenerCertGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAlbListenerCertGroupCreate,
		Read:   resourceKsyunAlbListenerCertGroupRead,
		Update: resourceKsyunAlbListenerCertGroupUpdate,
		Delete: resourceKsyunAlbListenerCertGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"alb_listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the ALB Listener.",
			},
			"alb_listener_cert_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the ALB Listener Cert Group.",
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				// Computed:    true,
				Description: "The certificate included in the cert group.",
				// DiffSuppressFunc: func(k, oldV, newV string, d *schema.ResourceData) bool {
				//	logger.Debug(logger.ReqFormat, "certificateDiffSuppressFunc k", k)
				//	logger.Debug(logger.ReqFormat, "certificateDiffSuppressFunc old", oldV)
				//	logger.Debug(logger.ReqFormat, "certificateDiffSuppressFunc new", newV)
				//	if k == "certificate.#" {
				//		if oldV == "0" && newV == "0" {
				//			return true
				//		}
				//	}
				//	return false
				// },
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the certificate.",
						},
						"certificate_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the certificate.",
						},
						"cert_authority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "certificate authority.",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The common name on the certificate.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expire time of the certificate.",
						},
					},
				},
			},
		},
	}
}

func resourceKsyunAlbListenerCertGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.CreateCertGroup(d, resourceKsyunAlbListenerCertGroup())
	if err != nil {
		return fmt.Errorf("error on creating ALB listener cert group %q, %s", d.Id(), err)
	}
	err = resourceKsyunAlbListenerCertGroupRead(d, meta)
	return
}
func resourceKsyunAlbListenerCertGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.ReadAndSetCertGroup(d, resourceKsyunAlbListenerCertGroup())
	if err != nil {
		return fmt.Errorf("error on reading ALB listener cert group %q, %s", d.Id(), err)
	}
	return
}
func resourceKsyunAlbListenerCertGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.ModifyCertGroup(d, resourceKsyunAlbListenerCertGroup())
	if err != nil {
		return fmt.Errorf("error on updating listener cert group %q, %s", d.Id(), err)
	}
	err = resourceKsyunAlbListenerCertGroupRead(d, meta)
	return
}
func resourceKsyunAlbListenerCertGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	s := AlbListenerCertGroupService{meta.(*KsyunClient)}
	err = s.RemoveCertGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener cert group %q, %s", d.Id(), err)
	}
	return
}
