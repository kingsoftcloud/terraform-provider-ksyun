package network

import (
	"net"
	"net/url"
	"os"
	"syscall"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
)

func TestShouldRetry(t *testing.T) {

	syscallError := os.SyscallError{
		Err:     syscall.ECONNRESET,
		Syscall: "read",
	}

	opError := net.OpError{
		Op:     "dial",
		Net:    "tcp",
		Source: net.Addr(nil),
		Err:    &syscallError,
	}

	urlError := url.Error{
		Op:  "Post",
		URL: "https://localhost:52398",
		Err: &opError,
	}
	origError := awserr.New("RequestError", "send request failed", &urlError)
	if a, e := shouldRetryError(origError), request.IsErrorRetryable(origError); e != a {
		t.Errorf("Expected to return %v to retry when error occured, got %v instead", a, e)
	}

}
