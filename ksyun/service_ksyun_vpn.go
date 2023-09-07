package ksyun

var (
	// vpnV1Attribute the flowing fields are invalid when vpn1.0
	vpnV2Attribute = []string{"ha_mode", "open_health_check", "local_peer_ip", "customer_peer_ip", "ike_version"}

	// vpnV1Attribute the flowing fields are invalid when vpn2.0
	vpnV1Attribute     = []string{"vpn_gre_ip", "ha_vpn_gre_ip", "customer_gre_ip", "ha_customer_gre_ip"}
	vpnCommonAttribute = []string{""}
)
