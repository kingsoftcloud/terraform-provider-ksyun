package network

import (
	"net"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/hashicorp/go-uuid"
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

// ReqUniqueIdHandler deal with add trace id
var ReqUniqueIdHandler = request.NamedHandler{
	Name: "ksyun.ReqUniqueIdHandler",
	Fn: func(r *request.Request) {
		httpReq := r.HTTPRequest
		body, _ := url.ParseQuery(httpReq.URL.RawQuery)
		body.Set("TraceId", GenerateTraceId())
		r.HTTPRequest.URL.RawQuery = body.Encode()
	},
}

var DebugTraceError = request.NamedHandler{
	Name: "ksyun.DebugTraceError",
	Fn: func(r *request.Request) {
		if r.Error == nil {
			return
		}

		httpReq := r.HTTPRequest
		body, _ := url.ParseQuery(httpReq.URL.RawQuery)

		if body.Has("TraceId") {
			switch err := r.Error.(type) {
			case awserr.Error:
				traceId := body.Get("TraceId")
				r.Error = awserr.NewBatchError(err.Code(), "Batch error caused by following:", []error{err, infraerrs.AssembleTraceIdError(traceId)})
			}
		}

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

func GenerateTraceId() string {
	uid, err := uuid.GenerateUUID()
	if err != nil {
		return ""
	}
	return uid
}
