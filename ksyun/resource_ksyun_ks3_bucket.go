/*
Provides an KS3 Bucket resource.

# Example Usage

```hcl
# 创建一个名为 "bucket-new" 的 ks3 存储桶资源

	resource "ksyun_ks3_bucket" "bucket-new" {
	  provider = ks3.bj-prod
	  #指定要创建的虚拟存储桶的名称
	  bucket = "YOUR_BUCKET"
	  #桶权限
	  acl    = "public-read"
	}

	resource "ksyun_ks3_bucket" "bucket-attr" {
	  provider = ksyun.bj-prod
	  bucket = "bucket-20230324-1"
	  #访问日志的存储路径
	  logging {
	    target_bucket = "bucket-20230324-2"
	    target_prefix = "log/"
	  }
	  #文件生命周期
	  lifecycle_rule {
	    id      = "id1"
	    #如果filter.prefix有值情况下会覆盖这个选项
	    prefix  = "path/1"
	    enabled = true

	    #过期删除时间 days or date 只能二选一
	    expiration {
	      #设置过期日期
	      date = "2023-04-10"
	      #设置过期时间
	      #days = "40"
	  }
	  #添加标签
	    filter {
	    prefix = "example_prefix"
	    #每个标签不能与其他标签重复
	    tag {
	      key   = "example_key1"
	      value = "example_value"
	    }
	    tag {
	      key   = "example_key2"
	      value = "example_value"
	    }
	  }

}

	lifecycle_rule {
	  id      = "id2"
	  prefix  = "path/365"
	  enabled = true

	  expiration {
	    date = "2023-04-10"
	  }
	}

#跨域规则

	  cors_rule {
	  allowed_origins =  [var.allow-origins-star]
	  allowed_methods = split(",", var.allow-methods-put)
	  allowed_headers = [var.allowed_headers]
	}

	cors_rule {
	  allowed_origins = split(",", var.allow-origins-ksyun)
	  allowed_methods = split(",", var.allow-methods-get)
	  allowed_headers = [var.allowed_headers]
	  expose_headers  = [var.expose_headers]
	  max_age_seconds = var.max_age_seconds
	}

}
```

# 创建variables.tf (上面var的变量才能引用到)
```hcl

	variable "bucket-new" {
	  default = "bucket-20180423-1"
	}

	variable "bucket-attr" {
	  default = "bucket-20180423-2"
	}

	variable "acl-bj" {
	  default = "private"
	}

	variable "target-prefix" {
	  default = "log3/"
	}

	variable "role-days" {
	  default = "expirationByDays"
	}

	variable "rule-days" {
	  default = 365
	}

	variable "my_variable" {
	  type = number
	  default = 42
	}

	variable "role-date" {
	  default = "expirationByDate"
	}

	variable "rule-date" {
	  default = "2023-03-18"
	}

	variable "rule-prefix" {
	  default = "path"
	}

	variable "allow-empty" {
	  default = true
	}

	variable "referers" {
	  default = "http://www.ksyun.com, https://www.ksyun.com, http://?.ksyun.com"
	}

	variable "allow-origins-star" {
	  default = "*"
	}

	variable "allow-origins-ksyun" {
	  default = "http://www.ksyun.com,http://*.ksyun.com"
	}

	variable "allow-methods-get" {
	  default = "GET"
	}

	variable "allow-methods-put" {
	  default = "PUT,GET"
	}

	variable "allowed_headers" {
	  default = "*"
	}

	variable "expose_headers" {
	  default = "x-ks3-test, x-ks3-test1"
	}

	variable "max_age_seconds" {
	  default = 100
	}

```
*/
package ksyun

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/wilac-pv/ksyun-ks3-go-sdk/ks3"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceKsyunKs3Bucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKs3BucketCreate,
		Read:   resourceKsyunKs3BucketRead,
		Update: resourceKsyunKs3BucketUpdate,
		Delete: resourceKsyunKs3BucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 63),
				Default:      resource.PrefixedUniqueId("ks3-bucket-"),
				Description:  "The name of the bucket. If omitted, Terraform will assign a random, unique name.",
			},

			"acl": {
				Type:         schema.TypeString,
				Default:      ks3.ACLPrivate,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
				Description:  "The canned ACL to apply. Defaults to private. Valid values are private, public-read, and public-read-write.",
			},
			"cors_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The cross domain resource sharing (CORS) configuration rules of the bucket. Each rule can contain up to 100 headers, 100 methods, and 100 origins. For more information, see Cross-domain resource sharing (CORS) in the KS3 Developer Guide.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The HTTP header that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, Authorization. Additionally, you can use \"*\" to represent all headers.",
						},
						"allowed_methods": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The HTTP method that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, GET, PUT, POST, DELETE, HEAD, OPTIONS.",
						},
						"allowed_origins": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The source address that users are allowed to access through cross domain resource sharing, which can contain up to one \"*\" wildcard character. Each CORSRule must define at least one source address and one method. For example, http://*. example. com. Additionally, you can use \"*\" to represent all sources.",
						},
						"expose_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Identify response headers that allow customers to access from applications, such as JavaScript XMLHttpRequest data element.",
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum time in seconds that the browser can cache the preflight response.",
						},
					},
				},
				MaxItems: 10,
			},

			"logging": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Call this interface to set the bucket logging configuration. If the configuration already exists, KS3 will replace it.\nTo use this interface, you need to have permission to perform the ks3: PutBucketLogging operation. The space owner has this permission by default and can grant corresponding permissions to others.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the bucket where you want KS3 to store server access logs. You can have your logs delivered to any bucket that you own, including the same bucket that is being logged. You can also configure multiple buckets to deliver their logs to the same target bucket. In this case, you should assign each bucket a unique prefix.",
						},
						"target_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A prefix for all log object keys. If you store log files from multiple buckets in a single bucket, you can use a prefix to distinguish which log files came from which bucket.",
						},
					},
				},
				MaxItems: 1,
			},

			//"logging_isenable": {
			//	Type:       schema.TypeBool,
			//	Optional:   true,
			//	Deprecated: "Use the logging block instead",
			//},

			"lifecycle_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Call this interface to set the bucket lifecycle configuration. If the configuration already exists, KS3 will replace it.\nTo use this interface, you need to have permission to perform the ks3: PutBucketLifecycle operation. The space owner has this permission by default and can grant corresponding permissions to others.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
							Description:  "Unique identifier for the rule. The value cannot be longer than 255 characters.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prefix identifying one or more objects to which the rule applies. The rule applies only to objects with key names beginning with the specified prefix. If you specify multiple rules in a lifecycle configuration, the rule with the longest prefix is used to determine which rule applies to a given object. If you specify a prefix for the rule, you cannot specify a tag for the rule.",
						},
						"filter": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Container for the filter of lifecycle rule. If you specify a filter, you cannot specify a prefix for the rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Prefix identifying one or more objects to which the rule applies. The rule applies only to objects with key names beginning with the specified prefix. If you specify multiple rules in a lifecycle configuration, the rule with the longest prefix is used to determine which rule applies to a given object. If you specify a prefix for the rule, you cannot specify a tag for the rule.",
									},
									"tag": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Container for the tag of lifecycle rule. If you specify a tag, you cannot specify a prefix for the rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The key of the tag.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of the tag.",
												},
											},
										},
									},
								},
							},
						},

						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If 'Enabled', the rule is currently being applied. If 'Disabled', the rule is not currently being applied.",
						},
						"expiration": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Specifies when an object transitions to a specified storage class. If you specify multiple rules in a lifecycle configuration, the rule with the earliest expiration date is applied to the object. If you specify multiple transition rules for the same object, the rule with the earliest date is applied.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Indicates at what date the object is to be moved or deleted.example:2016-01-02.",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Indicates the lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer.",
									},
								},
							},
						},
						"transitions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Specifies when an object transitions to a specified storage class. If you specify multiple rules in a lifecycle configuration, the rule with the earliest transition date is applied to the object. If you specify multiple transition rules for the same object, the rule with the earliest date is applied.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Indicates at what date the object is to be moved or deleted.example:2016-01-02.",
									},
									"days": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Indicates the lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer.",
									},
									"storage_class": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(ks3.TypeNormal),
											string(ks3.StorageIA),
											string(ks3.StorageArchive),
											string(ks3.StorageDeepIA),
										}, false),
										Description: "The class of storage used to store the object.",
									},
								},
							},
						},
					},
				},
				MaxItems: 10,
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ks3.TypeNormal),
				ValidateFunc: validation.StringInSlice([]string{
					string(ks3.TypeNormal),
					string(ks3.TypeIA),
					string(ks3.TypeArchive),
					string(ks3.TypeDeepIA),
				}, false),
				Description: "The class of storage used to store the object.",
			},
		},
	}
}

