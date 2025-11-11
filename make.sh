#bin/sh
go build
rm $GOPATH/bin/terraform-provider-ksyun
cp $GOPATH/src/github.com/kingsoftcloud/terraform-provider-ksyun/terraform-provider-ksyun $GOPATH/bin/