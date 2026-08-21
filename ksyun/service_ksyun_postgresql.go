package ksyun

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func readPostgresqlInstances(d *schema.ResourceData, meta interface{}, condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp                      *map[string]interface{}
		postgresqlInstanceResults interface{}
	)
	conn := meta.(*KsyunClient).postgresqlconn
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeDBInstances(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeDBInstances(&condition)
		if err != nil {
			return data, err
		}
	}

	postgresqlInstanceResults, err = getSdkValue("Data.Instances", *resp)
	if err != nil {
		return data, err
	}
	data = postgresqlInstanceResults.([]interface{})
	return data, err
}

func readPostgresqlInstance(d *schema.ResourceData, meta interface{}, instanceId string) (data map[string]interface{}, err error) {
	var (
		postgresqlInstanceResults []interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"DBInstanceIdentifier": instanceId,
	}
	postgresqlInstanceResults, err = readPostgresqlInstances(d, meta, req)
	if err != nil {
		if isPostgresqlInstanceNotFoundError(err) {
			return data, nil
		}
		return data, err
	}
	for _, v := range postgresqlInstanceResults {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, nil
	}
	return data, err
}

func isPostgresqlInstanceNotFoundError(err error) bool {
	if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 404 {
		return true
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "dbproduct not found")
}

func readAndSetPostgresqlInstance(d *schema.ResourceData, meta interface{}, isRR bool) (err error) {
	data, err := readPostgresqlInstance(d, meta, "")
	if err != nil {
		return err
	}
	if len(data) == 0 {
		d.SetId("")
		return nil
	}
	// check rr or master
	dbInstanceType := data["DBInstanceType"]
	if dbInstanceType == "RR_PG" && !isRR {
		return fmt.Errorf("postgresql instance is read replica, please use ksyun_postgresql_rr ")
	}
	if dbInstanceType != "RR_PG" && isRR {
		return fmt.Errorf("postgresql instance is not  read replica, please use ksyun_postgresql ")
	}
	if _, ok := data["Eip"]; ok {
		data["instance_has_eip"] = true
	} else {
		data["instance_has_eip"] = false
	}
	extra := map[string]SdkResponseMapping{
		"DBInstanceClass": {
			Field: "db_instance_class",
			FieldRespFunc: func(i interface{}) interface{} {
				value := i.(map[string]interface{})
				return fmt.Sprintf("db.ram.%v|db.disk.%v", value["Ram"], value["Disk"])
			},
		},
		"DBInstanceName": {
			Field: "db_instance_name",
		},
		"DBInstanceType": {
			Field: "db_instance_type",
		},
		"DBParameterGroupId": {
			Field: "db_parameter_group_id",
		},
	}
	if _, ok := d.GetOk("tags"); ok {
		err = mergeTagsData(d, &data, meta.(*KsyunClient), "postgresql")
		if err != nil {
			return fmt.Errorf("reading tags error: %s", err)
		}
	}
	if isRR {
		extra["DBSource"] = SdkResponseMapping{
			Field: "db_instance_identifier",
			FieldRespFunc: func(i interface{}) interface{} {
				return i.(map[string]interface{})["DBInstanceIdentifier"]
			},
		}
		extra["AvailabilityZone"] = SdkResponseMapping{
			Field: "availability_zone_1",
		}
		delete(data, "DBInstanceIdentifier")
		SdkResponseAutoResourceData(d, resourceKsyunPostgresqlRr(), data, extra)
	} else {
		extra["DBInstanceIdentifier"] = SdkResponseMapping{
			Field: "db_instance_identifier",
		}
		extra["MasterAvailabilityZone"] = SdkResponseMapping{
			Field: "availability_zone_1",
		}
		extra["SlaveAvailabilityZone"] = SdkResponseMapping{
			Field: "availability_zone_2",
		}
		SdkResponseAutoResourceData(d, resourceKsyunPostgresql(), data, extra)
	}

	if d.Get("force_restart") != nil {
		_ = d.Set("force_restart", d.Get("force_restart"))
	} else {
		_ = d.Set("force_restart", false)
	}
	if dbInstanceClass, ok := data["DBInstanceClass"]; ok {
		if m, ok := dbInstanceClass.(map[string]interface{}); ok {
			if v, ok := m["Vcpus"]; ok {
				_ = d.Set("vcpus", v)
			}
		}
	}
	return err
}

func readPostgresqlInstanceParameters(d *schema.ResourceData, meta interface{}, instanceId string) (data map[string]interface{}, err error) {
	var (
		resp                              *map[string]interface{}
		postgresqlInstanceParameterResult interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"DBInstanceIdentifier": instanceId,
	}
	conn := meta.(*KsyunClient).postgresqlconn
	action := "DescribeDBInstanceParameters"

	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.DescribeDBInstanceParameters(&req)
	if err != nil {
		return data, err
	}

	postgresqlInstanceParameterResult, err = getSdkValue("Data.Parameters", *resp)
	if err != nil {
		return data, err
	}
	data = postgresqlInstanceParameterResult.(map[string]interface{})
	return data, err
}

// readPostgresqlParameterGroup will read the parameters of postgresql parameter group by id
func readPostgresqlParameterGroup(d *schema.ResourceData, meta interface{}, parameterGroupId string, isAll bool) (data map[string]interface{}, err error) {
	var (
		resp []interface{}
	)
	if parameterGroupId == "" {
		parameterGroupId = d.Get("db_parameter_group_id").(string)
	}
	req := map[string]interface{}{
		"DBParameterGroupId": parameterGroupId,
	}
	postgresqlParameterSrv := NewPostgresqlParameterSrv(meta.(*KsyunClient))
	action := "DescribeDBParameterGroup"

	logger.Debug(logger.ReqFormat, action, req)
	resp, err = postgresqlParameterSrv.describeDBParameterGroupById(req)
	if err != nil {
		return data, err
	}
	engineVersionIf, _ := getSdkValue("0.EngineVersion", resp)
	engineVersion, _ := If2String(engineVersionIf)
	// if ok is false that means this resource is not exist.
	if _engineVersion, ok := d.GetOk("engine_version"); ok {
		if !reflect.DeepEqual(engineVersion, _engineVersion) {
			return nil, fmt.Errorf("db parameter group engine version must be matched postgresql instance engine version")
		}
	}
	sdkValue, _ := getSdkValue("0.Parameters", resp)
	data, err = If2Map(sdkValue)
	if !isAll {
		return data, err
	}
	defaultParameter, err := readPostgresqlDefaultParametersWithMap(d, meta)
	if err != nil {
		return nil, err
	}
	// merge default
	for k, v := range data {
		defaultParameter[k] = v
	}
	return defaultParameter, nil
}