func resourceKsyunKs3BucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	request := map[string]string{"bucketName": d.Get("bucket").(string)}
	var requestInfo *ks3.Client
	type Request struct {
		BucketName         string
		StorageClassOption ks3.Option
		AclTypeOption      ks3.Option
	}

	req := Request{
		d.Get("bucket").(string),
		ks3.BucketTypeClass(ks3.BucketType(d.Get("storage_class").(string))),
		ks3.ACL(ks3.ACLType(d.Get("acl").(string))),
	}
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		return nil, ks3Client.CreateBucket(req.BucketName, req.StorageClassOption, req.AclTypeOption)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ksyun_ks3_bucket", "CreateBucket", KsyunKs3GoSdk)
	}
	addDebug("CreateBucket", raw, requestInfo, req)
	d.SetId(request["bucketName"])

	return resourceKsyunKs3BucketUpdate(d, meta)
}

func resourceKsyunKs3BucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	ks3Service := Ks3Service{client}
	object, err := ks3Service.DescribeKs3Bucket(d.Id())
	if err != nil {
		if notFoundErrorNew(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bucket", d.Id())
	d.Set("acl", object.BucketInfo.ACL)
	d.Set("creation_date", object.BucketInfo.CreationDate.Format("2006-01-02"))
	d.Set("location", object.BucketInfo.Location)
	d.Set("owner", object.BucketInfo.Owner.ID)
	d.Set("storage_class", object.BucketInfo.StorageClass)

	request := map[string]string{"bucketName": d.Id()}
	var requestInfo *ks3.Client

	// Read the CORS
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		return ks3Client.GetBucketCORS(request["bucketName"])
	})
	if err != nil && !IsExpectedErrors(err, []string{"NoSuchCORSConfiguration"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketCORS", KsyunKs3GoSdk)
	}
	addDebug("GetBucketCORS", raw, requestInfo, request)
	cors, _ := raw.(ks3.GetBucketCORSResult)
	rules := make([]map[string]interface{}, 0, len(cors.CORSRules))
	for _, r := range cors.CORSRules {
		rule := make(map[string]interface{})
		rule["allowed_headers"] = r.AllowedHeader
		rule["allowed_methods"] = r.AllowedMethod
		rule["allowed_origins"] = r.AllowedOrigin
		rule["expose_headers"] = r.ExposeHeader
		rule["max_age_seconds"] = r.MaxAgeSeconds

		rules = append(rules, rule)
	}
	if err := d.Set("cors_rule", rules); err != nil {
		return WrapError(err)
	}

	// Read the logging configuration
	raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		return ks3Client.GetBucketLogging(d.Id())
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", KsyunKs3GoSdk)
	}
	addDebug("GetBucketLogging", raw, requestInfo, request)
	logging, _ := raw.(ks3.GetBucketLoggingResult)

	if &logging != nil {
		enable := logging.LoggingEnabled
		if &enable != nil && len(enable.TargetBucket) > 0 {
			lgs := make([]map[string]interface{}, 0)
			tb := logging.LoggingEnabled.TargetBucket
			tp := logging.LoggingEnabled.TargetPrefix
			if tb != "" || tp != "" {
				lgs = append(lgs, map[string]interface{}{
					"target_bucket": tb,
					"target_prefix": tp,
				})
			}
			if err := d.Set("logging", lgs); err != nil {
				return WrapError(err)
			}
		}
	}

	// Read the lifecycle rule configuration
	raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		return ks3Client.GetBucketLifecycle(d.Id())
	})
	if err != nil && !ks3NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLifecycle", KsyunKs3GoSdk)
	}
	log.Printf("[DEBUG] Ks3 bucket:  %s, raw: %#v", d.Id(), raw)
	addDebug("GetBucketLifecycle", raw, requestInfo, request)
	lrules := make([]map[string]interface{}, 0)
	lifecycle, _ := raw.(ks3.GetBucketLifecycleResult)
	for _, lifecycleRule := range lifecycle.Rules {
		rule := make(map[string]interface{})
		rule["id"] = lifecycleRule.ID
		rule["prefix"] = lifecycleRule.Prefix
		if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
			rule["enabled"] = true
		} else {
			rule["enabled"] = false
		}
		// Expiration
		if lifecycleRule.Expiration != nil {
			log.Printf("[DEBUG] lifecycleRule.Expiration start: %#v", lifecycleRule.Expiration)
			m := make(map[string]interface{})
			if &lifecycleRule.Expiration.Date != nil && lifecycleRule.Expiration.Date != "" {
				lifecycleRule.Expiration.Date = strings.ReplaceAll(lifecycleRule.Expiration.Date, ".000", "")
				t, err := time.Parse(Iso8601DateFormat, lifecycleRule.Expiration.Date)
				if err == nil {
					m["date"] = t.Format("2006-01-02")
				}
			}
			if &lifecycleRule.Expiration.Days != nil {
				m["days"] = lifecycleRule.Expiration.Days
			}
			rule["expiration"] = schema.NewSet(expirationHash, []interface{}{m})
		}
		// transitions
		if len(lifecycleRule.Transitions) != 0 {
			var eSli []interface{}
			for _, transition := range lifecycleRule.Transitions {
				e := make(map[string]interface{})
				if &transition.Date != nil && transition.Date != "" {
					transition.Date = strings.ReplaceAll(transition.Date, ".000", "")
					t, err := time.Parse(Iso8601DateFormat, transition.Date)
					if err != nil {
						return WrapError(err)
					}
					e["date"] = t.Format("2006-01-02")
				}
				e["days"] = fmt.Sprintf("%d", transition.Days)
				e["storage_class"] = string(transition.StorageClass)
				eSli = append(eSli, e)
			}
			rule["transitions"] = schema.NewSet(transitionsHash, eSli)
		}
		lrules = append(lrules, rule)
	}

	if err := d.Set("lifecycle_rule", lrules); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceKsyunKs3BucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)

	d.Partial(true)

	if d.HasChange("acl") && !d.IsNewResource() {
		request := map[string]string{"bucketName": d.Id(), "bucketACL": d.Get("acl").(string)}
		var requestInfo *ks3.Client
		raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return nil, ks3Client.SetBucketACL(d.Id(), ks3.ACLType(d.Get("acl").(string)))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketACL", KsyunKs3GoSdk)
		}
		addDebug("SetBucketACL", raw, requestInfo, request)

	}

	if d.HasChange("cors_rule") {
		if err := resourceKsyunKs3BucketCorsUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("cors_rule")
	}
	if d.HasChange("logging") {
		if err := resourceKsyunKs3BucketLoggingUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("logging")
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceKsyunKs3BucketLifecycleRuleUpdate(client, d); err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)
	return resourceKsyunKs3BucketRead(d, meta)
}

