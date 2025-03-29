package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type KpfsService struct {
	client *KsyunClient
}

// start kpfs

func (s *KpfsService) readPerformanceOnePosixAcl(d *schema.ResourceData, r *schema.Resource) (err error) {

	aclData, err := s.getPerformanceOnePosixAclList(d, r)
	if err != nil {
		return err
	}
	data, err := getSdkValue("Data", *aclData)
	if err != nil {
		return err
	}
	datas := data.([]interface{})
	if len(datas) == 0 {
		return fmt.Errorf("error on add acl: acl is empty")
	}

	firstData := datas[0].(map[string]interface{})

	// extra := map[string]SdkResponseMapping{
	// 	"Ips": {
	// 		Field: "ips",
	// 		FieldRespFunc: func(i interface{}) interface{} {
	// 			v, _ := strconv.Atoi(i.(string))
	// 			return v
	// 		},
	// 	},
	// }

	SdkResponseAutoResourceData(d, r, firstData, nil)

	return

}

func (s *KpfsService) addPosixAclIp(d *schema.ResourceData, r *schema.Resource) (err error) {
	sorceIp, err := s.getStorageIp(d, r)
	if err != nil {
		return err
	}
	transform := map[string]SdkReqTransform{
		"kpfs_acl_id": {mapping: "PosixAclId"},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return err
	}
	req["Ip"] = sorceIp
	//追加 acl 的ip属性
	return s.addPerformanceOnePosixAclIp(req)
}

func (s *KpfsService) deletePerformanceOnePosixAclIp(d *schema.ResourceData, r *schema.Resource) (err error) {
	sorceIp, err := s.getStorageIp(d, r)
	if err != nil {
		return err
	}
	transform := map[string]SdkReqTransform{
		"kpfs_acl_id": {mapping: "PosixAclId"},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return err
	}
	req["Ip"] = sorceIp

	return s.removePerformanceOnePosixAclIp(req)
}

func (s *KpfsService) removePerformanceOnePosixAclIp(aclData map[string]interface{}) error {
	kpfsconn := s.client.kpfsconn
	_, err := kpfsconn.RemovePerformanceOnePosixAclIp(&aclData)
	if err != nil {
		return fmt.Errorf("failed to delete acl ip: %v", err)
	}
	return nil
}

// 定义 updatePerformanceOnePosixAclIps 方法
func (s *KpfsService) addPerformanceOnePosixAclIp(aclData map[string]interface{}) error {
	// 调用 SDK 更新 ACL 的 IP 列表
	kpfsconn := s.client.kpfsconn
	_, err := kpfsconn.AddPerformanceOnePosixAclIp(&aclData)
	if err != nil {
		return fmt.Errorf("failed to update AddPerformanceOnePosixAclIp: %v", err)
	}
	return nil
}

func (s *KpfsService) getPerformanceOnePosixAclList(d *schema.ResourceData, r *schema.Resource) (*map[string]interface{}, error) {
	transform := map[string]SdkReqTransform{
		"kpfs_acl_id": {mapping: "PosixAclId"},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)
	if err != nil {
		return nil, err
	}
	kpfsconn := s.client.kpfsconn
	return kpfsconn.DescribePerformanceOnePosixAclList(&req)
}

func (s *KpfsService) getStorageIp(d *schema.ResourceData, r *schema.Resource) (string, error) {
	transform := map[string]SdkReqTransform{
		"epc_id": {mapping: "HostId.1"},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil)

	epcconn := s.client.epcconn
	resp, err := epcconn.DescribeEpcs(&req)

	results, err := getSdkValue("HostSet", *resp)
	if err != nil {
		return "", err
	}
	data := results.([]interface{})

	if len(data) == 0 {
		return "", fmt.Errorf("epc %s storage ip  not exist ", req["HostId.1"])
	}
	for _, v := range data {
		subData := v.(map[string]interface{})
		roces := subData["Roces"].([]interface{})
		for _, roce := range roces {
			roceData := roce.(map[string]interface{})
			if roceData["Type"] == "storage" {
				return roceData["Ip"].(string), nil
			}
		}
	}
	return "", fmt.Errorf("epc %s storage ip not exist ", req["HostId.1"])
}
