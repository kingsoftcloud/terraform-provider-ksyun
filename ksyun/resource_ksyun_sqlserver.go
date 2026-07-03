package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/sqlserver"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

var getSqlserverInTheCar = map[string]bool{
	"db_instance_identifier": true,
	"instance_create_time":   true,
	"port":                   true,
	"sub_order_id":           true,
	"availability_zone_1":    true,
	"availability_zone_2":    true,
	"region":                 true,
}

func resourceKsyunSqlServer() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunSqlServerCreate,
		Update: resourceKsyunSqlServerUpdate,
		Read:   resourceKsyunSqlServerRead,
		Delete: resourceKsyunSqlServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance identifier",
				Computed:    true,
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				Description: "this value regex db.ram.d{1,3}|db.disk.d{1,5} , " +
					"db.ram is rds random access memory size, db.disk is disk size",
			},
			"db_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HRDS_SS",
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "engine is db type, only support SQLServer",
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db engine version only support 2012sp4,2014sp2,2016sp2,2017",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"master_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"master_user_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bill_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "proprietary security group id for krds",
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"availability_zone_1": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"availability_zone_2": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Computed: true,
			},
			"sub_order_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_create_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func resourceKsyunSqlServerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	var err error
	creates := []string{
		"DBInstanceClass",
		"DBInstanceName",
		"DBInstanceType",
		"Engine",
		"EngineVersion",
		"MasterUserName",
		"MasterUserPassword",
		"VpcId",
		"SubnetId",
		"BillType",
		"Duration",
		"SecurityGroupId",
		"PreferredBackupTime",
		"AvailabilityZone.1",
		"AvailabilityZone.2",
		"ProjectId",
		"Port",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(Camel2Hungarian(v)); ok {
			createReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateDBInstance"
	logger.Debug(logger.RespFormat, action, createReq)
	resp, err = conn.CreateDBInstance(&createReq)
	logger.Debug(logger.AllFormat, action, createReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Instance(sqlserver): %s", err)
	}

	if resp != nil {
		bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
		if !dataOk {
			return fmt.Errorf("error on creating Instance(sqlserver): response body missing Data")
		}
		instances, instancesOk := bodyData["Instances"].([]interface{})
		if !instancesOk || len(instances) == 0 {
			return fmt.Errorf("error on creating Instance(sqlserver): response body missing Instances")
		}
		sqlserverInstance := instances[0].(map[string]interface{})
		instanceId := sqlserverInstance["DBInstanceIdentifier"].(string)
		logger.DebugInfo(" DBInstanceIdentifier : %v", instanceId)
		d.SetId(instanceId)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{tCreatingStatus},
		Target:     []string{tActiveStatus, tFailedStatus, tDeletedStatus, tStopedStatus},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    sqlserverInstanceStateRefreshForCreate(conn, d.Id(), []string{tCreatingStatus}),
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceKsyunSqlServerRead(d, meta)
}

func sqlserverInstanceStateRefreshForCreate(client *sqlserver.Sqlserver, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"DBInstanceIdentifier": instanceId}
		action := "DescribeDBInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeDBInstances(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return nil, "", err
		}
		bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
		if !dataOk {
			return nil, tCreatingStatus, nil
		}
		instances, instancesOk := bodyData["Instances"].([]interface{})
		if !instancesOk || len(instances) == 0 {
			return nil, tCreatingStatus, nil
		}
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)

		return resp, state, nil

	}
}

func resourceKsyunSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	req := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeDBInstances(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading Instance(sqlserver) %q, %s", d.Id(), err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		d.SetId("")
		return nil
	}
	instances, instancesOk := bodyData["Instances"].([]interface{})
	if !instancesOk || len(instances) == 0 {
		d.SetId("")
		return nil
	}

	sqlserverIds := make([]string, len(instances))
	sqlserverMap := make([]map[string]interface{}, len(instances))
	for k, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		instInfo := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "DBInstanceClass" {
			} else if k == "ReadReplicaDBInstanceIdentifiers" {
			} else if k == "DBSource" {
			} else {
				instInfo[Camel2Hungarian(k)] = v
			}
		}
		sqlserverMap[k] = instInfo
		logger.DebugInfo(" converted ---- %+v ", instInfo)

		sqlserverIds[k] = instInfo["db_instance_identifier"].(string)
		logger.DebugInfo("sqlserverIds fuck : %v", sqlserverIds)
	}

	logger.DebugInfo(" converted ---- %+v ", sqlserverMap)
	_ = SetDByFkResp(d, sqlserverMap[0], getSqlserverInTheCar)
	return nil
}

func resourceKsyunSqlServerUpdate(d *schema.ResourceData, meta interface{}) error {
	// 关闭事务，允许部分属性被修改  d.Partial(true) d.Partial(false)
	updateField := []string{
		"db_instance_class",
		"db_instance_name",
		"db_instance_type",
		"engine",
		"engine_version",
		"master_user_name",
		"master_user_password",
		"vpc_id",
		"subnet_id",
		"bill_type",
		"duration",
		"security_group_id",
		"preferred_backup_time",
		"availability_zone_1",
		"availability_zone_2",
		"project_id",
		"port",
	}
	d.Partial(true)
	for _, v := range updateField {
		if d.HasChange(v) {
			return fmt.Errorf("error on updating instance , sqlserver is not support update")
		}
	}
	d.Partial(false)
	return nil
}

func resourceKsyunSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	deleteReq := make(map[string]interface{})
	deleteReq["DBInstanceIdentifier"] = d.Id()

	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		readReq := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
		discribeAction := "DescribeInstances"
		logger.Debug(logger.ReqFormat, discribeAction, readReq)
		desResp, desErr := conn.DescribeDBInstances(&readReq)
		logger.Debug(logger.AllFormat, discribeAction, readReq, *desResp, desErr)

		if desErr != nil {
			if notFoundError(desErr) {
				return nil
			} else {
				return resource.NonRetryableError(desErr)
			}
		}

		bodyData, dataOk := (*desResp)["Data"].(map[string]interface{})
		if !dataOk {
			return nil
		}
		instances, instancesOk := bodyData["Instances"].([]interface{})
		if !instancesOk || len(instances) == 0 {
			return nil
		}
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)

		if state != tDeletedStatus {
			deleteAction := "DeleteDBInstance"
			logger.Debug(logger.ReqFormat, deleteAction, deleteReq)
			deleteResp, deleteErr := conn.DeleteDBInstance(&deleteReq)
			logger.Debug(logger.AllFormat, deleteAction, deleteReq, *deleteResp, deleteErr)
			if deleteErr == nil || notFoundError(deleteErr) {
				return nil
			}
			if deleteErr != nil {
				return resource.RetryableError(deleteErr)
			}

			logger.Debug(logger.ReqFormat, discribeAction, readReq)
			postDesResp, postDesErr := conn.DescribeDBInstances(&readReq)
			logger.Debug(logger.AllFormat, discribeAction, readReq, *postDesResp, postDesErr)

			if desErr != nil {
				if notFoundError(desErr) {
					return nil
				} else {
					return resource.NonRetryableError(fmt.Errorf("error on  reading kec when delete %q, %s", d.Id(), desErr))
				}
			}
		}

		return resource.RetryableError(desErr)
	})
}
