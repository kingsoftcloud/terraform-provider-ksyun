package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/bws"
	"github.com/KscSDK/ksc-sdk-go/service/ebs"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/KscSDK/ksc-sdk-go/service/iam"
	"github.com/KscSDK/ksc-sdk-go/service/kce"
	"github.com/KscSDK/ksc-sdk-go/service/kcev2"
	"github.com/KscSDK/ksc-sdk-go/service/kcm"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv2"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/KscSDK/ksc-sdk-go/service/mongodb"
	"github.com/KscSDK/ksc-sdk-go/service/rabbitmq"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/slb"
	"github.com/KscSDK/ksc-sdk-go/service/sqlserver"
	"github.com/KscSDK/ksc-sdk-go/service/tag"
	"github.com/KscSDK/ksc-sdk-go/service/tagv2"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	"github.com/ks3sdklib/aws-sdk-go/aws"
	"github.com/ks3sdklib/aws-sdk-go/aws/credentials"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"
)

// Config is the configuration of ksyun meta data
type Config struct {
	AccessKey     string
	SecretKey     string
	Region        string
	Insecure      bool
	Domain        string
	DryRun        bool
	IgnoreService bool
}

// Client will returns a client with connections for all product
func (c *Config) Client() (*KsyunClient, error) {
	var client KsyunClient
	//init ksc client info
	client.region = c.Region
	cli := ksc.NewClient(c.AccessKey, c.SecretKey)

	// 重试去掉
	var MaxRetries int = 0
	cli.Config.MaxRetries = &MaxRetries
	cfg := &ksc.Config{
		Region: &c.Region,
	}
	url := &utils.UrlInfo{
		UseSSL:                      false,
		Locate:                      false,
		CustomerDomain:              c.Domain,
		CustomerDomainIgnoreService: c.IgnoreService,
	}

	client.dryRun = c.DryRun
	client.vpcconn = vpc.SdkNew(cli, cfg, url)
	client.eipconn = eip.SdkNew(cli, cfg, url)
	client.slbconn = slb.SdkNew(cli, cfg, url)
	client.kecconn = kec.SdkNew(cli, cfg, url)
	client.sqlserverconn = sqlserver.SdkNew(cli, cfg, url)
	client.krdsconn = krds.SdkNew(cli, cfg, url)
	client.kcmconn = kcm.SdkNew(cli, cfg, url)
	client.sksconn = sks.SdkNew(cli, cfg, url)
	client.kcsv1conn = kcsv1.SdkNew(cli, cfg, url)
	client.kcsv2conn = kcsv2.SdkNew(cli, cfg, url)
	client.epcconn = epc.SdkNew(cli, cfg, url)
	client.ebsconn = ebs.SdkNew(cli, cfg, url)
	client.mongodbconn = mongodb.SdkNew(cli, cfg, url)
	client.iamconn = iam.SdkNew(cli, cfg, url)
	client.rabbitmqconn = rabbitmq.SdkNew(cli, cfg, url)
	client.bwsconn = bws.SdkNew(cli, cfg, url)
	client.tagconn = tagv2.SdkNew(cli, cfg, url)
	client.tagv1conn = tag.SdkNew(cli, cfg, url)
	client.kceconn = kce.SdkNew(cli, cfg, url)
	client.kcev2conn = kcev2.SdkNew(cli, cfg, url)

	credentials := credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	client.ks3conn = s3.New(&aws.Config{
		Region:           "BEIJING",
		Credentials:      credentials,
		Endpoint:         c.Region,
		DisableSSL:       true,
		LogLevel:         1,
		S3ForcePathStyle: true,
		LogHTTPBody:      true,
	})
	return &client, nil
}
