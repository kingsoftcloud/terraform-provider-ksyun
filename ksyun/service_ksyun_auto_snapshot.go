// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/service/kec"
)

type AutoSnapshotSrv struct {
	client *KsyunClient
}

// querySnapshotPolicyByID will query snapshot policy from ksyun open-api
func (s *AutoSnapshotSrv) querySnapshotPolicyByID(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)

	resp, err = s.GetConn().DescribeAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) createAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().CreateAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) deleteAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().DeleteAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) modifyAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().ModifyAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) associatedAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().ApplyAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) unassociatedAutoSnapshotPolicy(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = s.GetConn().CancelAutoSnapshotPolicy(&input)
	return *resp, err
}

func (s *AutoSnapshotSrv) GetConn() *kec.Kec {
	return s.client.kecconn
}
