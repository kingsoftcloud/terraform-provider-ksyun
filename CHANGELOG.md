## 1.18.5 (Mar 21, 2025)

BUGFIX:

- `ksyun_kpfs_acl`:  修复terraform import 的问题，更新import示例

## 1.18.4 (Mar 18, 2025)

IMPROVEMENTS:

- `ksyun_bare_metal`: 增加`password_inherit` `data_disk_mount` `storage_roce_network_card_name`


## 1.18.3 (Mar 18, 2025)

FEATURES:

- - **New Resource:** `ksyun_kpfs_acl` KPFS posix授权管理
 
## 1.18.2 (Mar 14, 2025)

BUGFIX:

- `ksyun_bare_metal`: 兼容openapi返回tags

## 1.18.1 (Mar 13, 2025)

IMPROVEMENT:

- `ksyun_bare_metal`: 变更raid枚举值，支持指定试用时间
- `ksyun_bare_metal_hot_standby_action`: 新增热备机替换操作
  
## 1.18.0 (Mar 13, 2025)

IMPROVEMENT:

- `ksyun_bare_metal`: 新增热备机处理，更新实例创建字段

## 1.17.9 (Mar 11, 2025)

BUGFIX:

- `ksyun_iam_policy`: 修复change策略文档不生效问题

## 1.17.8 (Feb 24, 2025)

BUGFIX:

- `ksyun_instance`: 修复sync_tag字段设置为false时无法获取字段信息
  
## 1.17.7 (Jan 02, 2025)

BUGFIX:

- `ksyun_alb_backend_server_group`: 修复当删除触发异常时，无法正确处理异常，误认为服务器组已被删除

## 1.17.6 (Dec 20, 2024)

IMPROVEMENT:

- `ksyun_alb_listener`: 新增默认转发规则返回固定响应和重写配置
- `ksyun_alb_rule_group`: 新增默认转发规则返回固定响应和重写配置， rule_set 新增header、method、query、sourceip等类型
- `ksyun_krds`: 增加tags字段
- `ksyun_krds_rr`: 增加tags字段
- `ksyun_redis_instance`: 增加tags字段
- `ksyun_mongodb_instance`: 增加tags字段

BUGFIX:

- `ksyun_redis_instance`: 修复查询实例的配置参数的响应变更问题

## 1.17.5(Dec 16, 2024)

IMPROVEMENTS:

- `ksyun_iam_relation_policy`: 新增policy_type(system、custom)类型，支持创建系统或自定义策略

## 1.17.4(Dec 13, 2024)

FEATURES:

- - **New Resource:** `ksyun_iam_relation_policy` IAM策略关联

## 1.17.3 (Nov 26, 2024)

IMPROVEMENTS:

- `ksyun_instance`: 新增ebs tag同步，优化data_disks管理
- `ksyun_nat`: 新增tags嵌入管理
- `ksyun_alb`: 新增tags嵌入管理
- `ksyun_ebs`: 新增tags嵌入管理
- `ksyun_bws`: 新增tags嵌入管理

## 1.17.2 (Nov 22, 2024)

BUGFIX:

- `ksyun_tag_v2_attachment` 修复批量创建时出现的资源找不到的问题

## 1.17.1 (Oct 21, 2024)

IMPROVEMENTS:

- `ksyun_alb`: 新增`enable_hpa`字段，支持开启HPA、`delete_protection`字段，支持开启删除保护、`modification_protection`字段，支持开启修改保护、`enabled_quic`字段，支持开启QUIC.
- `ksyun_alb`: 支持创建中阶版
- `ksyun_alb_listener`: 新增`ca_enabled` `ca_certificate_id`开启双向验证，`enable_quic_upgrade`开启QUIC升级, `quic_listener_id` QUIC监听器ID

## 1.17.0(Oct 08, 2024)

FEATURES:

