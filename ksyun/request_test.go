package ksyun_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/request"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
	"github.com/terraform-providers/terraform-provider-ksyun/ksyun/internal/pkg/network"
)

type connResetCloser struct {
	Err error
}

func (rc *connResetCloser) Read(b []byte) (int, error) {
	return 0, rc.Err
}

func (rc *connResetCloser) Close() error {
	return nil
}

type tempNetworkError struct {
	op     string
	msg    string
	isTemp bool
}

func (e *tempNetworkError) Temporary() bool { return e.isTemp }
func (e *tempNetworkError) Error() string {
	return fmt.Sprintf("%s: %s", e.op, e.msg)
}

var (
	// net.OpError accept, are always temporary
	errAcceptConnectionResetStub = &tempNetworkError{
		isTemp: true, op: "accept", msg: "connection reset",
	}

	// net.OpError read for ECONNRESET is not temporary.
	errReadConnectionResetStub = &tempNetworkError{
		isTemp: false, op: "read", msg: "connection reset",
	}

	// net.OpError write for ECONNRESET may not be temporary, but is treaded as
	// temporary by the SDK.
	errWriteConnectionResetStub = &tempNetworkError{
		isTemp: false, op: "write", msg: "connection reset",
	}

	// net.OpError write for broken pipe may not be temporary, but is treaded as
	// temporary by the SDK.
	errWriteBrokenPipeStub = &tempNetworkError{
		isTemp: false, op: "write", msg: "broken pipe",
	}

	// Generic connection reset error
	errConnectionResetStub = errors.New("connection reset")
)

type mockAddr struct {
	addr string
}

func (a mockAddr) Network() string {
	return "tcp"
}
func (a mockAddr) String() string {
	return a.addr
}

var (
	errNetOpErrorResetStub = &net.OpError{
		Op:  "read",
		Net: "tcp",
		Addr: mockAddr{
			addr: "120.92.15.230:80",
		},
		Source: mockAddr{
			addr: "127.0.0.1:23333",
		},
		Err: errors.New("read: connection reset by peer"),
	}

	errNetOpErrorDialStub = &net.OpError{
		Op:  "dial",
		Net: "tcp",
		Addr: mockAddr{
			addr: "120.92.15.230:80",
		},
		Source: mockAddr{
			addr: "127.0.0.1:23333",
		},
		Err: errors.New("connection reset"),
	}

	errNetOpErrorBrokenPipeStub = &net.OpError{
		Op:  "write",
		Net: "tcp",
		Addr: mockAddr{
			addr: "120.92.15.230:80",
		},
		Source: mockAddr{
			addr: "127.0.0.1:23333",
		},
		Err: errors.New("broken pipe"),
	}

	errUrlError = &url.Error{
		Op:  "Get",
		URL: "http://krds.api.ksyun.com/?Action=DescribeDBParameterGroup&Version=2016-07-01",
		Err: errNetOpErrorResetStub,
	}
)

func TestConnectionReset(t *testing.T) {
	cases := map[string]struct {
		Err            error
		ExpectAttempts int
	}{
		"accept with temporary": {
			Err:            errAcceptConnectionResetStub,
			ExpectAttempts: 6,
		},
		"read not temporary": {
			Err:            errReadConnectionResetStub,
			ExpectAttempts: 1,
		},
		"write with temporary": {
			Err:            errWriteConnectionResetStub,
			ExpectAttempts: 6,
		},
		"write broken pipe with temporary": {
			Err:            errWriteBrokenPipeStub,
			ExpectAttempts: 6,
		},
		"generic connection reset": {
			Err:            errConnectionResetStub,
			ExpectAttempts: 6,
		},
		"write broken pipe with OpError": {
			Err:            errNetOpErrorBrokenPipeStub,
			ExpectAttempts: 6,
		},
		"read with OpError": {
			Err: errNetOpErrorResetStub,

			ExpectAttempts: 1,
		},
		"dial with OpError": {
			Err:            errNetOpErrorDialStub,
			ExpectAttempts: 6,
		},
		"read with UrlError": {
			Err:            errUrlError,
			ExpectAttempts: 1,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			count := 0
			handlers := request.Handlers{}
			handlers.Send.PushBack(func(r *request.Request) {
				count++
				r.HTTPResponse = &http.Response{}
				r.HTTPResponse.Body = &connResetCloser{
					Err: c.Err,
				}
				r.Error = awserr.New("RequestError", "send request failed", c.Err)
			})

			handlers.Sign.PushBackNamed(v4.SignRequestHandler)
			handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
			handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
			handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
			handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)
			handlers.AfterRetry.PushBackNamed(corehandlers.AfterRetryHandler)
			handlers.CompleteAttempt.PushBackNamed(network.NetErrorHandler)

			cfg := unit.Session.Config.Copy()
			cfg.MaxRetries = aws.Int(5)
			cfg.SleepDelay = func(time.Duration) {}

			op := &request.Operation{
				Name:       "op",
				HTTPMethod: "POST",
				HTTPPath:   "/",
			}

			meta := metadata.ClientInfo{
				ServiceName:   "fooService",
				SigningName:   "foo",
				SigningRegion: "foo",
				Endpoint:      "localhost",
				APIVersion:    "2001-01-01",
				JSONVersion:   "1.1",
				TargetPrefix:  "Foo",
			}

			req := request.New(
				*cfg,
				meta,
				handlers,
				network.GetKsyunRetryer(*cfg.MaxRetries),
				op,
				&struct{}{},
				&struct{}{},
			)

			osErr := c.Err
			req.ApplyOptions(request.WithResponseReadTimeout(time.Second))
			err := req.Send()
			if err == nil {
				t.Error("Expected error 'RequestError', but received nil")
			}
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "RequestError" {
				t.Errorf("Expected 'RequestError', but received %q", aerr.Code())
			} else if !ok {
				t.Errorf("Expected 'awserr.Error', but received %v", reflect.TypeOf(err))
			} else if aerr.OrigErr().Error() != osErr.Error() {
				t.Logf("Raw OpError %q, Handled OpError %q", osErr.Error(), aerr.OrigErr().Error())
			}

			t.Logf(err.Error())
			if e, a := c.ExpectAttempts, count; e != a {
				t.Errorf("Expected %v, but received %v", e, a)
			}

		})
	}
}
