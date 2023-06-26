## 1.4.1 (Jun 20, 2023)

IMPROVEMENTS:

- 新增mysql参数组独立管理并支持实例使用临时参数(ksyun_krds_parameter_group)
- 新增kec容灾组管理(ksyun_data_guard_group)
- 新增kec自动快照策略管理(ksyun_auto_snapshot_policy)
- 修复kec实例升降配时的重启逻辑(ksyun_instance: instance_type)
  - NOTE: 实例降级操作会触发实例关机

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
