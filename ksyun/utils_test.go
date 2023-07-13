// Copyright 2022 NotOne Lv <aiphalv0010@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ksyun

import "testing"

func TestTransformWithFilter(t *testing.T) {
	m := map[string]string{
		"k": "v",
	}

	req := make(map[string]interface{})
	index := 1
	var err error
	for k, v := range m {

		index, err = transformWithFilter(v, k, SdkReqTransform{}, index, &req)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("req +%v", req)
}

func TestRouteFilter(t *testing.T) {
	routeFilter := RouteFilter{
		VpcId:         "25c78605-faa2-43fd-8483-6fb27885eb60",
		InstanceId:    "6f2b44ed-e93d-4075-a92c-48b9007dbdd5",
		DestCidrBlock: "10.10.10.0/21",
	}

	routeMap := routeFilter.Filler()
	t.Logf("route map: +%v", routeMap)

	reqParams, err := routeFilter.GetFilterParams()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("route request parameters: +%v", reqParams)
}
