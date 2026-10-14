package ksyun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kes "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/kes/v20201215"
	"github.com/mitchellh/mapstructure"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type KesService struct {
	client *KsyunClient
}

// derefPointerFields turns *string/*int/*bool values produced by mapstructure
// into plain values so that SdkResponseAutoMapping can consume them.
func derefPointerFields(m map[string]interface{}) {
	for k, v := range m {
		if ptr, ok := v.(*string); ok && ptr != nil {
			m[k] = *ptr
		} else if ptrInt, ok := v.(*int); ok && ptrInt != nil {
			m[k] = *ptrInt
		} else if ptrBool, ok := v.(*bool); ok && ptrBool != nil {
			m[k] = *ptrBool
		}
	}
}

// ReadClusters queries the KES cluster list via ListClusters.
func (s *KesService) ReadClusters(req map[string]interface{}) (resp *kes.ListClustersResponse, data []interface{}, err error) {
	kesReq := kes.NewListClustersRequest()
	if marker, ok := req["Marker"].(string); ok {
		kesReq.Marker = &marker
	}
	respInterface, err := s.client.WithKesClient(func(conn *kes.Client) (interface{}, error) {
		return conn.ListClustersSend(kesReq)
	})
	if err != nil {
		return nil, nil, err
	}
	resp = respInterface.(*kes.ListClustersResponse)
	data = make([]interface{}, len(resp.Clusters))
	for i, cluster := range resp.Clusters {
		var m map[string]interface{}
		mapstructure.Decode(cluster, &m)
		derefPointerFields(m)
		// manually expand the nested instance groups (lite shape)
		igInterfaces := make([]interface{}, len(cluster.InstanceGroups))
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
			igInterfaces[j] = igMap
		}
		m["instance_groups"] = igInterfaces
		data[i] = m
	}
	return resp, data, err
}

// ReadAndSetClusters reads the cluster list and writes it into schema.
func (s *KesService) ReadAndSetClusters(d *schema.ResourceData, r *schema.Resource) (err error) {
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
	if err != nil {
		return err
	}
	if resp.Total != nil {
		d.Set("total", *resp.Total)
	}
	return err
}

// ReadCluster queries a single KES cluster detail via DescribeCluster.
// It returns a one-element collection so that mergeDataSourcesResp can be reused.
func (s *KesService) ReadCluster(req map[string]interface{}) (resp *kes.DescribeClusterResponse, data []interface{}, err error) {
	kesReq := kes.NewDescribeClusterRequest()
	if clusterId, ok := req["ClusterId"].(string); ok {
		kesReq.ClusterId = &clusterId
	}
	respInterface, err := s.client.WithKesClient(func(conn *kes.Client) (interface{}, error) {
		return conn.DescribeClusterSend(kesReq)
	})
	if err != nil {
		return nil, nil, err
	}
	resp = respInterface.(*kes.DescribeClusterResponse)
	var m map[string]interface{}
	mapstructure.Decode(resp, &m)
	derefPointerFields(m)
	// manually expand the nested instance groups (with volume info)
	igInterfaces := make([]interface{}, len(resp.InstanceGroups))
	for j, group := range resp.InstanceGroups {
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
		if group.VolumeSize != nil {
			igMap["volume_size"] = *group.VolumeSize
		}
		if group.VolumeType != nil {
			igMap["volume_type"] = *group.VolumeType
		}
		igInterfaces[j] = igMap
	}
	m["instance_groups"] = igInterfaces
	// manually expand the tags
	tagInterfaces := make([]interface{}, len(resp.Tags))
	for j, tag := range resp.Tags {
		tagMap := map[string]interface{}{}
		if tag.TagKey != nil {
			tagMap["tag_key"] = *tag.TagKey
		}
		if tag.TagValue != nil {
			tagMap["tag_value"] = *tag.TagValue
		}
		tagInterfaces[j] = tagMap
	}
	m["tags"] = tagInterfaces
	data = []interface{}{m}
	return resp, data, err
}

