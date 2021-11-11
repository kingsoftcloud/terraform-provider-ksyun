package ksyun

import (
	"errors"
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/kcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"strings"
	"time"
)

// instance node
func resourceRedisInstanceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisInstanceNodeCreate,
		Delete: resourceRedisInstanceNodeDelete,
		Read:   resourceRedisInstanceNodeRead,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				var err error
				importParts := strings.Split(d.Id(), ":")
				if len(importParts) < 2 {
					return nil, fmt.Errorf("import too few parts,must CacheId:NodeId")
				}
				d.SetId(importParts[1])
				err = d.Set("cache_id", importParts[0])
				if err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, err
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cache_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRedisInstanceNodeCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	//read cluster
	_, err = readRedisInstanceNodeCluster(d, meta)
	if err != nil {
		return fmt.Errorf("error on add Instance node: %s", err)
	}
	// create
	resp, err = createRedisInstanceNode(d, meta)
	if err != nil {
		return fmt.Errorf("error on add instance node: %s", err)
	}
	if resp != nil {
		_ = d.Set("instance_id", (*resp)["Data"].(map[string]interface{})["NodeId"].(string))
	}
	d.SetId(d.Get("instance_id").(string))

	err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutCreate), d.Get("cache_id").(string))
	if err != nil {
		return fmt.Errorf("error on add Instance node: %s", err)
	}

	return resourceRedisInstanceNodeRead(d, meta)
}

func resourceRedisInstanceNodeDelete(d *schema.ResourceData, meta interface{}) error {
	// delete
	deleteParamReq := make(map[string]interface{})
	deleteParamReq["CacheId"] = d.Get("cache_id")
	deleteParamReq["NodeId"] = d.Get("instance_id")

	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		var (
			resp *map[string]interface{}
			err  error
		)
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &deleteParamReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv2conn
				return conn.DeleteCacheSlaveNode(&deleteParamReq)
			},
		}
		action := "DeleteCacheSlaveNode"
		logger.Debug(logger.ReqFormat, action, deleteParamReq)
		resp, err = integrationAzConf.integrationRedisAz()
		logger.Debug(logger.RespFormat, action, deleteParamReq, resp)
		if err == nil {
			return nil
		}
		_, err = readRedisInstanceNode(d, meta)
		if err != nil {
			if validateRedisNodeExists(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(errors.New("deleting"))
	})

}

func validateRedisNodeExists(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not exist")
}

func readRedisInstanceNodeCluster(d *schema.ResourceData, meta interface{}) (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Get("cache_id")

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &readReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheCluster(&readReq)
		},
	}
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return resp, fmt.Errorf("error on reading instance node Cluster %q, %s", d.Id(), err)
	}
	return resp, err
}

func createRedisInstanceNode(d *schema.ResourceData, meta interface{}) (*map[string]interface{}, error) {
	var (
		createNodeReq map[string]interface{}
		resp          *map[string]interface{}
		err           error
	)
	createNodeReq, err = SdkRequestAutoMapping(d, resourceRedisInstanceNode(), false, nil, nil)
	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &createNodeReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv2conn
			return conn.AddCacheSlaveNode(&createNodeReq)
		},
	}
	action := "AddCacheSlaveNode"
	logger.Debug(logger.ReqFormat, action, createNodeReq)
	resp, err = integrationAzConf.integrationRedisAz()
	return resp, err
}

func readRedisInstanceNode(d *schema.ResourceData, meta interface{}) (*map[string]interface{}, error) {
	var (
		item interface{}
		resp *map[string]interface{}
		err  error
		ok   bool
	)
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Get("cache_id")

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &readReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv2conn
			return conn.DescribeCacheReadonlyNode(&readReq)
		},
	}
	action := "DescribeCacheReadonlyNode"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return resp, fmt.Errorf("error on reading instance node %q, %s", d.Id(), err)
	}
	if item, ok = (*resp)["Data"]; !ok {
		return resp, fmt.Errorf("error on reading instance node %s not exist", d.Id())
	}
	items, ok := item.([]interface{})
	if !ok || len(items) == 0 {
		return resp, fmt.Errorf("error on reading instance node %s not exist", d.Id())
	}
	for _, v := range items {
		vMap := v.(map[string]interface{})
		if d.Id() == vMap["instanceId"] {
			return &vMap, err
		}
	}
	return resp, fmt.Errorf("error on reading instance node %s not exist", d.Id())
}

func resourceRedisInstanceNodeRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	_, err = readRedisInstanceNodeCluster(d, meta)
	if err != nil {
		return err
	}

	resp, err = readRedisInstanceNode(d, meta)
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, resourceRedisInstanceNode(), *resp, nil)
	return nil
}

func stateRefreshForOperateNodeFunc(client *kcsv1.Kcsv1, az, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			resp *map[string]interface{}
			item map[string]interface{}
			ok   bool
			err  error
		)

		queryReq := map[string]interface{}{"CacheId": instanceId}
		queryReq["AvailableZone"] = az
		action := "DescribeCacheCluster"
		logger.Debug(logger.ReqFormat, action, queryReq)
		if resp, err = client.DescribeCacheCluster(&queryReq); err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, action, queryReq, *resp)
		if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
			return nil, "", fmt.Errorf("no instance information was queried.%s", "")
		}
		status := int(item["status"].(float64))
		if status == 0 || status == 99 {
			return nil, "", fmt.Errorf("instance operate error,status:%v", status)
		}
		state := strconv.Itoa(status)
		for k, v := range target {
			if v == state && int(item["transition"].(float64)) != 1 {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}

		return resp, state, nil
	}
}