func parameterGroupTransformerPostgresql(d *schema.ResourceData, meta interface{}, parameters map[string]interface{}) (data map[string]interface{}, err error) {
	if parameters == nil {
		return nil, fmt.Errorf("the parameters of db parameter group is nil")
	}
	defaultParameter, err := readPostgresqlDefaultParameters(d, nil, meta)
	if err != nil {
		return data, err
	}
	data = make(map[string]interface{})
	for k, v := range parameters {
		t := "string"
		var val string
		if v1, ok := defaultParameter[k]; ok {
			t = v1.(map[string]interface{})["Type"].(string)
		}
		if vf, ok := v.(float64); ok {
			switch t {
			case "float":
				val = fmt.Sprintf("%v", strconv.FormatFloat(vf, 'f', 0, 64))
			default:
				val = fmt.Sprintf("%v", strconv.FormatInt(int64(vf), 10))
			}
		} else {
			val = fmt.Sprintf("%v", v)
		}

		data[k] = val
	}
	return data, err
}

func readPostgresqlDefaultParametersWithMap(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	defaultParameter, err := readPostgresqlDefaultParameters(d, nil, meta)
	if err != nil {
		return nil, err
	}

	defaultParameterMap := make(map[string]interface{})
	for k, v := range defaultParameter {
		val, _ := If2Map(v)
		defaultParameterMap[k] = val["Default"]
	}
	return defaultParameterMap, nil
}

func readAndSetPostgresqlInstanceParameters(d *schema.ResourceData, meta interface{}) (err error) {
	parameter, err := readPostgresqlInstanceParameters(d, meta, "")
	if err != nil {
		return err
	}
	defaultParameter, err := readPostgresqlDefaultParameters(d, nil, meta)
	if err != nil {
		return err
	}
	remote := make(map[string]map[string]interface{})
	var parameters []map[string]interface{}
	// 参考aws的做法，由于配置有很多默认项,这里只列出展示用户可能修改的配置参数
	local := d.Get("parameters").(*schema.Set)
	for k, v := range parameter {
		t := "string"
		m := make(map[string]interface{})
		if v1, ok := defaultParameter[k]; ok {
			t = v1.(map[string]interface{})["Type"].(string)
		}
		m["name"] = k
		if vf, ok := v.(float64); ok {
			switch t {
			case "float":
				m["value"] = fmt.Sprintf("%v", strconv.FormatFloat(vf, 'g', 1, 64))
			default:
				m["value"] = fmt.Sprintf("%v", strconv.FormatInt(int64(vf), 10))
			}
		} else {
			m["value"] = fmt.Sprintf("%v", v)
		}

		remote[k] = m
	}
	// if local.Len() < 1 {
	// 	for k, v := range remote {
	// 		if v1, ok := defaultParameter[k]; ok {
	// 			switch v1.(map[string]interface{})["Type"] {
	// 			case "integer":
	// 				if strconv.FormatInt(int64(v1.(map[string]interface{})["Default"].(float64)), 10) != v["value"] {
	// 					parameters = append(parameters, v)
	// 				}
	// 			case "float":
	// 				if strconv.FormatFloat(v1.(map[string]interface{})["Default"].(float64), 'g', 1, 64) != v["value"] {
	// 					parameters = append(parameters, v)
	// 				}
	// 			case "expression":
	// 				if v1.(map[string]interface{})["Variable"] == "instance_memory" {
	// 					ramStr := strings.Replace(strings.Split(d.Get("db_instance_class").(string), "|")[0], "db.ram.", "", -1)
	// 					scaleStr := strings.Replace(v1.(map[string]interface{})["DefaultScaleFactor"].(string), "%", "", -1)
	// 					ram, _ := strconv.ParseInt(ramStr, 10, 64)
	// 					scale, _ := strconv.ParseFloat(scaleStr, 64)
	// 					defaultScaleValue := strconv.FormatInt(int64(float64(ram*1024*1024*1024)*(scale/100)), 10)
	// 					if defaultScaleValue != v["value"] {
	// 						parameters = append(parameters, v)
	// 					}
	// 				} else {
	// 					parameters = append(parameters, v)
	// 				}
	// 			default:
	// 				if v1.(map[string]interface{})["Default"] != v["value"] {
	// 					parameters = append(parameters, v)
	// 				}
	// 			}
	//
	// 		}
	// 	}
	// } else {

	// temporary parameter should keep identical with online instance parameters
	for _, value := range local.List() {
		name := value.(map[string]interface{})["name"]
		for k, v := range remote {
			if k == name {
				parameters = append(parameters, v)
				break
			}
		}
	}
	// }

	// if local, ok := d.GetOk("parameters"); ok {
	// }
	err = d.Set("parameters", parameters)
	return err
}

func readPostgresqlDefaultParameters(d *schema.ResourceData, diff *schema.ResourceDiff, meta interface{}) (data map[string]interface{}, err error) {
	var (
		resp       *map[string]interface{}
		parameters interface{}
		req        map[string]interface{}
	)
	if d != nil {
		req = map[string]interface{}{
			"Engine":        d.Get("engine"),
			"EngineVersion": d.Get("engine_version"),
		}
	} else if diff != nil {
		req = map[string]interface{}{
			"Engine":        diff.Get("engine"),
			"EngineVersion": diff.Get("engine_version"),
		}
	}

	conn := meta.(*KsyunClient).postgresqlconn
	data = make(map[string]interface{})
	action := "DescribeEngineDefaultParameters"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.DescribeEngineDefaultParameters(&req)
	if err != nil {
		return data, err
	}
	parameters, err = getSdkValue("Data.Parameters", *resp)
	if err != nil {
		return data, err
	}
	for k, v := range parameters.(map[string]interface{}) {
		data[k] = v
	}
	return data, err
}

func checkAndProcessPostgresqlParameters(d *schema.ResourceData, meta interface{}) (req map[string]interface{}, needRestart bool, err error) {
	req = make(map[string]interface{})
	if d.HasChange("parameters") {
		var (
			oldParameters interface{}
			newParameters interface{}
			needToDefault []string
			needToSet     []string
		)
		oldKv := make(map[string]string)
		newKv := make(map[string]string)
		index := 1
		oldParameters, newParameters = d.GetChange("parameters")
		// init from list to kv
		if params, ok := oldParameters.(*schema.Set); ok {
			for _, i := range params.List() {
				oldKv[i.(map[string]interface{})["name"].(string)] = i.(map[string]interface{})["value"].(string)
			}
		}
		if params, ok := newParameters.(*schema.Set); ok {
			for _, i := range params.List() {
				newKv[i.(map[string]interface{})["name"].(string)] = i.(map[string]interface{})["value"].(string)
			}
		}
		// compare add or remove
		for k := range oldKv {
			if _, ok := newKv[k]; !ok {
				needToDefault = append(needToDefault, k)
			}
		}
		for k, n := range newKv {
			if o, ok := oldKv[k]; ok {
				if n != o {
					needToSet = append(needToSet, k)
				}
			} else {
				needToSet = append(needToSet, k)
			}

		}

		// check and prepare if value equals default ,the param will ignore
		needRestart, index, err = prepareModifyDbParameterParamsPostgresql(d, meta, needToDefault, &req, oldKv, index, true)
		if err != nil {
			return req, needRestart, err
		}
		needRestart, index, err = prepareModifyDbParameterParamsPostgresql(d, meta, needToSet, &req, newKv, index, false)
		if err != nil {
			return req, needRestart, err
		}
	}
	return req, needRestart, err
}

