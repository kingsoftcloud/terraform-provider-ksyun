package ksyun

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func readMongodbSupportRegions(meta interface{}) (regions []interface{}, err error) {
	var (
		resp *map[string]interface{}
		data interface{}
	)
	conn := meta.(*KsyunClient).mongodbconn
	resp, err = conn.DescribeRegions(nil)
	if err != nil {
		return regions, err
	}
	data, err = getSdkValue("Data.Regions", *resp)
	if err != nil {
		return regions, err
	}
	regions = data.([]interface{})
	return regions, err
}

func readMongodbSupportAzMappings(meta interface{}) (mappings map[string]string, err error) {
	var (
		current string
		exist   bool
		regions []interface{}
	)
	mappings = make(map[string]string)
	regions, err = readMongodbSupportRegions(meta)
	if err != nil {
		return mappings, err
	}
	current = *(meta.(*KsyunClient).mongodbconn.Config.Region)
	for _, region := range regions {
		if region.(map[string]interface{})["Code"].(string) == current {
			exist = true
			azs := region.(map[string]interface{})["AvailabilityZones"].([]interface{})
			for _, az := range azs {
				code := az.(map[string]interface{})["Code"].(string)
				name := az.(map[string]interface{})["Name"].(string)
				mappings[name] = code
			}
			break
		}
	}
	if !exist {
		return mappings, fmt.Errorf("region %s not support", current)
	}
	return mappings, err
}

func readMongodbSupportAvailabilityZones(meta interface{}) (availabilityZones []string, err error) {
	var (
		current string
		exist   bool
		regions []interface{}
	)
	regions, err = readMongodbSupportRegions(meta)
	if err != nil {
		return availabilityZones, err
	}
	current = *(meta.(*KsyunClient).mongodbconn.Config.Region)
	for _, region := range regions {
		if region.(map[string]interface{})["Code"].(string) == current {
			exist = true
			azs := region.(map[string]interface{})["AvailabilityZones"].([]interface{})
			for _, az := range azs {
				code := az.(map[string]interface{})["Code"].(string)
				availabilityZones = append(availabilityZones, code)
			}
			break
		}
	}
	if !exist {
		return availabilityZones, fmt.Errorf("region %s not support", current)
	}
	return availabilityZones, err
}

