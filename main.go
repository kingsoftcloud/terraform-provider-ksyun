package main

import (
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun"
)

var (
	version string
)

func main() {
	ksc.SDKName = "terraform-provider-ksyun"
	ksc.SDKVersion = version
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ksyun.Provider,
	})
}
