---
subcategory: "KS3"
layout: "ksyun"
page_title: "ksyun: ksyun_ks3_bucket"
sidebar_current: "docs-ksyun-resource-ks3_bucket"
description: |-
  Provides an KS3 Bucket resource.
---

# ksyun_ks3_bucket

Provides an KS3 Bucket resource.

#

## Example Usage

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
    allowed_methods = ["GET", "PUT"]
    allowed_origins = ["*"]
    expose_headers  = ["header1", "header2"]
    max_age_seconds = 3600
  }
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["POST", "DELETE"]
    allowed_origins = ["https://www.abc.com", "https://www.def.com"]
    expose_headers  = ["header1", "header2"]
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
      days          = 10
      storage_class = "STANDARD_IA"
    }
    transitions {
      days          = 40
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
      date          = "2024-05-10"
      storage_class = "STANDARD_IA"
    }
    transitions {
      date          = "2024-06-10"
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
      date          = "2024-05-10"
      storage_class = "STANDARD_IA"
    }
    transitions {
      date          = "2024-06-10"
      storage_class = "ARCHIVE"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) The canned ACL to apply. Defaults to private. Valid values are private, public-read, and public-read-write.
* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random, unique name.
* `cors_rule` - (Optional) The cross domain resource sharing (CORS) configuration rules of the bucket. Each rule can contain up to 100 headers, 100 methods, and 100 origins. For more information, see Cross-domain resource sharing (CORS) in the KS3 Developer Guide.
* `lifecycle_rule` - (Optional) Call this interface to set the bucket lifecycle configuration. If the configuration already exists, KS3 will replace it.
To use this interface, you need to have permission to perform the ks3: PutBucketLifecycle operation. The space owner has this permission by default and can grant corresponding permissions to others.
* `logging` - (Optional) Call this interface to set the bucket logging configuration. If the configuration already exists, KS3 will replace it.
To use this interface, you need to have permission to perform the ks3: PutBucketLogging operation. The space owner has this permission by default and can grant corresponding permissions to others. If you want to turn off this setting, just leave it blank in the configuration.
* `policy` - (Optional) Bucket Policy is an authorization policy for Bucket introduced by KS3. You can authorize other users to access the KS3 resources you specify through the space policy. If you want to turn off this setting, just leave it blank in the configuration.
* `storage_class` - (Optional) The class of storage used to store the object.

The `abort_incomplete_multipart_upload` object supports the following:

* `date` - (Optional) Specifies a date and KS3 will execute lifecycle rules on data that was last modified earlier than that date. If Date is a future date, the rule will not take effect until that date.
* `days_after_initiation` - (Optional) Specifies the number of days after the start of the multipart upload that the upload must be completed, otherwise it will be deleted. The value must be a non-zero positive integer.

The `and` object supports the following:

* `prefix` - (Optional) Prefix identifying one or more objects to which the rule applies.
* `tag` - (Optional) Container for the tag of lifecycle rule. If you specify a tag, you cannot specify a prefix for the rule.

The `cors_rule` object supports the following:

* `allowed_methods` - (Required) The HTTP method that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, GET, PUT, POST, DELETE, HEAD, OPTIONS.
* `allowed_origins` - (Required) The source address that users are allowed to access through cross domain resource sharing, which can contain up to one "*" wildcard character. Each CORSRule must define at least one source address and one method. For example, http://*. example. com. Additionally, you can use "*" to represent all sources.
* `allowed_headers` - (Optional) The HTTP header that users are allowed to access through cross domain resource sharing. Each CORSRule must define at least one source address and one method. For example, Authorization. Additionally, you can use "*" to represent all headers.
* `expose_headers` - (Optional) Identify response headers that allow customers to access from applications, such as JavaScript XMLHttpRequest data element.
* `max_age_seconds` - (Optional) The maximum time in seconds that the browser can cache the preflight response.

The `expiration` object supports the following:

* `date` - (Optional) Indicates at what date the object is to be moved or deleted.example:2016-01-02.
* `days` - (Optional) Indicates the lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer.

The `filter` object supports the following:

* `and` - (Optional) A subset of the object filter that is required when specifying tag.
* `prefix` - (Optional) Prefix identifying one or more objects to which the rule applies.

The `lifecycle_rule` object supports the following:

* `enabled` - (Required) If 'Enabled', the rule is currently being applied. If 'Disabled', the rule is not currently being applied.
* `abort_incomplete_multipart_upload` - (Optional) Specifies expiration attributes for incomplete multipart uploads.
* `expiration` - (Optional) Specifies when an object transitions to a specified storage class. If you specify multiple rules in a lifecycle configuration, the rule with the earliest expiration date is applied to the object. If you specify multiple transition rules for the same object, the rule with the earliest date is applied.
* `filter` - (Optional) Container for the filter of lifecycle rule. If you specify a filter, you cannot specify a prefix for the rule.
* `id` - (Optional) Unique identifier for the rule. The value cannot be longer than 255 characters.
* `prefix` - (Optional) Prefix identifying one or more objects to which the rule applies.
* `transitions` - (Optional) Specifies when an object transitions to a specified storage class. If you specify multiple rules in a lifecycle configuration, the rule with the earliest transition date is applied to the object. If you specify multiple transition rules for the same object, the rule with the earliest date is applied.

The `logging` object supports the following:

* `target_bucket` - (Required) The name of the bucket where you want KS3 to store server access logs. You can have your logs delivered to any bucket that you own, including the same bucket that is being logged. You can also configure multiple buckets to deliver their logs to the same target bucket. In this case, you should assign each bucket a unique prefix.
* `target_prefix` - (Optional) A prefix for all log object keys. If you store log files from multiple buckets in a single bucket, you can use a prefix to distinguish which log files came from which bucket.

The `tag` object supports the following:

* `key` - (Required) The key of the tag.
* `value` - (Required) The value of the tag.

The `transitions` object supports the following:

* `date` - (Optional) Indicates at what date the object is to be moved or deleted.example:2016-01-02.
* `days` - (Optional) Indicates the lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer.
* `storage_class` - (Optional) The class of storage used to store the object.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



