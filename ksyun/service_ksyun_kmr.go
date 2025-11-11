package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kmr "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/kmr/v20210902"
	"github.com/mitchellh/mapstructure"
)

type KmrService struct {
	client *KsyunClient
}

func (s *KmrService) ReadClusters(req map[string]interface{}) (resp *kmr.ListClustersResponse, data []interface{}, err error) {
	kmrReq := kmr.NewListClustersRequest()
	if marker, ok := req["Marker"].(string); ok {
		kmrReq.Marker = &marker
	}
	respInterface, err := s.client.WithKmrClient(func(conn *kmr.Client) (interface{}, error) {
		return conn.ListClustersSend(kmrReq)
	})
	if err != nil {
		return nil, nil, err
	}
	resp = respInterface.(*kmr.ListClustersResponse)
	data = make([]interface{}, len(resp.Clusters))
	for i, cluster := range resp.Clusters {
		var m map[string]interface{}
		mapstructure.Decode(cluster, &m)
		// 处理指针字段和嵌套
		for k, v := range m {
			if ptr, ok := v.(*string); ok && ptr != nil {
				m[k] = *ptr
			} else if ptrInt, ok := v.(*int); ok && ptrInt != nil {
				m[k] = *ptrInt
			} else if ptrBool, ok := v.(*bool); ok && ptrBool != nil {
				m[k] = *ptrBool
			}
		}
		// 手动处理instance_groups
		igMaps := make([]map[string]interface{}, len(cluster.InstanceGroups))
		for j, group := range cluster.InstanceGroups {
			igMap := map[string]interface{}{}
			if group.Id != nil {
				igMap["id"] = *group.Id
			}
			if group.InstanceGroupType != nil {
				igMap["instance_group_type"] = *group.InstanceGroupType
			}
			if group.ResourceType != nil {
				igMap["resource_type"] = *group.ResourceType
			}
			if group.InstanceType != nil {
				igMap["instance_type"] = *group.InstanceType
			}
			igMaps[j] = igMap
		}
		var igInterfaces []interface{}
		for _, ig := range igMaps {
			igInterfaces = append(igInterfaces, ig)
		}
		m["instance_groups"] = igInterfaces
		//fmt.Printf("cluster %d, instance_groups type: %T\n", i, m["instance_groups"])
		data[i] = m
	}
	//fmt.Printf("Final clusters data: %+v\n", data)
	return resp, data, err
}

func (s *KmrService) ReadAndSetClusters(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"marker": {
			mapping: "Marker",
			Type:    TransformDefault,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	resp, data, err := s.ReadClusters(req)
	if err != nil {
		return err
	}
	err = mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "ClusterName",
		idFiled:     "ClusterId",
		targetField: "clusters",
		extra:       map[string]SdkResponseMapping{},
	})
	d.Set("total", resp.Total)
	return err
}
