package ksyun

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	kes "github.com/kingsoftcloud/sdk-go/v2/ksyun/client/kes/v20201215"
	ksyunhttp "github.com/kingsoftcloud/sdk-go/v2/ksyun/common/http"
	"github.com/mitchellh/mapstructure"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// KesVectorService manages KES8 (8.x/9.x) clusters.
//
// Authoritative path (confirmed by the Ksyun console capture + a production
// probe against the inner OpenAPI):
//
//   - Create -> the SHARED kmr/2021-09-02 LaunchCluster backend, switched into
//     KES8 mode by the body field ClusterType="KES". The console POSTs to
//     /kmr/v2/LaunchCluster (Version 2021-09-02) with ClusterType="KES". A probe
//     with Service="kmr", Version="2021-09-02", ClusterType="KES" reached the
//     account-permission layer ("no permission to buy kmr resource") — i.e. the
//     body shape was ACCEPTED — whereas the KES kes/2020-12-15 LaunchCluster
//     rejects every 8.x MainVersion with "'MainVersion' invalid".
//   - Read   -> kes/2020-12-15 DescribeCluster/ListClusters, the shipped and
//     tested KES query actions (the same ones the ksyun_kes_clusters data source
//     uses). KES clusters are one product line; ListClusters lists all of them
//     regardless of which creation path was used. Switching Read to kmr/2021-09-02
//     is untested and would risk an unknown response shape, so Read stays on the
//     known-good KES query path.
//
// The provider initialises a single KES client whose endpoint is the KES host
// (config.Endpoint, e.g. kes.inner.api.ksyun.com — see WithKesClient). The kmr
// /2021-09-02 Create action is served by the kmr host, so CreateVectorCluster
// MUST call req.SetDomain(kmrEndpointFor(...)); otherwise SendV2 builds the URL
// from the KES client's endpoint (the kes host), which does not serve the
// kmr/2021-09-02 action. The signing/transport channel (common.Client.SendV2
// embedded in *kes.Client) is shared, so the KES client can sign and dispatch a
// Service="kmr" request — only the URL host differs. The response body is parsed
// by hand from the raw response string returned by SendV2.
type KesVectorService struct {
	client *KsyunClient
}

// kmrEndpointFor derives the kmr service host from the provider's configured KES
// endpoint by replacing the first "kes" segment with "kmr"
// (kes.inner.api.ksyun.com -> kmr.inner.api.ksyun.com, kes.api.ksyun.com ->
// kmr.api.ksyun.com). An empty endpoint falls back to the inner kmr host, which is
// what the internal test account uses.
func kmrEndpointFor(endpoint string) string {
	if endpoint == "" {
		return "kmr.inner.api.ksyun.com"
	}
	if strings.Contains(endpoint, "kes") {
		return strings.Replace(endpoint, "kes", "kmr", 1)
	}
	return "kmr." + endpoint
}

// kesVectorSendResult carries the two non-error return values of SendV2 out of
// the WithKesClient closure (which can only return a single interface{}).
type kesVectorSendResult struct {
	statusCode int
	msg        string
}

// kesVectorSendV2 dispatches a raw request through the KES client's generic
// signing/transport channel and returns the HTTP status code together with the
// raw response body string. The caller is responsible for parsing the body and
// enforcing a 2xx status check.
func (s *KesVectorService) kesVectorSendV2(req ksyunhttp.Request) (int, string, error) {
	respInterface, err := s.client.WithKesClient(func(conn *kes.Client) (interface{}, error) {
		sc, m, e := conn.SendV2(req, &ksyunhttp.BaseResponse{})
		return kesVectorSendResult{statusCode: sc, msg: m}, e
	})
	if err != nil {
		return 0, "", err
	}
	r, ok := respInterface.(kesVectorSendResult)
	if !ok {
		return 0, "", fmt.Errorf("unexpected return type from KES SendV2")
	}
	return r.statusCode, r.msg, nil
}

// kesVectorStringSlice converts a TypeList of strings from the schema (returned
// by d.Get as []interface{}) into a []string for the API body. The project has no
// generic string-list expander, so this is local to KES8.
func kesVectorStringSlice(raw []interface{}) []string {
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}