func prepareModifyDbParameterParamsPostgresql(d *schema.ResourceData, meta interface{}, keys []string, req *map[string]interface{}, kv map[string]string, index int, toDefault bool) (needRestart bool, num int, err error) {
	var (
		defaults map[string]interface{}
		currents map[string]interface{}
	)
	defaults = make(map[string]interface{})
	currents = make(map[string]interface{})
	num = index
	// readPostgresqlDefaultParameters
	defaults, err = readPostgresqlDefaultParameters(d, nil, meta)
	if err != nil {
		return needRestart, num, err
	}
	// read current
	if d.Id() == "" {
		currents = make(map[string]interface{})
	} else if v, ok := d.GetOk("resource_name"); ok {
		resourceName := v.(string)
		if resourceName == ResourcePostgresqlParameterGroup {
			currents, err = readPostgresqlParameterGroup(d, meta, "", true)
			if err != nil {
				return needRestart, num, err
			}
		} // added: to deal with the current parameters of postgresql parameter group
	} else {
		currents, err = readPostgresqlInstanceParameters(d, meta, "")
		if err != nil {
			return needRestart, num, err
		}
	}

	for _, key := range keys {
		if defaultObj, ok := defaults[key]; ok {
			dv := defaultObj.(map[string]interface{})["Default"]
			dt := defaultObj.(map[string]interface{})["Type"].(string)
			dr := defaultObj.(map[string]interface{})["RestartRequired"].(bool)
			switch dt {
			case "integer":
				dMin := int64(defaultObj.(map[string]interface{})["Min"].(float64))
				dMax := int64(defaultObj.(map[string]interface{})["Max"].(float64))
				valueNum, err := strconv.ParseInt(kv[key], 10, 64)
				if err != nil {
					return needRestart, num, err
				}
				if toDefault && valueNum == int64(dv.(float64)) {
					continue
				}
				if !toDefault && len(currents) > 0 && valueNum == int64(currents[key].(float64)) {
					continue
				}
				if valueNum < dMin || valueNum > dMax {
					return needRestart, num, fmt.Errorf("parameter %s must in [%d , %d]", key, dMin, dMax)
				}
			case "float":
				dMin := defaultObj.(map[string]interface{})["Min"].(float64)
				dMax := defaultObj.(map[string]interface{})["Max"].(float64)
				valueNum, err := strconv.ParseFloat(kv[key], 64)

				if err != nil {
					return needRestart, num, err
				}
				if toDefault && valueNum == dv.(float64) {
					continue
				}
				if !toDefault && len(currents) > 0 && valueNum == currents[key].(float64) {
					continue
				}
				if valueNum < dMin || valueNum > dMax {
					return needRestart, num, fmt.Errorf("parameter %s must in [%f , %f]", key, dMin, dMax)
				}
				// continue
			case "expression":
				if defaultObj.(map[string]interface{})["Variable"] == "instance_memory" {
					dMin := int64(defaultObj.(map[string]interface{})["Min"].(float64))
					scaleStr := strings.Replace(defaultObj.(map[string]interface{})["DefaultScaleFactor"].(string), "%", "", -1)
					maxScaleStr := strings.Replace(defaultObj.(map[string]interface{})["MaxScaleFactor"].(string), "%", "", -1)
					ramStr := strings.Replace(strings.Split(d.Get("db_instance_class").(string), "|")[0], "db.ram.", "", -1)
					ram, _ := strconv.ParseInt(ramStr, 10, 64)
					maxScale, _ := strconv.ParseFloat(maxScaleStr, 64)
					scale, _ := strconv.ParseFloat(scaleStr, 64)
					maxScaleValue := int64(float64(ram*1024*1024*1024) * (maxScale / 100))
					defaultScaleValue := int64(float64(ram*1024*1024*1024) * (scale / 100))
					valueNum, err := strconv.ParseInt(kv[key], 10, 64)
					if err != nil {
						return needRestart, num, err
					}
					if toDefault && valueNum == defaultScaleValue {
						continue
					}
					if !toDefault && len(currents) > 0 && valueNum == int64(currents[key].(float64)) {
						continue
					}
					if valueNum < dMin || valueNum > maxScaleValue {
						return needRestart, num, fmt.Errorf("parameter %s must in [%d , %d]", key, dMin, maxScaleValue)
					}
					dv = defaultScaleValue
				} else {
					return needRestart, num, fmt.Errorf("parameter %s must not support in terraform-provider-ksyun now", key)
				}
			default:
				if toDefault && dv.(string) == kv[key] {
					continue
				}
				if !toDefault && len(currents) > 0 && kv[key] == currents[key].(string) {
					continue
				}
				dE := defaultObj.(map[string]interface{})["Enums"].([]interface{})
				valid := false
				for _, s := range dE {
					if s.(string) == kv[key] {
						valid = true
					}
				}
				if !valid {
					return needRestart, num, fmt.Errorf("parameter %s must in %s", key, dE)
				}
			}
			(*req)["Parameters.Name."+strconv.Itoa(num)] = key
			if toDefault {
				(*req)["Parameters.Value."+strconv.Itoa(num)] = dv
			} else {
				(*req)["Parameters.Value."+strconv.Itoa(num)] = kv[key]
			}

			num = num + 1
			if dr {
				needRestart = true
			}
		} else {
			return needRestart, num, fmt.Errorf("parameter %s not support", key)
		}
	}
	return needRestart, num, err
}

func createPostgresqlRrParameterGroup(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var (
		data map[string]interface{}
	)
	// read master
	data, err = readPostgresqlInstance(d, meta, d.Get("db_instance_identifier").(string))
	if err != nil {
		return call, err
	}
	err = d.Set("engine", data["Engine"])
	if err != nil {
		return call, err
	}
	err = d.Set("engine_version", data["EngineVersion"])
	if err != nil {
		return call, err
	}
	return modifyPostgresqlParameterGroup(d, meta, true)
}

func createPostgresqlTempParameterGroup(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	paramsReq, _, err := checkAndProcessPostgresqlParameters(d, meta)
	if _, ok := d.GetOk("db_parameter_template_id"); ok {
		// merge parameters
		paramsReq, _, err = DbGroupParameterMergeAndCheckProcessPostgresql(d, meta)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return call, err
	}
	if len(paramsReq) > 0 {
		paramsReq["DBParameterGroupName"] = d.Get("db_instance_name").(string) + "_param"
		paramsReq["Description"] = d.Get("db_instance_name").(string) + "_desc"
		paramsReq["Engine"] = d.Get("engine")
		paramsReq["EngineVersion"] = d.Get("engine_version")

		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			action := "CreateDBParameterGroup"
			logger.Debug(logger.RespFormat, action, paramsReq)
			paramResp, err := conn.CreateDBParameterGroup(&paramsReq)
			logger.Debug(logger.AllFormat, action, paramsReq, paramResp, err)
			if err != nil {
				return fmt.Errorf("error on create Instance(postgresql) DBParameterGroup : %s", err)
			}
			parameterId, err := getSdkValue("Data.DBParameterGroup.DBParameterGroupId", *paramResp)
			if err != nil {
				return fmt.Errorf("error on create Instance(postgresql) DBParameterGroup : %s", err)
			}
			return d.Set("db_parameter_group_id", parameterId)
		}
	}
	return call, err
}
func removePostgresqlParameterGroup(d *schema.ResourceData, meta interface{}) (err error) {
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		conn := meta.(*KsyunClient).postgresqlconn
		if d.Get("db_parameter_group_id") != nil && d.Get("db_parameter_group_id").(string) != "" {
			delParam := make(map[string]interface{})
			delParam["DBParameterGroupId"] = d.Get("db_parameter_group_id").(string)
			_, deleteErr := conn.DeleteDBParameterGroup(&delParam)
			// logger.Debug("test %s %s %s", "DeleteDBParameterGroup", inUseError(deleteErr), deleteErr)
			if deleteErr == nil || notFoundErrorNew(deleteErr) || inUseError(deleteErr) {
				return nil
			} else {
				return resource.RetryableError(deleteErr)
			}
		}
		// 没有参数组，不执行删除
		return nil // resource.RetryableError(nil)
	})
}

