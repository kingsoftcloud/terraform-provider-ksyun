package ksyun

import (
	"encoding/json"
	"fmt"
	"github.com/ks3sdklib/ksyun-ks3-go-sdk/ks3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LifecycleRuleStatus string

const (
	ExpirationStatusEnabled  = LifecycleRuleStatus("Enabled")
	ExpirationStatusDisabled = LifecycleRuleStatus("Disabled")
)

type ErrorSource string

const (
	KsyunKs3GoSdk = ErrorSource("[SDK ksyun-ks3-go-sdk ERROR]")
	ProviderERROR = ErrorSource("[Provider ERROR]")
)

const Iso8601DateFormat = "2006-01-02T00:00:00+08:00"

type Status string

const (
	Deleted = Status("Deleted")
)

func ks3NotFoundError(err error) bool {
	if e, ok := err.(ks3.ServiceError); ok &&
		(e.StatusCode == 404 || strings.HasPrefix(e.Code, "NoSuch") || strings.HasPrefix(e.Message, "No Row found")) {
		return true
	}
	return false
}

type ListenerErr struct {
	ErrType string
	Err     error
}

func (e *ListenerErr) Error() string {
	return e.ErrType + " " + e.Err.Error()

}

// ComplexError is a format error which including origin error, extra error message, error occurred file and line
// Cause: a error is a origin error that comes from SDK, some exceptions and so on
// Err: a new error is built from extra message
// Path: the file path of error occurred
// Line: the file line of error occurred
type ComplexError struct {
	Cause error
	Err   error
	Path  string
	Line  int
}

func (e ComplexError) Error() string {
	if e.Cause == nil {
		e.Cause = Error("<nil cause>")
	}
	if e.Err == nil {
		return fmt.Sprintf("\u001B[31m[ERROR]\u001B[0m %s:%d:\n%s", e.Path, e.Line, e.Cause.Error())
	}
	return fmt.Sprintf("\u001B[31m[ERROR]\u001B[0m %s:%d: %s:\n%s", e.Path, e.Line, e.Err.Error(), e.Cause.Error())
}

func Error(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func IsExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*ComplexError); ok {
		return IsExpectedErrors(e.Cause, expectCodes)
	}

	if e, ok := err.(*ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}

	if e, ok := err.(ks3.ServiceError); ok {
		for _, code := range expectCodes {
			if e.Code == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}

	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

// Return a ComplexError which including error occurred file and path
func WrapError(cause error) error {
	if cause == nil {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("\u001B[31m[ERROR]\u001B[0m runtime.Caller error in WrapError.")
		return WrapComplexError(cause, nil, "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	return WrapComplexError(cause, nil, filepath, line)
}

// Return a ComplexError which including extra error message, error occurred file and path
func WrapErrorf(cause error, msg string, args ...interface{}) error {
	if cause == nil && strings.TrimSpace(msg) == "" {
		return nil
	}
	_, filepath, line, ok := runtime.Caller(1)
	if !ok {
		log.Printf("\u001B[31m[ERROR]\u001B[0m runtime.Caller error in WrapErrorf.")
		return WrapComplexError(cause, Error(msg), "", -1)
	}
	parts := strings.Split(filepath, "/")
	if len(parts) > 3 {
		filepath = strings.Join(parts[len(parts)-3:], "/")
	}
	// The second parameter of args is requestId, if the error message is NotFoundMsg the requestId need to be returned.
	if msg == NotFoundMsg && len(args) == 2 {
		msg += RequestIdMsg
	}
	return WrapComplexError(cause, fmt.Errorf(msg, args...), filepath, line)
}

func WrapComplexError(cause, err error, filepath string, fileline int) error {
	return &ComplexError{
		Cause: cause,
		Err:   err,
		Path:  filepath,
		Line:  fileline,
	}
}

// A default message of ComplexError's Err. It is format to Resource <resource-id> <operation> Failed!!! <error source>
const DefaultErrorMsg = "Resource %s %s Failed!!! %s"
const RequestIdMsg = "RequestId: %s"
const NotFoundMsg = ResourceNotfound + "!!! %s"
const WaitTimeoutMsg = "Resource %s %s Timeout In %d Seconds. Got: %s Expected: %s !!! %s"
const DefaultDebugMsg = "\n*************** %s Response *************** \n%s\n%s******************************\n\n"
const DataDefaultErrorMsg = "Datasource %s %s Failed!!! %s"

// default region for all resource
const DEFAULT_REGION = "BEIJING"
const DEFAULT_ENDPOINT = "ks3-cn-beijing.ksyuncs.com"

type TagResourceType string

type Invoker struct {
	catchers []*Catcher
}

type Catcher struct {
	Reason           string
	RetryCount       int
	RetryWaitSeconds int
}

var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 10, 5}

func (a *Invoker) AddCatcher(catcher Catcher) {
	a.catchers = append(a.catchers, &catcher)
}

func (a *Invoker) Run(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for _, catcher := range a.catchers {
		if IsExpectedErrors(err, []string{catcher.Reason}) {
			catcher.RetryCount--

			if catcher.RetryCount <= 0 {
				return fmt.Errorf("Retry timeout and got an error: %#v.", err)
			} else {
				time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
				return a.Run(f)
			}
		}
	}
	return err
}

func debugOn() bool {
	for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
		if strings.TrimSpace(part) == "terraform" {
			return true
		}
	}
	return false
}

func addDebug(action, content interface{}, requestInfo ...interface{}) {
	if debugOn() {
		trace := "[DEBUG TRACE]:\n"
		for skip := 1; skip < 5; skip++ {
			_, filepath, line, _ := runtime.Caller(skip)
			trace += fmt.Sprintf("%s:%d\n", filepath, line)
		}

		if len(requestInfo) > 0 {
			var request = struct {
				Domain     string
				Version    string
				UserAgent  string
				ActionName string
				Method     string
				Product    string
				Region     string
				AK         string
			}{}
			switch requestInfo[0].(type) {
			case *ks3.Client:
				request.Product = "KS3"
				request.ActionName = fmt.Sprintf("%s", action)
			}
			requestContent := ""
			if len(requestInfo) > 1 {
				requestContent = fmt.Sprintf("%#v", requestInfo[1])
			}

			if len(requestInfo) == 1 {
				if v, ok := requestInfo[0].(map[string]interface{}); ok {
					if res, err := json.Marshal(&v); err == nil {
						requestContent = string(res)
					}
					if res, err := json.Marshal(&content); err == nil {
						content = string(res)
					}
				}
			}

			content = fmt.Sprintf("%vDomain:%v, Version:%v, ActionName:%v, Method:%v, Product:%v, Region:%v\n\n"+
				"*************** %s Request ***************\n%#v\n",
				content, request.Domain, request.Version, request.ActionName,
				request.Method, request.Product, request.Region, request.ActionName, requestContent)
		}

		log.Printf(DefaultDebugMsg, action, content, trace)
	}
}

// Return a ComplexError which including extra error message, error occurred file and path
func GetFunc(level int) string {
	pc, _, _, ok := runtime.Caller(level)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in GetFuncName.")
		return ""
	}
	return strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
}
