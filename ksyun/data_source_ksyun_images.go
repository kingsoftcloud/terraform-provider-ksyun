/*
This data source providers a list of available image resources according to their availability zone, image ID and other fields.

# Example Usage

```hcl

	data "ksyun_images" "default" {
	  output_file="output_result"
	  is_public=true
	  image_source="system"
	}

```
*/
package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func dataSourceKsyunImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunImagesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of image IDs.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				Description:  "A regex string to filter resulting images by name. (Such as: `^CentOS 7.[1-2] 64` means CentOS 7.1 of 64-bit operating system or CentOS 7.2 of 64-bit operating system, \"^Ubuntu 16.04 64\" means Ubuntu 16.04 of 64-bit operating system).",
			},
			"platform": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Platform type of the image system.",
			},
			"is_public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If ksyun provide the image.",
			},
			"image_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Valid values are import, copy, share, extend, system.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results (after running `terraform plan`).",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of image that satisfy the condition.",
			},

			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "It is a nested type which documented below.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the image.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of image.",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time of creation.",
						},
						"is_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If ksyun provide the image.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform type of the image system.",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the image.",
						},
						"is_npe": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether networking enhancement is support or not.",
						},
						"user_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User defined category.",
						},
						"sys_disk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "size of system disk.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the id of the instance which the image based on.",
						},
						"progress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "image creation progress percentage.",
						},
						"image_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image source of the image.",
						},
						"cloud_init_support": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support cloud-init.",
						},
						"ipv6_support": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support ipv6.",
						},
						"is_modify_type": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support live upgrade.",
						},
						"is_cloud_market": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether image is from cloud market or not.",
						},
						"real_image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The real id of the image.",
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunImagesRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).kecconn
	req := make(map[string]interface{})
	var imageIds []string
	var allImages []interface{}
	resp, err := conn.DescribeImages(&req)
	if err != nil {
		return fmt.Errorf("error on reading Image list req(%v):%v", req, err)
	}
	//logger.Debug("%v", "DescribeImages", resp, err)
	itemSet, ok := (*resp)["ImagesSet"]
	if !ok {
		return fmt.Errorf("error on reading Image set")
	}
	items, ok := itemSet.([]interface{})
	if !ok {
		return nil
	}
	if items == nil || len(items) < 1 {
		return nil
	}
	allImages = append(allImages, items...)
	datas := GetSubSliceDByRep(allImages, imageKeys)
	if name, ok := d.GetOk("platform"); ok {
		var dataFilter []map[string]interface{}
		for _, v := range datas {
			if v["platform"] == name {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	if standard, ok := d.GetOk("is_public"); ok {
		var dataFilter []map[string]interface{}
		for _, v := range datas {
			if v["is_public"] == standard {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	if imageSource, ok := d.GetOk("image_source"); ok {
		var dataFilter []map[string]interface{}
		for _, v := range datas {
			if v["image_source"] == imageSource {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		var dataFilter []map[string]interface{}
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range datas {
			if r == nil || r.MatchString(v["name"].(string)) {
				dataFilter = append(dataFilter, v)
			}
		}
		datas = dataFilter
	}
	err = dataSourceKscSave(d, "images", imageIds, datas)
	if err != nil {
		return fmt.Errorf("error on save Images list, %s", err)
	}
	return nil
}