// ReadAndSetCluster reads a single cluster detail and writes it into schema.
func (s *KesService) ReadAndSetCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
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
	_, data, err := s.ReadCluster(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "ClusterName",
		idFiled:     "ClusterId",
		targetField: "clusters",
		extra:       map[string]SdkResponseMapping{},
	})
}

// ReadInstanceGroups queries the instance groups of a KES cluster via ListInstanceGroups.
func (s *KesService) ReadInstanceGroups(req map[string]interface{}) (resp *kes.ListInstanceGroupsResponse, data []interface{}, err error) {
	kesReq := kes.NewListInstanceGroupsRequest()
	if clusterId, ok := req["ClusterId"].(string); ok {
		kesReq.ClusterId = &clusterId
	}
	respInterface, err := s.client.WithKesClient(func(conn *kes.Client) (interface{}, error) {
		return conn.ListInstanceGroupsSend(kesReq)
	})
	if err != nil {
		return nil, nil, err
	}
	resp = respInterface.(*kes.ListInstanceGroupsResponse)
	data = make([]interface{}, len(resp.InstanceGroups))
	for i, group := range resp.InstanceGroups {
		var m map[string]interface{}
		mapstructure.Decode(group, &m)
		derefPointerFields(m)
		// manually expand the nested instances
		instInterfaces := make([]interface{}, len(group.Instances))
		for j, inst := range group.Instances {
			instMap := map[string]interface{}{}
			if inst.Id != nil {
				instMap["id"] = *inst.Id
			}
			if inst.InstanceGroupId != nil {
				instMap["instance_group_id"] = *inst.InstanceGroupId
			}
			if inst.InstanceId != nil {
				instMap["instance_id"] = *inst.InstanceId
			}
			if inst.InstanceName != nil {
				instMap["instance_name"] = *inst.InstanceName
			}
			if inst.InternalIp != nil {
				instMap["internal_ip"] = *inst.InternalIp
			}
			volInterfaces := make([]interface{}, len(inst.Volumes))
			for k, vol := range inst.Volumes {
				if vol != nil {
					volInterfaces[k] = *vol
				}
			}
			instMap["volumes"] = volInterfaces
			instInterfaces[j] = instMap
		}
		m["instances"] = instInterfaces
		data[i] = m
	}
	return resp, data, err
}

// ReadAndSetInstanceGroups reads instance groups and writes them into schema.
func (s *KesService) ReadAndSetInstanceGroups(d *schema.ResourceData, r *schema.Resource) (err error) {
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
	_, data, err := s.ReadInstanceGroups(req)
	if err != nil {
		return err
	}
	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "InstanceGroupType",
		idFiled:     "Id",
		targetField: "instance_groups",
		extra:       map[string]SdkResponseMapping{},
	})
}

// ptrString returns a pointer to the given string. The KES LaunchCluster SDK
// request uses pointer fields, and the project has no generic string->*string
// helper (only the reverse derefPointerFields).
func ptrString(s string) *string { return &s }

// ptrInt returns a pointer to the given int.
func ptrInt(i int) *int { return &i }

// ptrBool returns a pointer to the given bool.
func ptrBool(b bool) *bool { return &b }

