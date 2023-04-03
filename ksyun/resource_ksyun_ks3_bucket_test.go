package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKS3_Resource_new(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKS3BucketConfig,
			},
		},
	})
}

const testAccDataKS3BucketConfig = `


resource "ksyun_ks3_bucket" "bucket-new" {
  bucket = "ks3-bucket-tf2"
}
`

func TestAccKsyunKS3_Resource_attr(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKS3BucketConfigAttr,
			},
		},
	})
}

const testAccDataKS3BucketConfigAttr = `


resource "ksyun_ks3_bucket" "bucket-attr" {

  bucket   = var.bucket-new
  acl      = var.acl-bj

  logging {
    target_bucket = var.bucket-new
    target_prefix = var.target-prefix
  }
  lifecycle_rule {
    id      = var.role-date
    enabled = false

    transitions {
      date          = "2023-04-18"
      storage_class = "ARCHIVE"
    }
#    expiration  {
#      date = "2023-05-18"
#    }
  }

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

`
