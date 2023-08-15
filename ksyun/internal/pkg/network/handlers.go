package network

import (
	"net"
	"net/url"

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
		if err.Code() == "RequestError" && isNetworkError(origErr) {
			newErr := awserr.New(err.Code(), infraerrs.GetKsyunNetworkOpErrorMessage(err.Message()), err.OrigErr())
			return newErr
		}
	}

	return origErr
}

func isNetworkError(origErr error) bool {
	switch err := origErr.(type) {
	case awserr.Error:
		return isNetworkError(err.OrigErr())

	case *url.Error:
		return isNetworkError(err.Err)
	case temporary:
		if _, ok := err.(*net.OpError); ok {
			return true
		}

	}
	return false
}