- - **New Resource:** `ksyun_iam_user` IAM用户管理
- - **New Resource:** `ksyun_iam_role` IAM角色管理
- - **New Resource:** `ksyun_iam_group` IAM用户组管理
- - **New Resource:** `ksyun_iam_policy` IAM策略管理
- - **New Data Source:** `ksyun_iam_users` IAM用户列表
- - **New Data Source:** `ksyun_iam_roles` IAM角色列表
- - **New Data Source:** `ksyun_iam_groups` IAM角色列表

## 1.16.3 (Aug 22, 2024)

BUGFIX:

- `ksyun_instance`: 修复云主机变更镜像时`user_data`不生效

## 1.16.2 (Jul 30, 2024)

FEATURES:

- - **Deprecated Resource:** `ksyun_lb_listener_associate_backendgroup` 使用`backend_server_group_mounted` 字段替代

IMPROVEMENTS:

- `ksyun_lb_listener`: 新增`backend_server_group_mounted`字段，支持挂载后端服务器组
- `ksyun_nat`: `nat_line_id` 标记为Computed


## 1.16.1 (Jul 23, 2024)

FEATURES:

- - New Resource: `ksyun_lb_listener_associate_backendgroup` 监听器挂载后端服务器组

IMPROVEMENTS:

- `ksyun_lb_listener`: 新增`bind_type`
- `ksyun_lb_backend_server_group`: 新增`protocol`

## 1.16.0 (Jul 05, 2024)

FEATURES:

- - **New Resource:** `ksyun_kce_cluster_attachment` KCE集群添加新节点
- - **New Resource:** `ksyun_kce_cluster_attach_existence` KCE集群添加已有节点 替换 `ksyun_kce_worker`
- - **New Resource:** `ksyun_kce_cluster` KCE集群支持创建多机型节点，支持托管集群

## 1.15.7 (May 28, 2024)

IMPROVEMENTS:

- `ksyun_alb_rule_group` 新增`http_method` 健康检查HTTP方法，`health_protocol` 健康检查协议，`health_port` 健康检查端口;

## 1.15.6 (May 13, 2024)

BUGFIX:
- `ksyun_security_group_entry` 修复了在创建安全组规则时，触发的not exist问题


## 1.15.5 (May 13, 2024)

BUGFIX:
- `ksyun_security_group_entry_lite` 增加cidr_block字段重复判定，修复部分创建成功但terraform无法感知的问题

## 1.15.4 (Apr 23, 2024)

BUGFIX:
- `ksyun_alb_backend_server_group` 修复`protocol` change
  
## 1.15.3 (Mar 27, 2024)

BUGFIX:
- `ksyun_lb_listener` 弃用`http_protocol`字段
- `ksyun_alb_backend_server_group` 新增`protocol`字段，支持`http`和`grpc`两种协议

## 1.15.2 (Mar 18, 2024)

BUGFIX:
- `ksyun_lb_listener` 修复`health_check`在stop状态下触发的plan diff
- `ksyun_alb_backend_server_group` 移除health_check配置

## 1.15.1 (Mar 4, 2024)

BUGFIX:
- `ksyun_lb_listener` 支持`method`参数被修改

## 1.15.0 (Feb 19, 2024)

FEATURES:

- - **New Data Source:** `ksyun_ks3_buckets`

IMPROVEMENTS:

- `ksyun_ks3_bucket` 新增`policy` 空间策略; 新增`abort_incomplete_multipart_upload` 碎片管理
- `provider` 新增`endpoint` 配置参数

BUGFIX:

- `ksyun_ks3_bucket` 修复plan not empty

## 1.14.3 (Jan 26, 2024)

IMPROVEMENTS:

- -**New Resource:** `ksyun_security_group_entry_lite` 安全组批量管理资源

## 1.14.2 (Jan 17, 2024)

BUGFIX:

- `resource_ksyun_ssh_key` 修复public key comments引发的plan diff

## 1.14.1 (Jan 12, 2024)

