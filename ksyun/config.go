package ksyun

import (
	"fmt"
	klog "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/klog/v20200731"
	kmr "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/kmr/v20210902"
	"github.com/kingsoftcloud/sdk-go/v2/ksyun/common"
	"github.com/kingsoftcloud/sdk-go/v2/ksyun/common/profile"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ks3sdklib/ksyun-ks3-go-sdk/ks3"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/network"
)

// Config is the configuration of ksyun meta data
type Config struct {
	AccessKey     string
	SecretKey     string
	Region        string
	Insecure      bool
	Domain        string
	Endpoint      string
	DryRun        bool
	IgnoreService bool
	HttpKeepAlive bool
	MaxRetries    int
	HttpProxy     string
	UseSSL        bool
}

// Client will returns a client with connections for all product
func (c *Config) Client() (*KsyunClient, error) {
	var client KsyunClient
	var err error
	// init ksc client info
	client.region = c.Region
	cli := ksc.NewClient(c.AccessKey, c.SecretKey)

	registerClient(cli, c)
	// 重试去掉
	var MaxRetries = c.MaxRetries
	cli.Config.MaxRetries = &MaxRetries
	cfg := &ksc.Config{
		Region: &c.Region,
	}
	client.config = c
	url := &utils.UrlInfo{
		UseSSL:                      c.UseSSL,
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
	client.knadconn = knad.SdkNew(cli, cfg, url)
	client.pdnsconn = pdns.SdkNew(cli, cfg, url)
	client.kcrsconn = kcrs.SdkNew(cli, cfg, url)
	client.kpfsconn = kpfs.SdkNew(cli, cfg, url)
	if client.klogconn, err = klogSdkNew(c); err != nil {
		return nil, err
	}
	client.clickhouseconn = clickhouse.SdkNew(cli, cfg, url)
	client.monitorconn = monitor.SdkNew(cli, cfg, url)
	client.monitorv4conn = monitorv4.SdkNew(cli, cfg, url)
	client.cenconn = cen.SdkNew(cli, cfg, url)

	// 懒加载ks3-client 所以不在此初始
	return &client, nil
}

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe
var loadSdkfromRemoteMutex = sync.Mutex{}
var loadSdkEndpointMutex = sync.Mutex{}
var tagsMutex = sync.Mutex{}

func (client *KsyunClient) WithKs3BucketByName(bucketName string, do func(*ks3.Bucket) (interface{}, error)) (interface{}, error) {
	return client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		bucket, err := client.ks3conn.Bucket(bucketName)
		if err != nil {
			return nil, fmt.Errorf("unable to get the bucket %s: %#v", bucketName, err)
		}
		return do(bucket)
	})
}
func (client *KsyunClient) WithKs3Client(do func(*ks3.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()
	// Initialize the KS3 client if necessary
	if client.ks3conn == nil {
		ks3conn, err := ks3.New(client.config.Endpoint, client.config.AccessKey, client.config.SecretKey)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the KS3 client: %#v", err)
		}
		client.ks3conn = ks3conn
	}
	return do(client.ks3conn)
}

func registerClient(cli *session.Session, c *Config) {

	// register http client
	httpClient := getKsyunClient(c)
	cli.Config.WithHTTPClient(httpClient)

	cli.Config.Retryer = network.GetKsyunRetryer(c.MaxRetries)

	cli.Handlers.CompleteAttempt.PushBackNamed(network.NetErrorHandler)

	// TODO: output request's information when it encounters the special error that reset connection.
	// cli.Handlers.CompleteAttempt.PushBackNamed(network.OutputResetError)

	cli.Handlers.Sign.PushBackNamed(network.HandleRequestBody)
}

func getKsyunClient(c *Config) *http.Client {
	tp := &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			if c.HttpProxy != "" {
				proxyUrl, err := url.Parse(c.HttpProxy)
				if err != nil {
					return nil, fmt.Errorf("error parsing HTTP proxy URL: %w", err)
				}
				return http.ProxyURL(proxyUrl)(r)
			}
			return http.ProxyFromEnvironment(r)
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // dial connect timeout
			KeepAlive: 30 * time.Second, // the interval probes time between the ends of network
		}).DialContext,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     !c.HttpKeepAlive,
	}

	httpClient := &http.Client{
		Timeout:   3 * time.Minute, // a completed request, includes tcp connect, received response, elapsed time.
		Transport: tp,
	}
	return httpClient
}

func klogSdkNew(c *Config) (*klog.Client, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 20
	cpf.HttpProfile.Endpoint = c.Endpoint
	return klog.NewClient(common.NewCredential(c.AccessKey, c.SecretKey), c.Region, cpf)
}

func (client *KsyunClient) WithKmrClient(do func(*kmr.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()
	// Initialize the KMR client if necessary
	if client.kmrconn == nil {
		credential := common.NewCredential(client.config.AccessKey, client.config.SecretKey)
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = client.config.Endpoint
		kmrconn, err := kmr.NewClient(credential, client.config.Region, cpf)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the KMR client: %#v", err)
		}
		client.kmrconn = kmrconn
	}
	return do(client.kmrconn)
}
