package ksyun

var lineKeys = map[string]bool{
	"LineName": true,
	"LineId":   true,
	"LineType": true,
}
var eipKeys = map[string]bool{
	"CreateTime":         true,
	"ProjectId":          true,
	"PublicIp":           true,
	"AllocationId":       true,
	"State":              true,
	"LineId":             true,
	"BandWidth":          true,
	"InstanceType":       true,
	"InstanceId":         true,
	"NetworkInterfaceId": true,
	"InternetGatewayId":  true,
	"BandWidthShareId":   true,
	"IsBandWidthShare":   true,
}
var slbKeys = map[string]bool{
	"CreateTime":        true,
	"LoadBalancerName":  true,
	"VpcId":             true,
	"LoadBalancerId":    true,
	"Type":              true,
	"SubnetId":          true,
	"PublicIp":          true,
	"State":             true,
	"LoadBalancerState": true,
}
var listenerKeys = map[string]bool{
	"CreateTime":       true,
	"LoadBalancerId":   true,
	"ListenerName":     true,
	"ListenerId":       true,
	"ListenerState":    true,
	"CertificateId":    true,
	"ListenerProtocol": true,
	"ListenerPort":     true,
	"Method":           true,
	"HealthCheck":      true,
	"Session":          true,
	"RealServer":       true,
}
var healthCheckKeys = map[string]bool{
	"HealthCheckId":      true,
	"HealthCheckState":   true,
	"HealthyThreshold":   true,
	"Interval":           true,
	"Timeout":            true,
	"UnhealthyThreshold": true,
}
var sessionKeys = map[string]bool{
	"SessionPersistencePeriod": true,
	"SessionState":             true,
	"CookieType":               true,
	"CookieName":               true,
}

var serverKeys = map[string]bool{
	"RegisterId":      true,
	"RealServerIp":    true,
	"RealServerPort":  true,
	"RealServerType":  true,
	"InstanceId":      true,
	"RealServerState": true,
	"Weight":          true,
}
var lbAclKeys = map[string]bool{
	"CreateTime":              true,
	"LoadBalancerAclName":     true,
	"LoadBalancerAclId":       true,
	"LoadBalancerAclEntrySet": true,
}

var lbAclEntryKeys = map[string]bool{
	"LoadBalancerAclId":      true,
	"LoadBalancerAclEntryId": true,
	"CidrBlock":              true,
	"RuleNumber":             true,
	"RuleAction":             true,
	"Protocol":               true,
}
var availabilityZoneKeys = map[string]bool{
	"AvailabilityZoneName":  true,
	"AvailabilityZoneState": true,
}
var vpcNetworkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"MacAddress":           true,
	"SecurityGroupSet":     true,
	"InstanceId":           true,
	"InstanceType":         true,
	"PrivateIpAddress":     true,
	"DNS1":                 true,
	"DNS2":                 true,
}

/*
var subnetAvailableAddresseKeys = map[string]bool{
	"AvailableIpAddress": true,
}
var subnetAllocatedIpAddressesKeys = map[string]bool{
	"AvailableIpAddress": true,
}


*/
var vpcKeys = map[string]bool{
	"CidrBlock":  true,
	"CreateTime": true,
	"IsDefault":  true,
	"VpcName":    true,
	//"VpcId":      true,
}

var groupIdentifierKeys = map[string]bool{
	"SecurityGroupId":   true,
	"SecurityGroupName": true,
}
var subnetKeys = map[string]bool{
	"CreateTime":           true,
	"VpcId":                true,
	"SubnetId":             true,
	"SubnetType":           true,
	"SubnetName":           true,
	"CidrBlock":            true,
	"DhcpIpFrom":           true,
	"DhcpIpTo":             true,
	"GatewayIp":            true,
	"Dns1":                 true,
	"Dns2":                 true,
	"NetworkAclId":         true,
	"NatId":                true,
	"AvailbleIPNumber":     true,
	"AvailabilityZoneName": true,
}

var routeKeys = map[string]SdkResponseData{
	"VpcId":                {},
	"DestinationCidrBlock": {},
	"RouteType":            {},
	"CreateTime":           {},
	"NextHopSet":           {Next: routeNextHop},
}

