/*
The Ksyun provider is used to interact with many resources supported by Ksyun.
The provider needs to be configured with the proper credentials before it can be used.

~> **NOTE:** This guide requires an available Ksyun account or sub-account with project to create resources.

The Ksyun provider is used to interact with the
resources supported by Ksyun. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

# Example Usage

```hcl

	terraform {
	  required_providers {
	    ksyun = {
	      source  = "kingsoftcloud/ksyun"
	    }
	  }
	}

# Configure the ksyun Provider

	provider "ksyun" {
	   access_key = "your ak"
	   secret_key = "your sk"
	   region = "cn-beijing-6"
	}

```

# Resources List

Provider Data Sources

	ksyun_availability_zones

EIP

	Data Source
		ksyun_eips
		ksyun_bwses
		ksyun_lines

	Resource
		ksyun_eip
		ksyun_eip_associate
		ksyun_bws
		ksyun_bws_associate

VPC

	Data Source
		ksyun_vpcs
		ksyun_nats
		ksyun_network_acls
		ksyun_subnets
		ksyun_network_interfaces
		ksyun_routes
		ksyun_security_groups
		ksyun_subnet_allocated_ip_addresses
		ksyun_subnet_available_addresses

	Resource
		ksyun_vpc
		ksyun_subnet
		ksyun_nat
		ksyun_nat_associate
		ksyun_network_acl
		ksyun_network_acl_entry
		ksyun_network_acl_associate
		ksyun_route
		ksyun_security_group
		ksyun_security_group_entry
		ksyun_kec_network_interface

VPN

	Data Source
		ksyun_vpn_gateways
		ksyun_vpn_customer_gateways
		ksyun_vpn_tunnels

	Resource
		ksyun_vpn_gateway
		ksyun_vpn_customer_gateway
		ksyun_vpn_tunnel

SLB

	Data Source
		ksyun_lb_backend_server_groups
		ksyun_health_checks
		ksyun_lb_acls
		ksyun_lb_host_headers
		ksyun_lb_listener_servers
		ksyun_lb_rules
		ksyun_lbs
		ksyun_listeners
		ksyun_lb_register_backend_servers

	Resource
		ksyun_healthcheck
		ksyun_lb
		ksyun_lb_acl
		ksyun_lb_acl_entry
		ksyun_lb_backend_server_group
		ksyun_lb_host_header
		ksyun_lb_register_backend_server
		ksyun_lb_listener
		ksyun_lb_listener_associate_acl
		ksyun_lb_listener_server
		ksyun_lb_rule

SSH key

	Data Source
		ksyun_ssh_keys

	Resource
		ksyun_ssh_key

Instance(KEC)

	Data Source
		ksyun_images
		ksyun_instances
		ksyun_local_volumes
		ksyun_local_snapshots
		ksyun_auto_snapshot_policy
		ksyun_auto_snapshot_policy_volume_association
		ksyun_data_guard_group

	Resource
		ksyun_instance
		ksyun_kec_network_interface_attachment
		ksyun_auto_snapshot_policy
		ksyun_auto_snapshot_volume_association
		ksyun_data_guard_group

Volume(EBS)

	Data Source
		ksyun_volumes
		ksyun_snapshots

	Resource
		ksyun_volume
		ksyun_volume_attach

Bare Metal

	Data Source
		ksyun_bare_metal_images
		ksyun_bare_metal_raid_attributes
		ksyun_bare_metals

	Resource
		ksyun_bare_metal

KCE

	Data Source
		ksyun_kce_clusters
		ksyun_kce_instance_images

	Resource
		ksyun_kce_cluster
		ksyun_kce_worker

KCM

	Data Source
		ksyun_certificates

	Resource
		ksyun_certificate

KRDS

	Data Source
		ksyun_krds
		ksyun_krds_security_groups
		ksyun_krds_parameter_group

	Resource
		ksyun_krds
		ksyun_krds_rr
		ksyun_krds_security_group
		ksyun_krds_security_group_rule
		ksyun_krds_parameter_group

SQLServer

	Data Source
		ksyun_sqlservers

MongoDB

	Data Source
		ksyun_mongodbs

	Resource
		ksyun_mongodb_instance

RabbitMQ

	Data Source
		ksyun_rabbitmqs

	Resource
		ksyun_rabbitmq_instance
		ksyun_rabbitmq_security_rule

Redis

	Data Source
		ksyun_redis_instances
		ksyun_redis_security_groups

	Resource
		ksyun_redis_instance
		ksyun_redis_instance_node
		ksyun_redis_sec_group

Auto Scaling

	Data Source
		ksyun_scaling_activities
		ksyun_scaling_configurations
		ksyun_scaling_groups
		ksyun_scaling_instances
		ksyun_scaling_notifications
		ksyun_scaling_policies
		ksyun_scaling_scheduled_tasks

	Resource
		ksyun_scaling_configuration
		ksyun_scaling_group
		ksyun_scaling_instance
		ksyun_scaling_notification
		ksyun_scaling_policy
		ksyun_scaling_scheduled_task

Tag

	Data Source
		ksyun_tags

	Resource
		ksyun_tag

KS3

	Resource
		ksyun_ks3_bucket
*/
package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_ACCESS_KEY", nil),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_SECRET_KEY", nil),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_REGION", nil),
				Description: descriptions["region"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_INSECURE", true),
				Description: descriptions["insecure"],
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_DOMAIN", ""),
				Description: descriptions["domain"],
			},
			"dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_DRY_RUN", false),
				Description: descriptions["dry_run"],
			},
			"ignore_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_DOMAIN_IGNORE_SERVICE", false),
				Description: descriptions["ignore_service"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ksyun_lines":         dataSourceKsyunLines(),
			"ksyun_eips":          dataSourceKsyunEips(),
			"ksyun_slbs":          dataSourceKsyunLbs(),
			"ksyun_lbs":           dataSourceKsyunLbs(),
			"ksyun_listeners":     dataSourceKsyunListeners(),
			"ksyun_health_checks": dataSourceKsyunHealthChecks(),
			// 注册两个同样的data，应该去掉一个。。。文档保留ksyun_lb_listener_servers
			"ksyun_listener_servers":                        dataSourceKsyunLbListenerServers(),
			"ksyun_lb_listener_servers":                     dataSourceKsyunLbListenerServers(),
			"ksyun_lb_acls":                                 dataSourceKsyunSlbAcls(),
			"ksyun_availability_zones":                      dataSourceKsyunAvailabilityZones(),
			"ksyun_network_interfaces":                      dataSourceKsyunNetworkInterfaces(),
			"ksyun_network_acls":                            dataSourceKsyunNetworkAcls(),
			"ksyun_vpcs":                                    dataSourceKsyunVpcs(),
			"ksyun_subnets":                                 dataSourceKsyunSubnets(),
			"ksyun_subnet_available_addresses":              dataSourceKsyunSubnetAvailableAddresses(),
			"ksyun_subnet_allocated_ip_addresses":           dataSourceKsyunSubnetAllocatedIpAddresses(),
			"ksyun_security_groups":                         dataSourceKsyunSecurityGroups(),
			"ksyun_instances":                               dataSourceKsyunInstances(),
			"ksyun_local_volumes":                           dataSourceKsyunLocalVolumes(),
			"ksyun_local_snapshots":                         dataSourceKsyunLocalSnapshots(),
			"ksyun_images":                                  dataSourceKsyunImages(),
			"ksyun_sqlservers":                              dataSourceKsyunSqlServer(),
			"ksyun_krds":                                    dataSourceKsyunKrds(),
			"ksyun_krds_security_groups":                    dataSourceKsyunKrdsSecurityGroup(),
			"ksyun_certificates":                            dataSourceKsyunCertificates(),
			"ksyun_ssh_keys":                                dataSourceKsyunSSHKeys(),
			"ksyun_redis_instances":                         dataSourceRedisInstances(),
			"ksyun_redis_security_groups":                   dataSourceRedisSecurityGroups(),
			"ksyun_volumes":                                 dataSourceKsyunVolumes(),
			"ksyun_snapshots":                               dataSourceKsyunSnapshots(),
			"ksyun_mongodbs":                                dataSourceKsyunMongodbs(),
			"ksyun_lb_host_headers":                         dataSourceKsyunListenerHostHeaders(),
			"ksyun_lb_rules":                                dataSourceKsyunSlbRules(),
			"ksyun_lb_backend_server_groups":                dataSourceKsyunBackendServerGroups(),
			"ksyun_lb_register_backend_servers":             dataSourceKsyunRegisterBackendServers(),
			"ksyun_routes":                                  dataSourceKsyunRoutes(),
			"ksyun_nats":                                    dataSourceKsyunNats(),
			"ksyun_scaling_configurations":                  dataSourceKsyunScalingConfigurations(),
			"ksyun_scaling_groups":                          dataSourceKsyunScalingGroups(),
			"ksyun_scaling_activities":                      dataSourceKsyunScalingActivities(),
			"ksyun_scaling_instances":                       dataSourceKsyunScalingInstances(),
			"ksyun_scaling_policies":                        dataSourceKsyunScalingPolicies(),
			"ksyun_scaling_scheduled_tasks":                 dataSourceKsyunScalingScheduledTasks(),
			"ksyun_scaling_notifications":                   dataSourceKsyunScalingNotifications(),
			"ksyun_rabbitmqs":                               dataSourceKsyunRabbitmqs(),
			"ksyun_vpn_gateways":                            dataSourceKsyunVpnGateways(),
			"ksyun_vpn_customer_gateways":                   dataSourceKsyunVpnCustomerGateways(),
			"ksyun_vpn_tunnels":                             dataSourceKsyunVpnTunnels(),
			"ksyun_bwses":                                   dataSourceKsyunBandWidthShares(),
			"ksyun_bare_metals":                             dataSourceKsyunBareMetals(),
			"ksyun_bare_metal_images":                       dataSourceKsyunBareMetalImages(),
			"ksyun_bare_metal_raid_attributes":              dataSourceKsyunBareMetalRaidAttributes(),
			"ksyun_kce_clusters":                            dataSourceKsyunKceClusters(),
			"ksyun_kce_instance_images":                     dataSourceKsyunKceInstanceImages(),
			"ksyun_tags":                                    dataSourceKsyunTags(),
			"ksyun_auto_snapshot_policy":                    dataSourceKsyunAutoSnapshotPolicy(),
			"ksyun_data_guard_group":                        dataSourceKsyunDataGuardGroup(),
			"ksyun_krds_parameter_group":                    dataSourceKsyunKrdsParameterGroup(),
			"ksyun_auto_snapshot_policy_volume_association": dataSourceKsyunAutoSnapshotVolumeAssociation(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ksyun_eip":                              resourceKsyunEip(),
			"ksyun_eip_associate":                    resourceKsyunEipAssociation(),
			"ksyun_lb":                               resourceKsyunLb(),
			"ksyun_healthcheck":                      resourceKsyunHealthCheck(),
			"ksyun_lb_listener":                      resourceKsyunListener(),
			"ksyun_lb_listener_server":               resourceKsyunInstancesWithListener(),
			"ksyun_lb_acl":                           resourceKsyunLoadBalancerAcl(),
			"ksyun_lb_acl_entry":                     resourceKsyunLoadBalancerAclEntry(),
			"ksyun_lb_listener_associate_acl":        resourceKsyunListenerAssociateAcl(),
			"ksyun_vpc":                              resourceKsyunVpc(),
			"ksyun_subnet":                           resourceKsyunSubnet(),
			"ksyun_security_group":                   resourceKsyunSecurityGroup(),
			"ksyun_security_group_entry":             resourceKsyunSecurityGroupEntry(),
			"ksyun_instance":                         resourceKsyunInstance(),
			"ksyun_sqlserver":                        resourceKsyunSqlServer(),
			"ksyun_kec_network_interface":            resourceKsyunKecNetworkInterface(),
			"ksyun_kec_network_interface_attachment": resourceKsyunKecNetworkInterfaceAttachment(),
			"ksyun_krds":                             resourceKsyunKrds(),
			"ksyun_krds_rr":                          resourceKsyunKrdsRr(),
			"ksyun_krds_security_group":              resourceKsyunKrdsSecurityGroup(),
			"ksyun_krds_security_group_rule":         resourceKsyunKrdsSecurityGroupRule(),
			"ksyun_certificate":                      resourceKsyunCertificate(),
			"ksyun_ssh_key":                          resourceKsyunSSHKey(),
			"ksyun_redis_instance":                   resourceRedisInstance(),
			"ksyun_redis_instance_node":              resourceRedisInstanceNode(),
			"ksyun_redis_sec_group":                  resourceRedisSecurityGroup(),
			"ksyun_redis_sec_group_rule":             resourceRedisSecurityGroupRule(),
			"ksyun_redis_sec_group_allocate":         resourceRedisSecurityGroupAllocate(),
			"ksyun_mongodb_instance":                 resourceKsyunMongodbInstance(),
			"ksyun_mongodb_shard_instance":           resourceKsyunMongodbShardInstance(),
			"ksyun_mongodb_shard_instance_node":      resourceKsyunMongodbShardInstanceNode(),
			"ksyun_mongodb_security_rule":            resourceKsyunMongodbSecurityRule(),
			"ksyun_volume":                           resourceKsyunVolume(),
			"ksyun_volume_attach":                    resourceKsyunVolumeAttach(),
			"ksyun_lb_rule":                          resourceKsyunSlbRule(),
			"ksyun_lb_host_header":                   resourceKsyunListenerHostHeader(),
			"ksyun_lb_backend_server_group":          resourceKsyunBackendServerGroup(),
			"ksyun_lb_register_backend_server":       resourceKsyunRegisterBackendServer(),
			"ksyun_route":                            resourceKsyunRoute(),
			"ksyun_nat":                              resourceKsyunNat(),
			"ksyun_nat_associate":                    resourceKsyunNatAssociation(),
			"ksyun_scaling_configuration":            resourceKsyunScalingConfiguration(),
			"ksyun_scaling_group":                    resourceKsyunScalingGroup(),
			"ksyun_scaling_instance":                 resourceKsyunScalingInstance(),
			"ksyun_scaling_policy":                   resourceKsyunScalingPolicy(),
			"ksyun_scaling_scheduled_task":           resourceKsyunScalingScheduledTask(),
			"ksyun_scaling_notification":             resourceKsyunScalingNotification(),
			"ksyun_rabbitmq_instance":                resourceKsyunRabbitmq(),
			"ksyun_rabbitmq_security_rule":           resourceKsyunRabbitmqSecurityRule(),
			"ksyun_network_acl":                      resourceKsyunNetworkAcl(),
			"ksyun_network_acl_entry":                resourceKsyunNetworkAclEntry(),
			"ksyun_network_acl_associate":            resourceKsyunNetworkAclAssociate(),
			"ksyun_vpn_gateway":                      resourceKsyunVpnGateway(),
			"ksyun_vpn_customer_gateway":             resourceKsyunVpnCustomerGateway(),
			"ksyun_vpn_tunnel":                       resourceKsyunVpnTunnel(),
			"ksyun_bws":                              resourceKsyunBandWidthShare(),
			"ksyun_bws_associate":                    resourceKsyunBandWidthShareAssociate(),
			"ksyun_bare_metal":                       resourceKsyunBareMetal(),
			"ksyun_tag":                              resourceKsyunTag(),
			"ksyun_kce_cluster":                      resourceKsyunKceCluster(),
			"ksyun_kce_worker":                       resourceKsyunKceWorker(),
			"ksyun_ks3_bucket":                       resourceKsyunKs3Bucket(),
			"ksyun_auto_snapshot_policy":             resourceKsyunAutoSnapshotPolicy(),
			"ksyun_auto_snapshot_volume_association": resourceKsyunAutoSnapshotVolumeAssociation(),
			"ksyun_data_guard_group":                 resourceKsyunDataGuardGroup(),
			"ksyun_krds_parameter_group":             resourceKsyunKrdsParameterGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey:     d.Get("access_key").(string),
		SecretKey:     d.Get("secret_key").(string),
		Region:        d.Get("region").(string),
		Insecure:      d.Get("insecure").(bool),
		Domain:        d.Get("domain").(string),
		DryRun:        d.Get("dry_run").(bool),
		IgnoreService: d.Get("ignore_service").(bool),
	}
	client, err := config.Client()
	return client, err
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key":     "ak",
		"secret_key":     "sk",
		"region":         "cn-beijing-6",
		"insecure":       "true",
		"domain":         "",
		"dry_run":        "false",
		"ignore_service": "false",
	}
}
