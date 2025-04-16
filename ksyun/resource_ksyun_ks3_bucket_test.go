package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccKsyunKS3ResourceCreate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKS3BucketConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("ksyun_ks3_bucket.bucket-create"),
					testAccCheckIDExists("ksyun_ks3_bucket.bucket-create2"),
				),
			},
		},
	})
}

const testAccKS3BucketConfig = `
provider "ksyun" {
  #指定KS3服务的访问域名
  endpoint = "ks3-cn-beijing.ksyuncs.com"
}

resource "ksyun_ks3_bucket" "bucket-create" {
  #指定要创建的存储桶的名称
  bucket = "bucket-20240206-104450"
}

resource "ksyun_ks3_bucket" "bucket-create2" {
  #指定要创建的存储桶的名称
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
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`
