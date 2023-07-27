package ksyun

import (
	"testing"
)

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

func TestStructureConverter(t *testing.T) {
	bb := NewDescribeDnatsParams()
	bb.DnatIds = append(bb.DnatIds, "dasdasdas")
	m := map[string]interface{}{}
	if err := StructureConverter(bb, &m); err != nil {
		t.Fatal(err)
	}
	t.Logf("mapstructure: +%v", m)
}
func TestStructureConverter2(t *testing.T) {
	s := CreateDnatParams{
		NatIp: "dd",
	}
	m := map[string]interface{}{}

	if err := StructureConverter(s, &m); err != nil {
		t.Fatal(err)
	}
	t.Logf("+%v", m)
}