func readMongodbInstance(d *schema.ResourceData, meta interface{}, instanceId string) (data map[string]interface{}, err error) {
	var (
		resp                  *map[string]interface{}
		mongoDBInstanceResult interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	conn := meta.(*KsyunClient).mongodbconn
	action := "DescribeMongoDBInstance"
	req := map[string]interface{}{
		"InstanceId": instanceId,
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.DescribeMongoDBInstance(&req)
	if err != nil {
		return data, err
	}
	mongoDBInstanceResult, err = getSdkValue("MongoDBInstanceResult", *resp)
	if err != nil {
		return data, err
	}
	data = mongoDBInstanceResult.(map[string]interface{})
	if v, ok := data["InstanceType"]; ok && v.(string) == "Cluster" {
		data["TotalStorage"] = data["Storage"]
		data["Storage"] = d.Get("storage")
	}
	return data, err
}

func readMongodbSecurityGroupRules(d *schema.ResourceData, meta interface{}, instanceId string) (rules []string, err error) {
	var (
		resp                     *map[string]interface{}
		mongoDBSecurityGroupRule interface{}
	)
	conn := meta.(*KsyunClient).mongodbconn
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"InstanceId": instanceId,
	}
	action := "ListSecurityGroupRules"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.ListSecurityGroupRules(&req)
	if err != nil {
		return rules, err
	}
	mongoDBSecurityGroupRule, err = getSdkValue("MongoDBSecurityGroupRule", *resp)
	if err != nil {
		return rules, err
	}
	for _, rule := range mongoDBSecurityGroupRule.([]interface{}) {
		if r, ok := rule.(map[string]interface{}); ok {
			rules = append(rules, r["cidr"].(string))
		}
	}
	return rules, err
}

func checkMongodbAvailabilityZonesValid(d *schema.ResourceData, meta interface{}) (err error) {
	if availabilityZone, ok := d.GetOk("availability_zone"); ok {
		var (
			supportAzs []string
		)
		supportAzs, err = readMongodbSupportAvailabilityZones(meta)
		if err != nil {
			return err
		}
		azs := strings.Split(availabilityZone.(string), ",")
		for _, az := range azs {
			exist := false
			for _, supportAz := range supportAzs {
				if supportAz == az {
					exist = true
					break
				}
			}
			//az not support
			if !exist {
				return fmt.Errorf("availability_zone %s not support in region %s ", az, *meta.(*KsyunClient).mongodbconn.Config.Region)
			}
		}
	}
	return err
}

func createMongodbInstance(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	var (
		resp   *map[string]interface{}
		id     interface{}
		action string
	)
	transform := map[string]SdkReqTransform{
		"availability_zone": {
			mapping: "AvailabilityZone",
			Type:    TransformWithN,
		},
		"cidrs": {
			Ignore: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	conn := meta.(*KsyunClient).mongodbconn
	if _, ok := req["ShardClass"]; ok {
		action = "CreateMongoDBShardInstance"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.CreateMongoDBShardInstance(&req)
	} else {
		req["InstanceType"] = "HighIO"
		action = "CreateMongoDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.CreateMongoDBInstance(&req)
	}
	if err != nil {
		return err
	}
	id, err = getSdkValue("MongoDBInstanceResult.InstanceId", *resp)
	if err != nil {
		return err
	}
	d.SetId(id.(string))
	return err
}

func mongodbStateRefreshFunc(d *schema.ResourceData, meta interface{}, instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			err error
		)
		data, err := readMongodbInstance(d, meta, instanceId)
		if err != nil {
			return nil, "", err
		}

		status := data["Status"].(string)

		for _, v := range failStates {
			if v == status {
				return nil, "", fmt.Errorf("instance status  error, status:%v", status)
			}
		}
		return data, status, nil
	}
}

func checkMongodbState(d *schema.ResourceData, meta interface{}, instanceId string, timeout time.Duration) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"running"},
		Refresh:    mongodbStateRefreshFunc(d, meta, instanceId, []string{"error"}),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func checkMongodbSecurityGroupRulesDel(d *schema.ResourceData, meta interface{}, instanceId string, del string) (err error, delV4 string, delV6 string) {
	var (
		cidrs    []string
		delCidrs []string
	)
	cidrs, err = readMongodbSecurityGroupRules(d, meta, instanceId)
	if err != nil {
		return err, delV4, delV6
	}
	delCidrs = strings.Split(del, ",")
	for _, cidr := range delCidrs {
		exist := false
		for _, current := range cidrs {
			if current == cidr {
				exist = true
				break
			}
		}
		if exist {
			if strings.Contains(cidr, ":") {
				delV6 = delV6 + cidr + ","
			} else {
				delV4 = delV4 + cidr + ","
			}
		}
	}
	return err, delV4, delV6
}

func checkMongodbSecurityGroupRulesChange(d *schema.ResourceData, meta interface{}, field string, instanceId string) (err error, addV4 string, delV4 string, addV6 string, delV6 string) {
	if d.HasChange(field) && !d.IsNewResource() {
		var (
			cidrs    []string
			addCidrs []string
			delCidrs []string
		)
		cidrs, err = readMongodbSecurityGroupRules(d, meta, instanceId)
		if err != nil {
			return err, addV4, delV4, addV6, delV6
		}
		addCidrs, delCidrs = stringSplitChange(",", field, cidrs, d)
		for _, cidr := range addCidrs {
			if strings.Contains(cidr, ":") {
				addV6 = addV6 + cidr + ","
			} else {
				addV4 = addV4 + cidr + ","
			}

		}
		for _, cidr := range delCidrs {
			if strings.Contains(cidr, ":") {
				delV6 = delV6 + cidr + ","
			} else {
				delV4 = delV4 + cidr + ","
			}
		}
	} else if cidrs, ok := d.GetOk(field); ok && d.IsNewResource() && cidrs.(string) != "" {
		for _, cidr := range strings.Split(cidrs.(string), ",") {
			if strings.Contains(cidr, ":") {
				addV6 = addV6 + cidr + ","
			} else {
				addV4 = addV4 + cidr + ","
			}
		}
	}
	if len(addV4) > 0 {
		addV4 = addV4[0 : len(addV4)-1]
	}
	if len(addV6) > 0 {
		addV6 = addV6[0 : len(addV6)-1]
	}
	if len(delV4) > 0 {
		delV4 = delV4[0 : len(delV4)-1]
	}
	if len(delV6) > 0 {
		delV6 = delV6[0 : len(delV6)-1]
	}
	return err, addV4, delV4, addV6, delV6
}