func modifyPostgresqlParameterGroup(d *schema.ResourceData, meta interface{}, onCreate bool) (call ksyunApiCallFunc, err error) {
	paramsReq, restart, err := checkAndProcessPostgresqlParameters(d, meta)
	// logger.Debug("test", "test", paramsReq, restart, err)
	if err != nil {

		return call, err
	}
	if len(paramsReq) > 0 {
		// check force_restart
		if !d.Get("force_restart").(bool) && restart && !onCreate {
			return call, fmt.Errorf("some parameter change must set force_restart true")
		}
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			paramsReq["DBParameterGroupId"] = d.Get("db_parameter_group_id").(string)
			conn := meta.(*KsyunClient).postgresqlconn
			action := "ModifyDBParameterGroup"
			logger.Debug(logger.RespFormat, action, paramsReq)
			_, err = conn.ModifyDBParameterGroup(&paramsReq)
			if err != nil {
				return err
			}
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			if d.Get("force_restart").(bool) || (onCreate && restart) {
				restartParam := make(map[string]interface{})
				restartParam["DBInstanceIdentifier"] = d.Id()
				action = "RebootDBInstance"
				logger.Debug(logger.RespFormat, action, restartParam)
				_, err = conn.RebootDBInstance(&restartParam)
				if err != nil {
					return err
				}
			}
			return err
		}
	}
	return call, err
}

func createPostgresqlInstance(d *schema.ResourceData, meta interface{}, isRR bool) (err error) {
	var (
		api []ksyunApiCallFunc
	)
	if !isRR {
		// template parameter
		tempParameterGroupCall, err := createPostgresqlTempParameterGroup(d, meta)
		if err != nil {
			return err
		}
		api = append(api, tempParameterGroupCall)
		// create instance
		instanceCreateCall, err := createPostgresqlDbInstance(d, meta)
		if err != nil {
			return err
		}
		api = append(api, instanceCreateCall)
	} else {
		// create rr instance
		instanceCreateCall, err := createPostgresqlRrInstance(d, meta)
		if err != nil {
			return err
		}
		api = append(api, instanceCreateCall)

		tempParameterGroupCall, err := createPostgresqlRrParameterGroup(d, meta)
		if err != nil {
			return err
		}
		api = append(api, tempParameterGroupCall)
		// ModifySG
		modifyDBSg, err := modifyPostgresqlInstanceSg(d, meta, false)
		if err != nil {
			return err
		}
		api = append(api, modifyDBSg)
	}

	// process eip
	eipCall, err := allocateOrReleasePostgresqlInstanceEip(d, meta)
	if err != nil {
		return err
	}
	api = append(api, eipCall)
	// api call
	err = ksyunApiCall(api, d, meta)
	if err != nil {
		return err
	}
	return err
}

// parsePostgresqlCreateResp 从 CreateDBInstance / CreateDBInstanceReadReplica 响应中
// 解析实例 ID 与参数组 ID。postgresql 响应结构与 krds 不同 (Create 响应可能为
// Data.DBInstance.* / Data.Instances.0.* / 顶层等不同形态), 故递归遍历 resp 任意层级
// (map / slice) 查找第一个非空 string 值的 DBInstanceIdentifier / DBParameterGroupId,
// 兼容各种嵌套结构, 避免类型断言 panic。
func parsePostgresqlCreateResp(resp map[string]interface{}) (instanceId, dbParameterGroupId string, ok bool) {
	if id, found := findPostgresqlRespField(resp, "DBInstanceIdentifier"); found {
		instanceId = id
		ok = true
	}
	if pgId, found := findPostgresqlRespField(resp, "DBParameterGroupId"); found {
		dbParameterGroupId = pgId
	}
	return
}

// findPostgresqlRespField 递归遍历 resp 任意层级 (map / slice) 查找第一个非空 string 值的指定字段。
func findPostgresqlRespField(v interface{}, field string) (string, bool) {
	switch t := v.(type) {
	case map[string]interface{}:
		if val, exist := t[field]; exist {
			if s, ok2 := val.(string); ok2 && s != "" {
				return s, true
			}
		}
		for _, val := range t {
			if s, found := findPostgresqlRespField(val, field); found {
				return s, true
			}
		}
	case []interface{}:
		for _, val := range t {
			if s, found := findPostgresqlRespField(val, field); found {
				return s, true
			}
		}
	}
	return "", false
}

func createPostgresqlDbInstance(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	transform := map[string]SdkReqTransform{
		"db_instance_class":     {mapping: "DBInstanceClass"},
		"db_instance_name":      {mapping: "DBInstanceName"},
		"db_instance_type":      {mapping: "DBInstanceType"},
		"db_parameter_group_id": {mapping: "DBParameterGroupId"},
		"instance_has_eip":      {Ignore: true},
		"parameters":            {Ignore: true},
		"force_restart":         {Ignore: true},
		"availability_zone_1":   {mapping: "AvailabilityZone.1"},
		"availability_zone_2":   {mapping: "AvailabilityZone.2"},
		"vcpus":                 {mapping: "Vcpus"},
	}

	createReq, err := SdkRequestAutoMapping(d, resourceKsyunPostgresql(), false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return call, err
	}
	call = func(d *schema.ResourceData, meta interface{}) (err error) {
		conn := meta.(*KsyunClient).postgresqlconn
		action := "CreateDBInstance"

		// 如果创建了临时参数组，创建实例的时候使用该参数组
		if d.Get("db_parameter_group_id") != nil && d.Get("db_parameter_group_id").(string) != "" {
			createReq["DBParameterGroupId"] = d.Get("db_parameter_group_id")
		}
		logger.Debug(logger.RespFormat, action, createReq)
		resp, err := conn.CreateDBInstance(&createReq)
		if err != nil {

			// 由于临时参数组不被tf管理，创建实例失败，需要手动回收
			if d.Get("db_parameter_group_id") != nil && d.Get("db_parameter_group_id").(string) != "" {
				removePostgresqlParameterGroup(d, meta)
			}

			return err
		}
		logger.Debug(logger.AllFormat, action, createReq, *resp, err)
		if resp != nil {
			instanceId, _, ok := parsePostgresqlCreateResp(*resp)
			if !ok {
				return fmt.Errorf("create postgresql instance succeeded but failed to parse instance id from response, resp: %+v", *resp)
			}
			d.SetId(instanceId)
		}
		err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		err = removePostgresqlParameterGroup(d, meta)
		if err != nil {
			return err
		}
		return err
	}

	return call, err
}

