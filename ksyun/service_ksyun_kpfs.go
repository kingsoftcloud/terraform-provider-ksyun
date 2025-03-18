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

func (s *KpfsService) updatePerformanceOnePosixAcl(d *schema.ResourceData, r *schema.Resource) (err error) {
	sorceIp, err := s.getStorageIp(d, r)
	if err != nil {
		return err
	}
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
		return fmt.Errorf("acl is empty")
	}

	firstData := datas[0].(map[string]interface{})

	ips := firstData["Ips"].([]interface{})
	// 修改 ips 列表
	hasAdd := false
	for _, v := range ips {
		if v.(string) == sorceIp {
			hasAdd = true
		}
	}
	if hasAdd {
		return nil
	}
	ips = append(ips, sorceIp)

	//修改aclData 的ips属性
	updateAclData := datas[0].(map[string]interface{}) // 断言为 map[string]interface{}
	updateAclData["Ips"] = ips
	return s.updatePerformanceOnePosixAclIps(updateAclData)
}

func (s *KpfsService) deletePerformanceOnePosixAclIp(d *schema.ResourceData, r *schema.Resource) (err error) {
	sorceIp, err := s.getStorageIp(d, r)
	if err != nil {
		return err
	}
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
		return nil
	}

	firstData := datas[0].(map[string]interface{})
	ips := firstData["Ips"].([]interface{})
	// 修改 ips 列表
	var newIps []interface{}
	hasRemove := false
	for _, v := range ips {
		ip := v.(string)
		if ip != sorceIp {
			newIps = append(newIps, ip)
		}
		if v.(string) == sorceIp {
			hasRemove = true
		}
	}

	if !hasRemove {
		return nil
	}

	//修改aclData 的ips属性
	updateAclData := datas[0].(map[string]interface{}) // 断言为 map[string]interface{}
	updateAclData["Ips"] = newIps

	if len(newIps) == 0 {
		return s.deletePerformanceOnePosixAcl(updateAclData)
	}
	return s.updatePerformanceOnePosixAclIps(updateAclData)
}

func (s *KpfsService) deletePerformanceOnePosixAcl(aclData map[string]interface{}) error {
	kpfsconn := s.client.kpfsconn
	_, err := kpfsconn.DeletePerformanceOnePosixAcl(&aclData)
	if err != nil {
		return fmt.Errorf("failed to delete acl ip: %v", err)
	}
	return nil
}

// 定义 updatePerformanceOnePosixAclIps 方法
func (s *KpfsService) updatePerformanceOnePosixAclIps(aclData map[string]interface{}) error {
	// 调用 SDK 更新 ACL 的 IP 列表
	kpfsconn := s.client.kpfsconn
	_, err := kpfsconn.UpdatePerformanceOnePosixAcl(&aclData)
	if err != nil {
		return fmt.Errorf("failed to update PerformanceOnePosixAcl: %v", err)
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
