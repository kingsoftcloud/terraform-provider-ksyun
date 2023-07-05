package ksyun

import (
	"fmt"

	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DataGuardSrv struct {
	client *KsyunClient
}

func NewDataGuardSrv(client *KsyunClient) DataGuardSrv {
	return DataGuardSrv{
		client: client,
	}
}

func (d *DataGuardSrv) describeDataGuardGroup(input map[string]interface{}) (data []interface{}, err error) {
	var (
		resp *map[string]interface{}
	)

	resp, err = d.GetConn().DescribeDataGuardGroup(&input)

	if err != nil {
		return nil, err
	}
	results, err := getSdkValue("DataGuardsSet", *resp)
	if err != nil || results == nil {
		return nil, fmt.Errorf("the current available zone not exsits any data guard group")
	}

	return If2Slice(results)
}

// createDataGuardGroup will create data guard group and it returns this data guard group id
func (d *DataGuardSrv) createDataGuardGroup(input map[string]interface{}) (string, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = d.GetConn().CreateDataGuardGroup(&input)
	if err != nil {
		return "", err
	}

	results, err := getSdkValue("DataGuardId", *resp)
	if err != nil || results == nil {
		return "", err
	}
	return If2String(results)
}

func (d *DataGuardSrv) deleteDataGuardGroup(input map[string]interface{}) error {

	return resource.Retry(RetryTimeoutMinute, func() *resource.RetryError {
		if _, ok := input["DataGuardId.1"]; ok && input["DataGuardId.1"].(string) != "" {
			_, err := d.GetConn().DeleteDataGuardGroups(&input)
			if err != nil {
				return retryError(err)
			}
		}
		return nil
	})
}

func (d *DataGuardSrv) modifyModifyDataGuardGroups(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = d.GetConn().ModifyDataGuardGroups(&input)
	return *resp, err
}

func (d *DataGuardSrv) RemoveVmFromDataGuard(dataGuardId, instanceId string) error {
	removeParam := map[string]interface{}{
		"DataGuardId":  dataGuardId,
		"InstanceId.1": instanceId,
	}

	_, err := d.GetConn().RemoveVmFromDataGuard(&removeParam)
	return err
}

func (d *DataGuardSrv) AddVmIntoDataGuard(dataGuardId, instanceId string) error {
	addParam := map[string]interface{}{
		"DataGuardId":  dataGuardId,
		"InstanceId.1": instanceId,
	}

	_, err := d.GetConn().AddVmIntoDataGuard(&addParam)
	return err
}

func (d *DataGuardSrv) GetConn() *kec.Kec {
	return d.client.kecconn
}
