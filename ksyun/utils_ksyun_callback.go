package ksyun

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ApiCall struct {
	param         *map[string]interface{}
	action        string
	beforeCall    beforeCallFunc
	executeCall   executeCallFunc
	callError     callErrorFunc
	afterCall     afterCallFunc
	disableDryRun bool
	process       int // process represent the ApiCall's process
}

type ksyunApiCallFunc func(d *schema.ResourceData, meta interface{}) error
type callErrorFunc func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error
type executeCallFunc func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (*map[string]interface{}, error)
type afterCallFunc func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) error
type beforeCallFunc func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (bool, error)

type ApiProcess struct {
	DryRun bool
	Ctx    context.Context
	MulNum int

	d      *schema.ResourceData
	client *KsyunClient

	apiProcessQueue []ApiCall
	wg              *sync.WaitGroup
	errCh           chan error
	concurrentCh    chan struct{}
	earlyStop       chan struct{}
	errors          []error
}

func ksyunApiCall(api []ksyunApiCallFunc, d *schema.ResourceData, meta interface{}) (err error) {
	if api != nil {
		for _, f := range api {
			if f != nil {
				err = f(d, meta)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

func ksyunApiCallNew(api []ApiCall, d *schema.ResourceData, client *KsyunClient, isDryRun bool) (err error) {
	if !client.dryRun || !isDryRun {
		return ksyunApiCallProcess(api, d, client, false)
	} else {
		err = ksyunApiCallProcess(api, d, client, true)
		if err != nil {
			return err
		}
		return ksyunApiCallProcess(api, d, client, false)
	}
}

func ksyunApiCallProcess(api []ApiCall, d *schema.ResourceData, client *KsyunClient, isDryRun bool) (err error) {
	if api != nil {
		for _, f := range api {
			if f.executeCall != nil {
				var (
					resp *map[string]interface{}
				)
				doExecute := true
				if isDryRun {
					if f.disableDryRun {
						continue
					}
					(*(f.param))["DryRun"] = true
				} else if f.beforeCall != nil {
					doExecute, err = f.beforeCall(d, client, f)
					if err == nil {
						f.process++
					}
				}
				if doExecute || isDryRun {
					resp, err = f.executeCall(d, client, f)
					if err == nil {
						f.process++
					}
				}
				if isDryRun {
					delete(*(f.param), "DryRun")
					if ksyunError, ok := err.(awserr.RequestFailure); ok && ksyunError.StatusCode() == 412 {
						err = nil
					}
				} else {
					if err != nil {
						if f.callError == nil {
							return err
						} else {
							err = f.callError(d, client, f, err)
						}
					}
					if err != nil {
						return err
					}
					if doExecute && f.afterCall != nil {
						err = f.afterCall(d, client, resp, f)
						if err == nil {
							f.process++
						}
					}
				}
			}
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (c *ApiCall) RightNow(d *schema.ResourceData, client *KsyunClient, isDryRun bool) error {
	return ksyunApiCallNew([]ApiCall{*c}, d, client, isDryRun)
}

func (a *ApiProcess) PutCalls(candidate ...ApiCall) {
	a.apiProcessQueue = append(a.apiProcessQueue, candidate...)
}

func (a *ApiProcess) SetD(d *schema.ResourceData) {
	a.d = d
}
func (a *ApiProcess) GetD() *schema.ResourceData {
	return a.d
}

func (a *ApiProcess) Client() *KsyunClient {
	return a.client
}

func (a *ApiProcess) SetClient(client *KsyunClient) {
	a.client = client
}

// NewApiProcess returns a ApiProcess, but now it doesn't support concurrency
func NewApiProcess(ctx context.Context, d *schema.ResourceData, client *KsyunClient, dryRun bool) ApiProcess {
	mulNum := 1

	p := ApiProcess{
		apiProcessQueue: []ApiCall{},
		d:               d,
		client:          client,
		DryRun:          dryRun,
		Ctx:             ctx,
	}
	p.MulNum = mulNum
	p.errCh = make(chan error, p.MulNum)
	p.concurrentCh = make(chan struct{}, p.MulNum)
	p.wg = &sync.WaitGroup{}
	return p
}

// ConRun will process ApiProcess concurrently
func (a *ApiProcess) ConRun() []error {
	// receive errs
	go func() {
		for callErr := range a.errCh {
			a.errors = append(a.errors, callErr)
		}
	}()

	// concurrency api call run
	for _, call := range a.apiProcessQueue {
		select {
		case <-a.Ctx.Done():
			a.errors = append(a.errors, fmt.Errorf("stop api call early"))
			return a.errors
		case <-a.earlyStop:
			return a.errors
		default:
		}

		a.wg.Add(1)
		a.concurrentCh <- struct{}{}
		go func(call ApiCall) {
			defer func() {
				a.wg.Done()
				<-a.concurrentCh
			}()
			callErr := call.RightNow(a.d, a.client, a.DryRun)
			a.errCh <- callErr
		}(call)
	}
	a.wg.Wait()
	return a.errors
}

func (a *ApiProcess) Run() error {
	defer a.Clean()

	return ksyunApiCallNew(a.apiProcessQueue, a.d, a.client, a.DryRun)
}
func (a *ApiProcess) Clean() {
	a.apiProcessQueue = make([]ApiCall, 0, 5)
}

func (c *ApiCall) Copy() ApiCall {
	dst := ApiCall{
		param:         c.param,
		executeCall:   c.executeCall,
		afterCall:     c.afterCall,
		beforeCall:    c.beforeCall,
		action:        c.action,
		callError:     c.callError,
		disableDryRun: c.disableDryRun,
	}

	return dst
}
