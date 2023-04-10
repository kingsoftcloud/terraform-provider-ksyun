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
# 创建一个名为 "bucket-new" 的 ks3 存储桶资源
resource "ksyun_ks3_bucket" "bucket-new" {
  provider = ks3.bj-prod
  #指定要创建的虚拟存储桶的名称
  bucket = "YOUR_BUCKET"
  #桶权限
  acl = "public-read"
}

resource "ksyun_ks3_bucket" "bucket-attr" {
  provider = ksyun.bj-prod
  bucket   = "bucket-20230324-1"
  #访问日志的存储路径
  logging {
    target_bucket = "bucket-20230324-2"
    target_prefix = "log/"
  }
  #文件生命周期
  lifecycle_rule {
    id = "id1"
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
    allowed_origins = [var.allow-origins-star]
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
  type    = number
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

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) The canned ACL to apply. Defaults to private. Valid values are private, public-read, and public-read-write.
* `bucket` - (Optional, ForceNew) The name of the bucket. If omitted, Terraform will assign a random, unique name.
* `cors_rule` - (Optional) The cross domain resource sharing (CORS) configuration rules of the bucket. Each rule can contain up to 100 headers, 100 methods, and 100 origins. For more information, see Cross-domain resource sharing (CORS) in the KS3 Developer Guide.
* `lifecycle_rule` - (Optional) Call this interface to set the bucket lifecycle configuration. If the configuration already exists, KS3 will replace it.
To use this interface, you need to have permission to perform the ks3: PutBucketLifecycle operation. The space owner has this permission by default and can grant corresponding permissions to others.
* `logging` - (Optional) Call this interface to set the bucket logging configuration. If the configuration already exists, KS3 will replace it.
To use this interface, you need to have permission to perform the ks3: PutBucketLogging operation. The space owner has this permission by default and can grant corresponding permissions to others.
* `storage_class` - (Optional) The class of storage used to store the object.

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

* `prefix` - (Optional) Prefix identifying one or more objects to which the rule applies. The rule applies only to objects with key names beginning with the specified prefix. If you specify multiple rules in a lifecycle configuration, the rule with the longest prefix is used to determine which rule applies to a given object. If you specify a prefix for the rule, you cannot specify a tag for the rule.
* `tag` - (Optional) Container for the tag of lifecycle rule. If you specify a tag, you cannot specify a prefix for the rule.

The `lifecycle_rule` object supports the following:

* `enabled` - (Required) If 'Enabled', the rule is currently being applied. If 'Disabled', the rule is not currently being applied.
* `expiration` - (Optional) Specifies when an object transitions to a specified storage class. If you specify multiple rules in a lifecycle configuration, the rule with the earliest expiration date is applied to the object. If you specify multiple transition rules for the same object, the rule with the earliest date is applied.
* `filter` - (Optional) Container for the filter of lifecycle rule. If you specify a filter, you cannot specify a prefix for the rule.
* `id` - (Optional) Unique identifier for the rule. The value cannot be longer than 255 characters.
* `prefix` - (Optional) Prefix identifying one or more objects to which the rule applies. The rule applies only to objects with key names beginning with the specified prefix. If you specify multiple rules in a lifecycle configuration, the rule with the longest prefix is used to determine which rule applies to a given object. If you specify a prefix for the rule, you cannot specify a tag for the rule.
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



