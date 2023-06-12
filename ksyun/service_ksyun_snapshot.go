package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SnapshotService struct {
	client *KsyunClient
}

func (s *SnapshotService) readSnapshots(condition map[string]interface{}) (data []interface{}, err error) {

	var resp *map[string]interface{}
	var result interface{}
	var hasNext interface{}

	condition["PageNumber"] = 1
	condition["PageSize"] = 10
	//logger.Debug("readSnapshots", "start", 11111)
	for {
		//logger.Debug("readSnapshots", "condition", condition)
		resp, err = s.client.ebsconn.DescribeSnapshots(&condition)
		//logger.Debug("readSnapshots", "resp", resp)
		if err != nil {
			return
		}
		result, err = getSdkValue("Snapshots", *resp)
		if list, ok := result.([]interface{}); ok {
			data = append(data, list...)
		}
		hasNext, err = getSdkValue("Page.hasNext", *resp)
		//logger.Debug("readSnapshots", "hasNext", hasNext)
		if err != nil {
			return
		}
		if hasNextValue, ok := hasNext.(bool); ok && hasNextValue {
			condition["PageNumber"] = condition["PageNumber"].(int) + 1
		} else {
			return
		}
	}

}

func (s *SnapshotService) ReadAndSetSnapshots(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	var list []interface{}
	list, err = s.readSnapshots(req)
	if err != nil {
		return
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  list,
		idFiled:     "SnapshotId",
		nameField:   "SnapshotName",
		targetField: "snapshots",
	})
}