// ----- LaunchCluster (kmr/2021-09-02 + ClusterType:KES) -----

// kesVectorLaunchClusterRequest mirrors the body shape captured from the Ksyun
// console when creating a KES 8.x cluster (path /kmr/v2/LaunchCluster,
// Version 2021-09-02, Service="kmr"). The decisive field is ClusterType="KES",
// which switches the shared kmr/2021-09-02 backend from KMR mode into KES8 mode.
//
// DatabaseInfo has NO omitempty: the production kmr/v2 gateway validates the
// request against the KMR LaunchCluster shape (which includes DatabaseInfo as a
// non-omitempty dict) BEFORE routing by ClusterType. An absent DatabaseInfo is
// rejected as "'DatabaseInfo' should be dictionary, got None". KES never
// consumes DatabaseInfo (gringotts only reads it for CozeStudio/Ksgc), so an
// all-zero dict is a pure gateway-satisfier.
//
// Per-group placement (AvailabilityZone, VpcSubnetId, VpcId) and the per-group
// index/ChargeType live INSIDE each InstanceGroups entry on the kmr/2021-09-02
// backend; VpcSubnetId is also sent at top level (both are present in the console
// capture). PurchaseTime and VpcEpcSubnetId are NOT part of the KES8 shape and
// are omitted. All pointer fields use omitempty so unset optionals are omitted
// from the JSON body.
//
// *ksyunhttp.BaseRequest is embedded solely to satisfy the ksyunhttp.Request
// interface; its fields are all unexported, so it does not pollute the
// marshalled JSON body.
type kesVectorLaunchClusterRequest struct {
	*ksyunhttp.BaseRequest
	ClusterType            *string                `json:"ClusterType,omitempty" name:"ClusterType"`
	Distribution           *string                `json:"Distribution,omitempty" name:"Distribution"`
	ClusterName            *string                `json:"ClusterName,omitempty" name:"ClusterName"`
	MainVersion            *string                `json:"MainVersion,omitempty" name:"MainVersion"`
	ProjectId              *int                   `json:"ProjectId,omitempty" name:"ProjectId"`
	DatabaseInfo           map[string]interface{} `json:"DatabaseInfo" name:"DatabaseInfo"`
	Scenario               *string                `json:"Scenario,omitempty" name:"Scenario"`
	VpcDomainId            *string                `json:"VpcDomainId,omitempty" name:"VpcDomainId"`
	VpcSubnetId            *string                `json:"VpcSubnetId,omitempty" name:"VpcSubnetId"`
	SecurityGroupId        *string                `json:"SecurityGroupId,omitempty" name:"SecurityGroupId"`
	ChargeType             *string                `json:"ChargeType,omitempty" name:"ChargeType"`
	InstanceLoginType      *string                `json:"InstanceLoginType,omitempty" name:"InstanceLoginType"`
	InstancePassword       *string                `json:"InstancePassword,omitempty" name:"InstancePassword"`
	EnableBestDeliveryMode *bool                  `json:"EnableBestDeliveryMode,omitempty" name:"EnableBestDeliveryMode"`
	Services               []string               `json:"Services,omitempty" name:"Services"`
	XpackEnabled           *bool                  `json:"XpackEnabled,omitempty" name:"XpackEnabled"`
	EnableEip              *bool                  `json:"EnableEip,omitempty" name:"EnableEip"`
	XpackUsername          *string                `json:"XpackUsername,omitempty" name:"XpackUsername"`
	XpackPassword          *string                `json:"XpackPassword,omitempty" name:"XpackPassword"`
	Plugins                []string               `json:"Plugins,omitempty" name:"Plugins"`
	InstanceGroups         []map[string]interface{} `json:"InstanceGroups,omitempty" name:"InstanceGroups"`
}

// kesVectorLaunchClusterResponse is the minimal shape parsed from the
// LaunchCluster response body; only ClusterId is consumed.
type kesVectorLaunchClusterResponse struct {
	ClusterId *string `json:"ClusterId"`
	RequestId *string `json:"RequestId"`
}

