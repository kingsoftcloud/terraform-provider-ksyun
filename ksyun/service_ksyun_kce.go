package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type KceService struct {
	client *KsyunClient
}

func (s *KceService) readKceClusters(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp           *map[string]interface{}
		clusterResults interface{}
	)
	conn := s.client.kceconn
	//action := "DescribeCluster"

	if condition == nil {
		resp, err = conn.DescribeCluster(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeCluster(&condition)
		if err != nil {
			return data, err
		}
	}
	clusterResults, err = getSdkValue("ClusterSet", *resp)
	if err != nil {
		return data, err
	}
	data = clusterResults.([]interface{})
	logger.Debug("kce list", "123", data)
	return data, err
}

func (s *KceService) ReadAndSetKceClusters(d *schema.ResourceData, r *schema.Resource) (err error) {

	transform := map[string]SdkReqTransform{
		"cluster_id": {
			mapping: "ClusterId",
			Type:    TransformDefault,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readKceClusters(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "ClusterId",
		nameField:   "ClusterName",
		targetField: "cluster_set",
		extra: map[string]SdkResponseMapping{
			"EnableKMSE": {
				Field: "enable_kmse",
			},
		},
	})
}
