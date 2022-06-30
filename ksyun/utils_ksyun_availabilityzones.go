package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type requestFunc func() (*map[string]interface{}, error)
type isExist func(*map[string]interface{}) bool

type IntegrationRedisAzConf struct {
	resourceData *schema.ResourceData
	client       *KsyunClient
	requestFunc  requestFunc
	req          *map[string]interface{}
	field        string
	existFn      isExist
}

func (conf *IntegrationRedisAzConf) integrationRedisAz() (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	//currentRegion := *conf.client.kcsv1conn.Config.Region
	//if v, ok := conf.resourceData.GetOk(conf.field); ok {
	//	(*conf.req)[Downline2Hump(conf.field)] = v
	//	return conf.requestFunc()
	//}
	//conn := conf.client.kcsv1conn
	//req := make(map[string]interface{})
	//resp, err = conn.DescribeAvailabilityZones(&req)
	//if err != nil {
	//	return resp, err
	//}
	//obj, _ := getSdkValue("AvailabilityZoneSet", *resp)
	//for _, az := range obj.([]interface{}) {
	//	region := az.(map[string]interface{})["Region"]
	//	if region != currentRegion {
	//		continue
	//	}
	//	zone := az.(map[string]interface{})["AvailabilityZone"]
	//	(*conf.req)[Downline2Hump(conf.field)] = zone
	//	resp, err = conf.requestFunc()
	//	if err == nil && (conf.existFn == nil || (conf.existFn != nil && conf.existFn(resp))) {
	//		_ = conf.resourceData.Set(conf.field, zone)
	//		return resp, err
	//	} else {
	//		continue
	//	}
	//}
	if v, ok := conf.resourceData.GetOk(conf.field); ok {
		(*conf.req)[Downline2Hump(conf.field)] = v
	}
	resp, _ = conf.requestFunc()
	if err == nil && (conf.existFn == nil || (conf.existFn != nil && conf.existFn(resp))) {
		_ = conf.resourceData.Set(conf.field, (*conf.req)["AvailableZone"])
		return resp, err
	}
	return resp, fmt.Errorf(" not exist ")

}