func createPostgresqlRrInstance(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	transform := map[string]SdkReqTransform{
		"db_instance_identifier": {mapping: "DBInstanceIdentifier"},
		"db_instance_class":      {mapping: "DBInstanceClass"},
		"db_instance_name":       {mapping: "DBInstanceName"},
		"availability_zone_1":    {mapping: "AvailabilityZone.1"},
		"db_parameter_group_id":  {Ignore: true},
		"security_group_id":      {Ignore: true},
		"instance_has_eip":       {Ignore: true},
		"parameters":             {Ignore: true},
		"force_restart":          {Ignore: true},
		"vcpus":                  {mapping: "Vcpus"},
	}

	createReq, err := SdkRequestAutoMapping(d, resourceKsyunPostgresqlRr(), false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return call, err
	}

	config := strings.Split(createReq["DBInstanceClass"].(string), "|")
	createReq["Mem"] = strings.Replace(config[0], "db.ram.", "", -1)
	createReq["Disk"] = strings.Replace(config[1], "db.disk.", "", -1)
	delete(createReq, "DBInstanceClass")

	call = func(d *schema.ResourceData, meta interface{}) (err error) {
		conn := meta.(*KsyunClient).postgresqlconn
		action := "CreateDBInstanceReadReplica"
		logger.Debug(logger.RespFormat, action, createReq)
		resp, err := conn.CreateDBInstanceReadReplica(&createReq)
		if err != nil {
			return err
		}
		logger.Debug(logger.AllFormat, action, createReq, *resp, err)
		if resp != nil {
			instanceId, dbParameterGroupId, ok := parsePostgresqlCreateResp(*resp)
			if !ok {
				return fmt.Errorf("create postgresql rr instance succeeded but failed to parse instance id from response, resp: %+v", *resp)
			}
			d.SetId(instanceId)
			if dbParameterGroupId != "" {
				err = d.Set("db_parameter_group_id", dbParameterGroupId)
				if err != nil {
					return err
				}
			}
		}
		err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		return err
	}

	return call, err
}

func removePostgresqlInstance(d *schema.ResourceData, meta interface{}) (err error) {
	err = removePostgresqlDbInstance(d, meta)
	if err != nil {
		return err
	}
	return removePostgresqlParameterGroup(d, meta)
}

func removePostgresqlDbInstance(d *schema.ResourceData, meta interface{}) (err error) {
	conn := meta.(*KsyunClient).postgresqlconn
	req := make(map[string]interface{})
	req["DBInstanceIdentifier"] = d.Id()
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		action := "DeleteDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.DeleteDBInstance(&req)
		if err == nil {
			return nil
		}
		action = "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.DescribeDBInstances(&req)
		if err != nil {
			if notFoundError(err) {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading postgresql when delete %q, %s", d.Id(), err))
			}
		}
		return nil
	})
}

func modifyPostgresqlInstance(d *schema.ResourceData, meta interface{}, isRR bool) (err error) {
	var (
		call []ksyunApiCallFunc
	)
	// ModifyDBInstance
	modifyDBCall, err := modifyPostgresqlInstanceCommon(d, meta)
	if err != nil {
		return err
	}
	call = append(call, modifyDBCall)
	// ModifySG
	modifyDBSg, err := modifyPostgresqlInstanceSg(d, meta, true)
	if err != nil {
		return err
	}
	call = append(call, modifyDBSg)
	var (
		upgradeDBInstanceEngineVersionCall ksyunApiCallFunc
	)
	// ModifyDBInstanceSpec
	modifyDBInstanceSpecCall, err := modifyPostgresqlInstanceSpec(d, meta)
	if err != nil {
		return err
	}
	call = append(call, modifyDBInstanceSpecCall)
	if !isRR {
		// UpgradeDBInstanceEngineVersion
		upgradeDBInstanceEngineVersionCall, err = modifyPostgresqlInstanceEngineVersion(d, meta)
		if err != nil {
			return err
		}
		call = append(call, upgradeDBInstanceEngineVersionCall)
		// ModifyDBInstanceAvailabilityZone
		modifyDBInstanceAvailabilityZoneCall, err := modifyPostgresqlInstanceAvailabilityZone(d, meta)
		if err != nil {
			return err
		}
		call = append(call, modifyDBInstanceAvailabilityZoneCall)
	}

	// ModifyDBInstanceEip
	modifyDBInstanceEipCall, err := allocateOrReleasePostgresqlInstanceEip(d, meta)
	if err != nil {
		return err
	}
	call = append(call, modifyDBInstanceEipCall)
	// ModifyParameters
	if upgradeDBInstanceEngineVersionCall == nil {
		var (
			modifyParametersCall ksyunApiCallFunc
		)
		modifyParametersCall, err = modifyPostgresqlParameterGroup(d, meta, false)
		if err != nil {
			return err
		}
		call = append(call, modifyParametersCall)
	}
	err = ksyunApiCall(call, d, meta)
	if err != nil {
		return err
	}
	return err
}

func modifyPostgresqlInstanceCommon(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	oldType, _ := d.GetChange("db_instance_type")
	var modifyInstanceParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"db_instance_name":      {mapping: "DBInstanceName"},
		"master_user_password":  {},
		"preferred_backup_time": {},
		"project_id":            {},
	}
	modifyInstanceParam, err = SdkRequestAutoMapping(d, resourceKsyunPostgresql(), true, transform, nil)
	if err != nil {
		return call, err
	}
	if _, ok := modifyInstanceParam["PreferredBackupTime"]; ok && oldType == "RR_PG" {
		return call, fmt.Errorf("postgresql rr is not support update %s", "preferred_backup_time")
	}
	call = func(d *schema.ResourceData, meta interface{}) (err error) {
		conn := meta.(*KsyunClient).postgresqlconn
		// modify project
		err = ModifyProjectInstance(d.Id(), &modifyInstanceParam, meta)
		if err != nil {
			return err
		}
		if len(modifyInstanceParam) > 0 {
			modifyInstanceParam["DBInstanceIdentifier"] = d.Id()
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			action := "ModifyDBInstance"
			logger.Debug(logger.ReqFormat, action, modifyInstanceParam)
			_, err = conn.ModifyDBInstance(&modifyInstanceParam)
			if err != nil {
				return err
			}
		}
		return err
	}

	return call, err
}

func modifyPostgresqlInstanceSg(d *schema.ResourceData, meta interface{}, isUpdate bool) (call ksyunApiCallFunc, err error) {
	var modifyInstanceSgParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"security_group_id": {},
	}
	modifyInstanceSgParam, err = SdkRequestAutoMapping(d, resourceKsyunPostgresql(), isUpdate, transform, nil)
	if err != nil {
		return call, err
	}
	if len(modifyInstanceSgParam) > 0 {
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			modifyInstanceSgParam["DBInstanceIdentifier"] = d.Id()
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			action := "ModifyDBInstance"
			logger.Debug(logger.ReqFormat, action, modifyInstanceSgParam)
			_, err = conn.ModifyDBInstance(&modifyInstanceSgParam)
			if err != nil {
				return err
			}
			return err
		}
	}
	return call, err
}

