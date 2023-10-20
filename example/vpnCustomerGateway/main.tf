resource "ksyun_vpn_customer_gateway" "default" {
  customer_gateway_address   = "100.0.0.2"
  ha_customer_gateway_address = "100.0.2.2"
  customer_gateway_name = "ksyun_vpn_cus_gw"
}