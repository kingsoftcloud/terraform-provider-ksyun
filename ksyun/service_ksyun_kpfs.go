package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
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

func (s *KpfsService) CreateFileSystem(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.createFileSystemCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, false)
}

func (s *KpfsService) createFileSystemCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateFileSystem",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			kpfsconn := client.kpfsconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = kpfsconn.CreateFileSystem(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			var fileSystemId interface{}
			if resp != nil {
				fileSystemId, err = getSdkValue("FileSystemId", *resp)
				if err != nil {
					return err
				}
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				if err := s.checkFileSystemState(d, fileSystemId.(string), []string{"using"}, d.Timeout(schema.TimeoutCreate)); err != nil {
					return fmt.Errorf("waiting for kpfs fileSystem create caused an error: %s", err)
				}
				// set id
				d.SetId(fileSystemId.(string))
			}
			return err
		},
	}
	return callback, err
}

func (s *KpfsService) checkFileSystemState(d *schema.ResourceData, fileSystemId string, target []string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    s.fileSystemStateRefreshFunc(d, fileSystemId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func (s *KpfsService) fileSystemStateRefreshFunc(d *schema.ResourceData, fileSystemId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var err error
		resp, err := s.ReadKpfsFileSystem(d, fileSystemId)
		if err != nil {
			return nil, "", fmt.Errorf("Error DescribeFileSystemList : %s ", err)
		}
		l, err1 := getSdkValue("Status", resp)
		if err1 != nil {
			return nil, "", fmt.Errorf("Error DescribeFileSystemList : %s ", err)
		}
		return resp, l.(string), nil
	}
}

func (s *KpfsService) ReadKpfsFileSystem(d *schema.ResourceData, fileSystemId string) (data map[string]interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	if fileSystemId == "" {
		fileSystemId = d.Id()
	}
	req := map[string]interface{}{
		"FileSystemId": fileSystemId,
	}
	conn := s.client.kpfsconn
	resp, err = conn.DescribeFileSystemList(&req)
	if err != nil {
		return nil, err
	}
	results, err = getSdkValue("Data", *resp)
	if err != nil {
		return nil, err
	}

	resultList, ok := results.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: Data is not a list")
	}
	if len(resultList) == 0 {
		return nil, fmt.Errorf("kpfs fileSystem %s not exist ", fileSystemId)
	}

	// 取第一条数据作为结果返回
	firstItem, ok := resultList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected item format in Data list")
	}
	return firstItem, nil
}

func (s *KpfsService) ReadKpfsFileSystemList(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"id": {
			mapping: "FileSystemId",
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}

	data, err := s.ReadKpfsFs(req)
	if err != nil {
		return err
	}
	logger.Debug(logger.RespFormat, "ReadKpfsFs", "ReadKpfsFs", data)
	return mergeDataSourcesRespIdInObj(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "FileSystemInfo.FileSystemId",
		targetField: "data",
		nameField:   "FileSystemName",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *KpfsService) ReadKpfsFs(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	return pageQuery(condition, "PageSize", "PageNum", 1000, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kpfsconn
		action := "DescribeFileSystemList"
		logger.Debug(logger.ReqFormat, action, condition)
		resp, err = conn.DescribeFileSystemList(&condition)
		results, err = getSdkValue("Data", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
}

func (s *KpfsService) ReadKpfsFileSystemOne(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"id": {
			mapping: "FileSystemId",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}

	data, err := s.ReadKpfsFs(req)
	if err != nil {
		return err
	}
	logger.Debug(logger.RespFormat, "ReadKpfsFs", "ReadKpfsFs", data)
	SdkResponseAutoResourceData(d, r, data, nil)
	return
}

func (s *KpfsService) DeleteFileSystem(d *schema.ResourceData, meta interface{}) (err error) {
	conn := meta.(*KsyunClient).kpfsconn
	req := make(map[string]interface{})
	req["FileSystemId"] = d.Id()
	_, err = conn.DeleteFileSystem(&req)
	if err != nil {
		return fmt.Errorf("failed to delete fileSystem : %v", err)
	}
	return nil
}

func (s *KpfsService) ReadKpfsClusterList(d *schema.ResourceData, r *schema.Resource) (err error) {
	req, err := mergeDataSourcesReq(d, r, nil)
	if err != nil {
		return err
	}
	conn := s.client.kpfsconn
	resp, err := conn.DescribeClusterInfo(&req)
	if err != nil {
		logger.Debug(logger.RespFormat, "ReadKpfsClusterList------", req, *resp)
		return fmt.Errorf("failed to ReadKpfsClusterList : %v", err)
	}
	results, err := getSdkValue("Data", *resp)
	if err != nil {
		return fmt.Errorf("failed to ReadKpfsClusterList : %v", err)
	}
	var data = results.([]interface{})
	logger.Debug(logger.RespFormat, "ReadKpfsCluster", "ReadKpfsCluster", data)
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "data",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *KpfsService) ReadKpfsClientInstallInfo(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"id": {
			mapping: "FileSystemId",
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	conn := s.client.kpfsconn
	var resp, err1 = conn.DescribeClientInstallInfo(&req)
	if err1 != nil {
		logger.Debug(logger.RespFormat, "DescribeClientInstallInfo------", req, *resp)
		return fmt.Errorf("failed to DescribeClientInstallInfo : %v", err)
	}
	logger.Debug(logger.RespFormat, "ReadKpfsCluster", "DescribeClientInstallInfo", *resp)
	//将resp 变更为一个数组
	results, err := getSdkValue("Data", *resp)
	if err != nil {
		return fmt.Errorf("failed to DescribeClientInstallInfo : %v", err)
	}
	data := results.([]interface{})
	ip, _ := getSdkValue("ClusterDataIP", *resp)
	for _, item := range data {
		item.(map[string]interface{})["cluster_data_ip"] = ip
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		targetField: "data",
		extra:       map[string]SdkResponseMapping{},
	})
}
