package ksyun

import (
	"github.com/wilac-pv/ksyun-ks3-go-sdk/ks3"
)

func GetBucketInfo(client *ks3.Client, bucket string) (ks3.GetBucketInfoResult, error) {

	resp, err := client.ListBuckets()
	if err == nil {
		for _, bucketInfo := range resp.Buckets {
			if bucketInfo.Name == bucket {
				return ks3.GetBucketInfoResult{
					BucketInfo: ks3.BucketInfo{
						XMLName:      bucketInfo.XMLName,
						Name:         bucket,
						Region:       bucketInfo.Region,
						StorageClass: bucketInfo.Type,
					}}, nil
			}
		}
	}
	return ks3.GetBucketInfoResult{}, err
}
