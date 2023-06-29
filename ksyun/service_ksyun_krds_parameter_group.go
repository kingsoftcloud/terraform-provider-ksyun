package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type KrdsParameterSrv struct {
	client *KsyunClient
}

func NewKrdsParameterSrv(client *KsyunClient) KrdsParameterSrv {
	return KrdsParameterSrv{
		client: client,
	}
}

func (s *KrdsParameterSrv) GetConn() *krds.Krds {
	return s.client.krdsconn
}

// describeDBParameterGroup will query DB parameter groups from ksyun open-api
func (s *KrdsParameterSrv) describeDBParameterGroupById(input map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		err     error
		results interface{}
	)

	resp, err = s.GetConn().DescribeDBParameterGroup(&input)
	if err != nil {
		return nil, err
	}

	results, err = getSdkValue("Data.DBParameterGroups", *resp)
	if err != nil {
		return nil, err
	}
	data, err := If2Slice(results)
	if err != nil {
		return nil, err
	}

	return data, err
}

// describeDBParameterGroupAll will deal with all results of query by page limited
func (s *KrdsParameterSrv) describeDBParameterGroupAll(input map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return pageQuery(input, "MaxRecords", "Marker", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		conn := s.GetConn()
		action := "DescribeDBParameterGroup"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = conn.DescribeDBParameterGroup(nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = conn.DescribeDBParameterGroup(&condition)
			if err != nil {
				return data, err
			}
		}

		results, err = getSdkValue("Data.DBParameterGroups", *resp)
		if err != nil {
			return data, err
		}
		return If2Slice(results)
	})
}

func (s *KrdsParameterSrv) createDBParameterGroup(input map[string]interface{}) (string, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().CreateDBParameterGroup(&input)
	if err != nil {
		return "", err
	}

	results, err := getSdkValue("Data.DBParameterGroup.DBParameterGroupId", *resp)
	if err != nil || results == nil {
		return "", err
	}
	return If2String(results)
}

func (s *KrdsParameterSrv) deleteDBParameterGroup(input map[string]interface{}) (err error) {
	// _, err = s.GetConn().DeleteDBParameterGroup(&input)
	return resource.Retry(RetryTimeoutMinute, func() *resource.RetryError {
		_, err = s.GetConn().DeleteDBParameterGroup(&input)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
}

func (s *KrdsParameterSrv) modifyDBParameterGroup(input map[string]interface{}) (data map[string]interface{}, err error) {
	var resp *map[string]interface{}
	resp, err = s.GetConn().ModifyDBParameterGroup(&input)
	if err != nil {
		return data, err
	}
	results, err := getSdkValue("Data.DBParameterGroup", *resp)
	if err != nil || results == nil {
		return data, err
	}

	return If2Map(results)
}