IMPROVEMENTS:

- `ksyun_alb` 支持创建私网类型alb实例

## 1.14.0 (Jan 04, 2024)

FEATURES:

- - **New Resource:** `ksyun_kcrs_instance` 容器镜像服务-仓库实例
- - **New Resource:** `ksyun_kcrs_namespace` 容器镜像服务-命名空间
- - **New Resource:** `ksyun_kcrs_token`  容器镜像服务-访问凭证
- - **New Resource:** `ksyun_kcrs_webhook_trigger`  容器镜像服务-Webhook触发器
- - **New Resource:** `ksyun_kcrs_vpc_attachment`  容器镜像服务-内网访问
- - **New Data Source:** `ksyun_kcrs_instance`
- - **New Data Source:** `ksyun_kcrs_namespace`
- - **New Data Source:** `ksyun_kcrs_token`
- - **New Data Source:** `ksyun_kcrs_webhook_trigger`

BUGFIX:

- `ksyun_alb_backend_server_group` 弃用健康检查配置

## 1.13.1 (Dec 28, 2023) 

IMPROVEMENTS:

- `ksyun_alb_rule_group` 新增`fixed_response_config` 固定响应配置
- `ksyun_alb_rule_group` 新增`type` type类型须与后端资源类型对应
- `ksyun_alb_listener` 新增`fixed_response_config` 默认转发策略增加固定响应配置

## 1.13.0 (Dec 18, 2023) 

FEATURES:

- - **New Resource:** `ksyun_tag_v2` 标签管理
- - **New Resource:** `ksyun_tag_v2_attachment` 资源绑定标签，支持控制台所有的`resource_type` 

## 1.12.1 (Dec 14, 2023) 

IMPROVEMENTS:

- `ksyun_healthcheck` 新增`health_check_connect_port`健康检查端口
- `ksyun_alb_listener` 新增`config_content` 个性化配置


## 1.12.0 (Dec 7, 2023)

FEATURES:

- - **New Resource:** `ksyun_alb_backend_server_group` 应用型负载均衡上游服务器组
- - **New Resource:** `ksyun_alb_listener_associate_acl` 应用型负载均衡-Alb Acl绑定alb监听器
- - **New Resource:** `ksyun_alb_register_backend_server` 应用型负载均衡-服务器注册服务器组
- - **New Data Source:** `ksyun_alb_backend_server_groups`

IMPROVEMENTS:

- `ksyun_healthcheck` 新增`lb_type`兼容alb; 新增`http_method` 参数
- `ksyun_alb_rule_group` 新增`redirect_http_code` `redirect_alb_listener_id` 重定向alb监听器
- `ksyun_alb_listener` 新增`default_forward_rule`默认转发策略

BUGFIX:

- `ksyun_alb_listener` 更正tls 1.3参数

## 1.11.0 (Nov 23, 2023)

FEATURES:

- - **New Resource:** `ksyun_private_dns_zone` 内网DNS2.0
- - **New Resource:** `ksyun_private_dns_record` 内网DNS解析记录
- - **New Resource:** `ksyun_private_dns_zone_vpc_attachment` 内网DNS Zone绑定VPC
- - **New Data Source:** `ksyun_private_dns_records` 
- - **New Data Source:** `data_source_ksyun_private_dns_zone` 

## 1.10.1 (Oct 20, 2023)

IMPROVEMENTS:

- `provider` 新增`force_https`字段
- 优化http post请求发送tcp包数量，避免connection reset异常

## 1.10.0 (Oct 20, 2023)

FEATURES:

- - **New Resource:** `ksyun_vpn_gateway_route` Vpn网关路由，仅在Vpn2.0时生效
- - **New Data Source:** `ksyun_vpn_gateway_routes`

IMPROVEMENTS:

- `ksyun_vpn_gateway` 支持VPN2.0
- `ksyun_vpn_gateway_tunnel` 支持VPN2.0

