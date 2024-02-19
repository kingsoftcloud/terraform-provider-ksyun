/*
Provides an KS3 Bucket resource.

# Example Usage

```hcl

provider "ksyun" {
  #指定KS3服务的访问域名
  endpoint = "ks3-cn-beijing.ksyuncs.com"
}

resource "ksyun_ks3_bucket" "bucket-create" {
  #指定要创建的虚拟存储桶的名称
  bucket = "bucket-20240206-104450"
}

resource "ksyun_ks3_bucket" "bucket-create2" {
  #指定要创建的虚拟存储桶的名称
  bucket = "bucket-20240205-153413"
  #桶权限
  acl = "public-read"
  #桶存储类型
  storage_class = "IA"
  #空间策略
  policy = <<-EOT
  {
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "ks3:ListBucket",
          "ks3:ListBucketMultipartUploads",
          "ks3:GetObject",
          "ks3:GetObjectAcl",
          "ks3:ListMultipartUploadParts"
        ],
        "Principal": {
          "KSC": [
            "*"
          ]
        },
        "Resource": [
          "krn:ksc:ks3:::bucket-20240205-153413/*",
          "krn:ksc:ks3:::bucket-20240205-153413"
        ]
      }
    ]
  }
  EOT
  #日志
  logging {
    target_bucket = "bucket-20240206-104450"
    target_prefix = "log/"
  }
  #跨域资源共享
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET","PUT"]
    allowed_origins = ["*"]
    expose_headers = ["header1","header2"]
    max_age_seconds = 3600
  }
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["POST","DELETE"]
    allowed_origins = ["https://www.abc.com","https://www.def.com"]
    expose_headers = ["header1","header2"]
    max_age_seconds = 3600
  }
  #生命周期
  lifecycle_rule {
    id      = "id1"
    enabled = true
    filter {
      prefix = "logs"
    }
    expiration {
      days = 130
    }
    transitions {
      days = 10
      storage_class = "STANDARD_IA"
    }
    transitions {
      days = 40
      storage_class = "ARCHIVE"
    }
    abort_incomplete_multipart_upload {
      days_after_initiation = 50
    }
  }
  lifecycle_rule {
    id      = "id2"
    enabled = true
    filter {
      prefix = "test"
    }
    expiration {
      date = "2024-07-10"
    }
    transitions {
      date = "2024-05-10"
      storage_class = "STANDARD_IA"
    }
    transitions {
      date = "2024-06-10"
      storage_class = "ARCHIVE"
    }
    abort_incomplete_multipart_upload {
      date = "2024-07-10"
    }
  }
  lifecycle_rule {
    id      = "id3"
    enabled = true
    filter {
      and {
        prefix = "docs"
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
    expiration {
      date = "2024-07-10"
    }
    transitions {
      date = "2024-05-10"
      storage_class = "STANDARD_IA"
    }
    transitions {
      date = "2024-06-10"
      storage_class = "ARCHIVE"
    }
  }
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
	"github.com/ks3sdklib/ksyun-ks3-go-sdk/ks3"
	"log"
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

			"cors_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    10,
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
			},

			"logging": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Call this interface to set the bucket logging configuration. If the configuration already exists, KS3 will replace it.\nTo use this interface, you need to have permission to perform the ks3: PutBucketLogging operation. The space owner has this permission by default and can grant corresponding permissions to others. If you want to turn off this setting, just leave it blank in the configuration.",
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
			},

			"lifecycle_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				MaxItems:    10,
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
							Description: "Prefix identifying one or more objects to which the rule applies.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If 'Enabled', the rule is currently being applied. If 'Disabled', the rule is not currently being applied.",
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
										Description: "Prefix identifying one or more objects to which the rule applies.",
									},
									"and": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "A subset of the object filter that is required when specifying tag.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prefix": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Prefix identifying one or more objects to which the rule applies.",
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
								},
							},
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
										Type:        schema.TypeInt,
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
						"abort_incomplete_multipart_upload": {
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Description: "Specifies expiration attributes for incomplete multipart uploads.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies a date and KS3 will execute lifecycle rules on data that was last modified earlier than that date. If Date is a future date, the rule will not take effect until that date.",
									},
									"days_after_initiation": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Description:  "Specifies the number of days after the start of the multipart upload that the upload must be completed, otherwise it will be deleted. The value must be a non-zero positive integer.",
									},
								},
							},
						},
					},
				},
			},

			"policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bucket Policy is an authorization policy for Bucket introduced by KS3. You can authorize other users to access the KS3 resources you specify through the space policy. If you want to turn off this setting, just leave it blank in the configuration.",
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
	d.Set("region", object.BucketInfo.Region)
	d.Set("owner", object.BucketInfo.Owner.ID)
	d.Set("storage_class", object.BucketInfo.StorageClass)

	request := map[string]string{"bucketName": d.Id()}
	var requestInfo *ks3.Client

	// Read the CORS configuration
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
	lifecycleRules := make([]map[string]interface{}, 0)
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

		// Filter
		if lifecycleRule.Filter != nil {
			filter := make(map[string]interface{})
			if lifecycleRule.Filter.Prefix != "" {
				filter["prefix"] = lifecycleRule.Filter.Prefix
			}
			// And
			if lifecycleRule.Filter.And != nil {
				and := make(map[string]interface{})
				if lifecycleRule.Filter.And.Prefix != "" {
					and["prefix"] = lifecycleRule.Filter.And.Prefix
				}
				if len(lifecycleRule.Filter.And.Tag) != 0 {
					var tagList []interface{}
					for _, tagVal := range lifecycleRule.Filter.And.Tag {
						tag := make(map[string]interface{})
						tag["key"] = tagVal.Key
						tag["value"] = tagVal.Value
						tagList = append(tagList, tag)
					}
					and["tag"] = tagList
				}
				filter["and"] = []interface{}{and}
			}
			rule["filter"] = []interface{}{filter}
		}

		// Expiration
		if lifecycleRule.Expiration != nil {
			log.Printf("[DEBUG] lifecycleRule.Expiration start: %#v", lifecycleRule.Expiration)
			expiration := make(map[string]interface{})
			if lifecycleRule.Expiration.Date != "" {
				lifecycleRule.Expiration.Date = strings.ReplaceAll(lifecycleRule.Expiration.Date, ".000", "")
				t, err := time.Parse(Iso8601DateFormat, lifecycleRule.Expiration.Date)
				if err != nil {
					return WrapError(err)
				}
				expiration["date"] = t.Format("2006-01-02")
			}
			if lifecycleRule.Expiration.Days != 0 {
				expiration["days"] = lifecycleRule.Expiration.Days
			}
			rule["expiration"] = schema.NewSet(expirationHash, []interface{}{expiration})
		}

		// Transitions
		if len(lifecycleRule.Transitions) != 0 {
			var transitionList []interface{}
			for _, transitionItem := range lifecycleRule.Transitions {
				transition := make(map[string]interface{})
				if transitionItem.Date != "" {
					transitionItem.Date = strings.ReplaceAll(transitionItem.Date, ".000", "")
					t, err := time.Parse(Iso8601DateFormat, transitionItem.Date)
					if err != nil {
						return WrapError(err)
					}
					transition["date"] = t.Format("2006-01-02")
				}
				if transitionItem.Days != 0 {
					transition["days"] = transitionItem.Days
				}
				if transitionItem.StorageClass != "" {
					transition["storage_class"] = string(transitionItem.StorageClass)
				}
				transitionList = append(transitionList, transition)
			}
			rule["transitions"] = schema.NewSet(transitionsHash, transitionList)
		}

		// AbortIncompleteMultipartUpload
		if lifecycleRule.AbortIncompleteMultipartUpload != nil {
			log.Printf("[DEBUG] lifecycleRule.AbortIncompleteMultipartUpload start: %#v", lifecycleRule.AbortIncompleteMultipartUpload)
			abortMul := make(map[string]interface{})
			if lifecycleRule.AbortIncompleteMultipartUpload.Date != "" {
				lifecycleRule.AbortIncompleteMultipartUpload.Date = strings.ReplaceAll(lifecycleRule.AbortIncompleteMultipartUpload.Date, ".000", "")
				t, err := time.Parse(Iso8601DateFormat, lifecycleRule.AbortIncompleteMultipartUpload.Date)
				if err != nil {
					return WrapError(err)
				}
				abortMul["date"] = t.Format("2006-01-02")
			}
			if lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation != 0 {
				abortMul["days_after_initiation"] = lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation
			}
			rule["abort_incomplete_multipart_upload"] = schema.NewSet(abortMulHash, []interface{}{abortMul})
		}

		lifecycleRules = append(lifecycleRules, rule)
	}

	if err := d.Set("lifecycle_rule", lifecycleRules); err != nil {
		return WrapError(err)
	}

	// Read the policy configuration
	raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		params := map[string]interface{}{}
		params["policy"] = nil
		return ks3Client.GetBucketPolicy(d.Id())
	})

	if err != nil && !ks3NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketPolicy", KsyunKs3GoSdk)
	}
	addDebug("GetBucketPolicy", raw, requestInfo, request)
	policy := ""
	if err == nil {
		policy = raw.(string)
	}

	if err := d.Set("policy", policy); err != nil {
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

	if d.HasChange("policy") {
		if err := resourceKsyunKs3BucketPolicyUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("policy")
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
		// ID
		if val, ok := r["id"].(string); ok && val != "" {
			rule.ID = val
		}
		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rule.Status = string(ExpirationStatusEnabled)
		} else {
			rule.Status = string(ExpirationStatusDisabled)
		}
		// Prefix
		if val, ok := r["prefix"].(string); ok && val != "" {
			rule.Prefix = val
		}
		// Filter
		filterList := d.Get(fmt.Sprintf("lifecycle_rule.%d.filter", i)).(*schema.Set).List()
		if len(filterList) > 0 {
			filter := filterList[0].(map[string]interface{})
			filterVal := &ks3.LifecycleFilter{}
			filterVal.Prefix = filter["prefix"].(string)
			andList := filter["and"].([]interface{})
			if len(andList) > 0 {
				and := andList[0].(map[string]interface{})
				filterVal.And = &ks3.LifecycleAnd{}
				filterVal.And.Prefix = and["prefix"].(string)
				tagList := and["tag"].([]interface{})
				for _, tagItem := range tagList {
					tag := tagItem.(map[string]interface{})
					tagVal := ks3.Tag{}
					tagVal.Key = tag["key"].(string)
					tagVal.Value = tag["value"].(string)
					filterVal.And.Tag = append(filterVal.And.Tag, tagVal)
				}
			}
			rule.Filter = filterVal
		}

		// Expiration
		expirationList := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expirationList) > 0 {
			expiration := expirationList[0].(map[string]interface{})
			expirationVal := &ks3.LifecycleExpiration{}
			date, _ := expiration["date"].(string)
			days, _ := expiration["days"].(int)
			if date != "" {
				expirationVal.Date = fmt.Sprintf("%sT00:00:00.000+08:00", date)
			}
			if days > 0 {
				expirationVal.Days = days
			}
			rule.Expiration = expirationVal
		}

		// Transitions
		transitionList := d.Get(fmt.Sprintf("lifecycle_rule.%d.transitions", i)).(*schema.Set).List()
		if len(transitionList) > 0 {
			for _, transitionItem := range transitionList {
				transition := transitionItem.(map[string]interface{})
				transitionVal := ks3.LifecycleTransition{}
				date, _ := transition["date"].(string)
				days, _ := transition["days"].(int)
				storageClass, _ := transition["storage_class"].(string)
				if date != "" {
					transitionVal.Date = fmt.Sprintf("%sT00:00:00.000+08:00", date)
				}
				if days > 0 {
					transitionVal.Days = days
				}
				if storageClass != "" {
					transitionVal.StorageClass = ks3.StorageClassType(storageClass)
				}
				rule.Transitions = append(rule.Transitions, transitionVal)
			}
		}

		// AbortIncompleteMultipartUpload
		abortIncompleteMultipartUploadList := d.Get(fmt.Sprintf("lifecycle_rule.%d.abort_incomplete_multipart_upload", i)).(*schema.Set).List()
		if len(abortIncompleteMultipartUploadList) > 0 {
			abortMul := abortIncompleteMultipartUploadList[0].(map[string]interface{})
			abortMulVal := &ks3.LifecycleAbortIncompleteMultipartUpload{}
			date, _ := abortMul["date"].(string)
			daysAfterInitiation, _ := abortMul["days_after_initiation"].(int)
			if date != "" {
				abortMulVal.Date = fmt.Sprintf("%sT00:00:00.000+08:00", date)
			}
			if daysAfterInitiation > 0 {
				abortMulVal.DaysAfterInitiation = daysAfterInitiation
			}
			rule.AbortIncompleteMultipartUpload = abortMulVal
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

func resourceKsyunKs3BucketPolicyUpdate(client *KsyunClient, d *schema.ResourceData) error {
	bucket := d.Id()
	policy := d.Get("policy").(string)
	var requestInfo *ks3.Client
	if len(policy) == 0 {
		raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return nil, ks3Client.DeleteBucketPolicy(bucket)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketPolicy", KsyunKs3GoSdk)
		}
		addDebug("DeleteBucketPolicy", raw, requestInfo, map[string]string{"bucketName": bucket})
		return nil
	}
	raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		buffer := new(bytes.Buffer)
		buffer.Write([]byte(policy))
		return nil, ks3Client.SetBucketPolicy(bucket, policy)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketPolicy", KsyunKs3GoSdk)
	}
	addDebug("SetBucketPolicy", raw, requestInfo, map[string]string{"bucketName": bucket})
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
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
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

func abortMulHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days_after_initiation"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}
