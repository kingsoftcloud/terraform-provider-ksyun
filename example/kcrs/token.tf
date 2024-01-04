resource "ksyun_kcrs_token" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	token_type = "NeverExpire"
	token_time = 11
	desc = "test23"
	enable = true
}