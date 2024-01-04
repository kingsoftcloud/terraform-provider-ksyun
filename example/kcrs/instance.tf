resource "ksyun_kcrs_instance" "foo" {
	instance_name = "tfacceptancetest"
	instance_type = "basic"
	open_public_operation = true
	delete_bucket = true


	# open public access with external policy that permits an address, ip or cidr, to access this repository
	external_policy {
		entry = "192.168.2.133"
		desc = "d22"
	}
	external_policy {
		entry = "192.168.2.123/32"
		desc = "ddd"
	}
}