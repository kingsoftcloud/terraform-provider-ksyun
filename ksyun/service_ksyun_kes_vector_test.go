package ksyun

import "testing"

// KES8 is implemented by reusing the KMR framework: the console POSTs to the
// shared kmr/2021-09-02 LaunchCluster action with X-Ksc-Source: kes, which
// makes Oasis resolve product='kes' and, because the source is an inner-api
// source, skip complete_params. The body's own ClusterType="KES" then reaches
// gringotts unchanged and ListClusterVersions is queried with
// ClusterType='KES'.
//
// Without the header Oasis falls into the openapi 'user' branch, rewrites
// cluster_type to 'KMR' and queries KMR cluster versions, which yields
// "main_version must be in [] for ElasticSearch". A previous commit deleted the
// header believing it was a meaningless decoration, so this test pins it.

func TestKesVectorKmrHeadersSource(t *testing.T) {
	headers := kesVectorKmrHeaders()
	got, ok := headers["X-Ksc-Source"]
	if !ok {
		t.Fatalf(`kesVectorKmrHeaders() is missing X-Ksc-Source; without it Oasis
rewrites cluster_type to KMR and LaunchCluster fails with
"main_version must be in [] for ElasticSearch"`)
	}
	// "kmr" would resolve PRODUCT_GROUP_ID_MAP['KMR']=107 and reject KES
	// instance types (ES.I3.2C4G, ...) with "'InstanceType' not found".
	if got != "kes" {
		t.Errorf(`X-Ksc-Source = %q, want "kes"
(Oasis resolves PRODUCT_GROUP_ID_MAP[self.product.upper()] to validate
InstanceType, and KES instance types only exist in the KES product group)`, got)
	}
}

// The kmr/2021-09-02 action is served by the kmr host while the KES client is
// pinned to the KES host (config.go WithKesClient), so the domain override is
// load-bearing. Mirror the endpoints used by the two Oasis deployments.
func TestKmrEndpointFor(t *testing.T) {
	cases := map[string]string{
		"":                            "kmr.inner.api.ksyun.com",
		"kes.inner.api.ksyun.com":     "kmr.inner.api.ksyun.com",
		"kes.api.ksyun.com":           "kmr.api.ksyun.com",
		"kes.inner.api.ksyun.com:443": "kmr.inner.api.ksyun.com:443",
		// An already-kmr endpoint must stay reachable rather than double-prefixing.
		"kmr.inner.api.ksyun.com": "kmr.kmr.inner.api.ksyun.com",
	}
	for in, want := range cases {
		if got := kmrEndpointFor(in); got != want {
			t.Errorf("kmrEndpointFor(%q) = %q, want %q", in, got, want)
		}
	}
}