func resourceKsyunKs3BucketCorsUpdate(client *KsyunClient, d *schema.ResourceData) error {
	cors := d.Get("cors_rule").([]interface{})
	var requestInfo *ks3.Client
	if cors == nil || len(cors) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
				requestInfo = ks3Client
				return nil, ks3Client.DeleteBucketCORS(d.Id())
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			addDebug("DeleteBucketCORS", raw, requestInfo, map[string]string{"bucketName": d.Id()})
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketCORS", KsyunKs3GoSdk)
		}
		return nil
	}
	// Put CORS
	rules := make([]ks3.CORSRule, 0, len(cors))
	for _, c := range cors {
		corsMap := c.(map[string]interface{})
		rule := ks3.CORSRule{}
		for k, v := range corsMap {
			log.Printf("[DEBUG] KS3 bucket: %s, put CORS: %#v, %#v", d.Id(), k, v)
			if k == "max_age_seconds" {
				rule.MaxAgeSeconds = v.(int)
			} else {
				rMap := make([]string, len(v.([]interface{})))
				for i, vv := range v.([]interface{}) {
					rMap[i] = vv.(string)
				}
				switch k {
				case "allowed_headers":
					rule.AllowedHeader = rMap
				case "allowed_methods":
					rule.AllowedMethod = rMap
				case "allowed_origins":
					rule.AllowedOrigin = rMap
				case "expose_headers":
					rule.ExposeHeader = rMap
				}
			}
		}
		rules = append(rules, rule)
	}

	log.Printf("[DEBUG] Ks3 bucket: %s, put CORS: %#v", d.Id(), cors)
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		return nil, ks3Client.SetBucketCORS(d.Id(), rules)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketCORS", KsyunKs3GoSdk)
	}
	addDebug("SetBucketCORS", raw, requestInfo, map[string]interface{}{
		"bucketName": d.Id(),
		"corsRules":  rules,
	})
	return nil
}
func resourceKsyunKs3BucketLoggingUpdate(client *KsyunClient, d *schema.ResourceData) error {
	logging := d.Get("logging").([]interface{})
	var requestInfo *ks3.Client
	if logging == nil || len(logging) == 0 {
		raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return nil, ks3Client.DeleteBucketLogging(d.Id())
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketLogging", KsyunKs3GoSdk)
		}
		addDebug("DeleteBucketLogging", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	}

	c := logging[0].(map[string]interface{})
	var target_bucket, target_prefix string
	if v, ok := c["target_bucket"]; ok {
		target_bucket = v.(string)
	}
	if v, ok := c["target_prefix"]; ok {
		target_prefix = v.(string)
	}
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		return nil, ks3Client.SetBucketLogging(d.Id(), target_bucket, target_prefix, target_bucket != "" || target_prefix != "")
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketLogging", KsyunKs3GoSdk)
	}
	addDebug("SetBucketLogging", raw, requestInfo, map[string]interface{}{
		"bucketName":   d.Id(),
		"targetBucket": target_bucket,
		"targetPrefix": target_prefix,
		"isEnable":     target_bucket != "",
	})
	return nil
}

