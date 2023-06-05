// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type KrdsParameterSrv struct {
	client *KsyunClient
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

	results, err = getSdkValue("Data.DBParameterGroups", *resp)
	if err != nil {
		return nil, err
	}
	data := results.([]interface{})
	if err := transformMapValue2String("0.Parameters", data); err != nil {
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

	return pageQuery(input, "MaxRecords", "Marker", 10, 0, func(condition map[string]interface{}) ([]interface{}, error) {
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
		data = results.([]interface{})
		return data, err
	})
}
