package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type LocalVolumeService struct {
	client *KsyunClient
}

func (s *LocalVolumeService) ReadAndSetLocalSnapshots(d *schema.ResourceData, r *schema.Resource) (err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		data    []interface{}
		list    []interface{}
	)
	condition := map[string]interface{}{}
	if v, ok := d.GetOk("local_volume_name"); ok {
		condition["LocalVolumeName"] = v
	}
	if v, ok := d.GetOk("source_local_volume_id"); ok {
		condition["SourceLocalVolumeId"] = v
	}
	if v, ok := d.GetOk("local_volume_snapshot_id"); ok {
		condition["LocalVolumeSnapshotId"] = v
	}

	list, err = pageQuery(condition, "MaxResults", "Marker", 200, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kecconn
		action := "DescribeLocalVolumeSnapshots"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeLocalVolumeSnapshots(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeLocalVolumeSnapshots(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("LocalVolumeSnapshotSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  list,
		idFiled:     "LocalVolumeSnapshotId",
		nameField:   "LocalVolumeSnapshotName",
		targetField: "local_snapshot_set",
		//extra: map[string]SdkResponseMapping{
		//	"LocalVolumeId": {
		//		Field:    "id",
		//		KeepAuto: true,
		//	},
		//},
	})
}

func (s *LocalVolumeService) ReadAndSetLocalVolumes(d *schema.ResourceData, r *schema.Resource) (err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		data    []interface{}
	)

	condition := map[string]interface{}{}
	if v, ok := d.GetOk("instance_name"); ok {
		condition["InstanceName"] = v
	}
	list, err := pageQuery(condition, "MaxResults", "Marker", 200, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.client.kecconn
		action := "DescribeLocalVolumes"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeLocalVolumes(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeLocalVolumes(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("LocalVolumeSet", *resp)
		if err != nil {
			return data, err
		}
		data = results.([]interface{})
		return data, err
	})
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  list,
		idFiled:     "LocalVolumeId",
		nameField:   "LocalVolumeName",
		targetField: "local_volume_set",
		//extra: map[string]SdkResponseMapping{
		//	"LocalVolumeId": {
		//		Field:    "id",
		//		KeepAuto: true,
		//	},
		//},
	})
}