## 1.9.1 (Sep 25, 2023)

BUGFIX:

- `ksyun_bws_associate` 修复读取remote associations未处理异常而出现的空指针异常

## 1.9.0 (Sep 5, 2023)

FEATURES:

- - **New Resource:** `ksyun_alb` 应用型负载均衡
- - **New Resource:** `ksyun_alb_listener` 应用型负载均衡-监听器
- - **New Resource:** `ksyun_alb_rule_group` 应用型负载均衡-转发策略
- - **New Resource:** `ksyun_alb_listener_cert_group` 应用型负载均衡-证书组
- - **New Data Source:** `ksyun_albs`
- - **New Data Source:** `ksyun_alb_listeners`
- - **New Data Source:** `ksyun_alb_rule_groups`
- - **New Data Source:** `ksyun_alb_listener_cert_groups`

BUGFIX:

- `ksyun_eip_associate` 修复在ForceNew更新时触发的eip绑定状态不一致

## 1.8.0 (Aug 16, 2023)

FEATURES:

- `provider` 新增`http_httpalive` `max_retries` `http_proxy` 配置参数

IMPROVEMENTS:

- `provider` 增加网络异常重试

## 1.7.2 (Aug 14, 2023)

BUGFIX:
- `ksyun_vpc` 修复可能触发的无法删除vpc的操作

## 1.7.1 (Aug 3, 2023)

IMPROVEMENTS:

- `ksyun_bare_metal` 新增裸金属开机/关机操作
- `ksyun_bare_metal` 支持创建roce网络机型


## 1.7.0 (Jul 26, 2023)

FEATURES:

- `ksyun_nat` 新增Nat Ip 管理 `nat_ip`
- `ksyun_nat_associate` 支持nat绑定主机网卡
- **New Resource:** `ksyun_nat_instance_bandwidth_limit` 新增nat 服务器限速
- **New Resource:** `ksyun_dnat` 新增dnat规则管理
- **New Data Source:** `ksyun_dnats`

## 1.6.1 (Jul 14, 2023)

BUGFIX:

- `ksyun_kec_network_interface` 修复弹性网卡修改子网、安全组和主私网IP操作并兼容辅助私网IP

## 1.6.0 (Jul 13, 2023)

FEATURES:

- `ksyun_kec_network_interface` 弹性网卡新增辅助私网IP管理`secondary_private_ips` and `secondary_private_ip_address_count`
- `ksyun_vpc` 支持IPv6网段
- `ksyun_subnet` 支持IPv6网段

## 1.5.0 (Jul 12, 2023)

FEATURES:

- - **New Resource:** `ksyun_knad` knad 原生高防
- - **New Resource:** `ksyun_knad_associate`
- - **New Data Source:** `ksyun_knads` 

## 1.4.2 (Jul 07, 2023)

IMPROVEMENTS:

- `kec instance` 新增系统盘类型`ESSD_SYSTEM_PL0`,`ESSD_SYSTEM_PL1`,`ESSD_SYSTEM_PL2`

## 1.4.1 (Jul 06, 2023)

IMPROVEMENTS:

- 新增kec实例移入移出容灾组操作
- 支持krds实例使用参数组作为参数创建模版

BUG FIXES:

- 修复krds_parameter_group参数组参数修改时可能引发的异常

## 1.4.0 (Jun 29, 2023)

FEATURES:

- - **New Resource:** `ksyun_krds_parameter_group` krds 参数组管理
- - **New Resource:** `ksyun_data_guard_group`
- - **New Resource:** `ksyun_auto_snapshot_policy` 自动快照策略管理
- - **New Resource:** `ksyun_auto_snapshot_policy_volume_association` 自动快照策略应用到云硬盘
- - **New Resource:** `ksyun_snapshot` EBS快照
- - **New Data Source:** `ksyun_krds_parameter_group`
- - **New Data Source:** `ksyun_data_guard_group`
- - **New Data Source:** `ksyun_auto_snapshot_policy`
- - **New Data Source:** `ksyun_auto_snapshot_policy_volume_association`

