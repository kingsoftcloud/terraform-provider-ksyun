package ksyun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiCall_Retry(t *testing.T) {
	type fields struct {
		param          *map[string]interface{}
		action         string
		beforeCall     beforeCallFunc
		executeCall    executeCallFunc
		callError      callErrorFunc
		afterCall      afterCallFunc
		disableDryRun  bool
		process        int
		isRetry        bool
		client         *KsyunClient
		resp           *map[string]interface{}
		retryCondition []string
		retryCount     int
		Err            error
	}
	type args struct {
		d      *schema.ResourceData
		client *KsyunClient
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ApiCall{
				param:          tt.fields.param,
				action:         tt.fields.action,
				beforeCall:     tt.fields.beforeCall,
				executeCall:    tt.fields.executeCall,
				callError:      tt.fields.callError,
				afterCall:      tt.fields.afterCall,
				disableDryRun:  tt.fields.disableDryRun,
				process:        tt.fields.process,
				isRetry:        tt.fields.isRetry,
				client:         tt.fields.client,
				resp:           tt.fields.resp,
				retryCondition: tt.fields.retryCondition,
				retryCount:     tt.fields.retryCount,
				Err:            tt.fields.Err,
			}
			tt.wantErr(t, c.Retry(tt.args.d, tt.args.client), fmt.Sprintf("Retry(%v, %v)", tt.args.d, tt.args.client))
		})
	}
}
