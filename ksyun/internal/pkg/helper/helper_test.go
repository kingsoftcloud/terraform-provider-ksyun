package helper

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/structor/v1/kce"
)

func TestMapstructureFiller(t *testing.T) {
	var (
		dd *kce.DataDisk
	)

	var val map[string]interface{}
	advanced := kce.AdvancedSetting{
		DataDisk:      dd,
		PreUserScript: "adds",
		UserScript:    "sdasad",
		Schedulable:   false,
		Label:         []kce.Label{{Key: "tf_assembly_kce", Value: "advanced_setting"}},
		ExtraArg: &kce.ExtraArg{
			Kubelet: []string{"custom_arg"},
		},
		ContainerRuntime:     "containerd",
		ContainerPath:        "/data/containerd",
		DockerPath:           "/data/docker",
		ContainerLogMaxFiles: 0,
		ContainerLogMaxSize:  0,
		Taint: []kce.Taint{
			{Key: "key1", Value: "value1", Effect: "effect1"},
			{Key: "key2", Value: "value2", Effect: "effect2"},
		},
	}
	err := MapstructureFiller(advanced, &val, "tf-schema")
	if err != nil {
		t.Errorf("MapstructureFiller() error = %v", err)
		return
	}
	t.Logf("val = %v", val)

}

func TestGetDiffMap(t *testing.T) {
	a := map[string]interface{}{
		"a": "c",
		"v": "v",
		"c": "c",
		"d": "d",
	}
	b := map[string]interface{}{
		"a": "a",
		"v": "v",
		"c": "c",
		"d": "d",
		"e": "e",
	}

	diff := GetDiffMap(a, b)
	t.Logf("diff: %+v", diff)
}