BUG FIXES:

- 修复kec实例升降配时的重启逻辑(ksyun_instance: instance_type)
  - NOTE: 实例降级操作会触发实例关机

## 1.3.79 (Jun 15, 2023)

IMPROVEMENTS:

- EBS数据盘增加ESSD_PL0类型

## 1.3.78 (Jun 12, 2023)

IMPROVEMENTS:

- data_ksyun_snapshots: EBS快照
- data_ksyun_local_snapshots: 本地盘快照

## 1.3.77 (Jun 12, 2023)

IMPROVEMENTS:

- 删除redis安全组的机制优化
- SLB监听器的健康检查字段设为Deprecated，使用ksyun_healthcheck代替

## 1.3.76 (May 15, 2023)

IMPROVEMENTS:

- data_ksyun_images返回增加real_image_id字段

## 1.3.75 (May 5, 2023)

BUG FIXES:

- 修复kec更配不调整套餐时参数缺少instanceType的问题

## 1.3.74 (April 28, 2023)

BUG FIXES:

- 修复slice类型的interface映射不准确导致panic的问题

## 1.3.73 (April 19, 2023)

- user-agent增加provider版本

## 1.3.71 (April 10, 2023)

- 增加ks3文档

## 1.3.70 (April 6, 2023)

- ### KS3

  RESOURCES:

  * bucket create
  * bucket read
  * bucket update 
  * bucket delete

## 1.3.68 (Mar 20, 2023)

BUG FIXES:

- 修复EIP字段设置的BUG。

## 1.3.67 (Mar 14, 2023)

IMPROVEMENTS:

- 基于代码description生成文档

## 1.3.66 (Mar 3, 2023)

BUG FIXES:

- 修复redis清理安全组由于缓存数据未更新导致删除失败的问题

## 1.3.65 (Mar 3, 2023)

BUG FIXES:

- 修复redis清理安全组panic的问题

## 1.3.64 (Mar 3, 2023)

IMPROVEMENTS:

- redis支持delete_directly参数

## 1.3.63 (Mar 1, 2023)

BUG FIXES:

- 去掉KRDS计费方式校验

## 1.3.62 (Dec 29, 2022)

BUG FIXES:

- 修复KRDS创建失败后，临时参数组没有清理的问题

## 1.3.61 (Dec 27, 2022)

BUG FIXES:

- 修复KRDS修改parameters，float类型参数不生效的问题
- 修复KRDS新建实例，parameters不生效，必须创建完成后再modify一次的问题
- 增加KRDS的force_restart参数说明

## 1.3.60 (Dec 15, 2022)

BUG FIXES:

- 修复EIP通过其他方式修改项目制后，无法从tf获取EIP数据的问题。


## 1.3.59 (Dec 2, 2022)

BUG FIXES:

- 修复krds只读实例不能更配的问题：ksyun_krds_rr的db_instance_class字段，ForceNew置为false


## 1.3.58 (Nov 9, 2022)

BUG FIXES:

- 修复LB监听器绑定ACL后，如果通过其他方式解绑，tf不能正常获取到解绑状态的问题


## 1.3.57 (Nov 9, 2022)

IMPROVEMENTS:

- 使用整机镜像创建主机，默认不自动使用镜像关联的快照创建数据盘，需要手动配置盘资源。


## 1.3.56 (Nov 1, 2022)

IMPROVEMENTS:

- 支持一键三连，更配会自动关机，自动重启
- instance的security_group_id增加MinItems限制
- volume增加snapshot_id字段
- 优化instance和volume的snapshot_id字段的diff判断（由于api不返回该字段，diff默认忽略）
- 新增EBS类型：ESSD_PL1、ESSD_PL2、ESSD_PL3

