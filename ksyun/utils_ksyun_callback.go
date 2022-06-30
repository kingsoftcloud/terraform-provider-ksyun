package ksyun

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ApiCall struct {
	param         *map[string]interface{}
	action        string
	beforeCall    beforeCallFunc
	executeCall   executeCallFunc
	callError     callErrorFunc
	afterCall     afterCallFunc
	disableDryRun bool
}

type ksyunApiCallFunc func(d *schema.ResourceData, meta interface{}) error
type callErrorFunc func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error
type executeCallFunc func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error)
type afterCallFunc func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error
type beforeCallFunc func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (bool, error)

func ksyunApiCall(api []ksyunApiCallFunc, d *schema.ResourceData, meta interface{}) (err error) {
	if api != nil {
		for _, f := range api {
			if f != nil {
				err = f(d, meta)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

func ksyunApiCallNew(api []ApiCall, d *schema.ResourceData, client *KsyunClient, isDryRun bool) (err error) {
	if !client.dryRun || !isDryRun {
		return ksyunApiCallProcess(api, d, client, false)
	} else {
		err = ksyunApiCallProcess(api, d, client, true)
		if err != nil {
			return err
		}
		return ksyunApiCallProcess(api, d, client, false)
	}
}

func ksyunApiCallProcess(api []ApiCall, d *schema.ResourceData, client *KsyunClient, isDryRun bool) (err error) {
	if api != nil {
		for _, f := range api {
			if f.executeCall != nil {
				var (
					resp *map[string]interface{}
				)
				doExecute := true
				if isDryRun {
					if f.disableDryRun {
						continue
					}
					(*(f.param))["DryRun"] = true
				} else if f.beforeCall != nil {
					doExecute, err = f.beforeCall(d, client, f)
				}
				if doExecute || isDryRun {
					resp, err = f.executeCall(d, client, f)
				}
				if isDryRun {
					delete(*(f.param), "DryRun")
					if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 412 {
						err = nil
					}
				} else {
					if err != nil {
						if f.callError == nil {
							return err
						} else {
							err = f.callError(d, client, f, err)
						}
					}
					if err != nil {
						return err
					}
					if doExecute && f.afterCall != nil {
						err = f.afterCall(d, client, resp, f)
					}
				}
			}
			if err != nil {
				return err
			}
		}
	}
	return err
}
