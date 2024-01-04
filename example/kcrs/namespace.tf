# Create a namespace under the repository instance
resource "ksyun_kcrs_namespace" "foo" {
	instance_id = ksyun_kcrs_instance.foo.id
	namespace = "tftest"
	public = true
}

# Create a namespace under the repository instance
resource "ksyun_kcrs_namespace" "foo1" {
	instance_id = ksyun_kcrs_instance.foo.id
	namespace = "tftest2"
	public = true
}