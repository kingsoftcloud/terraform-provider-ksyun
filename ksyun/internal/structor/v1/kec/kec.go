package kec

type DescribeInstancesResponse struct {
	Marker        int        `json:"Marker" mapstructure:"Marker"`
	InstanceCount int        `json:"InstanceCount" mapstructure:"InstanceCount"`
	RequestId     string     `json:"RequestId" mapstructure:"RequestId"`
	InstancesSet  []Instance `json:"InstancesSet" mapstructure:"InstancesSet"`
}
type InstanceConfigure struct {
	VCPU         int    `json:"VCPU" mapstructure:"VCPU"`
	GPU          int    `json:"GPU" mapstructure:"GPU"`
	MemoryGb     int    `json:"MemoryGb" mapstructure:"MemoryGb"`
	DataDiskGb   int    `json:"DataDiskGb" mapstructure:"DataDiskGb"`
	RootDiskGb   int    `json:"RootDiskGb" mapstructure:"RootDiskGb"`
	DataDiskType string `json:"DataDiskType" mapstructure:"DataDiskType"`
	VGPU         string `json:"VGPU" mapstructure:"VGPU"`
	Spec         string `json:"Spec" mapstructure:"Spec"`
}
type InstanceState struct {
	Name               string `json:"Name" mapstructure:"Name"`
	OnMigrate          bool   `json:"OnMigrate" mapstructure:"OnMigrate"`
	TaskState          string `json:"TaskState" mapstructure:"TaskState"`
	PreMigrationStatus string `json:"PreMigrationStatus" mapstructure:"PreMigrationStatus"`
	CostTime           string `json:"CostTime" mapstructure:"CostTime"`
	TimeStamp          string `json:"TimeStamp" mapstructure:"TimeStamp"`
}
type VMConfig struct {
	InstanceType    string `json:"InstanceType" mapstructure:"InstanceType"`
	CPU             int    `json:"Cpu" mapstructure:"Cpu"`
	Mem             int    `json:"Mem" mapstructure:"Mem"`
	ProductType     int    `json:"ProductType" mapstructure:"ProductType"`
	SystemType      string `json:"SystemType" mapstructure:"SystemType"`
	SystemSize      int    `json:"SystemSize" mapstructure:"SystemSize"`
	DedicatedHostId string `json:"DedicatedHostId" mapstructure:"DedicatedHostId"`
	Model           string `json:"Model" mapstructure:"Model"`
}
type VolumeConfig struct {
	VolumeType        string `json:"VolumeType" mapstructure:"VolumeType"`
	VolumeSize        int    `json:"VolumeSize" mapstructure:"VolumeSize"`
	VolumeProductType int    `json:"VolumeProductType" mapstructure:"VolumeProductType"`
	VolumeId          string `json:"VolumeId" mapstructure:"VolumeId"`
}
type PreMigrateConfig struct {
	VMConfig     VMConfig     `json:"VmConfig" mapstructure:"VMConfig"`
	VolumeConfig VolumeConfig `json:"VolumeConfig" mapstructure:"VolumeConfig"`
}
type Monitoring struct {
	State string `json:"State" mapstructure:"State"`
}
type GroupSet struct {
	GroupId string `json:"GroupId" mapstructure:"GroupId"`
}
type SecurityGroupSet struct {
	SecurityGroupId string `json:"SecurityGroupId" mapstructure:"SecurityGroupId"`
}
type NetworkInterfaceSet struct {
	NetworkInterfaceId   string             `json:"NetworkInterfaceId" mapstructure:"NetworkInterfaceId"`
	NetworkInterfaceType string             `json:"NetworkInterfaceType" mapstructure:"NetworkInterfaceType"`
	VpcId                string             `json:"VpcId" mapstructure:"VpcId"`
	SubnetId             string             `json:"SubnetId" mapstructure:"SubnetId"`
	MacAddress           string             `json:"MacAddress" mapstructure:"MacAddress"`
	PrivateIpAddress     string             `json:"PrivateIpAddress" mapstructure:"PrivateIpAddress"`
	GroupSet             []GroupSet         `json:"GroupSet" mapstructure:"GroupSet"`
	SecurityGroupSet     []SecurityGroupSet `json:"SecurityGroupSet" mapstructure:"SecurityGroupSet"`
	NetworkInterfaceName string             `json:"NetworkInterfaceName" mapstructure:"NetworkInterfaceName"`
}
type DataGuardSet struct {
	DataGuardId   string `json:"DataGuardId" mapstructure:"DataGuardId"`
	DataGuardName string `json:"DataGuardName" mapstructure:"DataGuardName"`
	DataGuardType string `json:"DataGuardType" mapstructure:"DataGuardType"`
}
type Metadata struct {
	Id           string `json:"Id" mapstructure:"Id"`
	Key          bool   `json:"Key" mapstructure:"Key"`
	Value        string `json:"Value" mapstructure:"Value"`
	InstanceUUId string `json:"Instance_uuid" mapstructure:"InstanceUUId"`
}
type SystemDisk struct {
	DiskType string `json:"DiskType" mapstructure:"DiskType" tf-schema:"disk_type"`
	DiskSize int    `json:"DiskSize" mapstructure:"DiskSize" tf-schema:"disk_size"`
}
type DataDisks struct {
	DiskId             string `json:"DiskId" mapstructure:"DiskId"`
	DiskType           string `json:"DiskType" mapstructure:"DiskType"`
	DiskSize           int    `json:"DiskSize" mapstructure:"DiskSize"`
	DeleteWithInstance bool   `json:"DeleteWithInstance" mapstructure:"DeleteWithInstance"`
	Encrypted          bool   `json:"Encrypted" mapstructure:"Encrypted"`
}
type Instance struct {
	InstanceId            string                `json:"InstanceId" mapstructure:"InstanceId"`
	ProjectId             float64               `json:"ProjectId" mapstructure:"ProjectId"`
	ShutdownNoCharge      bool                  `json:"ShutdownNoCharge" mapstructure:"ShutdownNoCharge"`
	IsDistributeIpv6      bool                  `json:"IsDistributeIpv6" mapstructure:"IsDistributeIpv6"`
	InstanceName          string                `json:"InstanceName" mapstructure:"InstanceName"`
	InstanceType          string                `json:"InstanceType" mapstructure:"InstanceType"`
	InstanceConfigure     InstanceConfigure     `json:"InstanceConfigure" mapstructure:"InstanceConfigure"`
	ImageId               string                `json:"ImageId" mapstructure:"ImageId"`
	SubnetId              string                `json:"SubnetId" mapstructure:"SubnetId"`
	PrivateIPAddress      string                `json:"PrivateIpAddress" mapstructure:"PrivateIPAddress"`
	InstanceState         InstanceState         `json:"InstanceState" mapstructure:"InstanceState"`
	IsRunOptimised        bool                  `json:"IsRunOptimised" mapstructure:"IsRunOptimised"`
	PreMigrateConfig      PreMigrateConfig      `json:"PreMigrateConfig" mapstructure:"PreMigrateConfig"`
	Monitoring            Monitoring            `json:"Monitoring" mapstructure:"Monitoring"`
	NetworkInterfaceSet   []NetworkInterfaceSet `json:"NetworkInterfaceSet" mapstructure:"NetworkInterfaceSet"`
	SriovNetSupport       string                `json:"SriovNetSupport" mapstructure:"SriovNetSupport"`
	IsShowSriovNetSupport bool                  `json:"IsShowSriovNetSupport" mapstructure:"IsShowSriovNetSupport"`
	CreationDate          string                `json:"CreationDate" mapstructure:"CreationDate"`
	AvailabilityZone      string                `json:"AvailabilityZone" mapstructure:"AvailabilityZone"`
	AvailabilityZoneName  string                `json:"AvailabilityZoneName" mapstructure:"AvailabilityZoneName"`
	DataGuardSet          []DataGuardSet        `json:"DataGuardSet" mapstructure:"DataGuardSet"`
	DedicatedName         string                `json:"DedicatedName" mapstructure:"DedicatedName"`
	DedicatedClusterId    string                `json:"DedicatedClusterId" mapstructure:"DedicatedClusterId"`
	DedicatedUUId         string                `json:"DedicatedUuid" mapstructure:"DedicatedUUId"`
	AutoScalingType       string                `json:"AutoScalingType" mapstructure:"AutoScalingType"`
	LargeEBSCapacity      string                `json:"LargeEBSCapacity" mapstructure:"LargeEBSCapacity"`
	Metadata              []Metadata            `json:"Metadata" mapstructure:"Metadata"`
	ProductType           int                   `json:"ProductType" mapstructure:"ProductType"`
	ProductWhat           int                   `json:"ProductWhat" mapstructure:"ProductWhat"`
	LiveUpgradeSupport    bool                  `json:"LiveUpgradeSupport" mapstructure:"LiveUpgradeSupport"`
	ChargeType            string                `json:"ChargeType" mapstructure:"ChargeType"`
	SystemDisk            SystemDisk            `json:"SystemDisk" mapstructure:"SystemDisk" tf-schema:"system_disk"`
	HostName              string                `json:"HostName" mapstructure:"HostName"`
	UserData              string                `json:"UserData" mapstructure:"UserData"`
	Migration             int                   `json:"Migration" mapstructure:"Migration"`
	DataDisks             []DataDisks           `json:"DataDisks" mapstructure:"DataDisks"`
	VncSupport            bool                  `json:"VncSupport" mapstructure:"VncSupport"`
	Platform              string                `json:"Platform" mapstructure:"Platform"`
}
