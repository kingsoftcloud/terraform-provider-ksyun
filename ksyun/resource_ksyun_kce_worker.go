package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

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
				Description: "The ID of the kce cluster",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the kec instance. The instance will be shut down while being added to the kce cluster.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the image which support KCE.",
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
	return
}
func resourceKsyunKceWorkerRead(d *schema.ResourceData, meta interface{}) (err error) {
	return
}
func resourceKsyunKceWorkerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return
}
