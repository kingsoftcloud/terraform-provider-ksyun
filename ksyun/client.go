package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/service/bws"
	"github.com/KscSDK/ksc-sdk-go/service/cen"
	"github.com/KscSDK/ksc-sdk-go/service/clickhouse"
	"github.com/KscSDK/ksc-sdk-go/service/ebs"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/KscSDK/ksc-sdk-go/service/iam"
	"github.com/KscSDK/ksc-sdk-go/service/kce"
	"github.com/KscSDK/ksc-sdk-go/service/kcev2"
	"github.com/KscSDK/ksc-sdk-go/service/kcm"
	"github.com/KscSDK/ksc-sdk-go/service/kcrs"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv2"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/KscSDK/ksc-sdk-go/service/knad"
	"github.com/KscSDK/ksc-sdk-go/service/kpfs"
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/KscSDK/ksc-sdk-go/service/mongodb"
	"github.com/KscSDK/ksc-sdk-go/service/monitor"
	"github.com/KscSDK/ksc-sdk-go/service/monitorv4"
	"github.com/KscSDK/ksc-sdk-go/service/pdns"
	"github.com/KscSDK/ksc-sdk-go/service/rabbitmq"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/slb"
	"github.com/KscSDK/ksc-sdk-go/service/sqlserver"
	"github.com/KscSDK/ksc-sdk-go/service/tag"
	"github.com/KscSDK/ksc-sdk-go/service/tagv2"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	klog "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/klog/v20200731"
	kmr "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/kmr/v20210902" // 别名导入kmr SDK
	"github.com/ks3sdklib/ksyun-ks3-go-sdk/ks3"
)

type KsyunClient struct {
	region         string                 `json:"region,omitempty"`
	dryRun         bool                   `json:"dry_run,omitempty"`
	eipconn        *eip.Eip               `json:"eipconn,omitempty"`
	slbconn        *slb.Slb               `json:"slbconn,omitempty"`
	vpcconn        *vpc.Vpc               `json:"vpcconn,omitempty"`
	kecconn        *kec.Kec               `json:"kecconn,omitempty"`
	sqlserverconn  *sqlserver.Sqlserver   `json:"sqlserverconn,omitempty"`
	krdsconn       *krds.Krds             `json:"krdsconn,omitempty"`
	kcmconn        *kcm.Kcm               `json:"kcmconn,omitempty"`
	sksconn        *sks.Sks               `json:"sksconn,omitempty"`
	kcsv1conn      *kcsv1.Kcsv1           `json:"kcsv_1_conn,omitempty"`
	kcsv2conn      *kcsv2.Kcsv2           `json:"kcsv_2_conn,omitempty"`
	epcconn        *epc.Epc               `json:"epcconn,omitempty"`
	ebsconn        *ebs.Ebs               `json:"ebsconn,omitempty"`
	mongodbconn    *mongodb.Mongodb       `json:"mongodbconn,omitempty"`
	ks3conn        *ks3.Client            `json:"ks_3_conn,omitempty"`
	iamconn        *iam.Iam               `json:"iamconn,omitempty"`
	rabbitmqconn   *rabbitmq.Rabbitmq     `json:"rabbitmqconn,omitempty"`
	bwsconn        *bws.Bws               `json:"bwsconn,omitempty"`
	tagconn        *tagv2.Tagv2           `json:"tagconn,omitempty"`
	tagv1conn      *tag.Tag               `json:"tagv1conn,omitempty"`
	kceconn        *kce.Kce               `json:"kceconn,omitempty"`
	kcev2conn      *kcev2.Kcev2           `json:"kcev2conn,omitempty"`
	knadconn       *knad.Knad             `json:"knadconn,omitempty"`
	pdnsconn       *pdns.Pdns             `json:"pdnsconn,omitempty"`
	kcrsconn       *kcrs.Kcrs             `json:"kcrsconn,omitempty"`
	kpfsconn       *kpfs.Kpfs             `json:"kpfsconn,omitempty"`
	monitorconn    *monitor.Monitor       `json:"monitorconn,omitempty"`
	monitorv4conn  *monitorv4.Monitorv4   `json:"monitor_4_conn,omitempty"`
	cenconn        *cen.Cen               `json:"cenconn,omitempty"`
	clickhouseconn *clickhouse.Clickhouse `json:"clickhouseconn,omitempty"`
	kmrconn        *kmr.Client            `json:"kmrconn,omitempty"`
	klogconn       *klog.Client           `json:"klogconn,omitempty"`

	config *Config
}

func (client *KsyunClient) GetVpcClient() *vpc.Vpc {
	return client.vpcconn
}

func (client *KsyunClient) GetKecClient() *kec.Kec {
	return client.kecconn
}

func (client *KsyunClient) GetEipClient() *eip.Eip {
	return client.eipconn
}

func (client *KsyunClient) GetIamClient() *iam.Iam {
	return client.iamconn
}
