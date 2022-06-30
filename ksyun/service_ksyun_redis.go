package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func resourceRedisInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Id()

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &readReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheParameters(&readReq)
		},
	}

	action := "DescribeCacheParameters"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return fmt.Errorf("error on reading instance parameter %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		return nil
	}

	readReq = make(map[string]interface{})
	readReq["ParamVersion"] = d.Get("protocol")
	integrationAzConf = &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &readReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheDefaultParameters(&readReq)
		},
	}
	action = "DescribeCacheDefaultParameters"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return fmt.Errorf("error on reading default parameter %q, %s", d.Id(), err)
	}
	defaultData := (*resp)["Data"].([]interface{})
	if len(defaultData) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	defaultResult := make(map[string]interface{})
	parameter := make(map[string]interface{})
	for _, d1 := range data {
		param := d1.(map[string]interface{})
		result[param["name"].(string)] = fmt.Sprintf("%v", param["currentValue"])
	}
	for _, d1 := range defaultData {
		param := d1.(map[string]interface{})
		defaultResult[param["name"].(string)] = fmt.Sprintf("%v", param["defaultValue"])
	}
	localParams := d.Get("parameters").(map[string]interface{})
	if len(localParams) < 1 {
		for k, v := range result {
			if v1, ok := defaultResult[k]; ok {
				if v != v1 {
					parameter[k] = v
				}
			}
		}
	} else {
		for k := range localParams {
			if v, ok := result[k]; ok {
				parameter[k] = v
			}
		}
	}

	if err := d.Set("parameters", parameter); err != nil {
		return fmt.Errorf("error set data %v :%v", result, err)
	}
	return nil
}

func checkRedisInstanceStatus(d *schema.ResourceData, meta interface{}, timeout time.Duration, id string) error {
	var err error
	if id == "" {
		id = d.Id()
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"2"},
		Refresh:    stateRefreshForRedis(d, meta, []string{"2"}, id),
		Timeout:    timeout,
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func stateRefreshForRedis(d *schema.ResourceData, meta interface{}, target []string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			resp *map[string]interface{}
			item map[string]interface{}
			ok   bool
			err  error
		)

		resp, err = describeRedisInstance(d, meta, id)
		if err != nil {
			return nil, "", err
		}
		if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
			return nil, "", fmt.Errorf("no instance information was queried.%s", "")
		}
		status := int(item["status"].(float64))
		serviceStatus := int(item["serviceStatus"].(float64))
		// instance status error
		if status == 0 || status == 99 {
			return nil, "", fmt.Errorf("instance create error,status:%v", status)
		}
		// trade instance status error
		if serviceStatus == 3 {
			return nil, "", fmt.Errorf("instance create error,serviceStatus:%v", serviceStatus)
		}
		state := strconv.Itoa(status)
		for k, v := range target {
			if v == state && serviceStatus == 2 {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func resourceRedisInstanceSgRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	querySg := make(map[string]interface{})
	querySg["CacheId"] = d.Id()

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &querySg,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeSecurityGroups(&querySg)
		},
	}

	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return err
	}
	logger.Debug(logger.ReqFormat, "Demo", *resp)
	if item, ok := (*resp)["Data"].(map[string]interface{}); ok {
		var itemSetSlice []string
		sgIds := ""
		if sgs, ok := item["list"].([]interface{}); ok {
			for _, sg := range sgs {
				if info, ok := sg.(map[string]interface{}); ok {
					sgIds = sgIds + info["securityGroupId"].(string) + ","
					itemSetSlice = append(itemSetSlice, info["securityGroupId"].(string))
				}
			}
		}
		if len(sgIds) > 0 {
			err = d.Set("security_group_id", sgIds[0:len(sgIds)-1])
		}
		if err != nil {
			return err
		}
		//err =  d.Set("security_group_ids", itemSetSlice)
	}
	return err
}

