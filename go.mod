module github.com/terraform-providers/terraform-provider-ksyun

go 1.12

require (
	github.com/KscSDK/ksc-sdk-go v0.1.43
	github.com/aws/aws-sdk-go v1.25.3
	github.com/client9/misspell v0.3.4
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/golangci/golangci-lint v1.23.7
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/ks3sdklib/aws-sdk-go v1.0.4
	github.com/pkg/errors v0.9.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/mod v0.2.0 // indirect
)

//replace github.com/KscSDK/ksc-sdk-go => ../../KscSDK/ksc-sdk-go
