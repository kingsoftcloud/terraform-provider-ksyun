package kce

import (
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/structor/v1/kec"
)

type DescribeClusterInstanceResponse struct {
	InstanceSet []InstanceSet `json:"InstanceSet" mapstructure:"InstanceSet"`
	RequestID   string        `json:"RequestId" mapstructure:"RequestID"`
	MaxResults  int           `json:"MaxResults" mapstructure:"MaxResults"`
	Marker      int           `json:"Marker" mapstructure:"Marker"`
	TotalCount  int           `json:"TotalCount" mapstructure:"TotalCount"`
}

type KecInstancePara struct {
	DedicatedID         string                    `json:"DedicatedId" mapstructure:"DedicatedID"`
	ProjectID           int                       `json:"ProjectId" mapstructure:"ProjectID"`
	InstanceType        string                    `json:"InstanceType" mapstructure:"InstanceType"`
	InstanceConfigure   kec.InstanceConfigure     `json:"InstanceConfigure" mapstructure:"InstanceConfigure"`
	SystemDisk          kec.SystemDisk            `json:"SystemDisk" mapstructure:"SystemDisk"`
	ImageID             string                    `json:"ImageId" mapstructure:"ImageID"`
	PrivateIPAddress    string                    `json:"PrivateIpAddress" mapstructure:"PrivateIPAddress"`
	ChargeType          string                    `json:"ChargeType" mapstructure:"ChargeType"`
	CreateTime          interface{}               `json:"CreateTime" mapstructure:"CreateTime"`
	AvailabilityZone    string                    `json:"AvailabilityZone" mapstructure:"AvailabilityZone"`
	SubnetID            string                    `json:"SubnetId" mapstructure:"SubnetID"`
	VpcID               string                    `json:"VpcId" mapstructure:"VpcID"`
	NetworkInterfaceSet []kec.NetworkInterfaceSet `json:"NetworkInterfaceSet" mapstructure:"NetworkInterfaceSet"`
}
type DataDisk struct {
	AutoFormatAndMount bool   `json:"AutoFormatAndMount" mapstructure:"AutoFormatAndMount" tf-schema:"auto_format_and_mount"`
	FileSystem         string `json:"FileSystem" mapstructure:"FileSystem" tf-schema:"file_system"`
	MountTarget        string `json:"MountTarget" mapstructure:"MountTarget" tf-schema:"mount_target"`
}
type Label struct {
	Key   string `json:"Key" mapstructure:"Key" tf-schema:"key"`
	Value string `json:"Value" mapstructure:"Value" tf-schema:"value"`
}
type AdvancedSetting struct {
	DataDisk             *DataDisk `json:"DataDisk" mapstructure:"DataDisk" tf-schema:"data_disk,tf-list"`
	PreUserScript        string    `json:"PreUserScript" mapstructure:"PreUserScript" tf-schema:"pre_user_script"`
	UserScript           string    `json:"UserScript" mapstructure:"UserScript" tf-schema:"user_script"`
	Schedulable          bool      `json:"Schedulable" mapstructure:"Schedulable" tf-schema:"schedulable"`
	Label                []Label   `json:"Label" mapstructure:"Label" tf-schema:"label"`
	ExtraArg             *ExtraArg `json:"ExtraArg" mapstructure:"ExtraArg"`
	ContainerRuntime     string    `json:"ContainerRuntime" mapstructure:"ContainerRuntime" tf-schema:"container_runtime"`
	ContainerPath        string    `json:"ContainerPath" mapstructure:"ContainerPath" tf-schema:"container_path"`
	DockerPath           string    `json:"DockerPath" mapstructure:"DockerPath" tf-schema:"docker_path"`
	ContainerLogMaxFiles int       `json:"ContainerLogMaxFiles" mapstructure:"ContainerLogMaxFiles" tf-schema:"container_log_max_files"`
	ContainerLogMaxSize  int       `json:"ContainerLogMaxSize" mapstructure:"ContainerLogMaxSize" tf-schema:"container_log_max_size"`
	Taint                []Taint   `json:"Taint" mapstructure:"Taint"`
}

type Taint struct {
	Key    string `json:"Key" mapstructure:"Key" tf-schema:"key"`
	Value  string `json:"Value" mapstructure:"Value" tf-schema:"value"`
	Effect string `json:"Effect" mapstructure:"Effect" tf-schema:"effect"`
}

type ExtraArg struct {
	Kubelet []string `json:"Kubelet" mapstructure:"Kubelet" tf-schema:"kubelet"`
}

// InstanceSet for describeClusterInstance
type InstanceSet struct {
	InstanceID        string           `json:"InstanceId" mapstructure:"InstanceID"`
	Type              string           `json:"Type" mapstructure:"Type"`
	InstanceName      string           `json:"InstanceName" mapstructure:"InstanceName"`
	InstanceRole      string           `json:"InstanceRole" mapstructure:"InstanceRole"`
	InstanceStatus    string           `json:"InstanceStatus" mapstructure:"InstanceStatus"`
	KecInstancePara   KecInstancePara  `json:"KecInstancePara" mapstructure:"KecInstancePara"`
	AdvancedSetting   *AdvancedSetting `json:"AdvancedSetting" mapstructure:"AdvancedSetting"`
	UnSchedulable     bool             `json:"UnSchedulable" mapstructure:"UnSchedulable"`
	DrainStatus       string           `json:"DrainStatus" mapstructure:"DrainStatus"`
	NodePoolID        string           `json:"NodePoolId" mapstructure:"NodePoolID"`
	Available         bool             `json:"Available" mapstructure:"Available"`
	UnavailableReason string           `json:"UnavailableReason" mapstructure:"UnavailableReason"`
	ErrorMessage      string           `json:"ErrorMessage" mapstructure:"ErrorMessage"`
}