// CreateVectorCluster launches a KES8 (8.x/9.x) cluster via the shared
// kmr/2021-09-02 LaunchCluster backend with ClusterType="KES", then reads the
// result back into the resource state.
func (s *KesVectorService) CreateVectorCluster(d *schema.ResourceData, r *schema.Resource) (err error) {
	req := &kesVectorLaunchClusterRequest{BaseRequest: &ksyunhttp.BaseRequest{}}
	req.Init().WithApiInfo("kmr", "2021-09-02", "LaunchCluster")
	req.SetContext(context.Background())
	req.SetContentType("application/json")
	// The provider's single endpoint is the KES host; the kmr/2021-09-02 action
	// is served by the kmr host, so override the request domain explicitly.
	req.SetDomain(kmrEndpointFor(s.client.config.Endpoint))

	// ClusterType="KES" switches the shared kmr backend into KES8 mode.
	req.ClusterType = ptrString("KES")
	req.ClusterName = ptrString(d.Get("cluster_name").(string))
	req.MainVersion = ptrString(d.Get("main_version").(string))
	req.Distribution = ptrString(d.Get("distribution").(string))
	req.Scenario = ptrString(d.Get("scenario").(string))
	req.ChargeType = ptrString(d.Get("charge_type").(string))
	req.ProjectId = ptrInt(d.Get("project_id").(int))
	req.VpcDomainId = ptrString(d.Get("vpc_domain_id").(string))
	req.VpcSubnetId = ptrString(d.Get("vpc_subnet_id").(string))
	if v := d.Get("security_group_id").(string); v != "" {
		req.SecurityGroupId = ptrString(v)
	}
	req.XpackUsername = ptrString(d.Get("xpack_username").(string))
	req.XpackPassword = ptrString(d.Get("xpack_password").(string))

	// KES8-specific switches/credentials (defaults come from the schema).
	req.XpackEnabled = ptrBool(d.Get("xpack_enabled").(bool))
	req.EnableBestDeliveryMode = ptrBool(d.Get("enable_best_delivery_mode").(bool))
	req.EnableEip = ptrBool(d.Get("enable_eip").(bool))
	if v := d.Get("instance_login_type").(string); v != "" {
		req.InstanceLoginType = ptrString(v)
	}
	if v := d.Get("instance_password").(string); v != "" {
		req.InstancePassword = ptrString(v)
	}

	// Services defaults to ["elasticsearch"] when the user omits it; Plugins is
	// only sent when set (omitempty).
	services := kesVectorStringSlice(d.Get("services").([]interface{}))
	if len(services) == 0 {
		services = []string{"elasticsearch"}
	}
	req.Services = services
	if raw := d.Get("plugins").([]interface{}); len(raw) > 0 {
		req.Plugins = kesVectorStringSlice(raw)
	}

	// DatabaseInfo: all-zero placeholder. The kmr/v2 gateway type-checks this dict
	// BEFORE routing by ClusterType; KES never consumes it.
	req.DatabaseInfo = map[string]interface{}{
		"Type":          "",
		"RdsInstanceId": "",
		"Endpoint":      "",
		"Username":      "",
		"Password":      "",
		"DatabaseName":  "",
		"IsVerified":    false,
	}

	// Build instance groups from the schema list. On the kmr/2021-09-02 backend
	// each group carries its own placement (AvailabilityZone, VpcSubnetId, VpcId)
	// and index fields, plus the per-group ChargeType.
	igRaw := d.Get("instance_groups").([]interface{})
	igList := make([]map[string]interface{}, 0, len(igRaw))
	for _, raw := range igRaw {
		if raw == nil {
			continue
		}
		m := raw.(map[string]interface{})
		g := map[string]interface{}{}
		if v, ok := m["instance_group_type"].(string); ok && v != "" {
			g["InstanceGroupType"] = v
		}
		if v, ok := m["resource_type"].(string); ok && v != "" {
			g["ResourceType"] = v
		}
		if v, ok := m["instance_type_code"].(string); ok && v != "" {
			g["InstanceTypeCode"] = v
		}
		if v, ok := m["instance_type"].(string); ok && v != "" {
			g["InstanceType"] = v
		}
		g["InstanceCount"] = m["instance_count"].(int)
		if v, ok := m["vpc_id"].(string); ok && v != "" {
			g["VpcId"] = v
		}
		// Index fields are always sent (incl. 0) so the MASTER group (index 0) is
		// not silently dropped; the console capture includes them for every group.
		g["InstanceGroupIndex"] = m["instance_group_index"].(int)
		if v, ok := m["availability_zone"].(string); ok && v != "" {
			g["AvailabilityZone"] = v
		}
		g["AvailabilityZoneIndex"] = m["availability_zone_index"].(int)
		if v, ok := m["vpc_subnet_id"].(string); ok && v != "" {
			g["VpcSubnetId"] = v
		}
		if v, ok := m["volume_type"].(string); ok && v != "" {
			g["VolumeType"] = v
		}
		g["VolumeSize"] = m["volume_size"].(int)
		g["VolumeCount"] = m["volume_count"].(int)
		if v, ok := m["charge_type"].(string); ok && v != "" {
			g["ChargeType"] = v
		}
		if v, ok := m["multi_instance_count"].(int); ok && v != 0 {
			g["MultiInstanceCount"] = v
		}
		if v, ok := m["raid_type"].(string); ok && v != "" {
			g["RaidType"] = v
		}
		g["SystemDiskType"] = m["system_disk_type"].(string)
		g["SystemDiskSize"] = m["system_disk_size"].(int)
		igList = append(igList, g)
	}
	req.InstanceGroups = igList

	logger.Debug(logger.ReqFormat, "LaunchCluster", req)
	statusCode, msg, err := s.kesVectorSendV2(req)
	if err != nil {
		return fmt.Errorf("error on creating KES vector cluster: %s", err)
	}
	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("error on creating KES vector cluster: [HttpCode:%d] %s", statusCode, msg)
	}
	var resp kesVectorLaunchClusterResponse
	if err := json.Unmarshal([]byte(msg), &resp); err != nil {
		return fmt.Errorf("error on parsing KES vector cluster create response: %s, body: %s", err, msg)
	}
	if resp.ClusterId == nil || *resp.ClusterId == "" {
		return fmt.Errorf("KES vector cluster was created but no ClusterId was returned: %s", msg)
	}
	d.SetId(*resp.ClusterId)
	return s.ReadAndSetVectorClusterResource(d, r)
}