var routeNextHop = map[string]SdkResponseData{
	"GatewayId":   {},
	"GatewayName": {},
}

var natKeys = map[string]SdkResponseData{
	"ProjectId":       {},
	"CreateTime":      {},
	"VpcId":           {},
	"NatName":         {},
	"NatMode":         {},
	"NatType":         {},
	"BandWidth":       {},
	"NatIpSet":        {Next: natIp},
	"AssociateNatSet": {Next: associateNat},
}

var natIp = map[string]SdkResponseData{
	"NatIp":   {},
	"NatIpId": {},
}

var associateNat = map[string]SdkResponseData{
	"SubnetId": {},
}

var vpcSecurityGroupKeys = map[string]bool{
	"CreateTime":            true,
	"VpcId":                 true,
	"SecurityGroupName":     true,
	"SecurityGroupId":       true,
	"SecurityGroupType":     true,
	"SecurityGroupEntrySet": true,
}
var vpcSecurityGroupEntrySetKeys = map[string]bool{
	"Description":          true,
	"SecurityGroupEntryId": true,
	"CidrBlock":            true,
	"Direction":            true,
	"Protocol":             true,
	"IcmpType":             true,
	"IcmpCode":             true,
	"PortRangeFrom":        true,
	"PortRangeTo":          true,
}
var vpcSecurityGroupEntryKeys = map[string]SdkResponseData{
	"Description":          {},
	"SecurityGroupEntryId": {},
	"CidrBlock":            {},
	"Direction":            {},
	"Protocol":             {},
	"IcmpType":             {},
	"IcmpCode":             {},
	"PortRangeFrom":        {},
	"PortRangeTo":          {},
}

var instanceKeys = map[string]bool{
	"InstanceId":            true,
	"ProjectId":             true,
	"InstanceName":          true,
	"InstanceType":          true,
	"InstanceConfigure":     true,
	"ImageId":               true,
	"SubnetId":              true,
	"PrivateIpAddress":      true,
	"InstanceState":         true,
	"Monitoring":            true,
	"NetworkInterfaceSet":   true,
	"SriovNetSupport":       true,
	"IsShowSriovNetSupport": true,
	"CreationDate":          true,
	"AvailabilityZone":      true,
	"AvailabilityZoneName":  true,
	"AutoScalingType":       true,
	"ProductWhat":           true,
	"ChargeType":            true,
	"SystemDisk":            true,
	"KeySet":                true,
	"DataDisks":             true,
}
var instanceConfigureKeys = map[string]bool{
	"VCPU":       true,
	"GPU":        true,
	"MemoryGb":   true,
	"DataDiskGb": true,
	//"RootDiskGb":   true,
	//"DataDiskType": true,
}
var dataDiskKeys = map[string]bool{
	"DiskId":             true,
	"DiskType":           true,
	"DiskSize":           true,
	"DeleteWithInstance": true,
}
var instanceStateKeys = map[string]bool{
	"Name": true,
}
var monitoringKeys = map[string]bool{
	"State": true,
}
var kecNetworkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"MacAddress":           true,
	"SecurityGroupSet":     true,
	"PrivateIpAddress":     true,
	"DNS1":                 true,
	"DNS2":                 true,
	"PublicIp":             true,
}
var kecNetworkInterfaceSetKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"PrivateIpAddress":     true,
	"GroupSet":             true,
	"SecurityGroupSet":     true,
}
var systemDiskKeys = map[string]bool{
	"DiskType": true,
	"DiskSize": true,
}
var groupSetKeys = map[string]bool{
	"GroupId": true,
}
var kecSecurityGroupSetKeys = map[string]bool{
	"SecurityGroupId": true,
}
var kecSecurityGroupKeys = map[string]bool{
	"CreateTime":            true,
	"VpcId":                 true,
	"SecurityGroupName":     true,
	"SecurityGroupId":       true,
	"SecurityGroupType":     true,
	"SecurityGroupEntrySet": true,
}
var imageKeys = map[string]bool{
	"ImageId":          true,
	"Name":             true,
	"ImageState":       true,
	"CreationDate":     true,
	"Platform":         true,
	"IsPublic":         true,
	"InstanceId":       true,
	"IsNpe":            true,
	"UserCategory":     true,
	"SysDisk":          true,
	"Progress":         true,
	"ImageSource":      true,
	"CloudInitSupport": true,
	"Ipv6Support":      true,
	"IsModifyType":     true,
	"IsCloudMarket":    true,
}
var networkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"MacAddress":           true,
	"SecurityGroupSet":     true,
	"InstanceId":           true,
	"InstanceType":         true,
	"PrivateIpAddress":     true,
	"SubnetId":             true,
	"ProjectId":            true,
	"DNS1":                 true,
	"DNS2":                 true,
}
var certificateKeys = map[string]bool{
	"CertificateId":   true,
	"CertificateName": true,
}
var sshKeyKeys = map[string]bool{
	"KeyId":      true,
	"KeyName":    true,
	"PublicKey":  true,
	"CreateTime": true,
}

