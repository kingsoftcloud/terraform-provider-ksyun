package ksyun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// The KES billing layer (oasis/utils/validation/syntax.py) accepts exactly
// Monthly, Daily and HourlyInstantSettlement; Oasis rewrites
// HourlyInstantSettlement to Minutely before forwarding to billing. Anything
// else is rejected with `'ChargeType' invalid, got <value>.`, so the provider
// must refuse it during plan instead of failing at apply time.

func assertChargeTypeValidator(t *testing.T, label string, s *schema.Schema) {
	t.Helper()
	if s == nil {
		t.Fatalf("%s: charge_type schema is missing", label)
	}
	if s.ValidateFunc == nil {
		t.Fatalf("%s: charge_type has no ValidateFunc", label)
	}
	for _, ok := range []string{"Monthly", "Daily", "HourlyInstantSettlement"} {
		if _, errs := s.ValidateFunc(ok, "charge_type"); len(errs) != 0 {
			t.Errorf("%s: charge_type %q should be accepted, got %v", label, ok, errs)
		}
	}
	for _, bad := range []string{"Minutely", "Hourly", "Yearly", "PrePaid", "PostPaid", "monthly", ""} {
		if _, errs := s.ValidateFunc(bad, "charge_type"); len(errs) == 0 {
			t.Errorf("%s: charge_type %q should be rejected", label, bad)
		}
	}
}

func TestKesClusterChargeTypeValidation(t *testing.T) {
	r := resourceKsyunKesCluster()
	assertChargeTypeValidator(t, "ksyun_kes_cluster", r.Schema["charge_type"])

	purchaseTime := r.Schema["purchase_time"]
	if purchaseTime == nil {
		t.Fatalf("ksyun_kes_cluster: purchase_time schema is missing")
	}
	if purchaseTime.Required {
		t.Errorf("ksyun_kes_cluster: purchase_time must not be Required; the backend accepts an absent value but rejects an explicit 0")
	}
	if !purchaseTime.Optional {
		t.Errorf("ksyun_kes_cluster: purchase_time must be Optional")
	}
}

func TestKesClusterVectorChargeTypeValidation(t *testing.T) {
	r := resourceKsyunKesClusterVector()
	assertChargeTypeValidator(t, "ksyun_kes_cluster_vector", r.Schema["charge_type"])

	ig := r.Schema["instance_groups"]
	if ig == nil {
		t.Fatalf("ksyun_kes_cluster_vector: instance_groups schema is missing")
	}
	igResource, ok := ig.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("ksyun_kes_cluster_vector: instance_groups is not a nested resource")
	}
	assertChargeTypeValidator(t, "ksyun_kes_cluster_vector.instance_groups", igResource.Schema["charge_type"])
}

// TestKesResourcesInternalValidate catches schema-level mistakes (Required and
// Optional both set, invalid Default type, ...) without contacting the API.
func TestKesResourcesInternalValidate(t *testing.T) {
	for name, r := range map[string]*schema.Resource{
		"ksyun_kes_cluster":        resourceKsyunKesCluster(),
		"ksyun_kes_cluster_vector": resourceKsyunKesClusterVector(),
	} {
		if err := r.InternalValidate(nil, true); err != nil {
			t.Errorf("%s: InternalValidate failed: %v", name, err)
		}
	}
}

// The KES API returns MainVersion with an "OpenSource-" prefix while the
// documentation and every user configuration write the bare version. Without
// normalisation main_version drifts on every refresh and, being ForceNew, forces
// a destroy+recreate of the cluster on every plan.
func TestKesMainVersionStateFunc(t *testing.T) {
	cases := map[string]string{
		"7.4.2":              "7.4.2",
		"OpenSource-7.4.2":   "7.4.2",
		"8.19.13":            "8.19.13",
		"OpenSource-8.19.13": "8.19.13",
		"":                   "",
		"OpenSource-":        "",
	}
	for in, want := range cases {
		if got := mainVersionStateFunc(in); got != want {
			t.Errorf("mainVersionStateFunc(%q) = %q, want %q", in, got, want)
		}
	}

	// Both resources must wire the normaliser in; a missing StateFunc brings the
	// permanent ForceNew diff back.
	for name, r := range map[string]*schema.Resource{
		"ksyun_kes_cluster":        resourceKsyunKesCluster(),
		"ksyun_kes_cluster_vector": resourceKsyunKesClusterVector(),
	} {
		s := r.Schema["main_version"]
		if s == nil {
			t.Fatalf("%s: main_version schema is missing", name)
		}
		if s.StateFunc == nil {
			t.Fatalf("%s: main_version has no StateFunc", name)
		}
		if got := s.StateFunc("OpenSource-7.4.2"); got != "7.4.2" {
			t.Errorf("%s: main_version StateFunc returned %q, want %q", name, got, "7.4.2")
		}
		if !s.ForceNew {
			t.Errorf("%s: main_version is expected to stay ForceNew", name)
		}
	}
}