func resourceKsyunKs3BucketLifecycleRuleUpdate(client *KsyunClient, d *schema.ResourceData) error {
	bucket := d.Id()
	lifecycleRules := d.Get("lifecycle_rule").([]interface{})
	var requestInfo *ks3.Client
	if lifecycleRules == nil || len(lifecycleRules) == 0 {
		raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return nil, ks3Client.DeleteBucketLifecycle(bucket)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketLifecycle", KsyunKs3GoSdk)

		}
		addDebug("DeleteBucketLifecycle", raw, requestInfo, map[string]interface{}{
			"bucketName": bucket,
		})
		return nil
	}

	rules := make([]ks3.LifecycleRule, 0, len(lifecycleRules))

	for i, lifecycleRule := range lifecycleRules {
		r := lifecycleRule.(map[string]interface{})
		rule := ks3.LifecycleRule{}
		// ID--有值
		if val, ok := r["id"].(string); ok && val != "" {
			rule.ID = val
		}
		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rule.Status = string(ExpirationStatusEnabled)
		} else {
			rule.Status = string(ExpirationStatusDisabled)
		}
		if filterSet, ok := r["filter"].(*schema.Set); ok {
			if filterSet.Len() > 0 {
				if filter, ok := filterSet.List()[0].(map[string]interface{}); ok {
					filterModel := &ks3.LifecycleFilter{
						And: ks3.LifecycleAnd{},
					}
					filterModel.And.Prefix = filter["prefix"].(string)
					tagList := filter["tag"].([]interface{})
					for _, tag := range tagList {
						tagMap := tag.(map[string]interface{})
						key := tagMap["key"].(string)
						value := tagMap["value"].(string)
						filterModel.And.Tag = append(filterModel.And.Tag, ks3.Tag{
							Key:   key,
							Value: value,
						})
					}
					rule.Filter = filterModel
				}
			}
		} else {
			if val, ok := r["prefix"].(string); ok && val != "" {
				rule.Prefix = val
			}
		}
		// Expiration
		// Expiration
		expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expiration) > 0 {
			e := expiration[0].(map[string]interface{})
			i := ks3.LifecycleExpiration{}
			valDate, _ := e["date"].(string)
			valDays, _ := e["days"].(int)
			cnt := 0
			if valDate != "" {
				i.Date = fmt.Sprintf("%sT00:00:00+08:00", valDate)
				cnt++
			}
			if valDays > 0 {
				i.Days = valDays
				cnt++
			}
			if cnt != 1 {
				return WrapError(Error("One and only one of 'date', 'date' and 'days' can be specified in one expiration configuration."))
			}

			rule.Expiration = &i
		}
		log.Printf("[DEBUG] Ks3 bucket:  %s,rule.Expiration: %#v", d.Id(), rule.Expiration)
		// Transitions
		transitionsRaw := r["transitions"]
		if transitionsRaw != nil {
			transitions := transitionsRaw.(*schema.Set).List()
			if len(transitions) > 0 {
				for _, transition := range transitions {
					transitionTmp := ks3.LifecycleTransition{}
					valStorageClass := transition.(map[string]interface{})["storage_class"].(string)
					date, ok := transition.(map[string]interface{})["date"].(string)
					if ok && date != "" {
						transitionTmp.Date = fmt.Sprintf("%sT00:00:00+08:00", date)
					}
					cnt := 0
					daysInterface, ok := transition.(map[string]interface{})["days"].(string)
					if ok && daysInterface != "" {
						valDays, err := strconv.Atoi(daysInterface)
						if err == nil && valDays > 0 {
							transitionTmp.Days = valDays
							cnt++
						}
					}
					if valStorageClass != "" {
						transitionTmp.StorageClass = ks3.StorageClassType(valStorageClass)
					}
					rule.Transitions = append(rule.Transitions, transitionTmp)
				}
			}
		}
		rules = append(rules, rule)
	}
	log.Printf("[DEBUG] Ks3 bucket: %s, put Lifecycle: %#v", d.Id(), rules)
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		return nil, ks3Client.SetBucketLifecycle(bucket, rules)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketLifecycle", KsyunKs3GoSdk)
	}
	addDebug("SetBucketLifecycle", raw, requestInfo, map[string]interface{}{
		"bucketName": bucket,
		"rules":      rules,
	})
	return nil
}

