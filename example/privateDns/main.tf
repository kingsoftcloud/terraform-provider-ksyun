provider "ksyun" {
  region = "cn-guangzhou-1"
}


resource "ksyun_private_dns_zone" "foo" {
  zone_name = "tf-pdns-zone-pdns.com"
  zone_ttl = 3200
  charge_type = "TrafficMonthly"
}

resource "ksyun_private_dns_record" "foo" {
  record_name = "tf-pdns-record"
  record_ttl = 3600
  zone_id = ksyun_private_dns_zone.foo.id
  type = "SRV"
  record_value = "tf-record.com"
  port = 2222
  priority = 102
  weight = 2222
}


resource "ksyun_vpc" "foo" {
  vpc_name      = "tf-pdns-binding-vpc"
  cidr_block    = "10.7.0.0/21"
}

resource "ksyun_private_dns_zone_vpc_attachment" "example" {
  zone_id = ksyun_private_dns_zone.foo.id
  vpc_set {
    region_name = "cn-guangzhou-1"
    vpc_id = ksyun_vpc.foo.id
  }
}