func addMongodbSecurityGroupRules(d *schema.ResourceData, meta interface{}, instanceId string, v4Cidrs string, v6Cidrs string) (err error) {
	if len(v4Cidrs) > 0 || len(v6Cidrs) > 0 {
		if instanceId == "" {
			instanceId = d.Id()
		}
		conn := meta.(*KsyunClient).mongodbconn
		req := make(map[string]interface{})
		action := "AddSecurityGroupRule"
		req["InstanceId"] = instanceId
		if len(v4Cidrs) > 0 {
			req["type"] = "IPV4"
			req["cidrs"] = v4Cidrs
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.AddSecurityGroupRule(&req)
			if err != nil {
				return err
			}
		}
		if len(v6Cidrs) > 0 {
			req["type"] = "IPV6"
			req["cidrs"] = v6Cidrs
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.AddSecurityGroupRule(&req)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func delMongodbSecurityGroupRules(d *schema.ResourceData, meta interface{}, instanceId string, v4Cidrs string, v6Cidrs string) (err error) {
	if len(v4Cidrs) > 0 || len(v6Cidrs) > 0 {
		if instanceId == "" {
			instanceId = d.Id()
		}
		conn := meta.(*KsyunClient).mongodbconn
		req := make(map[string]interface{})
		action := "DeleteSecurityGroupRules"
		req["InstanceId"] = instanceId
		if len(v4Cidrs) > 0 {
			req["type"] = "IPV4"
			req["cidrs"] = v4Cidrs
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.DeleteSecurityGroupRules(&req)
			if err != nil {
				return err
			}
		}
		if len(v6Cidrs) > 0 {
			req["type"] = "IPV6"
			req["cidrs"] = v6Cidrs
			logger.Debug(logger.ReqFormat, action, req)
			_, err = conn.DeleteSecurityGroupRules(&req)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func modifyMongodbInstanceNameAndProject(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"name":           {},
		"iam_project_id": {mapping: "ProjectId"},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return err
	}
	err = ModifyProjectInstance(d.Id(), &req, meta)
	if err != nil {
		return err
	}
	if len(req) > 0 {
		err = checkMongodbState(d, meta, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		req["InstanceId"] = d.Id()
		conn := meta.(*KsyunClient).mongodbconn
		action := "RenameMongoDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.RenameMongoDBInstance(&req)
	}
	return err
}

func modifyMongodbInstanceSpec(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"instance_class": {},
		"storage":        {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return err
	}
	if d.Get("instance_type").(string) == "HighIO" && len(req) > 0 {
		if _, ok := req["Storage"]; !ok {
			req["Storage"] = d.Get("storage")
		}
		if _, ok := req["InstanceClass"]; !ok {
			req["InstanceClass"] = d.Get("instance_class")
		}
		err = checkMongodbState(d, meta, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		req["InstanceId"] = d.Id()
		conn := meta.(*KsyunClient).mongodbconn
		action := "UpdateMongoDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.UpdateMongoDBInstance(&req)
	}
	return err
}

func modifyMongodbInstancePassword(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"instance_password": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return err
	}
	if len(req) > 0 {
		err = checkMongodbState(d, meta, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		req["InstanceId"] = d.Id()
		req["InstanceAccount"] = d.Get("instance_account")
		conn := meta.(*KsyunClient).mongodbconn
		action := "ResetPasswordMongoDBInstance"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ResetPasswordMongoDBInstance(&req)
	}
	return err
}

func modifyMongodbInstanceNodeNum(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"node_num": {},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil)
	if err != nil {
		return err
	}
	if len(req) > 0 {
		var (
			oldNodeNum interface{}
			newNodeNum interface{}
		)
		oldNodeNum, newNodeNum = d.GetChange("node_num")
		if newNodeNum.(int) < oldNodeNum.(int) {
			err = fmt.Errorf("cat set node_num %v less then %v", newNodeNum, oldNodeNum)
			return err
		}
		err = checkMongodbState(d, meta, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		req["InstanceId"] = d.Id()
		conn := meta.(*KsyunClient).mongodbconn
		action := "AddSecondaryInstance"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.AddSecondaryInstance(&req)
	}
	return err
}

func modifyMongodbInstanceCommon(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	//valid cidrs
	err, addV4, delV4, addV6, delV6 := checkMongodbSecurityGroupRulesChange(d, meta, "cidrs", "")
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	// modify name if need
	err = modifyMongodbInstanceNameAndProject(d, meta, r)
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	// modify password if need
	err = modifyMongodbInstancePassword(d, meta, r)
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	// modify node num if need
	err = modifyMongodbInstanceNodeNum(d, meta, r)
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	// resize if need
	err = modifyMongodbInstanceSpec(d, meta, r)
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	//wait mongodb state
	err = checkMongodbState(d, meta, "", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	// modify sg if need
	err = addMongodbSecurityGroupRules(d, meta, "", addV4, addV6)
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	// modify sg if need
	err = delMongodbSecurityGroupRules(d, meta, "", delV4, delV6)
	if err != nil {
		return fmt.Errorf("error on update instance %q, %s", d.Id(), err)
	}
	return err
}

func createMongodbInstanceCommon(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	//valid availability_zone
	err = checkMongodbAvailabilityZonesValid(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	//valid cidrs
	err, addV4, _, addV6, _ := checkMongodbSecurityGroupRulesChange(d, meta, "cidrs", "")
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	//create mongodb instance
	err = createMongodbInstance(d, meta, r)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	//wait mongodb state
	err = checkMongodbState(d, meta, "", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	// set sg if need
	err = addMongodbSecurityGroupRules(d, meta, "", addV4, addV6)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	return err
}

func readMongodbInstanceCommon(d *schema.ResourceData, meta interface{}, r *schema.Resource) (err error) {
	var (
		cidrs string
	)
	data, err := readMongodbInstance(d, meta, "")
	if err != nil {
		return err
	}
	mappings, err := readMongodbSupportAzMappings(meta)
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"IP": {Field: "ip"},
		"Version": {
			Field: "db_version",
			FieldRespFunc: func(i interface{}) interface{} {
				return strings.Replace(i.(string), "mongodb ", "", -1)
			},
		},
		"Area": {
			Field: "availability_zone",
			FieldRespFunc: func(i interface{}) interface{} {
				return mappings[i.(string)]
			},
		},
	}
	//rules
	cidrs, err = readMongodbSecurityGroupCidrs(d, meta, "cidrs", "")
	if err != nil {
		return err
	}
	if cidrs != "" {
		data["Cidrs"] = cidrs
	}
	//special
	if _, ok := data["InstanceAccount"]; !ok {
		err = d.Set("instance_account", "root")
	}
	SdkResponseAutoResourceData(d, r, data, extra)
	return err
}

func removeMongodbInstance(d *schema.ResourceData, meta interface{}) (err error) {
	conn := meta.(*KsyunClient).mongodbconn
	action := "DeleteMongoDBInstance"
	req := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.DeleteMongoDBInstance(&req)
	return err
}

func canNotFoundMongodbError(err error) bool {
	if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 404 {
		return true
	}
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		return true
	}
	return false
}

func mongodbShardInstanceSchemaDiffSuppressFunc() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		if old == "" {
			return false
		}
		return true
	}
}

func mongodbShardInstanceCustomizeDiffFunc() schema.CustomizeDiffFunc {
	return func(diff *schema.ResourceDiff, i interface{}) (err error) {
		extra := []string{"shard_num", "mongos_num", "shard_class", "mongos_class", "storage"}
		for _, v := range extra {
			if diff.HasChange(v) {
				o, _ := diff.GetChange(v)
				if s, ok := o.(string); ok && s != "" {
					return fmt.Errorf("%s not support update,please use ksyun_mongodb_shard_instance_node", v)
				} else if s, ok := o.(int); ok && s != 0 {
					return fmt.Errorf("%s not support update,please use ksyun_mongodb_shard_instance_node", v)
				}
			}
		}
		return err
	}
}

func mongodbShardInstanceNodeImportStateFunc() schema.StateFunc {
	return func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		var err error
		items := strings.Split(d.Id(), ":")
		if len(items) != 2 {
			return nil, fmt.Errorf("id must split with %s and size %v", ":", 2)
		}
		err = d.Set("instance_id", items[0])
		if err != nil {
			return nil, err
		}
		err = d.Set("node_id", items[1])
		if err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, err
	}
}

func readMongodbShardInstanceNodes(d *schema.ResourceData, meta interface{}) (mongosNodeResult interface{}, shardNodeResult interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	conn := meta.(*KsyunClient).mongodbconn
	req := map[string]interface{}{
		"InstanceId": d.Get("instance_id"),
	}
	resp, err = conn.DescribeMongoDBShardNode(&req)
	if err != nil {
		return mongosNodeResult, shardNodeResult, err
	}
	mongosNodeResult, err = getSdkValue("MongosNodeResult", *resp)
	if err != nil {
		return mongosNodeResult, shardNodeResult, err
	}
	shardNodeResult, err = getSdkValue("ShardNodeResult", *resp)
	if err != nil {
		return mongosNodeResult, shardNodeResult, err
	}
	if shardNodeResult == nil || mongosNodeResult == nil {
		err = fmt.Errorf("read shard mongo instance error")
	}
	return mongosNodeResult, shardNodeResult, err
}

func readMongodbShardInstanceNode(d *schema.ResourceData, meta interface{}) (data interface{}, extra map[string]SdkResponseMapping, err error) {
	var (
		mongosNodeResult interface{}
		shardNodeResult  interface{}
		exist            bool
	)

	mongosNodeResult, shardNodeResult, err = readMongodbShardInstanceNodes(d, meta)
	extra = make(map[string]SdkResponseMapping)
	if err != nil {
		return data, extra, fmt.Errorf("read shard instance node error $%s ", err)
	}
	for _, v := range mongosNodeResult.([]interface{}) {
		if v.(map[string]interface{})["NodeId"].(string) == d.Get("node_id").(string) {
			exist = true
			v.(map[string]interface{})["NodeType"] = "mongos"
			extra = map[string]SdkResponseMapping{
				"InstanceClass": {Field: "node_class"},
			}
			return v, extra, err
		}
	}
	if !exist {
		for _, v := range shardNodeResult.([]interface{}) {
			if v.(map[string]interface{})["NodeId"].(string) == d.Get("node_id").(string) {
				exist = true
				v.(map[string]interface{})["NodeType"] = "shard"
				extra = map[string]SdkResponseMapping{
					"InstanceClass": {Field: "node_class"},
					"Disk":          {Field: "node_storage"},
				}
				return v, extra, err
			}
		}
	}
	if !exist {
		return data, extra, fmt.Errorf("mongodb shard instance node %s not found", d.Get("node_id"))
	}
	return data, extra, err
}

func modifyMongodbShardInstanceNode(d *schema.ResourceData, meta interface{}) (err error) {
	transform := map[string]SdkReqTransform{
		"node_class":   {mapping: "InstanceClass"},
		"node_storage": {mapping: "Storage"},
	}
	req, err := SdkRequestAutoMapping(d, resourceKsyunMongodbShardInstanceNode(), false, transform, nil)
	if len(req) > 0 {
		req["InstanceId"] = d.Get("instance_id")
		req["NodeId"] = d.Get("node_id")
		req["NodeType"] = d.Get("node_type")
		conn := meta.(*KsyunClient).mongodbconn
		_, err = conn.UpdateMongoDBInstanceCluster(&req)
		if err != nil {
			return err
		}
		err = checkMongodbState(d, meta, d.Get("instance_id").(string), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}
	return err
}

func delMongodbShardInstanceNode(d *schema.ResourceData, meta interface{}) (err error) {
	instanceId := d.Get("instance_id").(string)
	err = checkMongodbState(d, meta, instanceId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	req := map[string]interface{}{
		"InstanceId": instanceId,
		"NodeId":     d.Get("node_id"),
	}
	conn := meta.(*KsyunClient).mongodbconn
	_, err = conn.DeleteClusterNode(&req)
	if err != nil {
		return err
	}
	err = checkMongodbState(d, meta, instanceId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return err
}

func createMongodbShardInstanceNode(d *schema.ResourceData, meta interface{}) (err error) {
	var (
		mongosNodeResultOld interface{}
		shardNodeResultOld  interface{}
		mongosNodeResultNew interface{}
		shardNodeResultNew  interface{}
	)

	req, err := SdkRequestAutoMapping(d, resourceKsyunMongodbShardInstanceNode(), false, nil, nil)
	if err != nil {
		return err
	}
	err = checkMongodbState(d, meta, d.Get("instance_id").(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	mongosNodeResultOld, shardNodeResultOld, err = readMongodbShardInstanceNodes(d, meta)
	if err != nil {
		return err
	}
	conn := meta.(*KsyunClient).mongodbconn
	action := "AddClusterNode"
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.AddClusterNode(&req)
	if err != nil {
		return err
	}
	err = checkMongodbState(d, meta, d.Get("instance_id").(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	mongosNodeResultNew, shardNodeResultNew, err = readMongodbShardInstanceNodes(d, meta)
	if err != nil {
		return err
	}
	if d.Get("node_type") == "mongos" {
		for _, n := range mongosNodeResultNew.([]interface{}) {
			exist := false
			nId := n.(map[string]interface{})["NodeId"].(string)
			for _, o := range mongosNodeResultOld.([]interface{}) {
				oId := o.(map[string]interface{})["NodeId"].(string)
				if oId == nId {
					exist = true
					break
				}
			}
			if !exist {
				d.SetId(d.Get("instance_id").(string) + ":" + nId)
				err = d.Set("node_id", nId)
			}
		}
	} else {
		for _, n := range shardNodeResultNew.([]interface{}) {
			exist := false
			nId := n.(map[string]interface{})["NodeId"].(string)
			for _, o := range shardNodeResultOld.([]interface{}) {
				oId := o.(map[string]interface{})["NodeId"].(string)
				if oId == nId {
					exist = true
					break
				}
			}
			if !exist {
				d.SetId(d.Get("instance_id").(string) + ":" + nId)
				err = d.Set("node_id", nId)
			}
		}
	}

	if err != nil {
		return err
	}
	return err
}

func readMongodbSecurityGroupCidrs(d *schema.ResourceData, meta interface{}, field string, instanceId string) (cidrs string, err error) {
	var (
		rules        []string
		currentRules []string
		ruleStr      string
	)
	rules, err = readMongodbSecurityGroupRules(d, meta, instanceId)
	if err != nil {
		return cidrs, err
	}
	for _, rule := range rules {
		currentRules = append(currentRules, rule)
	}
	ruleStr = stringSplitRead(",", field, currentRules, d)
	if ruleStr != "" {
		cidrs = ruleStr
	}
	return cidrs, err
}