func describeRedisInstance(d *schema.ResourceData, meta interface{}, id string) (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	queryReq := make(map[string]interface{})
	if id == "" {
		id = d.Id()
	}
	queryReq["CacheId"] = id

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &queryReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheCluster(&queryReq)
		},
	}
	action := "DescribeCacheCluster"
	logger.Debug(logger.ReqFormat, action, queryReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return resp, err
	}
	logger.Debug(logger.RespFormat, action, queryReq, *resp)
	return resp, err
}

func validateExists(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "cannot be found") || strings.Contains(strings.ToLower(err.Error()), "invalidaction")
}

func modifyRedisInstanceNameAndProject(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"name":           {},
		"iam_project_id": {mapping: "ProjectId"},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}
	//modify project
	err = ModifyProjectInstance(d.Id(), &req, meta)
	if err != nil {
		return fmt.Errorf("error on updating instance iam_project_id , error is %s", err)
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.RenameCacheCluster(&req)
			},
		}
		action := "RenameCacheCluster"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on rename instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
	}
	return err
}

func modifyRedisInstancePassword(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"pass_word": {mapping: "Password"},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.UpdatePassword(&req)
			},
		}
		action := "UpdatePassword"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on UpdatePassword instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
	}
	return err
}

func modifyRedisInstanceSpec(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"capacity":   {},
		"shard_size": {},
		"shard_num":  {},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.ResizeCacheCluster(&req)
			},
		}
		action := "ResizeCacheCluster"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on ResizeCacheCluster instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutUpdate), "")
		if err != nil {
			return fmt.Errorf("error on ResizeCacheCluster instance %q, %s", d.Id(), err)
		}
	}
	return err
}

func modifyRedisInstanceSg(d *schema.ResourceData, meta interface{}, isUpdate bool) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"security_group_id": {},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), isUpdate, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}

	if len(req) > 0 {
		var (
			arrayOldSg []interface{}
			arrayNewSg []interface{}
		)
		oldSg, newSg := d.GetChange("security_group_id")
		for _, v := range strings.Split(oldSg.(string), ",") {
			if v != "" {
				arrayOldSg = append(arrayOldSg, v)
			}
		}
		for _, v := range strings.Split(newSg.(string), ",") {
			if v != "" {
				arrayNewSg = append(arrayNewSg, v)
			}
		}
		oi := schema.NewSet(schema.HashString, arrayOldSg)
		ni := schema.NewSet(schema.HashString, arrayNewSg)

		removeSgs := oi.Difference(ni).List()
		newSgs := ni.Difference(oi).List()

		if len(removeSgs) > 0 {
			for _, sg := range removeSgs {
				removeReq := make(map[string]interface{})
				removeReq["CacheId.1"] = d.Id()
				removeReq["SecurityGroupId"] = sg
				integrationAzConf := &IntegrationRedisAzConf{
					resourceData: d,
					client:       meta.(*KsyunClient),
					req:          &removeReq,
					field:        "available_zone",
					requestFunc: func() (*map[string]interface{}, error) {
						conn := meta.(*KsyunClient).kcsv1conn
						return conn.DeallocateSecurityGroup(&removeReq)
					},
				}
				action := "DeallocateSecurityGroup"
				logger.Debug(logger.ReqFormat, action, removeReq)
				resp, err = integrationAzConf.integrationRedisAz()
				if err != nil {
					return fmt.Errorf("error on DeallocateSecurityGroup instance %q, %s", d.Id(), err)
				}
				logger.Debug(logger.RespFormat, action, removeReq, *resp)
			}

		}
		if len(newSgs) > 0 {
			addReq := make(map[string]interface{})
			addReq["CacheId.1"] = d.Id()
			count := 1
			for _, sg := range newSgs {
				addReq["SecurityGroupId."+strconv.Itoa(count)] = sg
				count = count + 1
			}
			integrationAzConf := &IntegrationRedisAzConf{
				resourceData: d,
				client:       meta.(*KsyunClient),
				req:          &addReq,
				field:        "available_zone",
				requestFunc: func() (*map[string]interface{}, error) {
					conn := meta.(*KsyunClient).kcsv1conn
					return conn.AllocateSecurityGroup(&addReq)
				},
			}
			action := "AllocateSecurityGroup"
			logger.Debug(logger.ReqFormat, action, addReq)
			resp, err = integrationAzConf.integrationRedisAz()
			if err != nil {
				return fmt.Errorf("error on AllocateSecurityGroup instance %q, %s", d.Id(), err)
			}
			logger.Debug(logger.RespFormat, action, addReq, *resp)
		}

	}
	return err
}