// ----- DescribeCluster / ListClusters (kes/2020-12-15, Read only) -----

// kesVectorDescribeClusterRequest queries a single cluster by id.
type kesVectorDescribeClusterRequest struct {
	*ksyunhttp.BaseRequest
	ClusterId *string `json:"ClusterId,omitempty" name:"ClusterId"`
}

// kesVectorDescribeClusterResponse mirrors the KES DescribeCluster response shape.
// It is a typed struct (rather than map[string]interface{}) so that JSON numbers
// decode into *int instead of float64, which the TypeInt schema fields require.
type kesVectorDescribeClusterResponse struct {
	ClusterId      *string `json:"ClusterId"`
	ClusterName    *string `json:"ClusterName"`
	ClusterType    *string `json:"ClusterType"`
	MainVersion    *string `json:"MainVersion"`
	InstanceGroups []struct {
		Id                *string `json:"Id"`
		InstanceGroupType *string `json:"InstanceGroupType"`
		ResourceType      *string `json:"ResourceType"`
		InstanceType      *string `json:"InstanceType"`
		VolumeSize        *int    `json:"VolumeSize"`
		VolumeType        *string `json:"VolumeType"`
	} `json:"InstanceGroups"`
	EnableEip        *bool   `json:"EnableEip"`
	Region           *string `json:"Region"`
	AvailabilityZone *string `json:"AvailabilityZone"`
	VpcDomainId      *string `json:"VpcDomainId"`
	VpcSubnetId      *string `json:"VpcSubnetId"`
	ChargeType       *string `json:"ChargeType"`
	ClusterStatus    *string `json:"ClusterStatus"`
	CreateTime       *string `json:"CreateTime"`
	UpdateTime       *string `json:"UpdateTime"`
	ServingMinutes   *int    `json:"ServingMinutes"`
	ProxyPort        *int    `json:"ProxyPort"`
	Tags             []struct {
		TagKey   *string `json:"TagKey"`
		TagValue *string `json:"TagValue"`
	} `json:"Tags"`
	RequestId  *string `json:"RequestId"`
	StatusCode *int    `json:"StatusCode"`
}

