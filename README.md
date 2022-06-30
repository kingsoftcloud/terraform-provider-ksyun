# terraform-provider-ksyun
Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.11.x
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-ksyun`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-ksyun
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ksyun
$ go build
```

Using the provider
----------------------

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-ksyun
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to run the single source of Acceptance tests, you can run them by entering the following instructions in a terminal:.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ cd ksyun
$ export TF_ACC=true
$  go test -test.run TestAccKsyunEip_basic -v
```

# 中文版介绍
该介绍包括三部分：
##### terraform-provider-ksyun开发

本地开发调试方法

针对terraform v0.12版本：

1、在项目根目录下执行 `make dev_v0.12`，会编译provider到$HOME/.terraform.d/plugins

2、到tf文件所在目录执行`terraform init`会获 取到 该本地版本的provider

针对terraform v0.13及以上版本：

1、在项目根目录下执行 `make dev version=1.3.4`，会编译provider到$HOME/.terraform.d/plugin-cache（**注意**:调命令的时候需要给编译的provider传入一个x.y.z的版本号，以便terraform可以加载到该文件）

2、 到tf文件所在目录执行`terraform init --plugin-dir ~/.terraform.d/plugin-cache/`会获取到本地版本的provider(**注意**:高版本的terraform需要指定provider所在路径，这里通过参数形式指定，其他指定方式可以参考官方开发文档：https://www.terraform.io/cli/config/config-file)


##### terraform-provider-ksyun使用

_云产品用户参考。_

##### terraform-provider-ksyun 属性介绍

_各产品线开发人员负责补充，云产品用户参考。_

## terraform-provider-ksyun开发
### 开发指南（以eip为例）
##### 1、client.go 各产品请求连接声明（eipconn *eip.Eip）
##### 2、config.go 各产品请求连接定义（client.eipconn = eip.SdkNew(cli, cfg, url)）
##### 3、provider.go 配置dataSource和resource
      配置dataSource: "ksyun_eips":dataSourceKsyunEips()
      配置resource: "ksyun_eip":resourceKsyunEip()
##### 4、data_source_ksyun_eips.go dataSource具体实现（根据具体过滤条件拉取eip列表）
##### 5、resource_ksyun_eip.go resource具体实现（单个eip的增删改查）
##### 6、添加文档和链接
      1、添加对应的文档website/docs/d/${datasource}s.html.markdown和website/docs/r/${resource}.html.markdown（注意：文档命名必须遵守格式 ${datasource}s.html.markdown和${resource}.html.markdown）
      2、添加产品链接目录pronamespace
      3、根据对应的操作系统，执行相应的文档生成程序terraform-index-build-*
      4、检查生成的website/docs/index.html.markdown文件是否符合预期

### 开发注意事项
##### 1、所有的入参和出参必须在schema.Resource中定义，否则terraform无法识别。
##### 2、terraform的出参和入参都是schema.Schema类型，其底层调用的sdk的出参和入参是map[string]interface{}类型，两者需进行转换。
##### 3、异步创建需要添加状态轮询（参照主机创建resource_ksyun_instance.go）。
##### 4、所有的删除都需添加重试机制（防止有资源依赖时删除失败）。
##### 5、所有的修改最好添加Partial机制（防止有些属性同步有问题，多个修改接口或查询接口不展示全部属性会有问题）。
##### 6、网络和主机的接口属性较多，开发采用了函数封装，初次开发，可先参考腾讯的开发指南，更有易读性，了解terraform的开发过程后可直接调用封装好的函数，减少代码冗余。
      腾讯云开发指南：（https://cloud.tencent.com/developer/article/1067230）
      *Note:* 腾讯的sdk入参和出参是class（struct）类型，而我们的sdk是map类型。
##### 7、每次提交都需要对代码进行交叉编译(mac、linux和windows),编译好的可执行文件需传到公司官网，请将编译文件传给庞雄伟(雄伟同学负责我们terraform的官网文档)，不然用户没法使用。
##### 8、每次编译前都需要重新拉取terraform-provider-ksyun和ksc-sdk-go,保证terraform-provider-ksyun和ksc-sdk-go都是最新版，不然会导致其他服务不可用。
##### 9、所有的*test.go文件里的参数必须是线上可以直接运行的，所有依赖的资源必须是新建的。必须保证test是可以跑通的，test见上面英文版test介绍。不然terraform 官网没法验证。
##### 10、如果底层open API的入参，零值（int(0)、float(0)、boolean(false)及string("")）和不传代表不同含义或者需要显示透传零值时，需要单独处理（例lbListenrServer 中weight传0时为0，不传时为1）。
      terraform 在读取main.tf 时，若配置了零值或不配置，getok都返回false，所以无法对两者进行区分，如果通过getok进行判断传参，两者都不会传给open API。
      解决方案：设置默认值进行区分，但创建和更新接口需要对参数做单独处理
      1、若只是需要显示透传零值，terraform里字段的默认值设为open API的默认值（参见lbListenrServer 中weight，注意weight=0在实际应用中无意义，此处只是作为示范开发）
      2、若零值和不传代表不同含义，terraform 里字段的默认值需设为open API不支持的字段（参见securityGroupEntry 中icmp_code 和icmp_type）
            
### 提交注意事项
##### 1、不要上传ak、sk等敏感信息
##### 2、new pull request 之前请确保自己fork的是master的最新版本(在kscSDK(master)项目下点击new pull request 到自己的项目)，不然覆盖其他产品线的提交，后果自负。
##### 3、new pull request 之后，请发邮件或在群里说一下，不然没人知道。
##### 4、new pull request fork请forkmaster，merge请提到trunk分支。
##### 5、terraform 对代码的格式和风格有严格要求，以下几点需要运行make file 文件（mac和linux可直接执行，windows自行百度）。
        a、make fmt（terraform 要求代码是go fmt之后的，mac和linux下开发直接在项目根目录下执行make fmt，windows需要针对你修改的每一个文件执行gofmt -w -s $(GOFMT_FILES)）。
        b、make tools (下载两个工具包，若下载不下来，可通过码云极速下载。)
        c、make lint（会执行golangci-lint run,需对报错代码进行修改(修改原则凡是出现的代码，必须有意义)。以下几种报错若看不懂可直接百度。）
              1、deadcode、unused:没有用到的代码，需要删除或注释掉。
              2、errcheck:函数func返回类型为error类型，需要对返回值进行处理。
              3、ineffassign：无效赋值(忽略没用到的返回值)。
              4、staticcheck：返回值需要进行处理。
              5、gosimple：代码需要进行简化处理。
              6、govet：可能的bug或者可疑的构造。

## terraform-provider-ksyun使用

### 快速安装


  若用户只是利用插件，不对插件进行二次开发，且不想安装go环境和编译代码，可直接安装编译好的插件。不同操作系统的插件（terraform-provider-ksyun）位于目录bin/下，将其直接解压拷贝至terraform的plugins默认目录即可。
不同操作系统terraform的plugins默认目录：(https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)。
>注意：首次拷贝，需手动创建plugins文件夹。

### Terraform-provider-ksyun 文件配置：

 #### 1、provider配置

 	provider "ksyun" {
       access_key = "你的ak"
       secret_key = "你的sk"
       region = "cn-beijing-6"
     }
 	
  可以在配置文件中指定，也可以在环境变量中配置，若两处都配置，以配置文件为主。
  在环境变量中配置：
  
```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ksyun
$ export KSYUN_ACCESS_KEY=xxx
$ export KSYUN_SECRET_KEY=xxx
$ export KSYUN_REGION=xxx
$ export TF_LOG=DEBUG
$ export TF_LOG_PATH=/var/log/terraform.log
$ export OS_DEBUG=1
```
#### 2、data-source使用
  以data+产品线命名的文件夹，代表不同资源，可根据配置条件导出全部符合的资源。若配置条件为空，导出全部。
  以dataEips下的main.tf为例：

 	data "ksyun_eips" "default" {
        //导出的资源会输出到output_result文件中
      output_file = “output_result” 
        //只导出eipId包含在ids中的eip信息
      ids = []
       //只导出project_id=1的eip默认只导出project_id=0的eip
      project_id = [“1"]
       //只导出instance_type=“Ipfwd”的eip
      instance_type = ["Ipfwd"]
      network_interface_id = []
      internet_gateway_id = []
      band_width_share_id = []
      line_id = []
      public_ip = []
    }
  在该目录下执行：
```sh
$ terraform init
$ terraform plan
```
  会将符合条件的eip导出到output_result中
#### 3、resource使用
  以产品线命名的文件夹代表对应产品线的单资源配置，以产品线+Service命名的文件夹代表对应该资源的服务编排。
##### 单资源配置：
  以eips下的main.tf为例：

 	resource "ksyun_eip" "default1" {
      line_id ="cf8b7b95-4651-b96c-db67-b38336f2fe70"
      band_width =1
      charge_type = "PostPaidByDay"
      purchase_time =1
      project_id=0
    }

  在该目录下执行：
```sh
$ terraform init
$ terraform plan //获取操作类型（新建，修改或重建）
$ terraform apply //执行操作
$ terraform destroy //删除eip
$ terraform import ksyun_eip.default1 eipId //导入该eipId的eip信息，一般用于对已有实例的修改
```

##### 服务编排：
  以terraform-provider-ksyun/example/instanceService为例：

###### 1、variables.tf定义变量

 	variable "instance_name" {
      default = "ksyun_instance_tf"
    }
    variable "subnet_name" {
      default = "ksyun_subnet_tf"
    }

  定义变量instance_name，其默认值为“ksyun_instance_tf”
  若用户想自定义变量的值，可在执行terraform 命令时指定:
```sh
$ terraform plan -var ‘instance_name=kec’ -var ‘subnet_name=sub’
```
###### 2、outputs.tf控制台输出

 	output "eip_id" {
      value = "${ksyun_eip.default.id}"
    }


  执行terraform apply 后，控制台会输出：
```sh
$ eip_id=‘e9587b84-0da7-4fd7-a26d-bc56df63b01e’
```
###### 3、main.tf 资源编排

  1、定义了provider为ksyun

  2、定义了三个资源，即镜像、线路和可用区

  3、定义了多个资源配置，包括创建vpc、子网、安全组、安全组规则、主机及绑定eip

### Terraform-provider-ksyun 版本升级：
  terraform v0.12 版本的配置文件与v0.11 版本的配置文件格式不同。该example下的配置文件是基于v0.11.13开发的。若想使用v0.12版本的terraform需对配置文件进行修改，该修改不需手动修改，terraform支持自动修改。可在配置文件目录下直接执行：
```sh
$ terraform 0.12upgrade
```
  terraform会询问是否确认修改，输入yes即可。
  
### Terraform-provider-ksyun 属性介绍：
 
 1、.tf文件的参数属性请直接参考官网openapi的接口介绍。.tf文件中的属性字段一般都可以在openapi中找到。
 
 2、terraform-provider-ksyun 尽量保持了原子性，openapi创建接口里若出现同时创建多个资源的情况，provider是不支持的，请分别配置。
   
_例：官网openapi里主机创建的接口里，可以同时创建eip和主机，在terraform里是不支持的。主机和eip需单独配置。_
 
##### 下面只介绍官网openapi文档和tf配置文件不一致的资源。
  
######  云主机
1、不支持单个resource(ksyun_instance)批量创建主机，即不支持openapi文档里的MaxCount，MinCount，InstanceNameSuffix。

2、官网openapi主机创建时，SecurityGroupId 目前仅支持绑定一个安全组，terraform 里可配置多个。

######  redis
1、实例参数配置功能在ksyun_redis_instance资源的parameters中配置。

2、实例主从模式节点不支持并行批量添加, 顺序批量添加多个节点, 需要在ksyun_redis_instance_node资源中配置pre_node_id属性, 属性值是上一个创建的节点ID。

######  mongodb
1、副本集实例增加Second节点功能在ksyun_mongodb_instance资源中的node_num属性配置。

# 金山云业务对应Terraform的Resource和DataSource

| 资源名          | terraform(Resource)    | terraform(Data) | 资源分类
|--------------| -------  | ---- | ----
| 弹性IP         | ksyun_eip | ksyun_eips | eip
| 链路           | Not_Support | ksyun_lines | eip
| 弹性IP绑定和解绑    | ksyun_eip\_associate | Not_Support | eip
| 云物理机         | ksyun_epc | ksyun_epcs | epc
| 证书           | ksyun_certificate | ksyun_certificates | kcm
| 健康检查         | ksyun_lb\_healthcheck | ksyun_lb\_healthchecks | slb
| 负载均衡         | ksyun_lb | ksyun_lbs | slb
| 负载均衡访问控制列表   | ksyun_lb\_acl | ksyun_lb\_acls | slb
| 负载均衡访问控制列表规则 | ksyun_lb\_acl\_entry| Not_Support | slb
| 健康检查         | ksyun_healthcheck | ksyun_healthchecks | slb
| 监听器          | ksyun_lb\_listener | ksyun_lb\_listeners | slb
| 真实服务器        | ksyun_lb\_listener\_server | ksyun_lb\_listener\_servers | slb
| 监听器绑定访问控制列表  | ksyun_lb\_listener\_associate\_acl | Not_Support | slb
| 云主机          | ksyun_instance | ksyun_instances | kec
| 云主机镜像        | Not_Support | ksyun_images | kec
| 云盘           | ksyun_volume | ksyun_volumes | ebs
| 云盘绑定         | ksyun_volume_attach | Not_Support | ebs
| RDS          | ksyun_krds | ksyun_krds | krds
| RDS只读实例      | ksyun_krds\_read\_replica | ksyun_krds | krds
| RDS安全组       | ksyun_krds\_security\_group | ksyun_krds\_security\_groups | krds
| SqlServer    | ksyun_sqlserver | ksyun_sqlservers | krds
| MongoDB实例    | ksyun_mongodb\_instance | ksyun_mongodb | mongodb
| MongoDB安全组   | ksyun_mongodb\_security\_rule | ksyun_mongodb | mongodb
| MongoDB实例分片  | ksyun_mongodb\_shard\_instance | ksyun_mongodb | mongodb
| Redis实例      | ksyun_redis\_instance | ksyun_redis | kcs
| Redis节点      | ksyun_redis\_instance\_node | ksyun_redis | kcs
| Redis安全组规则   | ksyun_redis\_sec\_rule | ksyun_redis | kcs
| 安全组          | ksyun_security\_group | ksyun_security\_groups | vpc
| 安全组规则        | ksyun_security\_group\_entry | ksyun_security\_groups | vpc
| 虚拟网卡         | Not_Support | ksyun_network\_interface | vpc
| 子网           | ksyun_subnet | ksyun_subnets | vpc
| 子网已用IP       | Not_Support | ksyun_subnet\_allocated\_ip\_addresses | vpc
| 子网可用IP       | Not_Support | ksyun_subnet\_available\_addresses | vpc
| 虚拟私有网络       | ksyun_vpc | ksyun_vpcs | vpc
| 登录SSHKEY     | ksyun_ssh\_key | ksyun_ssh\_keys | sks
| 对象存储         | ksyun_ks3 | Not_Support | ks3
|  标签管理        | ksyun_tag | ksyun_tags | tag
