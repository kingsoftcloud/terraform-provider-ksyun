resource "ksyun_kcrs_webhook_trigger" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	namespace = ksyun_kcrs_namespace.foo.namespace
	trigger {
		trigger_url = "http://www.test11.com"
		trigger_name = "tfunittest22"
		event_types = ["PushImage"]
		headers {
			key = "pp1"
			value = "22222"
		}
		headers {
			key = "pp1"
			value = "333"
		}
		headers {
			key = "pp122"
			value = "32223"
		}
	}
}