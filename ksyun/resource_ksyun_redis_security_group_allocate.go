package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// redis security group allocate
func resourceRedisSecurityGroupAllocate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupAllocateCreate,
		Delete: resourceRedisSecurityGroupAllocateDelete,
		Read:   resourceRedisSecurityGroupAllocateRead,
		Update: resourceRedisSecurityGroupAllocateUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				return conflictResourceImport("security_group_id", "cache_id", "cache_ids", d)
			},
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cache_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"cache_ids"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return conflictResourceDiffSuppressForSingle("cache_ids", old, new, d)
				},
			},
			"cache_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:           schema.HashString,
				Deprecated:    "Use resourceRedisSecurityGroup().cache_ids instead",
				ConflictsWith: []string{"cache_id"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return conflictResourceDiffSuppressForMultiple("cache_id", "cache_ids", d)
				},
			},
		},
	}
}

func resourceRedisSecurityGroupAllocateCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		use string
		err error
	)
	use, err = checkConflictOnCreate("cache_id", "cache_ids", d)
	if err != nil {
		return err
	}

	transform := map[string]SdkReqTransform{
		"cache_ids": {
			mapping: "CacheId",
			Type:    TransformWithN,
		},
		"cache_id": {
			mapping: "CacheId",
			Type:    TransformSingleN,
		},
	}
	err = processRedisSecurityGroupAllocate(d, meta, transform, false, d.Get("security_group_id").(string))
	if err != nil {
		return fmt.Errorf("error on allocate redis security group: %s", err)
	}

	conflictResourceSetId(use, "security_group_id", "cache_id", "cache_ids", d)
	return resourceRedisSecurityGroupAllocateRead(d, meta)
}

func resourceRedisSecurityGroupAllocateUpdate(d *schema.ResourceData, meta interface{}) error {
	err := updateRedisSecurityGroupAllocate(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	return resourceRedisSecurityGroupAllocateRead(d, meta)
}

func resourceRedisSecurityGroupAllocateDelete(d *schema.ResourceData, meta interface{}) error {
	if checkMultipleExist("cache_ids", d) {
		return deallocateSecurityGroup(d, meta, d.Get("security_group_id").(string), SchemaSetToStringSlice(d.Get("cache_ids")), false)
	} else {
		return deallocateSecurityGroup(d, meta, d.Get("security_group_id").(string), []string{d.Get("cache_id").(string)}, false)
	}

}

func resourceRedisSecurityGroupAllocateRead(d *schema.ResourceData, meta interface{}) error {
	data, err := readRedisSecurityGroupAllocate(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		if validateRedisSgExists(err) {
			d.SetId("")
			return nil
		}
	}
	extra := make(map[string]SdkResponseMapping)
	if checkMultipleExist("cache_ids", d) {
		if err != nil {
			return err
		}
		extra["list"] = SdkResponseMapping{
			Field:         "cache_ids",
			FieldRespFunc: redisSgAllocateFieldRespFunc(d),
		}
	} else {
		if !checkValueInSliceMap(data["list"].([]interface{}), "id", d.Get("cache_id")) {
			d.SetId("")
			return nil
		}
		extra["list"] = SdkResponseMapping{
			Field: "cache_id",
			FieldRespFunc: func(i interface{}) interface{} {
				var id string
				for _, v := range i.([]interface{}) {
					if v.(map[string]interface{})["id"].(string) == d.Get("cache_id").(string) {
						id = v.(map[string]interface{})["id"].(string)
						break
					}
				}
				return id
			},
		}
	}

	SdkResponseAutoResourceData(d, resourceRedisSecurityGroupAllocate(), data, extra)
	return nil
}