func modifyPostgresqlInstanceSpec(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var modifyDBInstanceSpecParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"db_instance_class": {mapping: "DBInstanceClass"},
	}
	modifyDBInstanceSpecParam, err = SdkRequestAutoMapping(d, resourceKsyunPostgresql(), true, transform, nil)
	if err != nil {
		return call, err
	}
	if len(modifyDBInstanceSpecParam) > 0 {
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			modifyDBInstanceSpecParam["DBInstanceIdentifier"] = d.Id()
			action := "ModifyDBInstanceSpec"
			logger.Debug(logger.ReqFormat, action, modifyDBInstanceSpecParam)
			_, err = conn.ModifyDBInstanceSpec(&modifyDBInstanceSpecParam)
			logger.Debug(logger.AllFormat, action, modifyDBInstanceSpecParam, err)
			if err != nil {
				return err
			}
			return err
		}
	}
	return call, err
}

func modifyPostgresqlInstanceEngineVersion(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var upgradeDBInstanceEngineVersionParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"engine":         {},
		"engine_version": {},
	}
	upgradeDBInstanceEngineVersionParam, err = SdkRequestAutoMapping(d, resourceKsyunPostgresql(), true, transform, nil)
	if err != nil {
		return call, err
	}
	if len(upgradeDBInstanceEngineVersionParam) > 0 {
		_, err = modifyPostgresqlParameterGroup(d, meta, false)
		if err != nil {
			return call, err
		}
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			upgradeDBInstanceEngineVersionParam["DBInstanceIdentifier"] = d.Id()
			if _, ok := upgradeDBInstanceEngineVersionParam["Engine"]; !ok {
				upgradeDBInstanceEngineVersionParam["Engine"] = d.Get("engine")
			}
			if _, ok := upgradeDBInstanceEngineVersionParam["EngineVersion"]; !ok {
				upgradeDBInstanceEngineVersionParam["EngineVersion"] = d.Get("engine_version")
			}
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}

			action := "UpdateDBInstanceVersion"
			logger.Debug(logger.ReqFormat, action, upgradeDBInstanceEngineVersionParam)
			_, err = conn.UpdateDBInstanceVersion(&upgradeDBInstanceEngineVersionParam)
			logger.Debug(logger.AllFormat, action, upgradeDBInstanceEngineVersionParam, err)
			if err != nil {
				return err
			}
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			// clean old db_parameter_group_id
			err = removePostgresqlParameterGroup(d, meta)
			if err != nil {
				return err
			}
			// query db
			req := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
			action = "DescribeDBInstances"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err := conn.DescribeDBInstances(&req)
			logger.Debug(logger.AllFormat, action, req, resp, err)
			if err != nil {
				return err
			}
			value, err := getSdkValue("Data.Instances.0.DBParameterGroupId", *resp)
			if err != nil {
				return err
			}
			return d.Set("db_parameter_group_id", value)
		}
	}
	return call, err
}

func modifyPostgresqlInstanceAvailabilityZone(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var modifyDBInstanceAvailabilityZoneParam map[string]interface{}
	oldType, _ := d.GetChange("db_instance_type")
	transform := map[string]SdkReqTransform{
		"availability_zone_1": {mapping: "AvailabilityZone.1"},
		"availability_zone_2": {mapping: "AvailabilityZone.2"},
	}
	modifyDBInstanceAvailabilityZoneParam, err = SdkRequestAutoMapping(d, resourceKsyunPostgresql(), true, transform, nil)
	if err != nil {
		return call, err
	}
	if len(modifyDBInstanceAvailabilityZoneParam) > 0 {
		if oldType == "RR_PG" {
			return call, fmt.Errorf("postgresql rr is not support update %s", "availability_zone")
		}
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			modifyDBInstanceAvailabilityZoneParam["DBInstanceIdentifier"] = d.Id()
			action := "ModifyDBInstanceAvailabilityZone"
			logger.Debug(logger.ReqFormat, action, modifyDBInstanceAvailabilityZoneParam)
			_, err = conn.ModifyDBInstanceAvailabilityZone(&modifyDBInstanceAvailabilityZoneParam)
			if err != nil {
				return err
			}
			return err
		}
	}
	return call, err
}

func allocateOrReleasePostgresqlInstanceEip(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	if _, ok := d.GetOk("instance_has_eip"); ok || (d.HasChange("instance_has_eip") && !d.IsNewResource()) {
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			err = checkPostgresqlInstanceState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
			req := map[string]interface{}{
				"DBInstanceIdentifier": d.Id(),
			}

			if d.Get("instance_has_eip") == true {
				action := "AllocateDBInstanceEip"
				logger.Debug(logger.ReqFormat, action, req)
				resp, err := conn.AllocateDBInstanceEip(&req)
				logger.Debug(logger.AllFormat, action, req, resp, err)

				if err != nil {
					return err
				}
			} else if d.Get("instance_has_eip") == false && !d.IsNewResource() {
				action := "ReleaseDBInstanceEip"
				logger.Debug(logger.ReqFormat, action, req)
				resp, err := conn.ReleaseDBInstanceEip(&req)
				logger.Debug(logger.AllFormat, action, req, resp, err)

				if err != nil {
					return err
				}
			}
			return err
		}
	}

	return call, err
}

func checkPostgresqlInstanceState(d *schema.ResourceData, meta interface{}, instanceId string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"ACTIVE"},
		Refresh:    postgresqlInstanceStateRefreshFunc(d, meta, instanceId, []string{"error"}),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func postgresqlInstanceStateRefreshFunc(d *schema.ResourceData, meta interface{}, instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := readPostgresqlInstance(d, meta, instanceId)
		if err != nil {
			return nil, "", err
		}

		status := data["DBInstanceStatus"].(string)

		for _, v := range failStates {
			if v == status {
				return nil, "", fmt.Errorf("instance status  error, status:%v", status)
			}
		}
		return data, status, nil
	}
}

func postgresqlInstanceCustomizeDiff() schema.CustomizeDiffFunc {
	return func(diff *schema.ResourceDiff, i interface{}) (err error) {
		if diff.HasChange("parameters") {
			var (
				data map[string]interface{}
			)
			_, n := diff.GetChange("parameters")
			data, err = readPostgresqlDefaultParameters(nil, diff, i)
			if err != nil {
				return err
			}
			for _, v := range n.(*schema.Set).List() {
				key := v.(map[string]interface{})["name"]
				exist := false
				for k := range data {
					if k == key.(string) {
						exist = true
						break
					}
				}
				if !exist {
					return fmt.Errorf("parameter %s is not support", key)
				}
			}
		}
		return err
	}
}

func createPostgresqlSecurityGroupRule(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		req map[string]interface{}
	)
	transform := map[string]SdkReqTransform{
		"security_group_id":            {},
		"security_group_rule_protocol": {mapping: "SecurityGroupRule.SecurityGroupRuleProtocol.1"},
		"security_group_rule_name":     {mapping: "SecurityGroupRule.SecurityGroupRuleName.1"},
	}
	req, err = SdkRequestAutoMapping(d, resourceKsyunPostgresqlSecurityGroupRule(), false, transform, nil)
	if err != nil {
		return err
	}
	conn := meta.(*KsyunClient).postgresqlconn
	req["SecurityGroupRuleAction"] = "Attach"
	action := "ModifySecurityGroupRule"
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.ModifySecurityGroupRule(&req)
	if err != nil {
		return err
	}
	d.SetId(d.Get("security_group_id").(string) + ":" + d.Get("security_group_rule_protocol").(string))
	return err
}