func modifyRedisInstanceAutoBackup(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"timing_switch": {},
		"timezone":      {},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}
	if req["Timezone"] != nil {
		if ts, ok := d.GetOk("timing_switch"); ok {
			req["TimingSwitch"] = ts
		} else {
			return fmt.Errorf("error on updating instance, error is %s", err)
		}
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.SetTimingSnapshot(&req)
			},
		}
		action := "SetTimingSnapshot"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on SetTimingSnapshot instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutUpdate), "")
		if err != nil {
			return fmt.Errorf("error on SetTimingSnapshot instance %q, %s", d.Id(), err)
		}
	}
	return err
}

func resourceRedisInstanceParameterCheckAndPrepare(d *schema.ResourceData, meta interface{}, isUpdate bool) (*map[string]interface{}, error) {
	var (
		reset bool
		resp  *map[string]interface{}
		err   error
		index int
	)
	req := make(map[string]interface{})

	parameters := make(map[string]string)
	if !isUpdate || (isUpdate && d.HasChange("parameters")) {
		if data, ok := d.GetOk("parameters"); ok {
			for k, v := range data.(map[string]interface{}) {
				parameters[k] = v.(string)
			}
		}
	}

	if mode, ok := d.GetOk("mode"); ok {
		if mode == 3 {
			if slave_num, ok := d.GetOk("slave_num"); ok {
				a := slave_num.(int)
				if a > 0 {
					return nil, fmt.Errorf("redis SelfDefineCluster do not support slave. when mode is 3, then slave_num must equal 0 ")
				}
			}
		}
	}

	//reset_all_parameters and parameters Conflict
	if r, ok := d.GetOk("reset_all_parameters"); ok && !isUpdate && r.(bool) && len(parameters) > 0 {
		err = fmt.Errorf("parameters is not empty,can not set reset_all_parameters true")
		return &req, err
	}
	if r, ok := d.GetOk("reset_all_parameters"); ok && isUpdate && r.(bool) {
		if data, ok := d.GetOk("parameters"); ok {
			if len(data.(map[string]interface{})) > 0 {
				err = fmt.Errorf("parameters is not empty,can not set reset_all_parameters true")
				return &req, err
			}
		}
		if d.HasChange("reset_all_parameters") {
			reset = true
		}

	}

	// condition on reset_all_parameters
	if isUpdate && reset {
		req["ResetAllParameters"] = reset

		return &req, d.Set("reset_all_parameters", reset)
	}

	//condition on set parameters, check parameter key and value valid
	defaultReq := map[string]interface{}{
		"ParamVersion": d.Get("protocol"),
	}
	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &defaultReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheDefaultParameters(&defaultReq)
		},
	}
	action := "DescribeCacheDefaultParameters"
	logger.Debug(logger.ReqFormat, action, defaultReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return &req, fmt.Errorf("error on DescribeCacheDefaultParameters: %s", err)
	}
	data, err := getSdkValue("Data", *resp)
	if err != nil {
		return &req, fmt.Errorf("error on DescribeCacheDefaultParameters: %s", err)
	}
	defaultParameters := make(map[string]interface{})
	for _, item := range data.([]interface{}) {
		name, err := getSdkValue("name", item)
		if err != nil {
			continue
		}
		defaultParameters[name.(string)] = item
	}
	// query current parameter
	cacheParameters := make(map[string]string)
	if d.Id() != "" {
		reqParam := make(map[string]interface{})
		reqParam["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &reqParam,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.DescribeCacheParameters(&reqParam)
			},
		}
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return &req, fmt.Errorf("error on DescribeCacheParameters: %s", err)
		}
		data, err = getSdkValue("Data", *resp)
		for _, item := range data.([]interface{}) {
			name, err := getSdkValue("name", item)
			if err != nil {
				continue
			}
			currentValue, err := getSdkValue("currentValue", item)
			cacheParameters[name.(string)] = currentValue.(string)
		}
	}
	for k, v := range parameters {
		if _, ok := defaultParameters[k]; !ok {
			return &req, fmt.Errorf("error on paramerter %v not support", k)
		}
		paramType, err := getSdkValue("validity.type", defaultParameters[k])
		if err != nil {
			continue
		}
		switch paramType.(string) {
		case "enum":
			values, err := getSdkValue("validity.values", defaultParameters[k])
			if err != nil {
				break
			}
			valid := false
			for _, v1 := range values.([]interface{}) {
				if v1.(string) == v {
					valid = true
				}
			}
			if !valid {
				return &req, fmt.Errorf("error on paramerter %v value must in  %v ", k, values)
			}
		case "range":
			minStr, err := getSdkValue("validity.min", defaultParameters[k])
			if err != nil {
				break
			}
			maxStr, err := getSdkValue("validity.max", defaultParameters[k])
			if err != nil {
				break
			}
			min, err := strconv.Atoi(minStr.(string))
			if err != nil {
				break
			}
			max, err := strconv.Atoi(maxStr.(string))
			if err != nil {
				break
			}
			current, err := strconv.Atoi(v)
			if err != nil {
				return &req, fmt.Errorf("error on paramerter %v value must number ", k)
			}
			if current > max || current < min {
				return &req, fmt.Errorf("error on paramerter %v value must in %v,%v ", k, minStr, maxStr)
			}
		case "regexp":
			value, err := getSdkValue("validity.value", defaultParameters[k])
			if err != nil {
				break
			}
			reg := regexp.MustCompile(value.(string))
			if reg == nil {
				continue
			}
			if !reg.MatchString(v) {
				return &req, fmt.Errorf("error on paramerter %v value must match %v ", k, value)
			}
		default:
			break
		}

		if cv, ok := cacheParameters[k]; ok && cv == v {
			continue
		}
		index = index + 1
		req[fmt.Sprintf("%v%v", "Parameters.ParameterName.", index)] = k
		req[fmt.Sprintf("%v%v", "Parameters.ParameterValue.", index)] = v
	}

	return &req, d.Set("reset_all_parameters", reset)
}

