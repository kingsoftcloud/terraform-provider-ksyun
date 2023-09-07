package ksyun

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/pkg/errors"
)

const (
	NotFound         = "Notfound"
	ResourceNotfound = "ResourceNotfound"
	InternalError    = "INTERNAL_FAILURE"
	ServiceTimeout   = "ServiceTimeout"
)

// retryableErrorCode is retryable error code
var retryableErrorCode = []string{ServiceTimeout}

type ProviderError struct {
	errorCode string
	message   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Terraform Ksyun Provider Error: Code: %s Message: %s", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

/*
	func newNotFoundError(str string) error {
		return &ProviderError{
			errorCode: NotFound,
			message:   str,
		}
	}

	func getNotFoundMessage(product, id string) string {
		return fmt.Sprintf("the specified %s %s is not found", product, id)
	}

	func isNotFoundError(err error) bool {
		if e, ok := err.(*ProviderError); ok &&
			(e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), NotFound)) {
			return true
		}

		return false
	}
*/
func notFoundError(err error) bool {
	if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 404 {
		return true
	}
	errMessage := strings.ToLower(err.Error())
	if strings.Contains(errMessage, "notfound") ||
		strings.Contains(errMessage, "not found") ||
		strings.Contains(errMessage, "not exist") ||
		strings.Contains(errMessage, "not associate") ||
		strings.Contains(errMessage, "invalid") ||
		strings.Contains(errMessage, "not_found") {
		// strings.Contains(errMessage,"notfound"){
		return true
	}
	return false
}

func notFoundErrorNew(err error) bool {
	if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 404 {
		return true
	}
	errMessage := strings.ToLower(err.Error())
	if strings.Contains(errMessage, "notfound") ||
		strings.Contains(errMessage, "not_found") {
		// strings.Contains(errMessage,"notfound"){
		return true
	}
	return false
}

func inUseError(err error) bool {
	if err == nil {
		return false
	}
	errMessage := strings.ToLower(err.Error())
	if strings.Contains(errMessage, "inuse") ||
		strings.Contains(errMessage, "in use") ||
		strings.Contains(errMessage, "invalid_action") ||
		strings.Contains(errMessage, "used") {
		return true
	}
	return false
}

func isServerError(err error) bool {
	if err == nil {
		return false
	}
	if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() >= 500 {
		errMessage := strings.ToLower(ksyunError.Code())
		if strings.Contains(errMessage, InternalError) {
			return true
		}
	}
	return false
}

// retryError check whether retry and returns retry error
func retryError(err error, additionRetryableError ...string) *resource.RetryError {
	switch realErr := errors.Cause(err).(type) {
	case awserr.RequestFailure:
		if isExpectError(realErr, retryableErrorCode) {
			return resource.RetryableError(err)
		}

		if len(additionRetryableError) > 0 {
			if isExpectError(realErr, additionRetryableError) {
				return resource.RetryableError(err)
			}
		}
	default:
	}

	return resource.NonRetryableError(err)
}

// isExpectError returns whether error is expected error
func isExpectError(err error, expectError []string) bool {
	e, ok := err.(awserr.RequestFailure)
	if !ok {
		return false
	}

	longCode := e.Code()
	if IsContains(expectError, longCode) {
		return true
	}

	if strings.Contains(longCode, ".") {
		shortCodeSlice := strings.Split(longCode, ".")
		for _, shortCode := range shortCodeSlice {
			if IsContains(expectError, shortCode) {
				return true
			}
		}

	}

	return false
}

// IsContains returns whether value is within object
func IsContains(obj interface{}, value interface{}) bool {
	vv := reflect.ValueOf(obj)
	if vv.Kind() == reflect.Ptr || vv.Kind() == reflect.Interface {
		if vv.IsNil() {
			return false
		}
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Slice:
		for i := 0; i < vv.Len(); i++ {
			if reflect.DeepEqual(value, vv.Index(i).Interface()) {
				return true
			}
		}
		return false
	case reflect.Map:
		s := vv.MapKeys()
		for i := 0; i < len(s); i++ {
			if reflect.DeepEqual(value, s[i].Interface()) {
				return true
			}
		}
		return false
	case reflect.String:
		ss := reflect.ValueOf(value)
		switch ss.Kind() {
		case reflect.String:
			return strings.Contains(vv.String(), ss.String())
		}
		return false
	default:
		return reflect.DeepEqual(obj, value)
	}
}