func readAndSetPostgresqlSecurityGroupRule(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		rules map[string]interface{}
	)
	protocol := d.Get("security_group_rule_protocol").(string)
	sgId := d.Get("security_group_id").(string)
	rules, err = readPostgresqlSecurityGroupRules(d, meta, sgId)
	if err != nil {
		return err
	}
	for k, v := range rules {
		if k == protocol {
			SdkResponseAutoResourceData(d, resourceKsyunPostgresqlSecurityGroupRule(), v, nil)
			return err
		}
	}
	d.SetId("")
	return nil
}

func readPostgresqlSecurityGroupRules(d *schema.ResourceData, meta interface{}, sgId string) (data map[string]interface{}, err error) {
	var (
		sg    map[string]interface{}
		rules interface{}
	)
	if sgId == "" {
		sgId = d.Id()
	}
	sg, err = readPostgresqlSecurityGroup(d, meta, sgId)
	if err != nil {
		return data, err
	}
	if len(sg) == 0 {
		return data, nil
	}
	data = make(map[string]interface{})
	rules, err = getSdkValue("SecurityGroupRules", sg)
	if err != nil {
		return data, err
	}
	for _, i := range rules.([]interface{}) {
		protocol := i.(map[string]interface{})["SecurityGroupRuleProtocol"].(string)
		data[protocol] = i
	}
	return data, err
}

func removePostgresqlSecurityGroupRule(d *schema.ResourceData, meta interface{}) (err error) {
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		conn := meta.(*KsyunClient).postgresqlconn
		req := map[string]interface{}{
			"SecurityGroupId":                         d.Get("security_group_id"),
			"SecurityGroupRuleAction":                 "Delete",
			"SecurityGroupRule.SecurityGroupRuleId.1": d.Get("security_group_rule_id"),
		}
		action := "ModifySecurityGroupRule"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifySecurityGroupRule(&req)
		if err == nil || notFoundErrorNew(err) {
			return nil
		} else {
			return resource.RetryableError(err)
		}
	})
}

func importPostgresqlSecurityGroupRule() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
			items := strings.Split(d.Id(), ":")
			if len(items) != 2 {
				return nil, fmt.Errorf("id must split with %s and size %v", ":", 2)
			}
			err = d.Set("security_group_id", items[0])
			if err != nil {
				return nil, err
			}
			err = d.Set("security_group_rule_protocol", items[1])
			if err != nil {
				return nil, err
			}
			data = append(data, d)
			return data, err
		},
	}
}

func readPostgresqlAndSetSecurityGroup(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		sg map[string]interface{}
	)
	sg, err = readPostgresqlSecurityGroup(d, meta, "")
	if err != nil {
		return err
	}
	if len(sg) == 0 {
		d.SetId("")
		return nil
	}
	extra := map[string]SdkResponseMapping{
		"SecurityGroupRules": {
			Field: "security_group_rule",
			FieldRespFunc: func(i interface{}) interface{} {
				var (
					result []interface{}
				)
				for _, v := range i.([]interface{}) {
					r := map[string]interface{}{
						"security_group_rule_id":       v.(map[string]interface{})["SecurityGroupRuleId"],
						"created":                      v.(map[string]interface{})["Created"],
						"security_group_rule_protocol": v.(map[string]interface{})["SecurityGroupRuleProtocol"],
					}
					if o, ok := v.(map[string]interface{})["SecurityGroupRuleName"]; ok {
						r["security_group_rule_name"] = o
					}
					result = append(result, r)
				}
				return result
			},
		},
	}
	SdkResponseAutoResourceData(d, resourceKsyunPostgresqlSecurityGroup(), sg, extra)
	return err
}

func readPostgresqlSecurityGroup(d *schema.ResourceData, meta interface{}, sgId string) (data map[string]interface{}, err error) {
	var (
		req  map[string]interface{}
		resp *map[string]interface{}
		sg   interface{}
	)
	req = make(map[string]interface{})
	if sgId == "" {
		sgId = d.Id()
	}
	conn := meta.(*KsyunClient).postgresqlconn
	req["SecurityGroupId"] = sgId
	action := "DescribeSecurityGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.DescribeSecurityGroup(&req)
	if err != nil {
		if notFoundError(err) {
			return data, nil
		}
		return data, err
	}
	sg, err = getSdkValue("Data.SecurityGroups.0", *resp)
	if err != nil {
		return data, nil
	}
	data = sg.(map[string]interface{})
	return data, err
}

func modifyPostgresqlSecurityGroup(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		call []ksyunApiCallFunc
	)
	commonCall, err := modifyPostgresqlSecurityGroupCommon(d, meta)
	if err != nil {
		return err
	}
	call = append(call, commonCall)
	rulesCall, err := modifyPostgresqlSecurityGroupRules(d, meta)
	if err != nil {
		return err
	}
	call = append(call, rulesCall)
	return ksyunApiCall(call, d, meta)
}

func modifyPostgresqlSecurityGroupCommon(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var (
		req map[string]interface{}
	)
	transform := map[string]SdkReqTransform{
		"security_group_name":        {},
		"security_group_description": {},
	}
	req, err = SdkRequestAutoMapping(d, resourceKsyunPostgresqlSecurityGroup(), true, transform, nil)
	if err != nil {
		return call, err
	}
	if len(req) > 0 {
		req["SecurityGroupId"] = d.Id()
		call = func(d *schema.ResourceData, meta interface{}) (err error) {
			conn := meta.(*KsyunClient).postgresqlconn
			action := "ModifySecurityGroup"
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.ModifySecurityGroup(&req)
			return err
		}
	}
	return call, err
}

