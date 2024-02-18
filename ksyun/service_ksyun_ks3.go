package ksyun

import (
	"github.com/ks3sdklib/ksyun-ks3-go-sdk/ks3"
	"strconv"
	"time"
)

// Ks3Service *connectivity.KsyunClient
type Ks3Service struct {
	client *KsyunClient
}

func (s *Ks3Service) DescribeKs3Bucket(id string) (response ks3.GetBucketInfoResult, err error) {
	request := map[string]string{"bucketName": id}
	var requestInfo *ks3.Client
	raw, err := s.client.WithKs3Client(func(ks3Client *ks3.Client) (interface{}, error) {
		requestInfo = ks3Client
		return GetBucketInfo(ks3Client, request["bucketName"])
	})
	if err != nil {
		if ks3NotFoundError(err) {
			return response, WrapErrorf(err, NotFoundMsg, KsyunKs3GoSdk)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetBucketInfo", KsyunKs3GoSdk)
	}

	addDebug("GetBucketInfo", raw, requestInfo, request)
	response, _ = raw.(ks3.GetBucketInfoResult)
	return
}

func (s *Ks3Service) WaitForKs3BucketObject(bucket *ks3.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		exist, err := bucket.IsObjectExist(id)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, "IsObjectExist", KsyunKs3GoSdk)
		}
		addDebug("IsObjectExist", exist)

		if !exist {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(exist), status, ProviderERROR)
		}
	}
}
