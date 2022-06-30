provider "ksyun" {
}

resource "ksyun_tag" "test_tag" {
    key = "test_tag_key"
    value = "test_tag_value"
    resource_type = "eip"
    resource_id = 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
}