// kesVectorListClustersRequest lists clusters (used as a Read fallback).
type kesVectorListClustersRequest struct {
	*ksyunhttp.BaseRequest
	Marker *string `json:"Marker,omitempty" name:"Marker"`
}

// kesVectorListClustersResponse is the list response; ListClusters has no
// next-marker field, so a single (large) page is fetched and filtered client-side.
type kesVectorListClustersResponse struct {
	Clusters   []kesVectorListCluster `json:"Clusters"`
	Total      *int                   `json:"Total"`
	RequestId  *string                `json:"RequestId"`
	StatusCode *int                   `json:"StatusCode"`
}

// kesVectorListCluster is the lite cluster shape returned by ListClusters.
type kesVectorListCluster struct {
	ClusterId        *string `json:"ClusterId"`
	ClusterName      *string `json:"ClusterName"`
	MainVersion      *string `json:"MainVersion"`
	EnableEip        *bool   `json:"EnableEip"`
	Region           *string `json:"Region"`
	AvailabilityZone *string `json:"AvailabilityZone"`
	VpcDomainId      *string `json:"VpcDomainId"`
	VpcSubnetId      *string `json:"VpcSubnetId"`
	ChargeType       *string `json:"ChargeType"`
	ClusterStatus    *string `json:"ClusterStatus"`
	CreateTime       *string `json:"CreateTime"`
	UpdateTime       *string `json:"UpdateTime"`
	ServingMinutes   *int    `json:"ServingMinutes"`
}

// describeVectorCluster calls DescribeCluster v2020-12-15 (Service="kes") and
// returns a normalized map suitable for SdkResponseAutoResourceData.
func (s *KesVectorService) describeVectorCluster(clusterId string) (map[string]interface{}, error) {
	req := &kesVectorDescribeClusterRequest{BaseRequest: &ksyunhttp.BaseRequest{}}
	req.Init().WithApiInfo("kes", "2020-12-15", "DescribeCluster")
	req.SetContext(context.Background())
	req.SetContentType("application/json")
	req.ClusterId = ptrString(clusterId)

	logger.Debug(logger.ReqFormat, "DescribeCluster", req)
	statusCode, msg, err := s.kesVectorSendV2(req)
	if err != nil {
		return nil, err
	}
	if statusCode < 200 || statusCode > 299 {
		return nil, fmt.Errorf("[KsyunSDKError] [HttpCode:%d] %s", statusCode, msg)
	}
	var resp kesVectorDescribeClusterResponse
	if err := json.Unmarshal([]byte(msg), &resp); err != nil {
		return nil, fmt.Errorf("error on parsing KES vector cluster describe response: %s", err)
	}
	return normalizeVectorClusterMap(&resp, resp.Tags)
}

// listVectorClusterById calls ListClusters v2020-12-15 (Service="kes") and
// returns the cluster matching clusterId. It is the Read fallback when
// DescribeCluster is unavailable (e.g. the action is not supported for this
// account/region).
func (s *KesVectorService) listVectorClusterById(clusterId string) (map[string]interface{}, error) {
	req := &kesVectorListClustersRequest{BaseRequest: &ksyunhttp.BaseRequest{}}
	req.Init().WithApiInfo("kes", "2020-12-15", "ListClusters")
	req.SetContext(context.Background())
	req.SetContentType("application/json")
	req.Marker = ptrString("limit=100&offset=0")

	logger.Debug(logger.ReqFormat, "ListClusters", req)
	statusCode, msg, err := s.kesVectorSendV2(req)
	if err != nil {
		return nil, err
	}
	if statusCode < 200 || statusCode > 299 {
		return nil, fmt.Errorf("[KsyunSDKError] [HttpCode:%d] %s", statusCode, msg)
	}
	var resp kesVectorListClustersResponse
	if err := json.Unmarshal([]byte(msg), &resp); err != nil {
		return nil, fmt.Errorf("error on parsing KES vector cluster list response: %s", err)
	}
	for _, c := range resp.Clusters {
		if c.ClusterId != nil && *c.ClusterId == clusterId {
			return normalizeVectorClusterMap(&c, nil)
		}
	}
	return nil, fmt.Errorf("KES vector cluster %s not found", clusterId)
}

