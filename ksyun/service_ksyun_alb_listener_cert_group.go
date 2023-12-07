package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type AlbListenerCertGroupService struct {
	client *KsyunClient
}

func (s *AlbListenerCertGroupService) createCertGroupCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"alb_listener_cert_set": {
			Ignore: true,
		},
	}
	var req map[string]interface{}
	req, err = SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateAlbListenerCertGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateAlbListenerCertGroup(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var id interface{}
			id, err = getSdkValue("AlbListenerCertGroup.AlbListenerCertGroupId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return d.Set("alb_listener_cert_group_id", d.Id())
		},
	}
	return
}
func (s *AlbListenerCertGroupService) CreateCertGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	var callbacks []ApiCall
	var createCall ApiCall
	createCall, err = s.createCertGroupCall(d, r)
	if err != nil {
		return err
	}
	callbacks = append(callbacks, createCall)

	// if certList, ok := d.GetOk("alb_listener_cert_set"); ok {
	// for _, item := range certList.([]interface{}) {
	//	logger.Debug(logger.RespFormat, "alb_listener_cert_set", item)
	// }
	certCalls, err := s.modifyCertSetCall(d, r, false)
	if err != nil {
		return
	}
	callbacks = append(callbacks, certCalls...)
	// }

	err = ksyunApiCallNew(callbacks, d, s.client, false)
	return
}

func (s *AlbListenerCertGroupService) readCertGroups(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(condition, "MaxResults", "NextToken", 200, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.slbconn
		action := "DescribeAlbListenerCertGroups"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeAlbListenerCertGroups(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeAlbListenerCertGroups(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("AlbListenerCertGroupSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *AlbListenerCertGroupService) readCertGroup(d *schema.ResourceData, certGroupId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if certGroupId == "" {
		certGroupId = d.Id()
	}
	req := map[string]interface{}{
		"AlbListenerCertGroupId.1": certGroupId,
	}
	results, err = s.readCertGroups(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("ALb listener cert group %s not exist ", certGroupId)
	}
	return
}

func (s *AlbListenerCertGroupService) ReadAndSetCertGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "AlbListenerCertGroupId",
			Type:    TransformWithN,
		},
		"alb_listener_id": {
			mapping: "alblistener-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readCertGroups(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection: data,
		// nameField:   "AlbListenerName",
		idFiled:     "AlbListenerCertGroupId",
		targetField: "listener_cert_groups",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *AlbListenerCertGroupService) ReadAndSetCertGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	var data map[string]interface{}
	data, err = s.readCertGroup(d, "")
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"AlbListenerCertSet": {
			Field: "certificate",
			FieldRespFunc: func(i interface{}) interface{} {
				var result []map[string]interface{}
				result = []map[string]interface{}{}
				v := i.([]interface{})
				for _, certItem := range v {
					d := certItem.(map[string]interface{})
					r := make(map[string]interface{})
					r["certificate_id"] = d["CertificateId"]
					r["certificate_name"] = d["CertificateName"]
					r["cert_authority"] = d["CertAuthority"]
					r["common_name"] = d["CommonName"]
					r["expire_time"] = d["ExpireTime"]
					result = append(result, r)
				}
				logger.Debug(logger.RespFormat, "AlbListenerCertSetMap", result, i)
				return result
			},
		},
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return
}

func (s *AlbListenerCertGroupService) removeCertGroupCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"AlbListenerCertGroupId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteAlbListenerCertGroup",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteAlbListenerCertGroup(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.readCertGroup(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on reading alb listener cert group when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return
}
func (s *AlbListenerCertGroupService) RemoveCertGroup(d *schema.ResourceData) (err error) {
	var call ApiCall
	call, err = s.removeCertGroupCall(d)
	if err != nil {
		return err
	}
	err = ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
	return
}

func (s *AlbListenerCertGroupService) ModifyCertGroup(d *schema.ResourceData, r *schema.Resource) (err error) {
	if d.HasChange("certificate") {
		var callbacks []ApiCall

		callbacks, err = s.modifyCertSetCall(d, r, true)
		if err != nil {
			return
		}

		err = ksyunApiCallNew(callbacks, d, s.client, false)
	}
	return
}

func (s *AlbListenerCertGroupService) modifyCertSetCall(d *schema.ResourceData, r *schema.Resource, isUpdate bool) (callbacks []ApiCall, err error) {

	newCertIds := make([]string, 0)
	removeCertIds := make([]string, 0)

	if isUpdate {
		oldList, newList := d.GetChange("certificate")

		for _, newCert := range newList.([]interface{}) {
			newCertId := newCert.(map[string]interface{})["certificate_id"].(string)
			hasIt := func() bool {
				for _, oldCert := range oldList.([]interface{}) {
					if newCertId == oldCert.(map[string]interface{})["certificate_id"].(string) {
						return true
					}
				}
				return false
			}()
			if !hasIt {
				newCertIds = append(newCertIds, newCertId)
			}
		}

		for _, oldCert := range oldList.([]interface{}) {
			oldCertId := oldCert.(map[string]interface{})["certificate_id"].(string)
			hasIt := func() bool {
				for _, newCert := range newList.([]interface{}) {
					if oldCertId == newCert.(map[string]interface{})["certificate_id"].(string) {
						return true
					}
				}
				return false
			}()
			if !hasIt {
				removeCertIds = append(removeCertIds, oldCertId)
			}
		}
		logger.Debug(logger.RespFormat, "modifyCertSetCall new:", newCertIds)
		logger.Debug(logger.RespFormat, "modifyCertSetCall remove:", removeCertIds)
		// return
	} else {
		if certSets, ok := d.GetOk("certificate"); ok {
			for _, cert := range certSets.([]interface{}) {
				newCertIds = append(newCertIds, cert.(map[string]interface{})["certificate_id"].(string))
			}
		}
	}
	logger.Debug(logger.RespFormat, "modifyCertSetCallList", newCertIds, removeCertIds)
	logger.Debug(logger.RespFormat, "modifyCertSetCall", d.Id(), d.Get("alb_listener_cert_group_id"))
	for _, removeCertId := range removeCertIds {
		callbacks = append(callbacks, ApiCall{
			param: &map[string]interface{}{
				// "AlbListenerCertGroupId": d.Id(),
				"CertificateId": removeCertId,
			},
			action: "DissociateCertificateWithGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				(*call.param)["AlbListenerCertGroupId"] = d.Id()
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.DissociateCertificateWithGroup(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return
			},
		})
	}

	for _, newCertId := range newCertIds {
		callbacks = append(callbacks, ApiCall{
			param: &map[string]interface{}{
				// "AlbListenerCertGroupId": d.Id(),
				"CertificateId": newCertId,
			},
			action: "AssociateCertificateWithGroup",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				(*call.param)["AlbListenerCertGroupId"] = d.Id()
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.AssociateCertificateWithGroup(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return
			},
		})
	}

	return
}
