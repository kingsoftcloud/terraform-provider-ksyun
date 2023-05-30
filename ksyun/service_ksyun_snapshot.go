// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/aws/aws-sdk-go/aws/request"
)

type SnapshotSrv struct {
	client *KsyunClient
}

var fakeDescribeAutoSnapshotPolicy = "DescribeAutoSnapshotPolicy"

// querySnapshotPolicyByID will query snapshot policy from ksyun open-api
func (s *SnapshotSrv) querySnapshotPolicyByID(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)

	resp, err = s.DescribeAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *SnapshotSrv) createAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.CreateAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *SnapshotSrv) deleteAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.DeleteAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *SnapshotSrv) modifyAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.ModifyAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *SnapshotSrv) associatedAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.ApplyAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *SnapshotSrv) unassociatedAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.CancelAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *SnapshotSrv) GetConn() *kec.Kec {
	return s.client.kecconn
}

// the below functions are fake to simulate the ksyun-go-sdk function
func (s *SnapshotSrv) CreateAutoSnapshotPolicy(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := s.AutoSnapshotPolicyRequestByAction("CreateAutoSnapshotPolicy", input)
	return out, req.Send()
}

func (s *SnapshotSrv) ModifyAutoSnapshotPolicy(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := s.AutoSnapshotPolicyRequestByAction("ModifyAutoSnapshotPolicy", input)
	return out, req.Send()
}

func (s *SnapshotSrv) ApplyAutoSnapshotPolicy(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := s.AutoSnapshotPolicyRequestByAction("ApplyAutoSnapshotPolicy", input)
	return out, req.Send()
}

func (s *SnapshotSrv) CancelAutoSnapshotPolicy(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := s.AutoSnapshotPolicyRequestByAction("CancelAutoSnapshotPolicy", input)
	return out, req.Send()
}

func (s *SnapshotSrv) DeleteAutoSnapshotPolicy(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := s.AutoSnapshotPolicyRequestByAction("DeleteAutoSnapshotPolicy", input)
	return out, req.Send()
}

func (s *SnapshotSrv) DescribeAutoSnapshotPolicy(input *map[string]interface{}) (*map[string]interface{}, error) {
	req, out := s.AutoSnapshotPolicyRequestByAction(fakeDescribeAutoSnapshotPolicy, input)
	return out, req.Send()
}

func (s *SnapshotSrv) AutoSnapshotPolicyRequestByAction(action string, input *map[string]interface{}) (req *request.Request, output *map[string]interface{}) {
	op := &request.Operation{
		Name:       action,
		HTTPMethod: "GET",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &map[string]interface{}{}
	}

	output = &map[string]interface{}{}
	req = s.GetConn().NewRequest(op, input, output)

	return
}