// readVectorClusterById fetches a single KES8 cluster. It prefers DescribeCluster;
// if that action is unavailable it falls back to ListClusters + client-side
// filter.
func (s *KesVectorService) readVectorClusterById(clusterId string) (map[string]interface{}, error) {
	m, err := s.describeVectorCluster(clusterId)
	if err == nil {
		return m, nil
	}
	if isActionUnavailableError(err) {
		return s.listVectorClusterById(clusterId)
	}
	return nil, err
}

// ReadAndSetVectorClusterResource reads a single KES8 cluster by id and writes
// it into the resource state. Missing clusters are removed from state.
func (s *KesVectorService) ReadAndSetVectorClusterResource(d *schema.ResourceData, r *schema.Resource) (err error) {
	m, err := s.readVectorClusterById(d.Id())
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading KES vector cluster %s: %s", d.Id(), err)
	}
	if len(m) == 0 {
		d.SetId("")
		return nil
	}
	SdkResponseAutoResourceData(d, r, m, nil)
	return nil
}

// DeleteVectorCluster is a no-op: the KES SDK does not provide a DeleteCluster
// API. The cluster must be deleted manually in the Ksyun console. terraform
// destroy only removes the resource from Terraform state.
func (s *KesVectorService) DeleteVectorCluster(d *schema.ResourceData, meta interface{}) (err error) {
	d.SetId("")
	return nil
}

// ----- helpers -----

// normalizeVectorClusterMap turns a typed cluster response struct into the
// map shape that SdkResponseAutoResourceData expects: top-level PascalCase keys
// (Hump2Downline converts them to the snake_case schema fields) with tags
// expanded into snake_case sub-maps. instance_groups is intentionally dropped
// because the API only returns a partial instance-group shape (it omits
// InstanceCount and the optional disk fields that are Required/ForceNew in the
// schema); overwriting it would zero out instance_count in state and trigger a
// perpetual diff. Since every instance_groups field is ForceNew, preserving the
// config-provided values is both safe and correct.
func normalizeVectorClusterMap(cluster interface{}, tags []struct {
	TagKey   *string `json:"TagKey"`
	TagValue *string `json:"TagValue"`
}) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := mapstructure.Decode(cluster, &m); err != nil {
		return nil, fmt.Errorf("error on decoding KES vector cluster response: %s", err)
	}
	derefPointerFields(m)
	// drop the partial/raw nested collections so config values are preserved.
	delete(m, "InstanceGroups")
	delete(m, "Tags")
	// expand tags into snake_case sub-maps consumed by the nested schema.
	tagInterfaces := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagMap := map[string]interface{}{}
		if tag.TagKey != nil {
			tagMap["tag_key"] = *tag.TagKey
		}
		if tag.TagValue != nil {
			tagMap["tag_value"] = *tag.TagValue
		}
		tagInterfaces = append(tagInterfaces, tagMap)
	}
	m["tags"] = tagInterfaces
	return m, nil
}

// isActionUnavailableError reports whether err indicates that the action itself
// is not supported on this API version (so a list-based fallback should be used).
// It must be checked before notFoundError, because notFoundError also matches
// the substring "invalid" (e.g. "InvalidAction").
func isActionUnavailableError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "invalidaction") ||
		strings.Contains(msg, "invalid action") ||
		strings.Contains(msg, "unknownaction") ||
		strings.Contains(msg, "unknown action") ||
		strings.Contains(msg, "nosuchaction") ||
		strings.Contains(msg, "no such action") ||
		strings.Contains(msg, "actionnotfound") ||
		strings.Contains(msg, "action not found") ||
		strings.Contains(msg, "nosuchversion") ||
		strings.Contains(msg, "no such version") ||
		strings.Contains(msg, "unsupportedaction") ||
		strings.Contains(msg, "unsupported action") ||
		strings.Contains(msg, "notsupport")
}
