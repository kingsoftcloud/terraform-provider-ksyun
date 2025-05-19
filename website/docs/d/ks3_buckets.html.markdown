---
subcategory: "KS3"
layout: "ksyun"
page_title: "ksyun: ksyun_ks3_buckets"
sidebar_current: "docs-ksyun-datasource-ks3_buckets"
description: |-
  This data source provides a list of Ks3 resources according to bucket name they belong to.
---

# ksyun_ks3_buckets

This data source provides a list of Ks3 resources according to bucket name they belong to.

#

## Example Usage

```hcl
provider "ksyun" {
  #指定KS3服务的访问域名
  endpoint = "ks3-cn-beijing.ksyuncs.com"
}

data "ksyun_ks3_buckets" "default" {
  #匹配全部包涵该字符串的bucket
  name_regex = "bucket-202402"
  #输出文件路径
  output_file = "bucket_info.txt"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) The string used to match buckets.
* `output_file` - (Optional) The path of the output file.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `buckets` - The Matched bucket list.
  * `acl` - The canned ACL to apply.
  * `cors_rule` - The cross domain resource sharing (CORS) configuration rules of the bucket. Each rule can contain up to 100 headers, 100 methods, and 100 origins. For more information, see Cross-domain resource sharing (CORS) in the KS3 Developer Guide.
    * `allowed_headers` - The HTTP header that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, Authorization. Additionally, you can use "*" to represent all headers.
    * `allowed_methods` - The HTTP method that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, GET, PUT, POST, DELETE, HEAD, OPTIONS.
    * `allowed_origins` - The source address that users are allowed to access through cross domain resource sharing, which can contain up to one "*" wildcard character. Each CORSRule must define at least one source address and one method. For example, http://*. example. com. Additionally, you can use "*" to represent all sources.
    * `expose_headers` - Identify response headers that allow customers to access from applications, such as JavaScript XMLHttpRequest data element.
    * `max_age_seconds` - The maximum time in seconds that the browser can cache the preflight response.
  * `creation_date` - The creation time of the bucket.
  * `lifecycle_rule` - Call this interface to set the bucket lifecycle configuration. If the configuration already exists, KS3 will replace it.
To use this interface, you need to have permission to perform the ks3: PutBucketLifecycle operation. The space owner has this permission by default and can grant corresponding permissions to others.
    * `id` - Unique identifier for the rule. The value cannot be longer than 255 characters.
  * `logging` - Call this interface to set the bucket logging configuration. If the configuration already exists, KS3 will replace it.
To use this interface, you need to have permission to perform the ks3: PutBucketLogging operation. The space owner has this permission by default and can grant corresponding permissions to others. If you want to turn off this setting, just leave it blank in the configuration.
  * `name` - The name of the bucket.
  * `policy` - Bucket Policy is an authorization policy for Bucket introduced by KS3. You can authorize other users to access the KS3 resources you specify through the space policy. If you want to turn off this setting, just leave it blank in the configuration.
  * `region` - The region where the bucket is located.
  * `storage_class` - The class of storage used to store the object.
  * `tags` - the tags of the resource.


