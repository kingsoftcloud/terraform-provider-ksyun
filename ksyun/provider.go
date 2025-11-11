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
		ksyun_dnats
		ksyun_private_dns_records
		ksyun_private_dns_zones
		ksyun_direct_connects

	Resource
		ksyun_vpc
		ksyun_subnet
		ksyun_nat
		ksyun_nat_associate
		ksyun_nat_instance_bandwidth_limit
		ksyun_dnat
		ksyun_network_acl
		ksyun_network_acl_entry
		ksyun_network_acl_associate
		ksyun_route
		ksyun_security_group
		ksyun_security_group_entry
		ksyun_security_group_entry_lite
		ksyun_kec_network_interface
		ksyun_private_dns_zone
		ksyun_private_dns_record
		ksyun_private_dns_zone_vpc_attachment
		ksyun_direct_connect_gateway
		ksyun_direct_connect_gateway_route
		ksyun_direct_connect_interface
		ksyun_direct_connect_bfd_config
		ksyun_dc_interface_associate

VPN

	Data Source
		ksyun_vpn_gateways
		ksyun_vpn_customer_gateways
		ksyun_vpn_tunnels
		ksyun_vpn_gateway_routes

	Resource
		ksyun_vpn_gateway
		ksyun_vpn_customer_gateway
		ksyun_vpn_tunnel
		ksyun_vpn_gateway_route

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

ALB

	Data Source
		ksyun_albs
		ksyun_alb_listeners
		ksyun_alb_rule_groups
		ksyun_alb_listener_cert_groups
        ksyun_alb_backend_server_groups

	Resource
		ksyun_alb
		ksyun_alb_listener
		ksyun_alb_rule_group
		ksyun_alb_listener_cert_group
		ksyun_alb_backend_server_group
		ksyun_alb_register_backend_server
		ksyun_alb_listener_associate_acl

CEN

	Data Source
		ksyun_cens

	Resource
		ksyun_cen

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
		ksyun_auto_snapshot_volume_association
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
		ksyun_snapshot

Bare Metal

	Data Source
		ksyun_bare_metal_images
		ksyun_bare_metal_raid_attributes
		ksyun_bare_metals

	Resource
		ksyun_bare_metal
    ksyun_bare_metal_hot_standby_action

KCE

	Data Source
		ksyun_kce_clusters
		ksyun_kce_instance_images

	Resource
		ksyun_kce_cluster
		ksyun_kce_cluster_attach_existence
		ksyun_kce_cluster_attachment
		ksyun_kce_auth_attachment

KCR

	Data Source
		ksyun_kcrs_instances
		ksyun_kcrs_tokens
		ksyun_kcrs_namespaces
		ksyun_kcrs_webhook_triggers

	Resource
		ksyun_kcrs_instance
		ksyun_kcrs_namespace
		ksyun_kcrs_token
		ksyun_kcrs_webhook_trigger
		ksyun_kcrs_vpc_attachment

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

Clickhouse

	Data Source
		ksyun_clickhouse

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
		ksyun_tag_v2
		ksyun_tag_v2_attachment

KS3

	Data Source
		ksyun_ks3_buckets

	Resource
		ksyun_ks3_bucket

KNAD

	Data Source
		ksyun_knads

	Resource
		ksyun_knad
		ksyun_knad_associate

IAM

	Data Source
		ksyun_iam_users
		ksyun_iam_roles
		ksyun_iam_groups

	Resource
		ksyun_iam_user
		ksyun_iam_role
		ksyun_iam_group
		ksyun_iam_policy
		ksyun_iam_relation_policy
KPFS
	Resource
		ksyun_kpfs_acl

KFW

	Data Source
		ksyun_kfw_instances
		ksyun_kfw_acls
		ksyun_kfw_addrbooks
		ksyun_kfw_service_groups

	Resource
		ksyun_kfw_instance
		ksyun_kfw_acl

		ksyun_kfw_addrbook
		ksyun_kfw_service_group

Monitor

	Resource
		ksyun_monitor_alarm_policy

KMR

	Data Source
    ksyun_kmr_clusters

KLog

	Data Source
    ksyun_klog_projects
*/

package ksyun