// CreateCluster launches a KES cluster (v20201215 standard / KES7) and reads
// the result back into the resource state.
func (s *KesService) CreateCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
	kesReq := kes.NewLaunchClusterRequest()
	kesReq.ClusterName = ptrString(d.Get("cluster_name").(string))
	kesReq.MainVersion = ptrString(d.Get("main_version").(string))
	kesReq.Scenario = ptrString(d.Get("scenario").(string))
	kesReq.AvailabilityZone = ptrString(d.Get("availability_zone").(string))
	kesReq.ChargeType = ptrString(d.Get("charge_type").(string))
	kesReq.PurchaseTime = ptrInt(d.Get("purchase_time").(int))
	kesReq.ProjectId = ptrInt(d.Get("project_id").(int))
	kesReq.VpcDomainId = ptrString(d.Get("vpc_domain_id").(string))
	kesReq.VpcSubnetId = ptrString(d.Get("vpc_subnet_id").(string))
	if v := d.Get("vpc_epc_subnet_id").(string); v != "" {
		kesReq.VpcEpcSubnetId = ptrString(v)
	}
	if v := d.Get("security_group_id").(string); v != "" {
		kesReq.SecurityGroupId = ptrString(v)
	}

	// build instance groups from the schema list
	igRaw := d.Get("instance_groups").([]interface{})
	igList := make([]*kes.LaunchClusterInstanceGroups, 0, len(igRaw))
	for _, raw := range igRaw {
		if raw == nil {
			continue
		}
		m := raw.(map[string]interface{})
		g := &kes.LaunchClusterInstanceGroups{}
		if v, ok := m["instance_group_type"].(string); ok && v != "" {
			g.InstanceGroupType = ptrString(v)
		}
		if v, ok := m["instance_type"].(string); ok && v != "" {
			g.InstanceType = ptrString(v)
		}
		if v, ok := m["resource_type"].(string); ok && v != "" {
			g.ResourceType = ptrString(v)
		}
		g.InstanceCount = ptrInt(m["instance_count"].(int))
		if v, ok := m["multi_instance_count"].(int); ok && v > 0 {
			g.MultiInstanceCount = ptrInt(v)
		}
		if v, ok := m["raid_type"].(string); ok && v != "" {
			g.RaidType = ptrString(v)
		}
		if v, ok := m["volume_type"].(string); ok && v != "" {
			g.VolumeType = ptrString(v)
		}
		g.VolumeSize = ptrInt(m["volume_size"].(int))
		g.VolumeCount = ptrInt(m["volume_count"].(int))
		g.SystemDiskType = ptrString(m["system_disk_type"].(string))
		g.SystemDiskSize = ptrInt(m["system_disk_size"].(int))
		igList = append(igList, g)
	}
	kesReq.InstanceGroups = igList

	logger.Debug(logger.ReqFormat, "LaunchCluster", kesReq)
	respInterface, err := s.client.WithKesClient(func(conn *kes.Client) (interface{}, error) {
		return conn.LaunchClusterSend(kesReq)
	})
	if err != nil {
		return fmt.Errorf("error on creating KES cluster: %s", err)
	}
	resp := respInterface.(*kes.LaunchClusterResponse)
	if resp.ClusterId == nil || *resp.ClusterId == "" {
		return fmt.Errorf("KES cluster was created but no ClusterId was returned")
	}
	d.SetId(*resp.ClusterId)
	return s.ReadAndSetClusterResource(d, r)
}

// ReadAndSetClusterResource reads a single KES cluster by id and writes it
// into the resource state. Missing clusters are removed from state.
func (s *KesService) ReadAndSetClusterResource(d *schema.ResourceData, r *schema.Resource) (err error) {
	req := map[string]interface{}{"ClusterId": d.Id()}
	_, data, err := s.ReadCluster(req)
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading KES cluster %s: %s", d.Id(), err)
	}
	if len(data) == 0 {
		d.SetId("")
		return nil
	}
	// The DescribeCluster API only returns a partial instance_groups shape
	// (InstanceGroupType/InstanceType/ResourceType/VolumeSize/VolumeType) and
	// omits InstanceCount and the optional disk fields that are Required/ForceNew
	// in the schema. Overwriting instance_groups from this partial response would
	// zero out instance_count in state and trigger a perpetual diff (ForceNew).
	// Since every instance_groups field is ForceNew, drift on it can only ever
	// result in a destructive recreate, so we intentionally preserve the
	// config-provided values and only refresh the top-level Computed fields.
	m := data[0].(map[string]interface{})
	delete(m, "instance_groups")
	SdkResponseAutoResourceData(d, r, m, nil)
	return nil
}

// DeleteCluster is a no-op: the KES SDK does not provide a DeleteCluster API.
// The cluster must be deleted manually in the Ksyun console. terraform destroy
// only removes the resource from Terraform state.
func (s *KesService) DeleteCluster(d *schema.ResourceData, meta interface{}) (err error) {
	d.SetId("")
	return nil
}
