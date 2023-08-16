package network

import (
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/infraerrs"
)

const (
	NotFound                = "Notfound"
	ResourceNotfound        = "ResourceNotfound"
	InternalError           = "INTERNAL_FAILURE"
	ServiceTimeout          = "ServiceTimeout"
	CanceledErrorCode       = "RequestCanceled"
	RequestTimeout          = "RequestTimeout"
	ErrCodeResponseTimeout  = "ErrCodeResponseTimeout"
	RequestTimeoutException = "RequestTimeoutException"
)

// retryableErrorCodes is a list of retryable error code
var retryableErrorCodes = []string{ServiceTimeout, RequestTimeout, ErrCodeResponseTimeout, RequestTimeoutException}

type temporary interface {
	Temporary() bool
}

// custom retry
type KsyunRetryer struct {
	// ...
	NumMaxRetries int
}

var _ request.Retryer = (*KsyunRetryer)(nil)

func GetKsyunRetryer(maxRetries int) request.Retryer {
	return &KsyunRetryer{
		NumMaxRetries: maxRetries,
	}
}

func (k *KsyunRetryer) RetryRules(r *request.Request) time.Duration {
	// retry delay
	return 500 * time.Millisecond
}

func (k *KsyunRetryer) ShouldRetry(r *request.Request) bool {
	// indicates whether retry the request

	// ShouldRetry returns false if number of max retries is 0.
	if k.NumMaxRetries == 0 {
		return false
	}

	// If one of the other handlers already set the retry state
	// we don't want to override it based on the service's state
	if r.Retryable != nil {
		return *r.Retryable
	}

	if isErrCode(r.Error, r.RetryErrorCodes) {
		return true
	}

	// customs retry condition
	return shouldRetryError(r.Error) || isErrConnectionReset(r.Error)
}
func (k *KsyunRetryer) MaxRetries() int {
	return k.NumMaxRetries
}

func isErrConnectionReset(err error) bool {
	if strings.Contains(err.Error(), "read: connection reset") {
		return false
	}
	if strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "broken pipe") {
		return true
	}

	return false
}

func shouldRetryError(origErr error) bool {
	switch err := origErr.(type) {
	case awserr.Error:
		if err.Code() == CanceledErrorCode {
			return false
		}
		var shouldRetry bool
		origErr := err.OrigErr()
		if origErr != nil {
			shouldRetry = shouldRetryError(origErr)
			if err.Code() == "RequestError" && !shouldRetry {
				return false
			}
		}
		if infraerrs.IsExpectError(err, retryableErrorCodes) {
			return true
		}

		return shouldRetry
	case *url.Error:
		if strings.Contains(err.Error(), "connection refused") {
			// Refused connections should be retried as the service may not yet
			// be running on the port. Go TCP dial considers refused
			// connections as not temporary.
			return true
		}
		// *url.Error only implements Temporary after golang 1.6 but since
		// url.Error only wraps the error:
		return shouldRetryError(err.Err)
	case temporary:
		if netErr, ok := err.(*net.OpError); ok && netErr.Op == "dial" {
			return true
		}
		// If the error is temporary, we want to allow continuation of the
		// retry process
		return err.Temporary() || isErrConnectionReset(origErr)

	case nil:
		// `awserr.Error.OrigErr()` can be nil, meaning there was an error but
		// because we don't know the cause, it is marked as non-retryable.
		return false

	default:
		switch err.Error() {
		case "net/http: request canceled",
			"net/http: request canceled while waiting for connection":
			// known 1.5 error case when an http request is cancelled
			return false
		}
		// here we don't know the error; so we deny a retry.
		return false
	}
}

func isErrCode(err error, codes []string) bool {
	if aerr, ok := err.(awserr.Error); ok && aerr != nil {
		for _, code := range codes {
			if code == aerr.Code() {
				return true
			}
		}
	}

	return false
}
