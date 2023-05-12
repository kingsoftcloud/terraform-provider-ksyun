package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type KceWorkerService struct {
	client *KsyunClient
}

func (s *KceWorkerService) AddWorker(d *schema.ResourceData, resource *schema.Resource) (err error) {
	clusterId := d.Get("cluster_id")
	instanceId := d.Get("instance_id")
	imageId := d.Get("image_id")

	//var nodes []interface{}
	//nodes, err = s.getNode(clusterId.(string), instanceId.(string))
	//if err != nil {
	//	return
	//}
	//
	//// 检查是否已经绑定
	//if len(nodes) > 0 {
	//
	//}
	//logger.Debug("%s", "AddWorker", nodes, err)
	return
}

func (s *KceWorkerService) addNode() {}

func (s *KceWorkerService) getNode(clusterId, instanceId string) (list []interface{}, err error) {
	conn := s.client.kceconn
	condition := map[string]interface{}{
		"ClusterId":        clusterId,
		"Filter.1.Name":    "instance-id",
		"Filter.1.Value.1": instanceId,
	}

	var resp *map[string]interface{}
	resp, err = conn.DescribeClusterInstance(&condition)

	if err != nil {
		return
	}
	var clusterInstanceResults interface{}
	clusterInstanceResults, err = getSdkValue("InstanceSet", *resp)
	if err != nil {
		return
	}
	list = clusterInstanceResults.([]interface{})
	return
}
