package ksyun

import (
	"fmt"
	"time"

	"github.com/KscSDK/ksc-sdk-go/service/ebs"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type AutoSnapshotSrv struct {
	client *KsyunClient
}

func NewAutoSnapshotSrv(client *KsyunClient) AutoSnapshotSrv {
	return AutoSnapshotSrv{
		client: client,
	}
}

// querySnapshotPolicyByID will query snapshot policy from ksyun open-api
func (s *AutoSnapshotSrv) querySnapshotPolicyByID(input map[string]interface{}) ([]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)

	resp, err = s.GetConn().DescribeAutoSnapshotPolicy(&input)

	if err != nil {
		return nil, fmt.Errorf("error describing auto snapshot policy: %s", err)
	}

	results, err := getSdkValue("AutoSnapshotPolicySet", *resp)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	return If2Slice(results)
}

func (s *AutoSnapshotSrv) createAutoSnapshotPolicy(input map[string]interface{}) (string, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = s.GetConn().CreateAutoSnapshotPolicy(&input)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error create auto snapshot policy: %s", err)
	}

	results, err := getSdkValue("AutoSnapshotPolicyId", *resp)
	if err != nil {
		return "", fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	return If2String(results)
}

func (s *AutoSnapshotSrv) deleteAutoSnapshotPolicy(input map[string]interface{}) ([]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = s.GetConn().DeleteAutoSnapshotPolicy(&input)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error delete auto snapshot policy: %s", err)
	}

	retSet, err := getSdkValue("AutoSnapshotPolicySet", *resp)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	return If2Slice(retSet)
}

func (s *AutoSnapshotSrv) modifyAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().ModifyAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) associatedAutoSnapshotPolicy(input map[string]interface{}) (data map[string]interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	resp, err = s.GetConn().ApplyAutoSnapshotPolicy(&input)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	associationSet, err := getSdkValue("ReturnSet", *resp)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	return If2Map(associationSet)
}

func (s *AutoSnapshotSrv) unassociatedAutoSnapshotPolicy(input map[string]interface{}) (data map[string]interface{}, err error) {
	var (
		resp   *map[string]interface{}
		retSet interface{}
	)
	resp, err = s.GetConn().CancelAutoSnapshotPolicy(&input)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	retSet, err = getSdkValue("ReturnSet", *resp)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	return If2Map(retSet)
}

// readAutoSnapshotPolicyVolumeAssociationById will check the information whether contains auto_snapshot_policy_id
// after query volumes by volume_id
// if the return information not contains auto_snapshot_policy_id will return nil
func (s *AutoSnapshotSrv) readAutoSnapshotPolicyVolumeAssociationById(volumeId string) (data []interface{}, err error) {
	var (
		resp   *map[string]interface{}
		retSet interface{}
	)

	reqParameters := &map[string]interface{}{
		"VolumeId.1": volumeId,
	}

	resp, err = s.GetConnEbs().DescribeVolumes(reqParameters)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	retSet, err = getSdkValue("Volumes", *resp)
	if err != nil {
		return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
	}

	return If2Slice(retSet)
}

func (s *AutoSnapshotSrv) readAutoSnapshotPolicyVolumeAssociationAll() (data []interface{}, err error) {
	var (
		resp   *map[string]interface{}
		retSet interface{}
	)

	condition := make(map[string]interface{})

	return pageQuery(condition, "MaxResults", "Marker", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		resp, err = s.GetConnEbs().DescribeVolumes(&condition)
		if err != nil {
			return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
		}

		retSet, err = getSdkValue("Volumes", *resp)
		if err != nil {
			return nil, fmt.Errorf("the ksyun sdk internal error detail: %s", err)
		}

		return If2Slice(retSet)

	})
}

func (s *AutoSnapshotSrv) filterVolumesAutoSnapshotPolicyAssociations(
	volumes []interface{}, policyId string) (
	data []interface{}, err error) {

	if len(volumes) < 1 {
		return nil, fmt.Errorf("volumes doesn't include elements")
	}

	for _, v := range volumes {
		sdkPolicyMap, err := If2Map(v)
		if err != nil {
			return nil, err
		}
		policyIdIf, ok := sdkPolicyMap["AutoSnapshotPolicyId"]
		if !ok {
			continue
		}
		if policyId == "" {
			data = append(data, v)
		} else {
			sdkPolicyId := policyIdIf.(string)

			if sdkPolicyId == policyId {
				data = append(data, v)
			}
		}

	}
	return data, err
}

func (s *AutoSnapshotSrv) GetConn() *kec.Kec {
	return s.client.kecconn
}

func (s *AutoSnapshotSrv) GetConnEbs() *ebs.Ebs {
	return s.client.ebsconn
}