BUG FIXES:

- 修复更配无法触发的问题
- 修复机型其他属性修改触发网卡更新的问题


## 1.3.55 (Oct 8, 2022)

BUG FIXES:

- 修正LB日志参数的文档错误
- 修正LB的TAG资源类型问题


## 1.3.54 (Sep 28, 2022)

BUG FIXES:

- 修正创建LB同时开启日志不生效的问题


## 1.3.53 (Sep 28, 2022)

BUG FIXES:

- 更正LB日志功能的example和文档 


## 1.3.52 (Sep 20, 2022)

BUG FIXES:

- 修正EBS盘未修改但触发resize的问题


## 1.3.51 (Sep 15, 2022)

BUG FIXES:

- 修正KRDS的日志bug，增加KRDS绑定EIP的参数示例


## 1.3.50 (Sep 14, 2022)

BUG FIXES:

- 修正KRDS磁盘参数校验的正则


## 1.3.49 (Sep 6, 2022)

BUG FIXES:

- 修复云主机示例不能更新安全组的问题（不更新其他网络属性，只更新安全组）

## 1.3.48 (Sep 6, 2022)

BUG FIXES:

- releaser配置增加windows和arm的配置


## 1.3.47 (Sep 6, 2022)

BUG FIXES:

- 修复import子网不能获取到AZ属性的问题




## 1.1.0 (Dec 21,2020)

IMPROVEMENTS:

