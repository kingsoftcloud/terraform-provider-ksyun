#!/bin/bash

set -e

version=$1

# Detech current os category
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     OS_TYPE=linux;;
    Darwin*)    OS_TYPE=darwin;;
    CYGWIN*)    OS_TYPE=windows;;
    MINGW*)     OS_TYPE=windows;;
    *)          OS_TYPE="UNKNOWN:${unameOut}"
esac

unameOut="$(uname -m)"
case "${unameOut}" in
    x86_64*)     OS_ARCH=amd64;;
    arm64*)      OS_ARCH=arm64;;
    *)          OS_ARCH="UNKNOWN:${unameOut}"
esac

echo "OS ${OS_TYPE}-${OS_ARCH} is deteched."
echo "Compiling ..."

plugin_path=$HOME/.terraform.d/plugin-cache/registry.terraform.io/kingsoftcloud/ksyun/${version}/${OS_TYPE}_${OS_ARCH}
plugin_path_win=$HO$APPDATAME/terraform.d/plugin-cache/registry.terraform.io/kingsoftcloud/ksyun/${version}/${OS_TYPE}_${OS_ARCH}

echo $plugin_path

if [ $OS_TYPE == "linux" -o $OS_TYPE == "darwin" ]; then
	GOOS=$OS_TYPE GOARCH=$OS_ARCH go build -o bin/terraform-provider-ksyun
	chmod +x bin/terraform-provider-ksyun
    mkdir -p $plugin_path
    mv bin/terraform-provider-ksyun $plugin_path/terraform-provider-ksyun_v${version}
elif [ $OS_TYPE == "Windows" ]; then
	GOOS=$OS_TYPE GOARCH=$OS_ARCH go build -o bin/terraform-provider-ksyun
	chmod +x bin/terraform-provider-ksyun.exe
#    mkdir -p $APPDATA/terraform.d/plugins
    mkdir -p $plugin_path_win
    mv bin/terraform-provider-ksyun.exe $plugin_path_win/terraform-provider-ksyun_v${version}.exe
else
    echo "Invalid OS"
    exit 1
fi

echo "Installation of ksyun Terraform Provider is completed."
