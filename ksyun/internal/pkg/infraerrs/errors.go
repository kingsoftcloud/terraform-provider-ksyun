package infraerrs

import (
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

const (
	NetworkOpErrorMessage = "您的网络似乎不太稳定，请确认网络正常后重试"
)

// IsExpectError returns whether error is expected error
func IsExpectError(err error, expectError []string) bool {
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

func GetKsyunNetworkOpErrorMessage(origMessage string) string {
	return NetworkOpErrorMessage + ", Occur at: " + origMessage
}
