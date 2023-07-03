package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

type SnapshotService struct {
	client *KsyunClient
}

func (s *SnapshotService) CreateSnapshot(d *schema.ResourceData, r *schema.Resource) (err error) {

	transform := map[string]SdkReqTransform{}
	var req map[string]interface{}
	req, err = SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return
	}

	var resp *map[string]interface{}
	resp, err = s.client.ebsconn.CreateSnapshot(&req)
	if err != nil {
		return
	}
	id, err := getSdkValue("SnapshotId", *resp)
	if err != nil {
		return err
	}
	d.SetId(id.(string))
	_, err = s.checkSnapshotState(d, "", []string{"available"}, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return
}

func (s *SnapshotService) checkSnapshotState(d *schema.ResourceData, snapshotId string, target []string, timeout time.Duration) (state interface{}, err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    s.snapshotStateRefreshFunc(d, snapshotId, []string{"error"}),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	return stateConf.WaitForState()
}

func (s *SnapshotService) snapshotStateRefreshFunc(d *schema.ResourceData, snapshotId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := s.ReadSnapshot(d, snapshotId)
		if err != nil {
			return nil, "", err
		}

		status, err := getSdkValue("SnapshotStatus", data)
		if err != nil {
			return nil, "", err
		}

		for _, v := range failStates {
			if v == status.(string) {
				return nil, "", fmt.Errorf("snapshot status error, status:%v", status)
			}
		}
		return data, status.(string), nil
	}
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

func (s *SnapshotService) ModifySnapshot(d *schema.ResourceData, r *schema.Resource) (err error) {
	req := map[string]interface{}{
		"SnapshotId": d.Id(),
	}
	if d.HasChange("snapshot_name") {
		req["SnapshotName"] = d.Get("snapshot_name")
	}
	if d.HasChange("snapshot_desc") {
		req["SnapshotDesc"] = d.Get("snapshot_desc")
	}
	_, err = s.client.ebsconn.ModifySnapshot(&req)

	return
}

func (s *SnapshotService) Remove(d *schema.ResourceData) (err error) {
	removeReq := map[string]interface{}{
		"SnapshotId": d.Id(),
	}
	_, err = s.client.ebsconn.DeleteSnapshot(&removeReq)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadSnapshot(d, "")
		if callErr != nil {
			if notFoundError(callErr) {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on deleting snapshot %q, %s", d.Id(), callErr))
			}
		}
		if v, getErr := getSdkValue("SnapshotStatus", data); getErr == nil {
			if vStr, ok := v.(string); ok && vStr == "deleting" {
				return resource.RetryableError(errors.New("deleting"))
			}
		}
		return nil
	})
}

func (s *SnapshotService) ReadSnapshot(d *schema.ResourceData, snapshotId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if snapshotId == "" {
		snapshotId = d.Id()
	}
	req := map[string]interface{}{
		"SnapshotId": snapshotId,
	}

	results, err = s.readSnapshots(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("snapshot %s not exist ", snapshotId)
	}
	return data, err
}

func (s *SnapshotService) ReadAndSetSnapshot(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadSnapshot(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on reading snapshot %q, %s", d.Id(), callErr))
			}
		} else {
			SdkResponseAutoResourceData(d, r, data, nil)
			return nil
		}
	})
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
