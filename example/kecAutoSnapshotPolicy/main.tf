resource "ksyun_auto_snapshot_policy" "foo" {
  name   = "tf_combine_on_hcl_rename_1"
  auto_snapshot_date = [1,3,5]
  auto_snapshot_time = [1,3]
  retention_time = 3
}