func resourceKsyunKs3BucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	var requestInfo *ks3.Client
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		return ks3Client.IsBucketExist(d.Id())
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBucketExist", KsyunKs3GoSdk)
	}
	addDebug("IsBucketExist", raw, requestInfo, map[string]string{"bucketName": d.Id()})

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			return nil, ks3Client.DeleteBucket(d.Id())
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"BucketNotEmpty"}) {
				raw, er := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
					bucket, _ := ks3Client.Bucket(d.Get("bucket").(string))
					marker := ""
					retryTime := 3
					for retryTime > 0 {
						lsRes, err := bucket.ListObjects(ks3.Marker(marker))
						if err != nil {
							retryTime--
							time.Sleep(time.Duration(1) * time.Second)
							continue
						}
						for _, object := range lsRes.Objects {
							bucket.DeleteObject(object.Key)
						}
						if lsRes.IsTruncated {
							marker = lsRes.NextMarker
						} else {
							return true, nil
						}
					}
					return false, nil
				})
				if er != nil {
					return resource.NonRetryableError(er)
				}
				addDebug("DeleteObjects", raw, requestInfo, map[string]string{"bucketName": d.Id()})
				return resource.RetryableError(err)
			}
		}
		return resource.NonRetryableError(err)
		addDebug("DeleteBucket", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucket", KsyunKs3GoSdk)
	}
	return nil
}

func transitionsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return hashcode.String(buf.String())
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}
