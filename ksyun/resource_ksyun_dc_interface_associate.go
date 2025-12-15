/*
Associate an Direct connect gateway resource with interface.

# Example Usage

```hcl

resource "ksyun_dc_interface_associate" "test" {
  direct_connect_interface_id = ksyun_direct_connect_interface.test.id
  direct_connect_gateway_id   = ksyun_direct_connect_gateway.test.id
}

```
*/

package ksyun

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunDCInterfaceAssociate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceKsyunDCInterfaceAssociateCreate,
		Read:     resourceKsyunDCInterfaceAssociateRead,
		Delete:   resourceKsyunDCInterfaceAssociateDelete,
		Importer: nil,
		Schema: map[string]*schema.Schema{
			"direct_connect_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the direct connect interface to be associated with the gateway.",
			},
			"direct_connect_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Gateway ID of the direct connect gateway.",
			},
		},
	}
}

func resourceKsyunDCInterfaceAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	apiProcess := NewApiProcess(context.TODO(), d, srv.client, false)

	call, err := srv.attachDirectConnectGatewayCall(d)
	if err != nil {
		return fmt.Errorf("error on associating interface on gateway %q, %s", d.Id(), err)
	}

	apiProcess.PutCalls(call)
	err = apiProcess.Run()
	if err != nil {
		return fmt.Errorf("error on associating interface on gateway %q, %s", d.Id(), err)
	}
	d.SetId(d.Get("direct_connect_gateway_id").(string))
	return resourceKsyunDCInterfaceAssociateRead(d, meta)
}

func resourceKsyunDCInterfaceAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	var data map[string]interface{}
	data, err = srv.readDirectConnectGateway(d, d.Get("direct_connect_gateway_id").(string))
	if err != nil {
		return fmt.Errorf("error on reading association between gateway and interface %q, %s", d.Id(), err)
	}

	found := false
	dcinfoSet := data["DirectConnectInterfaceInfoSet"]
	switch dcinfoSet.(type) {
	case []interface{}:
		dcs := dcinfoSet.([]interface{})
		for _, dc := range dcs {
			dcc := dc.(map[string]interface{})
			inId := indirectString(dcc["DirectConnectInterfaceId"])
			d.Set("direct_connect_interface_id", inId)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("error on reading association between gateway and interface %q, %s", d.Id(), "no direct connect interface found")
	}

	d.Set("direct_connect_gateway_id", d.Id())
	return err
}

func resourceKsyunDCInterfaceAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	srv := VpcService{meta.(*KsyunClient)}
	apiProcess := NewApiProcess(context.TODO(), d, srv.client, false)
	detachCall, err := srv.detachDirectConnectGatewayCall(d)
	if err != nil {
		return fmt.Errorf("error on deleting direct connect gateway interface association %q, %s", d.Id(), err)
	}
	apiProcess.PutCalls(detachCall)
	return apiProcess.Run()
}
