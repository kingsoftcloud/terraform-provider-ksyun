/*
This data source provides a list of Ks3 resources according to bucket name they belong to.

# Example Usage

```hcl

	provider "ksyun" {
	  #指定KS3服务的访问域名
	  endpoint = "ks3-cn-beijing.ksyuncs.com"
	}

	data "ksyun_ks3_buckets" "default" {
	  #匹配全部包涵该字符串的bucket
	  name_regex  = "bucket-202402"
	  #输出文件路径
	  output_file = "bucket_info.txt"
	}

```
*/

package ksyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/ks3sdklib/ksyun-ks3-go-sdk/ks3"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func dataSourceKsyunKs3Buckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKs3BucketsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The string used to match buckets.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path of the output file.",
			},
			"buckets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Matched bucket list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the bucket.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the bucket is located.",
						},
						"acl": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The canned ACL to apply.",
						},
						"storage_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The class of storage used to store the object.",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the bucket.",
						},

						"cors_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cross domain resource sharing (CORS) configuration rules of the bucket. Each rule can contain up to 100 headers, 100 methods, and 100 origins. For more information, see Cross-domain resource sharing (CORS) in the KS3 Developer Guide.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The HTTP header that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, Authorization. Additionally, you can use \"*\" to represent all headers.",
									},
									"allowed_methods": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The HTTP method that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, GET, PUT, POST, DELETE, HEAD, OPTIONS.",
									},
									"allowed_origins": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The source address that users are allowed to access through cross domain resource sharing, which can contain up to one \"*\" wildcard character. Each CORSRule must define at least one source address and one method. For example, http://*. example. com. Additionally, you can use \"*\" to represent all sources.",
									},
									"expose_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Identify response headers that allow customers to access from applications, such as JavaScript XMLHttpRequest data element.",
									},
									"max_age_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The maximum time in seconds that the browser can cache the preflight response.",
									},
								},
							},
						},

						"logging": {
							Type:        schema.TypeList,
							Computed:    true,
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
							Computed:    true,
							Description: "Bucket Policy is an authorization policy for Bucket introduced by KS3. You can authorize other users to access the KS3 resources you specify through the space policy. If you want to turn off this setting, just leave it blank in the configuration.",
						},

						"tags": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "the tags of the resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKs3BucketsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	var requestInfo *ks3.Client
	var allBuckets []ks3.BucketProperties
	nextMarker := ""
	for {
		var options []ks3.Option
		if nextMarker != "" {
			options = append(options, ks3.Marker(nextMarker))
		}

		raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return ks3Client.ListBuckets(options...)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "ksyun_ks3_bucket", "CreateBucket", KsyunKs3GoSdk)
		}
		if debugOn() {
			addDebug("ListBuckets", raw, requestInfo, map[string]interface{}{"options": options})
		}
		response, _ := raw.(ks3.ListBucketsResult)

		if response.Buckets == nil || len(response.Buckets) < 1 {
			break
		}

		allBuckets = append(allBuckets, response.Buckets...)

		nextMarker = response.NextMarker
		if nextMarker == "" {
			break
		}
	}

	var filteredBucketsTemp []ks3.BucketProperties
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var ks3BucketNameRegex *regexp.Regexp
		if nameRegex != "" {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
			ks3BucketNameRegex = r
		}
		for _, bucket := range allBuckets {
			if ks3BucketNameRegex != nil && !ks3BucketNameRegex.MatchString(bucket.Name) {
				continue
			}
			filteredBucketsTemp = append(filteredBucketsTemp, bucket)
		}
	} else {
		filteredBucketsTemp = allBuckets
	}
	return bucketsDescriptionAttributes(d, filteredBucketsTemp, meta)
}

