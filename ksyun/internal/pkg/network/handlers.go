package network

import (
	"net"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/infraerrs"
)

// NetErrorHandler will deal with network error, returns a ksyun custom error.
var NetErrorHandler = request.NamedHandler{
	Name: "ksyun.NetErrorHandler",
	Fn: func(r *request.Request) {
		if r.Error == nil {
			return
		}
		r.Error = handeNetError(r.Error)
	},
}

func handeNetError(origErr error) error {
	switch err := origErr.(type) {
	case awserr.Error:
		nestErr := err.OrigErr()
		if orig, ok := nestErr.(*net.OpError); ok && err.Code() == "RequestError" {
			newErr := awserr.New(err.Code(), infraerrs.GetKsyunNetworkOpErrorMessage(err.Message()), orig)
			return newErr
		}
	}

	return origErr
}