- instance support public_ip and search ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- lb support name_grex ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- icmp_code , icmp_type and weight support zero value ([#20](https://github.com/terraform-providers/terraform-provider-ksyun/issues/20))
- instance remove the function of creating data_disk and support show data_disk ([#22](https://github.com/terraform-providers/terraform-provider-ksyun/issues/22))
- instance,eip and lb support project ([#26](https://github.com/terraform-providers/terraform-provider-ksyun/issues/26))
- add attach and release eip functions for krds instance ([#27](https://github.com/terraform-providers/terraform-provider-ksyun/issues/27))

FEATURES:

- - **New Resource:** `ksyun_lb_rule` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Data Source:** `ksyun_lb_rules` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Resource:** `ksyun_lb_backend_server_group` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Data Source:** `ksyun_lb_backend_server_groups` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Resource:** `ksyun_lb_host_header` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Data Source:** `ksyun_lb_host_headers` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Resource:** `ksyun_lb_register_backend_server` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))
- - **New Data Source:** `ksyun_lb_register_backend_servers` ([#16](https://github.com/terraform-providers/terraform-provider-ksyun/issues/16))

## 1.0.1 (July 01, 2020)

IMPROVEMENTS:

- data images support cloudinit ([#8](https://github.com/terraform-providers/terraform-provider-ksyun/issues/8))
- instance support update instance type ([#11](https://github.com/terraform-providers/terraform-provider-ksyun/issues/11))

BUG FIXES:

- fix eip read crash ([#6](https://github.com/terraform-providers/terraform-provider-ksyun/issues/6))
- fix instance test ([#6](https://github.com/terraform-providers/terraform-provider-ksyun/issues/6))
- fix redis param ([#10](https://github.com/terraform-providers/terraform-provider-ksyun/issues/10))
- fix instance update instance password ([#11](https://github.com/terraform-providers/terraform-provider-ksyun/issues/11))
- fix instance update instance type ([#12](https://github.com/terraform-providers/terraform-provider-ksyun/issues/12))
- using fixed schema to save db instance instead of all the data in response ([#28](https://github.com/terraform-providers/terraform-provider-ksyun/issues/28))
- fix bugs of using error az param when changing password for or renaming a redis instance ([#29](https://github.com/terraform-providers/terraform-provider-ksyun/issues/29))

## 1.0.0 (May 20, 2020)

FEATURES:

### KEC

RESOURCES:

* instance create
* instance read
* instance update
    * reset instance
    * reset password
    * reset keyid
    * update instance name
    * update host name
    * update security groups
    * update network interface
* instance delete

DATA SOURCES:

* image read
* instance read

### VPC

RESOURCES:

* vpc create
* vpc read
* vpc update (update vpc_name)
* vpc delete
* subnet create
* subnet read
* subnet update (update subnet_name,dns1,dns2)
* subnet delete
* security group create
* security group read
* security group update (update security_group_name)
* security group delete
* security group entry create
* security group entry read
* security group entry delete


DATA SOURCES:

* vpc read
* subnet read
* security group read
* subnet allocated ip addresses read
* subnet available addresses read
* network interface read


### EIP

RESOURCES:

* eip create
* eip read
* eip update (update band_width)
* eip delete
* associate address
* disassociate address

DATA SOURCES:

* eip read
* line read

### KCM

RESOURCES:

* certificate create
* certificate read
* certificate update (update certificate_name)
* certificate delete

DATA SOURCES:

* certificate read

### SLB

RESOURCES:

* health check create
* health check read
* health check update (update health_check_state,healthy_threshold,interval,timeout,unhealthy_threshold,is_default_host_name,host_name,url_path)
* health check delete
* lb create
* lb read
* lb update (update load_balancer_name,load_balancer_state)
* lb delete
* lb acl create
* lb acl read
* lb acl update (update load_balancer_acl_name)
* lb acl delete
* lb acl entry create
* lb acl entry read
* lb acl entry delete
* lb listener create
* lb listener read
* lb listener update (update certificate_id,listener_name,listener_state,method)
* lb listener delete
* lb listener server create
* lb listener server read
* lb listener server delete
* lb listener associate acl create
* lb listener associate acl read
* lb listener associate acl delete

DATA SOURCES:

* lb read
* lb health check read
* lb acl read
* lb listener read
* lb listener server read

### EBS

RESOURCES:

* volume create
* volume read
* volume update (update name,volume_desc,size)
* volume delete
* volume attach create
* volume attach read
* volume attach delete

DATA SOURCES:

* volume read

### KRDS

RESOURCES:
* krds create
* krds read
* krds update (update name,class,type,version,password,security_group,preferred_backup_time)
* krds delete
* krds read replica create
* krds read replica read
* krds read replica delete
* krds security group create
* krds security group read
* krds security group update (update name,security_group_description,security_group_rule)
* krds security group delete
* krds sqlserver create
* krds sqlserver read
* krds sqlserver delete

DATA SOURCES:

* krds read
* krds security groups read
* krds sqlservers read

### MONGODB

RESOURCES:

* mongodb instance create
* mongodb instance read
* mongodb instance update (update name,node_num)
* mongodb instance delete
* mongodb security rule create
* mongodb security rule read
* mongodb security rule update (update cidrs)
* mongodb security rule delete
* mongodb shard instance create
* mongodb shard instance read
* mongodb shard instance delete

DATA SOURCES:

* mongodb instance read

### KCS

RESOURCES:

* redis instance create
* redis instance read
* redis instance update (update name,pass_word,capacity)
* redis instance delete
* ~~redis security rule create~~
* ~~redis security rule read~~
* ~~redis security rule update (update rules)~~
* ~~redis security rule delete~~
* redis security group create
* redis security group read
* redis security group update (update name, description)
* redis security group delete
* redis security group allocate instance
* redis security group deallocate instance
* redis security group read allocate instance
* redis security group rule create
* redis security group rule read
* redis security group rule update (update rules)
* redis security group rule delete

DATA SOURCES:

* redis read
* redis security group read

###  IAM

DATA SOURCES:

* iam users read
* iam roles read
* iam groups read

RESOURCES:

* iam user create
* iam user read
* iam user delete
* iam role create
* iam role read
* iam role delete
* iam group create
* iam group read
* iam group delete
* iam policy create
* iam policy read
* iam policy delete