func checkAndProcessPostgresqlSgRules(d *schema.ResourceData, meta interface{}) (reqAttach map[string]interface{}, reqRemove map[string]interface{}, err error) {
	reqAttach = make(map[string]interface{})
	reqRemove = make(map[string]interface{})
	if d.HasChange("security_group_rule") {
		var (
			oldRules     interface{}
			newRules     interface{}
			needToRemove []interface{}
			needToAttach []interface{}
			current      map[string]interface{}
		)
		oldKv := make(map[string]interface{})
		newKv := make(map[string]interface{})
		if d.Id() == "" {
			current = make(map[string]interface{})
		} else {
			current, err = readPostgresqlSecurityGroupRules(d, meta, "")
			if err != nil {
				return reqAttach, reqRemove, err
			}
		}

		oldRules, newRules = d.GetChange("security_group_rule")
		// init from list to kv
		if params, ok := oldRules.(*schema.Set); ok {
			for _, i := range params.List() {
				oldKv[i.(map[string]interface{})["security_group_rule_protocol"].(string)] = i
			}
		}
		if params, ok := newRules.(*schema.Set); ok {
			for _, i := range params.List() {
				newKv[i.(map[string]interface{})["security_group_rule_protocol"].(string)] = i
			}
		}
		// compare add or remove
		for k := range oldKv {
			if _, ok := newKv[k]; !ok {
				needToRemove = append(needToRemove, k)
			}
		}
		for k, n := range newKv {
			if o, ok := oldKv[k]; ok {
				newName := n.(map[string]interface{})["security_group_rule_name"]
				oldName := o.(map[string]interface{})["security_group_rule_name"]
				if newName != "" && newName != oldName {
					needToRemove = append(needToRemove, k)
					needToAttach = append(needToAttach, k)
				}
			} else {
				if c, ok1 := current[k]; ok1 {
					logger.Debug(logger.ReqFormat, "Demo", "1111")
					newName := n.(map[string]interface{})["security_group_rule_name"]
					currentName := c.(map[string]interface{})["SecurityGroupRuleName"]
					if newName != "" && newName != currentName {
						needToRemove = append(needToRemove, k)
						needToAttach = append(needToAttach, k)
					}
				} else {
					needToAttach = append(needToAttach, k)
				}
			}

		}
		index := 1
		for _, v := range needToAttach {
			if newKv[v.(string)].(map[string]interface{})["security_group_rule_name"].(string) != "" {
				reqAttach["SecurityGroupRule.SecurityGroupRuleName."+strconv.Itoa(index)] = newKv[v.(string)].(map[string]interface{})["security_group_rule_name"]
			} else {
				reqAttach["SecurityGroupRule.SecurityGroupRuleName."+strconv.Itoa(index)] = ""
			}
			reqAttach["SecurityGroupRule.SecurityGroupRuleProtocol."+strconv.Itoa(index)] = newKv[v.(string)].(map[string]interface{})["security_group_rule_protocol"]
			index = index + 1
		}
		index = 1
		for _, v := range needToRemove {
			if _, ok := oldKv[v.(string)]; ok {
				reqRemove["SecurityGroupRule.SecurityGroupRuleId."+strconv.Itoa(index)] = oldKv[v.(string)].(map[string]interface{})["security_group_rule_id"]
				index = index + 1
			} else if _, ok1 := current[v.(string)]; ok1 {
				reqRemove["SecurityGroupRule.SecurityGroupRuleId."+strconv.Itoa(index)] = current[v.(string)].(map[string]interface{})["SecurityGroupRuleId"]
				index = index + 1
			}

		}

	}
	return reqAttach, reqRemove, err
}

func modifyPostgresqlSecurityGroupRules(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var (
		reqAttach map[string]interface{}
		reqRemove map[string]interface{}
	)
	reqAttach, reqRemove, err = checkAndProcessPostgresqlSgRules(d, meta)
	if err != nil {
		return call, err
	}
	call = func(d *schema.ResourceData, meta interface{}) (err error) {
		conn := meta.(*KsyunClient).postgresqlconn
		logger.Debug(logger.ReqFormat, "Demo", reqRemove)
		logger.Debug(logger.ReqFormat, "Demo", reqAttach)
		if len(reqRemove) > 0 {
			reqRemove["SecurityGroupRuleAction"] = "Delete"
			reqRemove["SecurityGroupId"] = d.Id()
			action := "ModifySecurityGroupRule-Delete"
			logger.Debug(logger.ReqFormat, action, reqRemove)
			_, err = conn.ModifySecurityGroupRule(&reqRemove)
			if err != nil {
				return err
			}
		}
		if len(reqAttach) > 0 {
			reqAttach["SecurityGroupRuleAction"] = "Attach"
			reqAttach["SecurityGroupId"] = d.Id()
			action := "ModifySecurityGroupRule-Attach"
			logger.Debug(logger.ReqFormat, action, reqAttach)
			_, err = conn.ModifySecurityGroupRule(&reqAttach)
			if err != nil {
				return err
			}
		}
		return err
	}
	return call, err
}

func removePostgresqlSecurityGroup(d *schema.ResourceData, meta interface{}) (err error) {
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		conn := meta.(*KsyunClient).postgresqlconn
		req := map[string]interface{}{
			"SecurityGroupId.1": d.Id(),
		}
		action := "DeleteSecurityGroup"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.DeleteSecurityGroup(&req)
		if err == nil || notFoundErrorNew(err) {
			return nil
		} else {
			return resource.RetryableError(err)
		}
	})
}

func createPostgresqlSecurityGroup(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		call []ksyunApiCallFunc
	)
	sgCall, err := createPostgresqlSecurityGroupCommon(d, meta)
	if err != nil {
		return err
	}
	call = append(call, sgCall)
	rulesCall, err := modifyPostgresqlSecurityGroupRules(d, meta)
	if err != nil {
		return err
	}
	call = append(call, rulesCall)
	return ksyunApiCall(call, d, meta)
}

func createPostgresqlSecurityGroupCommon(d *schema.ResourceData, meta interface{}) (call ksyunApiCallFunc, err error) {
	var (
		req  map[string]interface{}
		resp *map[string]interface{}
	)
	transform := map[string]SdkReqTransform{
		"security_group_description": {},
		"security_group_name":        {},
	}
	req, err = SdkRequestAutoMapping(d, resourceKsyunPostgresqlSecurityGroup(), false, transform, nil)
	if err != nil {
		return call, err
	}
	call = func(d *schema.ResourceData, meta interface{}) (err error) {
		conn := meta.(*KsyunClient).postgresqlconn
		action := "CreateSecurityGroup"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.CreateSecurityGroup(&req)
		if err != nil {
			return err
		}
		r, err := getSdkValue("Data.SecurityGroups.0.SecurityGroupId", *resp)
		if err != nil {
			return err
		}
		d.SetId(r.(string))
		return err
	}
	return call, err
}

func DbGroupParameterMergeAndCheckProcessPostgresql(d *schema.ResourceData,
	meta interface{},
) (map[string]interface{},
	bool,
	error) {
	var (
		templateId string

		templateParams map[string]interface{}
		err            error
	)

	templateId, _ = If2String(d.Get("db_parameter_template_id"))
	if templateId != "" {
		templateParams, err = readPostgresqlParameterGroup(d, meta, templateId, false)
		if err != nil {
			return nil, false, err
		}
		templateParams, err = parameterGroupTransformerPostgresql(d, meta, templateParams)
		if err != nil {
			return nil, false, err
		}
	} else {
		templateParams = make(map[string]interface{})
	}

	parameters := d.Get("parameters")
	unmergeParams, ok := parameters.(*schema.Set)
	if !ok {
		return nil, false, fmt.Errorf("parameters structure is invalid")
	}

	unMergeParamConverted := TfParametersConvert2Map(unmergeParams)

	for k, v := range unMergeParamConverted {
		templateParams[k] = v
	}

	kv := make(map[string]string)
	keys := make([]string, 0, len(templateParams))
	for k, v := range templateParams {
		kv[k], _ = If2String(v)
		keys = append(keys, k)
	}

	req := make(map[string]interface{})
	index := 1
	needRestart, index, err := prepareModifyDbParameterParamsPostgresql(d, meta, keys, &req, kv, index, false)
	if err != nil {
		return nil, false, err
	}

	return req, needRestart, nil
}
