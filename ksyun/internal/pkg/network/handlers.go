package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strings"

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

var OutputResetError = request.NamedHandler{
	Name: "ksyun.OutputResetError",
	Fn: func(r *request.Request) {
		if r.Error == nil {
			return
		}
		if isErrConnectionReset(r.Error) {
			fmt.Printf("Request: %v", r)
		}
	},
}

var HandleRequestBody = request.NamedHandler{
	Name: "ksyun.HandleRequestBody",
	Fn: func(r *request.Request) {
		if r.HTTPRequest.Method == "GET" {
			return
		}
		body := r.Body

		// r.HTTPRequest.Body = &CustomReader{body}
		r.HTTPRequest.Body = ioutil.NopCloser(body)

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

func IsReadConnectionReset(origErr error) bool {
	switch err := origErr.(type) {
	case awserr.Error:
		return IsReadConnectionReset(err.OrigErr())

	case *url.Error:
		return IsReadConnectionReset(err.Err)
	case temporary:
		if netErr, ok := err.(*net.OpError); ok &&
			strings.Contains(netErr.Error(), "read: connection reset") {
			return true
		}

	}
	return false
}
