package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type KceWorkerService struct {
	client *KsyunClient
}

func (s *KceWorkerService) AddWorker(d *schema.ResourceData, resource *schema.Resource) (err error) {
	clusterId := d.Get("cluster_id")
	instanceId := d.Get("instance_id")
	//imageId := d.Get("image_id")

	node, err := s.getNode(clusterId.(string), instanceId.(string))
	logger.Debug("%s", "AddWorker", node, err)
	return
}

func (s *KceWorkerService) getNode(clusterId, instanceId string) (list []interface{}, err error) {
	conn := s.client.kceconn
	condition := map[string]interface{}{
		"ClusterId":        clusterId,
		"Filter.1.Name":    "instance-id",
		"Filter.1.Value.1": instanceId,
	}

	fmt.Println(condition)
	var resp *map[string]interface{}
	resp, err = conn.DescribeClusterInstance(&condition)

	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Debug("%v", "getNode", resp, err, err != nil)
	var clusterInstanceResults interface{}
	clusterInstanceResults, err = getSdkValue("InstanceSet", *resp)
	if err != nil {
		return
	}
	list = clusterInstanceResults.([]interface{})
	return
}