func bucketsDescriptionAttributes(d *schema.ResourceData, buckets []ks3.BucketProperties, meta interface{}) error {
	client := meta.(*KsyunClient)

	var ids []string
	var s []map[string]interface{}
	var requestInfo *ks3.Client
	for _, bucket := range buckets {
		mapping := map[string]interface{}{
			"name":          bucket.Name,
			"creation_date": bucket.CreationDate.Format("2006-01-02"),
		}

		// bucket info
		raw, err := client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return GetBucketInfo(ks3Client, bucket.Name)
		})
		if err == nil {
			if debugOn() {
				addDebug("GetBucketInfo", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			response, _ := raw.(ks3.GetBucketInfoResult)
			mapping["acl"] = response.BucketInfo.ACL
			mapping["region"] = response.BucketInfo.Region
			mapping["storage_class"] = response.BucketInfo.StorageClass
		} else {
			log.Printf("[WARN] Unable to get additional information for the bucket %s: %v", bucket.Name, err)
		}

		// CORS
		var ruleMappings []map[string]interface{}
		raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return ks3Client.GetBucketCORS(bucket.Name)
		})
		if err == nil {
			if debugOn() {
				addDebug("GetBucketCORS", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			cors, _ := raw.(ks3.GetBucketCORSResult)
			if cors.CORSRules != nil {
				for _, rule := range cors.CORSRules {
					ruleMapping := make(map[string]interface{})
					ruleMapping["allowed_headers"] = rule.AllowedHeader
					ruleMapping["allowed_methods"] = rule.AllowedMethod
					ruleMapping["allowed_origins"] = rule.AllowedOrigin
					ruleMapping["expose_headers"] = rule.ExposeHeader
					ruleMapping["max_age_seconds"] = rule.MaxAgeSeconds
					ruleMappings = append(ruleMappings, ruleMapping)
				}
			}
		} else if !IsExpectedErrors(err, []string{"NoSuchCORSConfiguration"}) {
			log.Printf("[WARN] Unable to get CORS information for the bucket %s: %v", bucket.Name, err)
		}
		mapping["cors_rule"] = ruleMappings

		// logging
		var loggingMappings []map[string]interface{}
		raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			return ks3Client.GetBucketLogging(bucket.Name)
		})
		if err == nil {
			addDebug("GetBucketLogging", raw)
			logging, _ := raw.(ks3.GetBucketLoggingResult)
			if logging.LoggingEnabled.TargetBucket != "" || logging.LoggingEnabled.TargetPrefix != "" {
				loggingMapping := map[string]interface{}{
					"target_bucket": logging.LoggingEnabled.TargetBucket,
					"target_prefix": logging.LoggingEnabled.TargetPrefix,
				}
				loggingMappings = append(loggingMappings, loggingMapping)
			}
		} else {
			log.Printf("[WARN] Unable to get logging information for the bucket %s: %v", bucket.Name, err)
		}
		mapping["logging"] = loggingMappings

		// lifecycle
		var lifecycleRuleMappings []map[string]interface{}
		raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return ks3Client.GetBucketLifecycle(bucket.Name)
		})
		if err == nil {
			if debugOn() {
				addDebug("GetBucketLifecycle", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			lifecycle, _ := raw.(ks3.GetBucketLifecycleResult)
			if len(lifecycle.Rules) > 0 {
				for _, lifecycleRule := range lifecycle.Rules {
					ruleMapping := make(map[string]interface{})
					ruleMapping["id"] = lifecycleRule.ID
					if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
						ruleMapping["enabled"] = true
					} else {
						ruleMapping["enabled"] = false
					}
					if lifecycleRule.Prefix != "" {
						ruleMapping["prefix"] = lifecycleRule.Prefix
					}
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
								var tags []interface{}
								for _, tag := range lifecycleRule.Filter.And.Tag {
									e := make(map[string]interface{})
									e["key"] = tag.Key
									e["value"] = tag.Value
									tags = append(tags, e)
								}
								and["tag"] = tags
							}
							filter["and"] = []interface{}{and}
						}
						ruleMapping["filter"] = []interface{}{filter}
					}
					// Expiration
					if lifecycleRule.Expiration != nil {
						expirationMapping := make(map[string]interface{})
						if lifecycleRule.Expiration.Date != "" {
							lifecycleRule.Expiration.Date = strings.ReplaceAll(lifecycleRule.Expiration.Date, ".000", "")
							t, err := time.Parse(Iso8601DateFormat, lifecycleRule.Expiration.Date)
							if err != nil {
								return WrapError(err)
							}
							expirationMapping["date"] = t.Format("2006-01-02")
						}
						if lifecycleRule.Expiration.Days != 0 {
							expirationMapping["days"] = lifecycleRule.Expiration.Days
						}
						ruleMapping["expiration"] = []interface{}{expirationMapping}
						lifecycleRuleMappings = append(lifecycleRuleMappings, ruleMapping)
					}
					// Transitions
					if len(lifecycleRule.Transitions) > 0 {
						var transitionList []interface{}
						for _, transition := range lifecycleRule.Transitions {
							transitionMapping := make(map[string]interface{})
							if transition.Date != "" {
								transition.Date = strings.ReplaceAll(transition.Date, ".000", "")
								t, err := time.Parse(Iso8601DateFormat, transition.Date)
								if err != nil {
									return WrapError(err)
								}
								transitionMapping["date"] = t.Format("2006-01-02")
							}
							if transition.Days != 0 {
								transitionMapping["days"] = transition.Days
							}
							if transition.StorageClass != "" {
								transitionMapping["storage_class"] = transition.StorageClass
							}
							transitionList = append(transitionList, transitionMapping)
						}
						ruleMapping["transitions"] = transitionList
					}
					// AbortIncompleteMultipartUpload
					if lifecycleRule.AbortIncompleteMultipartUpload != nil {
						abortMulMapping := make(map[string]interface{})
						if lifecycleRule.AbortIncompleteMultipartUpload.Date != "" {
							lifecycleRule.AbortIncompleteMultipartUpload.Date = strings.ReplaceAll(lifecycleRule.AbortIncompleteMultipartUpload.Date, ".000", "")
							t, err := time.Parse(Iso8601DateFormat, lifecycleRule.AbortIncompleteMultipartUpload.Date)
							if err != nil {
								return WrapError(err)
							}
							abortMulMapping["date"] = t.Format("2006-01-02")
						}
						if lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation != 0 {
							abortMulMapping["days_after_initiation"] = lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation
						}
						ruleMapping["abort_incomplete_multipart_upload"] = []interface{}{abortMulMapping}
						lifecycleRuleMappings = append(lifecycleRuleMappings, ruleMapping)
					}
				}
			}
		} else {
			log.Printf("[WARN] Unable to get lifecycle information for the bucket %s: %v", bucket.Name, err)
		}
		mapping["lifecycle_rule"] = lifecycleRuleMappings

		// policy
		var policy string
		raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return ks3Client.GetBucketPolicy(bucket.Name)
		})

		if err == nil {
			if debugOn() {
				addDebug("GetBucketPolicy", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			policy = raw.(string)
		} else {
			log.Printf("[WARN] Unable to get policy information for the bucket %s: %v", bucket.Name, err)
		}
		mapping["policy"] = policy

		// tags
		tagsMap := make(map[string]string)
		raw, err = client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
			requestInfo = ks3Client
			return ks3Client.GetBucketTagging(bucket.Name)
		})
		if err == nil {
			if debugOn() {
				addDebug("GetBucketTagging", raw, requestInfo, map[string]string{"bucketName": bucket.Name})
			}
			tagging, _ := raw.(ks3.GetBucketTaggingResult)
			for _, t := range tagging.Tags {
				tagsMap[t.Key] = t.Value
			}
		} else {
			log.Printf("[WARN] Unable to get tagging information for the bucket %s: %v", bucket.Name, err)
		}
		mapping["tags"] = tagsMap

		ids = append(ids, bucket.Name)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("buckets", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeBucketInfoToFile(output.(string), s)
	}
	return nil
}

func dataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func GetBucketInfo(client *ks3.Client, bucket string) (ks3.GetBucketInfoResult, error) {
	acl, err := client.GetBucketACL(bucket)
	if err != nil {
		return ks3.GetBucketInfoResult{}, err
	}
	resp, err := client.ListBuckets()
	if err != nil {
		return ks3.GetBucketInfoResult{}, err
	}
	for _, bucketInfo := range resp.Buckets {
		if bucketInfo.Name == bucket {
			return ks3.GetBucketInfoResult{
				BucketInfo: ks3.BucketInfo{
					XMLName:      bucketInfo.XMLName,
					Name:         bucket,
					Region:       bucketInfo.Region,
					ACL:          string(acl.GetCannedACL()),
					StorageClass: bucketInfo.Type,
				}}, nil
		}
	}
	return ks3.GetBucketInfoResult{}, nil
}

func writeBucketInfoToFile(filePath string, data interface{}) error {
	absPath, err := getAbsPath(filePath)
	if err != nil {
		return err
	}
	os.Remove(absPath)
	var bs []byte
	switch data := data.(type) {
	case string:
		bs = []byte(data)
	default:
		bs, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v and got an error: %#v", data, err)
		}
	}
	return os.WriteFile(absPath, bs, 0422)
}
