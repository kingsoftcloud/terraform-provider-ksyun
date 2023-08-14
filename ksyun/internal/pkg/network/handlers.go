package network

import (
	"errors"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
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
		if orig, ok := nestErr.(*net.OpError); ok {
			dealErr := &net.OpError{
				Op:     orig.Op,
				Addr:   orig.Addr,
				Net:    orig.Net,
				Source: orig.Source,
				Err:    errors.New(fmt.Sprintf("%s, 您的网络似乎不太稳定.", orig.Err.Error())),
			}
			newErr := awserr.New(err.Code(), err.Message(), dealErr)
			return newErr
		}
	}

	return origErr
}