import (
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KSYUN_ENDPOINT", ""),
				Description: descriptions["endpoint"],
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
			"http_keepalive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether use http keepalive, if false, disables HTTP keep-alives and will only use the connection to the server for a single HTTP request",
			},
			"force_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force use https protocol for communication between sdk and remote server.",
			},
			"max_retries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 99),
			},
			"http_proxy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(i interface{}, s string) (ret []string, errs []error) {
					ip := i.(string)
					_, err := url.Parse(ip)
					if err != nil {
						errs = append(errs, err)
					}
					return
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ksyun_albs":                     dataSourceKsyunAlbs(),
			"ksyun_alb_listeners":            dataSourceKsyunAlbListeners(),
			"ksyun_alb_rule_groups":          dataSourceKsyunAlbRuleGroups(),
			"ksyun_alb_listener_cert_groups": dataSourceKsyunAlbListenerCertGroups(),
			"ksyun_lines":                    dataSourceKsyunLines(),
			"ksyun_eips":                     dataSourceKsyunEips(),
			"ksyun_slbs":                     dataSourceKsyunLbs(),
			"ksyun_lbs":                      dataSourceKsyunLbs(),
			"ksyun_listeners":                dataSourceKsyunListeners(),
			"ksyun_health_checks":            dataSourceKsyunHealthChecks(),
			// 注册两个同样的data，应该去掉一个。。。文档保留ksyun_lb_listener_servers
			"ksyun_listener_servers":                 dataSourceKsyunLbListenerServers(),
			"ksyun_lb_listener_servers":              dataSourceKsyunLbListenerServers(),
			"ksyun_lb_acls":                          dataSourceKsyunSlbAcls(),
			"ksyun_availability_zones":               dataSourceKsyunAvailabilityZones(),
			"ksyun_network_interfaces":               dataSourceKsyunNetworkInterfaces(),
			"ksyun_network_acls":                     dataSourceKsyunNetworkAcls(),
			"ksyun_vpcs":                             dataSourceKsyunVpcs(),
			"ksyun_subnets":                          dataSourceKsyunSubnets(),
			"ksyun_subnet_available_addresses":       dataSourceKsyunSubnetAvailableAddresses(),
			"ksyun_subnet_allocated_ip_addresses":    dataSourceKsyunSubnetAllocatedIpAddresses(),
			"ksyun_security_groups":                  dataSourceKsyunSecurityGroups(),
			"ksyun_instances":                        dataSourceKsyunInstances(),
			"ksyun_local_volumes":                    dataSourceKsyunLocalVolumes(),
			"ksyun_local_snapshots":                  dataSourceKsyunLocalSnapshots(),
			"ksyun_images":                           dataSourceKsyunImages(),
			"ksyun_sqlservers":                       dataSourceKsyunSqlServer(),
			"ksyun_krds":                             dataSourceKsyunKrds(),
			"ksyun_krds_security_groups":             dataSourceKsyunKrdsSecurityGroup(),
			"ksyun_ks3_buckets":                      dataSourceKsyunKs3Buckets(),
			"ksyun_certificates":                     dataSourceKsyunCertificates(),
			"ksyun_ssh_keys":                         dataSourceKsyunSSHKeys(),
			"ksyun_redis_instances":                  dataSourceRedisInstances(),
			"ksyun_redis_security_groups":            dataSourceRedisSecurityGroups(),
			"ksyun_volumes":                          dataSourceKsyunVolumes(),
			"ksyun_snapshots":                        dataSourceKsyunSnapshots(),
			"ksyun_mongodbs":                         dataSourceKsyunMongodbs(),
			"ksyun_lb_host_headers":                  dataSourceKsyunListenerHostHeaders(),
			"ksyun_lb_rules":                         dataSourceKsyunSlbRules(),
			"ksyun_lb_backend_server_groups":         dataSourceKsyunBackendServerGroups(),
			"ksyun_lb_register_backend_servers":      dataSourceKsyunRegisterBackendServers(),
			"ksyun_routes":                           dataSourceKsyunRoutes(),
			"ksyun_nats":                             dataSourceKsyunNats(),
			"ksyun_scaling_configurations":           dataSourceKsyunScalingConfigurations(),
			"ksyun_scaling_groups":                   dataSourceKsyunScalingGroups(),
			"ksyun_scaling_activities":               dataSourceKsyunScalingActivities(),
			"ksyun_scaling_instances":                dataSourceKsyunScalingInstances(),
			"ksyun_scaling_policies":                 dataSourceKsyunScalingPolicies(),
			"ksyun_scaling_scheduled_tasks":          dataSourceKsyunScalingScheduledTasks(),
			"ksyun_scaling_notifications":            dataSourceKsyunScalingNotifications(),
			"ksyun_rabbitmqs":                        dataSourceKsyunRabbitmqs(),
			"ksyun_vpn_gateways":                     dataSourceKsyunVpnGateways(),
			"ksyun_vpn_customer_gateways":            dataSourceKsyunVpnCustomerGateways(),
			"ksyun_vpn_tunnels":                      dataSourceKsyunVpnTunnels(),
			"ksyun_bwses":                            dataSourceKsyunBandWidthShares(),
			"ksyun_bare_metals":                      dataSourceKsyunBareMetals(),
			"ksyun_bare_metal_images":                dataSourceKsyunBareMetalImages(),
			"ksyun_bare_metal_raid_attributes":       dataSourceKsyunBareMetalRaidAttributes(),
			"ksyun_kce_clusters":                     dataSourceKsyunKceClusters(),
			"ksyun_kce_instance_images":              dataSourceKsyunKceInstanceImages(),
			"ksyun_tags":                             dataSourceKsyunTags(),
			"ksyun_auto_snapshot_policy":             dataSourceKsyunAutoSnapshotPolicy(),
			"ksyun_data_guard_group":                 dataSourceKsyunDataGuardGroup(),
			"ksyun_krds_parameter_group":             dataSourceKsyunKrdsParameterGroup(),
			"ksyun_auto_snapshot_volume_association": dataSourceKsyunAutoSnapshotVolumeAssociation(),
			"ksyun_knads":                            dataSourceKsyunKnads(),
			"ksyun_dnats":                            dataSourceKsyunDnats(),
			"ksyun_alb_backend_server_groups":        dataSourceKsyunAlbBackendServerGroups(),
			"ksyun_vpn_gateway_routes":               dataSourceKsyunVpnGatewayRoutes(),
			"ksyun_kmr_clusters":                     dataSourceKsyunKmrClusters(),

			// private_dns
			"ksyun_private_dns_zones":   dataSourceKsyunPrivateDnsZones(),
			"ksyun_private_dns_records": dataSourceKsyunPrivateDnsRecords(),

			// kcrs
			"ksyun_kcrs_instances":        dataSourceKsyunKcrsInstances(),
			"ksyun_kcrs_tokens":           dataSourceKsyunKcrsTokens(),
			"ksyun_kcrs_namespaces":       dataSourceKsyunKcrsNamespaces(),
			"ksyun_kcrs_webhook_triggers": dataSourceKsyunKcrsWebhookTriggers(),

			// kfw
			"ksyun_kfw_instances":      dataSourceKsyunKfwInstances(),
			"ksyun_kfw_addrbooks":      dataSourceKsyunKfwAddrbooks(),
			"ksyun_kfw_acls":           dataSourceKsyunKfwAcls(),
			"ksyun_kfw_service_groups": dataSourceKsyunKfwServiceGroups(),

			// iam
			"ksyun_iam_users":  dataSourceKsyunIamUsers(),
			"ksyun_iam_roles":  dataSourceKsyunIamRoles(),
			"ksyun_iam_groups": dataSourceKsyunIamGroups(),

			// klog
			"ksyun_klog_projects": dataSourceKsyunKlogProjects(),
			// direct connect
			"ksyun_direct_connects": dataSourceKsyunDirectConnects(),

			// clickhouse
			"ksyun_clickhouse": dataSourceKsyunClickhouse(),
			// cen
			"ksyun_cens": dataSourceKsyunCens(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ksyun_alb":                              resourceKsyunAlb(),
			"ksyun_alb_listener":                     resourceKsyunAlbListener(),
			"ksyun_alb_rule_group":                   resourceKsyunAlbRuleGroup(),
			"ksyun_alb_listener_cert_group":          resourceKsyunAlbListenerCertGroup(),
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
			"ksyun_snapshot":                         resourceKsyunSnapshot(),
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

			"ksyun_ks3_bucket":                       resourceKsyunKs3Bucket(),
			"ksyun_auto_snapshot_policy":             resourceKsyunAutoSnapshotPolicy(),
			"ksyun_auto_snapshot_volume_association": resourceKsyunAutoSnapshotVolumeAssociation(),
			"ksyun_data_guard_group":                 resourceKsyunDataGuardGroup(),
			"ksyun_krds_parameter_group":             resourceKsyunKrdsParameterGroup(),
			"ksyun_knad":                             resourceKsyunKnad(),
			"ksyun_knad_associate":                   resourceKsyunKnadAssociate(),
			"ksyun_nat_instance_bandwidth_limit":     resourceKsyunNatInstanceBandwidthLimit(),
			"ksyun_dnat":                             resourceKsyunDnat(),
			"ksyun_alb_backend_server_group":         resourceKsyunAlbBackendServerGroup(),
			"ksyun_alb_register_backend_server":      resourceKsyunRegisterAlbBackendServer(),
			"ksyun_alb_listener_associate_acl":       resourceKsyunAlbListenerAssociateAcl(),
			"ksyun_vpn_gateway_route":                resourceKsyunVpnGatewayRoute(),

			// kce
			"ksyun_kce_cluster":                  resourceKsyunKceCluster(),
			"ksyun_kce_cluster_attachment":       resourceKsyunKceClusterAttachment(),
			"ksyun_kce_cluster_attach_existence": resourceKsyunKceClusterAttachExistence(),
			"ksyun_kce_auth_attachment":          resourceKsyunKceAuthAttachment(),

			// private dns
			"ksyun_private_dns_zone":                resourceKsyunPrivateDnsZone(),
			"ksyun_private_dns_record":              resourceKsyunPrivateDnsRecord(),
			"ksyun_private_dns_zone_vpc_attachment": resourceKsyunPrivateDnsZoneVpcAttachment(),

			// kcrs
			"ksyun_kcrs_instance":        resourceKsyunKcrsInstance(),
			"ksyun_kcrs_namespace":       resourceKsyunKcrsNamespace(),
			"ksyun_kcrs_token":           resourceKsyunKcrsToken(),
			"ksyun_kcrs_webhook_trigger": resourceKsyunKcrsWebhookTrigger(),
			"ksyun_kcrs_vpc_attachment":  resourceKsyunKcrsVpcAttachment(),

			// tag
			"ksyun_tag_v2":            resourceKsyunTagv2(),
			"ksyun_tag_v2_attachment": resourceKsyunTagv2Attachment(),

			// iam
			"ksyun_iam_user":            resourceKsyunIamUser(),
			"ksyun_iam_role":            resourceKsyunIamRole(),
			"ksyun_iam_group":           resourceKsyunIamGroup(),
			"ksyun_iam_policy":          resourceKsyunIamPolicy(),
			"ksyun_iam_relation_policy": resourceKsyunIamRelationPolicy(),

			// security group
			"ksyun_security_group":            resourceKsyunSecurityGroup(),
			"ksyun_security_group_entry":      resourceKsyunSecurityGroupEntry(),
			"ksyun_security_group_entry_lite": resourceKsyunSecurityGroupEntryLite(),
			// "ksyun_security_group_entry_set":  resourceKsyunSecurityGroupEntrySet(),

			"ksyun_bare_metal_hot_standby_action": resourceKsyunBareMetalHotStandbyAction(),
			// lb
			// "ksyun_lb_listener_associate_backendgroup": resourceKsyunLbListenerAssociateBackendgroup(),

			// kpfs
			"ksyun_kpfs_acl": resourceKsyunKpfsAcl(),

			// direct connect
			"ksyun_direct_connect_gateway":       resourceKsyunDirectConnectGateway(),
			"ksyun_direct_connect_gateway_route": resourceKsyunDirectConnectGatewayRoute(),
			"ksyun_direct_connect_interface":     resourceKsyunDirectConnectInterface(),
			"ksyun_direct_connect_bfd_config":    resourceKsyunDirectConnectBfdConfig(),
			"ksyun_dc_interface_associate":       resourceKsyunDCInterfaceAssociate(),

			"ksyun_kfw_instance":      resourceKsyunCfwInstance(),
			"ksyun_kfw_acl":           resourceKsyunKfwAcl(),
			"ksyun_kfw_addrbook":      resourceKsyunKfwAddrbook(),
			"ksyun_kfw_service_group": resourceKsyunKfwServiceGroup(),

			// monitor
			"ksyun_monitor_alarm_policy": resourceKsyunMonitorAlarmPolicy(),
			// cen
			"ksyun_cen": resourceKsyunCen(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	retryNum := 0
	if mr, ok := d.GetOk("max_retries"); ok {
		retryNum = mr.(int)
	}
	config := Config{
		AccessKey:     d.Get("access_key").(string),
		SecretKey:     d.Get("secret_key").(string),
		Region:        d.Get("region").(string),
		Insecure:      d.Get("insecure").(bool),
		Domain:        d.Get("domain").(string),
		Endpoint:      d.Get("endpoint").(string),
		DryRun:        d.Get("dry_run").(bool),
		IgnoreService: d.Get("ignore_service").(bool),
		HttpKeepAlive: d.Get("http_keepalive").(bool),
		MaxRetries:    retryNum,
		HttpProxy:     d.Get("http_proxy").(string),
		UseSSL:        d.Get("force_https").(bool),
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
		"endpoint":       "",
		"dry_run":        "false",
		"ignore_service": "false",
	}
}