var redisInstanceKeys = map[string]bool{
	"cacheId":          true,
	"region":           true,
	"az":               true,
	"name":             true,
	"engine":           true,
	"mode":             true,
	"size":             true,
	"port":             true,
	"vip":              true,
	"status":           true,
	"createTime":       true,
	"netType":          true,
	"vpcId":            true,
	"vnetId":           true,
	"billType":         true,
	"orderType":        true,
	"source":           true,
	"serviceStatus":    true,
	"serviceBeginTime": true,
	"serviceEndTime":   true,
	"iamProjectId":     true,
	"iamProjectName":   true,
	"protocol":         true,
	"slaveVip":         true,
	"slaveNum":         true,
	"timingSwitch":     true,
	"timezone":         true,
	"usedMemory":       true,
	"subOrderId":       true,
	"productId":        true,
	"orderUse":         true,
	"readonlyNode":     true,
	//"rules":            true,
	"parameters": true,
	"shardSize":  true,
	"shardNum":   true,
	"eip":        true,
	"eipRo":      true,
}

/*
var epcInstanceKeys = map[string]bool{
	"CreateTime":                   true,
	"HostName":                     true,
	"HostType":                     true,
	"AllowModifyHyperThreading":    true,
	"ReleasableTime":               true,
	"TorName":                      true,
	"CabinetName":                  true,
	"RackName":                     true,
	"HostId":                       true,
	"Sn":                           true,
	"CabinetId":                    true,
	"AvailabilityZone":             true,
	"Raid":                         true,
	"ImageId":                      true,
	"KeyId":                        true,
	"NetworkInterfaceMode":         true,
	"EnableBond":                   true,
	"SecurityAgent":                true,
	"CloudMonitorAgent":            true,
	"ProductType":                  true,
	"OsName":                       true,
	"Memory":                       true,
	"HostStatus":                   true,
	"ClusterId":                    true,
	"EnableContainer":              true,
	"SystemFileType":               true,
	"DataFileType":                 true,
	"DataDiskCatalogue":            true,
	"DataDiskCatalogueSuffix":      true,
	"Cpu":                          true,
	"Gpu":                          true,
	"DiskSet":                      true,
	"NetworkInterfaceAttributeSet": true,
}

var epcNetworkInterfaceKeys = map[string]bool{
	"NetworkInterfaceId":   true,
	"NetworkInterfaceType": true,
	"SubnetId":             true,
	"PrivateIpAddress":     true,
	"DNS1":                 true,
	"DNS2":                 true,
	"Mac":                  true,
	"SecurityGroupSet":     true,
}

var epcDiskSetKeys = map[string]bool{
	"DiskType":        true,
	"Raid":            true,
	"Space":           true,
	"SystemDiskSpace": true,
	"DiskAttribute":   true,
	"DiskCount":       true,
}

var epcCpuKeys = map[string]bool{
	"Model":     true,
	"Frequence": true,
	"Count":     true,
	"CoreCount": true,
}

var epcGpuKeys = map[string]bool{
	"Model":     true,
	"Frequence": true,
	"Count":     true,
	"CoreCount": true,
	"GpuCount":  true,
}

var epcSecurityGroupKeys = map[string]bool{
	"SecurityGroupId": true,
}


*/
var volumeKeys = map[string]bool{
	"VolumeId":         true,
	"VolumeName":       true,
	"VolumeDesc":       true,
	"Size":             true,
	"VolumeStatus":     true,
	"VolumeType":       true,
	"VolumeCategory":   true,
	"InstanceId":       true,
	"CreateTime":       true,
	"AvailabilityZone": true,
	"ProjectId":        true,
}