func setResourceRedisInstanceParameter(d *schema.ResourceData, meta interface{}, createReq *map[string]interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	(*createReq)["CacheId"] = d.Id()
	(*createReq)["Protocol"] = d.Get("protocol")

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          createReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.SetCacheParameters(createReq)
		},
	}

	action := "SetCacheParameters"
	logger.Debug(logger.ReqFormat, action, *createReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return fmt.Errorf("error on set instance parameter: %s", err)
	}
	logger.Debug(logger.RespFormat, action, *createReq, *resp)
	err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutUpdate), "")
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	return err
}

func redisSgAllocateFieldRespFunc(d *schema.ResourceData) FieldRespFunc {
	return func(i interface{}) interface{} {
		var caches []string
		if cacheSet, ok := d.Get("cache_ids").(*schema.Set); ok {
			for _, v := range i.([]interface{}) {
				c := v.(map[string]interface{})["id"].(string)
				if cacheSet.Contains(c) {
					caches = append(caches, v.(map[string]interface{})["id"].(string))
				}
			}
		}
		return caches
	}
}

func validateRedisSgExists(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not exist")
}

func processRedisSecurityGroupRule(d *schema.ResourceData, meta interface{}, transform map[string]SdkReqTransform, isUpdate bool, sgId string) error {
	var (
		req    map[string]interface{}
		resp   *map[string]interface{}
		err    error
		action string
	)
	req, err = SdkRequestAutoMapping(d, resourceRedisSecurityGroup(), false, transform, nil)
	if len(req) > 0 {
		conn := meta.(*KsyunClient).kcsv1conn
		if sgId == "" {
			sgId = d.Id()
		}
		//read one time and merge available_zone
		resp, err = readRedisSecurityGroup(d, meta, sgId)
		req["SecurityGroupId"] = sgId
		req["AvailableZone"] = d.Get("available_zone")
		if !isUpdate {
			action = "CreateSecurityGroupRule"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.CreateSecurityGroupRule(&req)
			if err != nil {
				return fmt.Errorf("error on add redis security group rules: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		} else {
			action = "DeleteSecurityGroupRule"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.DeleteSecurityGroupRule(&req)
			if err != nil {
				return fmt.Errorf("error on delete redis security group rules: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		}

	}
	return err
}

func processRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}, transform map[string]SdkReqTransform, isUpdate bool, sgId string) error {
	var (
		req    map[string]interface{}
		resp   *map[string]interface{}
		err    error
		action string
	)
	req, err = SdkRequestAutoMapping(d, resourceRedisSecurityGroup(), false, transform, nil)
	if len(req) > 0 {
		conn := meta.(*KsyunClient).kcsv1conn
		if sgId == "" {
			sgId = d.Id()
		}
		//read one time and merge available_zone
		resp, err = readRedisSecurityGroup(d, meta, sgId)
		if err != nil {
			return err
		}
		req["AvailableZone"] = d.Get("available_zone")
		if !isUpdate {
			req["SecurityGroupId.1"] = sgId
			action = "AllocateSecurityGroup"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.AllocateSecurityGroup(&req)
			if err != nil {
				return fmt.Errorf("error on allocate security group to redis: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		} else {
			req["SecurityGroupId"] = sgId
			action = "DeallocateSecurityGroup"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.DeallocateSecurityGroup(&req)
			if err != nil {
				return fmt.Errorf("error on deallocateSecurityGroup  security group to redis: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		}

	}
	return err
}

func deallocateSecurityGroup(d *schema.ResourceData, meta interface{}, sgId string, cacheIds []string, all bool) error {
	var (
		err  error
		resp *map[string]interface{}
	)
	_, err = readRedisSecurityGroup(d, meta, sgId)
	if err != nil {
		return err
	}
	if sgId == "" {
		sgId = d.Id()
	}
	deallocateReq := make(map[string]interface{})
	deallocateReq["SecurityGroupId"] = sgId
	if !all {
		for i, id := range cacheIds {
			deallocateReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = id
		}
	} else {
		allocate, err := readRedisSecurityGroupAllocate(d, meta, "")
		if err != nil {
			return err
		}
		for i, a := range allocate["list"].([]interface{}) {
			deallocateReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = a.(map[string]interface{})["id"]
		}
	}

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &deallocateReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DeallocateSecurityGroup(&deallocateReq)
		},
	}
	action := "DeallocateSecurityGroup"
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deallocateReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err == nil {
			data := (*resp)["Data"].([]interface{})
			if len(data) == 0 {
				return nil
			}
		}
		if err != nil && inUseError(err) {
			return resource.RetryableError(err)
		}
		return nil
	})
}

func readRedisSecurityGroup(d *schema.ResourceData, meta interface{}, securityGroupId string) (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	req := make(map[string]interface{})
	if securityGroupId == "" {
		securityGroupId = d.Id()
	}
	req["SecurityGroupId"] = securityGroupId
	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &req,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeSecurityGroup(&req)
		},
		existFn: func(i *map[string]interface{}) bool {
			v, _ := getSdkValue("Data", *i)
			if v == nil || len(v.(map[string]interface{})) == 0 {
				return false
			}
			return true
		},
	}
	action := "DescribeSecurityGroup"
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return resp, fmt.Errorf("error on reading redis security group %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	return resp, err
}

func readRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}, sgId string) (map[string]interface{}, error) {
	var (
		resp      *map[string]interface{}
		err       error
		instances []interface{}
	)
	currentCount := int64(0)
	total := int64(0)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	if sgId == "" {
		sgId = d.Id()
	}
	readReq["SecurityGroupId"] = sgId
	readReq["Limit"] = 1000
	readReq["FilterCache"] = true
	for {
		readReq["Offset"] = currentCount
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &readReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				return conn.DescribeInstances(&readReq)
			},
		}

		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return nil, fmt.Errorf("error on reading redis security group allocate instances %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, readReq, *resp)
		data := (*resp)["Data"].(map[string]interface{})
		total = int64(data["total"].(float64))
		lists := data["list"].([]interface{})
		for _, v := range lists {
			instances = append(instances, v)
		}
		currentCount = int64(len(instances))
		if currentCount == total {
			data["list"] = instances
			return data, err
		}
	}
}

func updateRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}, sgId string) (err error) {
	//cache_ids
	if d.HasChange("cache_ids") {
		var (
			err      error
			oldArray []string
			newArray []string
			add      []interface{}
			del      []interface{}
		)
		_, err = readRedisSecurityGroup(d, meta, sgId)
		if err != nil {
			return err
		}

		o, n := d.GetChange("cache_ids")
		for _, v := range o.(*schema.Set).List() {
			oldArray = append(oldArray, v.(string))
		}
		for _, v := range n.(*schema.Set).List() {
			newArray = append(newArray, v.(string))
		}
		for _, a := range oldArray {
			exist := false
			for _, b := range newArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				del = append(del, a)
			}
		}
		for _, a := range newArray {
			exist := false
			for _, b := range oldArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				add = append(add, a)
			}
		}
		transformAdd := map[string]SdkReqTransform{
			"cache_ids": {
				mapping: "CacheId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(add) > 0 {
						return add, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupAllocate(d, meta, transformAdd, false, sgId)
		if err != nil {
			return err
		}
		transformDel := map[string]SdkReqTransform{
			"cache_ids": {
				mapping: "CacheId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(del) > 0 {
						return del, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupAllocate(d, meta, transformDel, true, sgId)
		if err != nil {
			return err
		}

	}
	return err
}

func updateRedisSecurityGroupRules(d *schema.ResourceData, meta interface{}, sgId string) (err error) {
	//rule
	if d.HasChange("rules") {
		var (
			oldArray []string
			newArray []string
			add      []interface{}
			del      []interface{}
			resp     *map[string]interface{}
		)
		resp, err = readRedisSecurityGroup(d, meta, sgId)
		if err != nil {
			return err
		}
		data := (*resp)["Data"].(map[string]interface{})
		rulesMap := make(map[string]interface{})
		//get rule id for del
		if rules, ok := data["rules"]; ok {
			for _, r := range rules.([]interface{}) {
				rule := r.(map[string]interface{})
				rulesMap[rule["cidr"].(string)] = rule["id"]
			}
		}
		o, n := d.GetChange("rules")
		for _, v := range o.(*schema.Set).List() {
			oldArray = append(oldArray, v.(string))
		}
		for _, v := range n.(*schema.Set).List() {
			newArray = append(newArray, v.(string))
		}
		for _, a := range oldArray {
			if _, ok := rulesMap[a]; !ok {
				continue
			}
			exist := false
			for _, b := range newArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				del = append(del, rulesMap[a])
			}
		}
		for _, a := range newArray {
			exist := false
			for _, b := range oldArray {
				if a == b {
					exist = true
					break
				}
			}
			if _, ok := rulesMap[a]; !ok && !exist {
				add = append(add, a)
			}
		}
		transformAdd := map[string]SdkReqTransform{
			"rules": {
				mapping: "Cidrs",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(add) > 0 {
						return add, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupRule(d, meta, transformAdd, false, sgId)
		if err != nil {
			return err
		}
		transformDel := map[string]SdkReqTransform{
			"rules": {
				mapping: "SecurityGroupRuleId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(del) > 0 {
						return del, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupRule(d, meta, transformDel, true, sgId)
		if err != nil {
			return err
		}

	}
	return err
}

func setRedisSgCidrs(current []interface{}, d *schema.ResourceData) (cidrs []string) {
	if ruleSet, ok := d.Get("rules").(*schema.Set); ok {
		for _, v := range current {
			c := v.(map[string]interface{})["cidr"].(string)
			if ruleSet.Contains(c) {
				cidrs = append(cidrs, c)
			}
		}
	}
	return cidrs
}