var mongodbInstanceKeys = map[string]bool{
	"UserId":          true,
	"Region":          true,
	"Name":            true,
	"InstanceId":      true,
	"Status":          true,
	"IP":              true,
	"InstanceType":    true,
	"Version":         true,
	"InstanceClass":   true,
	"Storage":         true,
	"totalStorage":    true,
	"SecurityGroupId": true,
	"Port":            true,
	"NetworkType":     true,
	"VpcId":           true,
	"VnetId":          true,
	"TimingSwitch":    true,
	"Timezone":        true,
	"TimeCycle":       true,
	"ProductId":       true,
	"PayType":         true,
	"ProductWhat":     true,
	"CreateDate":      true,
	"ExpirationDate":  true,
	"IamProjectId":    true,
	"IamProjectName":  true,
	"NodeNum":         true,
	"MongosNum":       true,
	"ShardNum":        true,
	"Mode":            true,
	"Config":          true,
	"Area":            true,
}

var mongodbInstanceNodeKeys = map[string]bool{
	"NodeId": true,
	"Name":   true,
	"Role":   true,
	"IP":     true,
	"Port":   true,
	"Status": true,
}

var mongodbShardInstanceMongosNodeKeys = map[string]bool{
	"NodeId":        true,
	"Name":          true,
	"Role":          true,
	"Endpoint":      true,
	"Status":        true,
	"Connections":   true,
	"InstanceClass": true,
}

var mongodbShardInstanceShardNodeKeys = map[string]bool{
	"NodeId":        true,
	"Name":          true,
	"Status":        true,
	"Disk":          true,
	"Iops":          true,
	"InstanceClass": true,
}
var hostHeaderKeys = map[string]bool{
	"CreateTime":    true,
	"HostHeader":    true,
	"HostHeaderId":  true,
	"ListenerId":    true,
	"CertificateId": true,
}
var slbRuleKeys = map[string]bool{
	"CreateTime":           true,
	"Path":                 true,
	"HostHeaderId":         true,
	"BackendServerGroupId": true,
	"Method":               true,
	"ListenerSync":         true,
	"HealthCheck":          true,
	"Session":              true,
	"RuleId":               true,
}
var backendServerGroupKeys = map[string]bool{
	"CreateTime":             true,
	"BackendServerGroupName": true,
	"BackendServerGroupId":   true,
	"VpcId":                  true,
	"BackendServerNumber":    true,
	"BackendServerGroupType": true,
	"HealthCheck":            true,
}
var registerBackendServerKeys = map[string]bool{
	"CreateTime":           true,
	"BackendServerGroupId": true,
	"BackendServerIp":      true,
	"RegisterId":           true,
	"RealServerIp":         true,
	"RealServerPort":       true,
	"RealServerType":       true,
	"MasterSlaveType":      true,
	"InstanceId":           true,
	"NetworkInterfaceId":   true,
	"RealServerState":      true,
	"Weight":               true,
}

var redisSecKeys = map[string]bool{
	"securityGroupId": true,
	"name":            true,
	"description":     true,
	"created":         true,
	"updated":         true,
}

var rabbitmqInstanceKeys = map[string]bool{
	"UserId":           true,
	"Region":           true,
	"InstanceName":     true,
	"InstanceId":       true,
	"Status":           true,
	"StatusName":       true,
	"Vip":              true,
	"WebVip":           true,
	"InstanceType":     true,
	"SsdDisk":          true,
	"Protocol":         true,
	"SecurityGroupId":  true,
	"Port":             true,
	"NetworkType":      true,
	"VpcId":            true,
	"SubnetId":         true,
	"ProductId":        true,
	"BillType":         true,
	"CreateDate":       true,
	"ExpirationDate":   true,
	"ProjectId":        true,
	"ProjectName":      true,
	"NodeNum":          true,
	"AvailabilityZone": true,
	"ProductWhat":      true,
	"Mode":             true,
	"ModeName":         true,
	"Eip":              true,
	"WebEip":           true,
	"EipEgress":        true,
}
